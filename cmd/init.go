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
	"io/ioutil"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var forceFlag bool

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "configure your system to use opw",
	Long:  "configure your system to use opw",
	Run: func(cmd *cobra.Command, args []string) {
		home, err := homedir.Dir()
		if err != nil {
			logrus.Errorf("error getting home directory: %s", err)
		}

		config := `subdomain: name # subdomain of your 1password account
email: me@example.com # your email address
cache: ` + home + `/.opw.cache # where you want to store your cached session
aliases: # set of aliases (key: value)
  nickname: uuid # find uuid by running 'opw list'

# For both passphrase and secret, you specify how to obtain the secret
# rather than specifying it directly.
# if the string starts with 'file:' then it will load from the file
# if the string starts with 'keychain:' then it will load from the keychain
# the keyring library supports MacOS, Linux (d-bus), and Windows
# for keychain, the value is specified as 'keychain:service:account'
# using the keychain is more secure

# The secret currently isn't used
#secret: file:/path/to/file/containing/secret
secret: keychain:opw-secret:subdomain

#passphrase: file:/path/to/file/containing/passphrase
passphrase: keychain:opw-passphrase:subdomain # keychain entry
`
		if _, err := os.Stat(home + "/.opw.yaml"); os.IsNotExist(err) || forceFlag {
			err = ioutil.WriteFile(home+"/.opw.yaml", []byte(config), 0600)
			if err != nil {
				logrus.Errorf("error writing config file: %s", err)
			}
		} else {
			logrus.Warnf("config file exists: %s", home+"/.opw.yaml")
		}

		logrus.Infof("use your package manager to install the 1password cli")
		logrus.Infof("or see: https://1password.com/downloads/command-line/")
		logrus.Infof("after installation, signin using: ")
		logrus.Infof("'op signin subdmain.1password.org email secret'")
		logrus.Infof("after you've done that, opw will handle managing a session for you")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	initCmd.Flags().BoolVarP(&forceFlag, "force", "f", false, "overwrite the configuration file if it exists")
}
