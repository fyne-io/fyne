package commands

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"

	"golang.org/x/sys/execabs"
)

type unixData struct {
	Name, Exec, Icon string
	Local            string
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

	tplData := unixData{Name: p.Name, Exec: filepath.Base(p.exe), Icon: p.Name + filepath.Ext(p.icon), Local: local}
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
		tarCmd := execabs.Command("tar", "-Jcf", p.Name+".tar.xz", "-C", tempDir, "usr", "Makefile")
		tarCmd.Stderr = &buf
		if err = tarCmd.Run(); err != nil {
			return fmt.Errorf("failed to create archive with tar: %s - %w", buf.String(), err)
		}
	}

	return nil
}
