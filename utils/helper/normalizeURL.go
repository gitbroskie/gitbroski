package helper

import (
	"strings"
)

// NormalizeToHTTPS converts SSH-style URLs to HTTPS.
// git@github.com:user/repo.git → https://github.com/user/repo.git
func NormalizeToHTTPS(url, host string) string {
	sshPrefix := "git@" + host + ":"
	if strings.HasPrefix(url, sshPrefix) {
		path := strings.TrimPrefix(url, sshPrefix)
		return "https://" + host + "/" + path
	}
	return url
}
