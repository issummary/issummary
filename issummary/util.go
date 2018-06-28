package issummary

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mpppk/gitany"
	"gopkg.in/russross/blackfriday.v2"
)

func toIssues(issues []gitany.Issue) (issummaryIssues []*Issue, err error) {
	issummaryIssues = []*Issue{}
	for _, issue := range issues {
		issue, err := toIssue(issue)
		if err != nil {
			return nil, err
		}
		issummaryIssues = append(issummaryIssues, issue)
	}

	return issummaryIssues, nil
}

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

func toDependLabel(labelName string, labels []gitany.Label, issues []gitany.Issue) (*DependLabel, error) {
	gitlabLabel, ok := findLabelByName(labels, labelName)
	if !ok {
		return nil, errors.New("label name not found: " + labelName)
	}

	label, err := toLabel(gitlabLabel, labels, issues)
	if err != nil {
		return nil, err
	}

	relatedGitLabIssues := findIssuesByLabelName(issues, labelName)

	relatedIssues, err := toIssues(relatedGitLabIssues)
	if err != nil {
		return nil, err
	}

	return &DependLabel{
		RelatedIssues: relatedIssues,
		Label:         label,
	}, nil
}

func toLabel(rawLabel gitany.Label, otherLabels []gitany.Label, issues []gitany.Issue) (label *Label, err error) {
	label = &Label{
		Label:        rawLabel,
		Description:  parseLabelDescription(rawLabel.GetDescription()),
		Dependencies: []*DependLabel{},
	}

	if label.Description.ParentName != "" {
		parentGitLabLabel, ok := findLabelByName(otherLabels, label.Description.ParentName)
		if !ok {
			return nil, fmt.Errorf("parent(%s) not found\n", label.Description.ParentName)
		}
		parentLabel, err := toLabel(parentGitLabLabel, otherLabels, issues)
		if err != nil {
			return nil, err
		}
		label.Parent = parentLabel
	}

	if len(label.Description.DependLabelNames) > 0 {
		for _, dependLabelName := range label.Description.DependLabelNames {
			dependLabel, err := toDependLabel(dependLabelName, otherLabels, issues)
			if err != nil {
				return nil, err
			}
			label.Dependencies = append(label.Dependencies, dependLabel)
		}
	}

	return label, nil
}

func toWorks(issues []gitany.Issue, projects []gitany.Repository, labels []gitany.Label, targetLabelPrefix, spLabelPrefix string) (works []*Work, err error) {
	for _, orgIssue := range issues {
		issue, err := toIssue(orgIssue)
		if err != nil {
			return nil, err
		}

		work := &Work{
			Issue: issue,
			Dependencies: &Dependencies{
				Issues: []*Issue{},
				Labels: []*DependLabel{},
			},
		}

		if project, ok := findProjectByID(projects, int64(orgIssue.GetRepositoryID())); ok {
			work.Repository = &Repository{Repository: project}
		}

		for _, depIssues := range issue.Description.Dependencies.Issues {
			findIssue, ok := findIssueByIIDAndProjectID(issues, depIssues.IID, orgIssue.GetRepositoryID())
			if !ok {
				continue
			}

			is, err := toIssue(findIssue)

			if err != nil {
				return nil, err
			}

			work.Dependencies.Issues = append(work.Dependencies.Issues, is)
		}

		for _, otherIssueDep := range issue.Description.Dependencies.OtherProjectIssues {
			if project, ok := findProjectByName(projects, otherIssueDep.ProjectName); ok {
				if otherIssue, ok := findIssueByIIDAndProjectID(issues, otherIssueDep.IID, project.GetID()); ok {
					ois, err := toIssue(otherIssue)
					if err != nil {
						return nil, err
					}

					ois.ProjectName = otherIssueDep.ProjectName
					ois.GroupName = otherIssueDep.GroupName
					work.Dependencies.Issues = append(work.Dependencies.Issues, ois)
				}
			}
		}

		for _, labelName := range issue.Description.Dependencies.LabelNames {
			dependLabel, err := toDependLabel(labelName, labels, issues)
			if err != nil {
				return nil, fmt.Errorf("failed to find depend label from '%v#%v(%v)': %v", issue.ProjectName, issue.GetNumber(), issue.GetTitle(), err)
			}
			work.Dependencies.Labels = append(work.Dependencies.Labels, dependLabel)
		}

		for _, labelName := range orgIssue.GetLabels() {
			if strings.HasPrefix(labelName, targetLabelPrefix) {
				if l, ok := findLabelByName(labels, labelName); ok {
					work.Label, err = toLabel(l, labels, issues)
					if err != nil {
						return nil, err
					}
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

func findIssuesByLabelName(issues []gitany.Issue, labelName string) (filteredIssues []gitany.Issue) {
	for _, issue := range issues {
		for _, issueLabelName := range issue.GetLabels() {
			if issueLabelName == labelName {
				filteredIssues = append(filteredIssues, issue)
				break
			}
		}
	}
	return
}

func findIssuesByIIDs(issues []gitany.Issue, iidList []int) (filteredIssues []gitany.Issue) {
	for _, iid := range iidList {
		if issue, ok := findIssueByIID(issues, iid); ok {
			filteredIssues = append(filteredIssues, issue)
		}
	}
	return
}

func findIssueByIID(issues []gitany.Issue, iid int) (gitany.Issue, bool) {
	for _, issue := range issues {
		if issue.GetNumber() == iid {
			return issue, true
		}
	}
	return nil, false
}

func findIssueByIIDAndProjectID(issues []gitany.Issue, iid int, projectId int64) (gitany.Issue, bool) {
	for _, issue := range issues {
		if issue.GetNumber() == iid && issue.GetRepositoryID() == projectId {
			return issue, true
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
						IID: depNum,
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
						IID:         issueIID,
						GroupName:   strings.Join(depStrList[:len(depStrList)-1], "/"),
					}

					issueDependencies.OtherProjectIssues = append(issueDependencies.OtherProjectIssues, opIssue)
				}
			}
		}
		childNode = childNode.Next
	}
}
