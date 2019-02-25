package github

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kvnxiao/sort-awesome-lists/logging"
	"github.com/kvnxiao/sort-awesome-lists/requests"
)

const (
	HostName      = "github.com"
	reposEndpoint = "https://api.github.com/repos"
)

func GetReposEndpoint(repoPath string) string {
	return reposEndpoint + repoPath
}

type Repository struct {
	StargazersCount int    `json:"stargazers_count"`
	Message         string `json:"message,omitempty"`
}

func fetchRepoJson(repoURL string, token string) Repository {
	resp, err := requests.Get(repoURL, map[string][]string{
		"Authorization": {"token " + token},
	})
	if err != nil {
		logging.Printlnf("an error occurred in fetching repository %s: %v", repoURL, err)
		return Repository{}
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var repo Repository
	err = decoder.Decode(&repo)
	if err != nil {
		logging.Printlnf("could not decode JSON body for repository %s", repoURL)
		return Repository{}
	}
	return repo
}

func GetRepoStars(repoURL string, token string) int {
	return getRepoStars(repoURL, token, 5)
}

func getRepoStars(repoURL string, token string, retries int) int {
	repo := fetchRepoJson(repoURL, token)
	if repo.Message != "" && repo.Message != "Not Found" {
		if retries > 0 {
			logging.Printlnf("temporary error message for repo %s: %s. Retrying...", repoURL, repo.Message)
			time.Sleep(500 * time.Millisecond)
			return getRepoStars(repoURL, token, retries-1)
		} else {
			fmt.Printf("failed to retrieve stats for %s after 5 retries", repoURL)
			return 0
		}
	}
	return repo.StargazersCount
}
