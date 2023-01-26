package graphite

import "fmt"

/** Used for looking up slot attributes. Most are already available in other functions **/
type attrCode uint8

const (
	/// adjusted glyph advance in x direction in design units
	acAdvX attrCode = iota
	/// adjusted glyph advance in y direction (usually 0) in design units
	acAdvY
	/// returns 0. Deprecated.
	acAttTo
	/// This slot attaches to its parent at the given design units in the x direction
	acAttX
	/// This slot attaches to its parent at the given design units in the y direction
	acAttY
	/// This slot attaches to its parent at the given glyph point (not implemented)
	_
	/// x-direction adjustment from the given glyph point (not implemented)
	acAttXOff
	/// y-direction adjustment from the given glyph point (not implemented)
	acAttYOff
	/// Where on this glyph should align with the attachment point on the parent glyph in the x-direction.
	acAttWithX
	/// Where on this glyph should align with the attachment point on the parent glyph in the y-direction
	acAttWithY
	/// Which glyph point on this glyph should align with the attachment point on the parent glyph (not implemented).
	_
	/// Adjustment to acWithGpt in x-direction (not implemented)
	acAttWithXOff
	/// Adjustment to acWithGpt in y-direction (not implemented)
	acAttWithYOff
	/// Attach at given nesting level (not implemented)
	acAttLevel
	/// Line break breakweight for this glyph
	acBreak
	/// Ligature component reference (not implemented)
	acCompRef
	/// bidi directionality of this glyph (not implemented)
	acDir
	/// Whether insertion is allowed before this glyph
	acInsert
	/// Final positioned position of this glyph relative to its parent in x-direction in pixels
	acPosX
	/// Final positioned position of this glyph relative to its parent in y-direction in pixels
	acPosY
	/// Amount to shift glyph by in x-direction design units
	acShiftX
	/// Amount to shift glyph by in y-direction design units
	acShiftY
	/// attribute user1
	acUserDefnV1
	/// not implemented
	acMeasureSol
	/// not implemented
	acMeasureEol
	/// Amount this slot can stretch (not implemented)
	acJStretch
	/// Amount this slot can shrink (not implemented)
	_
	/// Granularity by which this slot can stretch or shrink (not implemented)
	_
	/// Justification weight for this glyph (not implemented)
	_
	/// Amount this slot mush shrink or stretch in design units
	acJWidth
	/// SubSegment split point
	acSegSplit = acJStretch + 29
	/// User defined attribute, see subattr for user attr number
	acUserDefn = 55
	/// Bidi level
	acBidiLevel = 24 + iota
	/// Collision flags
	acColFlags
	/// Collision constraint rectangle left (bl.x)
	acColLimitblx
	/// Collision constraint rectangle lower (bl.y)
	acColLimitbly
	/// Collision constraint rectangle right (tr.x)
	acColLimittrx
	/// Collision constraint rectangle upper (tr.y)
	acColLimittry
	/// Collision shift x
	acColShiftx
	/// Collision shift y
	acColShifty
	/// Collision margin
	acColMargin
	/// Margin cost weight
	acColMarginWt
	// Additional glyph that excludes movement near this one:
	acColExclGlyph
	acColExclOffx
	acColExclOffy
	// Collision sequence enforcing attributes:
	acSeqClass
	acSeqProxClass
	acSeqOrder
	acSeqAboveXoff
	acSeqAboveWt
	acSeqBelowXlim
	acSeqBelowWt
	acSeqValignHt
	acSeqValignWt

	/// not implemented
	acMax
	/// not implemented
	_ = acMax + 1
)

type opcode uint8

// opcodes
const (
	ocNOP opcode = iota

	ocPUSH_BYTE
	ocPUSH_BYTEU
	ocPUSH_SHORT
	ocPUSH_SHORTU
	ocPUSH_LONG

	ocADD
	ocSUB
	ocMUL
	ocDIV

	ocMIN
	ocMAX

	ocNEG

	ocTRUNC8
	ocTRUNC16

	ocCOND

	ocAND
	ocOR
	ocNOT

	ocEQUAL
	ocNOT_EQ

	ocLESS
	ocGTR
	ocLESS_EQ
	ocGTR_EQ

	ocNEXT
	ocNEXT_N
	ocCOPY_NEXT

	ocPUT_GLYPH_8BIT_OBS
	ocPUT_SUBS_8BIT_OBS
	ocPUT_COPY

	ocINSERT
	ocDELETE

	ocASSOC

	ocCNTXT_ITEM

	ocATTR_SET
	ocATTR_ADD
	ocATTR_SUB

	ocATTR_SET_SLOT

	ocIATTR_SET_SLOT

	ocPUSH_SLOT_ATTR
	ocPUSH_GLYPH_ATTR_OBS

	ocPUSH_GLYPH_METRIC
	ocPUSH_FEAT

	ocPUSH_ATT_TO_GATTR_OBS
	ocPUSH_ATT_TO_GLYPH_METRIC

	ocPUSH_ISLOT_ATTR

	ocPUSH_IGLYPH_ATTR // not implemented

	ocPOP_RET
	ocRET_ZERO
	ocRET_TRUE

	ocIATTR_SET
	ocIATTR_ADD
	ocIATTR_SUB

	ocPUSH_PROC_STATE
	ocPUSH_VERSION

	ocPUT_SUBS
	ocPUT_SUBS2
	ocPUT_SUBS3

	ocPUT_GLYPH
	ocPUSH_GLYPH_ATTR
	ocPUSH_ATT_TO_GLYPH_ATTR

	ocBITOR
	ocBITAND
	ocBITNOT

	ocBITSET
	ocSET_FEAT

	ocMAX_OPCODE
	// private opcodes for internal use only, comes after all other on disk opcodes
	ocTEMP_COPY = ocMAX_OPCODE
)

func (opc opcode) String() string {
	if int(opc) < len(opcodeTable) {
		return opcodeTable[opc].name
	}
	return fmt.Sprintf("<unknown opcode: %d>", opc)
}

func (opc opcode) isReturn() bool {
	return opc == ocPOP_RET || opc == ocRET_ZERO || opc == ocRET_TRUE
}

type errorStatusCode uint8

const (
	invalidOpCode errorStatusCode = iota
	unimplementedOpCodeUsed
	outOfRangeData
	jumpPastEnd
	argumentsExhausted
	missingReturn
	nestedContextItem
	underfullStack
)

func (c errorStatusCode) Error() string {
	switch c {
	case invalidOpCode:
		return "invalid opcode"
	case unimplementedOpCodeUsed:
		return "unimplemented opcode used"
	case outOfRangeData:
		return "out of range data"
	case jumpPastEnd:
		return "jump past end"
	case argumentsExhausted:
		return "arguments exhausted"
	case missingReturn:
		return "missing return"
	case nestedContextItem:
		return "nested context item"
	case underfullStack:
		return "underfull stack"
	}
	return "<unknown error code>"
}

// instr is an op code with the selected implementation
type instr struct {
	fn   instrImpl
	code opcode
}

// represents loaded graphite stack machine code
type code struct {
	instrs []instr
	args   []byte // concatenated arguments for `instrs`

	// instr *     _code;
	// byte  *     _data;
	// size_t      _data_size,
	// instrCount int

	maxRef int // maximum index of slot encountered in the instructions

	// status                     codeStatus

	constraint, delete, modify bool

	// mutable bool _own;
}

// settings from the font needed to load a bytecode sequence
type codeContext struct {
	NumClasses, NumAttributes, NumFeatures uint16
	NumUserAttributes                      uint8
	Pt                                     passtype
}

// newCode decodes an input and returns the loaded instructions
// the errors returns are of type errorStatusCode
// If skipAnalysis is true, the required TEMP_COPY opcodes are not inserted.
// This is only useful to get a direct representation of the font file content required in tests.
func newCode(isConstraint bool, bytecode []byte,
	preContext uint8, ruleLength uint16, context codeContext, skipAnalysis bool,
) (code, error) {
	if len(bytecode) == 0 {
		return code{}, nil
	}

	lims := decoderLimits{
		preContext: uint16(preContext),
		ruleLength: ruleLength,
		classes:    context.NumClasses,
		glyfAttrs:  context.NumAttributes,
		features:   context.NumFeatures,
		attrid: [acMax]byte{
			1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 1, 255,
			1, 1, 1, 1, 1, 1, 1, 1,
			1, 1, 1, 1, 1, 1, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, context.NumUserAttributes, // 0, 0, etc...
		},
	}

	dec := newDecoder(isConstraint, lims, context.Pt)
	// parse the bytecodes
	lastOpcode, err := dec.load(bytecode)
	if err != nil {
		return code{}, err
	}

	// Is this an empty program?
	if len(dec.code.instrs) == 0 {
		return dec.code, nil
	}

	// When we reach the end check we've terminated it correctly
	if !lastOpcode.isReturn() {
		return dec.code, missingReturn
	}

	if !skipAnalysis {
		dec.code.instrs = dec.applyAnalysis(dec.code.instrs)
	}

	return dec.code, nil
}

const decoderNUMCONTEXTS = 256

type decoderLimits struct {
	preContext                   uint16
	ruleLength                   uint16
	classes, glyfAttrs, features uint16
	attrid                       [acMax]byte
}

type context struct {
	codeRef             uint8
	changed, referenced bool
}

// decoder is responsible for reading a sequence of instructions.
// This is the first step of the use of the code sequence from a font:
//   - at font creation (see computeRules), a first analyze of the sequence if done,
//     checking for some errors like out of bounds access
//   - at runtime, the code is run against an input
type decoder struct {
	code code // resulting loaded code

	stackDepth          int // the number of element needed in stack
	outIndex, outLength int
	slotRef             int
	inCtxtItem          bool
	passtype            passtype

	contexts [decoderNUMCONTEXTS]context

	max decoderLimits
}

func newDecoder(isConstraint bool, max decoderLimits, pt passtype) *decoder {
	out := decoder{code: code{constraint: isConstraint}, max: max, passtype: pt}
	out.outLength = 1
	if !out.code.constraint {
		out.outIndex = int(max.preContext)
		out.outLength = int(max.ruleLength)
	}
	return &out
}

// parses the input byte code and checks that operations
// and arguments are valid
// note that the instructions are NOT executed
// the last opcode is returned
func (dec *decoder) load(bytecode []byte) (opcode, error) {
	var lastOpcode opcode
	for len(bytecode) != 0 {
		opc, err := dec.fetchOpcode(bytecode)
		if err != nil {
			return 0, err
		}

		bytecode = bytecode[1:]
		dec.analyseOpcode(opc, bytecode)

		bytecode, err = dec.emitOpcode(opc, bytecode)
		if err != nil {
			return 0, err
		}

		lastOpcode = opc
	}

	return lastOpcode, nil
}

func (dec *decoder) validateOpcode(opc opcode, bc []byte) error {
	if opc >= ocMAX_OPCODE {
		return invalidOpCode
	}
	op := opcodeTable[opc]
	if op.impl[boolToInt(dec.code.constraint)] == nil {
		return unimplementedOpCodeUsed
	}
	if op.paramSize == varArgs && len(bc) == 0 {
		return argumentsExhausted
	}
	paramSize := op.paramSize
	if op.paramSize == varArgs { // read the number of additional args as first arg
		paramSize = bc[0] + 1
	}
	if len(bc) < int(paramSize) {
		return argumentsExhausted
	}
	return nil
}

// bc is not empty
func (dec *decoder) fetchOpcode(bc []byte) (opcode, error) {
	opc := opcode(bc[0])
	bc = bc[1:]

	// Do some basic sanity checks based on what we know about the opcode
	if err := dec.validateOpcode(opc, bc); err != nil {
		return 0, err
	}

	// And check its arguments as far as possible
	switch opc {
	case ocNOP:
	case ocPUSH_BYTE, ocPUSH_BYTEU, ocPUSH_SHORT, ocPUSH_SHORTU, ocPUSH_LONG:
		dec.stackDepth++
	case ocADD, ocSUB, ocMUL, ocDIV, ocMIN, ocMAX, ocAND, ocOR,
		ocEQUAL, ocNOT_EQ, ocLESS, ocGTR, ocLESS_EQ, ocGTR_EQ, ocBITOR, ocBITAND:
		dec.stackDepth--
		if dec.stackDepth <= 0 {
			return 0, underfullStack
		}
	case ocNEG, ocTRUNC8, ocTRUNC16, ocNOT, ocBITNOT, ocBITSET:
		if dec.stackDepth <= 0 {
			return 0, underfullStack
		}
	case ocCOND:
		dec.stackDepth -= 2
		if dec.stackDepth <= 0 {
			return 0, underfullStack
		}
	case ocNEXT_N:
	// runtime checked
	case ocNEXT, ocCOPY_NEXT:
		dec.outIndex++
		if dec.outIndex < -1 || dec.outIndex > dec.outLength || dec.slotRef > int(dec.max.ruleLength) {
			return 0, outOfRangeData
		}
	case ocPUT_GLYPH_8BIT_OBS:
		if err := validUpto(dec.max.classes, uint16(bc[0])); err != nil {
			return 0, err
		}
		if err := dec.testContext(); err != nil {
			return 0, err
		}
	case ocPUT_SUBS_8BIT_OBS:
		if err := dec.testRef(bc[0]); err != nil {
			return 0, err
		}
		if err := validUpto(dec.max.classes, uint16(bc[1])); err != nil {
			return 0, err
		}
		if err := validUpto(dec.max.classes, uint16(bc[2])); err != nil {
			return 0, err
		}
		if err := dec.testContext(); err != nil {
			return 0, err
		}
	case ocPUT_COPY:
		if err := dec.testRef(bc[0]); err != nil {
			return 0, err
		}
		if err := dec.testContext(); err != nil {
			return 0, err
		}
	case ocINSERT:
		if dec.passtype >= ptPOSITIONING {
			return 0, invalidOpCode
		}
		dec.outLength++
		if dec.outIndex < 0 {
			dec.outIndex++
		}
		if dec.outIndex < -1 || dec.outIndex >= dec.outLength {
			return 0, outOfRangeData
		}
	case ocDELETE:
		if dec.passtype >= ptPOSITIONING {
			return 0, invalidOpCode
		}
		if dec.outIndex < int(dec.max.preContext) {
			return 0, outOfRangeData
		}
		dec.outIndex--
		dec.outLength--
		if dec.outIndex < -1 || dec.outIndex > dec.outLength {
			return 0, outOfRangeData
		}
	case ocASSOC:
		if bc[0] == 0 {
			return 0, outOfRangeData
		}
		for num := bc[0]; num != 0; num-- {
			if err := dec.testRef(bc[num]); err != nil {
				return 0, err
			}
		}
		if err := dec.testContext(); err != nil {
			return 0, err
		}
	case ocCNTXT_ITEM:
		err := validUpto(dec.max.ruleLength, uint16(int(dec.max.preContext)+int(int8(bc[0]))))
		if err != nil {
			return 0, err
		}
		if len(bc) < 2+int(bc[1]) {
			return 0, jumpPastEnd
		}
		if dec.inCtxtItem {
			return 0, nestedContextItem
		}
	case ocATTR_SET, ocATTR_ADD, ocATTR_SUB, ocATTR_SET_SLOT:
		dec.stackDepth--
		if dec.stackDepth < 0 {
			return 0, underfullStack
		}
		err := validUpto(acMax, uint16(bc[0]))
		if err != nil {
			return 0, err
		}
		if attrCode(bc[0]) == acUserDefn { // use IATTR for user attributes
			return 0, outOfRangeData
		}
		if err := dec.testAttr(attrCode(bc[0])); err != nil {
			return 0, err
		}
		if err := dec.testContext(); err != nil {
			return 0, err
		}
	case ocIATTR_SET_SLOT:
		dec.stackDepth--
		if dec.stackDepth < 0 {
			return 0, underfullStack
		}
		if err := validUpto(acMax, uint16(bc[0])); err != nil {
			return 0, err
		}
		if err := validUpto(uint16(dec.max.attrid[bc[0]]), uint16(bc[1])); err != nil {
			return 0, err
		}
		if err := dec.testAttr(attrCode(bc[0])); err != nil {
			return 0, err
		}
		if err := dec.testContext(); err != nil {
			return 0, err
		}
	case ocPUSH_SLOT_ATTR:
		dec.stackDepth++
		if err := validUpto(acMax, uint16(bc[0])); err != nil {
			return 0, err
		}
		if err := dec.testRef(bc[1]); err != nil {
			return 0, err
		}
		if attrCode(bc[0]) == acUserDefn { // use IATTR for user attributes
			return 0, outOfRangeData
		}
		if err := dec.testAttr(attrCode(bc[0])); err != nil {
			return 0, err
		}
	case ocPUSH_GLYPH_ATTR_OBS, ocPUSH_ATT_TO_GATTR_OBS:
		dec.stackDepth++
		if err := validUpto(dec.max.glyfAttrs, uint16(bc[0])); err != nil {
			return 0, err
		}
		if err := dec.testRef(bc[1]); err != nil {
			return 0, err
		}
	case ocPUSH_ATT_TO_GLYPH_METRIC, ocPUSH_GLYPH_METRIC:
		dec.stackDepth++
		if err := validUpto(kgmetDescent, uint16(bc[0])); err != nil {
			return 0, err
		}
		if err := dec.testRef(bc[1]); err != nil {
			return 0, err
		}
	// level: dp[2] no check necessary
	case ocPUSH_FEAT:
		dec.stackDepth++
		if err := validUpto(dec.max.features, uint16(bc[0])); err != nil {
			return 0, err
		}
		if err := dec.testRef(bc[1]); err != nil {
			return 0, err
		}
	case ocPUSH_ISLOT_ATTR:
		dec.stackDepth++
		if err := validUpto(acMax, uint16(bc[0])); err != nil {
			return 0, err
		}
		if err := dec.testRef(bc[1]); err != nil {
			return 0, err
		}
		if err := validUpto(uint16(dec.max.attrid[bc[0]]), uint16(bc[2])); err != nil {
			return 0, err
		}
		if err := dec.testAttr(attrCode(bc[0])); err != nil {
			return 0, err
		}
	case ocPUSH_IGLYPH_ATTR:
		// not implemented
		dec.stackDepth++
	case ocPOP_RET:
		dec.stackDepth--
		if dec.stackDepth < 0 {
			return 0, underfullStack
		}
		fallthrough
	case ocRET_ZERO, ocRET_TRUE:
	case ocIATTR_SET, ocIATTR_ADD, ocIATTR_SUB:
		dec.stackDepth--
		if dec.stackDepth < 0 {
			return 0, underfullStack
		}
		if err := validUpto(acMax, uint16(bc[0])); err != nil {
			return 0, err
		}
		if err := validUpto(uint16(dec.max.attrid[bc[0]]), uint16(bc[1])); err != nil {
			return 0, err
		}
		if err := dec.testAttr(attrCode(bc[0])); err != nil {
			return 0, err
		}
		if err := dec.testContext(); err != nil {
			return 0, err
		}
	case ocPUSH_PROC_STATE, ocPUSH_VERSION:
		dec.stackDepth++
	case ocPUT_SUBS:
		if err := dec.testRef(bc[0]); err != nil {
			return 0, err
		}
		if err := validUpto(dec.max.classes, uint16(bc[1])<<8|uint16(bc[2])); err != nil {
			return 0, err
		}
		if err := validUpto(dec.max.classes, uint16(bc[3])<<8|uint16(bc[4])); err != nil {
			return 0, err
		}
		if err := dec.testContext(); err != nil {
			return 0, err
		}
	case ocPUT_SUBS2, ocPUT_SUBS3:
	// not implemented
	case ocPUT_GLYPH:
		if err := validUpto(dec.max.classes, uint16(bc[0])<<8|uint16(bc[1])); err != nil {
			return 0, err
		}
		if err := dec.testContext(); err != nil {
			return 0, err
		}
	case ocPUSH_GLYPH_ATTR, ocPUSH_ATT_TO_GLYPH_ATTR:
		dec.stackDepth++
		if err := validUpto(dec.max.glyfAttrs, uint16(bc[0])<<8|uint16(bc[1])); err != nil {
			return 0, err
		}
		if err := dec.testRef(bc[2]); err != nil {
			return 0, err
		}
	case ocSET_FEAT:
		if err := validUpto(dec.max.features, uint16(bc[0])); err != nil {
			return 0, err
		}
		if err := dec.testRef(bc[1]); err != nil {
			return 0, err
		}
	default:
		return 0, invalidOpCode
	}

	return opc, nil
}

func validUpto(limit, x uint16) error {
	if (limit != 0) && (x < limit) {
		return nil
	}
	return outOfRangeData
}

func (dec *decoder) testContext() error {
	if dec.outIndex >= dec.outLength || dec.outIndex < 0 || dec.slotRef >= decoderNUMCONTEXTS-1 {
		return outOfRangeData
	}
	return nil
}

func (dec *decoder) testRef(index_ byte) error {
	index := int8(index_)
	if dec.code.constraint && !dec.inCtxtItem {
		if index > 0 || uint16(-index) > dec.max.preContext {
			return outOfRangeData
		}
	} else {
		if L := dec.slotRef + int(dec.max.preContext) + int(index); dec.max.ruleLength == 0 ||
			L >= int(dec.max.ruleLength) || L < 0 {
			return outOfRangeData
		}
	}
	return nil
}

func (dec *decoder) testAttr(attr attrCode) error {
	if dec.passtype < ptPOSITIONING {
		if attr != acBreak && attr != acDir && attr != acUserDefn && attr != acCompRef {
			return outOfRangeData
		}
	}
	return nil
}

// the length of arg as been checked
func (dec *decoder) analyseOpcode(opc opcode, arg []byte) {
	switch opc {
	case ocDELETE:
		dec.code.delete = true
	case ocASSOC:
		dec.setChanged(0)
		//      for (uint8 num = arg[0]; num; --num)
		//        _analysis.setNoref(num);
	case ocPUT_GLYPH_8BIT_OBS, ocPUT_GLYPH:
		dec.code.modify = true
		dec.setChanged(0)
	case ocATTR_SET, ocATTR_ADD, ocATTR_SUB, ocATTR_SET_SLOT, ocIATTR_SET_SLOT, ocIATTR_SET, ocIATTR_ADD, ocIATTR_SUB:
		dec.setNoref(0)
	case ocNEXT, ocCOPY_NEXT:
		dec.slotRef++
		dec.contexts[dec.slotRef] = context{codeRef: uint8(len(dec.code.instrs) + 1)}
		// if (_analysis.slotRef > _analysis.max_ref) _analysis.max_ref = _analysis.slotRef;
	case ocINSERT:
		if dec.slotRef >= 0 {
			dec.slotRef--
		}
		dec.code.modify = true
	case ocPUT_SUBS_8BIT_OBS /* slotRef on 1st parameter */, ocPUT_SUBS:
		dec.code.modify = true
		dec.setChanged(0)
		fallthrough
	case ocPUT_COPY:
		if arg[0] != 0 {
			dec.setChanged(0)
			dec.code.modify = true
		}
		dec.setRef(int8(arg[0]))
	case ocPUSH_GLYPH_ATTR_OBS, ocPUSH_SLOT_ATTR, ocPUSH_GLYPH_METRIC, ocPUSH_ATT_TO_GATTR_OBS, ocPUSH_ATT_TO_GLYPH_METRIC, ocPUSH_ISLOT_ATTR, ocPUSH_FEAT, ocSET_FEAT:
		dec.setRef(int8(arg[1]))
	case ocPUSH_ATT_TO_GLYPH_ATTR, ocPUSH_GLYPH_ATTR:
		dec.setRef(int8(arg[2]))
	}
}

func (dec *decoder) setRef(arg int8) {
	index := int(arg)
	if index+dec.slotRef < 0 || index+dec.slotRef >= decoderNUMCONTEXTS {
		return
	}
	dec.contexts[index+dec.slotRef].referenced = true
	if index+dec.slotRef > dec.code.maxRef {
		dec.code.maxRef = index + dec.slotRef
	}
}

func (dec *decoder) setNoref(index int) {
	if index+dec.slotRef < 0 || index+dec.slotRef >= decoderNUMCONTEXTS {
		return
	}
	if index+dec.slotRef > dec.code.maxRef {
		dec.code.maxRef = index + dec.slotRef
	}
}

func (dec *decoder) setChanged(index int) {
	if index+dec.slotRef < 0 || index+dec.slotRef >= decoderNUMCONTEXTS {
		return
	}
	dec.contexts[index+dec.slotRef].changed = true
	if index+dec.slotRef > dec.code.maxRef {
		dec.code.maxRef = index + dec.slotRef
	}
}

// implements one op code.
// `args` has already been checked for its length (see `fetchOpcode` for exceptions)
// the function returns the truncated `args` slice and `true` if no error occured
type instrImpl func(reg *regbank, st *stack, args []byte) ([]byte, bool)

// length of bc has been checked
// the `code` item will be updated, and the remaining bytecode
// input is returned
func (dec *decoder) emitOpcode(opc opcode, bc []byte) ([]byte, error) {
	op := opcodeTable[opc]
	fn := op.impl[boolToInt(dec.code.constraint)]
	if fn == nil {
		return nil, unimplementedOpCodeUsed
	}

	paramSize := op.paramSize
	if op.paramSize == varArgs {
		paramSize = bc[0] + 1
	}

	// Add this instruction
	dec.code.instrs = append(dec.code.instrs, instr{
		fn:   fn,
		code: opc,
	})

	// Grab the parameters
	if paramSize != 0 {
		dec.code.args = append(dec.code.args, bc[:paramSize]...)
		bc = bc[paramSize:]
	}

	// recursively decode a context item so we can split the skip into
	// instruction and data portions.
	if opc == ocCNTXT_ITEM {
		// assert(_out_index == 0);
		dec.inCtxtItem = true
		dec.slotRef = int(int8(dec.code.args[len(dec.code.args)-2]))
		dec.outIndex = int(dec.max.preContext) + dec.slotRef
		dec.outLength = int(dec.max.ruleLength)

		// instrSkip takes into account the opcodes
		instrSkipIndex := len(dec.code.args) - 1
		instrSkip := dec.code.args[instrSkipIndex]
		dec.code.args = append(dec.code.args, 0) // filled later

		// save the current number of instructions
		nbOpcodesStart := len(dec.code.instrs)

		_, err := dec.load(bc[:instrSkip])
		if err != nil {
			return nil, err
		}
		nbOpcodesContext := byte(len(dec.code.instrs) - nbOpcodesStart)
		// update the args slice (see the opcode implementation)
		dec.code.args[instrSkipIndex] = nbOpcodesContext
		dec.code.args[instrSkipIndex+1] = instrSkip - nbOpcodesContext // without op codes

		bc = bc[instrSkip:]

		dec.outLength = 1
		dec.outIndex = 0
		dec.slotRef = 0
		dec.inCtxtItem = false
	}

	return bc, nil
}

// insert TEMP_COPY commands for slots that need them (that change and are referenced later)
func (dec *decoder) applyAnalysis(code []instr) []instr {
	if dec.code.constraint {
		return code
	}

	tempcount := 0
	tempCopy := opcodeTable[ocTEMP_COPY].impl[0]
	for _, c := range dec.contexts[:dec.slotRef] {
		if !c.referenced || !c.changed {
			continue
		}

		code = append(code, instr{})
		tip := code[int(c.codeRef)+tempcount:]
		copy(tip[1:], tip)
		tip[0] = instr{fn: tempCopy, code: ocTEMP_COPY}
		dec.code.delete = true

		tempcount++
	}

	return code
}
