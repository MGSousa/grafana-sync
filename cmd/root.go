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
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/mpostument/grafana-sync/grafana"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile       string
	customHeaders map[string]string

	gs grafana.Grafana
)

var rootCmd = &cobra.Command{
	Use:     "grafana-sync",
	Short:   "Root command for grafana interaction",
	Long:    `Root command for grafana interaction.`,
	Version: "1.5.0",
}

// Execute parse all defined flags
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Errorln(err)
		os.Exit(1)
	}
}

// initialize Grafana client
func initGrafana(cmd *cobra.Command, args []string) {
	apiKey := viper.GetString("apikey")
	url, _ := cmd.Flags().GetString("url")
	directory, _ := cmd.Flags().GetString("directory")

	gs = grafana.New(url, apiKey, directory)
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.grafana-sync.yaml)")
	rootCmd.PersistentFlags().StringP("url", "u", "http://localhost:3000", "Grafana URI")
	rootCmd.PersistentFlags().StringP("directory", "d", ".", "Local directory where to save dashboards, datasources, etc.")
	rootCmd.PersistentFlags().StringP("apikey", "a", "", "Grafana ServiceAccount API Key")
	rootCmd.PersistentFlags().StringToStringVar(&customHeaders, "customHeaders", map[string]string{}, "Key-value pairs of custom HTTP headers (key1=value1,key2=value2)")

	if err := viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey")); err != nil {
		log.Errorln(err)
	}

	if err := viper.BindPFlag("customHeaders", rootCmd.PersistentFlags().Lookup("customHeaders")); err != nil {
		log.Warningln(err)
	}

	rootCmd.AddCommand(
		pullDashboardsCmd, pushDashboardsCmd,
		pullFoldersCmd, pushFoldersCmd,
		pullNotificationsCmd, pushNotificationsCmd,
		pullDataSourcesCmd, pushDataSourcesCmd,
	)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".grafana-sync" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".grafana-sync")
	}

	// read environment variables that match flags
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	grafana.InitHttpClient(viper.GetStringMapString("customHeaders"))
}
