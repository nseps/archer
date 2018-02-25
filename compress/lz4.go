package compress

import (
	"io"
	"io/ioutil"

	"github.com/pierrec/lz4"
)

func init() {
	RegisterCompressor("lz4", lz4Compressor{})
}

type lz4Compressor struct{}

func (comp lz4Compressor) Compress(w io.Writer) (io.WriteCloser, error) {
	return lz4.NewWriter(w), nil
}

func (comp lz4Compressor) Decompress(r io.Reader) (io.ReadCloser, error) {
	return ioutil.NopCloser(lz4.NewReader(r)), nil
}
