package androidbinary_test

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/shogo82148/androidbinary"
	"github.com/shogo82148/androidbinary/apk"
)

func ExampleNewXMLFile() {
	f, _ := os.Open("testdata/AndroidManifest.xml")
	xmlFile, err := androidbinary.NewXMLFile(f)
	if err != nil {
		panic(err)
	}

	var v apk.Manifest
	dec := xml.NewDecoder(xmlFile.Reader())
	dec.Decode(&v)

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "\t")
	enc.Encode(v)

	// Output:
	// <Manifest package="net.sorablue.shogo.FWMeasure" versionCode="1" versionName="テスト版">
	// 	<application allowTaskReparenting="false" allowBackup="false" backupAgent="" debuggable="false" description="" enabled="false" hasCode="false" hardwareAccelerated="false" icon="@0x7F020000" killAfterRestore="false" largeHeap="false" label="@0x7F040000" logo="0" manageSpaceActivity="" name="" permission="" persistent="false" process="" restoreAnyVersion="false" requiredAccountType="" restrictedAccountType="" supportsRtl="false" taskAffinity="" testOnly="false" theme="" uiOptions="" vmSafeMode="false">
	// 		<activity theme="" name="FWMeasureActivity" label="">
	// 			<intent-filter>
	// 				<action name="android.intent.action.MAIN"></action>
	// 				<category name="android.intent.category.LAUNCHER"></category>
	// 			</intent-filter>
	// 		</activity>
	// 		<activity theme="" name="MapActivity" label=""></activity>
	// 		<activity theme="" name="SettingActivity" label=""></activity>
	// 		<activity theme="" name="PlaceSettingActivity" label=""></activity>
	// 	</application>
	// 	<instrumentation name="" targetPackage="" handleProfiling="false" functionalTest="false"></instrumentation>
	// 	<uses-sdk minSdkVersion="0" targetSdkVersion="0" maxSdkVersion="0"></uses-sdk>
	// </Manifest>
}

func ExampleNewTableFile() {
	f, err := os.Open("testdata/resources.arsc")
	if err != nil {
		panic(err)
	}
	tableFile, err := androidbinary.NewTableFile(f)
	if err != nil {
		panic(err)
	}

	val, err := tableFile.GetResource(0x7f040000, &androidbinary.ResTableConfig{})
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
	// Output:
	// FireworksMeasure
}
