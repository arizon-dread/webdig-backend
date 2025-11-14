//go:build windows
// +build windows

package platform

import (
	"os"
)

type windowsFindPath struct{}

func (w *windowsFindPath) FindPath() string {
	return os.Getenv("APPDATA")
}

func NewFindPath() Pathfinder {
	return &windowsFindPath{}
}
