package analyzer

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/cheggaaa/pb"
	"github.com/nimakaviani/github-contributors/pkg/scraper"
	"github.com/olekukonko/tablewriter"
)

type EnhancedIssue struct {
	scraper.Issue
	userDetails *Details
}

type issues struct {
	issues  []*EnhancedIssue
	charter *charter
	count   int
}

func NewIssues(c *charter, count int) *issues {
	return &issues{
		charter: c,
		count:   count,
		issues:  make([]*EnhancedIssue, 0),
	}
}

func (i *issues) Process(repo string) error {
	bar := pb.StartNew(i.count)
	issues, err := scraper.Issues(repo, i.count)
	if err != nil {
		return err
	}

	scraper.Log("> reviewing issues ...")
	for _, issue := range issues {
		bar.Increment()

		details, err := i.charter.Associate(issue.User.Login, issue.AuthorAssociation)
		if err != nil {
			continue
		}

		i.issues = append(i.issues, &EnhancedIssue{Issue: issue, userDetails: details})

	}

	bar.Finish()
	return nil
}

func (i *issues) Write() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Org",
		"GitHubId",
		"Email",
		"Issue / PR",
		"Association",
		"State",
	})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoMergeCells(true)

	issues := i.issues
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].userDetails.org < issues[j].userDetails.org {
			return true
		}

		if issues[i].userDetails.org == issues[j].userDetails.org {
			return issues[i].User.Login < issues[j].User.Login
		}

		return false
	})

	data := make([][]string, 0)
	for _, issue := range issues {
		row := []string{
			issue.userDetails.org,
			issue.User.Login,
			issue.userDetails.email,
			fmt.Sprintf("Issue(%s) : %s \n\n %s", strconv.Itoa(issue.Id), issue.Title, issue.Url),
			issue.userDetails.association,
			issue.State,
		}
		data = append(data, row)
	}

	table.AppendBulk(data)
	table.Render()
}
