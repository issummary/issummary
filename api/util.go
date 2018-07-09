package api

import (
	"github.com/issummary/issummary/issummary"
)

type Issue struct {
	ID          int64
	IID         int
	DueDate     string
	Title       string
	Description *IssueDescription
	URL         string
	ProjectName string
	Milestone   *Milestone
}

type IssueDescription struct {
	Raw     string
	Summary string
	Note    string
	Details string
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
	Relation        *issummary.WorkRelation
	Issue           *Issue
	Label           *Label
	DependWorks     []*Work
	TotalStoryPoint int
	StoryPoint      int
}

type Label struct {
	ID          int64
	Name        string
	ParentNames []string
	Description *issummary.LabelDescription
}

func toWork(work *issummary.Work) *Work {
	if work == nil {
		return nil
	}

	return &Work{
		Relation:        work.Relation,
		Issue:           toIssue(work.Issue),
		Label:           toLabel(work.Label),
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
		Description: toIssueDescription(issue.Description),
		URL:         issue.GetHTMLURL(),
		ProjectName: issue.ProjectName,
		Milestone:   toMilestone(issue.GetMilestone()),
	}
}

func toIssueDescription(description *issummary.IssueDescription) *IssueDescription {
	return &IssueDescription{
		Raw:     description.Raw,
		Summary: description.Summary,
		Note:    description.Note,
		Details: description.Details,
	}
}

func toLabel(label *issummary.Label) *Label {
	if label == nil {
		return nil
	}

	parentNames := []string{}
	for _, parentLabel := range label.ParentLabels {
		parentNames = append(parentNames, parentLabel.GetName())
	}

	return &Label{
		ID:          label.GetID(),
		Name:        label.GetName(),
		ParentNames: parentNames,
		Description: label.Description,
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
