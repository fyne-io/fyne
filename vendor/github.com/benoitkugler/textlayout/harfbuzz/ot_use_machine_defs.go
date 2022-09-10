package harfbuzz

// logic needed by the USE rl parser

func notCCSDefaultIgnorable(i GlyphInfo) bool {
	return !(i.complexCategory == useSyllableMachine_ex_CGJ && i.isDefaultIgnorable())
}

type pairUSE struct {
	i int // index in the original info slice
	v GlyphInfo
}

type machineIndexUSE struct {
	j int // index in the filtered slice
	p pairUSE
}

func preprocessInfoUSE(info []GlyphInfo) []machineIndexUSE {
	filterMark := func(p pairUSE) bool {
		if p.v.complexCategory == useSyllableMachine_ex_ZWNJ {
			for i := p.i + 1; i < len(info); i++ {
				if notCCSDefaultIgnorable(info[i]) {
					return !info[i].isUnicodeMark()
				}
			}
		}
		return true
	}
	var tmp []pairUSE
	for i, v := range info {
		if notCCSDefaultIgnorable(v) {
			p := pairUSE{i, v}
			if filterMark(p) {
				tmp = append(tmp, p)
			}
		}
	}
	data := make([]machineIndexUSE, len(tmp))
	for j, p := range tmp {
		data[j] = machineIndexUSE{j: j, p: p}
	}
	return data
}

func foundSyllableUSE(syllableType uint8, data []machineIndexUSE, ts, te int, info []GlyphInfo, syllableSerial *uint8) {
	start := data[ts].p.i
	end := len(info) // te might right after the end of data
	if te < len(data) {
		end = data[te].p.i
	}
	for i := start; i < end; i++ {
		info[i].syllable = (*syllableSerial << 4) | syllableType
	}
	*syllableSerial++
	if *syllableSerial == 16 {
		*syllableSerial = 1
	}
}
