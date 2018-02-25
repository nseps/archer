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

	cpio "github.com/nseps/go-cpio"
)

func init() {
	RegisterArchiver("cpio", cpioArchiver{})
}

type cpioArchiver struct{}

func (a cpioArchiver) Pack(src string, target io.Writer) error {
	cpiow := cpio.NewWriter(target)

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

		hdr, err := cpio.FileInfoHeader(info, link)
		if err != nil {
			return err
		}
		// path relative to src
		hdr.Name, err = filepath.Rel(filepath.Clean(src), path)
		if err != nil {
			return err
		}

		// is dir. just write header
		if mode.IsDir() {
			hdr.Size = 0
			return cpiow.WriteHeader(hdr)
		}
		cpiow.WriteHeader(hdr)

		if mode.IsRegular() {
			body, err := os.Open(path)
			if err != nil {
				return err
			}
			_, err = io.Copy(cpiow, body)
			if err != nil {
				return err
			}
			body.Close()
		}

		if mode.IsSymlink() {
			_, err := cpiow.Write([]byte(link))
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return cpiow.Close()
}

func (a cpioArchiver) Unpack(src io.Reader, target string) error {
	if err := os.Mkdir(target, 0755); err != nil {
		return err
	}

	cpior := cpio.NewReader(src)

	for {
		hdr, err := cpior.Next()
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
		case mode.IsDir():
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case mode.IsRegular():
			f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, hdr.FileInfo().Mode())
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, cpior); err != nil {
				return err
			}
		case mode.IsSymlink():
			if err := os.Symlink(hdr.Linkname, path); err != nil {
				return err
			}
		}
	}
}
