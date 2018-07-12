package usecase

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/issummary/issummary/issummary"
	"github.com/mpppk/gitany"
	"golang.org/x/sync/errgroup"
)

type WorkUseCase struct {
	Client       *issummary.Client
	Issues       []*issummary.Issue
	Repositories []*issummary.Repository
	Labels       []*issummary.Label
}

func (wu *WorkUseCase) getListAllGroupIssuesByLabelAsyncFunc(ctx context.Context, org string, targetLabels []string) (func() error, chan []*issummary.Issue) {
	issuesChan := make(chan []*issummary.Issue, 1)
	return func() error {
		log.Printf("Fetch Issues with %v as label from %v", targetLabels, org)
		allIssues, err := wu.Client.ListAllGroupIssuesByLabel(ctx, org, targetLabels)
		issuesChan <- allIssues
		return err
	}, issuesChan
}

func (wu *WorkUseCase) getListAllGroupRepositoriesAndLabelsAsyncFunc(ctx context.Context, org string) (func() error, chan []*issummary.Repository, chan []*issummary.Label) {
	repositoriesChan := make(chan []*issummary.Repository, 1)
	labelsChan := make(chan []*issummary.Label, 1)

	return func() error {
		log.Printf("Fetch Repositories from %v", org)
		repositories, err := wu.Client.ListAllRepositories(ctx, org)
		if err != nil {
			return err
		}

		repositoriesChan <- repositories

		var projectNames []string
		for _, project := range repositories {
			projectNames = append(projectNames, project.GetName())
		}

		log.Printf("Fetch Labels from repositories of %v (%v)", org, projectNames)
		labels, err := wu.Client.ListAllProjectsLabels(ctx, org, repositories)
		labelsChan <- labels
		return err
	}, repositoriesChan, labelsChan
}

func (wu *WorkUseCase) Fetch(ctx context.Context, org string, targetLabels []string) error {
	eg := errgroup.Group{}

	f, issuesChan := wu.getListAllGroupIssuesByLabelAsyncFunc(ctx, org, targetLabels)
	eg.Go(f)

	f2, repositoriesChan, labelsChan := wu.getListAllGroupRepositoriesAndLabelsAsyncFunc(ctx, org)
	eg.Go(f2)

	if err := eg.Wait(); err != nil {
		return err
	}

	wu.Repositories = <-repositoriesChan
	wu.Issues = <-issuesChan
	wu.Labels = <-labelsChan
	return nil
}

func (wu *WorkUseCase) ListGroupWorks(org string, prefix, spLabelPrefix string) (works []*issummary.Work, err error) {
	var filteredIssues []*issummary.Issue

	for _, issue := range wu.Issues {
		for _, project := range wu.Repositories {
			if issue.GetRepositoryID() == project.GetID() {
				filteredIssues = append(filteredIssues, issue)
				break
			}
		}
	}
	works, err = toWorks(org, filteredIssues, wu.Repositories, wu.Labels, prefix, spLabelPrefix)
	if err != nil {
		var projectNames []string
		for _, project := range wu.Repositories {
			projectNames = append(projectNames, project.GetName())
		}
		return nil, fmt.Errorf("failed to convert to works from %v Issues(wu.Repositories are %v): %v", org, projectNames, err)
	}
	return works, nil
}

func toWorks(org string, issues []*issummary.Issue, repositories []*issummary.Repository, labels []*issummary.Label, targetLabelPrefix, spLabelPrefix string) (works []*issummary.Work, err error) {
	for _, issue := range issues {
		work, err := toWork(org, targetLabelPrefix, spLabelPrefix, issue, repositories, labels)
		if err != nil {
			return nil, err
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
				return nil, err
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

func findProjectByID(projects []*issummary.Repository, id int64) (gitany.Repository, bool) {
	for _, project := range projects {
		if project.GetID() == id {
			return project, true
		}
	}

	return nil, false
}
