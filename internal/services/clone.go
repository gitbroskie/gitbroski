package services

import (
	"fmt"
	"gitbroski/utils/git"
	"gitbroski/utils/helper"
	"gitbroski/utils/logger"
	"os"
	"os/exec"
)

// Clone clones a repository using the active profile's credentials.
// Usage: gitbroski clone <repo-url> [destination]
func Clone(args ...string) {
	if len(args) == 0 {
		logger.Error("Usage: clone <repo-url> [destination]")
		return
	}

	repoURL := args[0]
	destination := ""
	if len(args) > 1 {
		destination = args[1]
	}

	// ── 1. Load active profile ────────────────────────────────────────────────
	profile, name, err := GetActiveProfile()
	if err != nil {
		logger.Error(err.Error())
		return
	}
	logger.Text(fmt.Sprintf("Using profile '%s' (%s)", name, profile.Username))

	// ── 2. Normalize URL to HTTPS ─────────────────────────────────────────────
	// If someone passes an SSH url (git@github.com:user/repo.git),
	// convert it to HTTPS so GIT_ASKPASS works.
	repoURL = helper.NormalizeToHTTPS(repoURL, profile.Host)

	// ── 3. Clear any cached credentials for this host ────────────────────────
	// This ensures git doesn't silently fall back to a previously stored token
	// and is forced to use our GIT_ASKPASS instead.
	helper.ClearCachedCredentials(profile.Host)

	// ── 5. Build git clone command ────────────────────────────────────────────
	gitArgs := []string{"clone", repoURL}
	if destination != "" {
		gitArgs = append(gitArgs, destination)
	}

	cmd := exec.Command("git", gitArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// ── 6. Write a temp askpass script ───────────────────────────────────────
	// GIT_ASKPASS must point to a standalone executable. Git calls it as:
	//   <script> "Username for '...'"
	//   <script> "Password for '...'"
	// We write a small shell script that echoes the right value.
	askpassPath, err := helper.WriteAskpassScript(profile.Username, profile.Token)
	if err != nil {
		logger.Error("Could not write askpass script: " + err.Error())
		return
	}
	defer os.Remove(askpassPath) // clean up after clone

	cmd.Env = append(os.Environ(),
		"GIT_ASKPASS="+askpassPath, // git calls our script for credentials
		"GIT_TERMINAL_PROMPT=0",    // disable git's own prompt fallback
		"GIT_CREDENTIAL_HELPER=",   // bypass credential helpers entirely
	)

	// ── 7. Run clone ──────────────────────────────────────────────────────────
	logger.Text(fmt.Sprintf("Cloning %s ...", repoURL))
	if err := cmd.Run(); err != nil {
		logger.Error("git clone failed: " + err.Error())
		return
	}

	// ── 8. Set git config user.name / user.email locally in cloned repo ───────
	repoDir := destination
	if repoDir == "" {
		repoDir = helper.RepoNameFromURL(repoURL)
	}

	if err := git.SetLocalGitConfig(repoDir, profile.Username, profile.Email); err != nil {
		// Non-fatal — clone succeeded, just warn
		logger.Text(fmt.Sprintf("⚠ Could not set local git config: %s", err.Error()))
	} else {
		logger.Text(fmt.Sprintf("✓ Set git user.name='%s' and user.email='%s' locally in %s",
			profile.Username, profile.Email, repoDir))
	}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// setLocalGitConfig sets user.name and user.email inside the cloned repo
// using --local so it doesn't affect global git config.
