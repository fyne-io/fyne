//go:build !windows

package glfw

func isDark() bool {
	return true // this is really a no-op placeholder for a windows menu workaround
}
