package issummary

import "time"

type Label struct {
	ID           int
	Name         string
	Description  *LabelDescription
	Parent       *Label
	Dependencies []*DependLabel
}

type LabelDescription struct {
	Raw              string
	DependLabelNames []string
	ParentName       string // TODO: 複数の親を持てるようにする
}

type Issue struct {
	ID          int
	IID         int
	DueDate     *time.Time
	Title       string
	Description *IssueDescription
	URL         string
	ProjectName string
	GroupName   string
	Milestone   *Milestone
}

type IssueDescription struct {
	Raw          string
	Dependencies *IssueDependencies
	Summary      string
	Note         string
	Details      string
}

type IssueDependencies struct { // FIXME merge Dependencies
	Issues             []*DependIssue
	LabelNames         []string
	OtherProjectIssues []*DependIssue
}

type DependIssue struct {
	GroupName   string
	ProjectName string
	IID         int
	ID          int
}

type Work struct {
	Issue           *Issue
	Label           *Label
	Dependencies    *Dependencies // FIXME
	DependWorks     []*Work
	StoryPoint      int
	TotalStoryPoint int
	ManDay          int
	TotalManDay     int
	RemainManDays   int
}

func (w *Work) GetTotalStoryPoint() (totalSP int) {
	for _, dWork := range w.DependWorks {
		totalSP += dWork.StoryPoint
	}
	return
}

type DependLabel struct {
	Label         *Label
	RelatedIssues []*Issue
}

type Dependencies struct {
	OtherProjectIssues []*Issue
	Issues             []*Issue
	Labels             []*DependLabel
}
