//go:build windows
// +build windows

package platform

import (
	"os"

	"github.com/arizon-dread/webdig-backend/internal/cliconfig"
)

type windowsFindPath struct{}

func (w *windowsFindPath) FindPath() string {
	return os.Getenv("APPDATA")
}

func NewFindPath() cliconfig.Pathfinder {
	return &windowsFindPath{}
}
