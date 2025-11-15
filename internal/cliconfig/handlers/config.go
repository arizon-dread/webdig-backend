package handlers

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/arizon-dread/webdig-backend/internal/cliconfig/platform"
	"github.com/arizon-dread/webdig-backend/pkg/types"
	"gopkg.in/yaml.v3"
)

var (
	pathFinder platform.Pathfinder = platform.NewFindPath()
	f          *os.File
)

var (
	ErrMakingDirectory = errors.New("unable to make configDir")
	ErrOpenFIle        = errors.New("uanble to open file")
	ErrReadFile        = errors.New("unable to read file")
	ErrUnmarshal       = errors.New("unable to unmarshal file into go struct")
	ErrWriteFile       = errors.New("unable to write to config file")
)

func ensureAppDir() (string, error) {
	dir := filepath.Join(pathFinder.FindPath(), "webdig")
	err := os.MkdirAll(dir, os.FileMode(0o755))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrMakingDirectory, err)
	}
	return dir, nil
}

func EnsureConfig(dir *string) (*types.ServerConf, error) {
	// make dir path if it doesn't exist
	var err error
	if dir == nil {
		// initialize the string pointer to a string
		dir = new(string)
		*dir, err = ensureAppDir()
		if err != nil {
			log.Printf("error ensuring directory, %v", err)
			return nil, err
		}
	}

	// create a file reference
	file := fmt.Sprintf("%v%vserver.yaml", dir, os.PathSeparator)

	// open the file for reading
	f, err = os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrOpenFIle, err)
	}
	defer f.Close()
	// file exists!
	var b []byte
	_, err = f.Read(b)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrReadFile, err)
	}

	// unmarshal into go struct
	var conf types.ServerConf
	err = yaml.Unmarshal(b, &conf)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrUnmarshal, err)
	}
	// try to find server.yaml and marshal into go struct, otherwise return err and let user specify server
	return &conf, nil
}

func SaveConf(url string) error {
	ErrSaveConf := errors.New("error saving config")
	dir, err := ensureAppDir()
	var conf *types.ServerConf
	if err != nil {
		return err
	}
	if f == nil {
		conf, err = EnsureConfig(&dir)
		if err != nil {
			f, err = os.Create(filepath.Join())
			if err != nil {
				return fmt.Errorf("%w: %w", ErrWriteFile, err)
			}
			*conf = types.ServerConf{}
			defer f.Close()
		}
	}
	conf.Server = url

	c, err := yaml.Marshal(conf)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrUnmarshal, err)
	}
	_, err = f.Write(c)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSaveConf, err)
	}
	return nil
}
