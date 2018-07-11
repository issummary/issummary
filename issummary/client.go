package issummary

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mpppk/gitany"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	client       gitany.Client
	Repositories []*Repository
	Issues       []*Issue
	Labels       []*Label
}

func New(client gitany.Client) *Client {
	return &Client{
		client: client,
	}
}

func (c *Client) Fetch(ctx context.Context, org string, targetLabels []string) error {
	eg := errgroup.Group{}
	issuesChan := make(chan []gitany.Issue, 1)
	repositoriesChan := make(chan []gitany.Repository, 1)
	labelsChan := make(chan []gitany.Label, 1)

	eg.Go(func() error {
		log.Printf("Fetch Issues with %v as label from %v", targetLabels, org)
		allIssues, err := c.listAllGroupIssuesByLabel(ctx, org, targetLabels)
		issuesChan <- allIssues
		return err
	})

	eg.Go(func() error {
		log.Printf("Fetch Repositories from %v", org)
		projects, err := c.listAllProjects(ctx, org)
		if err != nil {
			return err
		}

		repositoriesChan <- projects

		log.Printf("Fetch Labels from %v", org)
		labels, err := c.listAllProjectsLabels(ctx, org, projects)
		labelsChan <- labels
		return err
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	c.Repositories = toRepositories(<-repositoriesChan)
	issues, err := toIssues(<-issuesChan)
	if err != nil {
		return err
	}
	c.Issues = issues

	c.Labels = toLabels(<-labelsChan)
	return nil
}

func (c *Client) ListGroupWorks(org string, prefix, spLabelPrefix string) (works []*Work, err error) {
	var filteredIssues []*Issue

	for _, issue := range c.Issues {
		for _, project := range c.Repositories {
			if issue.GetRepositoryID() == project.GetID() {
				filteredIssues = append(filteredIssues, issue)
				break
			}
		}
	}
	works, err = toWorks(org, filteredIssues, c.Repositories, c.Labels, prefix, spLabelPrefix)
	if err != nil {
		var projectNames []string
		for _, project := range c.Repositories {
			projectNames = append(projectNames, project.GetName())
		}
		return nil, fmt.Errorf("failed to convert to works from %v Issues(c.Repositories are %v): %v", org, projectNames, err)
	}
	return works, nil
}

func (c *Client) listLabelsByPrefix(ctx context.Context, owner, repo, prefix string) (prefixLabels []gitany.Label, err error) {
	labels, err := c.listAllLabels(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	for _, label := range labels {
		if strings.Contains(label.GetName(), prefix) {
			prefixLabels = append(prefixLabels, label)
		}
	}
	return prefixLabels, nil
}

func (c *Client) listAllGroupIssuesByLabel(ctx context.Context, org string, labels []string) ([]gitany.Issue, error) {
	issueOpt := &gitany.IssueListOptions{
		ListOptions: gitany.ListOptions{
			Page:    1,
			PerPage: 100,
		},
		Labels: labels,
		State:  "open", // FIXME
	}

	var allIssues []gitany.Issue

	for {
		log.Printf("fetch Issues from org:%v Page:%v", org, issueOpt.Page)
		issues, _, err := c.client.GetIssues().ListByOrg(ctx, org, issueOpt)
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

func (c *Client) listAllProjectIssuesByLabel(ctx context.Context, owner, repo string, labels []string) ([]gitany.Issue, error) {
	issueOpt := &gitany.IssueListByRepoOptions{
		ListOptions: gitany.ListOptions{
			Page:    1,
			PerPage: 100,
		},
		Labels: labels,
	}

	var allIssues []gitany.Issue

	for {
		issues, _, err := c.client.GetIssues().ListByRepo(ctx, owner, repo, issueOpt)
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

func (c *Client) listAllProjectsLabels(ctx context.Context, org string, projects []gitany.Repository) (allLabels []gitany.Label, err error) {
	labelChan := make(chan gitany.Label, 100)
	eg := errgroup.Group{}

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	for _, project := range projects {
		project := project
		eg.Go(func() error {
			labels, err := c.listAllLabels(ctx, org, project.GetName())
			if err != nil {
				return err
			}

			for _, label := range labels {
				labelChan <- label
			}

			return nil
		})
	}

	go func() {
		if err := eg.Wait(); err != nil {
			log.Fatal(err)
		} else {
			close(labelChan)
		}
	}()

	labelMap := map[int64]gitany.Label{}
	for label := range labelChan {
		labelMap[label.GetID()] = label
	}

	return labelMapToSlice(labelMap), nil
}

func labelMapToSlice(labelMap map[int64]gitany.Label) (labels []gitany.Label) {
	for _, label := range labelMap {
		labels = append(labels, label)
	}
	return labels
}

func (c *Client) listAllLabels(ctx context.Context, owner, repo string) ([]gitany.Label, error) {
	opt := &gitany.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	var allLabels []gitany.Label

	for {
		labels, _, err := c.client.GetIssues().ListLabels(ctx, owner, repo, opt)
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

func findProjectByID(projects []*Repository, id int64) (gitany.Repository, bool) {
	for _, project := range projects {
		if project.GetID() == id {
			return project, true
		}
	}

	return nil, false
}

func findProjectByName(projects []*Repository, name string) (gitany.Repository, bool) {
	for _, project := range projects {
		if project.GetName() == name {
			return project, true
		}
	}

	return nil, false
}

func (c *Client) listAllProjects(ctx context.Context, org string) ([]gitany.Repository, error) {
	opt := &gitany.RepositoryListByOrgOptions{
		ListOptions: gitany.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	var allProjects []gitany.Repository

	for {
		repositories, _, err := c.client.GetRepositories().ListByOrg(ctx, org, opt)
		if err != nil {
			return nil, err
		}

		if len(repositories) == 0 {
			break
		}

		allProjects = append(allProjects, repositories...)
		opt.Page = opt.Page + 1
	}
	return allProjects, nil
}

func (c *Client) ListGroupMilestones(ctx context.Context, org string) ([]*Milestone, error) {
	opt := &gitany.MilestoneListOptions{
		ListOptions: gitany.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	var allMilestones []gitany.Milestone

	for {
		milestones, _, err := c.client.GetIssues().ListMilestonesByOrg(ctx, org, opt)
		if err != nil {
			return nil, err
		}

		if len(milestones) == 0 {
			break
		}

		allMilestones = append(allMilestones, milestones...)
		opt.Page = opt.Page + 1
	}

	return toMilestones(allMilestones), nil
}
