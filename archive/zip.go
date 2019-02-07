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
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func init() {
	RegisterArchiver("zip", zipArchiver{})
}

type zipArchiver struct{}

func (a zipArchiver) Pack(src string, target io.Writer) error {
	zipw := zip.NewWriter(target)

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

		hdr, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		// NOTE: The sane thing to do is check for target file and dectet if
		// compressed or not and use header.Method = zip.Deflate to compress.
		// But this breaks our tar-like behaviour so we just use zip.Store for
		// evrything. Maybe use it like .zip.gz??
		hdr.Method = zip.Store
		// path relative to src
		hdr.Name, err = filepath.Rel(filepath.Clean(src), path)
		if err != nil {
			return err
		}

		// is dir. just write header
		if info.IsDir() {
			zipw.CreateHeader(hdr)
			return nil
		}
		w, err := zipw.CreateHeader(hdr)
		if err != nil {
			return err
		}

		if mode.IsRegular() {
			body, err := os.Open(path)
			if err != nil {
				return err
			}
			_, err = io.Copy(w, body)
			if err != nil {
				return err
			}
			body.Close()
		}

		if mode.IsSymlink() {
			_, err := w.Write([]byte(link))
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return zipw.Close()
}

func (a zipArchiver) Unpack(src io.Reader, target string) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}
	// NOTE: either we change the Unpack abstraction to include src
	// file size or we read all in memory to count. The later for now..
	data, err := ioutil.ReadAll(src)
	if err != nil {
		return err
	}
	r := bytes.NewReader(data)

	zipr, err := zip.NewReader(r, r.Size())
	if err != nil {
		return err
	}

	for _, zf := range zipr.File {

		mode := FileMode{zf.FileInfo().Mode()}
		path := filepath.Join(target, zf.Name)

		switch {
		case mode.IsDir():
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case mode.IsRegular():
			f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, zf.FileInfo().Mode())
			if err != nil {
				return err
			}
			rc, err := zf.Open()
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, rc); err != nil {
				return err
			}
			rc.Close()
		case mode.IsSymlink():
			rc, err := zf.Open()
			if err != nil {
				return err
			}
			ldata, err := ioutil.ReadAll(rc)
			if err != nil {
				return err
			}
			rc.Close()
			if err := os.Symlink(string(ldata), path); err != nil {
				return err
			}
		}
	}
	return nil
}
