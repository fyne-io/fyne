package harfbuzz

// Code generated with ragel -Z -o ot_khmer_machine.go ot_khmer_machine.rl ; sed -i '/^\/\/line/ d' ot_khmer_machine.go ; goimports -w ot_khmer_machine.go  DO NOT EDIT.

// ported from harfbuzz/src/hb-ot-shape-complex-khmer-machine.rl Copyright © 2015 Google, Inc. Behdad Esfahbod

const (
	khmerConsonantSyllable = iota
	khmerBrokenCluster
	khmerNonKhmerCluster
)

const khmerSyllableMachine_ex_C = 1
const khmerSyllableMachine_ex_Coeng = 14
const khmerSyllableMachine_ex_DOTTEDCIRCLE = 12
const khmerSyllableMachine_ex_PLACEHOLDER = 11
const khmerSyllableMachine_ex_Ra = 16
const khmerSyllableMachine_ex_Robatic = 20
const khmerSyllableMachine_ex_V = 2
const khmerSyllableMachine_ex_VAbv = 26
const khmerSyllableMachine_ex_VBlw = 27
const khmerSyllableMachine_ex_VPre = 28
const khmerSyllableMachine_ex_VPst = 29
const khmerSyllableMachine_ex_Xgroup = 21
const khmerSyllableMachine_ex_Ygroup = 22
const khmerSyllableMachine_ex_ZWJ = 6
const khmerSyllableMachine_ex_ZWNJ = 5

var _khmerSyllableMachine_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 1, 5,
	1, 6, 1, 7, 1, 8, 1, 9,
	1, 10, 1, 11, 2, 2, 3, 2,
	2, 4,
}

var _khmerSyllableMachine_key_offsets []byte = []byte{
	0, 5, 8, 12, 15, 18, 21, 25,
	28, 32, 35, 38, 42, 45, 48, 51,
	55, 58, 62, 65, 70, 84, 94, 103,
	109, 110, 115, 122, 130, 139, 142, 146,
	155, 161, 162, 167, 174, 182, 185, 195,
}

var _khmerSyllableMachine_trans_keys []byte = []byte{
	20, 21, 26, 5, 6, 21, 5, 6,
	21, 26, 5, 6, 21, 5, 6, 16,
	1, 2, 21, 5, 6, 21, 26, 5,
	6, 21, 5, 6, 21, 26, 5, 6,
	21, 5, 6, 21, 5, 6, 21, 26,
	5, 6, 21, 5, 6, 16, 1, 2,
	21, 5, 6, 21, 26, 5, 6, 21,
	5, 6, 21, 26, 5, 6, 21, 5,
	6, 20, 21, 26, 5, 6, 14, 16,
	21, 22, 26, 27, 28, 29, 1, 2,
	5, 6, 11, 12, 14, 20, 21, 22,
	26, 27, 28, 29, 5, 6, 14, 21,
	22, 26, 27, 28, 29, 5, 6, 14,
	21, 22, 29, 5, 6, 22, 14, 21,
	22, 5, 6, 14, 21, 22, 26, 29,
	5, 6, 14, 21, 22, 26, 27, 29,
	5, 6, 14, 21, 22, 26, 27, 28,
	29, 5, 6, 16, 1, 2, 21, 26,
	5, 6, 14, 21, 22, 26, 27, 28,
	29, 5, 6, 14, 21, 22, 29, 5,
	6, 22, 14, 21, 22, 5, 6, 14,
	21, 22, 26, 29, 5, 6, 14, 21,
	22, 26, 27, 29, 5, 6, 16, 1,
	2, 14, 20, 21, 22, 26, 27, 28,
	29, 5, 6, 14, 21, 22, 26, 27,
	28, 29, 5, 6,
}

var _khmerSyllableMachine_single_lengths []byte = []byte{
	3, 1, 2, 1, 1, 1, 2, 1,
	2, 1, 1, 2, 1, 1, 1, 2,
	1, 2, 1, 3, 8, 8, 7, 4,
	1, 3, 5, 6, 7, 1, 2, 7,
	4, 1, 3, 5, 6, 1, 8, 7,
}

var _khmerSyllableMachine_range_lengths []byte = []byte{
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 1, 1, 3, 1, 1, 1,
	0, 1, 1, 1, 1, 1, 1, 1,
	1, 0, 1, 1, 1, 1, 1, 1,
}

var _khmerSyllableMachine_index_offsets []byte = []byte{
	0, 5, 8, 12, 15, 18, 21, 25,
	28, 32, 35, 38, 42, 45, 48, 51,
	55, 58, 62, 65, 70, 82, 92, 101,
	107, 109, 114, 121, 129, 138, 141, 145,
	154, 160, 162, 167, 174, 182, 185, 195,
}

var _khmerSyllableMachine_indicies []byte = []byte{
	2, 3, 4, 1, 0, 3, 1, 0,
	3, 4, 1, 0, 4, 5, 0, 6,
	6, 0, 8, 7, 0, 10, 4, 9,
	0, 10, 9, 0, 12, 4, 11, 0,
	12, 11, 0, 15, 14, 13, 15, 17,
	14, 16, 17, 18, 16, 19, 19, 16,
	21, 20, 16, 23, 17, 22, 16, 23,
	22, 16, 25, 17, 24, 16, 25, 24,
	16, 26, 15, 17, 14, 16, 30, 28,
	15, 19, 17, 23, 25, 21, 28, 29,
	2, 27, 33, 2, 3, 6, 4, 10,
	12, 8, 32, 31, 35, 3, 6, 4,
	10, 12, 8, 34, 31, 35, 4, 6,
	8, 5, 31, 6, 31, 35, 8, 6,
	7, 31, 35, 10, 6, 4, 8, 36,
	31, 35, 12, 6, 4, 10, 8, 37,
	31, 33, 3, 6, 4, 10, 12, 8,
	34, 31, 28, 28, 31, 15, 17, 14,
	38, 41, 15, 19, 17, 23, 25, 21,
	40, 39, 41, 17, 19, 21, 18, 39,
	19, 39, 41, 21, 19, 20, 39, 41,
	23, 19, 17, 21, 42, 39, 41, 25,
	19, 17, 23, 21, 43, 39, 44, 44,
	39, 30, 26, 15, 19, 17, 23, 25,
	21, 45, 39, 30, 15, 19, 17, 23,
	25, 21, 40, 39,
}

var _khmerSyllableMachine_trans_targs []byte = []byte{
	20, 1, 28, 22, 23, 3, 24, 5,
	25, 7, 26, 9, 27, 20, 10, 31,
	20, 32, 12, 33, 14, 34, 16, 35,
	18, 36, 39, 20, 21, 30, 37, 20,
	0, 29, 2, 4, 6, 8, 20, 20,
	11, 13, 15, 17, 38, 19,
}

var _khmerSyllableMachine_trans_actions []byte = []byte{
	15, 0, 5, 5, 5, 0, 0, 0,
	5, 0, 5, 0, 5, 19, 0, 21,
	17, 5, 0, 0, 0, 5, 0, 5,
	0, 5, 21, 7, 5, 24, 0, 9,
	0, 0, 0, 0, 0, 0, 13, 11,
	0, 0, 0, 0, 21, 0,
}

var _khmerSyllableMachine_to_state_actions []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 1, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var _khmerSyllableMachine_from_state_actions []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 3, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}

var _khmerSyllableMachine_eof_trans []byte = []byte{
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 1, 14, 17, 17, 17, 17, 17,
	17, 17, 17, 17, 0, 32, 32, 32,
	32, 32, 32, 32, 32, 32, 39, 40,
	40, 40, 40, 40, 40, 40, 40, 40,
}

const khmerSyllableMachine_start int = 20
const khmerSyllableMachine_first_final int = 20
const khmerSyllableMachine_error int = -1

const khmerSyllableMachine_en_main int = 20

func findSyllablesKhmer(buffer *Buffer) {
	var p, ts, te, act, cs int
	info := buffer.Info

	{
		cs = khmerSyllableMachine_start
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
		_acts = int(_khmerSyllableMachine_from_state_actions[cs])
		_nacts = uint(_khmerSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _khmerSyllableMachine_actions[_acts-1] {
			case 1:
				ts = p

			}
		}

		_keys = int(_khmerSyllableMachine_key_offsets[cs])
		_trans = int(_khmerSyllableMachine_index_offsets[cs])

		_klen = int(_khmerSyllableMachine_single_lengths[cs])
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
				case (info[p].complexCategory) < _khmerSyllableMachine_trans_keys[_mid]:
					_upper = _mid - 1
				case (info[p].complexCategory) > _khmerSyllableMachine_trans_keys[_mid]:
					_lower = _mid + 1
				default:
					_trans += int(_mid - int(_keys))
					goto _match
				}
			}
			_keys += _klen
			_trans += _klen
		}

		_klen = int(_khmerSyllableMachine_range_lengths[cs])
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
				case (info[p].complexCategory) < _khmerSyllableMachine_trans_keys[_mid]:
					_upper = _mid - 2
				case (info[p].complexCategory) > _khmerSyllableMachine_trans_keys[_mid+1]:
					_lower = _mid + 2
				default:
					_trans += int((_mid - int(_keys)) >> 1)
					goto _match
				}
			}
			_trans += _klen
		}

	_match:
		_trans = int(_khmerSyllableMachine_indicies[_trans])
	_eof_trans:
		cs = int(_khmerSyllableMachine_trans_targs[_trans])

		if _khmerSyllableMachine_trans_actions[_trans] == 0 {
			goto _again
		}

		_acts = int(_khmerSyllableMachine_trans_actions[_trans])
		_nacts = uint(_khmerSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _khmerSyllableMachine_actions[_acts-1] {
			case 2:
				te = p + 1

			case 3:
				act = 2
			case 4:
				act = 3
			case 5:
				te = p + 1
				{
					foundSyllableKhmer(khmerNonKhmerCluster, ts, te, info, &syllableSerial)
				}
			case 6:
				te = p
				p--
				{
					foundSyllableKhmer(khmerConsonantSyllable, ts, te, info, &syllableSerial)
				}
			case 7:
				te = p
				p--
				{
					foundSyllableKhmer(khmerBrokenCluster, ts, te, info, &syllableSerial)
				}
			case 8:
				te = p
				p--
				{
					foundSyllableKhmer(khmerNonKhmerCluster, ts, te, info, &syllableSerial)
				}
			case 9:
				p = (te) - 1
				{
					foundSyllableKhmer(khmerConsonantSyllable, ts, te, info, &syllableSerial)
				}
			case 10:
				p = (te) - 1
				{
					foundSyllableKhmer(khmerBrokenCluster, ts, te, info, &syllableSerial)
				}
			case 11:
				switch act {
				case 2:
					{
						p = (te) - 1
						foundSyllableKhmer(khmerBrokenCluster, ts, te, info, &syllableSerial)
					}
				case 3:
					{
						p = (te) - 1
						foundSyllableKhmer(khmerNonKhmerCluster, ts, te, info, &syllableSerial)
					}
				}

			}
		}

	_again:
		_acts = int(_khmerSyllableMachine_to_state_actions[cs])
		_nacts = uint(_khmerSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _khmerSyllableMachine_actions[_acts-1] {
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
			if _khmerSyllableMachine_eof_trans[cs] > 0 {
				_trans = int(_khmerSyllableMachine_eof_trans[cs] - 1)
				goto _eof_trans
			}
		}

	}

	_ = act // needed by Ragel, but unused
}
