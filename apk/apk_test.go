package apk

import (
	_ "image/jpeg"
	_ "image/png"
	"testing"
)

func TestParseAPKFile(t *testing.T) {
	apk, err := OpenFile("testdata/helloworld.apk")
	if err != nil {
		t.Errorf("OpenFile error: %v", err)
	}
	defer apk.Close()

	icon, err := apk.Icon(nil)
	if err != nil {
		t.Errorf("Icon error: %v", err)
	}
	if icon == nil {
		t.Error("Icon is nil")
	}

	label, err := apk.Label(nil)
	if err != nil {
		t.Errorf("Label error: %v", err)
	}
	if label != "HelloWorld" {
		t.Errorf("Label is not HelloWorld: %s", label)
	}
	t.Log("app label:", label)

	if apk.PackageName() != "com.example.helloworld" {
		t.Errorf("PackageName is not com.example.helloworld: %s", apk.PackageName())
	}

	manifest := apk.Manifest()
	if manifest.SDK.Target.MustInt32() != int32(24) {
		t.Errorf("SDK target is not 24: %d", manifest.SDK.Target.MustInt32())
	}

	mainActivity, err := apk.MainActivity()
	if err != nil {
		t.Errorf("MainActivity error: %v", err)
	}
	if mainActivity != "com.example.helloworld.MainActivity" {
		t.Errorf("MainActivity is not com.example.helloworld.MainActivity: %s", mainActivity)
	}
}
