// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

package tables

// Trak is the tracking table.
// See - https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6trak.html
type Trak struct {
	version  uint32    // Version number of the tracking table (0x00010000 for the current version).
	format   uint16    // Format of the tracking table (set to 0).
	Horiz    TrackData `offsetSize:"Offset16"` // Offset from start of tracking table to TrackData for horizontal text (or 0 if none).
	Vert     TrackData `offsetSize:"Offset16"` // Offset from start of tracking table to TrackData for vertical text (or 0 if none).
	reserved uint16    // Reserved. Set to 0.
}

// IsEmpty return `true` it the table has no entries.
func (t Trak) IsEmpty() bool {
	return len(t.Horiz.TrackTable)+len(t.Vert.TrackTable) == 0
}

type TrackData struct {
	nTracks    uint16            // Number of separate tracks included in this table.
	nSizes     uint16            // Number of point sizes included in this table.
	SizeTable  []Float1616       `offsetSize:"Offset32" offsetRelativeTo:"Parent" arrayCount:"ComputedField-nSizes"` // Offset from start of the tracking table to the start of the size subtable.
	TrackTable []TrackTableEntry `arrayCount:"ComputedField-nTracks" arguments:"perSizeTrackingCount=.nSizes"`       // Array[nTracks] of TrackTableEntry records.
}

type TrackTableEntry struct {
	Track           Float1616 // Track value for this record.
	NameIndex       uint16    // The 'name' table index for this track (a short word or phrase like "loose" or "very tight"). NameIndex has a value greater than 255 and less than 32768.
	PerSizeTracking []int16   `offsetSize:"Offset16" offsetRelativeTo:"GrandParent"` // in font units, with length len(SizeTable)
}
