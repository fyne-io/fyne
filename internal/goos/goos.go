package goos

// OS constants (to compare to runtime.GOOS)
const (
	Android    = "android"
	Darwin     = "darwin"
	FreeBSD    = "freebsd"
	IOS        = "ios"
	JavaScript = "js"
	Linux      = "linux"
	NetBSD     = "netbsd"
	OpenBSD    = "openbsd"
	Windows    = "windows"
)

// IsBSD returns whether the specified OS is a supported BSD variant.
func IsBSD(os string) bool {
	return os == FreeBSD || os == NetBSD || os == OpenBSD
}
