package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/xanzy/go-gitlab"
)

type Client struct {
	git          *gitlab.Client
	pid          interface{}
	targetLabels []string
}

func NewClient(token string, pid interface{}, targetLabels []string) *Client {
	git := gitlab.NewClient(nil, token)
	return &Client{
		git:          git,
		pid:          pid,
		targetLabels: targetLabels,
	}
}

func (c *Client) GetIssues() ([]*gitlab.Issue, *gitlab.Response, error) {
	opt := &gitlab.ListProjectIssuesOptions{
		Labels: c.targetLabels,
	}
	return c.git.Issues.ListProjectIssues(c.pid, opt)
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

func CreateJsonHandleFunc(bodyFunc BodyFunc) http.HandlerFunc {
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
