package truetype

import (
	"encoding/binary"
	"errors"
)

type tableVorg struct {
	metrics []struct {
		glyph  GID
		origin int16
	}
	defaultOrigin int16
}

func (t tableVorg) getYOrigin(glyph GID) int16 {
	// binary search
	for i, j := 0, len(t.metrics); i < j; {
		h := i + (j-i)/2
		entry := t.metrics[h]
		if glyph < entry.glyph {
			j = h
		} else if entry.glyph < glyph {
			i = h + 1
		} else {
			return entry.origin
		}
	}
	return t.defaultOrigin
}

func parseTableVorg(data []byte) (out tableVorg, err error) {
	if len(data) < 8 {
		return out, errors.New("invalid 'vorg' table (EOF)")
	}

	out.defaultOrigin = int16(binary.BigEndian.Uint16(data[4:]))
	count := int(binary.BigEndian.Uint16(data[6:]))
	if len(data) < 8+4*count {
		return out, errors.New("invalid 'vorg' table (EOF)")
	}
	out.metrics = make([]struct {
		glyph  GID
		origin int16
	}, count)
	for i := range out.metrics {
		out.metrics[i].glyph = GID(binary.BigEndian.Uint16(data[8+4*i:]))
		out.metrics[i].origin = int16(binary.BigEndian.Uint16(data[8+4*i+2:]))
	}
	return out, nil
}
