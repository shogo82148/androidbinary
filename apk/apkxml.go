package apk

import (
	"github.com/shogo82148/androidbinary"
)

// Instrumentation is an application instrumentation code.
type Instrumentation struct {
	Name            androidbinary.String `xml:"name,attr"`
	Target          androidbinary.String `xml:"targetPackage,attr"`
	HandleProfiling androidbinary.Bool   `xml:"handleProfiling,attr"`
	FunctionalTest  androidbinary.Bool   `xml:"functionalTest,attr"`
}

// ActivityAction is an action of an activity.
type ActivityAction struct {
	Name androidbinary.String `xml:"name,attr"`
}

// ActivityCategory is a category of an activity.
type ActivityCategory struct {
	Name androidbinary.String `xml:"name,attr"`
}

// ActivityIntentFilter is an androidbinary.Int32ent filter of an activity.
type ActivityIntentFilter struct {
	Actions    []ActivityAction   `xml:"action"`
	Categories []ActivityCategory `xml:"category"`
}

// AppActivity is an activity in an application.
type AppActivity struct {
	Theme             androidbinary.String   `xml:"theme,attr"`
	Name              androidbinary.String   `xml:"name,attr"`
	Label             androidbinary.String   `xml:"label,attr"`
	ScreenOrientation androidbinary.String   `xml:"screenOrientation,attr"`
	IntentFilters     []ActivityIntentFilter `xml:"intent-filter"`
}

// AppActivityAlias https://developer.android.com/guide/topics/manifest/activity-alias-element
type AppActivityAlias struct {
	Name           androidbinary.String   `xml:"name,attr"`
	Label          androidbinary.String   `xml:"label,attr"`
	TargetActivity androidbinary.String   `xml:"targetActivity,attr"`
	IntentFilters  []ActivityIntentFilter `xml:"intent-filter"`
}

// MetaData is a metadata in an application.
type MetaData struct {
	Name  androidbinary.String `xml:"name,attr"`
	Value androidbinary.String `xml:"value,attr"`
}

// Application is an application in an APK.
type Application struct {
	AllowTaskReparenting  androidbinary.Bool   `xml:"allowTaskReparenting,attr"`
	AllowBackup           androidbinary.Bool   `xml:"allowBackup,attr"`
	BackupAgent           androidbinary.String `xml:"backupAgent,attr"`
	Debuggable            androidbinary.Bool   `xml:"debuggable,attr"`
	Description           androidbinary.String `xml:"description,attr"`
	Enabled               androidbinary.Bool   `xml:"enabled,attr"`
	HasCode               androidbinary.Bool   `xml:"hasCode,attr"`
	HardwareAccelerated   androidbinary.Bool   `xml:"hardwareAccelerated,attr"`
	Icon                  androidbinary.String `xml:"icon,attr"`
	KillAfterRestore      androidbinary.Bool   `xml:"killAfterRestore,attr"`
	LargeHeap             androidbinary.Bool   `xml:"largeHeap,attr"`
	Label                 androidbinary.String `xml:"label,attr"`
	Logo                  androidbinary.String `xml:"logo,attr"`
	ManageSpaceActivity   androidbinary.String `xml:"manageSpaceActivity,attr"`
	Name                  androidbinary.String `xml:"name,attr"`
	Permission            androidbinary.String `xml:"permission,attr"`
	Persistent            androidbinary.Bool   `xml:"persistent,attr"`
	Process               androidbinary.String `xml:"process,attr"`
	RestoreAnyVersion     androidbinary.Bool   `xml:"restoreAnyVersion,attr"`
	RequiredAccountType   androidbinary.String `xml:"requiredAccountType,attr"`
	RestrictedAccountType androidbinary.String `xml:"restrictedAccountType,attr"`
	SupportsRtl           androidbinary.Bool   `xml:"supportsRtl,attr"`
	TaskAffinity          androidbinary.String `xml:"taskAffinity,attr"`
	TestOnly              androidbinary.Bool   `xml:"testOnly,attr"`
	Theme                 androidbinary.String `xml:"theme,attr"`
	UIOptions             androidbinary.String `xml:"uiOptions,attr"`
	VMSafeMode            androidbinary.Bool   `xml:"vmSafeMode,attr"`
	Activities            []AppActivity        `xml:"activity"`
	ActivityAliases       []AppActivityAlias   `xml:"activity-alias"`
	MetaData              []MetaData           `xml:"meta-data"`
}

// UsesSDK is target SDK version.
type UsesSDK struct {
	Min    androidbinary.Int32 `xml:"minSdkVersion,attr"`
	Target androidbinary.Int32 `xml:"targetSdkVersion,attr"`
	Max    androidbinary.Int32 `xml:"maxSdkVersion,attr"`
}

// UsesPermission is user grant the system permission.
type UsesPermission struct {
	Name androidbinary.String `xml:"name,attr"`
	Max  androidbinary.Int32  `xml:"maxSdkVersion,attr"`
}

// Manifest is a manifest of an APK.
type Manifest struct {
	Package         androidbinary.String `xml:"package,attr"`
	VersionCode     androidbinary.Int32  `xml:"versionCode,attr"`
	VersionName     androidbinary.String `xml:"versionName,attr"`
	App             Application          `xml:"application"`
	Instrument      Instrumentation      `xml:"instrumentation"`
	SDK             UsesSDK              `xml:"uses-sdk"`
	UsesPermissions []UsesPermission     `xml:"uses-permission"`
}
