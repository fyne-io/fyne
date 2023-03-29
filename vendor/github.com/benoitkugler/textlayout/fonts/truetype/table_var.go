package truetype

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func fixed1616ToFloat(fi uint32) float32 {
	// value are actually signed integers
	return float32(int32(fi)) / (1 << 16)
}

func fixed214ToFloat(fi uint16) float32 {
	// value are actually signed integers
	return float32(int16(fi)) / (1 << 14)
}

func parseTableFvar(table []byte, names TableName) (out TableFvar, err error) {
	hd, err := parseFvarHeader(table)
	if err != nil {
		return out, fmt.Errorf("invalid 'fvar' table header: %s", err)
	}

	axis, instanceOffset, err := parseVarAxisList(table, int(hd.axesArrayOffset), int(hd.axisSize), hd.axisCount)
	if err != nil {
		return out, err
	}
	// the instance offset is at the end of the axis
	instances, err := parseVarInstance(table, instanceOffset, int(hd.instanceSize), hd.instanceCount, hd.axisCount)
	if err != nil {
		return out, err
	}

	out = TableFvar{Axis: axis, Instances: instances}
	out.checkDefaultInstance(names)
	return out, nil
}

func (out *fvarHeader) mustParse(data []byte) {
	_ = data[15] // early bound checking
	out.majorVersion = uint16(binary.BigEndian.Uint16(data[0:2]))
	out.minorVersion = uint16(binary.BigEndian.Uint16(data[2:4]))
	out.axesArrayOffset = uint16(binary.BigEndian.Uint16(data[4:6]))
	out.reserved = uint16(binary.BigEndian.Uint16(data[6:8]))
	out.axisCount = uint16(binary.BigEndian.Uint16(data[8:10]))
	out.axisSize = uint16(binary.BigEndian.Uint16(data[10:12]))
	out.instanceCount = uint16(binary.BigEndian.Uint16(data[12:14]))
	out.instanceSize = uint16(binary.BigEndian.Uint16(data[14:16]))
}

func parseFvarHeader(data []byte) (fvarHeader, error) {
	var out fvarHeader
	if L := len(data); L < 16 {
		return fvarHeader{}, fmt.Errorf("EOF: expected length: 16, got %d", L)
	}
	out.mustParse(data)
	return out, nil
}

func (out *VarAxis) mustParse(data []byte) {
	_ = data[19] // early bound checking
	out.Tag = Tag(binary.BigEndian.Uint32(data[0:4]))
	out.Minimum = Float1616FromUint(binary.BigEndian.Uint32(data[4:8]))
	out.Default = Float1616FromUint(binary.BigEndian.Uint32(data[8:12]))
	out.Maximum = Float1616FromUint(binary.BigEndian.Uint32(data[12:16]))
	out.flags = uint16(binary.BigEndian.Uint16(data[16:18]))
	out.strid = NameID(binary.BigEndian.Uint16(data[18:20]))
}

func parseVarAxisList(table []byte, offset, size int, count uint16) ([]VarAxis, int, error) {
	// we need at least 20 byte per axis ....
	if size < 20 {
		return nil, 0, errors.New("invalid 'fvar' table axis")
	}
	// ...but "implementations must use the axisSize and instanceSize fields
	// to determine the start of each record".
	end := offset + int(count)*size
	if len(table) < end {
		return nil, 0, errors.New("invalid 'fvar' table axis")
	}

	out := make([]VarAxis, count) // guarded by previous check
	for i := range out {
		out[i].mustParse(table[offset+i*size:])
	}

	return out, end, nil
}

func parseVarInstance(table []byte, offset, size int, count, axisCount uint16) ([]VarInstance, error) {
	// we need at least 4+4*axisCount byte per instance ....
	if size < 4+4*int(axisCount) {
		return nil, errors.New("invalid 'fvar' table instance")
	}
	withPs := size >= 4+4*int(axisCount)+2

	// ...but "implementations must use the axisSize and instanceSize fields
	// to determine the start of each record".
	if len(table) < offset+int(count)*size {
		return nil, errors.New("invalid 'fvar' table axis")
	}

	out := make([]VarInstance, count) // limited by 16 bit type
	for i := range out {
		out[i] = parseOneVarInstance(table[offset+i*size:], axisCount, withPs)
	}

	return out, nil
}

// do not check the size of data
func parseOneVarInstance(data []byte, axisCount uint16, withPs bool) VarInstance {
	var out VarInstance
	out.Subfamily = NameID(binary.BigEndian.Uint16(data))
	// _ = binary.BigEndian.Uint16(data[2:]) reserved flags
	out.Coords = make([]float32, axisCount)
	for i := range out.Coords {
		out.Coords[i] = fixed1616ToFloat(binary.BigEndian.Uint32(data[4+i*4:]))
	}
	// optional PostscriptName id
	if withPs {
		out.PSStringID = NameID(binary.BigEndian.Uint16(data[4+axisCount*4:]))
	}

	return out
}

// -------------------------- avar table --------------------------

// one segment map for each axis, in the order of axes specified in the 'fvar' table.
type tableAvar [][]axisValueMap

type axisValueMap struct {
	from, to float32 // found as int16 2.14 fixed point
}

func parseTableAvar(data []byte, axisCountRef int) (tableAvar, error) {
	const avarHeaderSize = 2 * 4
	if len(data) < avarHeaderSize {
		return nil, errors.New("invalid 'avar' table (EOF)")
	}
	// table.majorVersion = binary.BigEndian.Uint16(data)
	// table.minorVersion = binary.BigEndian.Uint16(data[2:])
	// reserved
	axisCount := binary.BigEndian.Uint16(data[6:])
	out := make([][]axisValueMap, axisCount) // guarded by 16-bit constraint

	if int(axisCount) != axisCountRef {
		return nil, errors.New("invalid 'avar' table axis count")
	}

	var err error
	data = data[avarHeaderSize:] // start at the first segment list
	for i := range out {
		out[i], data, err = parseSegmentList(data)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

// data is at the start of the segment, return value at the start of the next
func parseSegmentList(data []byte) ([]axisValueMap, []byte, error) {
	const mapSize = 4
	if len(data) < 2 {
		return nil, nil, errors.New("invalid segment in 'avar' table")
	}
	count := binary.BigEndian.Uint16(data)
	size := int(count) * mapSize
	if len(data) < 2+size {
		return nil, nil, errors.New("invalid segment in 'avar' table")
	}
	out := make([]axisValueMap, count) // guarded by 16-bit constraint
	for i := range out {
		out[i].from = fixed214ToFloat(binary.BigEndian.Uint16(data[2+i*mapSize:]))
		out[i].to = fixed214ToFloat(binary.BigEndian.Uint16(data[2+i*mapSize+2:]))
	}
	data = data[2+size:]
	return out, data, nil
}

// VariationStoreIndex reference an item in the variation store
type VariationStoreIndex struct {
	DeltaSetOuter, DeltaSetInner uint16
}

// VariationStore store variation data.
// After successful parsing, every region indexes in `Datas` elements are valid,
// that is, can safely be used as index into `Regions`.
type VariationStore struct {
	Regions [][]VariationRegion // for each region, for each axis
	Datas   []ItemVariationData
}

// GetDelta uses the variation store and the selected instance coordinates
// to compute the value at `index`.
func (store VariationStore) GetDelta(index VariationStoreIndex, coords []float32) float32 {
	if int(index.DeltaSetOuter) >= len(store.Datas) {
		return 0
	}
	varData := store.Datas[index.DeltaSetOuter]
	if int(index.DeltaSetInner) >= len(varData.Deltas) {
		return 0
	}
	deltaSet := varData.Deltas[index.DeltaSetInner]
	var delta float32
	for i, regionIndex := range varData.RegionIndexes {
		region := store.Regions[regionIndex]
		v := float32(1)
		for axis, coord := range coords {
			factor := region[axis].evaluate(coord)
			v *= factor
		}
		delta += float32(deltaSet[i]) * v
	}
	return delta
}

func parseVariationStore(data []byte, offset uint32, axisCount int) (out VariationStore, err error) {
	if len(data) < int(offset)+8 {
		return out, errors.New("invalid item variation store (EOF)")
	}
	data = data[offset:]
	// format is ignored
	regionsOffset := binary.BigEndian.Uint32(data[2:])
	count := binary.BigEndian.Uint16(data[6:])

	out.Regions, err = parseItemVariationRegions(data, regionsOffset, axisCount)
	if err != nil {
		return out, err
	}

	if len(data) < 8+4*int(count) {
		return out, errors.New("invalid item variation store (EOF)")
	}
	out.Datas = make([]ItemVariationData, count)
	for i := range out.Datas {
		subtableOffset := binary.BigEndian.Uint32(data[8+4*i:])
		out.Datas[i], err = parseItemVariationData(data, subtableOffset, uint16(len(out.Regions)))
		if err != nil {
			return out, err
		}
	}
	return out, nil
}

func parseItemVariationRegions(data []byte, offset uint32, axisCountRef int) ([][]VariationRegion, error) {
	if len(data) < int(offset)+4 {
		return nil, errors.New("invalid item variation regions list (EOF)")
	}
	data = data[offset:]
	axisCount := int(binary.BigEndian.Uint16(data))
	regionCount := int(binary.BigEndian.Uint16(data[2:]))

	if axisCount != axisCountRef {
		return nil, errors.New("invalid item variation regions list number of axis")
	}

	if len(data) < 4+6*axisCount*regionCount {
		return nil, errors.New("invalid item variation regions list (EOF)")
	}
	regions := make([][]VariationRegion, regionCount)
	for i := range regions {
		ri := make([]VariationRegion, axisCount)
		for j := range ri {
			start := fixed214ToFloat(binary.BigEndian.Uint16(data[4+(i*axisCount+j)*6:]))
			peak := fixed214ToFloat(binary.BigEndian.Uint16(data[4+(i*axisCount+j)*6+2:]))
			end := fixed214ToFloat(binary.BigEndian.Uint16(data[4+(i*axisCount+j)*6+4:]))

			if start > peak || peak > end {
				return nil, errors.New("invalid item variation regions list")
			}
			if start < 0 && end > 0 && peak != 0 {
				return nil, errors.New("invalid item variation regions list")
			}
			ri[j] = VariationRegion{start, peak, end}
		}
		regions[i] = ri
	}
	return regions, nil
}

type ItemVariationData struct {
	RegionIndexes []uint16  // Array of indices into the variation region list for the regions referenced by this item variation data table.
	Deltas        [][]int16 // Each row as the same length as `RegionIndexes`
}

func parseItemVariationData(data []byte, offset uint32, nbRegions uint16) (out ItemVariationData, err error) {
	if len(data) < int(offset)+6 {
		return out, errors.New("invalid item variation data subtable (EOF)")
	}
	data = data[offset:]
	itemCount := int(binary.BigEndian.Uint16(data))
	shortDeltaCount := int(binary.BigEndian.Uint16(data[2:]))
	regionIndexCount := int(binary.BigEndian.Uint16(data[4:]))

	out.RegionIndexes, err = parseUint16s(data[6:], regionIndexCount)
	if err != nil {
		return out, fmt.Errorf("invalid item variation data subtable: %s", err)
	}
	// sanitize the indexes
	for _, regionIndex := range out.RegionIndexes {
		if regionIndex >= nbRegions {
			return out, fmt.Errorf("invalid item variation region index: %d (for size %d)", regionIndex, nbRegions)
		}
	}

	data = data[6+2*regionIndexCount:] // length checked by the previous `parseUint16s` call
	rowLength := shortDeltaCount + regionIndexCount
	if len(data) < itemCount*rowLength {
		return out, errors.New("invalid item variation data subtable (EOF)")
	}
	if shortDeltaCount > regionIndexCount {
		return out, errors.New("invalid item variation data subtable")
	}
	out.Deltas = make([][]int16, itemCount)
	for i := range out.Deltas {
		vi := make([]int16, regionIndexCount)
		j := 0
		for ; j < shortDeltaCount; j++ {
			vi[j] = int16(binary.BigEndian.Uint16(data[2*j:]))
		}
		for ; j < regionIndexCount; j++ {
			vi[j] = int16(int8(data[shortDeltaCount+j]))
		}
		out.Deltas[i] = vi
		data = data[rowLength:]
	}
	return out, nil
}

// ---------------------------------- mvar table ----------------------------------

type TableMvar struct {
	Values []VarValueRecord // sorted by tag
	Store  VariationStore
}

// return 0 if `tag` is not found
func (t TableMvar) getVar(tag Tag, coords []float32) float32 {
	// binary search
	for i, j := 0, len(t.Values); i < j; {
		h := i + (j-i)/2
		entry := t.Values[h]
		if tag < entry.Tag {
			j = h
		} else if entry.Tag < tag {
			i = h + 1
		} else {
			return t.Store.GetDelta(entry.Index, coords)
		}
	}
	return 0
}

type VarValueRecord struct {
	Tag   Tag
	Index VariationStoreIndex
}

func parseTableMvar(data []byte, axisCount int) (out TableMvar, err error) {
	if len(data) < 12 {
		return out, errors.New("invalid 'mvar' table (EOF)")
	}
	recordSize := int(binary.BigEndian.Uint16(data[6:]))
	recordCount := binary.BigEndian.Uint16(data[8:])
	storeOffset := uint32(binary.BigEndian.Uint16(data[10:]))

	if recordSize < 8 {
		return out, fmt.Errorf("invalid 'mvar' table record size: %d", recordSize)
	}

	out.Store, err = parseVariationStore(data, storeOffset, axisCount)
	if err != nil {
		return out, err
	}

	if len(data) < 12+recordSize*int(recordCount) {
		return out, errors.New("invalid 'mvar' table (EOF)")
	}
	out.Values = make([]VarValueRecord, recordCount)
	for i := range out.Values {
		out.Values[i].Tag = Tag(binary.BigEndian.Uint32(data[12+recordSize*i:]))
		out.Values[i].Index.DeltaSetOuter = binary.BigEndian.Uint16(data[12+recordSize*i+4:])
		out.Values[i].Index.DeltaSetInner = binary.BigEndian.Uint16(data[12+recordSize*i+6:])
	}

	return out, nil
}

// ---------------------------------- HVAR/VVAR ----------------------------------

type tableHVvar struct {
	store VariationStore
	// optional
	advances         deltaSetMapping
	leftSideBearings deltaSetMapping
}

func (t tableHVvar) getAdvanceVar(glyph GID, coords []float32) float32 {
	index := t.advances.getIndex(glyph)
	return t.store.GetDelta(index, coords)
}

func (t tableHVvar) getSideBearingVar(glyph GID, coords []float32) float32 {
	if t.leftSideBearings == nil {
		return 0
	}
	index := t.leftSideBearings.getIndex(glyph)
	return t.store.GetDelta(index, coords)
}

func parseTableHVvar(data []byte, axisCount int) (out tableHVvar, err error) {
	if len(data) < 20 {
		return out, errors.New("invalid metrics variation table (EOF)")
	}
	storeOffset := binary.BigEndian.Uint32(data[4:])
	advanceOffset := binary.BigEndian.Uint32(data[8:])
	lsbOffset := binary.BigEndian.Uint32(data[12:])
	out.store, err = parseVariationStore(data, storeOffset, axisCount)
	if err != nil {
		return out, err
	}
	if advanceOffset != 0 {
		out.advances, err = parseDeltaSetMapping(data, advanceOffset)
		if err != nil {
			return out, err
		}
	}
	if lsbOffset != 0 {
		out.leftSideBearings, err = parseDeltaSetMapping(data, lsbOffset)
		if err != nil {
			return out, err
		}
	}
	// we don't use the right side bearings

	return out, nil
}

// may have a length < numGlyph
type deltaSetMapping []VariationStoreIndex

func (m deltaSetMapping) getIndex(glyph GID) VariationStoreIndex {
	// If a mapping table is not provided, glyph indices are used as implicit delta-set indices.
	// [...] the delta-set outer-level index is zero, and the glyph ID is used as the inner-level index.
	if len(m) == 0 {
		return VariationStoreIndex{DeltaSetInner: uint16(glyph)}
	}

	// If a given glyph ID is greater than mapCount - 1, then the last entry is used.
	if int(glyph) >= len(m) {
		glyph = GID(len(m) - 1)
	}

	return m[glyph]
}

func parseDeltaSetMapping(data []byte, offset uint32) (deltaSetMapping, error) {
	if len(data) < int(offset)+4 {
		return nil, errors.New("invalid delta-set mapping (EOF)")
	}
	format := binary.BigEndian.Uint16(data[offset:])
	count := int(binary.BigEndian.Uint16(data[offset+2:]))
	data = data[offset+4:]

	entrySize := int((format&0x0030)>>4 + 1)
	innerBitSize := format&0x0F + 1
	if entrySize > 4 || len(data) < entrySize*count {
		return nil, errors.New("invalid delta-set mapping (EOF)")
	}
	out := make(deltaSetMapping, count)
	for i := range out {
		var v uint32
		for _, b := range data[entrySize*i : entrySize*(i+1)] { // 1 to 4 bytes
			v = v<<8 + uint32(b)
		}
		out[i].DeltaSetOuter = uint16(v >> innerBitSize)
		out[i].DeltaSetInner = uint16(v & (1<<innerBitSize - 1))
	}

	return out, nil
}

// ------------------------------------- GVAR -------------------------------------

type tableGvar struct {
	sharedTuples [][]float32          // N x axisCount
	variations   []glyphVariationData // length glyphCount
}

func nextIndex(i, start, end int) int {
	if i >= end {
		return start
	}
	return i + 1
}

func inferDelta(targetVal, prevVal, nextVal, prevDelta, nextDelta float32) float32 {
	if prevVal == nextVal {
		if prevDelta == nextDelta {
			return prevDelta
		}
		return 0
	} else if targetVal <= minF(prevVal, nextVal) {
		if prevVal < nextVal {
			return prevDelta
		}
		return nextDelta
	} else if targetVal >= maxF(prevVal, nextVal) {
		if prevVal > nextVal {
			return prevDelta
		}
		return nextDelta
	}

	/* linear interpolation */
	r := (targetVal - prevVal) / (nextVal - prevVal)
	return (1.-r)*prevDelta + r*nextDelta
}

// update `points` in place
func (t tableGvar) applyDeltasToPoints(glyph GID, coords []float32, points []contourPoint) {
	// adapted from harfbuzz/src/hb-ot-var-gvar-table.hh

	if int(glyph) >= len(t.variations) { // should not happend
		return
	}
	/* Save original points for inferred delta calculation */
	origPoints := append([]contourPoint(nil), points...)
	deltas := make([]contourPoint, len(points))

	var endPoints []int // index into points
	for i, p := range points {
		if p.isEndPoint {
			endPoints = append(endPoints, i)
		}
	}

	varData := t.variations[glyph]
	for _, tuple := range varData {
		scalar := tuple.calculateScalar(coords, t.sharedTuples)
		if scalar == 0 {
			continue
		}
		L := len(tuple.deltas)
		applyToAll := tuple.pointNumbers == nil
		xDeltas, yDeltas := tuple.deltas[:L/2], tuple.deltas[L/2:]

		// reset the current deltas
		for i := range deltas {
			deltas[i] = contourPoint{}
		}

		for i := range xDeltas {
			ptIndex := uint16(i)
			if !applyToAll {
				ptIndex = tuple.pointNumbers[i]
			}
			deltas[ptIndex].isExplicit = true
			deltas[ptIndex].X += float32(xDeltas[i]) * scalar
			deltas[ptIndex].Y += float32(yDeltas[i]) * scalar
		}

		/* infer deltas for unreferenced points */
		startPoint := 0
		for _, endPoint := range endPoints {
			// check the number of unreferenced points in a contour.
			// If no unref points or no ref points, nothing to do.
			unrefCount := 0
			for _, p := range deltas[startPoint : endPoint+1] {
				if !p.isExplicit {
					unrefCount++
				}
			}
			j := startPoint
			if unrefCount == 0 || unrefCount > endPoint-startPoint {
				goto noMoreGaps
			}

			for {
				/* Locate the next gap of unreferenced points between two referenced points prev and next.
				 * Note that a gap may wrap around at left (startPoint) and/or at right (endPoint).
				 */
				var prev, next, i int
				for {
					i = j
					j = nextIndex(i, startPoint, endPoint)
					if deltas[i].isExplicit && !deltas[j].isExplicit {
						break
					}
				}
				prev, j = i, i
				for {
					i = j
					j = nextIndex(i, startPoint, endPoint)
					if !deltas[i].isExplicit && deltas[j].isExplicit {
						break
					}
				}
				next = j
				/* Infer deltas for all unref points in the gap between prev and next */
				i = prev
				for {
					i = nextIndex(i, startPoint, endPoint)
					if i == next {
						break
					}
					deltas[i].X = inferDelta(origPoints[i].X, origPoints[prev].X, origPoints[next].X, deltas[prev].X, deltas[next].X)
					deltas[i].Y = inferDelta(origPoints[i].Y, origPoints[prev].Y, origPoints[next].Y, deltas[prev].Y, deltas[next].Y)
					unrefCount--
					if unrefCount == 0 {
						goto noMoreGaps
					}
				}
			}
		noMoreGaps:
			startPoint = endPoint + 1
		}

		/* apply specified / inferred deltas to points */
		for i, d := range deltas {
			points[i].X += d.X
			points[i].Y += d.Y
		}
	}
}

// axisCountRef, glyphCountRef are used to sanitize
func parseTableGvar(data []byte, axisCountRef int, glyphs TableGlyf) (out tableGvar, err error) {
	if len(data) < 20 {
		return out, errors.New("invalid 'gvar' table (EOF)")
	}
	axisCount := int(binary.BigEndian.Uint16(data[4:]))
	sharedTupleCount := binary.BigEndian.Uint16(data[6:])
	sharedTupleOffset := int(binary.BigEndian.Uint32(data[8:]))
	glyphCount := int(binary.BigEndian.Uint16(data[12:]))
	flags := binary.BigEndian.Uint16(data[14:])
	glyphVariationDataArrayOffset := int(binary.BigEndian.Uint32(data[16:]))

	if axisCount != axisCountRef {
		return out, errors.New("invalid 'gvar' table (EOF)")
	}
	if glyphCount != len(glyphs) {
		return out, errors.New("invalid 'gvar' table (EOF)")
	}

	offsets, err := parseTableLoca(data[20:], glyphCount, flags&1 != 0)
	if err != nil {
		return out, fmt.Errorf("invalid 'gvar' table: %s", err)
	}

	out.sharedTuples, err = parseSharedTuples(data, sharedTupleOffset, axisCount, int(sharedTupleCount))
	if err != nil {
		return out, err
	}

	if len(data) < glyphVariationDataArrayOffset {
		return out, errors.New("invalid 'gvar' table (EOF)")
	}
	startDataVariations := data[glyphVariationDataArrayOffset:]
	out.variations = make([]glyphVariationData, glyphCount)
	for i := range out.variations {
		if offsets[i] == offsets[i+1] {
			continue
		}

		out.variations[i], err = parseOneGlyphVariationData(startDataVariations[:offsets[i+1]], offsets[i], false,
			axisCount, glyphs[i].pointNumbersCount()+phantomCount)
		if err != nil {
			return out, err
		}
	}

	return out, nil
}

func parseSharedTuples(data []byte, offset, axisCount, sharedTupleCount int) ([][]float32, error) {
	if len(data) < offset+axisCount*2*sharedTupleCount {
		return nil, errors.New("invalid 'gvar' table (EOF)")
	}
	out := make([][]float32, sharedTupleCount)
	for i := range out {
		out[i] = parseTupleRecord(data[offset+axisCount*2*i:], axisCount)
	}
	return out, nil
}

// length as already been checked
func parseTupleRecord(data []byte, axisCount int) []float32 {
	vi := make([]float32, axisCount)
	for j := range vi {
		vi[j] = fixed214ToFloat(binary.BigEndian.Uint16(data[2*j:]))
	}
	return vi
}

type glyphVariationData []tupleVariation

// offset is at the beginning of the table
// if isCvar is true, the version fields are ignored
// pointNumbersCount includes the phantom points
func parseOneGlyphVariationData(data []byte, offset uint32, isCvar bool, axisCount, pointNumbersCount int) (glyphVariationData, error) {
	headerSize := 4
	if isCvar {
		headerSize = 8
	}

	if len(data) < int(offset)+headerSize {
		return nil, errors.New("invalid glyph variation data (EOF)")
	}
	data = data[offset:]
	tupleVariationCount := binary.BigEndian.Uint16(data[headerSize-4:]) // 0 or 4
	dataOffset := binary.BigEndian.Uint16(data[headerSize-2:])          // 2 or 6
	if len(data) < int(dataOffset) {
		return nil, errors.New("invalid glyph variation data (EOF)")
	}
	serializedData := data[dataOffset:]

	const (
		sharedPointNumbers = 0x8000
		countMask          = 0x0FFF
	)
	tupleCount := tupleVariationCount & countMask

	out := make(glyphVariationData, tupleCount) // allocation guarded by countMask
	data = data[headerSize:]
	var err error
	for i := range out {
		out[i].tupleVariationHeader, data, err = parseTupleVariationHeader(data, isCvar, axisCount)
		if err != nil {
			return out, err
		}
	}

	hasSharedPointNumbers := tupleVariationCount&sharedPointNumbers != 0
	err = parseGlyphVariationSerializedData(serializedData, hasSharedPointNumbers, pointNumbersCount, isCvar, out)

	return out, err
}

type tupleVariationHeader struct {
	peakTuple              []float32 // nil or with length axisCount
	intermediateStartTuple []float32 // nil or with length axisCount
	intermediateEndTuple   []float32 // nil or with length axisCount

	variationDataSize uint16 // usefull only during parsing

	tupleIndex uint16
}

func (t *tupleVariationHeader) hasPrivatePointNumbers() bool {
	const privatePointNumbers = 0x2000
	return t.tupleIndex&privatePointNumbers != 0
}

func (t *tupleVariationHeader) getIndex() uint16 {
	const TupleIndexMask = 0x0FFF
	return t.tupleIndex & TupleIndexMask
}

// sharedTuples has length _ x axisCount
func (t tupleVariationHeader) calculateScalar(coords []float32, sharedTuples [][]float32) float32 {
	peakTuple := t.peakTuple
	if peakTuple == nil { // use shared tuple
		index := t.getIndex()
		if int(index) >= len(sharedTuples) { // should not happend
			return 0.
		}
		peakTuple = sharedTuples[index]
	}

	startTuple, endTuple := t.intermediateStartTuple, t.intermediateEndTuple
	hasIntermediate := startTuple != nil

	var scalar float32 = 1.
	for i, v := range coords {
		peak := peakTuple[i]
		if peak == 0 || v == peak {
			continue
		}

		if hasIntermediate {
			start := startTuple[i]
			end := endTuple[i]
			if start > peak || peak > end || (start < 0 && end > 0 && peak != 0) {
				continue
			}
			if v < start || v > end {
				return 0.
			}
			if v < peak {
				if peak != start {
					scalar *= (v - start) / (peak - start)
				}
			} else {
				if peak != end {
					scalar *= (end - v) / (end - peak)
				}
			}
		} else if v == 0 || v < minF(0, peak) || v > maxF(0, peak) {
			return 0.
		} else {
			scalar *= v / peak
		}
	}
	return scalar
}

// return data after the tuple header
func parseTupleVariationHeader(data []byte, isCvar bool, axisCount int) (out tupleVariationHeader, _ []byte, err error) {
	if len(data) < 4 {
		return out, nil, errors.New("invalid tuple variation header (EOF)")
	}
	out.variationDataSize = binary.BigEndian.Uint16(data)
	out.tupleIndex = binary.BigEndian.Uint16(data[2:])

	const (
		embeddedPeakTuple  = 0x8000
		intermediateRegion = 0x4000
	)
	hasPeak := out.tupleIndex&embeddedPeakTuple != 0
	hasRegions := out.tupleIndex&intermediateRegion != 0
	if isCvar && !hasPeak {
		return out, nil, errors.New("invalid tuple variation header for 'cvar' table")
	}

	data = data[4:]

	if hasPeak {
		if len(data) < 2*axisCount {
			return out, nil, errors.New("invalid glyph variation data (EOF)")
		}
		out.peakTuple = parseTupleRecord(data, axisCount)
		data = data[2*axisCount:]
	}
	if hasRegions {
		if len(data) < 4*axisCount {
			return out, nil, errors.New("invalid glyph variation data (EOF)")
		}
		out.intermediateStartTuple = parseTupleRecord(data, axisCount)
		out.intermediateEndTuple = parseTupleRecord(data[2*axisCount:], axisCount)
		data = data[4*axisCount:]
	}
	return out, data, nil
}

type tupleVariation struct {
	pointNumbers []uint16 // nil means allPointsNumbers
	// length 2*len(pointNumbers) for gvar table or 2*allPointsNumbers if zero
	deltas []int16
	tupleVariationHeader
}

// complete `out`, which contains the parsed tuple headers.
// pointNumbersCountAll is used when the tuple variation data provides deltas for all glyph points
func parseGlyphVariationSerializedData(data []byte, hasSharedPoints bool, pointNumbersCountAll int, isCvar bool, out []tupleVariation) error {
	var (
		sharedPointNumbers []uint16
		err                error
	)
	if hasSharedPoints {
		sharedPointNumbers, data, err = parsePointNumbers(data)
		if err != nil {
			return err
		}
	}

	for i, h := range out {
		// adjust for the next iteration
		if len(data) < int(h.variationDataSize) {
			return errors.New("invalid glyph variation serialized data (EOF)")
		}
		nextData := data[h.variationDataSize:]

		// default to shared points
		privatePointNumbers := sharedPointNumbers
		if h.hasPrivatePointNumbers() {
			privatePointNumbers, data, err = parsePointNumbers(data)
			if err != nil {
				return err
			}
		}
		// the number of point is precised or defaut to all the points
		pointCount := pointNumbersCountAll
		if privatePointNumbers != nil {
			pointCount = len(privatePointNumbers)
		}

		out[i].pointNumbers = privatePointNumbers

		if !isCvar {
			pointCount *= 2 // for X and Y
		}

		out[i].deltas, err = unpackDeltas(data, pointCount)
		if err != nil {
			return err
		}

		data = nextData
	}
	return nil
}

// the returned slice is nil if all glyph points are used
func parsePointNumbers(data []byte) ([]uint16, []byte, error) {
	count, data, err := getPackedPointCount(data)
	if err != nil {
		return nil, nil, err
	}
	if count == 0 {
		return nil, data, nil
	}

	var lastPoint uint16
	points := make([]uint16, 0, count) // max value of count is 32767
	for len(points) < int(count) {     // loop through the runs
		if len(data) == 0 {
			return nil, nil, errors.New("invalid glyph variation points numbers (EOF)")
		}
		control := data[0]
		is16bit := control&0x80 != 0
		runLength := int(control&0x7F + 1)
		if is16bit {
			pts, err := parseUint16s(data[1:], runLength)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid glyph variation points numbers: %s", err)
			}
			for _, pt := range pts {
				actualValue := pt + lastPoint
				points = append(points, actualValue)
				lastPoint = actualValue
			}
			data = data[1+2*runLength:]
		} else {
			if len(data) < 1+runLength {
				return nil, nil, errors.New("invalid glyph variation points numbers (EOF)")
			}
			for _, b := range data[1 : 1+runLength] {
				actualValue := uint16(b) + lastPoint
				points = append(points, actualValue)
				lastPoint = actualValue
			}
			data = data[1+runLength:]
		}
	}

	return points, data, nil
}

// return the remaining data and special case of 00
func getPackedPointCount(data []byte) (uint16, []byte, error) {
	const highOrderBit byte = 1 << 7
	if len(data) < 1 {
		return 0, nil, errors.New("invalid glyph variation points numbers (EOF)")
	}
	if data[0] == 0 {
		return 0, data[1:], nil
	} else if data[0]&highOrderBit == 0 {
		count := uint16(data[0])
		return count, data[1:], nil
	} else {
		if len(data) < 2 {
			return 0, nil, errors.New("invalid glyph variation points numbers (EOF)")
		}
		count := uint16(data[0]&^highOrderBit)<<8 | uint16(data[1])
		return count, data[2:], nil
	}
}

func unpackDeltas(data []byte, pointNumbersCount int) ([]int16, error) {
	const (
		deltasAreZero     = 0x80
		deltasAreWords    = 0x40
		deltaRunCountMask = 0x3F
	)
	var out []int16
	// The data is read until the expected logic count of deltas is obtained.
	for len(out) < pointNumbersCount {
		if len(data) == 0 {
			return nil, errors.New("invalid packed deltas (EOF)")
		}
		control := data[0]
		count := control&deltaRunCountMask + 1
		if isZero := control&deltasAreZero != 0; isZero {
			//  no additional value to read, just fill with zeros
			out = append(out, make([]int16, count)...)
			data = data[1:]
		} else {
			isInt16 := control&deltasAreWords != 0
			if isInt16 {
				if len(data) < 1+2*int(count) {
					return nil, errors.New("invalid packed deltas (EOF)")
				}
				for i := byte(0); i < count; i++ { // count < 64 -> no overflow
					out = append(out, int16(binary.BigEndian.Uint16(data[1+2*i:])))
				}
				data = data[1+2*count:]
			} else {
				if len(data) < 1+int(count) {
					return nil, errors.New("invalid packed deltas (EOF)")
				}
				for i := byte(0); i < count; i++ { // count < 64 -> no overflow
					out = append(out, int16(int8(data[1+i])))
				}
				data = data[1+count:]
			}
		}
	}
	return out, nil
}
