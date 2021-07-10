package main

import (
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
const problemsUrl = "http://13.114.46.162:8800/problems.json"
const problemsFetchUrl = "http://13.114.46.162:8800/fetch_problems"

var lastFetchTime time.Time = time.Now()
var latestDislike = map[string]*Int{}


func fetchProblemsJson() {
	// コンテストポータルをクロールしてくるので制限きつめ
	if time.Minute < time.Since(lastFetchTime) {
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

func submitSolution(problem, filePath string) (*SubmitResponse, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("os.Open err:", err)
		return nil, err
	}
	defer file.Close()

	url := contestUrl + "/api/problems/" + problem + "/solutions"
	bearer := "Bearer " + os.Getenv("YOUR_API_TOKEN")

	req, err := http.NewRequest("POST", url, file)
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

func batchSubmission(solutionsDir string) {
	if os.Getenv("YOUR_API_TOKEN") == "" {
		log.Fatal("YOUR_API_TOKEN not set")
	}

	for {
		fetchProblemsJson()
		updateDislikes()

		dirEntries, err := os.ReadDir(solutionsDir)
		if err != nil {
			log.Println("ReadDir", err)
			time.Sleep(30 * time.Second)
			continue
		}

		ratelimit := false

		for _, entry := range dirEntries {
			if entry.IsDir() || strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			sp := strings.Split(entry.Name(), "-")
			if len(sp) == 0 {
				log.Println("Invalid Name:", entry.Name())
				continue
			}

			if !strings.HasPrefix(sp[len(sp)-1], "dislike") {
				log.Println("Invalid name", entry.Name())
				continue
			}

			dislike := new(Int)
			_, ok := dislike.SetString(strings.TrimPrefix(sp[len(sp)-1], "dislike"), 10)
			if !ok {
				log.Println("dislike parse failed", entry.Name())
				continue
			}

			if latestDislike[sp[0]] == nil || dislike.Cmp(latestDislike[sp[0]]) == -1 { // dislike < latest
				result, err := submitSolution(sp[0], path.Join(solutionsDir, entry.Name()))
				if err != nil {
					log.Fatal("submission error")
				}

				if result.Error == "" {
					latestDislike[sp[0]] = dislike

					go func() {
						prob := sp[0]
						id := result.ID
						if waitPoseValidation(prob, id) {
							log.Println("Submission accepted", prob, id)
						} else {
							log.Fatal("Invalid submission")
						}
					}()
				} else if result.isSubmissionRateLimit() {
					log.Println(result.Error)
					ratelimit = true
				} else {
					log.Fatal(result.Error)
				}
			}
		}

		if ratelimit {
			time.Sleep(60 * time.Second)
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}

func batchMode(solutionsDir string, submit bool) {
	if submit {
		go batchSubmission(solutionsDir)
	}

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
				}

				ans, err := ioutil.ReadFile(path.Join(queueDir, entry.Name()))
				if err != nil {
					log.Println("ReadFile failed:", err)
				}

				result, valid, msg := eval(prob, ans)
				if valid {
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
