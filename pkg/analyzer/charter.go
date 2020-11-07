package analyzer

import (
	"fmt"
	"os"
	"regexp"

	"github.com/nimakaviani/github-contributors/pkg/scraper"
	"github.com/olekukonko/tablewriter"
)

const (
	Unknown = "unknown"
)

type Details struct {
	alias string
	org   string
	email string
}

type charter struct {
	charterMap map[string]interface{}
}

func NewCharter() *charter {
	c := &charter{
		charterMap: make(map[string]interface{}),
	}

	c.charterMap[Unknown] = make(map[string]Details)
	return c
}

func (c *charter) Build(user string) error {
	email, err := scraper.Find(user)
	if err != nil {
		unknowns := c.charterMap[Unknown].(map[string]Details)
		unknowns[user] = Details{org: Unknown}
		c.charterMap[Unknown] = unknowns
		return err
	}

	return c.parse(user, email)
}

func (c *charter) parse(login, email string) error {
	details, err := extract(login, email)
	if err != nil {
		return err
	}

	users := c.charterMap[details.org]
	if users == nil {
		users = make(map[string]Details)
	}
	indexedUsers := users.(map[string]Details)
	indexedUsers[login] = details
	c.charterMap[details.org] = indexedUsers

	return nil
}

func (c *charter) Write(expand bool) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Org", "GitHub Id", "email"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetBorder(false)

	data := make([][]string, 0)

	for org, users := range c.charterMap {
		if !expand {
			continue
		}

		count := len(users.(map[string]Details))
		for login, details := range users.(map[string]Details) {
			data = append(data, []string{fmt.Sprintf("%s (%d)", org, count), login, details.email})
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
