package commands

import (
	"gitbroski/internal/services"
	"gitbroski/utils/logger"
)

func init() {
	Register("profile", Profile)
}

func Profile(args ...string) {
	logger.Text("Profile command executed")
	services.Profile(args...)
}
