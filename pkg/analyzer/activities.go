package analyzer

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/cheggaaa/pb/v3"
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

	fmt.Fprintf(os.Stderr, "Analyzig the last %d %ss on %s\n", len(activities), i.activityName, repo)
	bar := pb.StartNew(len(activities))

	utils.Log("> reviewing activities ...")
	for _, activity := range activities {
		bar.Increment()

		details, err := i.charter.Associate(repo, activity.User.Login, activity.AuthorAssociation)
		if err != nil {
			continue
		}

		i.activities = append(i.activities, &EnhancedActivity{Activity: activity, userDetails: details})

	}

	bar.Finish()
	return nil
}

func (i *activities) Write(expand bool, format string) error {
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

	type issue struct {
		Issue string `json:"issue"`
		State string `json:"state"`
	}

	type contributor struct {
		Login       string  `json:"login"`
		Email       string  `json:"email,omitempty"`
		Association string  `json:"association,omitempty"`
		Activities  []issue `json:"activities,omitempty"`
	}

	data := make([][]string, 0)

	contributors := make(map[string]*contributor)
	output := make(map[string][]*contributor)

	for _, activity := range activities {
		if _, ok := output[activity.userDetails.org]; !ok {
			output[activity.userDetails.org] = []*contributor{}
		}

		var (
			c  *contributor
			ok bool
		)

		if c, ok = contributors[activity.User.Login]; !ok {
			c = &contributor{
				Login:       activity.User.Login,
				Email:       activity.userDetails.email,
				Association: activity.userDetails.association,
			}
			contributors[activity.User.Login] = c
			output[activity.userDetails.org] = append(output[activity.userDetails.org], c)
		}

		var row []string
		if expand {
			issueString := fmt.Sprintf("%s(%s) : %s \n\n %s", i.activityName, strconv.Itoa(activity.Number), activity.Title, activity.Url)
			row = []string{
				activity.userDetails.org,
				activity.User.Login,
				activity.userDetails.email,
				issueString,
				activity.userDetails.association,
				activity.State,
			}

			c.Activities = append(c.Activities, issue{
				State: activity.State,
				Issue: issueString,
			})

		} else {
			row = []string{
				activity.userDetails.org,
				activity.User.Login,
				activity.userDetails.association,
			}
		}
		data = append(data, row)
	}

	if format != JsonFormat {
		table := tablewriter.NewWriter(os.Stdout)
		if expand {
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
		} else {
			table.SetHeader([]string{
				"Org",
				"GitHubId",
				"Association",
			})
			table.SetHeaderColor(
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
			)
		}
		table.SetAutoMergeCells(true)
		table.SetBorder(true)
		table.SetRowLine(true)

		table.AppendBulk(data)
		table.Render()
		return nil
	}

	type jsonOutput struct {
		Org          string         `json:"org"`
		Percentage   string         `json:"percentage"`
		Contributors []*contributor `json:"contributors"`
	}

	jsonStruct := []jsonOutput{}
	for org, contribs := range output {
		percentage := float64(float64(len(contribs))/float64(i.count)) * 100.0
		jsonStruct = append(jsonStruct, jsonOutput{
			Org:          org,
			Contributors: contribs,
			Percentage:   fmt.Sprintf("%.1f%%", percentage),
		})
	}

	jsonData, err := json.Marshal(jsonStruct)
	if err != nil {
		return err
	}
	fmt.Println(string(jsonData))
	return nil
}
