package harfbuzz

// Code generated with ragel -Z -o ot_myanmar_machine.go ot_myanmar_machine.rl ; sed -i '/^\/\/line/ d' ot_myanmar_machine.go ; goimports -w ot_myanmar_machine.go  DO NOT EDIT.

// ported from harfbuzz/src/hb-ot-shape-complex-myanmar-machine.rl Copyright © 2015 Mozilla Foundation. Google, Inc. Behdad Esfahbod

// myanmar_syllable_type_t
const (
	myanmarConsonantSyllable = iota
	myanmarPunctuationCluster
	myanmarBrokenCluster
	myanmarNonMyanmarCluster
)

const myanmarSyllableMachine_ex_A = 10
const myanmarSyllableMachine_ex_As = 18
const myanmarSyllableMachine_ex_C = 1
const myanmarSyllableMachine_ex_CS = 19
const myanmarSyllableMachine_ex_D = 32
const myanmarSyllableMachine_ex_D0 = 20
const myanmarSyllableMachine_ex_DB = 3
const myanmarSyllableMachine_ex_GB = 11
const myanmarSyllableMachine_ex_H = 4
const myanmarSyllableMachine_ex_IV = 2
const myanmarSyllableMachine_ex_MH = 21
const myanmarSyllableMachine_ex_ML = 33
const myanmarSyllableMachine_ex_MR = 22
const myanmarSyllableMachine_ex_MW = 23
const myanmarSyllableMachine_ex_MY = 24
const myanmarSyllableMachine_ex_P = 31
const myanmarSyllableMachine_ex_PT = 25
const myanmarSyllableMachine_ex_Ra = 16
const myanmarSyllableMachine_ex_V = 8
const myanmarSyllableMachine_ex_VAbv = 26
const myanmarSyllableMachine_ex_VBlw = 27
const myanmarSyllableMachine_ex_VPre = 28
const myanmarSyllableMachine_ex_VPst = 29
const myanmarSyllableMachine_ex_VS = 30
const myanmarSyllableMachine_ex_ZWJ = 6
const myanmarSyllableMachine_ex_ZWNJ = 5

var _myanmarSyllableMachine_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 1, 3,
	1, 4, 1, 5, 1, 6, 1, 7,
	1, 8, 1, 9,
}

var _myanmarSyllableMachine_key_offsets []int16 = []int16{
	0, 25, 43, 49, 52, 57, 64, 69,
	73, 84, 91, 100, 108, 118, 121, 137,
	149, 159, 168, 176, 187, 198, 211, 224,
	239, 253, 270, 276, 279, 284, 291, 296,
	300, 311, 318, 327, 335, 345, 348, 366,
	382, 394, 404, 413, 421, 432, 443, 456,
	469, 484, 498, 515, 533, 550, 573, 578,
}

var _myanmarSyllableMachine_trans_keys []byte = []byte{
	3, 4, 8, 10, 11, 16, 18, 19,
	21, 22, 23, 24, 25, 26, 27, 28,
	29, 30, 31, 32, 33, 1, 2, 5,
	6, 3, 4, 8, 10, 18, 21, 22,
	23, 24, 25, 26, 27, 28, 29, 30,
	33, 5, 6, 8, 18, 25, 29, 5,
	6, 8, 5, 6, 8, 25, 29, 5,
	6, 3, 8, 10, 18, 25, 5, 6,
	8, 18, 25, 5, 6, 8, 25, 5,
	6, 3, 8, 10, 18, 21, 25, 26,
	29, 33, 5, 6, 3, 8, 10, 25,
	29, 5, 6, 3, 8, 10, 18, 25,
	26, 29, 5, 6, 3, 8, 10, 25,
	26, 29, 5, 6, 3, 8, 10, 18,
	25, 26, 29, 33, 5, 6, 16, 1,
	2, 3, 8, 10, 18, 21, 22, 23,
	24, 25, 26, 27, 28, 29, 33, 5,
	6, 3, 8, 10, 18, 25, 26, 27,
	28, 29, 33, 5, 6, 3, 8, 10,
	25, 26, 27, 28, 29, 5, 6, 3,
	8, 10, 25, 26, 27, 29, 5, 6,
	3, 8, 10, 25, 27, 29, 5, 6,
	3, 8, 10, 25, 26, 27, 28, 29,
	30, 5, 6, 3, 8, 10, 18, 25,
	26, 27, 28, 29, 5, 6, 3, 8,
	10, 21, 23, 25, 26, 27, 28, 29,
	33, 5, 6, 3, 8, 10, 18, 21,
	25, 26, 27, 28, 29, 33, 5, 6,
	3, 8, 10, 18, 21, 22, 23, 25,
	26, 27, 28, 29, 33, 5, 6, 3,
	8, 10, 21, 22, 23, 25, 26, 27,
	28, 29, 33, 5, 6, 3, 4, 8,
	10, 18, 21, 22, 23, 24, 25, 26,
	27, 28, 29, 33, 5, 6, 8, 18,
	25, 29, 5, 6, 8, 5, 6, 8,
	25, 29, 5, 6, 3, 8, 10, 18,
	25, 5, 6, 8, 18, 25, 5, 6,
	8, 25, 5, 6, 3, 8, 10, 18,
	21, 25, 26, 29, 33, 5, 6, 3,
	8, 10, 25, 29, 5, 6, 3, 8,
	10, 18, 25, 26, 29, 5, 6, 3,
	8, 10, 25, 26, 29, 5, 6, 3,
	8, 10, 18, 25, 26, 29, 33, 5,
	6, 16, 1, 2, 3, 4, 8, 10,
	18, 21, 22, 23, 24, 25, 26, 27,
	28, 29, 30, 33, 5, 6, 3, 8,
	10, 18, 21, 22, 23, 24, 25, 26,
	27, 28, 29, 33, 5, 6, 3, 8,
	10, 18, 25, 26, 27, 28, 29, 33,
	5, 6, 3, 8, 10, 25, 26, 27,
	28, 29, 5, 6, 3, 8, 10, 25,
	26, 27, 29, 5, 6, 3, 8, 10,
	25, 27, 29, 5, 6, 3, 8, 10,
	25, 26, 27, 28, 29, 30, 5, 6,
	3, 8, 10, 18, 25, 26, 27, 28,
	29, 5, 6, 3, 8, 10, 21, 23,
	25, 26, 27, 28, 29, 33, 5, 6,
	3, 8, 10, 18, 21, 25, 26, 27,
	28, 29, 33, 5, 6, 3, 8, 10,
	18, 21, 22, 23, 25, 26, 27, 28,
	29, 33, 5, 6, 3, 8, 10, 21,
	22, 23, 25, 26, 27, 28, 29, 33,
	5, 6, 3, 4, 8, 10, 18, 21,
	22, 23, 24, 25, 26, 27, 28, 29,
	33, 5, 6, 3, 4, 8, 10, 18,
	21, 22, 23, 24, 25, 26, 27, 28,
	29, 30, 33, 5, 6, 3, 4, 8,
	10, 18, 21, 22, 23, 24, 25, 26,
	27, 28, 29, 33, 5, 6, 3, 4,
	8, 10, 11, 16, 18, 21, 22, 23,
	24, 25, 26, 27, 28, 29, 30, 32,
	33, 1, 2, 5, 6, 11, 16, 32,
	1, 2, 8,
}

var _myanmarSyllableMachine_single_lengths []byte = []byte{
	21, 16, 4, 1, 3, 5, 3, 2,
	9, 5, 7, 6, 8, 1, 14, 10,
	8, 7, 6, 9, 9, 11, 11, 13,
	12, 15, 4, 1, 3, 5, 3, 2,
	9, 5, 7, 6, 8, 1, 16, 14,
	10, 8, 7, 6, 9, 9, 11, 11,
	13, 12, 15, 16, 15, 19, 3, 1,
}

var _myanmarSyllableMachine_range_lengths []byte = []byte{
	2, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 2, 1, 0,
}

var _myanmarSyllableMachine_index_offsets []int16 = []int16{
	0, 24, 42, 48, 51, 56, 63, 68,
	72, 83, 90, 99, 107, 117, 120, 136,
	148, 158, 167, 175, 186, 197, 210, 223,
	238, 252, 269, 275, 278, 283, 290, 295,
	299, 310, 317, 326, 334, 344, 347, 365,
	381, 393, 403, 412, 420, 431, 442, 455,
	468, 483, 497, 514, 532, 549, 571, 576,
}

var _myanmarSyllableMachine_indicies []byte = []byte{
	2, 3, 5, 6, 1, 7, 8, 9,
	10, 11, 12, 13, 14, 15, 16, 17,
	18, 19, 20, 1, 21, 1, 4, 0,
	23, 24, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39,
	25, 22, 26, 40, 33, 37, 25, 22,
	26, 25, 22, 26, 33, 37, 25, 22,
	41, 26, 33, 42, 33, 25, 22, 26,
	42, 33, 25, 22, 26, 33, 25, 22,
	23, 26, 27, 43, 44, 33, 45, 37,
	43, 25, 22, 23, 26, 27, 33, 37,
	25, 22, 23, 26, 27, 43, 33, 45,
	37, 25, 22, 23, 26, 27, 33, 45,
	37, 25, 22, 23, 26, 27, 43, 33,
	45, 37, 43, 25, 22, 1, 1, 22,
	23, 26, 27, 28, 29, 30, 31, 32,
	33, 34, 35, 36, 37, 39, 25, 22,
	23, 26, 27, 46, 33, 34, 35, 36,
	37, 39, 25, 22, 23, 26, 27, 33,
	34, 35, 36, 37, 25, 22, 23, 26,
	27, 33, 34, 35, 37, 25, 22, 23,
	26, 27, 33, 35, 37, 25, 22, 23,
	26, 27, 33, 34, 35, 36, 37, 46,
	25, 22, 23, 26, 27, 46, 33, 34,
	35, 36, 37, 25, 22, 23, 26, 27,
	29, 31, 33, 34, 35, 36, 37, 39,
	25, 22, 23, 26, 27, 46, 29, 33,
	34, 35, 36, 37, 39, 25, 22, 23,
	26, 27, 47, 29, 30, 31, 33, 34,
	35, 36, 37, 39, 25, 22, 23, 26,
	27, 29, 30, 31, 33, 34, 35, 36,
	37, 39, 25, 22, 23, 24, 26, 27,
	28, 29, 30, 31, 32, 33, 34, 35,
	36, 37, 39, 25, 22, 5, 50, 14,
	18, 49, 48, 5, 49, 48, 5, 14,
	18, 49, 48, 51, 5, 14, 52, 14,
	49, 48, 5, 52, 14, 49, 48, 5,
	14, 49, 48, 2, 5, 6, 53, 54,
	14, 55, 18, 53, 49, 48, 2, 5,
	6, 14, 18, 49, 48, 2, 5, 6,
	53, 14, 55, 18, 49, 48, 2, 5,
	6, 14, 55, 18, 49, 48, 2, 5,
	6, 53, 14, 55, 18, 53, 49, 48,
	56, 56, 48, 2, 3, 5, 6, 8,
	10, 11, 12, 13, 14, 15, 16, 17,
	18, 19, 21, 49, 48, 2, 5, 6,
	8, 10, 11, 12, 13, 14, 15, 16,
	17, 18, 21, 49, 48, 2, 5, 6,
	57, 14, 15, 16, 17, 18, 21, 49,
	48, 2, 5, 6, 14, 15, 16, 17,
	18, 49, 48, 2, 5, 6, 14, 15,
	16, 18, 49, 48, 2, 5, 6, 14,
	16, 18, 49, 48, 2, 5, 6, 14,
	15, 16, 17, 18, 57, 49, 48, 2,
	5, 6, 57, 14, 15, 16, 17, 18,
	49, 48, 2, 5, 6, 10, 12, 14,
	15, 16, 17, 18, 21, 49, 48, 2,
	5, 6, 57, 10, 14, 15, 16, 17,
	18, 21, 49, 48, 2, 5, 6, 58,
	10, 11, 12, 14, 15, 16, 17, 18,
	21, 49, 48, 2, 5, 6, 10, 11,
	12, 14, 15, 16, 17, 18, 21, 49,
	48, 2, 3, 5, 6, 8, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 21,
	49, 48, 23, 24, 26, 27, 59, 29,
	30, 31, 32, 33, 34, 35, 36, 37,
	38, 39, 25, 22, 23, 60, 26, 27,
	28, 29, 30, 31, 32, 33, 34, 35,
	36, 37, 39, 25, 22, 2, 3, 5,
	6, 1, 1, 8, 10, 11, 12, 13,
	14, 15, 16, 17, 18, 19, 1, 21,
	1, 49, 48, 1, 1, 1, 1, 61,
	62, 61,
}

var _myanmarSyllableMachine_trans_targs []byte = []byte{
	0, 1, 26, 37, 0, 27, 33, 51,
	39, 54, 40, 46, 47, 48, 29, 42,
	43, 44, 32, 50, 55, 45, 0, 2,
	13, 0, 3, 9, 14, 15, 21, 22,
	23, 5, 17, 18, 19, 8, 25, 20,
	4, 6, 7, 10, 12, 11, 16, 24,
	0, 0, 28, 30, 31, 34, 36, 35,
	38, 41, 49, 52, 53, 0, 0,
}

var _myanmarSyllableMachine_trans_actions []byte = []byte{
	13, 0, 0, 0, 7, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 15, 0,
	0, 5, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	17, 11, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 19, 9,
}

var _myanmarSyllableMachine_to_state_actions []byte = []byte{
	1, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var _myanmarSyllableMachine_from_state_actions []byte = []byte{
	3, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var _myanmarSyllableMachine_eof_trans []int16 = []int16{
	0, 23, 23, 23, 23, 23, 23, 23,
	23, 23, 23, 23, 23, 23, 23, 23,
	23, 23, 23, 23, 23, 23, 23, 23,
	23, 23, 49, 49, 49, 49, 49, 49,
	49, 49, 49, 49, 49, 49, 49, 49,
	49, 49, 49, 49, 49, 49, 49, 49,
	49, 49, 49, 23, 23, 49, 62, 62,
}

const myanmarSyllableMachine_start int = 0
const myanmarSyllableMachine_first_final int = 0
const myanmarSyllableMachine_error int = -1

const myanmarSyllableMachine_en_main int = 0

func findSyllablesMyanmar(buffer *Buffer) {
	var p, ts, te, act, cs int
	info := buffer.Info

	{
		cs = myanmarSyllableMachine_start
		ts = 0
		te = 0
		act = 0
	}

	pe := len(info)
	eof := pe

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
		_acts = int(_myanmarSyllableMachine_from_state_actions[cs])
		_nacts = uint(_myanmarSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _myanmarSyllableMachine_actions[_acts-1] {
			case 1:
				ts = p

			}
		}

		_keys = int(_myanmarSyllableMachine_key_offsets[cs])
		_trans = int(_myanmarSyllableMachine_index_offsets[cs])

		_klen = int(_myanmarSyllableMachine_single_lengths[cs])
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
				case (info[p].complexCategory) < _myanmarSyllableMachine_trans_keys[_mid]:
					_upper = _mid - 1
				case (info[p].complexCategory) > _myanmarSyllableMachine_trans_keys[_mid]:
					_lower = _mid + 1
				default:
					_trans += int(_mid - int(_keys))
					goto _match
				}
			}
			_keys += _klen
			_trans += _klen
		}

		_klen = int(_myanmarSyllableMachine_range_lengths[cs])
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
				case (info[p].complexCategory) < _myanmarSyllableMachine_trans_keys[_mid]:
					_upper = _mid - 2
				case (info[p].complexCategory) > _myanmarSyllableMachine_trans_keys[_mid+1]:
					_lower = _mid + 2
				default:
					_trans += int((_mid - int(_keys)) >> 1)
					goto _match
				}
			}
			_trans += _klen
		}

	_match:
		_trans = int(_myanmarSyllableMachine_indicies[_trans])
	_eof_trans:
		cs = int(_myanmarSyllableMachine_trans_targs[_trans])

		if _myanmarSyllableMachine_trans_actions[_trans] == 0 {
			goto _again
		}

		_acts = int(_myanmarSyllableMachine_trans_actions[_trans])
		_nacts = uint(_myanmarSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _myanmarSyllableMachine_actions[_acts-1] {
			case 2:
				te = p + 1
				{
					foundSyllableMyanmar(myanmarConsonantSyllable, ts, te, info, &syllableSerial)
				}
			case 3:
				te = p + 1
				{
					foundSyllableMyanmar(myanmarNonMyanmarCluster, ts, te, info, &syllableSerial)
				}
			case 4:
				te = p + 1
				{
					foundSyllableMyanmar(myanmarPunctuationCluster, ts, te, info, &syllableSerial)
				}
			case 5:
				te = p + 1
				{
					foundSyllableMyanmar(myanmarBrokenCluster, ts, te, info, &syllableSerial)
				}
			case 6:
				te = p + 1
				{
					foundSyllableMyanmar(myanmarNonMyanmarCluster, ts, te, info, &syllableSerial)
				}
			case 7:
				te = p
				p--
				{
					foundSyllableMyanmar(myanmarConsonantSyllable, ts, te, info, &syllableSerial)
				}
			case 8:
				te = p
				p--
				{
					foundSyllableMyanmar(myanmarBrokenCluster, ts, te, info, &syllableSerial)
				}
			case 9:
				te = p
				p--
				{
					foundSyllableMyanmar(myanmarNonMyanmarCluster, ts, te, info, &syllableSerial)
				}
			}
		}

	_again:
		_acts = int(_myanmarSyllableMachine_to_state_actions[cs])
		_nacts = uint(_myanmarSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _myanmarSyllableMachine_actions[_acts-1] {
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
			if _myanmarSyllableMachine_eof_trans[cs] > 0 {
				_trans = int(_myanmarSyllableMachine_eof_trans[cs] - 1)
				goto _eof_trans
			}
		}

	}

	_ = act // needed by Ragel, but unused
}
