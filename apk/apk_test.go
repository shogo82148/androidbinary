package apk

import (
	_ "image/jpeg"
	_ "image/png"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIcon(t *testing.T) {
	apk, err := OpenFile("testdata/helloworld.apk")
	assert.NoError(t, err)

	icon, err := apk.Icon(nil)
	assert.NoError(t, err)
	assert.NotNil(t, icon)

	label, err := apk.Label(nil)
	assert.NoError(t, err)
	assert.Equal(t, "HelloWorld", label)
	t.Log("app label:", label)

	assert.Equal(t, "com.example.helloworld", apk.PackageName())

	manifest := apk.Manifest()
	assert.Equal(t, manifest.SDK.Target, 24)

	mainActivity, err := apk.MainActivity()
	assert.NoError(t, err)
	assert.Equal(t, "com.example.helloworld.MainActivity", mainActivity)
}
