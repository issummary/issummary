package gitlab

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mpppk/gitany"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	*gitlab.Client
	gitanyClient gitany.Client
}

func New(token string, gitanyClient gitany.Client) *Client {
	return &Client{
		gitanyClient: gitanyClient,
		Client:       gitlab.NewClient(nil, token),
	}
}

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

func (c *Client) ListGroupWorks(ctx context.Context, gid string, prefix, spLabelPrefix string) (works []*Work, err error) {
	eg := errgroup.Group{}
	issuesChan := make(chan []gitany.Issue, 1)
	repositoriesChan := make(chan []gitany.Repository, 1)
	labelsChan := make(chan []gitany.Label, 1)

	eg.Go(func() error {
		allIssues, err := c.listAllGroupIssuesByLabel(ctx, gid, gitlab.Labels{"W"}) // TODO: 外から指定できるようにする
		issuesChan <- allIssues
		return err
	})

	eg.Go(func() error {
		projects, err := c.listAllProjects(ctx, gid)
		if err != nil {
			return err
		}

		repositoriesChan <- projects

		labels, err := c.listAllProjectsLabels(ctx, gid, projects)
		labelsChan <- labels
		return err
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	projects := <-repositoriesChan
	issues := <-issuesChan
	var filteredIssues []gitany.Issue

	for _, issue := range issues {
		for _, project := range projects {
			fmt.Printf("%#v\n", issue)
			if issue.GetRepositoryID() == project.GetID() {
				filteredIssues = append(filteredIssues, issue)
				break
			}
		}
	}

	labels := <-labelsChan
	works, err = toWorks(filteredIssues, projects, labels, prefix, spLabelPrefix)
	if err != nil {
		var projectNames []string
		for _, project := range projects {
			projectNames = append(projectNames, project.GetName())
		}
		return nil, fmt.Errorf("failed to convert to works from %v issues(projects are %v): %v", gid, projectNames, err)
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

func (c *Client) listAllGroupIssuesByLabel(ctx context.Context, gid string, labels []string) ([]gitany.Issue, error) {
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
		issues, _, err := c.gitanyClient.GetIssues().ListByOrg(ctx, gid, issueOpt)
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
		issues, _, err := c.gitanyClient.GetIssues().ListByRepo(ctx, owner, repo, issueOpt)
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

func (c *Client) listAllProjectsLabels(ctx context.Context, gid string, projects []gitany.Repository) (allLabels []gitany.Label, err error) {
	labelChan := make(chan gitany.Label, 100)
	eg := errgroup.Group{}

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	for _, project := range projects {
		project := project
		eg.Go(func() error {
			labels, err := c.listAllLabels(ctx, gid, project.GetName())
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
		labels, _, err := c.gitanyClient.GetIssues().ListLabels(ctx, owner, repo, opt)
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

func findProjectByID(projects []gitany.Repository, id int64) (gitany.Repository, bool) {
	for _, project := range projects {
		if project.GetID() == id {
			return project, true
		}
	}

	return nil, false
}

func findProjectByName(projects []gitany.Repository, name string) (gitany.Repository, bool) {
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
		repositories, _, err := c.gitanyClient.GetRepositories().ListByOrg(ctx, org, opt)
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
		milestones, _, err := c.gitanyClient.GetIssues().ListMilestonesByOrg(ctx, org, opt)
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
