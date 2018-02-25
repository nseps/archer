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
	"io"
	"io/ioutil"

	"github.com/ulikunitz/xz"
)

func init() {
	RegisterCompressor("xz", xzCompressor{})
}

type xzCompressor struct{}

func (comp xzCompressor) Compress(w io.Writer) (io.WriteCloser, error) {
	return xz.NewWriter(w)
}

func (comp xzCompressor) Decompress(r io.Reader) (io.ReadCloser, error) {
	rdr, err := xz.NewReader(r)
	return ioutil.NopCloser(rdr), err
}
