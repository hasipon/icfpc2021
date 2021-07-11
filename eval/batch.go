package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const contestUrl = "https://poses.live"

var problemsUrl = "http://13.114.46.162:8800/problems.json"
var problemsFetchUrl = "http://13.114.46.162:8800/fetch_problems"

func init() {
	host, _ := os.Hostname()
	if host == "ip-172-30-1-195" {
		problemsUrl = "http://localhost:8800/problems.json"
		problemsFetchUrl = "http://localhost:8800/fetch_problems"
	}
}

var lastFetchTime = time.Time{}
var latestDislike = map[int]*Int{}

func batchMode(solutionsDir string, submit bool) {
	if defaultDB.Ok() {
		go batchEvalDB()
	}
	if submit {
		go batchSubmission()
	}
	batchEvalDir(solutionsDir)
}

var errNoBonusUsed = fmt.Errorf("no bonus used")
var errDifferentBonus = fmt.Errorf("different bonus used")

func replaceBonusProblemID(poseJson []byte, bonusName string, problemID int) ([]byte, error) {
	var pose Pose
	err := json.Unmarshal(poseJson, &pose)
	if err != nil {
		return nil, err
	}

	if len(pose.Bonuses) == 0 {
		return nil, errNoBonusUsed
	}

	if pose.Bonuses[0].Bonus != bonusName {
		return nil, errDifferentBonus
	}

	pose.Bonuses[0].Problem = problemID
	result, err := json.Marshal(pose)
	return result, err
}

func fetchProblemsJson() {
	if 30*time.Second < time.Since(lastFetchTime) {
		resp, err := http.Get(problemsFetchUrl)
		if err != nil {
			log.Println(err)
		} else {
			defer resp.Body.Close()
			lastFetchTime = time.Now()
			io.Copy(io.Discard, resp.Body)
		}
	}
}

func updateDislikes() {
	resp, err := http.Get(problemsUrl)
	if err != nil {
		log.Println(err)
		return
	} else {
		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)
		var problems [][]string
		err = dec.Decode(&problems)
		if err != nil {
			log.Println(err)
			return
		}

		for _, p := range problems {
			if 3 <= len(p) {
				probId, _ := strconv.Atoi(p[0])
				dislike := new(Int)
				_, ok := dislike.SetString(p[1], 10)
				if ok {
					if latestDislike[probId] == nil || dislike.Cmp(latestDislike[probId]) != 1 { // dislike <= latest
						latestDislike[probId] = dislike
					} else {
						log.Println("Skipping update dislike: (problem old new) = ", p[0], latestDislike[probId], p[1])
					}
				}
			}
		}
	}
}

type PoseResponse struct {
	State    string `json:"state"`
	Dislikes *Int   `json:"dislikes"`
	Error    string `json:"error"`
}

type SubmitResponse struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

func (s *SubmitResponse) isSubmissionRateLimit() bool {
	return strings.HasPrefix(s.Error, "Submission rate limit exceeded")
}

func (s *SubmitResponse) rateLimitSecond() int {
	value := 0
	fmt.Sscanf(s.Error, "Submission rate limit exceeded, please wait %d seconds before trying again", &value)
	return value
}

func waitPoseValidation(problem, poseID string) bool {
	time.Sleep(5 * time.Second)

	url := contestUrl + "/api/problems/" + problem + "/solutions/" + poseID
	bearer := "Bearer " + os.Getenv("YOUR_API_TOKEN")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("NewRequest err:", err)
		return false
	}
	req.Header.Add("Authorization", bearer)

	log.Println(url)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Get submission err:", err)
		return false
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Get submission read err:", err)
		return false
	}

	log.Println(string(body))

	result := new(PoseResponse)
	err = json.Unmarshal(body, result)
	if err != nil {
		log.Println("Get submission result err:", err)
		return false
	}

	log.Println(result)

	if result.State == "PENDING" {
		resp.Body.Close()
		return waitPoseValidation(problem, poseID)
	}

	return result.State == "VALID"
}

func submitSolutionFile(problemID int, filePath string) (*SubmitResponse, error) {
	fileBody, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("os.Open err:", err)
		return nil, err
	}
	return submitSolution(problemID, fileBody)
}

func submitSolution(problemID int, jsonBytes []byte) (*SubmitResponse, error) {
	url := contestUrl + "/api/problems/" + fmt.Sprint(problemID) + "/solutions"
	bearer := "Bearer " + os.Getenv("YOUR_API_TOKEN")

	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		log.Println("NewRequest err:", err)
		return nil, err
	}

	req.Header.Add("Authorization", bearer)

	log.Println(url)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("Submit err:", err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Read response err:", err)
		return nil, err
	}

	result := new(SubmitResponse)
	err = json.Unmarshal(body, result)
	if err != nil {
		log.Println("Get submission result err:", err)
		return nil, err
	}

	log.Println("Submit result", result)

	return result, nil
}

func batchSubmission() {
	if os.Getenv("YOUR_API_TOKEN") == "" {
		log.Fatal("YOUR_API_TOKEN not set")
	}

	rateLimitTime := map[int]time.Time{}

	for {
		fetchProblemsJson()
		updateDislikes()

		problemIds, err := defaultDB.GetAllProblemIDsInSubmission()
		if err != nil {
			log.Println("GetAllProblemIDsInSubmission err", err)
			time.Sleep(30 * time.Second)
			continue
		}

		for _, problemID := range problemIds {
			if 0 < time.Until(rateLimitTime[problemID]) {
				continue
			}

			_, err := defaultDB.GetProblemSetting(problemID)
			if err == sql.ErrNoRows {
				// Insert Empty Setting
				_ = defaultDB.InsertProblemSetting(&ProblemSetting{ProblemID: problemID})
			}

			solution, err := defaultDB.FindBestSolution(problemID)
			if err != nil {
				log.Println("FindBestSolution err", err, problemID)
				continue
			}

			poseBytes := []byte(solution.Json)

			// どの問題でこのボーナスをアンロックするべきか調べて
			// ボーナスをアンロックする問題番号の差し替え
			if solution.UseBonus != "" {
				setting, err := defaultDB.GetWhichProblemUnlocksTheBonus(GenBonusKey(solution.ProblemID, solution.UseBonus))
				if err != nil {
					log.Println("GetWhichProblemUnlocksTheBonus err", err)
					continue
				}

				poseBytes, err = replaceBonusProblemID(poseBytes, solution.UseBonus, setting.ProblemID)
				if err != nil {
					log.Println("replaceBonusProblemID err", err)
					continue
				}
			}

			dislike := new(Int)
			dislike.SetString(solution.DislikeS, 10)

			submission, err := defaultDB.GetSubmission(problemID)
			if err != nil {
				if err == sql.ErrNoRows {
				} else {
					log.Println("GetSubmission err", err)
					continue
				}
			}

			if submission != nil && submission.Json == string(poseBytes) {
				// 同じサブミットはしない
				continue
			}

			submission = &Submission{
				ProblemID:   problemID,
				Json:        string(poseBytes),
				Dislike:     solution.Dislike,
				DislikeS:    solution.DislikeS,
				UseBonus:    solution.UseBonus,
				UnlockBonus: solution.UnlockBonus,
			}

			log.Println("Submitting", problemID, solution.ID, solution.UseBonus, solution.UnlockBonus)

			result, err := submitSolution(problemID, poseBytes)
			if err != nil {
				log.Println("submission error")
				continue
			}

			if result.Error == "" {
				err = defaultDB.ReplaceSubmission(submission)
				if err != nil {
					log.Println("ReplaceSubmission error", err)
				}
			} else if result.isSubmissionRateLimit() {
				log.Println(result.Error)
				sec := result.rateLimitSecond()
				rateLimitTime[problemID] = time.Now().Add(time.Second * time.Duration(sec))
				log.Println("RateLimit updated. problem ", problemID, " duration", time.Until(rateLimitTime[problemID]))
			} else {
				log.Println("submission error", result.Error)
				continue
			}

			time.Sleep(time.Second)
		}

		time.Sleep(10 * time.Second)
	}
}

// solutionsディレクトリにあるファイルたちをDBに登録します
func registerSolutionInDirectory(solutionsDir string) {
	if !defaultDB.Ok() {
		return
	}

	dirEntries, err := os.ReadDir(solutionsDir)
	if err != nil {
		log.Println("ReadDir failed", err)
		return
	}

	for _, entry := range dirEntries {
		if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		sp := strings.Split(entry.Name(), "-")

		if !strings.HasPrefix(sp[len(sp)-1], "dislike") {
			log.Println("Invalid name", entry.Name())
			continue
		}

		problemID, _ := strconv.Atoi(sp[0])
		solutionName := strings.Join(sp[1:len(sp)], "-")
		ans, err := ioutil.ReadFile(path.Join(solutionsDir, entry.Name()))
		if err != nil {
			log.Println("Read Solution failed", err)
			continue
		}

		_, err = defaultDB.RegisterSolution(solutionName, problemID, ans)
		if err != nil {
			log.Println("Register Solution failed", err)
		}
	}
}

func batchEvalDir(solutionsDir string) {
	queueDir := path.Join(solutionsDir, "queue")
	for {
		dirEntries, err := os.ReadDir(queueDir)
		if err != nil {
			log.Println("ReadDir", err)
			time.Sleep(30 * time.Second)
			continue
		}

		for _, entry := range dirEntries {
			if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				log.Println("entryInfo:", err)
				continue
			}

			if 3*time.Second < time.Since(info.ModTime()) {
				sp := strings.SplitN(entry.Name(), "-", 2)
				if len(sp) == 0 {
					log.Println("Invalid Name:", entry.Name())
					continue
				}

				log.Println("Processing:", entry.Name())

				prob, err := getProblem(sp[0])
				if err != nil {
					log.Println("getProblem failed:", err)
					continue
				}

				ans, err := ioutil.ReadFile(path.Join(queueDir, entry.Name()))
				if err != nil {
					log.Println("ReadFile failed:", err)
					continue
				}

				result, valid, msg := eval(prob, ans)
				if valid {
					if defaultDB.Ok() {
						probId, _ := strconv.Atoi(sp[0])
						_, err = defaultDB.RegisterSolution(entry.Name(), probId, ans)
						if err != nil {
							log.Println("Register Solution failed", err)
						}
					}

					log.Println("Valid", entry.Name(), "Dislike", result)
					solutionFileName := fmt.Sprint(entry.Name(), "-dislike", result)
					log.Println("Moving", entry.Name(), " -> ", solutionFileName)
					err = os.Rename(path.Join(queueDir, entry.Name()), path.Join(solutionsDir, solutionFileName))
					if err != nil {
						log.Println("os.Rename:", err)
					}
				} else {
					log.Println("Invalid:", entry.Name(), msg)
				}
			}
		}
		time.Sleep(5 * time.Second)
	}
}

func batchEvalDB() {
	if !defaultDB.Ok() {
		return
	}
	for {
		solution, err := defaultDB.FindNoEvalSolution()
		if err == sql.ErrNoRows {
			time.Sleep(10 * time.Second)
			continue
		}
		if err != nil {
			log.Println("FindNoEvalSolution err", err)
			time.Sleep(30 * time.Second)
			continue
		}

		prob, err := getProblem(fmt.Sprint(solution.ProblemID))
		if err != nil {
			log.Println("getProblem failed:", err)
			time.Sleep(30 * time.Second)
			continue
		}

		poseBytes := []byte(solution.Json)
		result, valid, msg := eval(prob, poseBytes)
		log.Println("batchEvalDB", solution.ID, result, valid, msg)

		var bonusKeys []BonusKey
		if valid {
			bonusKeys = obtainBonusKeys(prob, poseBytes)
		}
		err = defaultDB.UpdateSolutionEvalResult(solution, result, valid, msg, bonusKeys)
		if err != nil {
			fmt.Println("UpdateSolutionEvalResult err", err)
		}
	}
}

func obtainBonusKeys(problemBytes []byte, poseBytes []byte) []BonusKey {
	var problem *Problem
	if err := json.Unmarshal(problemBytes, &problem); err != nil {
		log.Fatal("problem:", err)
	}

	var pose *Pose
	if err := json.Unmarshal(poseBytes, &pose); err != nil {
		log.Fatal("pose:", err)
	}

	var obtainBonuses []BonusKey
	for _, b := range problem.Bonuses {
		for _, v := range pose.Vertices {
			if v[0].Cmp(b.Position[0]) == 0 && v[1].Cmp(b.Position[1]) == 0 {
				obtainBonuses = append(obtainBonuses, GenBonusKey(b.Problem, b.Bonus))
				break
			}
		}
	}

	return obtainBonuses
}
