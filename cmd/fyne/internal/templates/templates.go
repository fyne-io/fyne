package templates

import "text/template"

//go:generate fyne bundle -package templates -o bundled.go data

var (
	// MakefileUNIX is the template for the makefile on UNIX systems like Linux and BSD
	MakefileUNIX = template.Must(template.New("Makefile").Parse(string(resourceMakefile.StaticContent)))

	// DesktopFileUNIX is the template the desktop file on UNIX systems like Linux and BSD
	DesktopFileUNIX = template.Must(template.New("DesktopFile").Parse(string(resourceAppDesktop.StaticContent)))

	// EntitlementsDarwin is a plist file that lists build entitlements for darwin releases
	EntitlementsDarwin = template.Must(template.New("Entitlements").Parse(string(resourceEntitlementsPlist.StaticContent)))

	// ManifestWindows is the manifest file for windows packaging
	ManifestWindows = template.Must(template.New("Manifest").Parse(string(resourceAppManifest.StaticContent)))

	// AppxManifestWindows is the manifest file for windows packaging
	AppxManifestWindows = template.Must(template.New("ReleaseManifest").Parse(string(resourceAppxmanifestXML.StaticContent)))

	// PlistDarwin is the manifest file for darwin packaging
	PlistDarwin = template.Must(template.New("Manifest").Parse(string(resourceInfoPlist.StaticContent)))
)
