package androidbinary

import (
	"encoding/binary"
	"io"
	"unicode/utf16"
)

func readUTF16(sr *io.SectionReader) (string, error) {
	// read lenth of string
	var size int
	var first, second uint16
	if err := binary.Read(sr, binary.LittleEndian, &first); err != nil {
		return "", err
	}
	if (first & 0x8000) != 0 {
		if err := binary.Read(sr, binary.LittleEndian, &second); err != nil {
			return "", err
		}
		size = (int(first&0x7FFF) << 16) + int(second)
	} else {
		size = int(first)
	}

	// read string value
	buf := make([]uint16, size)
	if err := binary.Read(sr, binary.LittleEndian, buf); err != nil {
		return "", err
	}
	return string(utf16.Decode(buf)), nil
}
