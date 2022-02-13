package androidbinary

import (
	"bytes"
	"encoding/xml"
	"io"
	"os"
	"reflect"
	"testing"
)

type XMLManifest struct {
	XMLName         xml.Name             `xml:"manifest"`
	VersionName     string               `xml:"http://schemas.android.com/apk/res/android versionName,attr"`
	VersionCode     string               `xml:"http://schemas.android.com/apk/res/android versionCode,attr"`
	Package         string               `xml:"package,attr"`
	UsesPermissions []*XMLUsesPermission `xml:"uses-permission"`
	Applications    []*XMLApplication    `xml:"application"`
}

type XMLUsesPermission struct {
	XMLName xml.Name `xml:"uses-permission"`
	Name    string   `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

type XMLApplication struct {
	XMLName       xml.Name          `xml:"application"`
	Label         string            `xml:"http://schemas.android.com/apk/res/android label,attr"`
	Icon          string            `xml:"http://schemas.android.com/apk/res/android icon,attr"`
	Debuggable    string            `xml:"http://schemas.android.com/apk/res/android debuggable,attr"`
	UsesLibraries []*XMLUsesLibrary `xml:"uses-library"`
	Activities    []*XMLActivity    `xml:"activity"`
}

type XMLUsesLibrary struct {
	XMLName xml.Name `xml:"uses-library"`
	Name    string   `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

type XMLActivity struct {
	XMLName           xml.Name `xml:"activity"`
	Name              string   `xml:"http://schemas.android.com/apk/res/android name,attr"`
	ScreenOrientation string   `xml:"http://schemas.android.com/apk/res/android screenOrientation,attr"`
}

func TestNewXMLFile(t *testing.T) {
	f, _ := os.Open("testdata/AndroidManifest.xml")
	xmlFile, err := NewXMLFile(f)
	if err != nil {
		t.Errorf("got %v want no error", err)
	}
	decoder := xml.NewDecoder(xmlFile.Reader())
	xmlManifest := &XMLManifest{}
	err = decoder.Decode(xmlManifest)
	if err != nil {
		t.Errorf("got %v want no error", err)
	}
	expected := &XMLManifest{
		XMLName:     xml.Name{Local: "manifest"},
		VersionName: "テスト版",
		VersionCode: "1",
		Package:     "net.sorablue.shogo.FWMeasure",
		UsesPermissions: []*XMLUsesPermission{
			{XMLName: xml.Name{Local: "uses-permission"}, Name: "android.permission.CAMERA"},
			{XMLName: xml.Name{Local: "uses-permission"}, Name: "android.permission.WAKE_LOCK"},
			{XMLName: xml.Name{Local: "uses-permission"}, Name: "android.permission.ACCESS_FINE_LOCATION"},
			{XMLName: xml.Name{Local: "uses-permission"}, Name: "android.permission.INTERNET"},
			{XMLName: xml.Name{Local: "uses-permission"}, Name: "android.permission.ACCESS_MOCK_LOCATION"},
			{XMLName: xml.Name{Local: "uses-permission"}, Name: "android.permission.RECORD_AUDIO"},
		},
		Applications: []*XMLApplication{
			{
				XMLName:    xml.Name{Local: "application"},
				Label:      "@0x7F040000",
				Icon:       "@0x7F020000",
				Debuggable: "false",
				UsesLibraries: []*XMLUsesLibrary{
					{XMLName: xml.Name{Local: "uses-library"}, Name: "com.google.android.maps"},
				},
				Activities: []*XMLActivity{
					{
						XMLName:           xml.Name{Local: "activity"},
						ScreenOrientation: "0",
						Name:              "FWMeasureActivity",
					},
					{
						XMLName:           xml.Name{Local: "activity"},
						ScreenOrientation: "0",
						Name:              "MapActivity",
					},
					{
						XMLName: xml.Name{Local: "activity"},
						Name:    "SettingActivity",
					},
					{
						XMLName: xml.Name{Local: "activity"},
						Name:    "PlaceSettingActivity",
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(xmlManifest, expected) {
		t.Errorf("got %v want %v", xmlManifest, expected)
	}
}

func TestReadStartNamespace(t *testing.T) {
	input := []uint8{
		0x00, 0x01, // Type = RES_XML_START_NAMESPACE_TYPE
		0x10, 0x00, // HeadSize = 16 bytes
		0x18, 0x00, 0x00, 0x00, // Size
		0x01, 0x00, 0x00, 0x00, // LineNumber = 1
		0xFF, 0xFF, 0xFF, 0xFF, // Comment is none
		0x02, 0x00, 0x00, 0x00, // Prefix = 2
		0x01, 0x00, 0x00, 0x00, // Uri = 1
	}
	buf := bytes.NewReader(input)
	sr := io.NewSectionReader(buf, 0, int64(len(input)))
	f := new(XMLFile)
	err := f.readStartNamespace(sr)
	if err != nil {
		t.Errorf("got %v want no error", err)
	}
	if f.notPrecessedNS[ResStringPoolRef(1)] != ResStringPoolRef(2) {
		t.Errorf("got %v want %v", f.notPrecessedNS[ResStringPoolRef(1)], ResStringPoolRef(2))
	}
	if f.namespaces.get(ResStringPoolRef(1)) != ResStringPoolRef(2) {
		t.Errorf("got %v want %v", f.namespaces.get(ResStringPoolRef(1)), ResStringPoolRef(2))
	}
}

func TestReadEndNamespace(t *testing.T) {
	input := []uint8{
		0x01, 0x01, // Type = RES_XML_END_NAMESPACE_TYPE
		0x10, 0x00, // HeadSize = 16 bytes
		0x18, 0x00, 0x00, 0x00, // Size
		0x01, 0x00, 0x00, 0x00, // LineNumber = 1
		0xFF, 0xFF, 0xFF, 0xFF, // Comment is none
		0x02, 0x00, 0x00, 0x00, // Prefix = 2
		0x01, 0x00, 0x00, 0x00, // Uri = 1
	}
	buf := bytes.NewReader(input)
	sr := io.NewSectionReader(buf, 0, int64(len(input)))

	f := new(XMLFile)
	f.namespaces.add(ResStringPoolRef(1), ResStringPoolRef(2))
	err := f.readEndNamespace(sr)

	if err != nil {
		t.Errorf("got %v want no error", err)
	}
}

func TestAddNamespacePrefix(t *testing.T) {
	nameRef := ResStringPoolRef(1)
	prefixRef := ResStringPoolRef(2)
	uriRef := ResStringPoolRef(3)

	f := new(XMLFile)
	f.namespaces.add(uriRef, prefixRef)
	f.stringPool = new(ResStringPool)
	f.stringPool.Strings = []string{"", "name", "prefix", "http://example.com"}

	if actual := f.addNamespacePrefix(NilResStringPoolRef, nameRef); actual != "name" {
		t.Errorf("got %v want name", actual)
	}

	if actual := f.addNamespacePrefix(uriRef, nameRef); actual != "prefix:name" {
		t.Errorf("got %v want prefix:name", actual)
	}
}

func TestReadStartElement(t *testing.T) {
	input := []uint8{
		0x02, 0x01, // Type = RES_XML_START_ELEMENT_TYPE
		0x10, 0x00, // HeadSize = 16 bytes
		0x3C, 0x00, 0x00, 0x00, // Size = 60 bytes
		0x01, 0x00, 0x00, 0x00, // LineNumber = 1
		0xFF, 0xFF, 0xFF, 0xFF, // Comment is none
		0xFF, 0xFF, 0xFF, 0xFF, // Namespace is none
		0x01, 0x00, 0x00, 0x00, // Name = 1
		0x14, 0x00, // AttributeStart = 20 bytes
		0x14, 0x00, // AttributeSize = 20 bytes
		0x01, 0x00, // AttributeCount = 1
		0x00, 0x00, // IdIndex = 0
		0x00, 0x00, // ClassIndex = 0
		0x00, 0x00, // StyleIndex = 0

		// Attributes
		0xFF, 0xFF, 0xFF, 0xFF, // Namespace is none
		0x04, 0x00, 0x00, 0x00, // Name is 'attr'
		0x05, 0x00, 0x00, 0x00, // RawValue is 'value'
		0x08, 0x00, // size = 8
		0x00,                   // padding
		0x03,                   // data type is TYPE_STRING
		0x05, 0x00, 0x00, 0x00, // data
	}
	buf := bytes.NewReader(input)
	sr := io.NewSectionReader(buf, 0, int64(len(input)))

	prefixRef := ResStringPoolRef(2)
	uriRef := ResStringPoolRef(3)

	f := new(XMLFile)
	f.notPrecessedNS = make(map[ResStringPoolRef]ResStringPoolRef)
	f.notPrecessedNS[uriRef] = prefixRef
	f.namespaces.add(uriRef, prefixRef)
	f.stringPool = new(ResStringPool)
	f.stringPool.Strings = []string{"", "name", "prefix", "http://example.com", "attr", "value"}
	err := f.readStartElement(sr)

	if err != nil {
		t.Errorf("got %v want no error", err)
	}

	if actual := f.xmlBuffer.String(); actual != "<name xmlns:prefix=\"http://example.com\" attr=\"value\">" {
		t.Errorf("got %v want <name xmlns:prefix=\"http://example.com\" attr=\"value\">", actual)
	}
}

func TestReadEndElement(t *testing.T) {
	input := []uint8{
		0x03, 0x01, // Type = RES_XML_END_ELEMENT_TYPE
		0x10, 0x00, // HeadSize = 16 bytes
		0x18, 0x00, 0x00, 0x00, // Size
		0x01, 0x00, 0x00, 0x00, // LineNumber = 1
		0xFF, 0xFF, 0xFF, 0xFF, // Comment is none
		0xFF, 0xFF, 0xFF, 0xFF, // Namespace is none
		0x01, 0x00, 0x00, 0x00, // Name = 1
	}
	buf := bytes.NewReader(input)
	sr := io.NewSectionReader(buf, 0, int64(len(input)))

	f := new(XMLFile)
	f.stringPool = new(ResStringPool)
	f.stringPool.Strings = []string{"", "name"}
	err := f.readEndElement(sr)

	if err != nil {
		t.Errorf("got %v want no error", err)
	}

	if actual := f.xmlBuffer.String(); actual != "</name>" {
		t.Errorf("got %v want </name>", actual)
	}
}
