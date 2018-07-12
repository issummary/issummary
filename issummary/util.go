package issummary

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mpppk/gitany"
	"github.com/pkg/errors"
	"gopkg.in/russross/blackfriday.v2"
)

func toRepository(rawRepository gitany.Repository) *Repository {
	return &Repository{Repository: rawRepository}
}

func toRepositories(rawRepositories []gitany.Repository) (repositories []*Repository) {
	for _, rawRepository := range rawRepositories {
		repositories = append(repositories, toRepository(rawRepository))
	}
	return
}

func toIssue(rawIssue gitany.Issue) (*Issue, error) {
	issueDescription, err := parseIssueDescription(rawIssue.GetBody())
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse issue description which title is %v\n", rawIssue.GetTitle()))
	}

	return &Issue{
		Issue:       rawIssue,
		Description: issueDescription,
	}, nil
}

func toIssues(rawIssues []gitany.Issue) (issues []*Issue, err error) {
	for _, rawIssue := range rawIssues {
		issue, err := toIssue(rawIssue)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to convert gitany issue to issummary issue which title is %v\n", issue.GetTitle()))
		}
		issues = append(issues, issue)
	}
	return
}

func toLabel(rawLabel gitany.Label) (label *Label) {
	return &Label{
		Label:       rawLabel,
		Description: parseLabelDescription(rawLabel.GetDescription()),
	}
}

func toLabels(rawLabels []gitany.Label) (labels []*Label) {
	for _, rawLabel := range rawLabels {
		labels = append(labels, toLabel(rawLabel))
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

func FindLabelByName(labels []*Label, name string) (*Label, bool) {
	for _, label := range labels {
		if label.GetName() == name {
			return toLabel(label), true
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
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse issue dependencies(raw text: %v)\n", description))
	}

	issueDescription.Dependencies = issueDependencies
	summary, err := getMDContentByHeader(node, "Summary")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse issue summary(raw text: %v)\n", description))
	}
	issueDescription.Summary = summary

	note, err := getMDContentByHeader(node, "Note")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse issue note(raw text: %v)\n", description))
	}
	issueDescription.Note = note

	details, err := getMDContentByHeader(node, "Details")
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to parse issue details(raw text: %v)\n", description))
	}
	issueDescription.Details = details

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
						return nil, errors.Wrap(err, fmt.Sprintf("failed to parse issue dependency(%v) to number\n", depStr))
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
							return nil, errors.New(fmt.Sprintf("invalid label syntax in dependencies header: %v\n", depStr))
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
						return nil, errors.Wrap(err, fmt.Sprintf("failed to parse other issue project dependency(%s)", depStr))
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
