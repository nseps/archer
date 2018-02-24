package archive

import (
	"io"
	"os"
	"path/filepath"

	cpio "github.com/nseps/go-cpio"
)

func init() {
	RegisterArchiver("cpio", cpioArchiver{})
}

type cpioArchiver struct{}

func (t cpioArchiver) Pack(src string, target io.Writer) error {
	cpiow := cpio.NewWriter(target)

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

		hdr, err := cpio.FileInfoHeader(info, link)
		if err != nil {
			return err
		}
		// path relative to src
		hdr.Name, err = filepath.Rel(filepath.Clean(src), path)
		if err != nil {
			return err
		}

		// is dir. just write header
		if mode.IsDir() {
			hdr.Size = 0
			return cpiow.WriteHeader(hdr)
		}
		cpiow.WriteHeader(hdr)

		if mode.IsRegular() {
			body, err := os.Open(path)
			if err != nil {
				return err
			}
			_, err = io.Copy(cpiow, body)
			if err != nil {
				return err
			}
			body.Close()
		}

		if mode.IsSymlink() {
			_, err := cpiow.Write([]byte(link))
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return cpiow.Close()
}

func (t cpioArchiver) Unpack(src io.Reader, target string) error {
	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	cpior := cpio.NewReader(src)

	for {
		hdr, err := cpior.Next()
		// done
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		path := filepath.Join(target, hdr.Name)

		switch {
		case hdr.Mode.IsDir():
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case hdr.Mode.IsRegular():
			f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, cpior); err != nil {
				return err
			}
		case hdr.Mode.IsSymlink():
			if err := os.Symlink(hdr.Linkname, path); err != nil {
				return err
			}
		}
	}
}
