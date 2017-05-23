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
	styles  []ResStringPoolSpan
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
			0x01, 0x00, 0x00, 0x00,
			0x02, 0x00, 0x00, 0x00,
			0x03, 0x00, 0x00, 0x00,
			0x04, 0x00, 0x00, 0x00,
		},
		[]string{"a", "\3042"},
		[]ResStringPoolSpan{{1, 2}, {3, 4}},
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
		[]uint8{0x00, 0x00, 0x61},
		"",
	},
	{
		[]uint8{0x01, 0x01, 0x61},
		"a",
	},
	{
		[]uint8{0x01, 0x03, 0xE3, 0x81, 0x82},
		"\u3042",
	},
	{
		[]uint8{0x80, 0x01, 0x80, 0x01, 0x61},
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

var newZeroFilledReaderTests = []struct {
	input    []uint8
	actual   int64
	expected int64
	output   []uint8
}{
	{
		input:    []uint8{0x01, 0x23, 0x45, 0x67},
		actual:   4,
		expected: 4,
		output:   []uint8{0x01, 0x23, 0x45, 0x67},
	},
	{
		input:    []uint8{0x01, 0x23, 0x45, 0x67},
		actual:   4,
		expected: 8,
		output:   []uint8{0x01, 0x23, 0x45, 0x67, 0x00, 0x00, 0x00, 0x00},
	},
}

func TestNewZeroFilledReader(t *testing.T) {
	for _, tt := range newZeroFilledReaderTests {
		buf := bytes.NewReader(tt.input)
		r, err := newZeroFilledReader(buf, tt.actual, tt.expected)
		if err != nil {
			t.Errorf("got %v want no error", err)
		}
		actualBytes := make([]uint8, tt.expected)
		r.Read(actualBytes)
		for i, a := range tt.output {
			if actualBytes[i] != a {
				t.Errorf("got %v(position %d) wants %v", actualBytes[i], i, a)
			}
		}
	}
}
