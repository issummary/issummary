package usecase

import (
	"github.com/issummary/issummary/issummary"
	"github.com/mpppk/gitany"
)

func findProjectByID(projects []*issummary.Repository, id int64) (gitany.Repository, bool) {
	for _, project := range projects {
		if project.GetID() == id {
			return project, true
		}
	}

	return nil, false
}
