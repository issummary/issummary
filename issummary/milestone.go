package issummary

import (
	"time"

	"github.com/mpppk/gitany"
)

type Milestone struct {
	ID        int
	IID       int
	Title     string
	StartDate time.Time
	DueDate   time.Time
	State     string
}

func toMilestone(milestone gitany.Milestone) *Milestone {
	if milestone == nil {
		return nil
	}

	return &Milestone{
		ID:        int(milestone.GetID()),
		IID:       milestone.GetNumber(),
		Title:     milestone.GetTitle(),
		StartDate: time.Time(*milestone.GetStartDate()),
		DueDate:   time.Time(*milestone.GetDueDate()),
		State:     milestone.GetState(),
	}
}

func toMilestones(gitlabMilestones []gitany.Milestone) (milestones []*Milestone) {
	for _, gm := range gitlabMilestones {
		milestones = append(milestones, toMilestone(gm))
	}
	return
}
