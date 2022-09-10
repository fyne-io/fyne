package truetype

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type TableOS2Version0 struct {
	Version             uint16
	XAvgCharWidth       uint16
	USWeightClass       uint16
	USWidthClass        uint16
	FSType              uint16
	YSubscriptXSize     int16
	YSubscriptYSize     int16
	YSubscriptXOffset   int16
	YSubscriptYOffset   int16
	YSuperscriptXSize   int16
	YSuperscriptYSize   int16
	YSuperscriptXOffset int16
	YSuperscriptYOffset int16
	YStrikeoutSize      int16
	YStrikeoutPosition  int16
	SFamilyClass        int16
	Panose              [10]byte
	UlCharRange         [4]uint32
	AchVendID           Tag
	FsSelection         uint16
	USFirstCharIndex    uint16
	USLastCharIndex     uint16
	STypoAscender       int16
	STypoDescender      int16
	STypoLineGap        int16
	UsWinAscent         uint16
	UsWinDescent        uint16
}

type TableOS2Version1 struct {
	TableOS2Version0

	UlCodePageRange1 uint32
	UlCodePageRange2 uint32
}

// TableOS2Version4 is the OS2 format for versions 2,3 and 4
type TableOS2Version4 struct {
	TableOS2Version1

	SxHeigh       int16
	SCapHeight    int16
	UsDefaultChar uint16
	UsBreakChar   uint16
	UsMaxContext  uint16
}

type TableOS2 struct {
	TableOS2Version4

	UsLowerPointSize uint16
	UsUpperPointSize uint16
}

func parseTableOS2(buf []byte) (*TableOS2, error) {
	if len(buf) < 2 {
		return nil, errors.New("invalid 'os2' table (EOF)")
	}

	var (
		dst interface{}
		out TableOS2
	)
	version := binary.BigEndian.Uint16(buf)
	switch version {
	case 0:
		dst = &out.TableOS2Version0
	case 1, 2, 3:
		dst = &out.TableOS2Version1
	case 4:
		dst = &out.TableOS2Version4
	case 5:
		dst = &out
	default:
		return nil, fmt.Errorf("unsupported 'os2' table version: %d", version)
	}

	if err := binary.Read(bytes.NewReader(buf), binary.BigEndian, dst); err != nil {
		return nil, fmt.Errorf("invalid 'os2' table: %s", err)
	}

	return &out, nil
}

func (t *TableOS2) useTypoMetrics() bool {
	const useTypoMetrics = 1 << 7
	return t.FsSelection&useTypoMetrics != 0
}

func (t *TableOS2) hasData() bool {
	return t.USWeightClass != 0 || t.USWidthClass != 0 || t.USFirstCharIndex != 0 || t.USLastCharIndex != 0
}
