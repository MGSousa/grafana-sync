package grafana

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/grafana-tools/sdk"
)

// PullNotifications
func (g *Grafana) PullNotifications() error {
	var (
		notifications []sdk.AlertNotification
		err           error
	)

	if notifications, err = g.client.GetAllAlertNotifications(g.ctx); err != nil {
		return err
	}
	for _, notification := range notifications {
		b, err := json.MarshalIndent(notification, "", "  ")
		if err != nil {
			return err
		}
		if err = writeToFile(g.dir, b, notification.Name, ""); err != nil {
			return err
		}
	}
	return nil
}

// PushNotification
func (g *Grafana) PushNotification() error {
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

			var notification sdk.AlertNotification
			if err = json.Unmarshal(rawFolder, &notification); err != nil {
				log.Println(err)
				ExecutionErrorHappened = true
				continue
			}

			if _, err := g.client.CreateAlertNotification(g.ctx, notification); err != nil {
				log.Printf("error on importing notification %s", notification.Name)
				ExecutionErrorHappened = true
				continue
			}
		}
	}
	return nil
}
