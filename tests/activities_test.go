package tests

import (
	"errors"
	"testing"

	"github.com/nimakaviani/github-contributors/pkg/analyzer"
	"github.com/nimakaviani/github-contributors/pkg/models"
	"github.com/nimakaviani/github-contributors/pkg/scraper/scraperfakes"
)

func TestActivities(t *testing.T) {
	fakeScraper := &scraperfakes.FakeScraper{}
	fakeScraper.ActivitiesReturns([]models.Activity{
		{
			Id:     1234,
			Number: 33,
			Title:  "some issue",
			Url:    "some url",
			User: models.User{
				Login: "nimakaviani",
				Id:    3455,
			},
			State:             "open",
			AuthorAssociation: "member",
		},
	}, nil)

	fakeScraper.FindInRepoCalls(func(repo, login string) (string, error) {
		switch login {
		case "nimakaviani":
			return "something@gmail.com", nil
		default:
			return "", errors.New("unreachable")
		}
	})

	c := analyzer.NewCharter(fakeScraper)
	a := analyzer.NewActivity(fakeScraper, c, models.Issue, 5)

	err := a.Process("some-repo")
	if err != nil {
		t.Fatal(err)
	}

	r, tt, count := fakeScraper.ActivitiesArgsForCall(0)
	if r != "some-repo" || tt != 0 /* issue */ || count != 5 {
		t.Fatal("wrong argumens for activity")
	}

	repo, login := fakeScraper.FindInRepoArgsForCall(0)
	if repo != "some-repo" || login != "nimakaviani" {
		t.Fatal("didnt look for the user")
	}

	if c.Org("nimakaviani") != "gmail.com" {
		t.Fatal("didnt associate user correctly")
	}
}

func TestActivityFailure(t *testing.T) {
	fakeScraper := &scraperfakes.FakeScraper{}

	expectedError := errors.New("some error")
	fakeScraper.ActivitiesReturns([]models.Activity{}, expectedError)

	a := analyzer.NewActivity(fakeScraper, analyzer.NewCharter(fakeScraper), models.Issue, 5)
	if err := a.Process("some-repo"); err != expectedError {
		t.Fatal(err, expectedError)
	}
}
