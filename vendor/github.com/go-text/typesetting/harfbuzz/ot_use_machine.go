package harfbuzz

// Code generated with ragel -Z -o ot_use_machine.go ot_use_machine.rl ; sed -i '/^\/\/line/ d' ot_use_machine.go ; goimports -w ot_use_machine.go  DO NOT EDIT.

// ported from harfbuzz/src/hb-ot-shape-complex-use-machine.rl Copyright © 2015 Mozilla Foundation. Google, Inc. Jonathan Kew Behdad Esfahbod

const (
	useViramaTerminatedCluster = iota
	useSakotTerminatedCluster
	useStandardCluster
	useNumberJoinerTerminatedCluster
	useNumeralCluster
	useSymbolCluster
	useHieroglyphCluster
	useBrokenCluster
	useNonCluster
)

const useSyllableMachine_ex_B = 1
const useSyllableMachine_ex_CGJ = 6
const useSyllableMachine_ex_CMAbv = 31
const useSyllableMachine_ex_CMBlw = 32
const useSyllableMachine_ex_CS = 43
const useSyllableMachine_ex_FAbv = 24
const useSyllableMachine_ex_FBlw = 25
const useSyllableMachine_ex_FMAbv = 45
const useSyllableMachine_ex_FMBlw = 46
const useSyllableMachine_ex_FMPst = 47
const useSyllableMachine_ex_FPst = 26
const useSyllableMachine_ex_G = 49
const useSyllableMachine_ex_GB = 5
const useSyllableMachine_ex_H = 12
const useSyllableMachine_ex_HN = 13
const useSyllableMachine_ex_HVM = 44
const useSyllableMachine_ex_J = 50
const useSyllableMachine_ex_MAbv = 27
const useSyllableMachine_ex_MBlw = 28
const useSyllableMachine_ex_MPre = 30
const useSyllableMachine_ex_MPst = 29
const useSyllableMachine_ex_N = 4
const useSyllableMachine_ex_O = 0
const useSyllableMachine_ex_R = 18
const useSyllableMachine_ex_SB = 51
const useSyllableMachine_ex_SE = 52
const useSyllableMachine_ex_SMAbv = 41
const useSyllableMachine_ex_SMBlw = 42
const useSyllableMachine_ex_SUB = 11
const useSyllableMachine_ex_Sk = 48
const useSyllableMachine_ex_VAbv = 33
const useSyllableMachine_ex_VBlw = 34
const useSyllableMachine_ex_VMAbv = 37
const useSyllableMachine_ex_VMBlw = 38
const useSyllableMachine_ex_VMPre = 23
const useSyllableMachine_ex_VMPst = 39
const useSyllableMachine_ex_VPre = 22
const useSyllableMachine_ex_VPst = 35
const useSyllableMachine_ex_ZWNJ = 14

var _useSyllableMachine_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 1, 3,
	1, 4, 1, 5, 1, 6, 1, 7,
	1, 8, 1, 9, 1, 10, 1, 11,
	1, 12, 1, 13,
}

var _useSyllableMachine_key_offsets []int16 = []int16{
	0, 35, 37, 38, 62, 86, 87, 103,
	114, 120, 125, 129, 131, 132, 142, 151,
	159, 160, 167, 182, 196, 209, 227, 244,
	263, 286, 298, 299, 300, 326, 350, 351,
	367, 378, 384, 389, 393, 395, 396, 406,
	415, 423, 424, 431, 446, 460, 473, 491,
	508, 527, 550, 562, 563, 564, 593, 617,
	619, 620, 622, 624, 627,
}

var _useSyllableMachine_trans_keys []byte = []byte{
	0, 1, 4, 5, 11, 12, 13, 18,
	23, 24, 25, 26, 27, 28, 30, 31,
	32, 33, 34, 35, 37, 38, 39, 41,
	42, 43, 44, 45, 46, 47, 48, 49,
	51, 22, 29, 41, 42, 42, 11, 12,
	23, 24, 25, 26, 27, 28, 30, 31,
	32, 33, 34, 35, 37, 38, 39, 44,
	45, 46, 47, 48, 22, 29, 11, 12,
	23, 24, 25, 26, 27, 28, 30, 33,
	34, 35, 37, 38, 39, 44, 45, 46,
	47, 48, 22, 29, 31, 32, 1, 22,
	23, 24, 25, 26, 33, 34, 35, 37,
	38, 39, 44, 45, 46, 47, 48, 23,
	24, 25, 26, 37, 38, 39, 45, 46,
	47, 48, 24, 25, 26, 45, 46, 47,
	25, 26, 45, 46, 47, 26, 45, 46,
	47, 45, 46, 46, 24, 25, 26, 37,
	38, 39, 45, 46, 47, 48, 24, 25,
	26, 38, 39, 45, 46, 47, 48, 24,
	25, 26, 39, 45, 46, 47, 48, 1,
	24, 25, 26, 45, 46, 47, 48, 23,
	24, 25, 26, 33, 34, 35, 37, 38,
	39, 44, 45, 46, 47, 48, 23, 24,
	25, 26, 34, 35, 37, 38, 39, 44,
	45, 46, 47, 48, 23, 24, 25, 26,
	35, 37, 38, 39, 44, 45, 46, 47,
	48, 22, 23, 24, 25, 26, 28, 29,
	33, 34, 35, 37, 38, 39, 44, 45,
	46, 47, 48, 22, 23, 24, 25, 26,
	29, 33, 34, 35, 37, 38, 39, 44,
	45, 46, 47, 48, 23, 24, 25, 26,
	27, 28, 33, 34, 35, 37, 38, 39,
	44, 45, 46, 47, 48, 22, 29, 11,
	12, 23, 24, 25, 26, 27, 28, 30,
	32, 33, 34, 35, 37, 38, 39, 44,
	45, 46, 47, 48, 22, 29, 1, 23,
	24, 25, 26, 37, 38, 39, 45, 46,
	47, 48, 13, 4, 11, 12, 23, 24,
	25, 26, 27, 28, 30, 31, 32, 33,
	34, 35, 37, 38, 39, 41, 42, 44,
	45, 46, 47, 48, 22, 29, 11, 12,
	23, 24, 25, 26, 27, 28, 30, 33,
	34, 35, 37, 38, 39, 44, 45, 46,
	47, 48, 22, 29, 31, 32, 1, 22,
	23, 24, 25, 26, 33, 34, 35, 37,
	38, 39, 44, 45, 46, 47, 48, 23,
	24, 25, 26, 37, 38, 39, 45, 46,
	47, 48, 24, 25, 26, 45, 46, 47,
	25, 26, 45, 46, 47, 26, 45, 46,
	47, 45, 46, 46, 24, 25, 26, 37,
	38, 39, 45, 46, 47, 48, 24, 25,
	26, 38, 39, 45, 46, 47, 48, 24,
	25, 26, 39, 45, 46, 47, 48, 1,
	24, 25, 26, 45, 46, 47, 48, 23,
	24, 25, 26, 33, 34, 35, 37, 38,
	39, 44, 45, 46, 47, 48, 23, 24,
	25, 26, 34, 35, 37, 38, 39, 44,
	45, 46, 47, 48, 23, 24, 25, 26,
	35, 37, 38, 39, 44, 45, 46, 47,
	48, 22, 23, 24, 25, 26, 28, 29,
	33, 34, 35, 37, 38, 39, 44, 45,
	46, 47, 48, 22, 23, 24, 25, 26,
	29, 33, 34, 35, 37, 38, 39, 44,
	45, 46, 47, 48, 23, 24, 25, 26,
	27, 28, 33, 34, 35, 37, 38, 39,
	44, 45, 46, 47, 48, 22, 29, 11,
	12, 23, 24, 25, 26, 27, 28, 30,
	32, 33, 34, 35, 37, 38, 39, 44,
	45, 46, 47, 48, 22, 29, 1, 23,
	24, 25, 26, 37, 38, 39, 45, 46,
	47, 48, 4, 13, 1, 5, 11, 12,
	13, 23, 24, 25, 26, 27, 28, 30,
	31, 32, 33, 34, 35, 37, 38, 39,
	41, 42, 44, 45, 46, 47, 48, 22,
	29, 11, 12, 23, 24, 25, 26, 27,
	28, 30, 31, 32, 33, 34, 35, 37,
	38, 39, 44, 45, 46, 47, 48, 22,
	29, 41, 42, 42, 1, 5, 50, 52,
	49, 50, 52, 49, 51,
}

var _useSyllableMachine_single_lengths []byte = []byte{
	33, 2, 1, 22, 20, 1, 16, 11,
	6, 5, 4, 2, 1, 10, 9, 8,
	1, 7, 15, 14, 13, 18, 17, 17,
	21, 12, 1, 1, 24, 20, 1, 16,
	11, 6, 5, 4, 2, 1, 10, 9,
	8, 1, 7, 15, 14, 13, 18, 17,
	17, 21, 12, 1, 1, 27, 22, 2,
	1, 2, 2, 3, 2,
}

var _useSyllableMachine_range_lengths []byte = []byte{
	1, 0, 0, 1, 2, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 1,
	1, 0, 0, 0, 1, 2, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	1, 1, 0, 0, 0, 1, 1, 0,
	0, 0, 0, 0, 0,
}

var _useSyllableMachine_index_offsets []int16 = []int16{
	0, 35, 38, 40, 64, 87, 89, 106,
	118, 125, 131, 136, 139, 141, 152, 162,
	171, 173, 181, 197, 212, 226, 245, 263,
	282, 305, 318, 320, 322, 348, 371, 373,
	390, 402, 409, 415, 420, 423, 425, 436,
	446, 455, 457, 465, 481, 496, 510, 529,
	547, 566, 589, 602, 604, 606, 635, 659,
	662, 664, 667, 670, 674,
}

var _useSyllableMachine_indicies []byte = []byte{
	0, 1, 3, 4, 5, 6, 7, 8,
	10, 11, 12, 13, 14, 15, 16, 17,
	18, 19, 20, 21, 22, 23, 24, 25,
	26, 27, 28, 29, 30, 31, 6, 32,
	33, 9, 2, 0, 35, 34, 35, 34,
	37, 38, 40, 41, 42, 43, 44, 45,
	46, 1, 47, 48, 49, 50, 51, 52,
	53, 54, 55, 56, 57, 38, 39, 36,
	37, 38, 40, 41, 42, 43, 44, 45,
	46, 48, 49, 50, 51, 52, 53, 54,
	55, 56, 57, 38, 39, 47, 36, 37,
	58, 39, 40, 41, 42, 43, 48, 49,
	50, 51, 52, 53, 40, 55, 56, 57,
	59, 36, 40, 41, 42, 43, 51, 52,
	53, 55, 56, 57, 59, 36, 41, 42,
	43, 55, 56, 57, 36, 42, 43, 55,
	56, 57, 36, 43, 55, 56, 57, 36,
	55, 56, 36, 56, 36, 41, 42, 43,
	51, 52, 53, 55, 56, 57, 59, 36,
	41, 42, 43, 52, 53, 55, 56, 57,
	59, 36, 41, 42, 43, 53, 55, 56,
	57, 59, 36, 61, 60, 41, 42, 43,
	55, 56, 57, 59, 36, 40, 41, 42,
	43, 48, 49, 50, 51, 52, 53, 40,
	55, 56, 57, 59, 36, 40, 41, 42,
	43, 49, 50, 51, 52, 53, 40, 55,
	56, 57, 59, 36, 40, 41, 42, 43,
	50, 51, 52, 53, 40, 55, 56, 57,
	59, 36, 39, 40, 41, 42, 43, 45,
	39, 48, 49, 50, 51, 52, 53, 40,
	55, 56, 57, 59, 36, 39, 40, 41,
	42, 43, 39, 48, 49, 50, 51, 52,
	53, 40, 55, 56, 57, 59, 36, 40,
	41, 42, 43, 44, 45, 48, 49, 50,
	51, 52, 53, 40, 55, 56, 57, 59,
	39, 36, 37, 38, 40, 41, 42, 43,
	44, 45, 46, 47, 48, 49, 50, 51,
	52, 53, 54, 55, 56, 57, 38, 39,
	36, 37, 40, 41, 42, 43, 51, 52,
	53, 55, 56, 57, 59, 58, 63, 62,
	3, 64, 37, 38, 40, 41, 42, 43,
	44, 45, 46, 1, 47, 48, 49, 50,
	51, 52, 53, 0, 35, 54, 55, 56,
	57, 38, 39, 36, 5, 6, 10, 11,
	12, 13, 14, 15, 16, 19, 20, 21,
	22, 23, 24, 28, 29, 30, 31, 6,
	9, 18, 65, 5, 65, 9, 10, 11,
	12, 13, 19, 20, 21, 22, 23, 24,
	10, 29, 30, 31, 66, 65, 10, 11,
	12, 13, 22, 23, 24, 29, 30, 31,
	66, 65, 11, 12, 13, 29, 30, 31,
	65, 12, 13, 29, 30, 31, 65, 13,
	29, 30, 31, 65, 29, 30, 65, 30,
	65, 11, 12, 13, 22, 23, 24, 29,
	30, 31, 66, 65, 11, 12, 13, 23,
	24, 29, 30, 31, 66, 65, 11, 12,
	13, 24, 29, 30, 31, 66, 65, 67,
	65, 11, 12, 13, 29, 30, 31, 66,
	65, 10, 11, 12, 13, 19, 20, 21,
	22, 23, 24, 10, 29, 30, 31, 66,
	65, 10, 11, 12, 13, 20, 21, 22,
	23, 24, 10, 29, 30, 31, 66, 65,
	10, 11, 12, 13, 21, 22, 23, 24,
	10, 29, 30, 31, 66, 65, 9, 10,
	11, 12, 13, 15, 9, 19, 20, 21,
	22, 23, 24, 10, 29, 30, 31, 66,
	65, 9, 10, 11, 12, 13, 9, 19,
	20, 21, 22, 23, 24, 10, 29, 30,
	31, 66, 65, 10, 11, 12, 13, 14,
	15, 19, 20, 21, 22, 23, 24, 10,
	29, 30, 31, 66, 9, 65, 5, 6,
	10, 11, 12, 13, 14, 15, 16, 18,
	19, 20, 21, 22, 23, 24, 28, 29,
	30, 31, 6, 9, 65, 5, 10, 11,
	12, 13, 22, 23, 24, 29, 30, 31,
	66, 65, 68, 65, 7, 65, 1, 1,
	5, 6, 7, 10, 11, 12, 13, 14,
	15, 16, 17, 18, 19, 20, 21, 22,
	23, 24, 25, 26, 28, 29, 30, 31,
	6, 9, 65, 5, 6, 10, 11, 12,
	13, 14, 15, 16, 17, 18, 19, 20,
	21, 22, 23, 24, 28, 29, 30, 31,
	6, 9, 65, 25, 26, 65, 26, 65,
	1, 1, 69, 71, 32, 70, 32, 71,
	71, 70, 32, 33, 70,
}

var _useSyllableMachine_trans_targs []byte = []byte{
	1, 3, 0, 26, 28, 29, 30, 51,
	53, 31, 32, 33, 34, 35, 46, 47,
	48, 54, 49, 43, 44, 45, 38, 39,
	40, 55, 56, 57, 50, 36, 37, 0,
	58, 60, 0, 2, 0, 4, 5, 6,
	7, 8, 9, 10, 21, 22, 23, 24,
	18, 19, 20, 13, 14, 15, 25, 11,
	12, 0, 0, 16, 0, 17, 0, 27,
	0, 0, 41, 42, 52, 0, 0, 59,
}

var _useSyllableMachine_trans_actions []byte = []byte{
	0, 0, 9, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 7,
	0, 0, 21, 0, 15, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 5, 11, 0, 13, 0, 19, 0,
	17, 25, 0, 0, 0, 27, 23, 0,
}

var _useSyllableMachine_to_state_actions []byte = []byte{
	1, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

var _useSyllableMachine_from_state_actions []byte = []byte{
	3, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

var _useSyllableMachine_eof_trans []int16 = []int16{
	0, 35, 35, 37, 37, 59, 37, 37,
	37, 37, 37, 37, 37, 37, 37, 37,
	61, 37, 37, 37, 37, 37, 37, 37,
	37, 59, 63, 65, 37, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66,
	66, 66, 66, 66, 66, 66, 66, 66,
	66, 70, 71, 71, 71,
}

const useSyllableMachine_start int = 0
const useSyllableMachine_first_final int = 0
const useSyllableMachine_error int = -1

const useSyllableMachine_en_main int = 0

func findSyllablesUse(buffer *Buffer) {
	info := buffer.Info
	data := preprocessInfoUSE(info)
	p, pe := 0, len(data)
	eof := pe
	var cs, act, ts, te int

	{
		cs = useSyllableMachine_start
		ts = 0
		te = 0
		act = 0
	}

	var syllableSerial uint8 = 1

	{
		var _klen int
		var _trans int
		var _acts int
		var _nacts uint
		var _keys int
		if p == pe {
			goto _test_eof
		}
	_resume:
		_acts = int(_useSyllableMachine_from_state_actions[cs])
		_nacts = uint(_useSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _useSyllableMachine_actions[_acts-1] {
			case 1:
				ts = p

			}
		}

		_keys = int(_useSyllableMachine_key_offsets[cs])
		_trans = int(_useSyllableMachine_index_offsets[cs])

		_klen = int(_useSyllableMachine_single_lengths[cs])
		if _klen > 0 {
			_lower := int(_keys)
			var _mid int
			_upper := int(_keys + _klen - 1)
			for {
				if _upper < _lower {
					break
				}

				_mid = _lower + ((_upper - _lower) >> 1)
				switch {
				case ((data[p]).p.v.complexCategory) < _useSyllableMachine_trans_keys[_mid]:
					_upper = _mid - 1
				case ((data[p]).p.v.complexCategory) > _useSyllableMachine_trans_keys[_mid]:
					_lower = _mid + 1
				default:
					_trans += int(_mid - int(_keys))
					goto _match
				}
			}
			_keys += _klen
			_trans += _klen
		}

		_klen = int(_useSyllableMachine_range_lengths[cs])
		if _klen > 0 {
			_lower := int(_keys)
			var _mid int
			_upper := int(_keys + (_klen << 1) - 2)
			for {
				if _upper < _lower {
					break
				}

				_mid = _lower + (((_upper - _lower) >> 1) & ^1)
				switch {
				case ((data[p]).p.v.complexCategory) < _useSyllableMachine_trans_keys[_mid]:
					_upper = _mid - 2
				case ((data[p]).p.v.complexCategory) > _useSyllableMachine_trans_keys[_mid+1]:
					_lower = _mid + 2
				default:
					_trans += int((_mid - int(_keys)) >> 1)
					goto _match
				}
			}
			_trans += _klen
		}

	_match:
		_trans = int(_useSyllableMachine_indicies[_trans])
	_eof_trans:
		cs = int(_useSyllableMachine_trans_targs[_trans])

		if _useSyllableMachine_trans_actions[_trans] == 0 {
			goto _again
		}

		_acts = int(_useSyllableMachine_trans_actions[_trans])
		_nacts = uint(_useSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _useSyllableMachine_actions[_acts-1] {
			case 2:
				te = p + 1
				{
					foundSyllableUSE(useStandardCluster, data, ts, te, info, &syllableSerial)
				}
			case 3:
				te = p + 1
				{
					foundSyllableUSE(useBrokenCluster, data, ts, te, info, &syllableSerial)
				}
			case 4:
				te = p + 1
				{
					foundSyllableUSE(useNonCluster, data, ts, te, info, &syllableSerial)
				}
			case 5:
				te = p
				p--
				{
					foundSyllableUSE(useViramaTerminatedCluster, data, ts, te, info, &syllableSerial)
				}
			case 6:
				te = p
				p--
				{
					foundSyllableUSE(useSakotTerminatedCluster, data, ts, te, info, &syllableSerial)
				}
			case 7:
				te = p
				p--
				{
					foundSyllableUSE(useStandardCluster, data, ts, te, info, &syllableSerial)
				}
			case 8:
				te = p
				p--
				{
					foundSyllableUSE(useNumberJoinerTerminatedCluster, data, ts, te, info, &syllableSerial)
				}
			case 9:
				te = p
				p--
				{
					foundSyllableUSE(useNumeralCluster, data, ts, te, info, &syllableSerial)
				}
			case 10:
				te = p
				p--
				{
					foundSyllableUSE(useSymbolCluster, data, ts, te, info, &syllableSerial)
				}
			case 11:
				te = p
				p--
				{
					foundSyllableUSE(useHieroglyphCluster, data, ts, te, info, &syllableSerial)
				}
			case 12:
				te = p
				p--
				{
					foundSyllableUSE(useBrokenCluster, data, ts, te, info, &syllableSerial)
				}
			case 13:
				te = p
				p--
				{
					foundSyllableUSE(useNonCluster, data, ts, te, info, &syllableSerial)
				}
			}
		}

	_again:
		_acts = int(_useSyllableMachine_to_state_actions[cs])
		_nacts = uint(_useSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _useSyllableMachine_actions[_acts-1] {
			case 0:
				ts = 0

			}
		}

		p++
		if p != pe {
			goto _resume
		}
	_test_eof:
		{
		}
		if p == eof {
			if _useSyllableMachine_eof_trans[cs] > 0 {
				_trans = int(_useSyllableMachine_eof_trans[cs] - 1)
				goto _eof_trans
			}
		}

	}

	_ = act // needed by Ragel, but unused
}
