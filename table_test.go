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

var isMoreSpecificThanTests = []struct {
	me       *ResTableConfig
	other    *ResTableConfig
	expected bool
}{
	{
		me:       &ResTableConfig{},
		other:    &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{Mcc: 1},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{Mcc: 1, Mnc: 1},
		other:    &ResTableConfig{Mcc: 1},
		expected: true,
	},
	{
		me:       &ResTableConfig{Language: [2]byte{'j', 'a'}},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me: &ResTableConfig{
			Language: [2]uint8{'j', 'a'},
			Country:  [2]uint8{'J', 'P'},
		},
		other: &ResTableConfig{
			Language: [2]uint8{'j', 'a'},
		},
		expected: true,
	},
}

func TestIsMoreSpecificThan(t *testing.T) {
	for _, tt := range isMoreSpecificThanTests {
		actual := tt.me.IsMoreSpecificThan(tt.other)
		if actual != tt.expected {
			if tt.expected {
				t.Errorf("%v is more specific than %v, but get false", tt.me, tt.other)
			} else {
				t.Errorf("%v is not more specific than %v, but get true", tt.me, tt.other)
			}
		}

		if tt.expected {
			// If 'me' is more specific than 'other', 'other' is not more specific than 'me'
			if tt.other.IsMoreSpecificThan(tt.me) {
				t.Errorf("%v is not more specific than %v, but get true", tt.other, tt.me)
			}
		}
	}
}
