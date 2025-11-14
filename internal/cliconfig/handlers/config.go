package handlers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/arizon-dread/webdig-backend/internal/cliconfig/platform"
)

var pathFinder platform.Pathfinder = platform.NewFindPath()

var (
	ErrMakingDirectory = errors.New("unable to make configDir")
	ErrOpenFIle        = errors.New("uanble to open file")
)

func ensureAppDir() (string, error) {
	dir := filepath.Join(pathFinder.FindPath(), "webdig")
	err := os.MkdirAll(dir, os.FileMode(0o755))
	if err != nil {
		return "", fmt.Errorf("%w: %w", &ErrMakingDirectory, err)
	}
	return dir, nil
}

func EnsureConfig() error {
	dir, err := ensureAppDir()
	if err != nil {
		log.Printf("error ensuring directory, %v", err)
		return err
	}
	file := fmt.Sprintf("%v%vserver.conf", dir, os.PathSeparator)
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	// file exists!
	// try to find config.yaml and marshal into go struct, otherwise return err and let user specify server
	return nil
}
