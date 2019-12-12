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

	"github.com/shawncatz/opw/opw"
)

var passwordFlag bool

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get UUID",
	Short: "get password from 1password",
	Long:  "get password from 1password",
	Run: func(cmd *cobra.Command, args []string) {
		uuid := args[0]
		alias := cfg.Aliases[uuid]
		if alias != "" {
			uuid = alias
		}

		item, err := opw.GetItem(uuid)
		if err != nil {
			logrus.Errorf("error getting item: %s", err)
			return
		}

		if passwordFlag {
			fmt.Println(item.Password())
		} else {
			fmt.Printf("%s / %s\n", item.Username(), item.Password())
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCmd.Flags().BoolVarP(&passwordFlag, "password", "p", false, "only print password")
}
