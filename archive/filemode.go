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
	"os"
)

// FileMode wraps an os.FileMode
type FileMode struct {
	os.FileMode
}

// IsSymlink returns true if file is a sympolic link
func (fm *FileMode) IsSymlink() bool {
	return fm.FileMode&os.ModeSymlink != 0
}

// IsBlockDevice returns true if file is a block device
func (fm *FileMode) IsBlockDevice() bool {
	if fm.FileMode&os.ModeDevice != 0 {
		// if not CharDev
		if !(fm.FileMode&os.ModeCharDevice != 0) {
			return true
		}
	}
	return false
}

// IsCharDevice returns true if file is a charachter device
func (fm *FileMode) IsCharDevice() bool {
	if fm.FileMode&os.ModeDevice != 0 {
		// if is CharDev
		if fm.FileMode&os.ModeCharDevice != 0 {
			return true
		}
	}
	return false
}

// IsNamedPipe returns true if file is a named pipe
func (fm *FileMode) IsNamedPipe() bool {
	return fm.FileMode&os.ModeNamedPipe != 0
}

// IsSocket returns true if file is a socket
func (fm *FileMode) IsSocket() bool {
	return fm.FileMode&os.ModeSocket != 0
}

// HasSetUID returns true if file has setuid bit set
func (fm *FileMode) HasSetUID() bool {
	return fm.FileMode&os.ModeSetuid != 0
}

// HasSetGID returns true if file has setgid bit set
func (fm *FileMode) HasSetGID() bool {
	return fm.FileMode&os.ModeSetgid != 0
}

// HasSticky returns true if file has sticky bit set
func (fm *FileMode) HasSticky() bool {
	return fm.FileMode&os.ModeSticky != 0
}
