package archive

import (
	"os"
)

// FileMode wraps an os.FileMode
type FileMode struct {
	os.FileMode
}

// IsSymlink returns true if file is a sympolic link
func (fm *FileMode) IsSymlink() bool {
	return fm.FileMode&os.ModeSymlink != 0
}

// IsBlockDevice returns true if file is a block device
func (fm *FileMode) IsBlockDevice() bool {
	if fm.FileMode&os.ModeDevice != 0 {
		// if not CharDev
		if !(fm.FileMode&os.ModeCharDevice != 0) {
			return true
		}
	}
	return false
}

// IsCharDevice returns true if file is a charachter device
func (fm *FileMode) IsCharDevice() bool {
	if fm.FileMode&os.ModeDevice != 0 {
		// if is CharDev
		if fm.FileMode&os.ModeCharDevice != 0 {
			return true
		}
	}
	return false
}

// IsNamedPipe returns true if file is a named pipe
func (fm *FileMode) IsNamedPipe() bool {
	return fm.FileMode&os.ModeNamedPipe != 0
}

// IsSocket returns true if file is a socket
func (fm *FileMode) IsSocket() bool {
	return fm.FileMode&os.ModeSocket != 0
}

// HasSetUID returns true if file has setuid bit set
func (fm *FileMode) HasSetUID() bool {
	return fm.FileMode&os.ModeSetuid != 0
}

// HasSetGID returns true if file has setgid bit set
func (fm *FileMode) HasSetGID() bool {
	return fm.FileMode&os.ModeSetgid != 0
}

// HasSticky returns true if file has sticky bit set
func (fm *FileMode) HasSticky() bool {
	return fm.FileMode&os.ModeSticky != 0
}
