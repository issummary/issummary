package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/issummary/issummary/issummary"
	"github.com/issummary/issummary/usecase"
	"github.com/pkg/errors"
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
	workUseCase := usecase.NewWorkUseCase(client, config)
	worksBodyFunc := GetWorksBodyFunc(ctx, workUseCase)
	return CreateJsonHandleFunc(worksBodyFunc)
}

func GetWorksBodyFunc(ctx context.Context, workUseCase *usecase.WorkUseCase) func(body []byte) (interface{}, error) {
	worksBodyFunc := func(body []byte) (interface{}, error) {
		sortedWorks, err := workUseCase.GetSortedWorks(ctx)
		return ToWorks(sortedWorks), errors.Wrap(err, fmt.Sprintf("failed to get sorted works. config: %#v\n", workUseCase.Config))
	}

	return worksBodyFunc
}

func GetMilestonesJsonHandleFunc(ctx context.Context, client *issummary.Client, config *issummary.Config) http.HandlerFunc {
	milestoneBodyFunc := GetMilestonesBodyFunc(ctx, client, config)
	return CreateJsonHandleFunc(milestoneBodyFunc)
}

func GetMilestonesBodyFunc(ctx context.Context, client *issummary.Client, config *issummary.Config) func(body []byte) (interface{}, error) {
	milestonesBodyFunc := func(body []byte) (interface{}, error) { // FIXME Implement MilestoneUseCase
		var allMilestones []*issummary.Milestone
		for _, org := range config.Organizations {
			milestones, err := client.ListGroupMilestones(ctx, org)

			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch milestones which organization is %v and config is %#v\n", org, config))
			}

			allMilestones = append(allMilestones, milestones...)
		}

		return ToMilestones(allMilestones), nil
	}
	return milestonesBodyFunc
}
