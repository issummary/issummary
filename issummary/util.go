package issummary

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mpppk/gitany"
	"gopkg.in/russross/blackfriday.v2"
)

func toIssue(rawIssue gitany.Issue) (*Issue, error) {
	issueDescription, err := parseIssueDescription(rawIssue.GetBody())
	if err != nil {
		return nil, err
	}

	return &Issue{
		Issue:       rawIssue,
		Description: issueDescription,
	}, nil
}

func toLabel(rawLabel gitany.Label) (label *Label) {
	return &Label{
		Label:       rawLabel,
		Description: parseLabelDescription(rawLabel.GetDescription()),
	}
}

func toWorks(org string, issues []gitany.Issue, projects []gitany.Repository, labels []gitany.Label, targetLabelPrefix, spLabelPrefix string) (works []*Work, err error) {
	for _, orgIssue := range issues {
		issue, err := toIssue(orgIssue)
		if err != nil {
			return nil, err
		}

		work := &Work{
			Issue: issue,
		}

		if project, ok := findProjectByID(projects, int64(orgIssue.GetRepositoryID())); ok {
			work.Repository = &Repository{Repository: project}
			work.Issue.ProjectName = project.GetName()
			work.Issue.GroupName = org
		}

		for _, labelName := range orgIssue.GetLabels() {
			if strings.HasPrefix(labelName, targetLabelPrefix) {
				if l, ok := findLabelByName(labels, labelName); ok {
					work.Label = toLabel(l)
				}
				break
			}
		}

		for _, labelName := range orgIssue.GetLabels() {
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

		for _, project := range projects {
			if project.GetID() == orgIssue.GetRepositoryID() {
				work.Issue.ProjectName = project.GetName()
			}
			break
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

func findLabelByName(labels []gitany.Label, name string) (gitany.Label, bool) {
	for _, label := range labels {
		if label.GetName() == name {
			return label, true
		}
	}
	return nil, false
}

func parseIssueDescription(description string) (*IssueDescription, error) {
	issueDescription := &IssueDescription{Raw: description}

	md := blackfriday.New()
	node := md.Parse([]byte(description))

	issueDependencies, err := getIssueDependenciesFromMDNodes(node)
	if err != nil {
		return nil, err
	}

	issueDescription.Dependencies = issueDependencies
	summary, err := getMDContentByHeader(node, "Summary")
	if err != nil {
		return nil, err
	}
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
			return strs, nil
		}

		strs = strs + string(childNode.FirstChild.Literal)
		childNode = childNode.Next
	}
}

func getIssueDependenciesFromMDNodes(node *blackfriday.Node) (*IssueDependencies, error) {
	issueDependencies := &IssueDependencies{}
	childNode := node.FirstChild
	for {
		if childNode == nil {
			return issueDependencies, nil
		}

		if childNode.Type == blackfriday.Heading && string(childNode.FirstChild.Literal) == "dependencies" {
			nextChildNode := childNode.Next
			if nextChildNode == nil {
				return issueDependencies, nil
			}

			dependencyStrs := strings.Split(string(nextChildNode.FirstChild.Literal), " ")
			for i, depStr := range dependencyStrs {
				if i == 0 && nextChildNode.Type == blackfriday.Heading {
					depStr = "#" + depStr
				}

				// Issue dependency
				if strings.HasPrefix(depStr, "#") {
					trimmedDep := strings.TrimLeft(depStr, "#")
					depNum, err := strconv.Atoi(trimmedDep)

					depIssue := &DependIssue{
						Number: depNum,
					}

					if err != nil {
						return nil, err
					}
					issueDependencies.Issues = append(issueDependencies.Issues, depIssue)
					continue
				}

				// Label dependency
				if strings.HasPrefix(depStr, "~") {
					j := i
					for {
						if strings.Count(depStr, `"`) == 2 {
							break
						}
						j++
						if len(dependencyStrs) <= j {
							return nil, errors.New("invalid label syntax in dependencies header")
						}
						depStr += " " + dependencyStrs[j]
					}

					trimmedDep := strings.TrimLeft(depStr, "~")
					labelName := strings.Trim(trimmedDep, `"`)
					issueDependencies.LabelNames = append(issueDependencies.LabelNames, labelName)
					continue
				}

				// Other project issue dependency: ex) awesome_group/awesome_project#3
				if strings.Contains(depStr, "#") {
					depStrList := strings.Split(depStr, "/")
					projectNameAndIssueIIDStr := depStrList[len(depStrList)-1]
					projectNameAndIssueIID := strings.Split(projectNameAndIssueIIDStr, "#")
					projectName := projectNameAndIssueIID[0]
					issueIID, err := strconv.Atoi(projectNameAndIssueIID[1])
					if err != nil {
						return nil, fmt.Errorf("failed to parse other issue project dependency(%s): %s",
							depStr, err)
					}

					opIssue := &DependIssue{
						ProjectName: projectName,
						Number:      issueIID,
						GroupName:   strings.Join(depStrList[:len(depStrList)-1], "/"),
					}

					issueDependencies.Issues = append(issueDependencies.Issues, opIssue)
				}
			}
		}
		childNode = childNode.Next
	}
}
