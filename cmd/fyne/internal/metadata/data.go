package metadata

// FyneApp describes the top level metadata for building a fyne application
type FyneApp struct {
	Website string
	Details AppDetails
}

// AppDetails describes the build information, this group may be OS or arch specific
type AppDetails struct {
	Icon     string
	Name, ID string
	Version  string
	Build    int
}
