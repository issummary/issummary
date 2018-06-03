package gitlab

import (
	"log"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
	"golang.org/x/sync/errgroup"
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
}

type IssueDescription struct {
	Raw           string
	DependencyIDs *DependencyIDs
	Summary       string
	Note          string
	Details       string
}

type DependencyIDs struct {
	IssueIIDs          []int
	LabelNames         []string
	OtherProjectIssues []*OtherProjectIssueDependency
}

type OtherProjectIssueDependency struct {
	GroupName   string
	ProjectName string
	IssueIID    int
}

type Work struct {
	Issue           *Issue
	Label           *Label
	Dependencies    *Dependencies
	StoryPoint      int
	TotalStoryPoint int
	ManDay          int
	TotalManDay     int
	RemainManDays   int
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

func (c *Client) ListGroupWorks(gid interface{}, prefix, spLabelPrefix string) (works []*Work, err error) {
	eg := errgroup.Group{}
	issuesChan := make(chan []*gitlab.Issue, 1)
	projectsChan := make(chan []*gitlab.Project, 1)
	labelsChan := make(chan []*gitlab.Label, 1)

	eg.Go(func() error {
		allIssues, err := c.listAllGroupIssuesByLabel(gid, gitlab.Labels{"W"}) // TODO: 外から指定できるようにする
		issuesChan <- allIssues
		return err
	})

	eg.Go(func() error {
		projects, err := c.listAllProjects(gid)
		if err != nil {
			return err
		}

		projectsChan <- projects

		labels, err := c.listAllProjectsLabels(projects)
		labelsChan <- labels
		return err
	})

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	works, err = toWorks(<-issuesChan, <-projectsChan, <-labelsChan, prefix, spLabelPrefix)
	if err != nil {
		return nil, err
	}
	return works, nil
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

func (c *Client) listAllGroupIssuesByLabel(pid interface{}, labels gitlab.Labels) ([]*gitlab.Issue, error) {
	issueOpt := &gitlab.ListGroupIssuesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
		Labels: labels,
		State:  gitlab.String("opened"), // FIXME
	}

	var allIssues []*gitlab.Issue

	for {
		issues, _, err := c.Issues.ListGroupIssues(pid, issueOpt)
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

func (c *Client) listAllProjectIssuesByLabel(pid interface{}, labels gitlab.Labels) ([]*gitlab.Issue, error) {
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

func (c *Client) listAllProjectsLabels(projects []*gitlab.Project) (allLabels []*gitlab.Label, err error) {
	labelChan := make(chan *gitlab.Label, 100)
	eg := errgroup.Group{}

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}

	for _, project := range projects {
		project := project
		eg.Go(func() error {
			labels, err := c.listAllLabels(project.ID)
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

	labelMap := map[int]*gitlab.Label{}
	for label := range labelChan {
		labelMap[label.ID] = label
	}

	return labelMapToSlice(labelMap), nil
}

func issueMapToSlice(issueMap map[int]*gitlab.Issue) (issues []*gitlab.Issue) {
	for _, issue := range issueMap {
		issues = append(issues, issue)
	}
	return issues
}

func labelMapToSlice(labelMap map[int]*gitlab.Label) (labels []*gitlab.Label) {
	for _, label := range labelMap {
		labels = append(labels, label)
	}
	return labels
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

func findProjectByID(projects []*gitlab.Project, id int) (*gitlab.Project, bool) {
	for _, project := range projects {
		if project.ID == id {
			return project, true
		}
	}

	return nil, false
}

func findProjectByName(projects []*gitlab.Project, name string) (*gitlab.Project, bool) {
	for _, project := range projects {
		if project.Name == name {
			return project, true
		}
	}

	return nil, false
}

func (c *Client) listAllProjects(gid interface{}) ([]*gitlab.Project, error) {
	opt := &gitlab.ListGroupProjectsOptions{
		Page:    1,
		PerPage: 100,
	}

	var allProjects []*gitlab.Project

	for {
		labels, _, err := c.Groups.ListGroupProjects(gid, opt)
		if err != nil {
			return nil, err
		}

		if len(labels) == 0 {
			break
		}

		allProjects = append(allProjects, labels...)
		opt.Page = opt.Page + 1
	}
	return allProjects, nil
}

func (c *Client) ListGroupMilestones(gid string) ([]*Milestone, error) {
	opt := &gitlab.ListGroupMilestonesOptions{
		ListOptions: gitlab.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}

	var allMilestones []*gitlab.GroupMilestone

	for {
		milestones, _, err := c.GroupMilestones.ListGroupMilestones(gid, opt)
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
