package github

import (
	"encoding/json"
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
	repo := fetchRepoJson(repoURL, token)
	if repo.Message != "" {
		time.Sleep(1000 * time.Millisecond)
		return GetRepoStars(repoURL, token)
	}
	return repo.StargazersCount
}
