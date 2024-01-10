package grafana

import (
	"os"
	"path/filepath"
)

var ExecutionErrorHappened = false

func writeToFile(directory string, content []byte, name, tag string) error {
	var (
		err           error
		path          string
		dashboardFile *os.File
	)

	path = directory
	if tag != "" {
		path = filepath.Join(path, tag)
	}

	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0750); err != nil {
			return err
		}
	}

	dashboardFile, err = os.Create(
		filepath.Clean(
			filepath.Join(path, name+".json")))
	if err != nil {
		return err
	}
	defer dashboardFile.Close()

	if err = os.WriteFile(dashboardFile.Name(), content, os.FileMode(0755)); err != nil {
		return err
	}
	return nil
}
