package icns

import (
	"bytes"
	"image"
	"image/png"
	"io"
)

// Icon encodes an icns icon.
type Icon struct {
	Type  OsType
	Image image.Image

	header    [8]byte
	headerSet bool
	data      []byte
}

// WriteTo encodes the icon into wr.
func (i *Icon) WriteTo(wr io.Writer) (int64, error) {
	var written int64
	if err := i.encodeImage(); err != nil {
		return written, err
	}
	size, err := i.writeHeader(wr)
	written += size
	if err != nil {
		return written, err
	}
	size, err = i.writeData(wr)
	written += size
	if err != nil {
		return written, err
	}
	return written, nil
}

func (i *Icon) encodeImage() error {
	if len(i.data) > 0 {
		return nil
	}
	data, err := encodeImage(i.Image)
	if err != nil {
		return err
	}
	i.data = data
	return nil
}

func encodeImage(img image.Image) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (i *Icon) writeHeader(wr io.Writer) (int64, error) {
	if !i.headerSet {
		defer func() { i.headerSet = true }()
		i.header[0] = i.Type.ID[0]
		i.header[1] = i.Type.ID[1]
		i.header[2] = i.Type.ID[2]
		i.header[3] = i.Type.ID[3]
		length := uint32(len(i.data) + 8)
		writeUint32(i.header[4:8], length)
	}
	written, err := wr.Write(i.header[:8])
	return int64(written), err
}

func (i *Icon) writeData(wr io.Writer) (int64, error) {
	written, err := wr.Write(i.data)
	return int64(written), err
}

// IconSet encodes a set of icons into an ICNS file.
type IconSet struct {
	Icons []*Icon

	header    [8]byte
	headerSet bool
	data      []byte
}

// WriteTo writes the ICNS file to wr.
func (s *IconSet) WriteTo(wr io.Writer) (int64, error) {
	var written int64
	if err := s.encodeIcons(); err != nil {
		return written, err
	}
	size, err := s.writeHeader(wr)
	written += size
	if err != nil {
		return written, err
	}
	size, err = s.writeData(wr)
	written += size
	if err != nil {
		return written, err
	}
	return written, nil
}

func (s *IconSet) encodeIcons() error {
	if len(s.data) > 0 {
		return nil
	}
	buf := bytes.NewBuffer(nil)
	for _, icon := range s.Icons {
		if icon == nil {
			continue
		}
		if _, err := icon.WriteTo(buf); err != nil {
			return err
		}
	}
	s.data = buf.Bytes()
	return nil
}

func (s *IconSet) writeHeader(wr io.Writer) (int64, error) {
	if !s.headerSet {
		defer func() { s.headerSet = true }()
		s.header[0] = 'i'
		s.header[1] = 'c'
		s.header[2] = 'n'
		s.header[3] = 's'
		length := uint32(len(s.data) + 8)
		writeUint32(s.header[4:8], length)
	}
	written, err := wr.Write(s.header[:8])
	return int64(written), err
}

func (s *IconSet) writeData(wr io.Writer) (int64, error) {
	written, err := wr.Write(s.data)
	return int64(written), err
}
