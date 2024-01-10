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

var pullDashboardsCmd = &cobra.Command{
	Use:   "pull-dashboards",
	Short: "Pull grafana dashboards in to the directory",
	Long: `Save to the directory grafana dashboards.
Directory name specified by flag --directory. If flag --tag is used,
only dashboards with given tag are pulled`,
	PersistentPreRun: initGrafana,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			folderId int
			err      error
		)
		tag, _ := cmd.Flags().GetString("tag")
		folderName, _ := cmd.Flags().GetString("folderName")
		fetchVersions, _ := cmd.Flags().GetBool("fetchVersions")

		if folderName != "" {
			folderId, err = gs.FindFolderId(folderName)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			folderId, _ = cmd.Flags().GetInt("folderId")
		}

		if err := gs.PullDashboard(tag, folderId, fetchVersions); err != nil {
			log.Fatalln("Pull dashboards command failed", err)
		}
		if grafana.ExecutionErrorHappened {
			os.Exit(1)
		}
	},
}

var pushDashboardsCmd = &cobra.Command{
	Use:              "push-dashboards",
	Short:            "Push grafana dashboards from directory",
	Long:             `Read json with dashboards description and publish to grafana.`,
	PersistentPreRun: initGrafana,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			folderId int
			err      error
		)
		folderName, _ := cmd.Flags().GetString("folderName")

		if folderName != "" {
			folderId, err = gs.FindFolderId(folderName)
			if err != nil {
				log.Fatalln(err)
			}
		} else {
			folderId, _ = cmd.Flags().GetInt("folderId")
		}

		if err := gs.PushDashboard(folderId); err != nil {
			log.Fatalln("Push dashboards command failed", err)
		}
		if grafana.ExecutionErrorHappened {
			os.Exit(1)
		}
	},
}

func init() {
	pushDashboardsCmd.PersistentFlags().IntP("folderId", "i", 0, "Grafana dir ID to which push dashboards")
	pushDashboardsCmd.PersistentFlags().StringP("folderName", "n", "", "Grafana dir name to which push dashboards")
	pullDashboardsCmd.PersistentFlags().IntP("folderId", "i", -1, "Grafana dir ID from which pull dashboards")
	pullDashboardsCmd.PersistentFlags().StringP("folderName", "n", "", "Grafana dir name from which pull dashboards")
	pullDashboardsCmd.PersistentFlags().StringP("tag", "t", "", "Dashboard tag to p")
	pullDashboardsCmd.PersistentFlags().BoolP("fetchVersions", "v", false, "Also fetch Dashboard versions for git purposes")
}
