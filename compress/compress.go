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

package compress

import (
	"fmt"
	"io"
)

// Compressor interface should cover most cases
type Compressor interface {
	Compress(w io.Writer) (io.WriteCloser, error)
	Decompress(r io.Reader) (io.ReadCloser, error)
}

var compressors = map[string]Compressor{}

type AlreadyExistsError string

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("Ccmpressor with name \"%s\" aready registered", string(e))
}

type DoesNotExistError string

func (e DoesNotExistError) Error() string {
	return fmt.Sprintf("compressor with name \"%s\" does not exist", string(e))
}

func RegisterCompressor(name string, comp Compressor) error {
	if _, inMap := compressors[name]; inMap {
		return AlreadyExistsError(name)
	}
	compressors[name] = comp
	return nil
}

func GetCompressor(name string) (Compressor, error) {
	if comp, inMap := compressors[name]; inMap {
		return comp, nil
	}
	return nil, DoesNotExistError(name)
}

func ListSupported() []string {
	list := []string{}
	for comp := range compressors {
		list = append(list, comp)
	}
	return list
}
