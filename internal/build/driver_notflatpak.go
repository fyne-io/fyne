//go:build !flatpak

package build

// IsFlatpak is true if the binary is compiled for a Flatpak package.
const IsFlatpak = false
