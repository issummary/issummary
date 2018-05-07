package api

import (
	"github.com/xanzy/go-gitlab"
)

type Client struct {
	git          *gitlab.Client
	pid          interface{}
	targetLabels []string
}

func NewClient(token string, pid interface{}, targetLabels []string) *Client {
	git := gitlab.NewClient(nil, token)
	return &Client{
		git:          git,
		pid:          pid,
		targetLabels: targetLabels,
	}
}

func (c *Client) GetIssues() ([]*gitlab.Issue, *gitlab.Response, error) {
	opt := &gitlab.ListProjectIssuesOptions{
		Labels: c.targetLabels,
	}
	return c.git.Issues.ListProjectIssues(c.pid, opt)
}
