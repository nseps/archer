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
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nseps/archer/match"

	"github.com/nseps/archer/archive"
	"github.com/nseps/archer/compress"
	"github.com/spf13/cobra"
)

var packUse = "pack [--type|-t " + strings.Join(archive.ListSupported(), "|") + "] [--compression|-c " + strings.Join(compress.ListSupported(), "|") + "] SRCDIR TARGETFILE"

// packCmd represents the pack command
var packCmd = &cobra.Command{
	Use:   packUse,
	Short: "Pack a directory to an archive",
	Long:  ``,
	Run:   packRun,
	Args:  cobra.RangeArgs(1, 2),
}

func init() {
	rootCmd.AddCommand(packCmd)

	packCmd.Flags().StringP("type", "t", "", "Archive type")
	packCmd.Flags().StringP("compression", "c", "", "Compression type")
}

func packRun(cmd *cobra.Command, args []string) {
	// Get flags
	t, err := cmd.Flags().GetString("type")
	dieOnErr(err)
	c, err := cmd.Flags().GetString("compression")
	dieOnErr(err)

	var ar archive.Archiver
	var comp compress.Compressor
	var trgt string
	if t == "" {
		// check for target name
		if len(args) != 2 {
			dieOnErr(fmt.Errorf("a filename with proper extension must be provided if no explicit type is set "))
		}
		trgt = args[1]
		// try to guess from target filename
		ar, comp, err = match.FileName(args[1])
		if err != nil {
			dieOnErr(fmt.Errorf("could not infer archive and/or compression type from filename, %v", err))
		}
	} else {
		// target name is src dir basename
		trgt = filepath.Base(filepath.Clean(args[0]))
		// Find archiver
		ar, err = archive.GetArchiver(t)
		dieOnErr(err)
		// append file ext
		trgt = fmt.Sprintf("%s.%s", trgt, t)
		if c != "" {
			// find compressor
			comp, err = compress.GetCompressor(c)
			dieOnErr(err)
			// append file ext
			trgt = fmt.Sprintf("%s.%s", trgt, c)
		}
	}

	// Create target file
	f, err := os.Create(trgt)
	dieOnErr(err)
	defer f.Close()

	// Wrap if compresion requested
	var out io.WriteCloser
	if comp != nil {
		out = comp.Compress(f)
		defer out.Close()
	} else {
		out = f
	}

	// Ship it!
	err = ar.Pack(args[0], out)
	dieOnErr(err)
}
