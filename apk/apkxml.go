package apk

import (
	"github.com/shogo82148/androidbinary"
)

// Instrumentation is an application instrumentation code.
type Instrumentation struct {
	Name            androidbinary.String `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Target          androidbinary.String `xml:"http://schemas.android.com/apk/res/android targetPackage,attr"`
	HandleProfiling androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android handleProfiling,attr"`
	FunctionalTest  androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android functionalTest,attr"`
}

// ActivityAction is an action of an activity.
type ActivityAction struct {
	Name androidbinary.String `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

// ActivityCategory is a category of an activity.
type ActivityCategory struct {
	Name androidbinary.String `xml:"http://schemas.android.com/apk/res/android name,attr"`
}

// ActivityIntentFilter is an androidbinary.Int32ent filter of an activity.
type ActivityIntentFilter struct {
	Actions    []ActivityAction   `xml:"action"`
	Categories []ActivityCategory `xml:"category"`
}

// AppActivity is an activity in an application.
type AppActivity struct {
	Theme             androidbinary.String   `xml:"http://schemas.android.com/apk/res/android theme,attr"`
	Name              androidbinary.String   `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Label             androidbinary.String   `xml:"http://schemas.android.com/apk/res/android label,attr"`
	ScreenOrientation androidbinary.String   `xml:"http://schemas.android.com/apk/res/android screenOrientation,attr"`
	IntentFilters     []ActivityIntentFilter `xml:"intent-filter"`
}

// AppActivityAlias https://developer.android.com/guide/topics/manifest/activity-alias-element
type AppActivityAlias struct {
	Name           androidbinary.String   `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Label          androidbinary.String   `xml:"http://schemas.android.com/apk/res/android label,attr"`
	TargetActivity androidbinary.String   `xml:"http://schemas.android.com/apk/res/android targetActivity,attr"`
	IntentFilters  []ActivityIntentFilter `xml:"intent-filter"`
}

// MetaData is a metadata in an application.
type MetaData struct {
	Name  androidbinary.String `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Value androidbinary.String `xml:"http://schemas.android.com/apk/res/android value,attr"`
}

// Application is an application in an APK.
type Application struct {
	AllowTaskReparenting  androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android allowTaskReparenting,attr"`
	AllowBackup           androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android allowBackup,attr"`
	BackupAgent           androidbinary.String `xml:"http://schemas.android.com/apk/res/android backupAgent,attr"`
	Debuggable            androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android debuggable,attr"`
	Description           androidbinary.String `xml:"http://schemas.android.com/apk/res/android description,attr"`
	Enabled               androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android enabled,attr"`
	HasCode               androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android hasCode,attr"`
	HardwareAccelerated   androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android hardwareAccelerated,attr"`
	Icon                  androidbinary.String `xml:"http://schemas.android.com/apk/res/android icon,attr"`
	KillAfterRestore      androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android killAfterRestore,attr"`
	LargeHeap             androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android largeHeap,attr"`
	Label                 androidbinary.String `xml:"http://schemas.android.com/apk/res/android label,attr"`
	Logo                  androidbinary.String `xml:"http://schemas.android.com/apk/res/android logo,attr"`
	ManageSpaceActivity   androidbinary.String `xml:"http://schemas.android.com/apk/res/android manageSpaceActivity,attr"`
	Name                  androidbinary.String `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Permission            androidbinary.String `xml:"http://schemas.android.com/apk/res/android permission,attr"`
	Persistent            androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android persistent,attr"`
	Process               androidbinary.String `xml:"http://schemas.android.com/apk/res/android process,attr"`
	RestoreAnyVersion     androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android restoreAnyVersion,attr"`
	RequiredAccountType   androidbinary.String `xml:"http://schemas.android.com/apk/res/android requiredAccountType,attr"`
	RestrictedAccountType androidbinary.String `xml:"http://schemas.android.com/apk/res/android restrictedAccountType,attr"`
	SupportsRtl           androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android supportsRtl,attr"`
	TaskAffinity          androidbinary.String `xml:"http://schemas.android.com/apk/res/android taskAffinity,attr"`
	TestOnly              androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android testOnly,attr"`
	Theme                 androidbinary.String `xml:"http://schemas.android.com/apk/res/android theme,attr"`
	UIOptions             androidbinary.String `xml:"http://schemas.android.com/apk/res/android uiOptions,attr"`
	VMSafeMode            androidbinary.Bool   `xml:"http://schemas.android.com/apk/res/android vmSafeMode,attr"`
	Activities            []AppActivity        `xml:"activity"`
	ActivityAliases       []AppActivityAlias   `xml:"activity-alias"`
	MetaData              []MetaData           `xml:"meta-data"`
}

// UsesSDK is target SDK version.
type UsesSDK struct {
	Min    androidbinary.Int32 `xml:"http://schemas.android.com/apk/res/android minSdkVersion,attr"`
	Target androidbinary.Int32 `xml:"http://schemas.android.com/apk/res/android targetSdkVersion,attr"`
	Max    androidbinary.Int32 `xml:"http://schemas.android.com/apk/res/android maxSdkVersion,attr"`
}

// UsesPermission is user grant the system permission.
type UsesPermission struct {
	Name androidbinary.String `xml:"http://schemas.android.com/apk/res/android name,attr"`
	Max  androidbinary.Int32  `xml:"http://schemas.android.com/apk/res/android maxSdkVersion,attr"`
}

// Manifest is a manifest of an APK.
type Manifest struct {
	Package         androidbinary.String `xml:"package,attr"`
	VersionCode     androidbinary.Int32  `xml:"http://schemas.android.com/apk/res/android versionCode,attr"`
	VersionName     androidbinary.String `xml:"http://schemas.android.com/apk/res/android versionName,attr"`
	App             Application          `xml:"application"`
	Instrument      Instrumentation      `xml:"instrumentation"`
	SDK             UsesSDK              `xml:"uses-sdk"`
	UsesPermissions []UsesPermission     `xml:"uses-permission"`
}
