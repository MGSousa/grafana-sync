package grafana

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/grafana-tools/sdk"
	log "github.com/sirupsen/logrus"
)

// PullDashboard allows to fetch dashboards to synchronize on filesystem
func (g *Grafana) PullDashboard(tag string, folderID int, fetchVersions bool) error {
	var (
		boardLinks []sdk.FoundBoard
		rawBoard   sdk.Board
		properties sdk.BoardProperties
		err        error
	)

	searchParams := []sdk.SearchParam{sdk.SearchType(sdk.SearchTypeDashboard)}
	if folderID != -1 {
		searchParams = append(searchParams, sdk.SearchFolderID(folderID))
	}

	if tag != "" {
		searchParams = append(searchParams, sdk.SearchTag(tag))
	}

	if boardLinks, err = g.client.Search(g.ctx, searchParams...); err != nil {
		return err
	}

	for _, link := range boardLinks {
		if rawBoard, properties, err = g.client.GetDashboardByUID(g.ctx, link.UID); err != nil {
			log.Errorf("%s for %s\n", err, link.URI)

			ExecutionErrorHappened = true
			continue
		}

		b, err := json.MarshalIndent(rawBoard, "", "  ")
		if err != nil {
			return err
		}
		if err = writeToFile(g.dir, b, properties.Slug, tag); err != nil {
			return err
		}
		log.Printf("Dashboard <%s> saved!", properties.Slug)

		// if enabled fetch dashboard versions metadata
		if fetchVersions {
			meta, err := g.client.GetDashboardVersionsByDashboardID(g.ctx, link.ID)
			if err != nil {
				return err
			}

			if b, err = json.Marshal(meta[0]); err != nil {
				return err
			}

			if err = writeToFile(g.dir, b, fmt.Sprintf("%s_meta", properties.Slug), ""); err != nil {
				return err
			}
			log.Printf("Dashboard metadata <%s> saved!", properties.Slug)
		}
	}
	return nil
}

// PushDashboard allows to synchronize dashboards from filesystem to Grafana
func (g *Grafana) PushDashboard(folderId int) error {
	var (
		filesInDir []os.DirEntry
		rawBoard   []byte
		err        error
	)

	if filesInDir, err = os.ReadDir(g.dir); err != nil {
		return err
	}
	for _, file := range filesInDir {
		if filepath.Ext(file.Name()) == ".json" {
			if rawBoard, err = os.ReadFile(filepath.Join(g.dir, file.Name())); err != nil {
				log.Errorln(err)
				ExecutionErrorHappened = true
				continue
			}

			var board sdk.Board
			if err = json.Unmarshal(rawBoard, &board); err != nil {
				log.Errorln(err)
				ExecutionErrorHappened = true
				continue
			}

			params := sdk.SetDashboardParams{
				FolderID:  folderId,
				Overwrite: true,
			}
			if _, err := g.client.SetDashboard(g.ctx, board, params); err != nil {
				log.Errorf("error on importing dashboard %s", board.Title)
				ExecutionErrorHappened = true
				continue
			}
		}
	}
	return nil
}
