package analyzer

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type User struct {
	Login string
	Id    int
}

type GithubUser struct {
	Email string
}

func Run(user string) error {
	if _, err := QueryGithub("profile", fmt.Sprintf("https://api.github.com/users/%s", user)); err != nil {
		return err
	}

	if _, err := QueryGithub("events", fmt.Sprintf("https://api.github.com/users/%s/events", user)); err != nil {
		return err
	}

	if _, err := QueryGithub("activity", fmt.Sprintf("https://api.github.com/users/%s/repos?type=owner&sort=updated", user)); err != nil {
		return err
	}

	return nil
}

func GetContribs(repo string) ([]User, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/contributors", repo)
	println("here", url)
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

func QueryGithub(endpoint, url string) ([]string, error) {
	println("> query", endpoint, url)

	githubToken := os.Getenv("GH_EMAIL_TOKEN")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("token: %s", githubToken))

	println(githubToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return extractEmails(body), nil
}

func extractEmails(resp []byte) []string {
	emails := []string{}
	scanner := bufio.NewScanner(strings.NewReader(string(resp)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "email") {
			emails = append(emails, line)
		}
	}
	return emails
}
