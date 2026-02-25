package dialog

import (
	"os"
	"syscall"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
)

func driveMask() uint32 {
	dll, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		fyne.LogError("Error loading kernel32.dll", err)
		return 0
	}
	handle, err := syscall.GetProcAddress(dll, "GetLogicalDrives")
	if err != nil {
		fyne.LogError("Could not find GetLogicalDrives call", err)
		return 0
	}

	ret, _, err := syscall.SyscallN(uintptr(handle))
	if err != syscall.Errno(0) { // for some reason Syscall returns something not nil on success
		fyne.LogError("Error calling GetLogicalDrives", err)
		return 0
	}

	return uint32(ret)
}

func listDrives() []string {
	var drives []string
	mask := driveMask()

	for i := 0; i < 26; i++ {
		if mask&1 == 1 {
			letter := string('A' + rune(i))
			drives = append(drives, letter+":")
		}
		mask >>= 1
	}

	return drives
}

func (f *fileDialog) getPlaces() []favoriteItem {
	drives := listDrives()
	places := make([]favoriteItem, len(drives))
	for i, drive := range drives {
		driveRoot := drive + string(os.PathSeparator) // capture loop var
		driveRootURI, _ := storage.ListerForURI(storage.NewFileURI(driveRoot))
		places[i] = favoriteItem{
			drive,
			theme.StorageIcon(),
			driveRootURI,
		}
	}
	return places
}

func isHidden(file fyne.URI) bool {
	if file.Scheme() != "file" {
		fyne.LogError("Cannot check if non file is hidden", nil)
		return false
	}

	point, err := syscall.UTF16PtrFromString(file.Path())
	if err != nil {
		fyne.LogError("Error making string pointer", err)
		return false
	}
	attr, err := syscall.GetFileAttributes(point)
	if err != nil {
		fyne.LogError("Error getting file attributes", err)
		return false
	}

	return attr&syscall.FILE_ATTRIBUTE_HIDDEN != 0
}

func hideFile(filename string) (err error) {
	// git does not preserve windows hidden flag so we have to set it.
	filenameW, err := syscall.UTF16PtrFromString(filename)
	if err != nil {
		return err
	}
	return syscall.SetFileAttributes(filenameW, syscall.FILE_ATTRIBUTE_HIDDEN)
}

func fileOpenOSOverride(*FileDialog) bool {
	return false
}

func fileSaveOSOverride(*FileDialog) bool {
	return false
}

func getFavoriteLocation(homeURI fyne.URI, name string) (fyne.URI, error) {
	return storage.Child(homeURI, name)
}
