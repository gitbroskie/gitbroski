// Package commands provides CLI command handlers.
package commands

import (
	"gitbroski/internal/services/mr"
	"gitbroski/utils/logger"
)

func init() {
	Register("mr", MR)
}

// MR handles the mr subcommand with its various actions.
func MR(args ...string) {
	if len(args) == 0 {
		logger.Text("Usage: gitbroski mr <save|list|open|remove> [args]")
		return
	}

	subCmd := args[0]
	subArgs := args[1:]

	switch subCmd {
	case "save":
		if len(subArgs) == 0 {
			logger.Error("No URL provided. Usage: gitbroski mr save <url>")
			return
		}
		mr.Save(subArgs[0])
	case "list":
		mr.List()
	default:
		logger.Text("Unknown mr subcommand: " + subCmd)
		logger.Text("Available: save, list")
	}
}
