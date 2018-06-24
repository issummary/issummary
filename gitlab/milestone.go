package gitlab

import (
	"time"

	"github.com/mpppk/gitany"
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

func toMilestoneFromGroupMilestone(milestone *gitlab.GroupMilestone) *Milestone {
	return &Milestone{
		ID:        milestone.ID,
		IID:       milestone.IID,
		Title:     milestone.Title,
		StartDate: time.Time(*milestone.StartDate),
		DueDate:   time.Time(*milestone.DueDate),
		State:     milestone.State,
	}
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

func toMilestones(gitlabMilestones []*gitlab.GroupMilestone) (milestones []*Milestone) {
	for _, gm := range gitlabMilestones {
		milestones = append(milestones, toMilestoneFromGroupMilestone(gm))
	}
	return
}
