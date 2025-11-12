//go:build darwin
// +build darwin

package platform

import (
	"os"

	"github.com/arizon-dread/webdig-backend/internal/cliconfig"
)

type darwinFindPath struct{}

func (d *darwinFindPath) FindPath() string {
	return os.Getenv("HOME") + "/Library/Application Support"
}

func NewFindPath() cliconfig.Pathfinder {
	return &darwinFindPath{}
}
