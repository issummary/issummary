package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/xanzy/go-gitlab"
	"gopkg.in/russross/blackfriday.v2"
)

func toIssue(gitlabIssue *gitlab.Issue) (*Issue, error) {
	issueDescription, err := parseIssueDescription(gitlabIssue.Description)
	if err != nil {
		return nil, err
	}

	return &Issue{
		ID:          gitlabIssue.ID,
		IID:         gitlabIssue.IID,
		DueDate:     (*time.Time)(gitlabIssue.DueDate),
		Title:       gitlabIssue.Title,
		Description: issueDescription,
		URL:         gitlabIssue.WebURL,
	}, nil
}

func toLabel(gitlabLabel *gitlab.Label, otherLabels []*gitlab.Label) (label *Label, err error) {
	label = &Label{
		ID:          gitlabLabel.ID,
		Name:        gitlabLabel.Name,
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

func toWorks(issues []*gitlab.Issue, labels []*gitlab.Label, targetLabelPrefix, spLabelPrefix string) (works []*Work, err error) {
	for _, gitlabIssue := range issues {
		issue, err := toIssue(gitlabIssue)
		if err != nil {
			return nil, err
		}

		work := &Work{
			Issue: issue,
			Dependencies: &Dependencies{
				Issues: []*Issue{},
				Labels: []*Label{},
			},
		}

		for _, issue := range findIssuesByIIDs(issues, issue.Description.DependencyIIDs) {
			is, err := toIssue(issue)
			if err != nil {
				return nil, err
			}
			work.Dependencies.Issues = append(work.Dependencies.Issues, is)
		}

		for _, labelName := range gitlabIssue.Labels {
			if strings.HasPrefix(labelName, targetLabelPrefix) {
				if l, ok := findLabelByName(labels, labelName); ok {
					work.Label, err = toLabel(l, labels)
					if err != nil {
						return nil, err
					}
				}
				break
			}
		}

		for _, labelName := range gitlabIssue.Labels {
			if strings.HasPrefix(labelName, spLabelPrefix) {
				spStr := strings.TrimPrefix(labelName, spLabelPrefix)
				fmt.Println(spStr)
				sp, err := strconv.Atoi(spStr)
				if err != nil {
					return nil, err
				}
				work.StoryPoint = sp
				break
			}
		}

		works = append(works, work)
	}

	return
}

func parseLabelDescription(description string) *LabelDescription {
	ld := &LabelDescription{Raw: description}
	depsKey := "deps: "     // TODO: 別の場所で定義したほうがいい気がする
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
			ld.ParentName = strings.Trim(ld.ParentName, "\"")
		}
	}
	return ld
}

func findLabelByName(labels []*gitlab.Label, name string) (*gitlab.Label, bool) {
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

func parseIssueDescription(description string) (*IssueDescription, error) {
	issueDescription := &IssueDescription{Raw: description}

	md := blackfriday.New()
	node := md.Parse([]byte(description))

	depIIDs, err := getDependencyIIDsFromMDNodes(node)
	if err != nil {
		return nil, err
	}
	issueDescription.DependencyIIDs = depIIDs
	summary, err := getMDContentByHeader(node, "Summary")
	if err != nil {
		return nil, err
	}
	fmt.Println("summary")
	fmt.Println(summary)
	issueDescription.Summary = summary

	note, err := getMDContentByHeader(node, "Note")
	if err != nil {
		return nil, err
	}
	issueDescription.Note = note

	detail, err := getMDContentByHeader(node, "Details")
	if err != nil {
		return nil, err
	}
	issueDescription.Details = detail

	return issueDescription, nil
}

func getMDContentByHeader(node *blackfriday.Node, header string) (string, error) {
	childNode := node.FirstChild
	for {
		if childNode == nil {
			log.Printf("header(%s) not found\n", header)
			return "", nil
		}

		if childNode.Type == blackfriday.Heading && string(childNode.FirstChild.Literal) == header {
			break
		}
		childNode = childNode.Next
	}

	childNode = childNode.Next

	strs := ""
	for {

		if childNode == nil || childNode.Type == blackfriday.Heading {
			fmt.Println("strs")

			return strs, nil
		}

		if header == "Summary" {
			fmt.Println(childNode.String(), string(childNode.FirstChild.Literal))
		}

		strs = strs + string(childNode.FirstChild.Literal)
		childNode = childNode.Next
	}
}

func getDependencyIIDsFromMDNodes(node *blackfriday.Node) ([]int, error) {
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
