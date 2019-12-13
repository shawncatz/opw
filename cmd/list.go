/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

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

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/shawncatz/opw/wrapper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list all available logins",
	Long:  "list all available logins",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := wrapper.NewClient(cfg)
		if err != nil {
			logrus.Errorf("error getting client: %s", err)
			return
		}

		err = client.SignIn()
		if err != nil {
			logrus.Errorf("error getting signin: %s", err)
			return
		}

		items, err := client.List()
		if err != nil {
			logrus.Errorf("error getting list: %s", err)
			return
		}

		for _, i := range items {
			fmt.Printf("%-28s %-3.3s %s\n", i.UUID, i.TemplateUUID, i.Overview.Title)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
