package analyzer

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/cheggaaa/pb"
	"github.com/nimakaviani/github-contributors/pkg/models"
	"github.com/nimakaviani/github-contributors/pkg/scraper"
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
}

func NewActivity(a models.GHActivity, c *charter, count int) *activities {
	return &activities{
		charter:      c,
		count:        count,
		activities:   make([]*EnhancedActivity, 0),
		activityType: a,
	}
}

func (i *activities) Process(repo string) error {
	bar := pb.StartNew(i.count)
	activities, err := scraper.Activities(repo, int(i.activityType), i.count)
	if err != nil {
		return err
	}

	scraper.Log("> reviewing activities ...")
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

	var activityName string
	switch i.activityType {
	case models.Issue:
		activityName = "Issue"
	default:
		activityName = "PR"
	}

	data := make([][]string, 0)
	for _, activity := range activities {
		row := []string{
			activity.userDetails.org,
			activity.User.Login,
			activity.userDetails.email,
			fmt.Sprintf("%s(%s) : %s \n\n %s", activityName, strconv.Itoa(activity.Number), activity.Title, activity.Url),
			activity.userDetails.association,
			activity.State,
		}
		data = append(data, row)
	}

	table.AppendBulk(data)
	table.Render()
}
