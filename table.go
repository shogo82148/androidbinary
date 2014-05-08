package androidbinary

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

type ResId uint32

type TableFile struct {
	stringPool    *ResStringPool
	tablePackages map[uint32]*TablePackage
}

type ResTableHeader struct {
	Header       ResChunkHeader
	PackageCount uint32
}

type ResTablePackage struct {
	Header         ResChunkHeader
	Id             uint32
	Name           [128]uint16
	TypeStrings    uint32
	LastPublicType uint32
	KeyStrings     uint32
	LastPublicKey  uint32
}

type TablePackage struct {
	Header      ResTablePackage
	TypeStrings *ResStringPool
	KeyStrings  *ResStringPool
	TableTypes  []*TableType
}

type ResTableType struct {
	Header       ResChunkHeader
	Id           uint8
	Res0         uint8
	Res1         uint16
	EntryCount   uint32
	EntriesStart uint32
	Config       ResTableConfig
}

// ScreenLayout bits
const (
	MASK_SCREENSIZE   = 0x0f
	SCREENSIZE_ANY    = 0x01
	SCREENSIZE_SMALL  = 0x02
	SCREENSIZE_NORMAL = 0x03
	SCREENSIZE_LARGE  = 0x04
	SCREENSIZE_XLARGE = 0x05

	MASK_SCREENLONG  = 0x30
	SHIFT_SCREENLONG = 4
	SCREENLONG_ANY   = 0x00
	SCREENLONG_NO    = 0x10
	SCREENLONG_YES   = 0x20

	MASK_LAYOUTDIR  = 0xC0
	SHIFT_LAYOUTDIR = 6
	LAYOUTDIR_ANY   = 0x00
	LAYOUTDIR_LTR   = 0x40
	LAYOUTDIR_RTL   = 0x80
)

// UIMode bits
const (
	MASK_UI_MODE_TYPE   = 0x0f
	UI_MODE_TYPE_ANY    = 0x01
	UI_MODE_TYPE_NORMAL = 0x02
	UI_MODE_TYPE_DESK   = 0x03
	UI_MODE_TYPE_CAR    = 0x04

	MASK_UI_MODE_NIGHT  = 0x30
	SHIFT_UI_MODE_NIGHT = 4
	UI_MODE_NIGHT_ANY   = 0x00
	UI_MODE_NIGHT_NO    = 0x10
	UI_MODE_NIGHT_YES   = 0x20
)

type ResTableConfig struct {
	Size uint32
	// imsi
	Mcc uint16
	Mnc uint16

	// locale
	Language [2]uint8
	Country  [2]uint8

	// screen type
	Orientation uint8
	Touchscreen uint8
	Density     uint16

	// inout
	Keyboard   uint8
	Navigation uint8
	InputFlags uint8
	InputPad0  uint8

	// screen
	ScreenWidth  uint16
	ScreenHeight uint16

	// version
	SDKVersion   uint16
	MinorVersion uint16

	// screen config
	ScreenLayout          uint8
	UIMode                uint8
	SmallestScreenWidthDp uint16
}

type TableType struct {
	Header  *ResTableType
	Entries []TableEntry
}

type ResTableEntry struct {
	Size  uint16
	Flags uint16
	Key   ResStringPoolRef
}

type TableEntry struct {
	Key   *ResTableEntry
	Value *ResValue
	Flags uint32
}

type ResTableTypeSpec struct {
	Header     ResChunkHeader
	Id         uint8
	Res0       uint8
	Res1       uint16
	EntryCount uint32
}

func (id ResId) Package() int {
	return int(id) >> 24
}

func (id ResId) Type() int {
	return (int(id) >> 16) & 0xFF
}

func (id ResId) Entry() int {
	return int(id) & 0xFFFF
}

func NewTableFile(r io.ReaderAt) (*TableFile, error) {
	f := new(TableFile)
	sr := io.NewSectionReader(r, 0, 1<<63-1)

	header := new(ResTableHeader)
	binary.Read(sr, binary.LittleEndian, header)
	f.tablePackages = make(map[uint32]*TablePackage)

	offset := int64(header.Header.HeaderSize)
	for offset < int64(header.Header.Size) {
		chunkHeader, err := f.readChunk(sr, offset)
		if err != nil {
			return nil, err
		}
		offset += int64(chunkHeader.Size)
	}
	return f, nil
}

func (f *TableFile) findPackage(id int) *TablePackage {
	return f.tablePackages[uint32(id)]
}

func (p *TablePackage) findType(id int, config *ResTableConfig) *TableType {
	var best *TableType
	for _, t := range p.TableTypes {
		if int(t.Header.Id) != id {
			continue
		}
		if !t.Header.Config.Match(config) {
			continue
		}
		if best == nil || t.Header.Config.IsBetterThan(&best.Header.Config, config) {
			best = t
		}
	}
	return best
}

func (f *TableFile) GetResource(id ResId, config *ResTableConfig) (interface{}, error) {
	p := f.findPackage(id.Package())
	t := p.findType(id.Type(), config)
	e := t.Entries[id.Entry()]
	v := e.Value
	switch v.DataType {
	case TYPE_NULL:
		return nil, nil
	case TYPE_STRING:
		return f.GetString(ResStringPoolRef(v.Data)), nil
	case TYPE_INT_DEC:
		return v.Data, nil
	case TYPE_INT_HEX:
		return v.Data, nil
	case TYPE_INT_BOOLEAN:
		return v.Data != 0, nil
	}
	return v.Data, nil
}

func (f *TableFile) GetString(ref ResStringPoolRef) string {
	return f.stringPool.GetString(ref)
}

func (f *TableFile) readChunk(r io.ReaderAt, offset int64) (*ResChunkHeader, error) {
	sr := io.NewSectionReader(r, offset, 1<<63-1-offset)
	chunkHeader := &ResChunkHeader{}
	sr.Seek(0, os.SEEK_SET)
	if err := binary.Read(sr, binary.LittleEndian, chunkHeader); err != nil {
		return nil, err
	}

	var err error
	sr.Seek(0, os.SEEK_SET)
	switch chunkHeader.Type {
	case RES_STRING_POOL_TYPE:
		f.stringPool, err = readStringPool(sr)
	case RES_TABLE_PACKAGE_TYPE:
		var tablePackage *TablePackage
		tablePackage, err = readTablePackage(sr)
		f.tablePackages[tablePackage.Header.Id] = tablePackage
	}
	if err != nil {
		return nil, err
	}

	return chunkHeader, nil
}

func readTablePackage(sr *io.SectionReader) (*TablePackage, error) {
	tablePackage := new(TablePackage)
	header := new(ResTablePackage)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return nil, err
	}
	tablePackage.Header = *header

	srTypes := io.NewSectionReader(sr, int64(header.TypeStrings), int64(header.Header.Size-header.TypeStrings))
	if typeStrings, err := readStringPool(srTypes); err == nil {
		tablePackage.TypeStrings = typeStrings
	} else {
		return nil, err
	}

	srKeys := io.NewSectionReader(sr, int64(header.KeyStrings), int64(header.Header.Size-header.KeyStrings))
	if keyStrings, err := readStringPool(srKeys); err == nil {
		tablePackage.KeyStrings = keyStrings
	} else {
		return nil, err
	}

	offset := int64(header.Header.HeaderSize)
	for offset < int64(header.Header.Size) {
		chunkHeader := &ResChunkHeader{}
		sr.Seek(offset, os.SEEK_SET)
		if err := binary.Read(sr, binary.LittleEndian, chunkHeader); err != nil {
			return nil, err
		}

		var err error
		chunkReader := io.NewSectionReader(sr, offset, int64(chunkHeader.Size))
		sr.Seek(offset, os.SEEK_SET)
		switch chunkHeader.Type {
		case RES_TABLE_TYPE_TYPE:
			var tableType *TableType
			tableType, err = readTableType(chunkReader)
			tablePackage.TableTypes = append(tablePackage.TableTypes, tableType)
		case RES_TABLE_TYPE_SPEC_TYPE:
			_, err = readTableTypeSpec(chunkReader)
		}
		if err != nil {
			return nil, err
		}
		offset += int64(chunkHeader.Size)
	}

	return tablePackage, nil
}

func readTableType(sr *io.SectionReader) (*TableType, error) {
	header := new(ResTableType)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	entryIndexes := make([]uint32, header.EntryCount)
	sr.Seek(int64(header.Header.HeaderSize), os.SEEK_SET)
	if err := binary.Read(sr, binary.LittleEndian, entryIndexes); err != nil {
		return nil, err
	}

	entries := make([]TableEntry, header.EntryCount)
	for i, index := range entryIndexes {
		if index == 0xFFFFFFFF {
			continue
		}
		sr.Seek(int64(header.EntriesStart+index), os.SEEK_SET)
		var key ResTableEntry
		binary.Read(sr, binary.LittleEndian, &key)
		entries[i].Key = &key

		var val ResValue
		binary.Read(sr, binary.LittleEndian, &val)
		entries[i].Value = &val
	}
	return &TableType{
		header,
		entries,
	}, nil
}

func readTableTypeSpec(sr *io.SectionReader) ([]uint32, error) {
	header := new(ResTableTypeSpec)
	if err := binary.Read(sr, binary.LittleEndian, header); err != nil {
		return nil, err
	}

	flags := make([]uint32, header.EntryCount)
	sr.Seek(int64(header.Header.HeaderSize), os.SEEK_SET)
	if err := binary.Read(sr, binary.LittleEndian, flags); err != nil {
		return nil, err
	}
	return flags, nil
}

func (c *ResTableConfig) IsMoreSpecificThan(o *ResTableConfig) bool {
	// imsi
	if c.Mcc != o.Mcc {
		if c.Mcc != 0 {
			return false
		}
		if o.Mnc != 0 {
			return true
		}
	}
	if c.Mnc != o.Mnc {
		if c.Mnc != 0 {
			return false
		}
		if o.Mnc != 0 {
			return true
		}
	}

	// locale
	if c.Language[0] != o.Language[0] {
		if c.Language[0] != 0 {
			return false
		}
		if o.Language[0] != 0 {
			return true
		}
	}
	if c.Country[0] != o.Country[0] {
		if c.Country[0] != 0 {
			return false
		}
		if o.Country[0] != 0 {
			return true
		}
	}

	// screen layout
	if c.ScreenLayout != 0 || o.ScreenLayout != 0 {
		if ((c.ScreenLayout ^ o.ScreenLayout) & MASK_LAYOUTDIR) != 0 {
			if (c.ScreenLayout & MASK_LAYOUTDIR) != 0 {
				return false
			}
			if (o.ScreenLayout & MASK_LAYOUTDIR) != 0 {
				return true
			}
		}
	}

	// smallest screen width dp
	if c.SmallestScreenWidthDp != 0 || o.SmallestScreenWidthDp != 0 {
		if c.SmallestScreenWidthDp != o.SmallestScreenWidthDp {
			if c.SmallestScreenWidthDp != 0 {
				return false
			}
			if o.SmallestScreenWidthDp != 0 {
				return true
			}
		}
	}

	// screen layout
	if c.ScreenLayout != 0 || o.ScreenLayout != 0 {
		if ((c.ScreenLayout ^ o.ScreenLayout) & MASK_SCREENSIZE) != 0 {
			if (c.ScreenLayout & MASK_SCREENSIZE) != 0 {
				return false
			}
			if (o.ScreenLayout & MASK_SCREENSIZE) != 0 {
				return true
			}
		}
		if ((c.ScreenLayout ^ o.ScreenLayout) & MASK_SCREENLONG) != 0 {
			if (c.ScreenLayout & MASK_SCREENLONG) != 0 {
				return false
			}
			if (o.ScreenLayout & MASK_SCREENLONG) != 0 {
				return true
			}
		}
	}

	// orientation
	if c.Orientation != o.Orientation {
		if c.Orientation != 0 {
			return false
		}
		if o.Orientation != 0 {
			return true
		}
	}

	// uimode
	if c.UIMode != 0 || o.UIMode != 0 {
		diff := c.UIMode ^ o.UIMode
		if (diff & MASK_UI_MODE_TYPE) != 0 {
			if (c.UIMode & MASK_UI_MODE_TYPE) != 0 {
				return false
			}
			if (o.UIMode & MASK_UI_MODE_TYPE) != 0 {
				return true
			}
		}
		if (diff & MASK_UI_MODE_NIGHT) != 0 {
			if (c.UIMode & MASK_UI_MODE_NIGHT) != 0 {
				return false
			}
			if (o.UIMode & MASK_UI_MODE_NIGHT) != 0 {
				return true
			}
		}
	}

	// touchscreen
	if c.Touchscreen != o.Touchscreen {
		if c.Touchscreen != 0 {
			return false
		}
		if o.Touchscreen != 0 {
			return true
		}
	}

	// TODO: input

	// screen size
	if c.ScreenWidth != o.ScreenWidth {
		if c.ScreenWidth != 0 {
			return false
		}
		if o.ScreenWidth != 0 {
			return true
		}
	}
	if c.ScreenHeight != o.ScreenHeight {
		if c.ScreenHeight != 0 {
			return false
		}
		if o.ScreenHeight != 0 {
			return true
		}
	}

	//version
	if c.SDKVersion != o.SDKVersion {
		if c.SDKVersion != 0 {
			return false
		}
		if o.SDKVersion != 0 {
			return true
		}
	}
	if c.MinorVersion != o.MinorVersion {
		if c.MinorVersion != 0 {
			return false
		}
		if o.MinorVersion != 0 {
			return true
		}
	}

	return false
}

func (c *ResTableConfig) IsBetterThan(o *ResTableConfig, r *ResTableConfig) bool {
	if r == nil {
		return c.IsMoreSpecificThan(o)
	}

	// imsi
	if c.Mcc != 0 || c.Mnc != 0 || o.Mcc != 0 || o.Mnc != 0 {
		if c.Mcc != o.Mcc && r.Mcc != 0 {
			return c.Mcc != 0
		}
		if c.Mnc != o.Mnc && r.Mnc != 0 {
			return c.Mnc != 0
		}
	}

	// locale
	if c.Language[0] != 0 || c.Country[0] != 0 || o.Language[0] != 0 || o.Country[0] != 0 {
		if c.Language[0] != o.Language[0] && r.Language[0] != 0 {
			return c.Language[0] != 0
		}
		if c.Country[0] != o.Country[0] && r.Country[0] != 0 {
			return c.Country[0] != 0
		}
	}

	// screen layout
	if c.ScreenLayout != 0 || o.ScreenLayout != 0 {
		myLayoutdir := c.ScreenLayout & MASK_LAYOUTDIR
		oLayoutdir := o.ScreenLayout & MASK_LAYOUTDIR
		if (myLayoutdir^oLayoutdir) != 0 && (r.ScreenLayout&MASK_LAYOUTDIR) != 0 {
			return myLayoutdir > oLayoutdir
		}
	}

	// smallest screen width dp
	if c.SmallestScreenWidthDp != 0 || o.SmallestScreenWidthDp != 0 {
		if c.SmallestScreenWidthDp != o.SmallestScreenWidthDp {
			return c.SmallestScreenWidthDp > o.SmallestScreenWidthDp
		}
	}

	// screen layout
	if c.ScreenLayout != 0 || o.ScreenLayout != 0 {
		mySL := c.ScreenLayout & MASK_SCREENSIZE
		oSL := o.ScreenLayout & MASK_SCREENSIZE
		if (mySL^oSL) != 0 && (r.ScreenLayout&MASK_SCREENSIZE) != 0 {
			fixedMySL := mySL
			fixedOSL := oSL
			if (r.ScreenLayout & MASK_SCREENSIZE) >= SCREENSIZE_NORMAL {
				if fixedMySL == 0 {
					fixedMySL = SCREENSIZE_NORMAL
				}
				if fixedOSL == 0 {
					fixedOSL = SCREENSIZE_NORMAL
				}
			}
			if fixedMySL == fixedOSL {
				return mySL != 0
			} else {
				return fixedMySL > fixedOSL
			}
		}

		if ((c.ScreenLayout^o.ScreenLayout)&MASK_SCREENLONG) != 0 &&
			(r.ScreenLayout&MASK_SCREENLONG) != 0 {
			return (c.ScreenLayout & MASK_SCREENLONG) != 0
		}
	}

	// orientation
	if c.Orientation != o.Orientation || r.Orientation != 0 {
		return c.Orientation != 0
	}

	// uimode
	if c.UIMode != 0 || o.UIMode != 0 {
		diff := c.UIMode ^ o.UIMode
		if (diff&MASK_UI_MODE_TYPE) != 0 && (r.UIMode&MASK_UI_MODE_TYPE) != 0 {
			return (c.UIMode & MASK_UI_MODE_TYPE) != 0
		}
		if (diff&MASK_UI_MODE_NIGHT) != 0 && (r.UIMode&MASK_UI_MODE_NIGHT) != 0 {
			return (c.UIMode & MASK_UI_MODE_NIGHT) != 0
		}
	}

	// TODO: screen type

	// TODO: input

	// screen size
	if c.ScreenWidth != 0 || c.ScreenHeight != 0 || o.ScreenWidth != 0 || o.ScreenHeight != 0 {
		if c.ScreenWidth != o.ScreenWidth && r.ScreenWidth != 0 {
			return c.ScreenWidth != 0
		}
		if c.ScreenHeight != o.ScreenHeight && r.ScreenHeight != 0 {
			return c.ScreenHeight != 0
		}
	}

	// version
	if c.SDKVersion != 0 || c.SDKVersion != 0 || o.MinorVersion != 0 || o.MinorVersion != 0 {
		if c.SDKVersion != o.SDKVersion && r.SDKVersion != 0 {
			return c.SDKVersion > o.SDKVersion
		}
		if c.MinorVersion != o.MinorVersion {
			return c.MinorVersion != 0
		}
	}

	return false
}

func (c *ResTableConfig) Match(settings *ResTableConfig) bool {
	// match imsi
	if settings.Mcc == 0 {
		if c.Mcc != 0 {
			return false
		}
	} else {
		if c.Mcc != 0 && c.Mcc != settings.Mcc {
			return false
		}
	}
	if settings.Mnc == 0 {
		if c.Mnc != 0 {
			return false
		}
	} else {
		if c.Mnc != 0 && c.Mnc != settings.Mnc {
			return false
		}
	}

	// match locale
	if settings.Language[0] != 0 && c.Language[0] != 0 &&
		!(settings.Language[0] == c.Language[0] && settings.Language[1] == c.Language[1]) {
		return false
	}
	if settings.Country[0] != 0 && c.Country[0] != 0 &&
		!(settings.Country[0] == c.Country[0] && settings.Country[1] == c.Country[1]) {
		return false
	}

	// screen layout
	layoutDir := c.ScreenLayout & MASK_LAYOUTDIR
	setLayoutDir := settings.ScreenLayout & MASK_LAYOUTDIR
	if layoutDir != 0 && layoutDir != setLayoutDir {
		return false
	}

	screenSize := c.ScreenLayout & MASK_SCREENSIZE
	setScreenSize := settings.ScreenLayout & MASK_SCREENSIZE
	if screenSize != 0 && screenSize > setScreenSize {
		return false
	}

	screenLong := c.ScreenLayout & MASK_SCREENLONG
	setScreenLong := settings.ScreenLayout & MASK_SCREENLONG
	if screenLong != 0 && screenLong != setScreenLong {
		return false
	}

	// ui mode
	uiModeType := c.UIMode & MASK_UI_MODE_TYPE
	setUIModeType := settings.UIMode & MASK_UI_MODE_TYPE
	if uiModeType != 0 && uiModeType != setUIModeType {
		return false
	}

	// smallest screen width dp
	if c.SmallestScreenWidthDp != 0 &&
		c.SmallestScreenWidthDp > settings.SmallestScreenWidthDp {
		return false
	}

	// TODO: input

	// screen size
	if settings.ScreenWidth != 0 && c.ScreenWidth != 0 &&
		c.ScreenWidth != settings.ScreenWidth {
		return false
	}
	if settings.ScreenHeight != 0 && c.ScreenHeight != 0 &&
		c.ScreenHeight != settings.ScreenHeight {
		return false
	}

	// version
	if settings.SDKVersion != 0 && c.SDKVersion != 0 &&
		c.SDKVersion > settings.SDKVersion {
		return false
	}
	if settings.MinorVersion != 0 && c.MinorVersion != 0 &&
		c.MinorVersion != settings.MinorVersion {
		return false
	}

	return true
}

func (c *ResTableConfig) Locale() string {
	if c.Language[0] == 0 {
		return ""
	}
	if c.Country[0] == 0 {
		return fmt.Sprintf("%c%c", c.Language[0], c.Language[1])
	}
	return fmt.Sprintf("%c%c-%c%c", c.Language[0], c.Language[1], c.Country[0], c.Country[1])
}
