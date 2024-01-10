package grafana

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/grafana-tools/sdk"
)

func (g *Grafana) PullFolders() error {
	var (
		folders []sdk.Folder
		err     error
	)

	if folders, err = g.client.GetAllFolders(g.ctx); err != nil {
		return err
	}
	for _, folder := range folders {
		b, err := json.MarshalIndent(folder, "", "  ")
		if err != nil {
			return err
		}
		if err = writeToFile(g.dir, b, folder.Title, ""); err != nil {
			return err
		}
	}
	return nil
}

func (g *Grafana) PushFolder() error {
	var (
		filesInDir []os.DirEntry
		rawFolder  []byte
		err        error
	)

	if filesInDir, err = os.ReadDir(g.dir); err != nil {
		return err
	}
	for _, file := range filesInDir {
		if filepath.Ext(file.Name()) == ".json" {
			if rawFolder, err = os.ReadFile(filepath.Join(g.dir, file.Name())); err != nil {
				log.Println(err)
				ExecutionErrorHappened = true
				continue
			}
			var folder sdk.Folder
			if err = json.Unmarshal(rawFolder, &folder); err != nil {
				log.Println(err)
				ExecutionErrorHappened = true
				continue
			}
			if _, err := g.client.CreateFolder(g.ctx, folder); err != nil {
				log.Printf("error on importing folder %s", folder.Title)
				ExecutionErrorHappened = true
				continue
			}
		}
	}
	return nil
}

func (g *Grafana) FindFolderId(folderName string) (int, error) {
	allFolders, err := g.client.GetAllFolders(g.ctx)
	if err != nil {
		return 0, err
	}

	for _, folder := range allFolders {
		if folder.Title == folderName {
			return folder.ID, nil
		}
	}
	return 0, nil
}
