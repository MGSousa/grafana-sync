/*
Copyright © 2024

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/mpostument/grafana-sync/grafana"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var pullNotificationsCmd = &cobra.Command{
	Use:              "pull-notifications",
	Short:            "Pull grafana notifications json in to the directory",
	Long:             `Save to the directory grafana folders json. Directory name specified by flag --directory.`,
	PersistentPreRun: initGrafana,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gs.PullNotifications(); err != nil {
			log.Fatalln("Pull notifications command failed", err)
		}
		if grafana.ExecutionErrorHappened {
			os.Exit(1)
		}
	},
}

var pushNotificationsCmd = &cobra.Command{
	Use:              "push-notifications",
	Short:            "Read json and create grafana notifications",
	Long:             `Read json with notifications description and publish to grafana.`,
	PersistentPreRun: initGrafana,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gs.PushNotification(); err != nil {
			log.Fatalln("Push notifications command failed", err)
		}
		if grafana.ExecutionErrorHappened {
			os.Exit(1)
		}
	},
}
