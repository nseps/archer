package archive

import (
	"fmt"
	"io"
)

type Archiver interface {
	Pack(src string, target io.Writer) error
	Unpack(src io.Reader, target string) error
}

var archivers = map[string]Archiver{}

type AlreadyExistsError string

func (e AlreadyExistsError) Error() string {
	return fmt.Sprintf("Archiver with name \"%s\" aready registered", string(e))
}

type DoesNotExistError string

func (e DoesNotExistError) Error() string {
	return fmt.Sprintf("Archiver with name \"%s\" does not exist", string(e))
}

func RegisterArchiver(name string, ar Archiver) error {
	if _, inMap := archivers[name]; inMap {
		return AlreadyExistsError(name)
	}
	archivers[name] = ar
	return nil
}

func GetArchiver(name string) (Archiver, error) {
	if ar, inMap := archivers[name]; inMap {
		return ar, nil
	}
	return nil, DoesNotExistError(name)
}

func ListSupported() []string {
	list := []string{}
	for ar := range archivers {
		list = append(list, ar)
	}
	return list
}
