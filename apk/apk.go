package apk

import (
	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"io"
	"os"
	"strconv"

	"github.com/shogo82148/androidbinary"

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
	return nil, fmt.Errorf("File %s not found", strconv.Quote(name))
}
