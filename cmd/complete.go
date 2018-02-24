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

	"github.com/spf13/cobra"
)

// completeCmd represents the complete command
var completeCmd = &cobra.Command{
	Use:   "complete [--shell bash|zsh] OUTFILE",
	Short: "Generate completion for bash or zsh",
	Run:   completeRun,
	Args:  cobra.ExactArgs(1),
}

func init() {
	rootCmd.AddCommand(completeCmd)
	completeCmd.Flags().String("shell", "bash", "Target shell name")
}

func completeRun(cmd *cobra.Command, args []string) {

	shl, err := cmd.Flags().GetString("shell")
	dieOnErr(err)

	switch shl {
	case "bash":
		err = rootCmd.GenBashCompletionFile(args[0])
	case "zsh":
		err = rootCmd.GenZshCompletionFile(args[0])
	default:
		dieOnErr(fmt.Errorf("Shell not supported: %s", shl))
	}
	dieOnErr(err)
}
