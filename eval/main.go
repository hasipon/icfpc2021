package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

//go:embed problems/*
var problems embed.FS

func eval(problemBytes, poseBytes []byte) (string, bool, string) {
	var problem *Problem
	if err := json.Unmarshal(problemBytes, &problem); err != nil {
		log.Fatal("problem:", err)
	}

	var pose *Pose
	if err := json.Unmarshal(poseBytes, &pose); err != nil {
		log.Fatal("pose:", err)
	}
	problem = applyBonus(problem, pose)
	result := dislike(problem, pose)
	valid, msg := validate(problem, pose)
	return result.String(), valid, msg

}
func getProblem(id string) ([]byte, error) {
	p, err := problems.ReadFile("problems/" + id)
	return p, err
}

func convertAll() {
	var err error
	for i := 1; ; i++ {
		id := fmt.Sprintf("%d", i)
		_, err = getProblem(id)
		if err != nil {
			break
		}
		result := convert("", id)
		err = ioutil.WriteFile("./converted_problems/"+id, result, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func convert(problemFile, problemId string) []byte {
	var problemBytes []byte
	var err error
	if problemFile != "" {
		problemBytes, err = ioutil.ReadFile(problemFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if problemId != "" {
		problemBytes, err = getProblem(problemId)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("./eval --convert --problem-file <filename>\n" +
			"./eval --convert --problem-id <id>\n")
	}

	var problem Problem
	if err := json.Unmarshal(problemBytes, &problem); err != nil {
		log.Fatal(err)
	}

	n := len(problem.Hole)
	en := len(problem.Figure.Edges)
	vn := len(problem.Figure.Vertices)

	var ret bytes.Buffer
	ret.WriteString(fmt.Sprintf("%d %d %d %s\n", n, en, vn, problem.Epsilon.String()))

	for _, h := range problem.Hole {
		ret.WriteString(fmt.Sprintf("%s %s\n", h[0].String(), h[1].String()))
	}
	for _, e := range problem.Figure.Edges {
		ret.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	for _, v := range problem.Figure.Vertices {
		ret.WriteString(fmt.Sprintf("%s %s\n", v[0].String(), v[1].String()))
	}
	return ret.Bytes()
}

func cli(problemFile, problemId, poseFile string) {
	var problemBytes []byte
	var err error
	if problemFile != "" {
		problemBytes, err = ioutil.ReadFile(problemFile)
		if err != nil {
			log.Fatal(err)
		}
	} else if problemId != "" {
		problemBytes, err = getProblem(problemId)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("./eval --problem-file <filename> --pose-file <poseFile>\n" +
			"./eval --problem-id <id> --pose-file <poseFile>\n")
	}
	poseBytes, err := ioutil.ReadFile(poseFile)
	if err != nil {
		log.Fatal(err)
	}
	dislike, valid, msg := eval(problemBytes, poseBytes)
	if valid {
		fmt.Println(dislike)
		fmt.Println("valid")
	} else {
		fmt.Println(-1)
		fmt.Println(msg)
	}
}

func main() {
	/*
		inputProblem := "{\"hole\":[[45,80],[35,95],[5,95],[35,50],[5,5],[35,5],[95,95],[65,95],[55," +
			"80]],\"epsilon\":150000,\"figure\":{\"edges\":[[2,5],[5,4],[4,1],[1,0],[0,8],[8,3],[3,7],[7,11],[11,13],[13,12],[12,18],[18,19],[19,14],[14,15],[15,17],[17,16],[16,10],[10,6],[6,2],[8,12],[7,9],[9,3],[8,9],[9,12],[13,9],[9,11],[4,8],[12,14],[5,10],[10,15]],\"vertices\":[[20,30],[20,40],[30,95],[40,15],[40,35],[40,65],[40,95],[45,5],[45,25],[50,15],[50,70],[55,5],[55,25],[60,15],[60,35],[60,65],[60,95],[70,95],[80,30],[80,40]]}}"
		poseText := "{\n\"vertices\": [\n[21, 28], [31, 28], [31, 87], [29, 41], [44, 43], [58, 70],\n[38, 79], [32, 31], [36, 50], [39, 40], [66, 77], [42, 29],\n[46, 49], [49, 38], [39, 57], [69, 66], [41, 70], [39, 60],\n[42, 25], [40, 35]\n]\n}"
	*/
	var server = flag.String("server", "", "http server mode(specify listen addr)")
	var problemId = flag.String("problem-id", "", "problem id")
	var problemFile = flag.String("problem-file", "", "problem file")
	var poseFile = flag.String("pose-file", "", "pose file")
	var batch = flag.Bool("batch", false, "batch mode")
	var solutions = flag.String("solutions", "solutions", "solutions directory")
	var submit = flag.Bool("submit", false, "auto submit (batch mode)")
	var convertFlag = flag.Bool("convert", false, "input convert")
	var allFlag = flag.Bool("all", false, "convert all")
	var dbpath = flag.String("db", "", "sqlitedb path")
	var initdb = flag.Bool("initdb", false, "init db")
	var solutionsToDB = flag.Bool("solutions-to-db", false, "insert solution from solution folder")
	flag.Parse()
	if *dbpath != "" {
		conn, err := sqlx.Open("sqlite3", *dbpath)
		if err != nil {
			log.Fatal("failed to open database", err)
		}
		defaultDB = SQLiteDB{DB: conn}

		if *initdb {
			defaultDB.Init()
		}

		if *solutionsToDB {
			registerSolutionInDirectory(*solutions)
		}
	}

	if *server != "" {
		fmt.Println("server mode")
		if *batch {
			fmt.Println("also batch mode")
			go batchMode(*solutions, *submit)
		}

		http.HandleFunc("/health", func(w http.ResponseWriter, request *http.Request) {
			w.WriteHeader(200)
			_, _ = w.Write([]byte("ok"))
		})
		http.HandleFunc("/eval/", func(w http.ResponseWriter, request *http.Request) {
			id := strings.TrimPrefix(request.URL.Path, "/eval/")
			log.Printf("id = %s\n", id)
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			problem, err := getProblem(id)
			if err != nil {
				log.Print(err)
				w.WriteHeader(500)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			result, valid, msg := eval(problem, body)
			if !valid {
				result = "-1"
			}

			if *batch && valid {
				fileName := fmt.Sprintf("%v-eval-dislike%v", id, result)
				filePath := path.Join(*solutions, fileName)
				_, err = os.Stat(filePath)
				if os.IsNotExist(err) {
					os.WriteFile(filePath, body, 0644)
				}
			}

			w.WriteHeader(200)
			_, _ = w.Write([]byte(fmt.Sprintf("{\"dislike\": %s, \"valid\": %t, \"msg\": \"%s\"}",
				result, valid, msg)))
		})
		for {
			err := http.ListenAndServe(*server, nil)
			log.Print(err)
		}
	} else if *batch {
		batchMode(*solutions, *submit)
	} else if *allFlag {
		convertAll()
	} else if *convertFlag {
		fmt.Print(string(convert(*problemFile, *problemId)))
	} else {
		cli(*problemFile, *problemId, *poseFile)
	}

}
