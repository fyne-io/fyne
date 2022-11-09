package graphite

import (
	"fmt"
)

type machineStatus uint8

const (
	machineStackUnderflow machineStatus = iota
	machineStackNotEmpty
	machineStackOverflow
	machineSlotOffsetOutOfBounds
	machineDiedEarly
)

func (m machineStatus) Error() string {
	switch m {
	case machineStackUnderflow:
		return "stack_underflow"
	case machineStackNotEmpty:
		return "stack_not_empty"
	case machineStackOverflow:
		return "stack_overflow"
	case machineSlotOffsetOutOfBounds:
		return "slot_offset_out_bounds"
	case machineDiedEarly:
		return "died_early"
	}
	return "<unknown machine status>"
}

const (
	stackGuard = 2
	stackMax   = 1 << 10
)

type stack struct {
	vals [stackMax + 2*stackGuard]int32
	top  int // the top of the stack is at vals[top-1]
}

type machine struct {
	map_  *slotMap // shared with the finite state machine
	stack stack
}

func newMachine(map_ *slotMap) *machine { return &machine{map_: map_} }

// map_ is the index into m.slotMap.slots to start from
// return the result of the code and the index of the current slot
func (m *machine) run(co *code, map_ int) (int32, int, error) {
	if L := co.maxRef + int(m.map_.preContext); m.map_.size <= L || m.map_.get(L) == nil {
		return 1, map_, machineSlotOffsetOutOfBounds
	}

	// Declare virtual machine registers
	reg := regbank{
		smap:      m.map_,
		is:        m.map_.slots[map_],
		map_:      map_,
		mapb:      1 + int(m.map_.preContext),
		ip:        0,
		direction: m.map_.isRTL,
		flags:     0,
	}

	// Run the program
	program, args := co.instrs, co.args
	for ; reg.ip < len(program); reg.ip++ {
		var ok bool
		instruction := program[reg.ip]
		args, ok = instruction.fn(&reg, &m.stack, args)

		if debugMode >= 3 {
			fmt.Printf("FSM after index %d (opcode %s) : %v; args: %v\n", reg.ip, instruction.code, m.stack.vals[:m.stack.top], args)
		}

		// false is returned either for a return opcode
		// or for an error
		if !ok {
			if instruction.code.isReturn() {
				break
			} else {
				return 0, map_, machineDiedEarly
			}
		}

	}

	if debugMode >= 3 {
		fmt.Println()
	}

	if err := m.checkFinalStack(); err != nil {
		return 0, map_, err
	}

	var ret int32
	if m.stack.top > 0 {
		ret = m.stack.pop()
	}

	reg.smap.slots[reg.map_] = reg.is
	return ret, reg.map_, nil
}

func (m *machine) checkFinalStack() error {
	if m.stack.top < 1 {
		return machineStackUnderflow // This should be impossible now.
	} else if m.stack.top >= stackMax {
		return machineStackOverflow // So should this.
	} else if m.stack.top != 1 {
		return machineStackNotEmpty
	}
	return nil
}

type finiteStateMachine struct {
	ruleTable []rule // from the font file

	// indexes in ruleTable. storage allocated once
	// but cleared and filled in runFSM
	rules []uint16

	slots slotMap // shared with the machine
}

// merge `rules` into `fsm.rules`, according to the order
// defined in `compareRuleIndex`, without adding duplicates
func (fsm *finiteStateMachine) accumulateRules(rules []uint16) {
	a, b := fsm.rules, rules
	out := make([]uint16, 0, len(a)+len(b))

	var i, j int
	for i < len(a) && j < len(b) {
		if compareRuleIndex(fsm.ruleTable, a[i], b[j]) {
			out = append(out, a[i])
			i++
		} else if compareRuleIndex(fsm.ruleTable, b[j], a[i]) {
			out = append(out, b[j])
			j++
		} else { // do not create duplicates
			out = append(out, a[i])
			i++
			j++
		}
	}

	out = append(out, a[i:]...) // one of the tails
	out = append(out, b[j:]...) // is actually empty

	fsm.rules = out
}

// returns true if r1 < r2
func compareRuleIndex(table []rule, r1, r2 uint16) bool {
	lsort := table[r1].sortKey
	rsort := table[r2].sortKey
	return lsort > rsort || (lsort == rsort && r1 < r2)
}

// clears the active rules and slots and install the rule table of the pass
// adjust `slot`
func (fsm *finiteStateMachine) reset(slot *Slot, maxPreCtxt uint16, rules []rule) *Slot {
	fsm.rules = fsm.rules[:0]
	fsm.ruleTable = rules
	var ctxt uint16
	for ; ctxt != maxPreCtxt && slot.prev != nil; ctxt, slot = ctxt+1, slot.prev {
	}
	fsm.slots.reset(slot, ctxt)
	return slot
}
