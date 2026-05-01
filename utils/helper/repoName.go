package helper

import (
	"path/filepath"
	"strings"
)

// RepoNameFromURL extracts the repo folder name from a clone URL.
// https://github.com/user/my-repo.git → my-repo
func RepoNameFromURL(url string) string {
	base := filepath.Base(url)
	return strings.TrimSuffix(base, ".git")
}
