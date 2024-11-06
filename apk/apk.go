package apk

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/avast/apkparser"
	"github.com/shogo82148/androidbinary"
	_ "golang.org/x/image/webp"

	_ "image/jpeg" // handle jpeg format
	_ "image/png"  // handle png format
)

// Apk is an application package file for android.
type Apk struct {
	f         *os.File
	zipreader *zip.Reader
	manifest  Manifest
	table     *androidbinary.TableFile
}

// OpenFile will open the file specified by filename and return Apk
func OpenFile(filename string) (apk *Apk, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			f.Close()
		}
	}()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	apk, err = OpenZipReader(f, fi.Size())
	if err != nil {
		return nil, err
	}
	apk.f = f
	return
}

// OpenZipReader has same arguments like zip.NewReader
func OpenZipReader(r io.ReaderAt, size int64) (*Apk, error) {
	zipreader, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	apk := &Apk{
		zipreader: zipreader,
	}
	if err = apk.parseResources(); err != nil {
		return nil, err
	}
	if err = apk.parseManifest(); err != nil {
		return nil, errorf("parse-manifest: %w", err)
	}
	return apk, nil
}

// Close is avaliable only if apk is created with OpenFile
func (k *Apk) Close() error {
	if k.f == nil {
		return nil
	}
	return k.f.Close()
}

type AdaptiveIcon struct {
	Background struct {
		Drawable string `xml:"http://schemas.android.com/apk/res/android drawable,attr"`
	} `xml:"background"`
	Foreground struct {
		Drawable string `xml:"http://schemas.android.com/apk/res/android drawable,attr"`
	} `xml:"foreground"`
}

// Icon returns the icon image of the APK.
func (k *Apk) Icon(resConfig *androidbinary.ResTableConfig) (image.Image, error) {
	iconPath, err := k.manifest.App.Icon.WithResTableConfig(resConfig).String()
	if err != nil {
		return nil, err
	}
	if androidbinary.IsResID(iconPath) {
		return nil, newError("unable to convert icon-id to icon path")
	}
	imgData, err := k.readZipFile(iconPath)
	if err != nil {
		return nil, err
	}

	// is it an android XML file (maybe an adaptive-icon)?
	if filepath.Ext(iconPath) == ".xml" {
		// TODO: I absolutely do not like the below code but it worked for my sample application that used an adaptive icon; not sure if it would also work for others

		// read the raw XML data from the android xml file
		buf := new(bytes.Buffer)
		enc := xml.NewEncoder(buf)
		enc.Indent("", "    ")
		if err := apkparser.ParseXml(bytes.NewReader(imgData), enc, nil); err != nil {
			return nil, err
		}
		tmpStr := buf.String()
		log.Println(tmpStr)

		xmlfile, err := androidbinary.NewXMLFile(bytes.NewReader(imgData))
		if err != nil {
			return nil, errorf("failed to parse %q: %w", iconPath, err)
		}
		var adaptiveIcon AdaptiveIcon
		if err = xmlfile.Decode(&adaptiveIcon, k.table, nil); err != nil {
			return nil, errorf("failed to decode %q: %w", iconPath, err)
		}

		var innerForegroundImagePath string
		if androidbinary.IsResID(adaptiveIcon.Background.Drawable) {
			resID, err := androidbinary.ParseResID(adaptiveIcon.Foreground.Drawable)
			if err != nil {
				return nil, errorf("failed to parse resID %q: %w", adaptiveIcon.Foreground.Drawable, err)
			}

			innerForegroundImagePathTmp, err := k.table.GetResource(resID, nil)
			switch v := innerForegroundImagePathTmp.(type) {
			case string:
				innerForegroundImagePath = v
			default:
				return nil, errorf("failed to get resource %q: %w", adaptiveIcon.Foreground.Drawable, err)
			}

			innerForegroundImagePath = innerForegroundImagePathTmp.(string)
			if err != nil {
				return nil, errorf("failed to get resource %q: %w", adaptiveIcon.Foreground.Drawable, err)
			}
		}

		if innerForegroundImagePath == "" {
			return nil, errorf("failed to find inner foreground in %q", iconPath)
		}

		// read from the inner foreground image location
		imgData, err = k.readZipFile(innerForegroundImagePath)
		if err != nil {
			return nil, errorf("failed to read %q: %w", adaptiveIcon.Foreground, err)
		}
	}
	m, _, err := image.Decode(bytes.NewReader(imgData))
	return m, err
}

// Label returns the label of the APK.
func (k *Apk) Label(resConfig *androidbinary.ResTableConfig) (s string, err error) {
	s, err = k.manifest.App.Label.WithResTableConfig(resConfig).String()
	if err != nil {
		return
	}
	if androidbinary.IsResID(s) {
		err = newError("unable to convert label-id to string")
	}
	return
}

// Manifest returns the manifest of the APK.
func (k *Apk) Manifest() Manifest {
	return k.manifest
}

// PackageName returns the package name of the APK.
func (k *Apk) PackageName() string {
	return k.manifest.Package.MustString()
}

func isMainIntentFilter(intent ActivityIntentFilter) bool {
	ok := false
	for _, action := range intent.Actions {
		s, err := action.Name.String()
		if err == nil && s == "android.intent.action.MAIN" {
			ok = true
			break
		}
	}
	if !ok {
		return false
	}
	ok = false
	for _, category := range intent.Categories {
		s, err := category.Name.String()
		if err == nil && s == "android.intent.category.LAUNCHER" {
			ok = true
			break
		}
	}
	return ok
}

// MainActivity returns the name of the main activity.
func (k *Apk) MainActivity() (activity string, err error) {
	for _, act := range k.manifest.App.Activities {
		for _, intent := range act.IntentFilters {
			if isMainIntentFilter(intent) {
				return act.Name.String()
			}
		}
	}
	for _, act := range k.manifest.App.ActivityAliases {
		for _, intent := range act.IntentFilters {
			if isMainIntentFilter(intent) {
				return act.TargetActivity.String()
			}
		}
	}

	return "", newError("No main activity found")
}

func (k *Apk) parseManifest() error {
	xmlData, err := k.readZipFile("AndroidManifest.xml")
	if err != nil {
		return errorf("failed to read AndroidManifest.xml: %w", err)
	}
	xmlfile, err := androidbinary.NewXMLFile(bytes.NewReader(xmlData))
	if err != nil {
		return errorf("failed to parse AndroidManifest.xml: %w", err)
	}
	return xmlfile.Decode(&k.manifest, k.table, nil)
}

func (k *Apk) parseResources() (err error) {
	resData, err := k.readZipFile("resources.arsc")
	if err != nil {
		return
	}
	k.table, err = androidbinary.NewTableFile(bytes.NewReader(resData))
	return
}

func (k *Apk) readZipFile(name string) (data []byte, err error) {
	buf := bytes.NewBuffer(nil)
	for _, file := range k.zipreader.File {
		if file.Name != name {
			continue
		}
		rc, er := file.Open()
		if er != nil {
			err = er
			return
		}
		defer rc.Close()
		_, err = io.Copy(buf, rc)
		if err != nil {
			return
		}
		return buf.Bytes(), nil
	}
	return nil, fmt.Errorf("apk: file %q not found", name)
}
