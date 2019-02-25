package github

import (
	"errors"
	"strings"
)

func GetApiEndpoint(hostname string, path string) string {
	if hostname == HostName {
		split := strings.Split(path, "/")
		split = split[:3]
		return GetReposEndpoint(strings.Join(split, "/"))
	} else if strings.HasSuffix(hostname, ".github.io") {
		repoPath, err := convertGitHubIOToGitHubRepo(hostname, path)
		if err == nil {
			return GetReposEndpoint(repoPath)
		}
	}
	return ""
}

func convertGitHubIOToGitHubRepo(hostname string, path string) (string, error) {
	if path == "" || path == "/" {
		return "", errors.New("cannot parse a root github.io link without additional path")
	}
	user := hostname[:strings.Index(hostname, ".")]
	return user + strings.TrimRight(path, "/"), nil
}
