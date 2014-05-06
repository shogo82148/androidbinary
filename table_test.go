package androidbinary

import (
	"os"
	"testing"
)

func loadTestData() *TableFile {
	f, _ := os.Open("testdata/resources.arsc")
	tableFile, _ := NewTableFile(f)
	return tableFile
}

func TestFindPackage(t *testing.T) {
	tableFile := loadTestData()
	p := tableFile.findPackage(0x7F)
	if p == nil {
		t.Error("got nil want package(id: 0x7F)")
		t.Errorf("%v", tableFile.tablePackages)
	}
}
