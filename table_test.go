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
		if got := IsResID(c.input); got != c.want {
			t.Errorf("%s: want %v, got %v", c.input, got, c.want)
		}
	}
}

func TestParseResId(t *testing.T) {
	id, err := ParseResID("@0x12345678")
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

func TestGetResourceNil(t *testing.T) {
	tableFile := loadTestData()
	val, _ := tableFile.GetResource(ResID(0x7f040000), nil)
	if val != "花火距離計算" {
		t.Errorf(`got %v want "花火距離計算"`, val)
	}
}

func TestGetResourceDefault(t *testing.T) {
	tableFile := loadTestData()
	val, _ := tableFile.GetResource(ResID(0x7f040000), &ResTableConfig{})
	if val != "FireworksMeasure" {
		t.Errorf(`got %v want "FireworksMeasure"`, val)
	}
}

func TestGetResourceJA(t *testing.T) {
	tableFile := loadTestData()
	config := &ResTableConfig{
		Language: [2]uint8{'j', 'a'},
	}
	val, _ := tableFile.GetResource(ResID(0x7f040000), config)
	if val != "花火距離計算" {
		t.Errorf(`got %v want "花火距離計算"`, val)
	}
}

func TestGetResourceEN(t *testing.T) {
	tableFile := loadTestData()
	config := &ResTableConfig{
		Language: [2]uint8{'e', 'n'},
	}
	val, _ := tableFile.GetResource(ResID(0x7f040000), config)
	if val != "FireworksMeasure" {
		t.Errorf(`got %v want "FireworksMeasure"`, val)
	}
}

var isMoreSpecificThanTests = []struct {
	me       *ResTableConfig
	other    *ResTableConfig
	expected bool
}{
	{
		me:       nil,
		other:    nil,
		expected: false,
	},
	{
		me:       nil,
		other:    &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{},
		other:    nil,
		expected: false,
	},
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
		me: &ResTableConfig{
			Language: [2]uint8{'j', 'a'},
		},
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
		me:       &ResTableConfig{ScreenLayout: ScreenSizeNormal},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenLayout: ScreenLongYes},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenLayout: LayoutDirLTR},
		other:    &ResTableConfig{ScreenLayout: LayoutDirAny},
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
		me:       &ResTableConfig{UIMode: UIModeTypeAny},
		other:    &ResTableConfig{},
		expected: true,
	},
	{
		me:       &ResTableConfig{UIMode: UIModeNightYes},
		other:    &ResTableConfig{UIMode: UIModeNightAny},
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
		me:       &ResTableConfig{UIMode: UIModeTypeAny},
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
				t.Errorf("%+v is more specific than %+v, but get false", tt.me, tt.other)
			} else {
				t.Errorf("%+v is not more specific than %+v, but get true", tt.me, tt.other)
			}
		}

		if tt.expected {
			// If 'me' is more specific than 'other', 'other' is not more specific than 'me'
			if tt.other.IsMoreSpecificThan(tt.me) {
				t.Errorf("%+v is not more specific than %+v, but get true", tt.other, tt.me)
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
		me:       &ResTableConfig{ScreenLayout: ScreenSizeNormal},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{ScreenLayout: ScreenSizeNormal},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{ScreenLayout: ScreenSizeNormal},
		expected: true,
	},
	{
		me:       &ResTableConfig{},
		other:    &ResTableConfig{ScreenLayout: ScreenSizeSmall},
		require:  &ResTableConfig{ScreenLayout: ScreenSizeXLarge},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenLayout: ScreenSizeSmall},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{ScreenLayout: ScreenSizeSmall},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenLayout: ScreenLongYes},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{},
		expected: false,
	},
	{
		me:       &ResTableConfig{ScreenLayout: ScreenLongYes},
		other:    &ResTableConfig{},
		require:  &ResTableConfig{ScreenLayout: ScreenLongYes},
		expected: true,
	},
	{
		me:       &ResTableConfig{ScreenLayout: LayoutDirLTR},
		other:    &ResTableConfig{ScreenLayout: LayoutDirAny},
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
				t.Errorf("%+v is better than %+v, but get false (%+v)", tt.me, tt.other, tt.require)
			} else {
				t.Errorf("%+v is not better than %+v, but get true (%+v)", tt.me, tt.other, tt.require)
			}
		}

		if tt.expected {
			// If 'me' is better than 'other', 'other' is not better than 'me'
			if tt.other.IsBetterThan(tt.me, tt.require) {
				t.Errorf("%v is not better than %+v, but get true (%+v)", tt.other, tt.me, tt.require)
			}
		}
	}
}
