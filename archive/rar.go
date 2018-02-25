package archive

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/nwaples/rardecode"
)

func init() {
	RegisterArchiver("rar", rarArchiver{})
}

type rarArchiver struct{}

func (a rarArchiver) Pack(src string, target io.Writer) error {
	return errors.New("rar has no packing implementation")
}

func (a rarArchiver) Unpack(src io.Reader, target string) error {
	if err := os.Mkdir(target, 0755); err != nil {
		return err
	}
	// the second param is password. maybe if we ever support
	// archiver specific options
	rarr, err := rardecode.NewReader(src, "")
	if err != nil {
		return err
	}

	for {
		hdr, err := rarr.Next()
		// done
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		path := filepath.Join(target, hdr.Name)
		mode := FileMode{hdr.Mode()}

		switch {
		case mode.IsDir():
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case mode.IsRegular():
			f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, hdr.Mode())
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, rarr); err != nil {
				return err
			}
		}

	}

}
