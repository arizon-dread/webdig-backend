//go:build linux
// +build linux

package platform

import (
	"os"
	"path/filepath"
)

type linuxFindPath struct{}

func (l *linuxFindPath) FindPath() string {
	// user level
	var userConfig string
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		userConfig = os.Getenv("XDG_CONFIG_HOME")
	} else {
		userConfig = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return userConfig
	// os level
}

func NewFindPath() Pathfinder {
	return &linuxFindPath{}
}
