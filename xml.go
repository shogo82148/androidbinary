package androidbinary

import (
	"bytes"
	"encoding/binary"
	"io"
	"os"
)

type XMLFile struct {
	stringPool     *ResStringPool
	resourceMap    []uint32
	notPrecessedNS map[ResStringPoolRef]ResStringPoolRef
	namespaces     map[ResStringPoolRef]ResStringPoolRef
	xmlBuffer      bytes.Buffer
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

type ResXMLTreeEndElementExt struct {
	NS   ResStringPoolRef
	Name ResStringPoolRef
}

func (f *XMLFile) readStartNamespace(sr *io.SectionReader) error {
	header := new(ResXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return err
	}

	sr.Seek(int64(header.Header.HeaderSize), os.SEEK_SET)
	namespace := new(ResXMLTreeNamespaceExt)
	if err := binary.Read(sr, binary.LittleEndian, namespace); err != nil {
		return err
	}

	if f.notPrecessedNS == nil {
		f.notPrecessedNS = make(map[ResStringPoolRef]ResStringPoolRef)
	}
	f.notPrecessedNS[namespace.Uri] = namespace.Prefix

	if f.namespaces == nil {
		f.namespaces = make(map[ResStringPoolRef]ResStringPoolRef)
	}
	f.namespaces[namespace.Uri] = namespace.Prefix

	return nil
}

func (f *XMLFile) readEndNamespace(sr *io.SectionReader) error {
	header := new(ResXMLTreeNode)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return err
	}

	sr.Seek(int64(header.Header.HeaderSize), os.SEEK_SET)
	namespace := new(ResXMLTreeNamespaceExt)
	if err := binary.Read(sr, binary.LittleEndian, namespace); err != nil {
		return err
	}
	delete(f.namespaces, namespace.Uri)
	return nil
}
