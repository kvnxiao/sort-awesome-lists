package github

import (
	"testing"
)

func TestGetApiEndpointGitHubHost(t *testing.T) {
	path1 := "/username/repository-name/"
	path2 := "/username/repository-name/tree/master/cli/asdf/"
	path3 := "/username/repository-name/tree/master/cli/asdf"
	path4 := "/username/repository-name"

	expected := "https://api.github.com/repos/username/repository-name"
	r1 := GetApiEndpoint(HostName, path1)
	r2 := GetApiEndpoint(HostName, path2)
	r3 := GetApiEndpoint(HostName, path3)
	r4 := GetApiEndpoint(HostName, path4)

	if expected != r1 || expected != r2 || expected != r3 || expected != r4 {
		t.Errorf("failed to process repo path to repos api endpoint")
	}
}
