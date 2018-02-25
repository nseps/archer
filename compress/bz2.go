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

	"github.com/dsnet/compress/bzip2"
)

func init() {
	RegisterCompressor("bz2", bz2Compressor{})
}

type bz2Compressor struct{}

func (comp bz2Compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	return bzip2.NewWriter(w, &bzip2.WriterConfig{
		Level: bzip2.DefaultCompression,
	})
}

func (comp bz2Compressor) Decompress(r io.Reader) (io.ReadCloser, error) {
	// second param is conf and is not used internaly by the lib
	return bzip2.NewReader(r, nil)
}
