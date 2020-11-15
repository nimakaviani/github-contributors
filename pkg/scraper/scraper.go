package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/nimakaviani/github-contributors/pkg/models"
	"github.com/nimakaviani/github-contributors/pkg/utils"
)

var Anonymous bool

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

//counterfeiter:generate . Scraper
type Scraper interface {
	Find(user string) (string, error)
	FindInRepo(repo, user string) (string, error)
	Contributors(repo string, count int) ([]models.User, error)
	Activities(repo string, activity, count int) ([]models.Activity, error)
}

type githubScraper struct {
	url string
}

func NewGithubScraper(url string) Scraper {
	return &githubScraper{url: url}
}

func (g githubScraper) FindInRepo(repo, user string) (string, error) {
	// find from commits
	email, err := g.fromCommits(repo, user)
	if err == nil {
		utils.Log(">> found from commits", email)
		return email, nil
	}

	return g.Find(user)
}

func (g githubScraper) Find(user string) (string, error) {
	// find from profile
	ghUser, err := g.fromProfile(user)
	if err == nil {
		utils.Log(">> found from profile", ghUser.Email)
		return ghUser.Email, err
	}

	// find from recent activity
	email, err := g.fromEvents(user, ghUser.Name)
	if err == nil {
		utils.Log(">> found from events", email)
		return email, nil
	}

	// from repo activities
	email, err = g.fromRepos(user, ghUser.Name)
	if err == nil {
		utils.Log(">> found from repos", email)
		return email, nil
	}

	return "", err
}

func (g githubScraper) fromCommits(repo, user string) (string, error) {
	activity := []models.RepoCommits{}
	if err := queryGithub("activity", fmt.Sprintf("%s/repos/%s/commits?author=%s", g.url, repo, user), &activity); err != nil {
		return "", err
	}

	for _, a := range activity {
		author := a.Commit.Author
		if author.Email != "" && !strings.Contains(author.Email, "noreply") {
			return author.Email, nil
		}
	}
	return "", errors.New("not found")
}

func (g githubScraper) fromProfile(user string) (models.User, error) {
	ghUser := models.User{}
	if err := queryGithub("profile", fmt.Sprintf("%s/users/%s", g.url, user), &ghUser); err != nil {
		return ghUser, err
	}

	if ghUser.Email == "" {
		return ghUser, errors.New("not found")
	}

	return ghUser, nil
}

func (g githubScraper) fromEvents(user, userName string) (string, error) {
	ghEvents := []models.Event{}
	if err := queryGithub("events", fmt.Sprintf("%s/users/%s/events?per_page=100", g.url, user), &ghEvents); err != nil {
		return "", err
	}

	var count int = 5
	var email string
	for _, e := range ghEvents {
		if e.Type == "PushEvent" {
			commits := e.Payload.Commits
			if len(commits) == 1 {
				email := commits[0].Author.Email
				if email == "" || strings.Contains(email, "noreply") {
					return "", errors.New("not found")
				}
				return email, nil
			}

			for _, p := range commits {
				email = p.Author.Email
				if p.Author.Name == userName && email != "" && !strings.Contains(email, "noreply") {
					return email, nil
				}
			}

			if count -= 1; count == 0 {
				break
			}
		}
	}

	return "", errors.New("not found")
}

func (g githubScraper) fromRepos(user, userName string) (string, error) {
	repos := []models.Repo{}
	if err := queryGithub("repos", fmt.Sprintf("%s/users/%s/repos?type=owner&sort=updated&per_page=5", g.url, user), &repos); err != nil {
		return "", err
	}

	for _, r := range repos {
		email, err := g.fromCommits(r.FullName, user)
		if err == nil {
			return email, nil
		}
	}

	return "", errors.New("not found")
}

func (g githubScraper) Contributors(repo string, count int) ([]models.User, error) {
	users := []models.User{}
	err := queryGithub("contributors", fmt.Sprintf("%s/repos/%s/contributors?per_page=%d", g.url, repo, count), &users)
	return users, err
}

func (g githubScraper) Activities(repo string, activity, count int) ([]models.Activity, error) {
	var activityName string
	switch activity {
	case models.Issue:
		activityName = "issues"
	default:
		activityName = "pulls"
	}
	issues := []models.Activity{}
	err := queryGithub("activity", fmt.Sprintf("%s/repos/%s/%s?per_page=%d&state=all&sort=updated", g.url, repo, activityName, count), &issues)
	return issues, err
}

func queryGithub(endpoint, url string, content interface{}) error {
	utils.Log("> query", endpoint, url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if !Anonymous {
		req.Header.Set("Authorization", os.ExpandEnv("token $GH_EMAIL_TOKEN"))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("failed with status code %d \n %s", resp.StatusCode, string(body)))
	}

	return json.NewDecoder(resp.Body).Decode(content)
}
