package analyzer

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"

	"github.com/cheggaaa/pb/v3"
	"github.com/nimakaviani/github-contributors/pkg/scraper"
	"github.com/nimakaviani/github-contributors/pkg/utils"
	"github.com/olekukonko/tablewriter"
)

const (
	Unknown    = "unknown"
	JsonFormat = "json"
)

type Details struct {
	alias       string
	org         string
	email       string
	association string
}

type charter struct {
	charterMap map[string]interface{}
	userOrg    map[string]string
	scraper    scraper.Scraper
	total      float64
}

func NewCharter(scraper scraper.Scraper) *charter {
	c := &charter{
		charterMap: make(map[string]interface{}),
		userOrg:    make(map[string]string),
		scraper:    scraper,
	}

	c.charterMap[Unknown] = make(map[string]*Details)
	return c
}

func (c *charter) Process(repo string, count int) error {
	utils.Log("> pulling data from repo", repo)
	users, err := c.scraper.Contributors(repo, count)
	if err != nil {
		return err
	}

	c.total = math.Max(float64(len(users)), float64(count))
	fmt.Fprintf(os.Stderr, "Analyzig the top %d contributors on %s\n", int(c.total), repo)

	bar := pb.StartNew(len(users))
	defer bar.Finish()
	defer utils.Log("> done")
	defer utils.Log(">> RESULTS")

	utils.Log("> building charter ...")
	for _, user := range users {
		err := c.build(repo, user.Login)
		if err != nil {
			utils.Log(user.Login, err.Error())
		}
		bar.Increment()
	}

	return nil
}

func (c *charter) Associate(repo, login, association string) (*Details, error) {
	if _, ok := c.userOrg[login]; !ok {
		err := c.build(repo, login)
		if err != nil {
			return nil, err
		}
	}

	org := c.userOrg[login]
	userDetails := c.charterMap[org].(map[string]*Details)[login]
	userDetails.association = association
	return userDetails, nil
}

func (c *charter) build(repo, login string) error {
	email, err := c.scraper.FindInRepo(repo, login)
	if err != nil {
		return err
	}

	return c.parse(login, email)
}

func (c *charter) Org(login string) string {
	if org, ok := c.userOrg[login]; ok {
		return org
	}
	return Unknown
}

func (c *charter) parse(login, email string) error {
	details, err := extract(login, email)
	if err != nil || details.org == "" {
		unknowns := c.charterMap[Unknown].(map[string]*Details)
		unknowns[login] = &Details{org: Unknown}
		c.charterMap[Unknown] = unknowns
		c.userOrg[login] = Unknown
		return err
	}

	users := c.charterMap[details.org]
	if users == nil {
		users = make(map[string]*Details)
	}

	indexedUsers := users.(map[string]*Details)
	indexedUsers[login] = &details
	c.charterMap[details.org] = indexedUsers
	c.userOrg[login] = details.org

	return nil
}

func (c *charter) Write(expand bool, format string) error {
	type contributor struct {
		Login string `json:"login"`
		Email string `json:"email"`
	}

	type jsonOutput struct {
		OrgName      string        `json:"org"`
		Percentage   string        `json:"percentage"`
		Contributors []contributor `json:"contributors,omitempty"`
	}

	output := []jsonOutput{}
	data := make([][]string, 0)

	for org, users := range c.charterMap {
		count := len(users.(map[string]*Details))
		percentage := float64(float64(count)/c.total) * 100.0

		outputItem := jsonOutput{
			OrgName:    org,
			Percentage: fmt.Sprintf("%.1f%%", percentage),
		}

		if expand {
			for login, details := range users.(map[string]*Details) {
				data = append(data, []string{
					fmt.Sprintf("%s \n %.1f%% ", org, percentage),
					login,
					details.email,
				})

				outputItem.Contributors = append(
					outputItem.Contributors, contributor{
						Login: login,
						Email: details.email,
					})
			}
		}

		if count == 0 {
			continue
		}

		if !expand {
			data = append(data, []string{
				org,
				strconv.Itoa(count),
				fmt.Sprintf("%.1f%% ", percentage),
			})
		}

		output = append(output, outputItem)
	}

	if format != JsonFormat {
		table := tablewriter.NewWriter(os.Stdout)

		if expand {
			table.SetHeader([]string{
				"Org",
				"GitHubId",
				"Email",
			})
			table.SetHeaderColor(
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
			)
			table.SetAutoMergeCells(true)
		} else {
			table.SetHeader([]string{
				"Org",
				"Count",
				"Percentage",
			})
			table.SetHeaderColor(
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
				tablewriter.Colors{tablewriter.Bold},
			)
		}

		table.SetRowLine(true)
		table.SetBorder(true)
		table.AppendBulk(data)
		table.Render()
		return nil
	}

	jsonData, err := json.Marshal(output)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonData))
	return nil
}

func extract(login, email string) (Details, error) {
	rg, err := regexp.Compile("([a-z0-9]+)@([a-z0-9]+).([a-z]+)$")
	if err != nil {
		return Details{}, err
	}

	orgUser := rg.FindAllStringSubmatch(email, -1)

	if len(orgUser) < 1 || len(orgUser[0]) < 4 {
		return Details{}, nil
	}

	user, org, domain := orgUser[0][1], orgUser[0][2], orgUser[0][3]
	return Details{org: fmt.Sprintf("%s.%s", org, domain), alias: user, email: email}, nil
}
