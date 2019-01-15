package rsrc

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/akavel/rsrc/binutil"
	"github.com/akavel/rsrc/coff"
	"github.com/akavel/rsrc/ico"
	"github.com/akavel/rsrc/internal"
)

// on storing icons, see: http://blogs.msdn.com/b/oldnewthing/archive/2012/07/20/10331787.aspx
type _GRPICONDIR struct {
	ico.ICONDIR
	Entries []_GRPICONDIRENTRY
}

func (group _GRPICONDIR) Size() int64 {
	return int64(binary.Size(group.ICONDIR) + len(group.Entries)*binary.Size(group.Entries[0]))
}

type _GRPICONDIRENTRY struct {
	ico.IconDirEntryCommon
	Id uint16
}

func Embed(fnameout, arch, fnamein, fnameico string) error {
	lastid := uint16(0)
	newid := func() uint16 {
		lastid++
		return lastid
	}

	out := coff.NewRSRC()
	err := out.Arch(arch)
	if err != nil {
		return err
	}

	if fnamein != "" {
		manifest, err := binutil.SizedOpen(fnamein)
		if err != nil {
			return fmt.Errorf("rsrc: error opening manifest file '%s': %s", fnamein, err)
		}
		defer manifest.Close()

		id := newid()
		out.AddResource(coff.RT_MANIFEST, id, manifest)
		// TODO(akavel): reintroduce the Printlns in package main after Embed returns
		// fmt.Println("Manifest ID: ", id)
	}
	if fnameico != "" {
		for _, fnameicosingle := range strings.Split(fnameico, ",") {
			f, err := addIcon(out, fnameicosingle, newid)
			if err != nil {
				return err
			}
			defer f.Close()
		}
	}

	out.Freeze()

	return internal.Write(out, fnameout)
}

func addIcon(out *coff.Coff, fname string, newid func() uint16) (io.Closer, error) {
	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	icons, err := ico.DecodeHeaders(f)
	if err != nil {
		f.Close()
		return nil, err
	}

	if len(icons) > 0 {
		// RT_ICONs
		group := _GRPICONDIR{ICONDIR: ico.ICONDIR{
			Reserved: 0, // magic num.
			Type:     1, // magic num.
			Count:    uint16(len(icons)),
		}}
		for _, icon := range icons {
			id := newid()
			r := io.NewSectionReader(f, int64(icon.ImageOffset), int64(icon.BytesInRes))
			out.AddResource(coff.RT_ICON, id, r)
			group.Entries = append(group.Entries, _GRPICONDIRENTRY{icon.IconDirEntryCommon, id})
		}
		id := newid()
		out.AddResource(coff.RT_GROUP_ICON, id, group)
		// TODO(akavel): reintroduce the Printlns in package main after Embed returns
		// fmt.Println("Icon ", fname, " ID: ", id)
	}

	return f, nil
}
