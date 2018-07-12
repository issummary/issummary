package issummary

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/mpppk/gitany"
	"github.com/pkg/errors"
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

func (c *Client) listLabelsByPrefix(ctx context.Context, owner, repo, prefix string) (prefixLabels []gitany.Label, err error) {
	labels, err := c.listAllLabels(ctx, owner, repo)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to list labels (owner:%v repo: %v, prefix: %v)\n", owner, repo, prefix))
	}

	for _, label := range labels {
		if strings.Contains(label.GetName(), prefix) {
			prefixLabels = append(prefixLabels, label)
		}
	}
	return prefixLabels, nil
}

func (c *Client) ListAllGroupIssuesByLabel(ctx context.Context, org string, labels []string) ([]*Issue, error) {
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
		log.Printf("fetch issues from org:%v Page:%v", org, issueOpt.Page)
		issues, _, err := c.client.GetIssues().ListByOrg(ctx, org, issueOpt)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to get issues which org is %v from gitany client\n", org))
		}

		if len(issues) == 0 {
			break
		}

		allIssues = append(allIssues, issues...)
		issueOpt.Page = issueOpt.Page + 1
	}

	return toIssues(allIssues)
}

func (c *Client) ListAllProjectsLabels(ctx context.Context, org string, repositories []*Repository) (allLabels []*Label, err error) {
	labelChan := make(chan gitany.Label, 100)
	eg := errgroup.Group{}

	for _, repository := range repositories {
		repository := repository
		eg.Go(func() error {
			labels, err := c.listAllLabels(ctx, org, repository.GetName())
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to list labels (org: %v, repo: %v)\n", org, repository.GetName()))
			}

			for _, label := range labels {
				labelChan <- label
			}

			return nil
		})
	}

	go func() {
		if err := eg.Wait(); err != nil {
			log.Fatal(err) // FIXME
		} else {
			close(labelChan)
		}
	}()

	labelMap := map[int64]gitany.Label{}
	for label := range labelChan {
		labelMap[label.GetID()] = label
	}

	labels := labelMapToSlice(labelMap)
	return toLabels(labels), nil

}

func labelMapToSlice(labelMap map[int64]gitany.Label) (labels []gitany.Label) {
	for _, label := range labelMap {
		labels = append(labels, label)
	}
	return labels
}

func (c *Client) listAllLabels(ctx context.Context, owner, repo string) ([]*Label, error) {
	opt := &gitany.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	var allLabels []gitany.Label

	for {
		labels, res, err := c.client.GetIssues().ListLabels(ctx, owner, repo, opt)
		if err != nil {
			log.Printf("failed to get issues(owner: %v, repo: %v) from gitany client. response: %#v\n", owner, repo, res) // FIXME
			return nil, errors.Wrap(err, fmt.Sprintf("failed to get issues(owner: %v, repo: %v) from gitany client\n", owner, repo))
		}

		if len(labels) == 0 {
			break
		}

		allLabels = append(allLabels, labels...)
		opt.Page = opt.Page + 1
	}
	return toLabels(allLabels), nil
}

func (c *Client) ListAllRepositories(ctx context.Context, org string) ([]*Repository, error) {
	opt := &gitany.RepositoryListByOrgOptions{
		ListOptions: gitany.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	var allRepositories []gitany.Repository

	for {
		repositories, _, err := c.client.GetRepositories().ListByOrg(ctx, org, opt)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to list repositories(org: %v) from gitany client: \n", org))
		}

		if len(repositories) == 0 {
			break
		}

		allRepositories = append(allRepositories, repositories...)
		opt.Page = opt.Page + 1
	}

	return toRepositories(allRepositories), nil
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
			return nil, errors.Wrap(err, fmt.Sprintf("failed to list milestones(org: %v)\n", org))
		}

		if len(milestones) == 0 {
			break
		}

		allMilestones = append(allMilestones, milestones...)
		opt.Page = opt.Page + 1
	}

	return toMilestones(allMilestones), nil
}
