package templates

import "text/template"

//go:generate fyne bundle -package templates -o bundled.go data

var (
	// MakefileUNIX is the template for the makefile on UNIX systems like Linux and BSD
	MakefileUNIX = template.Must(template.New("Makefile").Parse(string(resourceMakefile.StaticContent)))

	// DesktopFileUNIX is the template the desktop file on UNIX systems like Linux and BSD
	DesktopFileUNIX = template.Must(template.New("DesktopFile").Parse(string(resourceAppDesktop.StaticContent)))

	// ManifestWindows is the manifest file for windows packaging
	ManifestWindows = template.Must(template.New("Manifest").Parse(string(resourceAppManifest.StaticContent)))

	// PlistDarwin is the manifest file for darwin packaging
	PlistDarwin = template.Must(template.New("Manifest").Parse(string(resourceInfoPlist.StaticContent)))
)
