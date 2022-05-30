package templates

import "text/template"

//go:generate go run ../.. bundle -package templates -o bundled.go data

var (
	// MakefileUNIX is the template for the makefile on UNIX systems like Linux and BSD
	MakefileUNIX = template.Must(template.New("Makefile").Parse(string(resourceMakefile.StaticContent)))

	// DesktopFileUNIX is the template the desktop file on UNIX systems like Linux and BSD
	DesktopFileUNIX = template.Must(template.New("DesktopFile").Parse(string(resourceAppDesktop.StaticContent)))

	// EntitlementsDarwin is a plist file that lists build entitlements for darwin releases
	EntitlementsDarwin = template.Must(template.New("Entitlements").Parse(string(resourceEntitlementsDarwinPlist.StaticContent)))

	// EntitlementsDarwinMobile is a plist file that lists build entitlements for iOS releases
	EntitlementsDarwinMobile = template.Must(template.New("EntitlementsMobile").Parse(string(resourceEntitlementsIosPlist.StaticContent)))

	// FyneMetadataInit is the metadata injecting file for fyne metadata
	FyneMetadataInit = template.Must(template.New("fyne_metadata_init.got").Parse(string(resourceFynemetadatainitGot.StaticContent)))

	// ManifestWindows is the manifest file for windows packaging
	ManifestWindows = template.Must(template.New("Manifest").Parse(string(resourceAppManifest.StaticContent)))

	// AppxManifestWindows is the manifest file for windows packaging
	AppxManifestWindows = template.Must(template.New("ReleaseManifest").Parse(string(resourceAppxmanifestXML.StaticContent)))

	// InfoPlistDarwin is the manifest file for darwin packaging
	InfoPlistDarwin = template.Must(template.New("InfoPlist").Parse(string(resourceInfoPlist.StaticContent)))

	// XCAssetsDarwin is the Contents.json file for darwin xcassets bundle
	XCAssetsDarwin = template.Must(template.New("XCAssets").Parse(string(resourceXcassetsJSON.StaticContent)))

	// IndexHTML is the index.html used to serve web package
	IndexHTML = template.Must(template.New("index.html").Parse(string(resourceIndexHtml.StaticContent)))

	// WebGLDebugJs is the content of https://raw.githubusercontent.com/KhronosGroup/WebGLDeveloperTools/b42e702487d02d5278814e0fe2e2888d234893e6/src/debug/webgl-debug.js
	WebGLDebugJs = resourceWebglDebugJs.StaticContent

	// SpinnerLigth is a spinning gif of Fyne logo with a light background
	SpinnerLight = resourceSpinnerlightGif.StaticContent

	// CSSLight is a CSS that define color for element on the web splash screen following the light theme
	CSSLight = resourceLightCss.StaticContent

	// SpinnerDark is a spinning gif of Fyne logo with a dark background
	SpinnerDark = resourceSpinnerdarkGif.StaticContent

	// CSSDark is a CSS that define color for element on the web splash screen following the dark theme
	CSSDark = resourceDarkCss.StaticContent
)
