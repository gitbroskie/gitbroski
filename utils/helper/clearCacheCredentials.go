package helper

import (
	"fmt"
	"os/exec"
	"strings"
)

// clearCachedCredentials flushes any stored credentials for the given host
// from git's credential helper so our GIT_ASKPASS is always used instead.
func ClearCachedCredentials(host string) {
	input := fmt.Sprintf("protocol=https\nhost=%s\n", host)
	cmd := exec.Command("git", "credential", "reject")
	cmd.Stdin = strings.NewReader(input)
	// Ignore errors — if there's nothing cached, this is a no-op
	cmd.Run()
}
