package gitlab

import (
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
)

type Client struct {
	*gitlab.Client
}

func New(token string) *Client {
	return &Client{
		Client: gitlab.NewClient(nil, token),
	}
}

type Label struct {
	ID           int
	Name         string
	Description  *LabelDescription
	Parent       *Label
	Dependencies []*Label
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
}

type IssueDescription struct {
	Raw            string
	DependencyIIDs []int
	Summary        string
	Note           string
	Details        string
}

type Work struct {
	Issue        *Issue
	Label        *Label
	Dependencies *Dependencies
	StoryPoint   int
}

type Dependencies struct {
	Issues []*Issue
	Labels []*Label
}

func (c *Client) ListWorks(pid interface{}, prefix, spLabelPrefix string) (works []*Work, err error) {
	allIssues, err := c.listAllIssuesByLabel(pid, gitlab.Labels{"W"}) // TODO: 外から指定できるようにする
	if err != nil {
		return nil, err
	}

	labels, err := c.listAllLabels(pid)

	// TODO: worksをオプティカルソート
	works, err = toWorks(allIssues, labels, prefix, spLabelPrefix)
	if err != nil {
		return nil, err
	}
	workManager := NewWorkManager()
	workManager.AddWorks(works)
	workManager.ConnectByDependencies()
	return workManager.GetSortedWorks()
}

func (c *Client) listLabelsByPrefix(pid interface{}, prefix string) (prefixLabels []*gitlab.Label, err error) {
	labels, err := c.listAllLabels(pid)
	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		if strings.Contains(label.Name, prefix) {
			prefixLabels = append(prefixLabels, label)
		}
	}
	return prefixLabels, nil
}

func (c *Client) listAllIssuesByLabel(pid interface{}, labels gitlab.Labels) ([]*gitlab.Issue, error) {
	issueOpt := &gitlab.ListProjectIssuesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
		Labels: labels,
	}

	var allIssues []*gitlab.Issue

	for {
		issues, _, err := c.Issues.ListProjectIssues(pid, issueOpt)
		if err != nil {
			return nil, err
		}

		if len(issues) == 0 {
			break
		}

		allIssues = append(allIssues, issues...)
		issueOpt.Page = issueOpt.Page + 1
	}
	return allIssues, nil
}

func (c *Client) listAllLabels(pid interface{}) ([]*gitlab.Label, error) {
	opt := &gitlab.ListLabelsOptions{
		Page:    1,
		PerPage: 100,
	}

	var allLabels []*gitlab.Label

	for {
		labels, _, err := c.Labels.ListLabels(pid, opt)
		if err != nil {
			return nil, err
		}

		if len(labels) == 0 {
			break
		}

		allLabels = append(allLabels, labels...)
		opt.Page = opt.Page + 1
	}
	return allLabels, nil
}
