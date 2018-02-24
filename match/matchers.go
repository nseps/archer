package match

import filetype "gopkg.in/h2non/filetype.v1"

func init() {
	filetype.AddMatcher(cpioType, cpioMatcher)
}

var cpioType = filetype.NewType("cpio", "application/x-cpio")

func cpioMatcher(buf []byte) bool {
	// 0707 in ASCII
	return len(buf) > 1 && buf[0] == 0x30 && buf[1] == 0x37 &&
		buf[2] == 0x30 && buf[3] == 0x37
}
