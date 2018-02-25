package archive

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

func init() {
	RegisterArchiver("tar", tarArchiver{})
}

type tarArchiver struct{}

func (a tarArchiver) Pack(src string, target io.Writer) error {
	tarw := tar.NewWriter(target)

	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		// ignore root dir
		if filepath.Clean(src) == path {
			return nil
		}

		mode := FileMode{info.Mode()}

		link := ""
		if mode.IsSymlink() {
			link, err = os.Readlink(path)
			if err != nil {
				return err
			}
		}

		hdr, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return err
		}
		// path relative to src
		hdr.Name, err = filepath.Rel(filepath.Clean(src), path)
		if err != nil {
			return err
		}

		// is dir. just write header
		if info.IsDir() {
			hdr.Size = 0
			return tarw.WriteHeader(hdr)
		}
		tarw.WriteHeader(hdr)

		if mode.IsRegular() {
			body, err := os.Open(path)
			if err != nil {
				return err
			}
			_, err = io.Copy(tarw, body)
			if err != nil {
				return err
			}
			body.Close()
		}
		return nil
	})

	if err != nil {
		return err
	}

	return tarw.Close()
}

func (a tarArchiver) Unpack(src io.Reader, target string) error {
	if err := os.Mkdir(target, 0755); err != nil {
		return err
	}

	tarr := tar.NewReader(src)

	for {
		hdr, err := tarr.Next()
		// done
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		path := filepath.Join(target, hdr.Name)
		mode := FileMode{hdr.FileInfo().Mode()}

		switch {
		case mode.IsDir():
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case mode.IsRegular():
			f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tarr); err != nil {
				return err
			}
		case mode.IsSymlink():
			if err := os.Symlink(hdr.Linkname, path); err != nil {
				return err
			}
		}
	}
}
