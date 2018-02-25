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
	"fmt"
	"io"
)

// Archiver interface may be altered for Unpack to take a
// writer intsted of a string
type Archiver interface {
	Pack(src string, target io.Writer) error
	Unpack(src io.Reader, target string) error
}

var archivers = map[string]Archiver{}

type AlreadyExistsError string

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("archiver with name \"%s\" aready registered", string(e))
}

type DoesNotExistError string

func (e DoesNotExistError) Error() string {
	return fmt.Sprintf("archiver with name \"%s\" does not exist", string(e))
}

func RegisterArchiver(name string, ar Archiver) error {
	if _, inMap := archivers[name]; inMap {
		return AlreadyExistsError(name)
	}
	archivers[name] = ar
	return nil
}

func GetArchiver(name string) (Archiver, error) {
	if ar, inMap := archivers[name]; inMap {
		return ar, nil
	}
	return nil, DoesNotExistError(name)
}

func ListSupported() []string {
	list := []string{}
	for ar := range archivers {
		list = append(list, ar)
	}
	return list
}
