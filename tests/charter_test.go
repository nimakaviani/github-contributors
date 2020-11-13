package tests

import (
	"errors"
	"testing"

	"github.com/nimakaviani/github-contributors/pkg/analyzer"
	"github.com/nimakaviani/github-contributors/pkg/models"
	"github.com/nimakaviani/github-contributors/pkg/scraper/scraperfakes"
)

func TestCharter(t *testing.T) {
	fakeScraper := &scraperfakes.FakeScraper{}
	fakeScraper.ContributorsReturns([]models.User{
		{
			Id:    1234,
			Name:  "Nima Kaviani",
			Login: "nimakaviani",
		},
		{
			Id:    1235,
			Name:  "Tom Hardy",
			Login: "tom",
		},
	}, nil)

	fakeScraper.FindCalls(func(login string) (string, error) {
		switch login {
		case "nimakaviani":
			return "something@gmail.com", nil
		case "tom":
			return "", nil
		default:
			return "", errors.New("unreachable")
		}
	})

	c := analyzer.NewCharter(fakeScraper)
	err := c.Process("test", 10)
	if err != nil {
		t.Fatal(err)
	}

	if fakeScraper.FindArgsForCall(0) != "nimakaviani" {
		t.Fatal("didnt query for the right user")
	}

	if c.Org("nimakaviani") != "gmail.com" {
		t.Fatal("found wrong org for nimakaviani", c.Org("nimakaviani"))
	}

	if c.Org("tom") != "" {
		t.Fatal("found wrong org for tom")
	}

	if fakeScraper.FindCallCount() != 2 || fakeScraper.ContributorsCallCount() != 1 {
		t.Fatal("incorrect call count")
	}
}

func TestContributorFails(t *testing.T) {
	fakeScraper := &scraperfakes.FakeScraper{}

	expectedError := errors.New("some error")
	fakeScraper.ContributorsReturns([]models.User{}, expectedError)

	c := analyzer.NewCharter(fakeScraper)
	if err := c.Process("test", 10); err != expectedError {
		t.Fatal(err)
	}
}

func TestFindFails(t *testing.T) {
	fakeScraper := &scraperfakes.FakeScraper{}
	fakeScraper.ContributorsReturns([]models.User{
		{
			Id:    1234,
			Name:  "Nima Kaviani",
			Login: "nimakaviani",
		},
	}, nil)

	expectedError := errors.New("some errors")
	fakeScraper.FindReturns("", expectedError)

	c := analyzer.NewCharter(fakeScraper)
	if err := c.Process("test", 10); err != nil {
		t.Fatal(err)
	}

	if c.Org("nimakaviani") != "" {
		t.Fatal("found wrong org for nimakaviani", c.Org("nimakaviani"))
	}
}
