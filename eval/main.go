package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed problems/*
var problems embed.FS

func eval(problemBytes, poseBytes []byte) string {
	var problem Problem
	if err := json.Unmarshal(problemBytes, &problem); err != nil {
		log.Fatal(err)
	}

	var pose Pose
	if err := json.Unmarshal(poseBytes, &pose); err != nil {
		log.Fatal(err)
	}
	result := dislike(&problem, &pose)
	return result.String()
}

func cli(){
	if len(os.Args) != 3 {
		log.Fatal("./eval <problem file> <pose file>")
	}
	problemBytes, err  := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	poseBytes, err  := ioutil.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(eval(problemBytes, poseBytes))
}

func main(){
/*
	inputProblem := "{\"hole\":[[45,80],[35,95],[5,95],[35,50],[5,5],[35,5],[95,95],[65,95],[55," +
		"80]],\"epsilon\":150000,\"figure\":{\"edges\":[[2,5],[5,4],[4,1],[1,0],[0,8],[8,3],[3,7],[7,11],[11,13],[13,12],[12,18],[18,19],[19,14],[14,15],[15,17],[17,16],[16,10],[10,6],[6,2],[8,12],[7,9],[9,3],[8,9],[9,12],[13,9],[9,11],[4,8],[12,14],[5,10],[10,15]],\"vertices\":[[20,30],[20,40],[30,95],[40,15],[40,35],[40,65],[40,95],[45,5],[45,25],[50,15],[50,70],[55,5],[55,25],[60,15],[60,35],[60,65],[60,95],[70,95],[80,30],[80,40]]}}"
	poseText := "{\n\"vertices\": [\n[21, 28], [31, 28], [31, 87], [29, 41], [44, 43], [58, 70],\n[38, 79], [32, 31], [36, 50], [39, 40], [66, 77], [42, 29],\n[46, 49], [49, 38], [39, 57], [69, 66], [41, 70], [39, 60],\n[42, 25], [40, 35]\n]\n}"
 */
	var server = flag.String("server", "", "http server mode")
	flag.Parse()
	if *server != "" {
		fmt.Println("server mode")
		http.HandleFunc("/health", func(w http.ResponseWriter, request *http.Request) {
			w.WriteHeader(200)
			_, _ = w.Write([]byte("ok"))
		})
		http.HandleFunc("/eval/", func(w http.ResponseWriter, request *http.Request) {
			id := strings.TrimPrefix(request.URL.Path, "/eval/")
			p, err := problems.ReadFile("problems/" + id)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(string(p))
			body, err := ioutil.ReadAll(request.Body)
			if err != nil {
				log.Println(err)
				w.WriteHeader(500)
				_, _ = w.Write([]byte(err.Error()))
				return
			}
			log.Println(string(body))

			result := eval(p, body)
			_, _ = w.Write([]byte(fmt.Sprintf("{\"dislike\": %s}", result)))
			w.WriteHeader(200)
		})
		http.ListenAndServe(*server, nil)
	} else {
		cli()
	}

}