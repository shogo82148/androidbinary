package main

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"unicode/utf16"
)

type File struct {
	stringPool  *ResStringPool
	resourceMap []uint32
	Namespace   *ResXMLTreeNamespaceExt
	XMLBuffer   bytes.Buffer
}

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

type ResStringPoolRef uint32

const NilResStringPoolRef = ResStringPoolRef(0xFFFFFFFF)

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

type ResXMLTreeNode struct {
	Header     ResChunkHeader
	LineNumber uint32
	Comment    ResStringPoolRef
}

type ResXMLTreeNamespaceExt struct {
	Prefix ResStringPoolRef
	Uri    ResStringPoolRef
}

type ResXMLTreeAttrExt struct {
	NS             ResStringPoolRef
	Name           ResStringPoolRef
	AttributeStart uint16
	AttributeSize  uint16
	AttributeCount uint16
	IdIndex        uint16
	ClassIndex     uint16
	StyleIndex     uint16
}

type ResXMLTreeAttribute struct {
	NS         ResStringPoolRef
	Name       ResStringPoolRef
	RawValue   ResStringPoolRef
	TypedValue ResValue
}

const (
	TYPE_NULL            = 0x00
	TYPE_REFERENCE       = 0x01
	TYPE_ATTRIBUTE       = 0x02
	TYPE_STRING          = 0x03
	TYPE_FLOAT           = 0x04
	TYPE_DIMENSION       = 0x05
	TYPE_FRACTION        = 0x06
	TYPE_FIRST_INT       = 0x10
	TYPE_INT_DEC         = 0x10
	TYPE_INT_HEX         = 0x11
	TYPE_INT_BOOLEAN     = 0x12
	TYPE_FIRST_COLOR_INT = 0x1c
	TYPE_INT_COLOR_ARGB8 = 0x1c
	TYPE_INT_COLOR_RGB8  = 0x1d
	TYPE_INT_COLOR_ARGB4 = 0x1e
	TYPE_INT_COLOR_RGB4  = 0x1f
	TYPE_LAST_COLOR_INT  = 0x1f
	TYPE_LAST_INT        = 0x1f
)

type ResValue struct {
	Size     uint16
	Res0     uint8
	DataType uint8
	Data     uint32
}

type ResXMLTreeEndElementExt struct {
	NS   ResStringPoolRef
	Name ResStringPoolRef
}

func NewFile(r io.ReaderAt) (*File, error) {
	f := new(File)
	sr := io.NewSectionReader(r, 0, 1<<63-1)

	header := new(ResChunkHeader)
	binary.Read(sr, binary.LittleEndian, header)
	offset := uint32(header.HeaderSize)

	for offset < header.Size {
		sr.Seek(int64(offset), os.SEEK_SET)
		chunkHeader := &ResChunkHeader{}
		binary.Read(sr, binary.LittleEndian, chunkHeader)

		var err error
		chunkReader := io.NewSectionReader(r, int64(offset), int64(chunkHeader.Size))
		switch chunkHeader.Type {
		case RES_STRING_POOL_TYPE:
			f.stringPool, err = ReadStringPool(chunkReader)
		case RES_XML_RESOURCE_MAP_TYPE:
			f.resourceMap, err = ReadResourceMap(chunkReader)
		case RES_XML_START_NAMESPACE_TYPE:
			err = f.ReadStartNamespace(chunkReader)
		case RES_XML_START_ELEMENT_TYPE:
			err = f.ReadStartElement(chunkReader)
		case RES_XML_END_ELEMENT_TYPE:
			err = f.ReadEndElement(chunkReader)
		}
		if err != nil {
			return nil, err
		}

		offset += chunkHeader.Size
	}
	return f, nil
}

func (f *File) GetString(ref ResStringPoolRef) string {
	if ref == NilResStringPoolRef {
		return ""
	}
	return f.stringPool.Strings[int(ref)]
}

func ReadStringPool(sr *io.SectionReader) (*ResStringPool, error) {
	sp := new(ResStringPool)
	binary.Read(sr, binary.LittleEndian, &sp.Header)

	stringStarts := make([]uint32, sp.Header.StringCount)
	binary.Read(sr, binary.LittleEndian, stringStarts)
	styleStarts := make([]uint32, sp.Header.StyleCount)
	binary.Read(sr, binary.LittleEndian, styleStarts)

	sp.Strings = make([]string, sp.Header.StringCount)
	for i, start := range stringStarts {
		var str string
		var err error
		if (sp.Header.Flags & UTF8_FLAG) == 0 {
			str, err = ReadUTF16(sr, int64(sp.Header.StringStart+start))
		} else {
			str, err = ReadUTF8(sr, int64(sp.Header.StringStart+start))
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
		if (sp.Header.Flags & UTF8_FLAG) == 0 {
			str, err = ReadUTF16(sr, int64(sp.Header.StylesStart+start))
		} else {
			str, err = ReadUTF8(sr, int64(sp.Header.StylesStart+start))
		}
		if err != nil {
			return nil, err
		}
		sp.Styles[i] = str
	}

	return sp, nil
}

func ReadUTF16(sr *io.SectionReader, offset int64) (string, error) {
	var size uint16
	sr.Seek(offset, os.SEEK_SET)
	if err := binary.Read(sr, binary.LittleEndian, &size); err != nil {
		return "", err
	}
	buf := make([]uint16, size)
	if err := binary.Read(sr, binary.LittleEndian, buf); err != nil {
		return "", err
	}
	return string(utf16.Decode(buf)), nil
}

func ReadUTF8(sr *io.SectionReader, offset int64) (string, error) {
	var size uint16
	sr.Seek(offset, os.SEEK_SET)
	if err := binary.Read(sr, binary.LittleEndian, &size); err != nil {
		return "", err
	}
	buf := make([]uint8, size)
	if err := binary.Read(sr, binary.LittleEndian, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func ReadResourceMap(sr *io.SectionReader) ([]uint32, error) {
	header := new(ResChunkHeader)
	binary.Read(sr, binary.LittleEndian, header)
	count := (header.Size - uint32(header.HeaderSize)) / 4
	resourceMap := make([]uint32, count)
	if err := binary.Read(sr, binary.LittleEndian, resourceMap); err != nil {
		return nil, err
	}
	return resourceMap, nil
}

func (f *File) ReadStartNamespace(sr *io.SectionReader) error {
	header := new(ResXMLTreeNode)
	binary.Read(sr, binary.LittleEndian, header)
	sr.Seek(int64(header.Header.HeaderSize), os.SEEK_SET)
	namespace := new(ResXMLTreeNamespaceExt)
	binary.Read(sr, binary.LittleEndian, namespace)
	f.Namespace = namespace
	return nil
}

func (f *File) ReadStartElement(sr *io.SectionReader) error {
	header := new(ResXMLTreeNode)
	binary.Read(sr, binary.LittleEndian, header)
	sr.Seek(int64(header.Header.HeaderSize), os.SEEK_SET)
	ext := new(ResXMLTreeAttrExt)
	binary.Read(sr, binary.LittleEndian, ext)

	fmt.Fprintf(&f.XMLBuffer, "<%s", f.AddNamespace(ext.NS, ext.Name))

	offset := int64(ext.AttributeStart + header.Header.HeaderSize)
	for i := 0; i < int(ext.AttributeCount); i++ {
		sr.Seek(offset, os.SEEK_SET)
		attr := new(ResXMLTreeAttribute)
		binary.Read(sr, binary.LittleEndian, attr)

		var value string
		if attr.RawValue != NilResStringPoolRef {
			value = f.GetString(attr.RawValue)
		} else {
			data := attr.TypedValue.Data
			switch attr.TypedValue.DataType {
			case TYPE_NULL:
				value = ""
			case TYPE_REFERENCE:
				value = fmt.Sprintf("@0x%08X", data)
			case TYPE_INT_DEC:
				value = fmt.Sprintf("%d", data)
			case TYPE_INT_HEX:
				value = fmt.Sprintf("0x%08X", data)
			case TYPE_INT_BOOLEAN:
				if data != 0 {
					value = "true"
				} else {
					value = "false"
				}
			default:
				value = fmt.Sprintf("@0x%08X", data)
			}
		}

		fmt.Fprintf(&f.XMLBuffer, " %s=\"", f.AddNamespace(attr.NS, attr.Name))
		xml.Escape(&f.XMLBuffer, []byte(value))
		fmt.Fprint(&f.XMLBuffer, "\"")
		offset += int64(ext.AttributeSize)
	}
	fmt.Fprint(&f.XMLBuffer, ">")
	return nil
}

func (f *File) ReadEndElement(sr *io.SectionReader) error {
	header := new(ResXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return err
	}
	sr.Seek(int64(header.Header.HeaderSize), os.SEEK_SET)
	ext := new(ResXMLTreeEndElementExt)
	if err := binary.Read(sr, binary.LittleEndian, ext); err != nil {
		return err
	}
	fmt.Fprintf(&f.XMLBuffer, "</%s>", f.AddNamespace(ext.NS, ext.Name))
	return nil
}

func (f *File) AddNamespace(ns, name ResStringPoolRef) string {
	if ns != NilResStringPoolRef {
		return fmt.Sprintf("%s:%s", f.GetString(f.Namespace.Prefix), f.GetString(name))
	} else {
		return f.GetString(name)
	}
}

func main() {
	f, _ := os.Open("AndroidManifest.xml")
	xml, _ := NewFile(f)
	fmt.Println(xml.XMLBuffer.String())
}
