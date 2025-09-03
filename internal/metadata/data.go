package metadata

// FyneApp describes the top level metadata for building a fyne application
type FyneApp struct {
	Website     string `toml:",omitempty"`
	Details     AppDetails
	Development map[string]string `toml:",omitempty"`
	Release     map[string]string `toml:",omitempty"`
	Source      *AppSource        `toml:",omitempty"`
	LinuxAndBSD *LinuxAndBSD      `toml:",omitempty"`
	Languages   []string          `toml:",omitempty"`
	Migrations  map[string]bool   `toml:",omitempty"`
}

// AppDetails describes the build information, this group may be OS or arch specific
type AppDetails struct {
	Icon     string `toml:",omitempty"`
	Name, ID string `toml:",omitempty"`
	Version  string `toml:",omitempty"`
	Build    int    `toml:",omitempty"`
}

type AppSource struct {
	Repo, Dir string `toml:",omitempty"`
}

// LinuxAndBSD describes specific metadata for desktop files on Linux and BSD.
type LinuxAndBSD struct {
	GenericName string   `toml:",omitempty"`
	Categories  []string `toml:",omitempty"`
	Comment     string   `toml:",omitempty"`
	Keywords    []string `toml:",omitempty"`
	ExecParams  string   `toml:",omitempty"`
}
