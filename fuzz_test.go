//go:build go1.18
// +build go1.18

package androidbinary

import (
	"bytes"
	"os"
	"testing"
)

func FuzzNewXMLFile(f *testing.F) {
	data, err := os.ReadFile("testdata/AndroidManifest.xml")
	if err != nil {
		f.Fatal(err)
	}
	f.Add(data)

	f.Fuzz(func(t *testing.T, data []byte) {
		_, err := NewXMLFile(bytes.NewReader(data))
		if err != nil {
			t.Skip(err)
		}
	})
}
