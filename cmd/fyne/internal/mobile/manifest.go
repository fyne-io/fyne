// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mobile

import (
	"encoding/xml"
	"errors"
	"fmt"
	"html/template"
)

type manifestXML struct {
	Activity activityXML `xml:"application>activity"`
}

type activityXML struct {
	Name     string        `xml:"name,attr"`
	MetaData []metaDataXML `xml:"meta-data"`
}

type metaDataXML struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// manifestLibName parses the AndroidManifest.xml and finds the library
// name of the NativeActivity.
func manifestLibName(data []byte) (string, error) {
	manifest := new(manifestXML)
	if err := xml.Unmarshal(data, manifest); err != nil {
		return "", err
	}
	if manifest.Activity.Name != "org.golang.app.GoNativeActivity" {
		return "", fmt.Errorf("can only build an .apk for GoNativeActivity, not %q", manifest.Activity.Name)
	}
	libName := ""
	for _, md := range manifest.Activity.MetaData {
		if md.Name == "android.app.lib_name" {
			libName = md.Value
			break
		}
	}
	if libName == "" {
		return "", errors.New("AndroidManifest.xml missing meta-data android.app.lib_name")
	}
	return libName, nil
}

type manifestTmplData struct {
	JavaPkgPath string
	Name        string
	Debug       bool
	LibName     string
	Version     string
	Build       int
}

var manifestTmpl = template.Must(template.New("manifest").Parse(`
<manifest
	xmlns:android="http://schemas.android.com/apk/res/android"
	package="{{.JavaPkgPath}}"
	android:versionCode="{{.Build}}"
	android:versionName="{{.Version}}">

	<application android:label="{{.Name}}" android:debuggable="{{.Debug}}">
	<activity android:name="org.golang.app.GoNativeActivity"
		android:label="{{.Name}}"
		android:configChanges="orientation|keyboardHidden|uiMode"
		android:theme="@android:style/Theme">
		<meta-data android:name="android.app.lib_name" android:value="{{.LibName}}" />
		<intent-filter>
			<action android:name="android.intent.action.MAIN" />
			<category android:name="android.intent.category.LAUNCHER" />
		</intent-filter>
	</activity>
	</application>

	<uses-permission android:name="android.permission.WRITE_EXTERNAL_STORAGE" />
	<uses-permission android:name="android.permission.READ_EXTERNAL_STORAGE" />
	<uses-permission android:name="android.permission.INTERNET" />
</manifest>`))
