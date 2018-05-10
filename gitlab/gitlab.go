package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"strings"
)

type Client struct {
	*gitlab.Client
}

func New(token string) *Client{
	return &Client{
		 Client: gitlab.NewClient(nil, token),
	}
}

type Label struct {
	Name        string
	Description string
	Parent      *Label
}

type Issue struct {
	Title       string
	Description string
	Summary     string
	Note        string
}

type Work struct {
	Issue        *Issue
	Label *Label
	Dependencies []*Issue
}

func (c *Client) ListWorks(pid interface{}, prefix string) (works []*Work, err error) {
	allIssues, err := c.listAllIssuesByLabel(pid, gitlab.Labels{"W"})
	if err != nil {
		return nil ,err
	}

	prefixLabels, err := c.listLabelsByPrefix(pid, prefix)

	// TODO: worksをオプティカルソート
	return toWorks(allIssues, prefixLabels, prefix), err
}


func (c *Client) listLabelsByPrefix(pid interface{}, prefix string) (prefixLabels []*gitlab.Label, err error){
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

func (c *Client) listAllIssuesByLabel(pid interface{}, labels gitlab.Labels) ([]*gitlab.Issue, error){
	issueOpt := &gitlab.ListProjectIssuesOptions{
		ListOptions: gitlab.ListOptions{
			Page: 1,
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
		Page: 1,
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
