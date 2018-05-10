package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"strings"
	"gopkg.in/russross/blackfriday.v2"
	"log"
	"strconv"
)

func toIssue(gitlabIssue *gitlab.Issue) *Issue {
	return &Issue{
		Title: gitlabIssue.Title,
		Description: gitlabIssue.Description,
	}
}

func toWorks(issues []*gitlab.Issue, prefixLabels []*gitlab.Label, prefix string) (works []*Work) {
	for _, issue := range issues {
		work := &Work{Issue: toIssue(issue)}

		IIDs, err := getDependencyIIDs(issue)
		if err != nil {
			panic(err)
		}

		for _, issue := range findIssuesByIIDs(issues, IIDs) {
			work.Dependencies = append(work.Dependencies, toIssue(issue))
		}

		for _, labelName := range issue.Labels {
			if strings.HasPrefix(labelName, prefix) {
				if l, ok := findLabelByName(prefixLabels, labelName); ok {
					work.Label = &Label {
						Name: l.Name,
						Description: l.Description,
						// TODO: parentを設定
					}
				}
				break
			}

		}
		works = append(works, work)
	}

	// TODO: 親を設定
	return
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
