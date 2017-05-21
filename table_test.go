package androidbinary

import (
	"os"
	"testing"
)

func TestIsResId(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"@0x00", true},
		{"foo", false},
	}
	for _, c := range cases {
		if got := IsResId(c.input); got != c.want {
			t.Errorf("%s: want %v, got %v", c.input, got, c.want)
		}
	}
}

func TestParseResId(t *testing.T) {
	id, err := ParseResId("@0x12345678")
	if err != nil {
		t.Error(err)
	}
	if id != 0x12345678 {
		t.Errorf("want 0x12345678, got %X", id)
	}
}

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
	{
		me:       &ResTableConfig{ScreenLayout: SCREENSIZE_NORMAL},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenLayout: SCREENLONG_YES},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenLayout: LAYOUTDIR_LTR},
		other:    &ResTableConfig{ScreenLayout: LAYOUTDIR_ANY},
		expected: true,
	},
	{
		me:       &ResTableConfig{SmallestScreenWidthDp: 72},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenWidthDp: 100},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenHeightDp: 100},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{Orientation: 1},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{UIMode: UI_MODE_TYPE_ANY},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{UIMode: UI_MODE_NIGHT_YES},
		other:    &ResTableConfig{UIMode: UI_MODE_NIGHT_ANY},
		expected: true,
	},
	{
		me:       &ResTableConfig{Keyboard: 1},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{Navigation: 1},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{UIMode: UI_MODE_TYPE_ANY},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{Touchscreen: 1},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenWidth: 100},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenHeight: 100},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{SDKVersion: 1},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{SDKVersion: 1, MinorVersion: 1},
		other:    &ResTableConfig{SDKVersion: 1},
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

func TestIsBetterThan_request_is_nil(t *testing.T) {
	// a.IsBetterThan(b, nil) is same as a.IsMoreSpecificThan(b)
	for _, tt := range isMoreSpecificThanTests {
		actual := tt.me.IsBetterThan(tt.other, nil)
		if actual != tt.expected {
			if tt.expected {
				t.Errorf("%v is better than %v, but get false", tt.me, tt.other)
			} else {
				t.Errorf("%v is better than %v, but get true", tt.me, tt.other)
			}
		}

		if tt.expected {
			// If 'me' is more specific than 'other', 'other' is not more specific than 'me'
			if tt.other.IsBetterThan(tt.me, nil) {
				t.Errorf("%v is better than %v, but get true", tt.other, tt.me)
			}
		}
	}
}

var isBetterThanTests = []struct {
	me       *ResTableConfig
	other    *ResTableConfig
	require  *ResTableConfig
	expected bool
}{
	{
		me:       &ResTableConfig{},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{Mcc: 1},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{Mcc: 1},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{Mcc: 1},
		expected: true,
	},
	{
		me:       &ResTableConfig{Mcc: 1, Mnc: 1},
		other:    &ResTableConfig{Mcc: 1},
		require:  &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{Language: [2]byte{'j', 'a'}},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{Language: [2]byte{'j', 'a'}},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{Language: [2]byte{'j', 'a'}},
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
		require:  &ResTableConfig{},
		expected: false,
	},
	{
		me: &ResTableConfig{
			Language: [2]uint8{'j', 'a'},
			Country:  [2]uint8{'J', 'P'},
		},
		other: &ResTableConfig{
			Language: [2]uint8{'j', 'a'},
		},
		require: &ResTableConfig{
			Language: [2]uint8{'j', 'a'},
			Country:  [2]uint8{'J', 'P'},
		},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenLayout: SCREENSIZE_NORMAL},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{ScreenLayout: SCREENSIZE_NORMAL},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{ScreenLayout: SCREENSIZE_NORMAL},
		expected: true,
	},
	{
		me:       &ResTableConfig{},
		other:    &ResTableConfig{ScreenLayout: SCREENSIZE_SMALL},
		require:  &ResTableConfig{ScreenLayout: SCREENSIZE_XLARGE},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenLayout: SCREENSIZE_SMALL},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{ScreenLayout: SCREENSIZE_SMALL},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenLayout: SCREENLONG_YES},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{ScreenLayout: SCREENLONG_YES},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{ScreenLayout: SCREENLONG_YES},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenLayout: LAYOUTDIR_LTR},
		other:    &ResTableConfig{ScreenLayout: LAYOUTDIR_ANY},
		expected: true,
	},
	{
		me:       &ResTableConfig{SmallestScreenWidthDp: 72},
		other:    &ResTableConfig{SmallestScreenWidthDp: 71},
		require:  &ResTableConfig{SmallestScreenWidthDp: 72},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenWidthDp: 100},
		other:    &ResTableConfig{ScreenWidthDp: 99},
		require:  &ResTableConfig{ScreenWidthDp: 100},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenHeightDp: 100},
		other:    &ResTableConfig{ScreenHeightDp: 99},
		require:  &ResTableConfig{ScreenHeightDp: 100},
		expected: true,
	},
	{
		me:       &ResTableConfig{Orientation: 1},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{Orientation: 1},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{Orientation: 1},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenWidth: 100},
		other:    &ResTableConfig{ScreenWidth: 99},
		require:  &ResTableConfig{ScreenWidth: 100},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenHeight: 100},
		other:    &ResTableConfig{ScreenHeight: 99},
		require:  &ResTableConfig{ScreenHeight: 100},
		expected: true,
	},
	{
		me:       &ResTableConfig{SDKVersion: 2},
		other:    &ResTableConfig{SDKVersion: 1},
		require:  &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{SDKVersion: 2},
		other:    &ResTableConfig{SDKVersion: 1},
		require:  &ResTableConfig{SDKVersion: 1},
		expected: true,
	},
	{
		me:       &ResTableConfig{SDKVersion: 1, MinorVersion: 1},
		other:    &ResTableConfig{SDKVersion: 1},
		require:  &ResTableConfig{SDKVersion: 1},
		expected: false,
	},
	{
		me:       &ResTableConfig{SDKVersion: 1, MinorVersion: 1},
		other:    &ResTableConfig{SDKVersion: 1},
		require:  &ResTableConfig{SDKVersion: 1, MinorVersion: 1},
		expected: true,
	},
}

func TestIsBetterThan(t *testing.T) {
	for _, tt := range isBetterThanTests {
		actual := tt.me.IsBetterThan(tt.other, tt.require)
		if actual != tt.expected {
			if tt.expected {
				t.Errorf("%v is better than %v, but get false", tt.me, tt.other)
			} else {
				t.Errorf("%v is not better than %v, but get true", tt.me, tt.other)
			}
		}

		if tt.expected {
			// If 'me' is better than 'other', 'other' is not better than 'me'
			if tt.other.IsBetterThan(tt.me, tt.require) {
				t.Errorf("%v is not better than %v, but get true", tt.other, tt.me)
			}
		}
	}
}
