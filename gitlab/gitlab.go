package gitlab

import (
	"github.com/xanzy/go-gitlab"
	"strings"
	"gopkg.in/russross/blackfriday.v2"
	"log"
	"strconv"
)

type Client struct {
	*gitlab.Client
}

func New(token string) *Client{
	return &Client{
		 Client: gitlab.NewClient(nil, token),
	}
}

func (c *Client) ListWorks(pid interface{}, prefixClasses *Classes) (works []*Work, err error) {
	issues, _, err := c.Issues.ListProjectIssues(pid, nil)
	return toWorks(issues, prefixClasses), err
}

type Issue struct {
	Title       string
	Description string
	Summary     string
	Note        string
}

type Work struct {
	Issue        *Issue
	Classes      *Classes
	Dependencies []*Issue
}

type Classes struct {
	Large  string
	Middle string
	Small  string
}

func toIssue(gitlabIssue *gitlab.Issue) *Issue {
	return &Issue{
		Title: gitlabIssue.Title,
		Description: gitlabIssue.Description,
	}
}

func toWorks(issues []*gitlab.Issue, classPrefix *Classes) (works []*Work) {
	for _, issue := range issues {
		work := &Work{Issue: toIssue(issue), Classes: &Classes{}}
		for _, label := range issue.Labels {
			if strings.Contains(label, classPrefix.Large) {
				work.Classes.Large = strings.Replace(label, classPrefix.Large, "", 1)
			}
			if strings.Contains(label, classPrefix.Middle) {
				work.Classes.Middle = strings.Replace(label, classPrefix.Middle, "", 1)
			}
			if strings.Contains(label, classPrefix.Small) {
				work.Classes.Small = strings.Replace(label, classPrefix.Small, "", 1)
			}
			IIDs, err := getDependencyIIDs(issue)
			if err != nil {
				panic(err)
			}

			for _, issue := range findIssuesByIIDs(issues, IIDs) {
				work.Dependencies = append(work.Dependencies, toIssue(issue))
			}
			works = append(works, work)
		}
	}
	return
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
