package usecase

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/issummary/issummary/issummary"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type WorkUseCase struct {
	Client       *issummary.Client
	Config       *issummary.Config
	issues       []*issummary.Issue
	repositories []*issummary.Repository
	labels       []*issummary.Label
}

func NewWorkUseCase(client *issummary.Client, config *issummary.Config) *WorkUseCase {
	return &WorkUseCase{
		Client: client,
		Config: config,
	}
}

func (wu *WorkUseCase) GetSortedWorks(ctx context.Context) ([]*issummary.Work, error) {
	workManager := issummary.NewWorkManager()

	for _, org := range wu.Config.Organizations {
		_, _, labels, err := wu.fetchOrgResources(ctx, org, wu.Config.TargetLabelPrefixes)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to fetch resources(org: %v, targetLabelPrefixes: %v)\n", org, wu.Config.TargetLabelPrefixes))
		}

		works, err := wu.listOrgWorks(org)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to list works which org is %v\n", org))
		}

		workManager.AddWorks(works)
		workManager.AddLabels(labels)
	}

	if err := workManager.ResolveDependencies(); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to resolve dependencies of work graph\n"))
	}

	sortedWorks, err := workManager.GetSortedWorks()
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get sorted works from work manager\n"))
	}
	return sortedWorks, nil
}

func (wu *WorkUseCase) getListAllGroupIssuesByLabelAsyncFunc(ctx context.Context, org string, targetLabels []string) (func() error, chan []*issummary.Issue) {
	issuesChan := make(chan []*issummary.Issue, 1)
	return func() error {
		log.Printf("fetchOrgResources issues with %v as label from %v", targetLabels, org)
		allIssues, err := wu.Client.ListAllGroupIssuesByLabel(ctx, org, targetLabels)
		issuesChan <- allIssues
		return errors.Wrap(err, fmt.Sprintf("failed to list issues(labels: %v, org: %v)\n", targetLabels, org))
	}, issuesChan
}

func (wu *WorkUseCase) getListAllGroupRepositoriesAndLabelsAsyncFunc(ctx context.Context, org string) (func() error, chan []*issummary.Repository, chan []*issummary.Label) {
	repositoriesChan := make(chan []*issummary.Repository, 1)
	labelsChan := make(chan []*issummary.Label, 1)

	return func() error {
		log.Printf("fetchOrgResources repositories from %v", org)
		repositories, err := wu.Client.ListAllRepositories(ctx, org)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to list repositories which org is %v\n", org))
		}

		repositoriesChan <- repositories

		var projectNames []string
		for _, project := range repositories {
			projectNames = append(projectNames, project.GetName())
		}

		log.Printf("fetchOrgResources labels from repositories of %v (%v)", org, projectNames)
		labels, err := wu.Client.ListAllProjectsLabels(ctx, org, repositories)
		labelsChan <- labels
		return errors.Wrap(err, fmt.Sprintf("failed to list labels which org is %v\n", org))
	}, repositoriesChan, labelsChan
}

func (wu *WorkUseCase) fetchOrgResources(ctx context.Context, org string, targetLabels []string) (repositories []*issummary.Repository, issues []*issummary.Issue, labels []*issummary.Label, err error) {
	eg := errgroup.Group{}

	f, issuesChan := wu.getListAllGroupIssuesByLabelAsyncFunc(ctx, org, targetLabels)
	eg.Go(f)

	f2, repositoriesChan, labelsChan := wu.getListAllGroupRepositoriesAndLabelsAsyncFunc(ctx, org)
	eg.Go(f2)

	if err := eg.Wait(); err != nil {
		return nil, nil, nil, errors.Wrap(err, fmt.Sprintf("failed to list resorces\n"))
	}

	repositories = <-repositoriesChan
	issues = <-issuesChan
	labels = <-labelsChan

	wu.repositories = repositories
	wu.issues = issues
	wu.labels = labels
	return
}

func (wu *WorkUseCase) listOrgWorks(org string) (works []*issummary.Work, err error) {
	var filteredIssues []*issummary.Issue

	for _, issue := range wu.issues {
		for _, project := range wu.repositories {
			if issue.GetRepositoryID() == project.GetID() {
				filteredIssues = append(filteredIssues, issue)
				break
			}
		}
	}
	works, err = toWorks(org, filteredIssues, wu.repositories, wu.labels, wu.Config.ClassLabelPrefix, wu.Config.SPLabelPrefix)
	if err != nil {
		var projectNames []string
		for _, project := range wu.repositories {
			projectNames = append(projectNames, project.GetName())
		}
		return nil, errors.Wrap(err, fmt.Sprintf("failed to convert to works from %v issues(wu.repositories are %v)", org, projectNames))
	}
	return works, nil
}

func toWorks(org string, issues []*issummary.Issue, repositories []*issummary.Repository, labels []*issummary.Label, targetLabelPrefix, spLabelPrefix string) (works []*issummary.Work, err error) {
	for _, issue := range issues {
		work, err := toWork(org, targetLabelPrefix, spLabelPrefix, issue, repositories, labels)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to create work(org: %v, targetLabelPrefix: %v, spLabelPrefix: %v, issue: %v)\n", org, targetLabelPrefix, spLabelPrefix, issue.GetTitle()))
		}
		works = append(works, work)
	}

	// TODO: Set work properties
	totalStoryPoint := 0
	for _, work := range works {
		totalStoryPoint += work.StoryPoint
		work.TotalStoryPoint = totalStoryPoint
		//work.ManDay = totalStoryPoint / velocity
		//work.TotalManDay = totalStoryPoint / velocity
		// work.CompletionDate = timeNow.Add(work.TotalManDay)
		//work.RemainManDays = date.CountBusinessDay(time.Now(), work.CompletionDate)
	}

	return
}

func toWork(org, targetLabelPrefix, spLabelPrefix string, issue *issummary.Issue, repositories []*issummary.Repository, labels []*issummary.Label) (*issummary.Work, error) {
	work := &issummary.Work{
		Issue: issue,
	}

	if project, ok := findProjectByID(repositories, int64(issue.GetRepositoryID())); ok {
		work.Repository = &issummary.Repository{Repository: project}
		work.Issue.ProjectName = project.GetName()
		work.Issue.GroupName = org
	}

	for _, labelName := range issue.GetLabels() {
		if strings.HasPrefix(labelName, targetLabelPrefix) {
			if l, ok := issummary.FindLabelByName(labels, labelName); ok {
				work.Label = l
			}
			break
		}
	}

	for _, labelName := range issue.GetLabels() {
		if strings.HasPrefix(labelName, spLabelPrefix) {
			spStr := strings.TrimPrefix(labelName, spLabelPrefix)
			sp, err := strconv.Atoi(spStr)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("failed to parse story point of label(%v)\n", labelName))
			}
			work.StoryPoint = sp
			break
		}
	}

	for _, project := range repositories {
		if project.GetID() == issue.GetRepositoryID() {
			work.Issue.ProjectName = project.GetName()
		}
		break
	}
	return work, nil
}
