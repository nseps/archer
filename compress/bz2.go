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
