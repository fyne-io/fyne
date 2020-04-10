package report

// OS returns the info about the OS distro
type OS struct{}

// Summary return the summary
func (i *OS) Summary() string {
	return "OS info"
}
