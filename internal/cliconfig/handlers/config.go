package handlers

import (
	"log"
	"os"
	"path/filepath"

	"github.com/arizon-dread/webdig-backend/internal/cliconfig"
	"github.com/arizon-dread/webdig-backend/internal/cliconfig/platform"
)

var pathFinder cliconfig.Pathfinder = platform.NewFindPath()

func ensureAppDir() error {
	err := os.MkdirAll(filepath.Join(pathFinder.FindPath(), "webdig"), os.FileMode(0755))
	if err != nil {
		log.Printf("Error making directory at config path: %v, error: %v", pathFinder.FindPath(), err)
		return err
	}
	return nil
}

func EnsureConfig() error {
	err := ensureAppDir()
	if err != nil {
		log.Printf("error ensuring directory, %v", err)
	}
	//try to find config.yaml and marshal into go struct, otherwise return err and let user specify server
	return nil
}
