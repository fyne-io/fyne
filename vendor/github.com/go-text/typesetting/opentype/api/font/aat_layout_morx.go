// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package font

import "github.com/go-text/typesetting/opentype/tables"

type Morx []MorxChain

func newMorx(table tables.Morx) Morx {
	if len(table.Chains) == 0 {
		return nil
	}
	out := make(Morx, len(table.Chains))
	for i, c := range table.Chains {
		out[i] = newMorxChain(c)
	}
	return out
}

type MorxChain struct {
	Features     []tables.AATFeature
	Subtables    []MorxSubtable
	DefaultFlags uint32
}

func newMorxChain(table tables.MorxChain) (out MorxChain) {
	out.DefaultFlags = table.Flags
	out.Features = table.Features
	out.Subtables = make([]MorxSubtable, len(table.Subtables))
	for i, s := range table.Subtables {
		out.Subtables[i] = newMorxSubtable(s)
	}
	return out
}

type MorxSubtable struct {
	Data     interface{ isMorxSubtable() }
	Coverage uint8  // high byte of the coverage flag
	Flags    uint32 // Mask identifying which subtable this is.
}

func (MorxRearrangementSubtable) isMorxSubtable() {}
func (MorxContextualSubtable) isMorxSubtable()    {}
func (MorxLigatureSubtable) isMorxSubtable()      {}
func (MorxNonContextualSubtable) isMorxSubtable() {}
func (MorxInsertionSubtable) isMorxSubtable()     {}

func newMorxSubtable(table tables.MorxChainSubtable) (out MorxSubtable) {
	out.Coverage = table.Coverage
	out.Flags = table.SubFeatureFlags
	switch data := table.Data.(type) {
	case tables.MorxSubtableRearrangement:
		out.Data = MorxRearrangementSubtable(newAATStableTable(data.AATStateTableExt))
	case tables.MorxSubtableContextual:
		out.Data = MorxContextualSubtable{
			Machine:       newAATStableTable(data.AATStateTableExt),
			Substitutions: data.Substitutions.Substitutions,
		}
	case tables.MorxSubtableLigature:
		s := MorxLigatureSubtable{
			Machine:        newAATStableTable(data.AATStateTableExt),
			LigatureAction: data.LigActions,
			Components:     data.Components,
			Ligatures:      make([]GID, len(data.Ligatures)),
		}
		for i, g := range data.Ligatures {
			s.Ligatures[i] = GID(g)
		}
		out.Data = s
	case tables.MorxSubtableNonContextual:
		out.Data = MorxNonContextualSubtable{Class: data.Class}
	case tables.MorxSubtableInsertion:
		s := MorxInsertionSubtable{
			Machine:    newAATStableTable(data.AATStateTableExt),
			Insertions: make([]GID, len(data.Insertions)),
		}
		for i, g := range data.Insertions {
			s.Insertions[i] = GID(g)
		}
		out.Data = s
	}
	return out
}

type MorxRearrangementSubtable AATStateTable

type MorxContextualSubtable struct {
	Substitutions []tables.AATLookup
	Machine       AATStateTable
}

type MorxLigatureSubtable struct {
	LigatureAction []uint32
	Components     []uint16
	Ligatures      []GID
	Machine        AATStateTable
}

type MorxNonContextualSubtable struct {
	Class tables.AATLookup // the lookup value is interpreted as a GlyphIndex
}

type MorxInsertionSubtable struct {
	// After successul parsing, this array may be safely
	// indexed by the indexes and counts from Machine entries.
	Insertions []GID
	Machine    AATStateTable
}
