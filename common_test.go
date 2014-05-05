package androidbinary

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

var readStringPoolTests = []struct {
	input   []uint8
	strings []string
	styles  []string
}{
	{
		[]uint8{
			0x01, 0x00, // Type = RES_STRING_POOL_TYPE
			0x1C, 0x00, // HeaderSize = 28 bytes
			0x3C, 0x00, 0x00, 0x00, // Size = 60
			0x02, 0x00, 0x00, 0x00, // StringCount = 2
			0x02, 0x00, 0x00, 0x00, // StyleScount = 2
			0x00, 0x00, 0x00, 0x00, // Flags = 0x00
			0x2C, 0x00, 0x00, 0x00, // StringStart = 44
			0x34, 0x00, 0x00, 0x00, // StylesStart = 52

			// StringIndexes
			0x00, 0x00, 0x00, 0x00,
			0x04, 0x00, 0x00, 0x00,

			// StyleIndexes
			0x00, 0x00, 0x00, 0x00,
			0x04, 0x00, 0x00, 0x00,

			// Strings
			0x01, 0x00, 0x61, 0x00,
			0x01, 0x00, 0x42, 0x30,

			// Styles
			0x01, 0x00, 0x63, 0x00,
			0x01, 0x00, 0x43, 0x30,
		},
		[]string{"a", "\3042"},
		[]string{"b", "\3043"},
	},
}

func TestReadStringPool(t *testing.T) {
	for _, tt := range readStringPoolTests {
		buf := bytes.NewReader(tt.input)
		sr := io.NewSectionReader(buf, 0, int64(len(tt.input)))
		actual, err := readStringPool(sr)
		if err != nil {
			t.Errorf("got %v want no error", err)
		}
		if reflect.DeepEqual(actual.Strings, tt.strings) {
			t.Errorf("got %v want %v", actual.Strings, tt.strings)
		}
		if reflect.DeepEqual(actual.Styles, tt.styles) {
			t.Errorf("got %v want %v", actual.Styles, tt.styles)
		}
	}
}

var readUTF16Tests = []struct {
	input  []uint8
	output string
}{
	{
		[]uint8{0x00, 0x00, 0x61, 0x00},
		"",
	},
	{
		[]uint8{0x01, 0x00, 0x61, 0x00},
		"a",
	},
	{
		[]uint8{0x01, 0x00, 0x42, 0x30},
		"\u3042",
	},
	{
		[]uint8{0x00, 0x80, 0x01, 0x00, 0x61, 0x00},
		"a",
	},
}

func TestReadUTF16(t *testing.T) {
	for _, tt := range readUTF16Tests {
		buf := bytes.NewReader(tt.input)
		sr := io.NewSectionReader(buf, 0, int64(len(tt.input)))
		actual, err := readUTF16(sr)
		if err != nil {
			t.Errorf("got %v want no error", err)
		}
		if actual != tt.output {
			t.Errorf("got %v want %v", actual, tt.output)
		}
	}
}

var readUTF8Tests = []struct {
	input  []uint8
	output string
}{
	{
		[]uint8{0x00, 0x61},
		"",
	},
	{
		[]uint8{0x01, 0x61},
		"a",
	},
	{
		[]uint8{0x03, 0xE3, 0x81, 0x82},
		"\u3042",
	},
	{
		[]uint8{0x80, 0x01, 0x61},
		"a",
	},
}

func TestReadUTF8(t *testing.T) {
	for _, tt := range readUTF8Tests {
		buf := bytes.NewReader(tt.input)
		sr := io.NewSectionReader(buf, 0, int64(len(tt.input)))
		actual, err := readUTF8(sr)
		if err != nil {
			t.Errorf("got %v want no error", err)
		}
		if actual != tt.output {
			t.Errorf("got %v want %v", actual, tt.output)
		}
	}
}
