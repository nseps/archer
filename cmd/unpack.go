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
	"os"
	"strings"

	"github.com/nseps/archer/match"
	"github.com/spf13/cobra"
)

// unpackCmd represents the unpack command
var unpackCmd = &cobra.Command{
	Use:   "unpack SRCFILE [TARGETDIR]",
	Short: "Unpack an archive to a directory",
	Long: `Unpack an archive to a directory using manna to detect the compression and archive type.
The file name extentions are not taken into consideration.

If TARGETDIR is not specified, SRCFILE name without extension will be used`,
	Run:  unpackRun,
	Args: cobra.RangeArgs(1, 2),
}

func init() {
	rootCmd.AddCommand(unpackCmd)
}

func unpackRun(cmd *cobra.Command, args []string) {
	// read source file
	f, err := os.Open(args[0])
	dieOnErr(err)
	defer f.Close()
	// detect compression if any
	in, err := match.Compression(f)
	dieOnErr(err)
	defer in.Close()
	// detect archive type
	ar, rdr, err := match.Archive(in)
	dieOnErr(err)
	// check for target
	var trgt string
	if len(args) == 2 {
		trgt = args[1]
	} else {
		parts := strings.Split(args[0], ".")
		if parts[0] == args[0] {
			dieOnErr(fmt.Errorf("cannot find a suitable directory name for given file"))
		}
		trgt = parts[0]
	}
	// do your thing
	err = ar.Unpack(rdr, trgt)
	dieOnErr(err)
}
