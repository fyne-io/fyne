// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package cff

import (
	"fmt"

	"github.com/go-text/typesetting/opentype/api"
	ps "github.com/go-text/typesetting/opentype/api/font/cff/interpreter"
	"github.com/go-text/typesetting/opentype/tables"
)

// LoadGlyph parses the glyph charstring to compute segments and path bounds.
// It returns an error if the glyph is invalid or if decoding the charstring fails.
func (f *Font) LoadGlyph(glyph tables.GlyphID) ([]api.Segment, ps.PathBounds, error) {
	var (
		psi    ps.Machine
		loader type2CharstringHandler
		index  byte = 0
		err    error
	)
	if f.fdSelect != nil {
		index, err = f.fdSelect.fontDictIndex(glyph)
		if err != nil {
			return nil, ps.PathBounds{}, err
		}
	}
	if int(glyph) >= len(f.Charstrings) {
		return nil, ps.PathBounds{}, fmt.Errorf("invalid glyph index %d", glyph)
	}

	subrs := f.localSubrs[index]
	err = psi.Run(f.Charstrings[glyph], subrs, f.globalSubrs, &loader)
	return loader.cs.Segments, loader.cs.Bounds, err
}

// type2CharstringHandler implements operators needed to fetch Type2 charstring metrics
type type2CharstringHandler struct {
	cs ps.CharstringReader

	// found in private DICT, needed since we can't differenciate
	// no width set from 0 width
	// `width` must be initialized to default width
	nominalWidthX int32
	width         int32
}

func (type2CharstringHandler) Context() ps.Context { return ps.Type2Charstring }

func (met *type2CharstringHandler) Apply(state *ps.Machine, op ps.Operator) error {
	var err error
	if !op.IsEscaped {
		switch op.Operator {
		case 11: // return
			return state.Return() // do not clear the arg stack
		case 14: // endchar
			if state.ArgStack.Top > 0 { // width is optional
				met.width = met.nominalWidthX + state.ArgStack.Vals[0]
			}
			met.cs.ClosePath()
			return ps.ErrInterrupt
		case 10: // callsubr
			return ps.LocalSubr(state) // do not clear the arg stack
		case 29: // callgsubr
			return ps.GlobalSubr(state) // do not clear the arg stack
		case 21: // rmoveto
			if state.ArgStack.Top > 2 { // width is optional
				met.width = met.nominalWidthX + state.ArgStack.Vals[0]
			}
			err = met.cs.Rmoveto(state)
		case 22: // hmoveto
			if state.ArgStack.Top > 1 { // width is optional
				met.width = met.nominalWidthX + state.ArgStack.Vals[0]
			}
			err = met.cs.Hmoveto(state)
		case 4: // vmoveto
			if state.ArgStack.Top > 1 { // width is optional
				met.width = met.nominalWidthX + state.ArgStack.Vals[0]
			}
			err = met.cs.Vmoveto(state)
		case 1, 18: // hstem, hstemhm
			met.cs.Hstem(state)
		case 3, 23: // vstem, vstemhm
			met.cs.Vstem(state)
		case 19, 20: // hintmask, cntrmask
			// variable number of arguments, but always even
			// for xxxmask, if there are arguments on the stack, then this is an impliied stem
			if state.ArgStack.Top&1 != 0 {
				met.width = met.nominalWidthX + state.ArgStack.Vals[0]
			}
			met.cs.Hintmask(state)
			// the stack is managed by the previous call
			return nil

		case 5: // rlineto
			met.cs.Rlineto(state)
		case 6: // hlineto
			met.cs.Hlineto(state)
		case 7: // vlineto
			met.cs.Vlineto(state)
		case 8: // rrcurveto
			met.cs.Rrcurveto(state)
		case 24: // rcurveline
			err = met.cs.Rcurveline(state)
		case 25: // rlinecurve
			err = met.cs.Rlinecurve(state)
		case 26: // vvcurveto
			met.cs.Vvcurveto(state)
		case 27: // hhcurveto
			met.cs.Hhcurveto(state)
		case 30: // vhcurveto
			met.cs.Vhcurveto(state)
		case 31: // hvcurveto
			met.cs.Hvcurveto(state)
		default:
			// no other operands are allowed before the ones handled above
			err = fmt.Errorf("invalid operator %s in charstring", op)
		}
	} else {
		switch op.Operator {
		case 34: // hflex
			err = met.cs.Hflex(state)
		case 35: // flex
			err = met.cs.Flex(state)
		case 36: // hflex1
			err = met.cs.Hflex1(state)
		case 37: // flex1
			err = met.cs.Flex1(state)
		default:
			// no other operands are allowed before the ones handled above
			err = fmt.Errorf("invalid operator %s in charstring", op)
		}
	}
	state.ArgStack.Clear()
	return err
}
