package charter

import (
	"fmt"
	"regexp"

	"github.com/nimakaviani/github-contributors/pkg/scraper"
)

const (
	Unknown = "unknown"
)

type Details struct {
	alias string
	org   string
}

type charter struct {
	charterMap map[string]interface{}
}

func Build(users []scraper.User) *charter {
	c := &charter{
		charterMap: make(map[string]interface{}),
	}

	c.charterMap[Unknown] = make(map[string]Details)

	for _, user := range users {
		println("->   user", user.Login)

		email, err := scraper.Find(user.Login)
		if err != nil {
			unknowns := c.charterMap[Unknown].(map[string]Details)
			unknowns[user.Login] = Details{org: Unknown}
			c.charterMap[Unknown] = unknowns
			continue
		}

		err = c.parse(user.Login, email)

		if err != nil {
			scraper.Log(user.Login, err.Error())
		}
	}

	return c
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
	for org, users := range c.charterMap {
		fmt.Printf("> org: %s (%d)\n", org, len(users.(map[string]Details)))
		if !expand {
			continue
		}

		for login, details := range users.(map[string]Details) {
			fmt.Printf("\t%s %s@\n", login, details.alias)
		}
	}
}

func extract(login, input string) (Details, error) {
	rg, err := regexp.Compile("([a-z0-9]+)@([a-z0-9]+).[a-z]+")
	if err != nil {
		return Details{}, err
	}

	orgUser := rg.FindAllStringSubmatch(input, -1)

	if len(orgUser) < 1 || len(orgUser[0]) < 3 {
		return Details{}, nil
	}

	user, org := orgUser[0][1], orgUser[0][2]
	return Details{org: org, alias: user}, nil
}
