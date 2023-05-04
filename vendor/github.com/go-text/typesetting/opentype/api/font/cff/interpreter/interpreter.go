// SPDX-License-Identifier: Unlicense OR BSD-3-Clause

// Package psinterpreter implement a Postscript interpreter
// required to parse .CFF files, and Type1 and Type2 Charstrings.
// This package provides the low-level mechanisms needed to
// read such formats; the data is consumed in higher level packages,
// which implement `PsOperatorHandler`.
// It also provides helpers to interpret glyph outline descriptions,
// shared between Type1 and CFF font formats.
package psinterpreter

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"strconv"
)

var (

	// ErrInterrupt signals the interpreter to stop early, without erroring.
	ErrInterrupt = errors.New("interruption")

	errInvalidCFFTable               = errors.New("invalid ps instructions")
	errUnsupportedRealNumberEncoding = errors.New("unsupported real number encoding")

	be = binary.BigEndian
)

const (
	// psArgStackSize is the argument stack size for a PostScript interpreter.
	// 5176.CFF.pdf section 4 "DICT Data" says that "An operator may be
	// preceded by up to a maximum of 48 operands". 5177.Type2.pdf Appendix B
	// "Type 2 Charstring Implementation Limits" says that "Argument stack 48".
	// T1_SPEC.pdf 6.1 Encoding as a limitation of 24.
	psArgStackSize = 48

	// Similarly, Appendix B says "Subr nesting, stack limit 10".
	psCallStackSize = 10

	maxRealNumberStrLen = 64 // Maximum length in bytes of the "-123.456E-7" representation.
)

// Context is the flavour of the Postcript language.
type Context uint32

const (
	TopDict         Context = iota // Top dict in CFF files
	PrivateDict                    // Private dict in CFF files
	Type2Charstring                // Charstring in CFF files
	Type1Charstring                // Charstring in Type1 font files
)

type ArgStack struct {
	Vals [psArgStackSize]int32
	// Effecive size currently in use. The first value to
	// pop is at index Top-1
	Top int32
}

// Uint16 returns the top level value as uint16,
// without popping the stack.
func (a *ArgStack) Uint16() uint16 { return uint16(a.Vals[a.Top-1]) }

// Float return the top level value as a real number (which is stored as its binary representation),
// without popping the stack.
func (a *ArgStack) Float() float32 {
	return math.Float32frombits(uint32(a.Vals[a.Top-1]))
}

// Pop returns the top level value and decrease `Top`
// It will panic if the stack is empty.
func (a *ArgStack) Pop() int32 {
	a.Top--
	return a.Vals[a.Top]
}

// Clear clears the stack
func (a *ArgStack) Clear() { a.Top = 0 }

// PopN check and remove the n top levels entries.
// Passing a negative `numPop` clears all the stack.
func (a *ArgStack) PopN(numPop int32) error {
	if a.Top < numPop {
		return fmt.Errorf("invalid number of operands in PS stack: %d", numPop)
	}
	if numPop < 0 { // pop all
		a.Top = 0
	} else {
		a.Top -= numPop
	}
	return nil
}

// Machine is a PostScript interpreter.
// A same interpreter may be re-used using muliples `Run` calls.
type Machine struct {
	localSubrs  [][]byte
	globalSubrs [][]byte

	instructions []byte

	callStack struct {
		vals [psCallStackSize][]byte // parent instructions
		top  int32                   // effecive size currently in use
	}
	ArgStack ArgStack

	parseNumberBuf [maxRealNumberStrLen]byte
	ctx            Context
}

// SkipBytes skips the next `count` bytes from the instructions, and clears the argument stack.
// It does nothing if `count` exceed the length of the instructions.
func (p *Machine) SkipBytes(count int32) {
	if int(count) >= len(p.instructions) {
		return
	}
	p.instructions = p.instructions[count:]
	p.ArgStack.Clear()
}

// 5176.CFF.pdf section 4 "DICT Data" says that "Two-byte operators have an
// initial escape byte of 12".
const escapeByte = 12

// Run runs the instructions in the PostScript context asked by `handler`.
// `localSubrs` and `globalSubrs` contains the subroutines that may be called in the instructions.
func (p *Machine) Run(instructions []byte, localSubrs, globalSubrs [][]byte, handler OperatorHandler) error {
	p.ctx = handler.Context()
	p.instructions = instructions
	p.localSubrs = localSubrs
	p.globalSubrs = globalSubrs
	p.ArgStack.Top = 0
	p.callStack.top = 0

	for len(p.instructions) > 0 {
		// Push a numeric operand on the stack, if applicable.
		if hasResult, err := p.parseNumber(); hasResult {
			if err != nil {
				return err
			}
			continue
		}

		// Otherwise, execute an operator.
		b := p.instructions[0]
		p.instructions = p.instructions[1:]

		// check for the escape byte
		escaped := b == escapeByte
		if escaped {
			if len(p.instructions) <= 0 {
				return errInvalidCFFTable
			}
			b = p.instructions[0]
			p.instructions = p.instructions[1:]
		}

		err := handler.Apply(p, Operator{Operator: b, IsEscaped: escaped})
		if err == ErrInterrupt { // stop cleanly
			return nil
		}
		if err != nil {
			return err
		}

	}
	return nil
}

// See 5176.CFF.pdf section 4 "DICT Data".
func (p *Machine) parseNumber() (hasResult bool, err error) {
	number := int32(0)
	switch b := p.instructions[0]; {
	case b == 28:
		if len(p.instructions) < 3 {
			return true, errInvalidCFFTable
		}
		number, hasResult = int32(int16(be.Uint16(p.instructions[1:]))), true
		p.instructions = p.instructions[3:]

	case b == 29 && p.ctx != Type2Charstring:
		if len(p.instructions) < 5 {
			return true, errInvalidCFFTable
		}
		number, hasResult = int32(be.Uint32(p.instructions[1:])), true
		p.instructions = p.instructions[5:]

	case b == 30 && p.ctx != Type2Charstring && p.ctx != Type1Charstring:
		// Parse a real number. This isn't listed in 5176.CFF.pdf Table 3
		// "Operand Encoding" but that table lists integer encodings. Further
		// down the page it says "A real number operand is provided in addition
		// to integer operands. This operand begins with a byte value of 30
		// followed by a variable-length sequence of bytes."

		s := p.parseNumberBuf[:0]
		p.instructions = p.instructions[1:]
	loop:
		for {
			if len(p.instructions) == 0 {
				return true, errInvalidCFFTable
			}
			by := p.instructions[0]
			p.instructions = p.instructions[1:]
			// Process by's two nibbles, high then low.
			for i := 0; i < 2; i++ {
				nib := by >> 4
				by = by << 4
				if nib == 0x0f {
					f, err := strconv.ParseFloat(string(s), 32)
					if err != nil {
						return true, errInvalidCFFTable
					}
					number, hasResult = int32(math.Float32bits(float32(f))), true
					break loop
				}
				if nib == 0x0d {
					return true, errInvalidCFFTable
				}
				if len(s)+maxNibbleDefsLength > len(p.parseNumberBuf) {
					return true, errUnsupportedRealNumberEncoding
				}
				s = append(s, nibbleDefs[nib]...)
			}
		}

	case b < 32:
		// not a number: no-op.
	case b < 247:
		p.instructions = p.instructions[1:]
		number, hasResult = int32(b)-139, true
	case b < 251:
		if len(p.instructions) < 2 {
			return true, errInvalidCFFTable
		}
		b1 := p.instructions[1]
		p.instructions = p.instructions[2:]
		number, hasResult = +int32(b-247)*256+int32(b1)+108, true
	case b < 255:
		if len(p.instructions) < 2 {
			return true, errInvalidCFFTable
		}
		b1 := p.instructions[1]
		p.instructions = p.instructions[2:]
		number, hasResult = -int32(b-251)*256-int32(b1)-108, true
	case b == 255 && (p.ctx == Type2Charstring || p.ctx == Type1Charstring):
		if len(p.instructions) < 5 {
			return true, errInvalidCFFTable
		}
		number, hasResult = int32(be.Uint32(p.instructions[1:])), true
		p.instructions = p.instructions[5:]
	}

	if hasResult {
		if p.ArgStack.Top == psArgStackSize {
			return true, errInvalidCFFTable
		}
		p.ArgStack.Vals[p.ArgStack.Top] = number
		p.ArgStack.Top++
	}
	return hasResult, nil
}

const maxNibbleDefsLength = len("E-")

// nibbleDefs encodes 5176.CFF.pdf Table 5 "Nibble Definitions".
var nibbleDefs = [16]string{
	0x00: "0",
	0x01: "1",
	0x02: "2",
	0x03: "3",
	0x04: "4",
	0x05: "5",
	0x06: "6",
	0x07: "7",
	0x08: "8",
	0x09: "9",
	0x0a: ".",
	0x0b: "E",
	0x0c: "E-",
	0x0d: "",
	0x0e: "-",
	0x0f: "",
}

// subrBias returns the subroutine index bias as per 5177.Type2.pdf section 4.7
// "Subroutine Operators".
func subrBias(numSubroutines int) int32 {
	if numSubroutines < 1240 {
		return 107
	}
	if numSubroutines < 33900 {
		return 1131
	}
	return 32768
}

// CallSubroutine calls the subroutine, identified by its index, as found
// in the instructions (that is, before applying the subroutine biased).
// `isLocal` controls whether the local or global subroutines are used.
// No argument stack modification is performed.
func (p *Machine) CallSubroutine(index int32, isLocal bool) error {
	subrs := p.globalSubrs
	if isLocal {
		subrs = p.localSubrs
	}

	// no bias in type1 fonts
	if p.ctx == Type2Charstring {
		index += subrBias(len(subrs))
	}

	if index < 0 || int(index) >= len(subrs) {
		return fmt.Errorf("invalid subroutine index %d (for length %d)", index, len(subrs))
	}
	if p.callStack.top == psCallStackSize {
		return errors.New("maximum call stack size reached")
	}
	// save the current instructions
	p.callStack.vals[p.callStack.top] = p.instructions
	p.callStack.top++

	// activate the subroutine
	p.instructions = subrs[index]
	return nil
}

// Return returns from a subroutine call.
func (p *Machine) Return() error {
	if p.callStack.top <= 0 {
		return errors.New("no subroutine has been called")
	}
	p.callStack.top--
	// restore the previous instructions
	p.instructions = p.callStack.vals[p.callStack.top]
	return nil
}

// Operator is a postcript command, which may be escaped.
type Operator struct {
	Operator  byte
	IsEscaped bool
}

func (p Operator) String() string {
	if p.IsEscaped {
		return fmt.Sprintf("2-byte operator (12 %d)", p.Operator)
	}
	return fmt.Sprintf("1-byte operator (%d)", p.Operator)
}

// OperatorHandler defines the behaviour of an operator.
type OperatorHandler interface {
	// Context defines the precise behaviour of the interpreter,
	// which has small nuances depending on the context.
	Context() Context

	// Apply implements the operator defined by `operator` (which is the second byte if `escaped` is true).
	//
	// Returning `ErrInterrupt` stop the parsing of the instructions, without reporting an error.
	// It can be used as an optimization.
	Apply(state *Machine, operator Operator) error
}
