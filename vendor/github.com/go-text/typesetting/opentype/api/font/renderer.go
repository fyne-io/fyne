// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package font

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"

	"github.com/go-text/typesetting/opentype/api"
)

// GlyphData returns the glyph content for [gid], or nil if
// not found.
func (f *Face) GlyphData(gid GID) api.GlyphData {
	// since outline may be specified for SVG and bitmaps, check it at the end
	outB, err := f.sbix.glyphData(gID(gid), f.XPpem, f.YPpem)
	if err == nil {
		outline, ok := f.outlineGlyphData(gID(gid))
		if ok {
			outB.Outline = &outline
		}
		return outB
	}

	outB, err = f.bitmap.glyphData(gID(gid), f.XPpem, f.YPpem)
	if err == nil {
		outline, ok := f.outlineGlyphData(gID(gid))
		if ok {
			outB.Outline = &outline
		}
		return outB
	}

	outS, ok := f.svg.glyphData(gID(gid))
	if ok {
		// Spec :
		// For every SVG glyph description, there must be a corresponding TrueType,
		// CFF or CFF2 glyph description in the font.
		outS.Outline, _ = f.outlineGlyphData(gID(gid))
		return outS
	}

	if out, ok := f.outlineGlyphData(gID(gid)); ok {
		return out
	}

	return nil
}

func (sb sbix) glyphData(gid gID, xPpem, yPpem uint16) (api.GlyphBitmap, error) {
	st := sb.chooseStrike(xPpem, yPpem)
	if st == nil {
		return api.GlyphBitmap{}, errors.New("empty 'sbix' table")
	}

	glyph := strikeGlyph(st, gid, 0)
	if glyph.GraphicType == 0 {
		return api.GlyphBitmap{}, fmt.Errorf("no glyph %d in 'sbix' table for resolution (%d, %d)", gid, xPpem, yPpem)
	}

	out := api.GlyphBitmap{Data: glyph.Data}
	var err error
	out.Width, out.Height, out.Format, err = decodeBitmapConfig(glyph)

	return out, err
}

func (bt bitmap) glyphData(gid gID, xPpem, yPpem uint16) (api.GlyphBitmap, error) {
	st := bt.chooseStrike(xPpem, yPpem)
	if st == nil || st.ppemX == 0 || st.ppemY == 0 {
		return api.GlyphBitmap{}, errors.New("empty bitmap table")
	}

	subtable := st.findTable(gid)
	if subtable == nil {
		return api.GlyphBitmap{}, fmt.Errorf("no glyph %d in bitmap table for resolution (%d, %d)", gid, xPpem, yPpem)
	}

	glyph := subtable.image(gid)
	if glyph == nil {
		return api.GlyphBitmap{}, fmt.Errorf("no glyph %d in bitmap table for resolution (%d, %d)", gid, xPpem, yPpem)
	}

	out := api.GlyphBitmap{
		Data:   glyph.image,
		Width:  int(glyph.metrics.Width),
		Height: int(glyph.metrics.Height),
	}
	switch subtable.imageFormat {
	case 17, 18, 19: // PNG
		out.Format = api.PNG
	case 2, 5:
		out.Format = api.BlackAndWhite
	default:
		return api.GlyphBitmap{}, fmt.Errorf("unsupported format %d in bitmap table", subtable.imageFormat)
	}

	return out, nil
}

// look for data in 'glyf' and 'cff' tables
func (f *Face) outlineGlyphData(gid gID) (api.GlyphOutline, bool) {
	out, err := f.glyphDataFromCFF1(gid)
	if err == nil {
		return out, true
	}

	out, err = f.glyphDataFromGlyf(gid)
	if err == nil {
		return out, true
	}

	return api.GlyphOutline{}, false
}

func (s svg) glyphData(gid gID) (api.GlyphSVG, bool) {
	data, ok := s.rawGlyphData(gid)
	if !ok {
		return api.GlyphSVG{}, false
	}

	// un-compress if needed
	if r, err := gzip.NewReader(bytes.NewReader(data)); err == nil {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, r); err == nil {
			data = buf.Bytes()
		}
	}

	return api.GlyphSVG{Source: data}, true
}

// this file converts from font format for glyph outlines to
// segments that rasterizer will consume
//
// adapted from snft/truetype.go

func midPoint(p, q api.SegmentPoint) api.SegmentPoint {
	return api.SegmentPoint{
		X: (p.X + q.X) / 2,
		Y: (p.Y + q.Y) / 2,
	}
}

// build the segments from the resolved contour points
func buildSegments(points []contourPoint) []api.Segment {
	var (
		firstOnCurveValid, firstOffCurveValid, lastOffCurveValid bool
		firstOnCurve, firstOffCurve, lastOffCurve                api.SegmentPoint
		out                                                      []api.Segment
	)

	for _, point := range points {
		p := point.SegmentPoint
		if !firstOnCurveValid {
			if point.isOnCurve {
				firstOnCurve = p
				firstOnCurveValid = true
				out = append(out, api.Segment{
					Op:   api.SegmentOpMoveTo,
					Args: [3]api.SegmentPoint{p},
				})
			} else if !firstOffCurveValid {
				firstOffCurve = p
				firstOffCurveValid = true

				if !point.isEndPoint {
					continue
				}
			} else {
				firstOnCurve = midPoint(firstOffCurve, p)
				firstOnCurveValid = true
				lastOffCurve = p
				lastOffCurveValid = true
				out = append(out, api.Segment{
					Op:   api.SegmentOpMoveTo,
					Args: [3]api.SegmentPoint{firstOnCurve},
				})
			}
		} else if !lastOffCurveValid {
			if !point.isOnCurve {
				lastOffCurve = p
				lastOffCurveValid = true

				if !point.isEndPoint {
					continue
				}
			} else {
				out = append(out, api.Segment{
					Op:   api.SegmentOpLineTo,
					Args: [3]api.SegmentPoint{p},
				})
			}
		} else {
			if !point.isOnCurve {
				out = append(out, api.Segment{
					Op: api.SegmentOpQuadTo,
					Args: [3]api.SegmentPoint{
						lastOffCurve,
						midPoint(lastOffCurve, p),
					},
				})
				lastOffCurve = p
				lastOffCurveValid = true
			} else {
				out = append(out, api.Segment{
					Op:   api.SegmentOpQuadTo,
					Args: [3]api.SegmentPoint{lastOffCurve, p},
				})
				lastOffCurveValid = false
			}
		}

		if point.isEndPoint {
			// closing the contour
			switch {
			case !firstOffCurveValid && !lastOffCurveValid:
				out = append(out, api.Segment{
					Op:   api.SegmentOpLineTo,
					Args: [3]api.SegmentPoint{firstOnCurve},
				})
			case !firstOffCurveValid && lastOffCurveValid:
				out = append(out, api.Segment{
					Op:   api.SegmentOpQuadTo,
					Args: [3]api.SegmentPoint{lastOffCurve, firstOnCurve},
				})
			case firstOffCurveValid && !lastOffCurveValid:
				out = append(out, api.Segment{
					Op:   api.SegmentOpQuadTo,
					Args: [3]api.SegmentPoint{firstOffCurve, firstOnCurve},
				})
			case firstOffCurveValid && lastOffCurveValid:
				out = append(out, api.Segment{
					Op: api.SegmentOpQuadTo,
					Args: [3]api.SegmentPoint{
						lastOffCurve,
						midPoint(lastOffCurve, firstOffCurve),
					},
				},
					api.Segment{
						Op:   api.SegmentOpQuadTo,
						Args: [3]api.SegmentPoint{firstOffCurve, firstOnCurve},
					},
				)
			}

			firstOnCurveValid = false
			firstOffCurveValid = false
			lastOffCurveValid = false
		}
	}

	return out
}

// apply variation when needed
func (f *Face) glyphDataFromGlyf(glyph gID) (api.GlyphOutline, error) {
	if int(glyph) >= len(f.glyf) {
		return api.GlyphOutline{}, fmt.Errorf("out of range glyph %d", glyph)
	}
	var points []contourPoint
	f.getPointsForGlyph(glyph, 0, &points)
	segments := buildSegments(points[:len(points)-phantomCount])
	return api.GlyphOutline{Segments: segments}, nil
}

var errNoCFFTable error = errors.New("no CFF table")

func (f *Font) glyphDataFromCFF1(glyph gID) (api.GlyphOutline, error) {
	if f.cff == nil {
		return api.GlyphOutline{}, errNoCFFTable
	}
	segments, _, err := f.cff.LoadGlyph(glyph)
	if err != nil {
		return api.GlyphOutline{}, err
	}
	return api.GlyphOutline{Segments: segments}, nil
}

// BitmapSizes returns the size of bitmap glyphs present in the font.
func (font *Font) BitmapSizes() []api.BitmapSize {
	upem := font.head.UnitsPerEm

	avgWidth := font.os2.xAvgCharWidth

	// handle invalid head/os2 tables
	if upem == 0 || font.os2.version == 0xFFFF {
		avgWidth = 1
		upem = 1
	}

	// adapted from freetype tt_face_load_sbit
	if font.bitmap != nil {
		return font.bitmap.availableSizes(avgWidth, upem)
	}

	if hori := font.hhea; hori != nil {
		return font.sbix.availableSizes(hori, avgWidth, upem)
	}

	return nil
}
