package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/mpppk/issummary/api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := api.NewClient(os.Getenv("GITLAB_TOKEN"), os.Getenv("GITLAB_PID"), []string{os.Getenv("GITLAB_TARGET_LABEL")})
	issues, _, err := client.GetIssues()
	if err != nil {
		panic(err)
	}

	for _, issue := range issues {
		fmt.Printf("%#v", *issue)
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/json", createJsonHandleFunc(echoBodyFunc))
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

func echoBodyFunc(body []byte) (interface{}, error) {
	input := Input{}
	if err := json.Unmarshal(body, &input); err != nil {
		return nil, err
	}
	return input, nil
}

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
