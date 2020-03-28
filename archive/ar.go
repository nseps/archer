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
	"io"
	"os"
	"path/filepath"

	"github.com/thegrumpylion/ar"
)

func init() {
	RegisterArchiver("ar", arArchiver{})
}

type arArchiver struct{}

func (a arArchiver) Pack(src string, target io.Writer) error {
	arw := ar.NewWriter(target)

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		// ignore root dir
		if filepath.Clean(src) == path {
			return nil
		}

		mode := FileMode{info.Mode()}

		link := ""
		if mode.IsSymlink() {
			link, err = os.Readlink(path)
			if err != nil {
				return err
			}
		}

		hdr, err := ar.FileInfoHeader(info, link)
		if err != nil {
			return err
		}
		// path relative to src
		hdr.Name, err = filepath.Rel(filepath.Clean(src), path)
		if err != nil {
			return err
		}

		// is dir. just write header
		if info.IsDir() {
			hdr.Size = 0
			return arw.WriteHeader(hdr)
		}
		arw.WriteHeader(hdr)

		if mode.IsRegular() {
			body, err := os.Open(path)
			if err != nil {
				return err
			}
			_, err = io.Copy(arw, body)
			if err != nil {
				return err
			}
			body.Close()
		}
		return nil
	})

	return err
}

func (a arArchiver) Unpack(src io.Reader, target string) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	arr := ar.NewReader(src)

	for {
		hdr, err := arr.Next()
		// done
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		path := filepath.Join(target, hdr.Name)
		mode := FileMode{hdr.FileInfo().Mode()}

		switch {
		case mode.IsRegular():
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}
			f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, arr); err != nil {
				return err
			}
		}
	}
}
