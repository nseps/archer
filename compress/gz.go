package compress

import (
	"compress/gzip"
	"io"
)

func init() {
	RegisterCompressor("gz", gzCompressor{})
}

type gzCompressor struct{}

func (gz gzCompressor) Compress(w io.Writer) io.WriteCloser {
	return gzip.NewWriter(w)
}

func (gz gzCompressor) Decompress(r io.Reader) (io.ReadCloser, error) {
	return gzip.NewReader(r)
}
