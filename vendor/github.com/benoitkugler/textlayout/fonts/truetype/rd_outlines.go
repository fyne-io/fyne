package truetype

import (
	"errors"
	"fmt"

	"github.com/benoitkugler/textlayout/fonts"
)

// this file converts from font format for glyph outlines to
// segments that rasterizer will consume
//
// adapted from snft/truetype.go

func midPoint(p, q fonts.SegmentPoint) fonts.SegmentPoint {
	return fonts.SegmentPoint{
		X: (p.X + q.X) / 2,
		Y: (p.Y + q.Y) / 2,
	}
}

// build the segments from the resolved contour points
func buildSegments(points []contourPoint) []fonts.Segment {
	var (
		firstOnCurveValid, firstOffCurveValid, lastOffCurveValid bool
		firstOnCurve, firstOffCurve, lastOffCurve                fonts.SegmentPoint
		out                                                      []fonts.Segment
	)

	for _, point := range points {
		p := point.SegmentPoint
		if !firstOnCurveValid {
			if point.isOnCurve {
				firstOnCurve = p
				firstOnCurveValid = true
				out = append(out, fonts.Segment{
					Op:   fonts.SegmentOpMoveTo,
					Args: [3]fonts.SegmentPoint{p},
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
				out = append(out, fonts.Segment{
					Op:   fonts.SegmentOpMoveTo,
					Args: [3]fonts.SegmentPoint{firstOnCurve},
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
				out = append(out, fonts.Segment{
					Op:   fonts.SegmentOpLineTo,
					Args: [3]fonts.SegmentPoint{p},
				})
			}
		} else {
			if !point.isOnCurve {
				out = append(out, fonts.Segment{
					Op: fonts.SegmentOpQuadTo,
					Args: [3]fonts.SegmentPoint{
						lastOffCurve,
						midPoint(lastOffCurve, p),
					},
				})
				lastOffCurve = p
				lastOffCurveValid = true
			} else {
				out = append(out, fonts.Segment{
					Op:   fonts.SegmentOpQuadTo,
					Args: [3]fonts.SegmentPoint{lastOffCurve, p},
				})
				lastOffCurveValid = false
			}
		}

		if point.isEndPoint {
			// closing the contour
			switch {
			case !firstOffCurveValid && !lastOffCurveValid:
				out = append(out, fonts.Segment{
					Op:   fonts.SegmentOpLineTo,
					Args: [3]fonts.SegmentPoint{firstOnCurve},
				})
			case !firstOffCurveValid && lastOffCurveValid:
				out = append(out, fonts.Segment{
					Op:   fonts.SegmentOpQuadTo,
					Args: [3]fonts.SegmentPoint{lastOffCurve, firstOnCurve},
				})
			case firstOffCurveValid && !lastOffCurveValid:
				out = append(out, fonts.Segment{
					Op:   fonts.SegmentOpQuadTo,
					Args: [3]fonts.SegmentPoint{firstOffCurve, firstOnCurve},
				})
			case firstOffCurveValid && lastOffCurveValid:
				out = append(out, fonts.Segment{
					Op: fonts.SegmentOpQuadTo,
					Args: [3]fonts.SegmentPoint{
						lastOffCurve,
						midPoint(lastOffCurve, firstOffCurve),
					},
				},
					fonts.Segment{
						Op:   fonts.SegmentOpQuadTo,
						Args: [3]fonts.SegmentPoint{firstOffCurve, firstOnCurve},
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
func (f *Font) glyphDataFromGlyf(glyph GID) (fonts.GlyphOutline, error) {
	if int(glyph) >= len(f.Glyf) {
		return fonts.GlyphOutline{}, fmt.Errorf("out of range glyph %d", glyph)
	}
	var points []contourPoint
	f.getPointsForGlyph(glyph, 0, &points)
	segments := buildSegments(points[:len(points)-phantomCount])
	return fonts.GlyphOutline{Segments: segments}, nil
}

var errNoCFFTable error = errors.New("no CFF table")

func (f *Font) glyphDataFromCFF1(glyph GID) (fonts.GlyphOutline, error) {
	if f.cff == nil {
		return fonts.GlyphOutline{}, errNoCFFTable
	}
	segments, _, err := f.cff.LoadGlyph(glyph)
	if err != nil {
		return fonts.GlyphOutline{}, err
	}
	return fonts.GlyphOutline{Segments: segments}, nil
}
