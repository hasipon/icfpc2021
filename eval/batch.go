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
var latestDislike = map[string]*Int{}

func batchMode(solutionsDir string, submit bool) {
	if defaultDB.Ok() {
		go batchEvalDB()
	}
	if submit {
		go batchSubmission()
	}
	batchEvalDir(solutionsDir)
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
				dislike := new(Int)
				_, ok := dislike.SetString(p[1], 10)
				if ok {
					if latestDislike[p[0]] == nil || dislike.Cmp(latestDislike[p[0]]) != 1 { // dislike <= latest
						latestDislike[p[0]] = dislike
					} else {
						log.Println("Skipping update dislike: (problem old new) = ", p[0], latestDislike[p[0]], p[1])
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

func submitSolutionFile(problem, filePath string) (*SubmitResponse, error) {
	fileBody, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("os.Open err:", err)
		return nil, err
	}
	return submitSolution(problem, fileBody)
}

func submitSolution(problem string, jsonBytes []byte) (*SubmitResponse, error) {
	url := contestUrl + "/api/problems/" + problem + "/solutions"
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

	rateLimitTime := map[string]time.Time{}

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

			solution, err := defaultDB.FindBestSolution(problemID)
			if err != nil {
				log.Println("FindBestSolution err", err, problemID)
				continue
			}

			if solution.UseBonus != "" {
				// TODO ボーナスつきの自動提出
				continue
			}

			dislike := new(Int)
			dislike.SetString(solution.DislikeS, 10)

			if latestDislike[problemID] == nil || dislike.Cmp(latestDislike[problemID]) == -1 { // dislike < latest
				log.Println("Submitting", problemID, solution.ID)
				log.Printf("%#v", solution)

				// TODO CONVERT BONUS PARAMETER
				result, err := submitSolution(problemID, []byte(solution.Json))
				if err != nil {
					log.Fatal("submission error")
				}

				if result.Error == "" {
					latestDislike[problemID] = dislike
				} else if result.isSubmissionRateLimit() {
					log.Println(result.Error)
					sec := result.rateLimitSecond()
					rateLimitTime[problemID] = time.Now().Add(time.Second * time.Duration(sec))
				} else {
					log.Fatal(result.Error)
				}
			}
		}

		time.Sleep(30 * time.Second)
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

		problemID := sp[0]
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
						_, err = defaultDB.RegisterSolution(entry.Name(), sp[0], ans)
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

		prob, err := getProblem(solution.ProblemID)
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
				obtainBonuses = append(obtainBonuses, GenBonusKey(fmt.Sprint(b.Problem), b.Bonus))
				break
			}
		}
	}

	return obtainBonuses
}
