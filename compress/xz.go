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
