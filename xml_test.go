package androidbinary

import (
	"bytes"
	"io"
	"testing"
)

func TestReadStartNamespace(t *testing.T) {
	input := []uint8{
		0x00, 0x01, // Type = RES_XML_START_NAMESPACE_TYPE
		0x10, 0x00, // HeadSize = 16 bytes
		0x18, 0x00, 0x00, 0x00, // Size
		0x01, 0x00, 0x00, 0x00, // LineNumber = 1
		0xFF, 0xFF, 0xFF, 0xFF, // Comment is none
		0x02, 0x00, 0x00, 0x00, // Prefix = 2
		0x01, 0x00, 0x00, 0x00, // Uri = 1
	}
	buf := bytes.NewReader(input)
	sr := io.NewSectionReader(buf, 0, int64(len(input)))
	f := new(XMLFile)
	err := f.readStartNamespace(sr)
	if err != nil {
		t.Errorf("got %v want no error", err)
	}
	if f.notPrecessedNS[ResStringPoolRef(1)] != ResStringPoolRef(2) {
		t.Errorf("got %v want %v", f.notPrecessedNS[ResStringPoolRef(1)], ResStringPoolRef(2))
	}
	if f.namespaces[ResStringPoolRef(1)] != ResStringPoolRef(2) {
		t.Errorf("got %v want %v", f.namespaces[ResStringPoolRef(1)], ResStringPoolRef(2))
	}
}
