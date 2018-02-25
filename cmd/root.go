// Copyright © 2018 Nikolas Sepos <nikolas.sepos@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "archer",
	Short: "An archiver fused with magic",
	Long:  ``,
	Run:   rootCmdRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("completionFile", "f", "", "Completion out file")
	rootCmd.Flags().String("shell", "bash", "Target shell name for completion")
}

func rootCmdRun(cmd *cobra.Command, args []string) {
	compf, err := cmd.Flags().GetString("completionFile")
	dieOnErr(err)

	if compf == "" {
		cmd.Usage()
		return
	}
	shl, err := cmd.Flags().GetString("shell")
	dieOnErr(err)

	switch shl {
	case "bash":
		err = cmd.GenBashCompletionFile(compf)
	case "zsh":
		err = cmd.GenZshCompletionFile(compf)
	default:
		dieOnErr(fmt.Errorf("shell not supported: %s", shl))
	}
	dieOnErr(err)
}

func dieOnErr(err error) {
	if err != nil {
		log.Fatalf("Fatal Error: %v\n", err)
	}
}
