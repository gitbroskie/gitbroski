package commands

import (
	"gitbroski/internal/services"
)

func init() {
	Register("clone", Clone)
}

func Clone(args ...string) {
	services.Clone(args...)
}
