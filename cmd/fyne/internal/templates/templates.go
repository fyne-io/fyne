package templates

import "text/template"

//go:generate fyne bundle -package templates -o bundled.go data

var (
	// MakefileUNIX is the template for the makefile on UNIX systems like Linux and BSD
	MakefileUNIX = template.Must(template.New("Makefile").Parse(string(resourceMakefile.StaticContent)))

	// DesktopFileUNIX is the template the desktop file on UNIX systems like Linux and BSD
	DesktopFileUNIX = template.Must(template.New("DesktopFile").Parse(string(resourceAppDesktop.StaticContent)))

	// EntitlementsDarwin is a plist file that lists build entitlements for darwin releases
	EntitlementsDarwin = template.Must(template.New("Entitlements").Parse(string(resourceEntitlementsDarwinPlist.StaticContent)))

	// EntitlementsDarwinMobile is a plist file that lists build entitlements for iOS releases
	EntitlementsDarwinMobile = template.Must(template.New("EntitlementsMobile").Parse(string(resourceEntitlementsIosPlist.StaticContent)))

	// ManifestWindows is the manifest file for windows packaging
	ManifestWindows = template.Must(template.New("Manifest").Parse(string(resourceAppManifest.StaticContent)))

	// AppxManifestWindows is the manifest file for windows packaging
	AppxManifestWindows = template.Must(template.New("ReleaseManifest").Parse(string(resourceAppxmanifestXML.StaticContent)))

	// InfoPlistDarwin is the manifest file for darwin packaging
	InfoPlistDarwin = template.Must(template.New("InfoPlist").Parse(string(resourceInfoPlist.StaticContent)))

	// XCAssetsDarwin is the Contents.json file for darwin xcassets bundle
	XCAssetsDarwin = template.Must(template.New("XCAssets").Parse(string(resourceXcassetsJSON.StaticContent)))
)
