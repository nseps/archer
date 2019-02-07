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

package match

import (
	"bufio"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/thegrumpylion/archer/archive"
	"github.com/thegrumpylion/archer/compress"

	"gopkg.in/h2non/filetype.v1"
)

// Compression will wrap reader r with the compresor if the file
// signature match. We get some overhead with wrapping with buffered
// readers but we gain magic convenience
func Compression(r io.Reader) (compress.Compressor, io.Reader, error) {
	// we wrap with buffered reader to do Peek
	br := bufio.NewReader(r)
	// 262 comes from the filetype lib
	hdr, err := br.Peek(262)
	if err != nil {
		return nil, nil, err
	}

	ty, err := filetype.Match(hdr)
	if err != nil {
		return nil, nil, err
	}
	if ty == filetype.Unknown {
		return nil, nil, errors.New("Cannot detect file type")
	}

	comp, err := compress.GetCompressor(ty.Extension)
	// we don't support this kind of compression
	if err != nil {
		return nil, br, nil
	}
	return comp, br, nil
}

// Archive will try to match file signature with registered archivers
// Same concept with buffered reader as above
func Archive(r io.Reader) (archive.Archiver, io.Reader, error) {
	// we wrap with buffered reader to do Peek
	br := bufio.NewReader(r)
	// 262 comes from the filetype lib
	hdr, err := br.Peek(262)
	if err != nil {
		return nil, nil, err
	}

	ty, err := filetype.Match(hdr)
	if err != nil {
		return nil, nil, err
	}
	if ty == filetype.Unknown {
		return nil, nil, errors.New("Cannot detect file type")
	}

	arch, err := archive.GetArchiver(ty.Extension)

	return arch, br, err
}

// FileName will try to infer the archiver and/or coressor from the filename
func FileName(fn string) (archive.Archiver, compress.Compressor, error) {
	// get file extention
	ext := filepath.Ext(fn)
	if ext == "" {
		return nil, nil, errors.New("No extention found")
	}
	// check for compressor. also remove leading . from ext
	comp, err := compress.GetCompressor(ext[1:])
	// compresor found
	if err == nil {
		// remove what we found
		fn = strings.TrimSuffix(fn, ext)
		// repeat for archive
		ext = filepath.Ext(fn)
		// return error if file was name.comp
		if ext == "" {
			return nil, nil, errors.New("No extention found")
		}
	}
	arch, err := archive.GetArchiver(ext[1:])

	return arch, comp, err
}
