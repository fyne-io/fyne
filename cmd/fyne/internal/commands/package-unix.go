package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/cmd/fyne/internal/metadata"
	"fyne.io/fyne/v2/cmd/fyne/internal/templates"

	"golang.org/x/sys/execabs"
)

type unixData struct {
	Name, Exec, Icon string
	Local            string
	GenericName      string
	Categories       string
	Comment          string
	Keywords         string
}

func (p *Packager) packageUNIX() error {
	var prefixDir string
	local := "local/"
	tempDir := "tmp-pkg"

	if p.install {
		tempDir = ""
	}

	if _, err := os.Stat(filepath.Join("/", "usr", "local")); os.IsNotExist(err) {
		prefixDir = util.EnsureSubDir(util.EnsureSubDir(p.dir, tempDir), "usr")
		local = ""
	} else {
		prefixDir = util.EnsureSubDir(util.EnsureSubDir(util.EnsureSubDir(p.dir, tempDir), "usr"), "local")
	}

	shareDir := util.EnsureSubDir(prefixDir, "share")

	binDir := util.EnsureSubDir(prefixDir, "bin")
	binName := filepath.Join(binDir, filepath.Base(p.exe))
	err := util.CopyExeFile(p.exe, binName)
	if err != nil {
		return fmt.Errorf("failed to copy application binary file: %w", err)
	}

	iconDir := util.EnsureSubDir(shareDir, "pixmaps")
	iconPath := filepath.Join(iconDir, p.Name+filepath.Ext(p.icon))
	err = util.CopyFile(p.icon, iconPath)
	if err != nil {
		return fmt.Errorf("failed to copy icon: %w", err)
	}

	appsDir := util.EnsureSubDir(shareDir, "applications")
	desktop := filepath.Join(appsDir, p.Name+".desktop")
	deskFile, _ := os.Create(desktop)

	data, err := metadata.LoadStandard(p.srcDir)
	if err != nil {
		return fmt.Errorf("failed to load data from FyneApp.toml: %w", err)
	}

	fmt.Println(data.LinuxAndBSD)

	tplData := unixData{
		Name:        p.Name,
		Exec:        filepath.Base(p.exe),
		Icon:        p.Name + filepath.Ext(p.icon),
		Local:       local,
		GenericName: data.LinuxAndBSD.GenericName,
		Keywords:    data.LinuxAndBSD.Keywords,
		Comment:     data.LinuxAndBSD.Comment,
		Categories:  data.LinuxAndBSD.Categories,
	}
	err = templates.DesktopFileUNIX.Execute(deskFile, tplData)
	if err != nil {
		return fmt.Errorf("failed to write desktop entry string: %w", err)
	}

	if !p.install {
		defer os.RemoveAll(filepath.Join(p.dir, tempDir))

		makefile, _ := os.Create(filepath.Join(p.dir, tempDir, "Makefile"))
		err := templates.MakefileUNIX.Execute(makefile, tplData)
		if err != nil {
			return fmt.Errorf("failed to write Makefile string: %w", err)
		}

		var buf bytes.Buffer
		tarCmd := execabs.Command("tar", "-Jcf", filepath.Join(p.dir, p.Name+".tar.xz"), "-C", filepath.Join(p.dir, tempDir), "usr", "Makefile")
		tarCmd.Stderr = &buf
		if err = tarCmd.Run(); err != nil {
			return fmt.Errorf("failed to create archive with tar: %s - %w", buf.String(), err)
		}
	}

	return nil
}
