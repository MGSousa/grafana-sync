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

var pullDataSourcesCmd = &cobra.Command{
	Use:              "pull-datasources",
	Short:            "Pull grafana datasources json in to the directory",
	Long:             `Export Datasources in JSON to the specified --directory.`,
	PersistentPreRun: initGrafana,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gs.PullDatasources(); err != nil {
			log.Fatalln("Pull datasources command failed", err)
		}
		if grafana.ExecutionErrorHappened {
			os.Exit(1)
		}
	},
}

var pushDataSourcesCmd = &cobra.Command{
	Use:              "push-datasources",
	Short:            "Read json and create grafana datasources",
	Long:             `Read json with datasources description and publish to grafana.`,
	PersistentPreRun: initGrafana,
	Run: func(cmd *cobra.Command, args []string) {
		if err := gs.PushDatasources(); err != nil {
			log.Fatalln("Push datasources command failed", err)
		}
		if grafana.ExecutionErrorHappened {
			os.Exit(1)
		}
	},
}

func init() {
	pullDataSourcesCmd.PersistentFlags().StringP("tag", "t", "", "Dashboard tag to read")
}
