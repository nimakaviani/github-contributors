package analyzer

import (
	"fmt"
	"os"
	"regexp"

	"github.com/cheggaaa/pb"
	"github.com/nimakaviani/github-contributors/pkg/scraper"
	"github.com/olekukonko/tablewriter"
)

const (
	Unknown = "unknown"
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
}

func NewCharter() *charter {
	c := &charter{
		charterMap: make(map[string]interface{}),
		userOrg:    make(map[string]string),
	}

	c.charterMap[Unknown] = make(map[string]*Details)
	return c
}

func (c *charter) Process(repo string) error {
	scraper.Log("> pulling data from repo", repo)
	users, err := scraper.Contributors(repo)
	if err != nil {
		return err
	}

	bar := pb.StartNew(len(users))
	scraper.Log("> building charter ...")
	for _, user := range users {
		err := c.build(user.Login)
		if err != nil {
			scraper.Log(user.Login, err.Error())
		}
		bar.Increment()
	}

	bar.Finish()
	scraper.Log("> done")
	scraper.Log(">> RESULTS")
	return nil
}

func (c *charter) Associate(login, association string) (*Details, error) {
	if _, ok := c.userOrg[login]; !ok {
		err := c.build(login)
		if err != nil {
			return nil, err
		}
	}

	org := c.userOrg[login]
	userDetails := c.charterMap[org].(map[string]*Details)[login]
	userDetails.association = association
	return userDetails, nil
}

func (c *charter) build(login string) error {
	email, err := scraper.Find(login)
	if err != nil {
		return err
	}

	return c.parse(login, email)
}

func (c *charter) parse(login, email string) error {
	details, err := extract(login, email)
	if err != nil {
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

func (c *charter) Write(expand bool) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{
		"Org",
		"Association",
		"GitHubId",
		"Email",
	})
	table.SetHeaderColor(
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
		tablewriter.Colors{tablewriter.Bold},
	)
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetBorder(true)

	data := make([][]string, 0)
	for org, users := range c.charterMap {
		if !expand {
			continue
		}

		count := len(users.(map[string]*Details))
		for login, details := range users.(map[string]*Details) {
			data = append(data, []string{
				fmt.Sprintf("%s (%d)", org, count),
				details.association,
				login,
				details.email,
			})
		}
	}
	table.AppendBulk(data)
	table.Render()
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
