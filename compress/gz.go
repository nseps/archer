package compress

import (
	"compress/gzip"
	"io"
)

func init() {
	RegisterCompressor("gz", gzCompressor{})
}

type gzCompressor struct{}

func (comp gzCompressor) Compress(w io.Writer) (io.WriteCloser, error) {
	return gzip.NewWriter(w), nil
}

func (comp gzCompressor) Decompress(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}
