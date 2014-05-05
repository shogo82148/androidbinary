package androidbinary

import (
	"encoding/binary"
	"io"
	"os"
	"unicode/utf16"
)

const (
	RES_NULL_TYPE        = 0x0000
	RES_STRING_POOL_TYPE = 0x0001
	RES_TABLE_TYPE       = 0x0002
	RES_XML_TYPE         = 0x0003

	// Chunk types in RES_XML_TYPE
	RES_XML_FIRST_CHUNK_TYPE     = 0x0100
	RES_XML_START_NAMESPACE_TYPE = 0x0100
	RES_XML_END_NAMESPACE_TYPE   = 0x0101
	RES_XML_START_ELEMENT_TYPE   = 0x0102
	RES_XML_END_ELEMENT_TYPE     = 0x0103
	RES_XML_CDATA_TYPE           = 0x0104
	RES_XML_LAST_CHUNK_TYPE      = 0x017f

	// This contains a uint32_t array mapping strings in the string
	// pool back to resource identifiers.  It is optional.
	RES_XML_RESOURCE_MAP_TYPE = 0x0180

	// Chunk types in RES_TABLE_TYPE
	RES_TABLE_PACKAGE_TYPE   = 0x0200
	RES_TABLE_TYPE_TYPE      = 0x0201
	RES_TABLE_TYPE_SPEC_TYPE = 0x0202
)

type ResChunkHeader struct {
	Type       uint16
	HeaderSize uint16
	Size       uint32
}

const SORTED_FLAG = 1 << 0
const UTF8_FLAG = 1 << 8

type ResStringPoolHeader struct {
	Header      ResChunkHeader
	StringCount uint32
	StyleCount  uint32
	Flags       uint32
	StringStart uint32
	StylesStart uint32
}

type ResStringPool struct {
	Header  ResStringPoolHeader
	Strings []string
	Styles  []string
}

const NilResStringPoolRef = ResStringPoolRef(0xFFFFFFFF)

type ResStringPoolRef uint32

func readStringPool(sr *io.SectionReader) (*ResStringPool, error) {
	sp := new(ResStringPool)
	if err := binary.Read(sr, binary.LittleEndian, &sp.Header); err != nil {
		return nil, err
	}

	stringStarts := make([]uint32, sp.Header.StringCount)
	if err := binary.Read(sr, binary.LittleEndian, stringStarts); err != nil {
		return nil, err
	}

	styleStarts := make([]uint32, sp.Header.StyleCount)
	if err := binary.Read(sr, binary.LittleEndian, styleStarts); err != nil {
		return nil, err
	}

	sp.Strings = make([]string, sp.Header.StringCount)
	for i, start := range stringStarts {
		var str string
		var err error
		sr.Seek(int64(sp.Header.StringStart+start), os.SEEK_SET)
		if (sp.Header.Flags & UTF8_FLAG) == 0 {
			str, err = readUTF16(sr)
		} else {
			str, err = readUTF8(sr)
		}
		if err != nil {
			return nil, err
		}
		sp.Strings[i] = str
	}

	sp.Styles = make([]string, sp.Header.StyleCount)
	for i, start := range styleStarts {
		var str string
		var err error
		sr.Seek(int64(sp.Header.StylesStart+start), os.SEEK_SET)
		if (sp.Header.Flags & UTF8_FLAG) == 0 {
			str, err = readUTF16(sr)
		} else {
			str, err = readUTF8(sr)
		}
		if err != nil {
			return nil, err
		}
		sp.Styles[i] = str
	}

	return sp, nil
}

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

func readUTF8(sr *io.SectionReader) (string, error) {
	// read lenth of string
	var size int
	var first, second uint8
	if err := binary.Read(sr, binary.LittleEndian, &first); err != nil {
		return "", err
	}
	if (first & 0x80) != 0 {
		if err := binary.Read(sr, binary.LittleEndian, &second); err != nil {
			return "", err
		}
		size = (int(first&0x7F) << 8) + int(second)
	} else {
		size = int(first)
	}

	buf := make([]uint8, size)
	if err := binary.Read(sr, binary.LittleEndian, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}
