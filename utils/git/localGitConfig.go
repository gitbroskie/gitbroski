package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func SetLocalGitConfig(repoDir, username, email string) error {
	cmds := [][]string{
		{"git", "-C", repoDir, "config", "--local", "user.name", username},
		{"git", "-C", repoDir, "config", "--local", "user.email", email},
	}
	for _, args := range cmds {
		out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
		if err != nil {
			return fmt.Errorf("%s: %s", strings.Join(args, " "), string(out))
		}
	}
	return nil
}
