//go:build darwin
// +build darwin

package platform

import (
	"os"
)

type darwinFindPath struct{}

func (d *darwinFindPath) FindPath() string {
	return os.Getenv("HOME") + "/Library/Application Support"
}

func NewFindPath() Pathfinder {
	return &darwinFindPath{}
}
