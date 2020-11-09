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
)

func Find(user string) (string, error) {
	// // find from profile
	ghUser, err := fromProfile(user)
	if err == nil {
		Log(">> found from email", ghUser.Email)
		return ghUser.Email, nil
	}

	// find from recent activity
	email, err := fromEvents(user, ghUser)
	if err == nil {
		Log(">> found from events", email)
		return email, nil
	}

	// from repo activities
	email, err = fromRepos(user, ghUser)
	if err == nil {
		Log(">> found from repos", email)
		return email, nil
	}

	return "", err
}

func fromProfile(user string) (models.User, error) {
	ghUser := models.User{}
	if err := QueryGithub("profile", fmt.Sprintf("https://api.github.com/users/%s", user), &ghUser); err != nil {
		return ghUser, err
	}

	if ghUser.Email == "" {
		return ghUser, errors.New("not found")
	}

	return ghUser, nil
}

func fromEvents(user string, ghUser models.User) (string, error) {
	ghEvents := []models.Event{}
	if err := QueryGithub("events", fmt.Sprintf("https://api.github.com/users/%s/events?per_page=10", user), &ghEvents); err != nil {
		return "", err
	}

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
				if p.Author.Name == ghUser.Name && email != "" && !strings.Contains(email, "noreply") {
					return email, nil
				}
			}
		}
	}

	return "", errors.New("not found")
}

func fromRepos(user string, ghUser models.User) (string, error) {
	repos := []models.Repo{}
	if err := QueryGithub("repos", fmt.Sprintf("https://api.github.com/users/%s/repos?type=owner&sort=updated&per_page=5", user), &repos); err != nil {
		return "", err
	}

	activity := []models.RepoCommits{}
	for _, r := range repos {
		if err := QueryGithub("activity", fmt.Sprintf("https://api.github.com/repos/%s/commits", r.FullName), &activity); err != nil {
			return "", err
		}

		for _, a := range activity {
			committer := a.Commit.Committer
			if committer.Name == ghUser.Name && committer.Email != "" && !strings.Contains(committer.Email, "noreply") {
				return committer.Email, nil
			}

			author := a.Commit.Author
			if author.Name == ghUser.Name && author.Email != "" && !strings.Contains(author.Email, "noreply") {
				return author.Email, nil
			}
		}
	}

	return "", errors.New("not found")
}

func Contributors(repo string) ([]models.User, error) {
	users := []models.User{}
	err := QueryGithub("contributors", fmt.Sprintf("https://api.github.com/repos/%s/contributors", repo), &users)
	return users, err
}

func Activities(repo string, activity, count int) ([]models.Activity, error) {
	var activityName string
	switch activity {
	case models.Issue:
		activityName = "issues"
	default:
		activityName = "pulls"
	}
	issues := []models.Activity{}
	err := QueryGithub("activity", fmt.Sprintf("https://api.github.com/repos/%s/%s?per_page=%d&state=all&sort=updated", repo, activityName, count), &issues)
	return issues, err
}

func QueryGithub(endpoint, url string, content interface{}) error {
	Log("> query", endpoint, url)

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

func Log(msg ...string) {
	if Debug {
		for _, m := range msg {
			fmt.Printf("%s ", m)
		}
		println()
	}
}
