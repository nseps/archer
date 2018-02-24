package compress

import (
	"fmt"
	"io"
)

type AlreadyExistsError string

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("Compressor with name \"%s\" aready registered", string(e))
}

type DoesNotExistError string

func (e DoesNotExistError) Error() string {
	return fmt.Sprintf("Compressor with name \"%s\" does not exist", string(e))
}

var compressors = map[string]Compressor{}

type Compressor interface {
	Compress(w io.Writer) io.WriteCloser
	Decompress(r io.Reader) (io.ReadCloser, error)
}

func RegisterCompressor(name string, comp Compressor) error {
	if _, inMap := compressors[name]; inMap {
		return AlreadyExistsError(name)
	}
	compressors[name] = comp
	return nil
}

func GetCompressor(name string) (Compressor, error) {
	if comp, inMap := compressors[name]; inMap {
		return comp, nil
	}
	return nil, DoesNotExistError(name)
}

func ListSupported() []string {
	list := []string{}
	for comp := range compressors {
		list = append(list, comp)
	}
	return list
}
