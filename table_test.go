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

func TestFindType(t *testing.T) {
	tableFile := loadTestData()
	p := tableFile.findPackage(0x7F)
	id := 0x04
	config := &ResTableConfig{}
	tableType := p.findType(id, config)
	if int(tableType.Header.Id) != id {
		t.Errorf("got %v want %v", tableType.Header.Id, id)
	}
	locale := tableType.Header.Config.Locale()
	if locale != "" {
		t.Errorf("got %v want \"\"", locale)
	}
}

func TestFindTypeJa(t *testing.T) {
	tableFile := loadTestData()
	p := tableFile.findPackage(0x7F)
	id := 0x04
	config := &ResTableConfig{}
	config.Language[0] = uint8('j')
	config.Language[1] = uint8('a')
	tableType := p.findType(id, config)
	if int(tableType.Header.Id) != id {
		t.Errorf("got %v want %v", tableType.Header.Id, id)
	}
	locale := tableType.Header.Config.Locale()
	if locale != "ja" {
		t.Errorf("got %v want ja", locale)
	}
}

func TestFindTypeEn(t *testing.T) {
	tableFile := loadTestData()
	p := tableFile.findPackage(0x7F)
	id := 0x04
	config := &ResTableConfig{}
	config.Language[0] = uint8('e')
	config.Language[1] = uint8('n')
	tableType := p.findType(id, config)
	if int(tableType.Header.Id) != id {
		t.Errorf("got %v want %v", tableType.Header.Id, id)
	}
	locale := tableType.Header.Config.Locale()
	if locale != "" {
		t.Errorf("got %v want \"\"", locale)
	}
}

func TestGetResource(t *testing.T) {
	tableFile := loadTestData()
	config := &ResTableConfig{}
	val, _ := tableFile.GetResource(ResId(0x7f040000), config)
	if val != "FireworksMeasure" {
		t.Errorf("got %v want \"\"", val)
	}
}
