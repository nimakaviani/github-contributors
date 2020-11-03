package analyzer

import (
	"bufio"
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"strings"

	"github.com/agnivade/levenshtein"
)

const (
	Unknown = "unknown"
)

type Details struct {
	alias   string
	org     string
	trusted bool
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
	cmd := exec.Command("hack/github-email.sh", user)
	out, err := cmd.Output()
	if err != nil {
		return err
	}

	return c.parse(user, string(out))
}

func (c *charter) parse(login, output string) error {
	var (
		details Details
		err     error
	)

	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "@") && !strings.Contains(line, "noreply") {
			details, err = extract(login, line)
			if err != nil {
				return err
			}

			if !details.trusted {
				continue
			}

			users := c.charterMap[details.org]
			if users == nil {
				users = make(map[string]Details)
			}
			indexedUsers := users.(map[string]Details)
			indexedUsers[login] = details
			c.charterMap[details.org] = indexedUsers
			break
		}
	}

	if !details.trusted {
		unknowns := c.charterMap[Unknown].(map[string]Details)
		unknowns[login] = Details{org: Unknown}
		c.charterMap[Unknown] = unknowns
	}

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

	distance := levenshtein.ComputeDistance(login, user)
	trusted := distance < int(math.Min(float64(len(login)), float64(len(user))))

	return Details{org: org, alias: user, trusted: trusted}, nil
}
