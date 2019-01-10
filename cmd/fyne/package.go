package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/jackmordaunt/icns"
)

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}

func ensureSubDir(parent, name string) string {
	path := filepath.Join(parent, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	return path
}

func copyFile(src, tgt string) error {
	return copyFileMode(src, tgt, 0644)
}

func copyExeFile(src, tgt string) error {
	return copyFileMode(src, tgt, 0755)
}

func copyFileMode(src, tgt string, perm os.FileMode) error {
	_, err := os.Stat(src)
	if err != nil {
		return err
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.OpenFile(tgt, os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	return err
}

type packager struct {
	os, name, dir, exe, icon string
}

func (p *packager) packageLinux() {
	prefixDir := ensureSubDir(ensureSubDir(p.dir, "usr"), "local")
	shareDir := ensureSubDir(prefixDir, "share")

	appsDir := ensureSubDir(shareDir, "applications")
	desktop := filepath.Join(appsDir, p.name+".desktop")
	deskFile, _ := os.Create(desktop)
	io.WriteString(deskFile, "[Desktop Entry]\n"+
		"Type=Application\n"+
		"Name="+p.name+"\n"+
		"Exec="+filepath.Base(p.exe)+"\n"+
		"Icon="+filepath.Base(p.icon)+"\n")

	binDir := ensureSubDir(prefixDir, "bin")
	binName := filepath.Join(binDir, filepath.Base(p.exe))
	copyExeFile(p.exe, binName)

	iconThemeDir := ensureSubDir(ensureSubDir(shareDir, "icons"), "hicolor")
	iconDir := ensureSubDir(ensureSubDir(iconThemeDir, "512x512"), "apps")
	iconPath := filepath.Join(iconDir, filepath.Base(p.icon))
	copyFile(p.icon, iconPath)

	tarCmd := exec.Command("tar", "cf", p.name+".tar.gz", "usr")
	tarCmd.Run()
}

func (p *packager) packageDarwin() {
	appDir := ensureSubDir(p.dir, p.name+".app")

	contentsDir := ensureSubDir(appDir, "Contents")
	info := filepath.Join(contentsDir, "Info.plist")
	infoFile, _ := os.Create(info)
	io.WriteString(infoFile, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"+
		"<!DOCTYPE plist PUBLIC \"-//Apple Computer//DTD PLIST 1.0//EN\" \"http://www.apple.com/DTDs/PropertyList-1.0.dtd\">"+
		"<plist version=\"1.0\">"+
		"<dict>"+
		"<key>CFBundleGetInfoString</key>"+
		"<string>"+p.name+"</string>"+
		"<key>CFBundleExecutable</key>"+
		"<string>"+p.name+"</string>"+
		"<key>CFBundleIdentifier</key>"+
		"<string>com.example."+p.name+"</string>"+
		"<key>CFBundleName</key>"+
		"<string>"+p.name+"</string>"+
		"<key>CFBundleIconFile</key>"+
		"<string>icon.icns</string>"+
		"<key>CFBundleShortVersionString</key>"+
		"<string>1.0</string>"+
		"<key>CFBundleInfoDictionaryVersion</key>"+
		"<string>6.0</string>"+
		"<key>CFBundlePackageType</key>"+
		"<string>APPL</string>"+
		"<key>IFMajorVersion</key>"+
		"<integer>0</integer>"+
		"<key>IFMinorVersion</key>"+
		"<integer>1</integer>"+
		"<key>NSHighResolutionCapable</key><true/>"+
		"<key>NSSupportsAutomaticGraphicsSwitching</key><true/>"+
		"</dict>"+
		"</plist>")

	macOSDir := ensureSubDir(contentsDir, "MacOS")
	binName := filepath.Join(macOSDir, p.name)
	copyFile(p.exe, binName)

	resDir := ensureSubDir(contentsDir, "Resources")
	icnsPath := filepath.Join(resDir, "icon.icns")

	img, err := os.Open(p.icon)
	if err != nil {
		log.Fatalf("opening source image: %v", err)
	}
	defer img.Close()
	srcImg, _, err := image.Decode(img)
	if err != nil {
		log.Fatalf("decoding source image: %v", err)
	}
	dest, err := os.Create(icnsPath)
	if err != nil {
		log.Fatalf("opening destination file: %v", err)
	}
	defer dest.Close()
	if err := icns.Encode(dest, srcImg); err != nil {
		log.Fatalf("encoding icns: %v", err)
	}
}

func (p *packager) packageWindows() {

}

func (p *packager) addFlags() {
	flag.StringVar(&p.os, "os", runtime.GOOS, "The operating system to target (like GOOS)")
	flag.StringVar(&p.exe, "executable", "", "The path to the executable, default is the current dir main binary")
	flag.StringVar(&p.name, "name", "", "The name of the application, default is the executable file name")
	flag.StringVar(&p.icon, "icon", "", "The name of the application icon file")
}

func (*packager) printHelp(indent string) {
	fmt.Println(indent, "The package command prepares an application for distribution.")
	fmt.Println(indent, "Before running this command you must to build the application for release.")
	fmt.Println(indent, "Command usage: fyne package [parameters]")
}

func (p *packager) run(_ []string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to get the current directory, needed to find main executable")
	}
	p.dir = dir

	if p.exe == "" {
		exeName := filepath.Base(dir)
		if p.os == "windows" {
			exeName = exeName + ".exe"
		}

		p.exe = filepath.Join(dir, exeName)
	}

	if !exists(p.exe) {
		fmt.Fprintln(os.Stderr, "Executable", p.exe, "could not be found")
		fmt.Fprintln(os.Stderr, "You may not have run \"go build\" for the target os")
		os.Exit(1)
	}
	if p.name == "" {
		name := filepath.Base(p.exe)
		p.name = strings.TrimSuffix(name, filepath.Ext(name))
	}
	if p.icon == "" {
		fmt.Fprintln(os.Stderr, "The -icon parameter is required for packaging")
		os.Exit(1)
	}

	switch p.os {
	case "linux":
		p.packageLinux()
	case "darwin":
		p.packageDarwin()
	case "windows":
		p.packageWindows()
	default:
		fmt.Fprintln(os.Stderr, "Unsupported os parameter,", p.os)
	}
}
