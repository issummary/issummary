package gitlab

import (
	"time"

	"github.com/xanzy/go-gitlab"
)

type Milestone struct {
	ID        int
	IID       int
	Title     string
	StartDate time.Time
	DueDate   time.Time
	State     string
}

func toMilestone(milestone *gitlab.GroupMilestone) *Milestone {
	return &Milestone{
		ID:        milestone.ID,
		IID:       milestone.IID,
		Title:     milestone.Title,
		StartDate: time.Time(*milestone.StartDate),
		DueDate:   time.Time(*milestone.DueDate),
		State:     milestone.State,
	}
}

func toMilestones(gitlabMilestones []*gitlab.GroupMilestone) (milestones []*Milestone) {
	for _, gm := range gitlabMilestones {
		milestones = append(milestones, toMilestone(gm))
	}
	return
}
