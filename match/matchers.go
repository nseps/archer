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
	"encoding/binary"

	filetype "gopkg.in/h2non/filetype.v1"
)

func init() {
	filetype.AddMatcher(cpioType, cpioMatcher)
	filetype.AddMatcher(lz4Type, lz4Matcher)
}

var cpioType = filetype.NewType("cpio", "application/x-cpio")

func cpioMatcher(buf []byte) bool {
	// 0707 in ASCII
	return len(buf) > 1 && buf[0] == 0x30 && buf[1] == 0x37 &&
		buf[2] == 0x30 && buf[3] == 0x37
}

var lz4Type = filetype.NewType("lz4", "application/x-lz4")

func lz4Matcher(buf []byte) bool {
	// 4 Bytes, Little endian format. Value : 0x184D2204
	// https://github.com/lz4/lz4/wiki/lz4_Frame_format.md
	if len(buf) < 1 {
		return false
	}
	val := binary.LittleEndian.Uint32(buf[0:])
	return val == 0x184D2204
}
