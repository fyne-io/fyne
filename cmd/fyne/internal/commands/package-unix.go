package commands

import (
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/cmd/fyne/internal/templates"
	"fyne.io/fyne/v2/cmd/fyne/internal/util"
	"github.com/pkg/errors"

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
		return errors.Wrap(err, "Failed to copy application binary file")
	}

	iconDir := util.EnsureSubDir(shareDir, "pixmaps")
	iconPath := filepath.Join(iconDir, p.name+filepath.Ext(p.icon))
	err = util.CopyFile(p.icon, iconPath)
	if err != nil {
		return errors.Wrap(err, "Failed to copy icon")
	}

	appsDir := util.EnsureSubDir(shareDir, "applications")
	desktop := filepath.Join(appsDir, p.name+".desktop")
	deskFile, _ := os.Create(desktop)

	tplData := unixData{Name: p.name, Exec: filepath.Base(p.exe), Icon: p.name + filepath.Ext(p.icon), Local: local}
	err = templates.DesktopFileUNIX.Execute(deskFile, tplData)
	if err != nil {
		return errors.Wrap(err, "Failed to write desktop entry string")
	}

	if !p.install {
		defer os.RemoveAll(filepath.Join(p.dir, tempDir))

		makefile, _ := os.Create(filepath.Join(p.dir, tempDir, "Makefile"))
		err := templates.MakefileUNIX.Execute(makefile, tplData)
		if err != nil {
			return errors.Wrap(err, "Failed to write Makefile string")
		}

		tarCmd := execabs.Command("tar", "-Jcf", p.name+".tar.xz", "-C", tempDir, "usr", "Makefile")
		if err = tarCmd.Run(); err != nil {
			return errors.Wrap(err, "Failed to create archive with tar")
		}
	}

	return nil
}
