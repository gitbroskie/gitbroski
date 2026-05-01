package helper

import (
	"fmt"
	"os"
	"path/filepath"
)

// writeAskpassScript writes a temp shell script that git calls for credentials.
// Git calls it twice:
//
//	script "Username for '...'" → prints username
//	script "Password for '...'" → prints token
func WriteAskpassScript(username, token string) (string, error) {
	script := fmt.Sprintf(`#!/bin/sh
case "$1" in
  Username*) echo "%s" ;;
  Password*) echo "%s" ;;
esac
`, username, token)

	path := filepath.Join(os.TempDir(), "gitbroski-askpass.sh")
	if err := os.WriteFile(path, []byte(script), 0700); err != nil {
		return "", err
	}
	return path, nil
}
