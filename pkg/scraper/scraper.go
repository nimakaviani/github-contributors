package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func Find(user string) (string, error) {
	// // find from profile
	ghUser, err := fromProfile(user)
	if err == nil {
		Log("found from email")
		return ghUser.Email, nil
	}

	// find from recent activity
	email, err := fromEvents(user, ghUser)
	if err == nil {
		Log("found from events")
		return email, nil
	}

	// from repo activities
	return fromRepos(user, ghUser)
}

func fromProfile(user string) (GithubUser, error) {
	ghUser := GithubUser{}
	if err := QueryGithub("profile", fmt.Sprintf("https://api.github.com/users/%s", user), &ghUser); err != nil {
		return ghUser, err
	}

	if ghUser.Email == "" {
		return ghUser, errors.New("not found")
	}

	return ghUser, nil
}

func fromEvents(user string, ghUser GithubUser) (string, error) {
	ghEvents := []GithubEvent{}
	if err := QueryGithub("events", fmt.Sprintf("https://api.github.com/users/%s/events?per_page=10", user), &ghEvents); err != nil {
		return "", err
	}

	var email string
	for _, e := range ghEvents {
		if e.Type == "PushEvent" {
			commits := e.Payload.Commits
			if len(commits) == 1 {
				email := commits[0].Author.Email
				if email == "" {
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

func fromRepos(user string, ghUser GithubUser) (string, error) {
	repos := []Repo{}
	if err := QueryGithub("repos", fmt.Sprintf("https://api.github.com/users/%s/repos?type=owner&sort=updated&per_page=5", user), &repos); err != nil {
		return "", err
	}

	activity := []RepoCommits{}
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

func GetContribs(repo string) ([]User, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contributors", repo)
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(body))
	}

	var users []User
	if err = json.Unmarshal(body, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func QueryGithub(endpoint, url string, content interface{}) error {
	Log("> query", endpoint, url)

	githubToken := os.Getenv("GH_EMAIL_TOKEN")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	if !Anonymous {
		req.Header.Set("Authorization", fmt.Sprintf("token: %s", githubToken))
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("failed with status code %d", resp.StatusCode))
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
