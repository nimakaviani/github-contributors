package analyzer

import (
	"os"
	"sort"
	"strconv"

	"github.com/nimakaviani/github-contributors/pkg/scraper"
	"github.com/olekukonko/tablewriter"
)

type EnhancedIssue struct {
	scraper.Issue
	org string
}

type issues struct {
	issues  []*EnhancedIssue
	charter *charter
	enabled bool
	count   int
}

func NewIssues(c *charter, enabled bool, count int) *issues {
	return &issues{
		charter: c,
		enabled: enabled,
		count:   count,
		issues:  make([]*EnhancedIssue, 0),
	}
}

func (i *issues) Process(repo string) error {
	if !i.enabled {
		return nil
	}

	issues, err := scraper.Issues(repo, i.count)
	if err != nil {
		return err
	}

	scraper.Log("> reviewing issues ...")
	for _, issue := range issues {
		details, err := i.charter.Associate(issue.User.Login, issue.AuthorAssociation)
		if err != nil {
			continue
		}

		i.issues = append(i.issues, &EnhancedIssue{Issue: issue, org: details.org})

	}

	return nil
}

func (i *issues) Write() {
	if !i.enabled {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Id", "Org", "GitHubId", "Issue", "State"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold}, tablewriter.Colors{tablewriter.Bold})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)
	table.SetAutoMergeCells(true)

	issues := i.issues
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].org < issues[j].org {
			return true
		}

		if issues[i].org == issues[j].org {
			return issues[i].User.Login < issues[j].User.Login
		}

		return false
	})

	data := make([][]string, 0)
	for _, issue := range issues {
		row := []string{strconv.Itoa(issue.Id), issue.org, issue.User.Login, issue.Title, issue.State}
		data = append(data, row)
	}

	table.AppendBulk(data)
	table.Render()
}
