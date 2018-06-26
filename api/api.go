package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

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
