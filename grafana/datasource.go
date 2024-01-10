package grafana

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/grafana-tools/sdk"
)

// PullDatasources
func (g *Grafana) PullDatasources() error {
	var (
		datasources []sdk.Datasource
		err         error
	)

	if datasources, err = g.client.GetAllDatasources(g.ctx); err != nil {
		return err
	}
	for _, datasource := range datasources {
		b, err := json.MarshalIndent(datasource, "", "  ")
		if err != nil {
			return err
		}
		if err = writeToFile(g.dir, b, datasource.Name, ""); err != nil {
			return err
		}
	}
	return nil
}

// PushDatasources
func (g *Grafana) PushDatasources() error {
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
				continue
			}

			var datasource sdk.Datasource
			if err = json.Unmarshal(rawFolder, &datasource); err != nil {
				log.Println(err)
				ExecutionErrorHappened = true
				continue
			}

			if _, err := g.client.CreateDatasource(g.ctx, datasource); err != nil {
				log.Printf("error on importing folder %s", datasource.Name)
				ExecutionErrorHappened = true
				continue
			}
		}
	}
	return nil
}
