package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/mpppk/issummary/gitlab"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println(os.Getenv("GITLAB_TOKEN"))

	client := gitlab.New(os.Getenv("GITLAB_TOKEN"))

	if os.Getenv("GITLAB_BASEURL") != "" {
		client.SetBaseURL(os.Getenv("GITLAB_BASEURL"))
	}

	works, err := client.ListWorks(os.Getenv("GITLAB_PID"), &gitlab.Classes{"LC", "MC", "SC"})
	if err != nil {
		panic(err)
	}

	worksBodyFunc := func(body []byte) (interface{}, error) {
		fmt.Println("")
		return works, nil
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/works", createJsonHandleFunc(worksBodyFunc))
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
