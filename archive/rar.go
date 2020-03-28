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

package archive

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/nwaples/rardecode"
)

func init() {
	RegisterArchiver("rar", rarArchiver{})
}

type rarArchiver struct{}

func (a rarArchiver) Pack(src string, target io.Writer) error {
	return errors.New("rar has no packing implementation")
}

func (a rarArchiver) Unpack(src io.Reader, target string) error {
	if err := os.Mkdir(target, 0755); err != nil {
		return err
	}
	// the second param is password. maybe if we ever support
	// archiver specific options
	rarr, err := rardecode.NewReader(src, "")
	if err != nil {
		return err
	}

	for {
		hdr, err := rarr.Next()
		// done
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		path := filepath.Join(target, hdr.Name)
		mode := FileMode{hdr.Mode()}

		switch {
		case mode.IsDir():
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case mode.IsRegular():
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, hdr.Mode())
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, rarr); err != nil {
				return err
			}
		}

	}

}
