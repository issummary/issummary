package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/issummary/issummary/issummary"
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

func GetWorksJsonHandleFunc(ctx context.Context, client *issummary.Client, config *issummary.Config) http.HandlerFunc {
	worksBodyFunc := GetWorksBodyFunc(ctx, client, config)
	return CreateJsonHandleFunc(worksBodyFunc)
}

func GetWorksBodyFunc(ctx context.Context, client *issummary.Client, config *issummary.Config) func(body []byte) (interface{}, error) {
	worksBodyFunc := func(body []byte) (interface{}, error) {
		workManager := issummary.NewWorkManager()
		for _, gid := range config.GIDs {
			if err := client.Fetch(ctx, gid); err != nil {
				return nil, err
			}
			works, err := client.ListGroupWorks(gid, config.ClassLabelPrefix, config.SPLabelPrefix)

			if err != nil {
				return nil, err
			}

			workManager.AddWorks(works)
			workManager.AddLabels(client.Labels)
		}

		if err := workManager.ResolveDependencies(); err != nil {
			return nil, err
		}
		sortedWorks, err := workManager.GetSortedWorks()
		if err != nil {
			return nil, err
		}

		return ToWorks(sortedWorks), nil
	}

	return worksBodyFunc
}

func GetMilestonesJsonHandleFunc(ctx context.Context, client *issummary.Client, config *issummary.Config) http.HandlerFunc {
	milestoneBodyFunc := GetMilestonesBodyFunc(ctx, client, config)
	return CreateJsonHandleFunc(milestoneBodyFunc)
}

func GetMilestonesBodyFunc(ctx context.Context, client *issummary.Client, config *issummary.Config) func(body []byte) (interface{}, error) {
	milestonesBodyFunc := func(body []byte) (interface{}, error) {
		var allMilestones []*issummary.Milestone
		for _, gid := range config.GIDs {
			milestones, err := client.ListGroupMilestones(ctx, gid)

			if err != nil {
				panic(err)
			}

			allMilestones = append(allMilestones, milestones...)
		}

		return ToMilestones(allMilestones), nil
	}
	return milestonesBodyFunc
}
