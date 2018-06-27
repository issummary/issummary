package api

import (
	"github.com/issummary/issummary/issummary"
)

type Issue struct {
	ID          int64
	IID         int
	DueDate     string
	Title       string
	Description string
	URL         string
	ProjectName string
	Milestone   *Milestone
}

type Milestone struct {
	ID        int
	IID       int
	parallel  int
	title     string
	StartDate string
	DueDate   string
}

type Work struct {
	Issue           *Issue
	Label           *Label
	Dependencies    *Dependencies
	DependWorks     []*Work
	TotalStoryPoint int
	StoryPoint      int
}

type Dependencies struct {
	Issues []*Issue
	Labels []*DependLabel
}

type Label struct {
	ID           int64
	Name         string
	Description  string
	Parent       *Label
	Dependencies []*DependLabel
}

type DependLabel struct {
	Label         *Label
	RelatedIssues []*Issue
}

func toWork(work *issummary.Work) *Work {
	if work == nil {
		return nil
	}

	return &Work{
		Issue:           toIssue(work.Issue),
		Label:           toLabel(work.Label),
		Dependencies:    toDependencies(work.Dependencies),
		DependWorks:     ToWorks(work.DependWorks),
		TotalStoryPoint: work.TotalStoryPoint,
		StoryPoint:      work.StoryPoint,
	}
}

func ToWorks(works []*issummary.Work) (apiWorks []*Work) {
	apiWorks = []*Work{}
	for _, work := range works {
		if work == nil {
			continue
		}

		apiWorks = append(apiWorks, toWork(work))
	}
	return
}

func toMilestone(milestone *issummary.Milestone) *Milestone {
	if milestone == nil {
		return nil
	}

	return &Milestone{
		ID: milestone.ID, // FIXME
	}
}

func ToMilestones(milestones []*issummary.Milestone) (apiMilestones []*Milestone) {
	apiMilestones = []*Milestone{}
	for _, milestone := range milestones {
		if milestone == nil {
			continue
		}
		apiMilestones = append(apiMilestones, toMilestone(milestone))
	}
	return
}

func toIssue(issue *issummary.Issue) *Issue {
	if issue == nil {
		return nil
	}

	dueDateString := ""
	if issue.GetDueDate() != nil {
		dueDateString = issue.GetDueDate().String()
	}

	return &Issue{
		ID:          issue.GetID(),
		IID:         issue.GetNumber(),
		DueDate:     dueDateString,
		Title:       issue.GetTitle(),
		Description: issue.GetBody(),
		URL:         issue.GetHTMLURL(),
		ProjectName: issue.ProjectName,
		Milestone:   toMilestone(issue.GetMilestone()),
	}
}

func toIssues(issues []*issummary.Issue) (apiIssues []*Issue) {
	apiIssues = []*Issue{}
	for _, issue := range issues {
		if issue == nil {
			continue
		}

		apiIssues = append(apiIssues, toIssue(issue))
	}
	return
}

func toLabel(label *issummary.Label) *Label {
	if label == nil {
		return nil
	}

	return &Label{
		ID:           label.GetID(),
		Name:         label.GetName(),
		Description:  label.GetDescription(),
		Parent:       toLabel(label.Parent),
		Dependencies: toDependLabels(label.Dependencies),
	}
}

func toLabels(labels []*issummary.Label) (apiLabels []*Label) {
	apiLabels = []*Label{}
	for _, label := range labels {
		if label == nil {
			continue
		}

		apiLabels = append(apiLabels, toLabel(label))
	}
	return
}

func toDependencies(dependencies *issummary.Dependencies) *Dependencies {
	if dependencies == nil {
		return nil
	}

	return &Dependencies{
		Issues: toIssues(dependencies.Issues),
		Labels: toDependLabels(dependencies.Labels),
	}
}

func toDependLabels(dependLabels []*issummary.DependLabel) (apiDependLabels []*DependLabel) {
	apiDependLabels = []*DependLabel{}
	for _, dependLabel := range dependLabels {
		if dependLabel == nil {
			continue
		}

		apiDependLabels = append(apiDependLabels, toDependLabel(dependLabel))
	}
	return
}

func toDependLabel(dependLabel *issummary.DependLabel) *DependLabel {
	if dependLabel == nil {
		return nil
	}

	return &DependLabel{
		Label:         toLabel(dependLabel.Label),
		RelatedIssues: toIssues(dependLabel.RelatedIssues),
	}
}
