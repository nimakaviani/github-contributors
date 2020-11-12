package analyzer

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/cheggaaa/pb"
	"github.com/nimakaviani/github-contributors/pkg/models"
	"github.com/nimakaviani/github-contributors/pkg/scraper"
	"github.com/nimakaviani/github-contributors/pkg/utils"
	"github.com/olekukonko/tablewriter"
)

type EnhancedActivity struct {
	models.Activity
	userDetails *Details
}

type activities struct {
	activities   []*EnhancedActivity
	charter      *charter
	count        int
	activityType models.GHActivity
	activityName string
	scraper      scraper.Scraper
}

func NewActivity(s scraper.Scraper, c *charter, a models.GHActivity, count int) *activities {
	var an string
	switch a {
	case models.Issue:
		an = "Issue"
	default:
		an = "PR"
	}

	return &activities{
		charter:      c,
		count:        count,
		activities:   make([]*EnhancedActivity, 0),
		activityType: a,
		activityName: an,
		scraper:      s,
	}
}

func (i *activities) Process(repo string) error {
	utils.Log("> fetch activity list ...")
	activities, err := i.scraper.Activities(repo, int(i.activityType), i.count)
	if err != nil {
		return err
	}

	fmt.Printf("Analyzig the last %d %ss on %s\n", len(activities), i.activityName, repo)
	bar := pb.StartNew(len(activities))

	utils.Log("> reviewing activities ...")
	for _, activity := range activities {
		bar.Increment()

		details, err := i.charter.Associate(activity.User.Login, activity.AuthorAssociation)
		if err != nil {
			continue
		}

		i.activities = append(i.activities, &EnhancedActivity{Activity: activity, userDetails: details})

	}

	bar.Finish()
	return nil
}

func (i *activities) Write() {
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

	activities := i.activities
	sort.Slice(activities, func(i, j int) bool {
		if activities[i].userDetails.org < activities[j].userDetails.org {
			return true
		}

		if activities[i].userDetails.org == activities[j].userDetails.org {
			return activities[i].User.Login < activities[j].User.Login
		}

		return false
	})

	data := make([][]string, 0)
	for _, activity := range activities {
		row := []string{
			activity.userDetails.org,
			activity.User.Login,
			activity.userDetails.email,
			fmt.Sprintf("%s(%s) : %s \n\n %s", i.activityName, strconv.Itoa(activity.Number), activity.Title, activity.Url),
			activity.userDetails.association,
			activity.State,
		}
		data = append(data, row)
	}

	table.AppendBulk(data)
	table.Render()
}
