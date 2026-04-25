package mr

import (
	"encoding/json"
	"os/exec"
	"regexp"
	"strings"
)

// Status represents the state of an MR/PR.
type Status struct {
	State string // "open", "merged", "closed", "unknown"
	Title string
	Repo  string
}

// Icon returns the emoji for the status.
func (s Status) Icon() string {
	switch s.State {
	case "open":
		return "🟢"
	case "merged":
		return "🟣"
	case "closed":
		return "🔴"
	default:
		return "⚪"
	}
}

// Label returns the text label for the status.
func (s Status) Label() string {
	switch s.State {
	case "open":
		return "Open"
	case "merged":
		return "Merged"
	case "closed":
		return "Closed"
	default:
		return "Unknown"
	}
}

// FetchStatus gets the status of an MR/PR using gh or glab CLI.
func FetchStatus(url string) Status {
	status := Status{
		State: "unknown",
		Repo:  ExtractRepoName(url),
	}

	if strings.Contains(url, "github.com") {
		return fetchGitHubPRStatus(url, status)
	} else if strings.Contains(url, "gitlab") {
		return fetchGitLabMRStatus(url, status)
	}

	return status
}

func fetchGitHubPRStatus(url string, status Status) Status {
	re := regexp.MustCompile(`github\.com/([^/]+)/([^/]+)/pull/(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) != 4 {
		return status
	}

	owner, repo, prNum := matches[1], matches[2], matches[3]
	status.Repo = repo

	cmd := exec.Command("gh", "pr", "view", prNum, "--repo", owner+"/"+repo, "--json", "state,title")
	output, err := cmd.Output()
	if err != nil {
		return status
	}

	var result struct {
		State string `json:"state"`
		Title string `json:"title"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return status
	}

	status.Title = result.Title
	switch strings.ToLower(result.State) {
	case "open":
		status.State = "open"
	case "merged":
		status.State = "merged"
	case "closed":
		status.State = "closed"
	}

	return status
}

func fetchGitLabMRStatus(url string, status Status) Status {
	re := regexp.MustCompile(`gitlab[^/]*/(.+)/-/merge_requests/(\d+)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) != 3 {
		return status
	}

	projectPath, mrNum := matches[1], matches[2]
	status.Repo = ExtractRepoName(url)

	cmd := exec.Command("glab", "mr", "view", mrNum, "--repo", projectPath, "--output", "json")
	output, err := cmd.Output()
	if err != nil {
		return status
	}

	var result struct {
		State string `json:"state"`
		Title string `json:"title"`
	}
	if err := json.Unmarshal(output, &result); err != nil {
		return status
	}

	status.Title = result.Title
	switch strings.ToLower(result.State) {
	case "opened":
		status.State = "open"
	case "merged":
		status.State = "merged"
	case "closed":
		status.State = "closed"
	}

	return status
}

// ExtractRepoName extracts the repo name from a GitLab/GitHub URL.
func ExtractRepoName(url string) string {
	parts := strings.Split(url, "/")
	for i, part := range parts {
		if part == "-" || part == "merge_requests" || part == "pull" {
			if i > 0 {
				return parts[i-1]
			}
		}
	}
	if len(parts) >= 2 {
		return parts[len(parts)-2]
	}
	return "unknown"
}
