package tests

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nimakaviani/github-contributors/pkg/models"
	"github.com/nimakaviani/github-contributors/pkg/scraper"
)

func TestFindFromProfile(t *testing.T) {
	data, err := ioutil.ReadFile("data/from_profile.txt")
	if err != nil {
		t.Fatal(err)
	}

	responseMap := map[string][]byte{"/users/nimakaviani": data}
	srv := serverMock(responseMap)
	defer srv.Close()

	scraper := scraper.NewGithubScraper(srv.URL)
	a, err := scraper.Find("nimakaviani")
	if err != nil {
		t.Fatal(err)
	}

	if a != "something@gmail.com" {
		t.Fatal("didnt find the email correctly")
	}
}

func TestFindFromEvents(t *testing.T) {
	data1, err := ioutil.ReadFile("data/from_profile_no_email.txt")
	if err != nil {
		t.Fatal(err)
	}

	data2, err := ioutil.ReadFile("data/events.txt")
	if err != nil {
		t.Fatal(err)
	}

	responseMap := map[string][]byte{
		"/users/nimakaviani":        data1,
		"/users/nimakaviani/events": data2,
	}
	srv := serverMock(responseMap)
	defer srv.Close()

	scraper := scraper.NewGithubScraper(srv.URL)
	a, err := scraper.Find("nimakaviani")
	if err != nil {
		t.Fatal(err)
	}

	if a != "something@gmail.com" {
		t.Fatal("didnt find the email correctly")
	}
}

func TestFindFromRepos(t *testing.T) {
	data1, err := ioutil.ReadFile("data/from_profile_no_email.txt")
	if err != nil {
		t.Fatal(err)
	}

	data2, err := ioutil.ReadFile("data/events_no_email.txt")
	if err != nil {
		t.Fatal(err)
	}

	data3, err := ioutil.ReadFile("data/repos.txt")
	if err != nil {
		t.Fatal(err)
	}

	data4, err := ioutil.ReadFile("data/repo.txt")
	if err != nil {
		t.Fatal(err)
	}

	responseMap := map[string][]byte{
		"/users/nimakaviani":                data1,
		"/users/nimakaviani/events":         data2,
		"/users/nimakaviani/repos":          data3,
		"/repos/nimakaviani/myrepo/commits": data4,
	}
	srv := serverMock(responseMap)
	defer srv.Close()

	scraper := scraper.NewGithubScraper(srv.URL)
	a, err := scraper.Find("nimakaviani")
	if err != nil {
		t.Fatal(err)
	}

	if a != "something@gmail.com" {
		t.Fatal("didnt find the email correctly")
	}
}

func TestScraperActivities(t *testing.T) {
	data, err := ioutil.ReadFile("data/issue.txt")
	if err != nil {
		t.Fatal(err)
	}

	responseMap := map[string][]byte{"/repos/test/issues": data}
	srv := serverMock(responseMap)
	defer srv.Close()

	scraper := scraper.NewGithubScraper(srv.URL)
	a, err := scraper.Activities("test", models.Issue, 2)
	if err != nil {
		t.Fatal(err)
	}

	if len(a) != 2 {
		t.Fatal("length doesnt match", len(a))
	}

	if a[0].Id != 741364683 || a[1].Id != 741355006 {
		t.Fatal("didnt parse correctly")
	}
}

func TestScraperContributors(t *testing.T) {
	data, err := ioutil.ReadFile("data/contributors.txt")
	if err != nil {
		t.Fatal(err)
	}

	responseMap := map[string][]byte{"/repos/test/contributors": data}
	srv := serverMock(responseMap)
	defer srv.Close()

	scraper := scraper.NewGithubScraper(srv.URL)
	a, err := scraper.Contributors("test", 2)
	if err != nil {
		t.Fatal(err)
	}

	if len(a) != 2 {
		t.Fatal("length doesnt match", len(a))
	}

	if a[0].Id != 20407524 || a[1].Id != 13653959 {
		t.Fatal("didnt parse correctly")
	}
}

func serverMock(endpointResponse map[string][]byte) *httptest.Server {
	handler := http.NewServeMux()

	for endpoint, response := range endpointResponse {
		response := response
		handler.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) { responseMock(w, r, response) })
	}

	return httptest.NewServer(handler)
}

func responseMock(w http.ResponseWriter, r *http.Request, response []byte) {
	_, _ = w.Write(response)
}
