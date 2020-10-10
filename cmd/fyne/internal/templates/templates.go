package templates

import "text/template"

var (
	// MakefileUNIX is the template for the makefile on UNIX systems like Linux and BSD
	MakefileUNIX = template.Must(template.New("Makefile").Parse(string(makefile.StaticContent)))

	// DesktopFileUNIX is the template the desktop file on UNIX systems like Linux and BSD
	DesktopFileUNIX = template.Must(template.New("DesktopFile").Parse(string(desktopFile.StaticContent)))

	// ManifestWindows is the manifest file for windows packaging
	ManifestWindows = template.Must(template.New("Manifest").Parse(string(manifest.StaticContent)))

	// PlistDarwin is the manifest file for darwin packaging
	PlistDarwin = template.Must(template.New("Manifest").Parse(string(plist.StaticContent)))
)
