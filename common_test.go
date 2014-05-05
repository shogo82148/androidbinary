package androidbinary

import (
	"bytes"
	"io"
	"testing"
)

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
