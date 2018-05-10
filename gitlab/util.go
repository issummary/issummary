package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"strings"
	"gopkg.in/russross/blackfriday.v2"
	"log"
	"strconv"
	"fmt"
)

func toIssue(gitlabIssue *gitlab.Issue) *Issue {
	return &Issue{
		ID: gitlabIssue.ID,
		IID: gitlabIssue.IID,
		Title: gitlabIssue.Title,
		Description: gitlabIssue.Description,
		URL: gitlabIssue.WebURL,
	}
}

func toLabel(gitlabLabel *gitlab.Label, otherLabels []*gitlab.Label) (label *Label, err error) {
	label = &Label{
		ID: gitlabLabel.ID,
		Name: gitlabLabel.Name,
		Description: parseLabelDescription(gitlabLabel.Description),
	}

	if label.Description.ParentName != "" {
		parentGitLabLabel, ok := findLabelByName(otherLabels, label.Description.ParentName)
		if !ok {
			return nil, fmt.Errorf("parent(%s) not found\n", label.Description.ParentName)
		}
		parentLabel, err := toLabel(parentGitLabLabel, otherLabels)
		if err != nil {
			return nil, err
		}
		label.Parent = parentLabel
	}

	if len(label.Description.DependLabelNames) > 0 {
		for _, dependLabelName := range label.Description.DependLabelNames {
			dependGitLabLabel, ok := findLabelByName(otherLabels, dependLabelName)
			if !ok {
				return nil, fmt.Errorf("depend label(%s) not found\n", dependLabelName)
			}
			dependLabel, err := toLabel(dependGitLabLabel, otherLabels)
			if err != nil {
				return nil, err
			}
			label.Dependencies = append(label.Dependencies, dependLabel)
		}
	}

	return label, nil
}

func toWorks(issues []*gitlab.Issue, labels []*gitlab.Label, prefix string) (works []*Work, err error) {
	for _, issue := range issues {
		work := &Work{
			Issue: toIssue(issue),
			Dependencies: &Dependencies{},
		}

		IIDs, err := getDependencyIIDs(issue)
		if err != nil {
			panic(err)
		}

		for _, issue := range findIssuesByIIDs(issues, IIDs) {
			work.Dependencies.Issues = append(work.Dependencies.Issues, toIssue(issue))
		}

		for _, labelName := range issue.Labels {
			if strings.HasPrefix(labelName, prefix) {
				if l, ok := findLabelByName(labels, labelName); ok {
					work.Label, err = toLabel(l, labels)
					if err != nil {
						return nil, err
					}
				}
				break
			}
		}
		works = append(works, work)
	}

	return
}

func parseLabelDescription(description string) (*LabelDescription){
	ld := &LabelDescription{Raw: description}
	depsKey := "deps: " // TODO: 別の場所で定義したほうがいい気がする
	parentKey := "parent: " // TODO: 別の場所で定義したほうがいい気がする
	lines := strings.Split(description, ";")
	for _, line := range lines {
		if strings.Contains(line, depsKey) {
			depLabelNamesStr := strings.TrimPrefix(line, depsKey)
			depLabelNamesStr = strings.Trim(depLabelNamesStr, "\"")
			ld.DependLabelNames = strings.Split(depLabelNamesStr, ",")
		}

		if strings.Contains(line, parentKey) {
			parentLabelNamesStr := strings.TrimPrefix(line, parentKey)
			ld.ParentName = strings.Split(parentLabelNamesStr, ",")[0] // FIXME
		}
	}
	return ld
}

func findLabelByName(labels []*gitlab.Label, name string) (*gitlab.Label, bool){
	for _, label := range labels {
		if label.Name == name {
			return label, true
		}
	}
	return nil, false
}

func findIssuesByIIDs(issues []*gitlab.Issue, iidList []int) (filteredIssues []*gitlab.Issue) {
	for _, iid := range iidList {
		if issue, ok := findIssueByIID(issues, iid); ok {
			filteredIssues = append(filteredIssues, issue)
		}
	}
	return
}

func findIssueByIID(issues []*gitlab.Issue, iid int) (*gitlab.Issue, bool) {
	for _, issue := range issues {
		if issue.IID == iid {
			return issue, true
		}
	}
	return nil, false
}

func getDependencyIIDs(issue *gitlab.Issue) ([]int, error) {
	description := issue.Description
	md := blackfriday.New()
	node := md.Parse([]byte(description))
	childNode := node.FirstChild
	for {
		if childNode == nil {
			log.Println("dependencies header not found")
			return []int{}, nil
		}

		if childNode.Type == blackfriday.Heading && string(childNode.FirstChild.Literal) == "dependencies" {
			nextChildNode := childNode.Next
			if nextChildNode == nil {
				log.Println("dependencies list not found")
				return []int{}, nil
			}

			if nextChildNode.Type == blackfriday.Heading {
				dependencyStrs := strings.Split(string(nextChildNode.FirstChild.Literal), " ")
				var dependencies []int
				for _, depStr := range dependencyStrs {
					trimmedDep := strings.TrimLeft(depStr, "#")
					depNum, err := strconv.Atoi(trimmedDep)
					if err != nil {
						return nil, err
					}
					dependencies = append(dependencies, depNum)
				}

				return dependencies, nil
			} else {
				log.Println("dependencies list not found")
				return []int{}, nil
			}
		}
		childNode = childNode.Next
	}
}
