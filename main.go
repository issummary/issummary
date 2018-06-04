//go:generate statik -src=./static/dist

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mpppk/issummary/gitlab"
	_ "github.com/mpppk/issummary/statik"
	"github.com/rakyll/statik/fs"
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
				return nil, err
			}

			workManager.AddWorks(works)
		}

		if err := workManager.ConnectByDependencies(); err != nil {
			return nil, err
		}
		sortedWorks, err := workManager.GetSortedWorks()

		if err != nil {
			return nil, err
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

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/", http.FileServer(statikFS))

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
	Error string
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
			fmt.Println(string(marshaledJson))
			fmt.Fprint(rw, string(marshaledJson))
		}()

		if req.Method != "POST" {
			retJson = ErrorOutput{"request is not post method"}
			return
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			retJson = ErrorOutput{err.Error()}
			return
		}

		input, err := bodyFunc(body)
		if err != nil {
			retJson = ErrorOutput{err.Error()}
			return
		}
		retJson = input
	}
	return jsonHandleFunc
}
