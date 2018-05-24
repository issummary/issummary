package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"strings"

	"github.com/joho/godotenv"
	"github.com/mpppk/issummary/gitlab"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := gitlab.New(os.Getenv("GITLAB_TOKEN"))

	if os.Getenv("GITLAB_BASEURL") != "" {
		client.SetBaseURL(os.Getenv("GITLAB_BASEURL"))
	}

	gidList := strings.Split(os.Getenv("GITLAB_PID"), ",")

	worksBodyFunc := func(body []byte) (interface{}, error) {
		workManager := gitlab.NewWorkManager()
		for _, gid := range gidList {
			works, err := client.ListGroupWorks(gid, "LC", "S")

			if err != nil {
				panic(err)
			}

			workManager.AddWorks(works)
		}

		workManager.ConnectByDependencies()
		sortedWorks, err := workManager.GetSortedWorks()

		if err != nil {
			panic(err)
		}
		return sortedWorks, nil
	}

	milestonesBodyFunc := func(body []byte) (interface{}, error) {
		var allMilestones []*gitlab.Milestone
		for _, gid := range gidList {
			milestones, err := client.ListGroupMilestones(gid)

			if err != nil {
				panic(err)
			}

			allMilestones = append(allMilestones, milestones...)
		}

		return allMilestones, nil
	}

	fs := http.FileServer(http.Dir("static/dist"))
	http.Handle("/", fs)
	http.HandleFunc("/api/works", createJsonHandleFunc(worksBodyFunc))
	http.HandleFunc("/api/milestones", createJsonHandleFunc(milestonesBodyFunc))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

type Input struct {
	In string
}

type Output struct {
	Out string
}

type ErrorOutput struct {
	Error error
}

type BodyFunc func(body []byte) (interface{}, error)

func createJsonHandleFunc(bodyFunc BodyFunc) http.HandlerFunc {
	jsonHandleFunc := func(rw http.ResponseWriter, req *http.Request) {
		var retJson interface{}
		defer func() {
			marshaledJson, err := json.Marshal(retJson)
			if err != nil {
				log.Fatal(err)
			}
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, string(marshaledJson))
		}()

		if req.Method != "POST" {
			retJson = ErrorOutput{errors.New("request is not post method")}
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			retJson = ErrorOutput{err}
			return
		}

		input, err := bodyFunc(body)
		if err != nil {
			retJson = ErrorOutput{err}
			return
		}
		retJson = input
	}
	return jsonHandleFunc
}
