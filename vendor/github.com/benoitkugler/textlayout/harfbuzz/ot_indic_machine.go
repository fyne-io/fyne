package harfbuzz

// Code generated with ragel -Z -o ot_indic_machine.go ot_indic_machine.rl ; sed -i '/^\/\/line/ d' ot_indic_machine.go ; goimports -w ot_indic_machine.go  DO NOT EDIT.

// ported from harfbuzz/src/hb-ot-shape-complex-indic-machine.rl Copyright © 2015 Google, Inc. Behdad Esfahbod

// indic_syllable_type_t
const (
	indicConsonantSyllable = iota
	indicVowelSyllable
	indicStandaloneCluster
	indicSymbolCluster
	indicBrokenCluster
	indicNonIndicCluster
)

const indicSyllableMachine_ex_A = 10
const indicSyllableMachine_ex_C = 1
const indicSyllableMachine_ex_CM = 17
const indicSyllableMachine_ex_CS = 19
const indicSyllableMachine_ex_DOTTEDCIRCLE = 12
const indicSyllableMachine_ex_H = 4
const indicSyllableMachine_ex_M = 7
const indicSyllableMachine_ex_N = 3
const indicSyllableMachine_ex_PLACEHOLDER = 11
const indicSyllableMachine_ex_RS = 13
const indicSyllableMachine_ex_Ra = 16
const indicSyllableMachine_ex_Repha = 15
const indicSyllableMachine_ex_SM = 8
const indicSyllableMachine_ex_Symbol = 18
const indicSyllableMachine_ex_V = 2
const indicSyllableMachine_ex_ZWJ = 6
const indicSyllableMachine_ex_ZWNJ = 5

var _indicSyllableMachine_actions []byte = []byte{
	0, 1, 0, 1, 1, 1, 2, 1, 6,
	1, 7, 1, 8, 1, 9, 1, 10,
	1, 11, 1, 12, 1, 13, 1, 14,
	1, 15, 1, 16, 1, 17, 1, 18,
	2, 2, 3, 2, 2, 4, 2, 2,
	5,
}

var _indicSyllableMachine_key_offsets []int16 = []int16{
	0, 1, 6, 9, 13, 18, 19, 20,
	25, 31, 36, 37, 40, 44, 49, 50,
	51, 56, 62, 68, 74, 75, 78, 82,
	87, 88, 89, 94, 99, 105, 106, 109,
	113, 118, 119, 120, 125, 130, 134, 135,
	152, 161, 169, 176, 182, 186, 189, 190,
	192, 199, 205, 211, 218, 224, 229, 235,
	239, 244, 248, 256, 265, 274, 282, 289,
	295, 304, 312, 319, 325, 328, 329, 331,
	338, 344, 351, 357, 362, 368, 372, 376,
	381, 385, 393, 402, 407, 415, 423, 430,
	436, 445, 451, 454, 455, 457, 464, 470,
	477, 483, 488, 496, 502, 506, 510, 515,
	519, 528, 534, 539, 548, 556, 563, 569,
	578, 584, 587, 588, 590, 597, 603, 610,
	616, 621, 629, 635, 639, 643, 648, 652,
	666, 675, 688, 695, 698, 699, 701, 710,
	715, 719, 722, 723, 725,
}

var _indicSyllableMachine_trans_keys []byte = []byte{
	8, 4, 7, 8, 5, 6, 7, 5,
	6, 7, 8, 5, 6, 4, 7, 8,
	5, 6, 6, 16, 4, 7, 8, 5,
	6, 4, 7, 8, 13, 5, 6, 4,
	7, 8, 5, 6, 8, 7, 5, 6,
	7, 8, 5, 6, 4, 7, 8, 5,
	6, 6, 16, 4, 7, 8, 5, 6,
	4, 7, 8, 13, 5, 6, 4, 7,
	8, 13, 5, 6, 4, 7, 8, 13,
	5, 6, 8, 7, 5, 6, 7, 8,
	5, 6, 4, 7, 8, 5, 6, 6,
	16, 4, 7, 8, 5, 6, 4, 7,
	8, 5, 6, 4, 7, 8, 13, 5,
	6, 8, 7, 5, 6, 7, 8, 5,
	6, 4, 7, 8, 5, 6, 6, 16,
	4, 7, 8, 5, 6, 4, 7, 8,
	5, 6, 7, 8, 5, 6, 8, 1,
	2, 3, 4, 5, 6, 7, 8, 10,
	13, 15, 16, 17, 18, 19, 11, 12,
	3, 4, 5, 6, 7, 8, 10, 13,
	17, 3, 4, 7, 8, 10, 17, 5,
	6, 4, 7, 8, 10, 17, 5, 6,
	1, 5, 6, 8, 10, 16, 8, 10,
	5, 6, 5, 8, 10, 10, 5, 10,
	1, 3, 8, 10, 16, 5, 6, 1,
	8, 10, 16, 5, 6, 1, 5, 6,
	8, 10, 16, 3, 4, 5, 6, 7,
	8, 10, 4, 5, 6, 7, 8, 10,
	7, 8, 10, 5, 6, 4, 7, 8,
	10, 5, 6, 5, 6, 8, 10, 3,
	8, 10, 5, 6, 5, 6, 8, 10,
	3, 4, 7, 8, 10, 17, 5, 6,
	3, 4, 5, 6, 7, 8, 10, 13,
	17, 3, 4, 5, 6, 7, 8, 10,
	13, 17, 3, 4, 5, 6, 7, 8,
	10, 17, 4, 5, 6, 7, 8, 10,
	17, 1, 5, 6, 8, 10, 16, 3,
	4, 5, 6, 7, 8, 10, 13, 17,
	3, 4, 7, 8, 10, 17, 5, 6,
	4, 7, 8, 10, 17, 5, 6, 1,
	5, 6, 8, 10, 16, 5, 8, 10,
	10, 5, 10, 1, 3, 8, 10, 16,
	5, 6, 1, 8, 10, 16, 5, 6,
	3, 4, 5, 6, 7, 8, 10, 4,
	5, 6, 7, 8, 10, 7, 8, 10,
	5, 6, 4, 7, 8, 10, 5, 6,
	5, 6, 8, 10, 8, 10, 5, 6,
	3, 8, 10, 5, 6, 5, 6, 8,
	10, 3, 4, 7, 8, 10, 17, 5,
	6, 3, 4, 5, 6, 7, 8, 10,
	13, 17, 4, 7, 8, 5, 6, 3,
	4, 5, 6, 7, 8, 10, 17, 3,
	4, 7, 8, 10, 17, 5, 6, 4,
	7, 8, 10, 17, 5, 6, 1, 5,
	6, 8, 10, 16, 3, 4, 5, 6,
	7, 8, 10, 13, 17, 1, 5, 6,
	8, 10, 16, 5, 8, 10, 10, 5,
	10, 1, 3, 8, 10, 16, 5, 6,
	1, 8, 10, 16, 5, 6, 3, 4,
	5, 6, 7, 8, 10, 4, 5, 6,
	7, 8, 10, 7, 8, 10, 5, 6,
	3, 4, 7, 8, 10, 17, 5, 6,
	4, 7, 8, 10, 5, 6, 5, 6,
	8, 10, 8, 10, 5, 6, 3, 8,
	10, 5, 6, 5, 6, 8, 10, 3,
	4, 5, 6, 7, 8, 10, 13, 17,
	4, 7, 8, 13, 5, 6, 4, 7,
	8, 5, 6, 3, 4, 5, 6, 7,
	8, 10, 13, 17, 3, 4, 7, 8,
	10, 17, 5, 6, 4, 7, 8, 10,
	17, 5, 6, 1, 5, 6, 8, 10,
	16, 3, 4, 5, 6, 7, 8, 10,
	13, 17, 1, 5, 6, 8, 10, 16,
	5, 8, 10, 10, 5, 10, 1, 3,
	8, 10, 16, 5, 6, 1, 8, 10,
	16, 5, 6, 3, 4, 5, 6, 7,
	8, 10, 4, 5, 6, 7, 8, 10,
	7, 8, 10, 5, 6, 3, 4, 7,
	8, 10, 17, 5, 6, 4, 7, 8,
	10, 5, 6, 5, 6, 8, 10, 8,
	10, 5, 6, 3, 8, 10, 5, 6,
	5, 6, 8, 10, 1, 2, 3, 4,
	5, 6, 7, 8, 10, 13, 16, 17,
	11, 12, 3, 4, 5, 6, 7, 8,
	10, 13, 17, 1, 2, 3, 4, 5,
	6, 7, 8, 10, 12, 13, 16, 17,
	4, 7, 8, 10, 13, 5, 6, 5,
	8, 10, 10, 5, 10, 1, 3, 4,
	7, 8, 10, 16, 5, 6, 3, 8,
	10, 5, 6, 8, 10, 5, 6, 5,
	8, 10, 10, 5, 10, 1, 11, 16,
}

var _indicSyllableMachine_single_lengths []byte = []byte{
	1, 3, 1, 2, 3, 1, 1, 3,
	4, 3, 1, 1, 2, 3, 1, 1,
	3, 4, 4, 4, 1, 1, 2, 3,
	1, 1, 3, 3, 4, 1, 1, 2,
	3, 1, 1, 3, 3, 2, 1, 15,
	9, 6, 5, 6, 2, 3, 1, 2,
	5, 4, 6, 7, 6, 3, 4, 4,
	3, 4, 6, 9, 9, 8, 7, 6,
	9, 6, 5, 6, 3, 1, 2, 5,
	4, 7, 6, 3, 4, 4, 2, 3,
	4, 6, 9, 3, 8, 6, 5, 6,
	9, 6, 3, 1, 2, 5, 4, 7,
	6, 3, 6, 4, 4, 2, 3, 4,
	9, 4, 3, 9, 6, 5, 6, 9,
	6, 3, 1, 2, 5, 4, 7, 6,
	3, 6, 4, 4, 2, 3, 4, 12,
	9, 13, 5, 3, 1, 2, 7, 3,
	2, 3, 1, 2, 3,
}

var _indicSyllableMachine_range_lengths []byte = []byte{
	0, 1, 1, 1, 1, 0, 0, 1,
	1, 1, 0, 1, 1, 1, 0, 0,
	1, 1, 1, 1, 0, 1, 1, 1,
	0, 0, 1, 1, 1, 0, 1, 1,
	1, 0, 0, 1, 1, 1, 0, 1,
	0, 1, 1, 0, 1, 0, 0, 0,
	1, 1, 0, 0, 0, 1, 1, 0,
	1, 0, 1, 0, 0, 0, 0, 0,
	0, 1, 1, 0, 0, 0, 0, 1,
	1, 0, 0, 1, 1, 0, 1, 1,
	0, 1, 0, 1, 0, 1, 1, 0,
	0, 0, 0, 0, 0, 1, 1, 0,
	0, 1, 1, 1, 0, 1, 1, 0,
	0, 1, 1, 0, 1, 1, 0, 0,
	0, 0, 0, 0, 1, 1, 0, 0,
	1, 1, 1, 0, 1, 1, 0, 1,
	0, 0, 1, 0, 0, 0, 1, 1,
	1, 0, 0, 0, 0,
}

var _indicSyllableMachine_index_offsets []int16 = []int16{
	0, 2, 7, 10, 14, 19, 21, 23,
	28, 34, 39, 41, 44, 48, 53, 55,
	57, 62, 68, 74, 80, 82, 85, 89,
	94, 96, 98, 103, 108, 114, 116, 119,
	123, 128, 130, 132, 137, 142, 146, 148,
	165, 175, 183, 190, 197, 201, 205, 207,
	210, 217, 223, 230, 238, 245, 250, 256,
	261, 266, 271, 279, 289, 299, 308, 316,
	323, 333, 341, 348, 355, 359, 361, 364,
	371, 377, 385, 392, 397, 403, 408, 412,
	417, 422, 430, 440, 445, 454, 462, 469,
	476, 486, 493, 497, 499, 502, 509, 515,
	523, 530, 535, 543, 549, 554, 558, 563,
	568, 578, 584, 589, 599, 607, 614, 621,
	631, 638, 642, 644, 647, 654, 660, 668,
	675, 680, 688, 694, 699, 703, 708, 713,
	727, 737, 751, 758, 762, 764, 767, 776,
	781, 785, 789, 791, 794,
}

var _indicSyllableMachine_indicies []byte = []byte{
	1, 0, 2, 4, 1, 3, 0, 4,
	3, 0, 4, 1, 3, 0, 5, 4,
	1, 3, 0, 6, 0, 7, 0, 8,
	4, 1, 3, 0, 2, 4, 1, 9,
	3, 0, 11, 13, 14, 12, 10, 14,
	10, 13, 12, 10, 13, 14, 12, 10,
	15, 13, 14, 12, 10, 16, 10, 17,
	10, 18, 13, 14, 12, 10, 11, 13,
	14, 19, 12, 10, 11, 13, 14, 20,
	12, 10, 22, 24, 25, 26, 23, 21,
	25, 21, 24, 23, 27, 24, 25, 23,
	21, 28, 24, 25, 23, 21, 29, 21,
	30, 21, 22, 24, 25, 23, 21, 31,
	24, 25, 23, 21, 33, 35, 36, 37,
	34, 32, 36, 32, 35, 34, 32, 35,
	36, 34, 32, 38, 35, 36, 34, 32,
	39, 32, 40, 32, 33, 35, 36, 34,
	32, 41, 35, 36, 34, 32, 24, 1,
	23, 0, 43, 42, 45, 46, 47, 48,
	49, 50, 24, 25, 51, 26, 53, 54,
	55, 56, 57, 52, 44, 59, 60, 61,
	62, 4, 1, 63, 9, 64, 58, 65,
	60, 4, 1, 63, 64, 66, 58, 60,
	4, 1, 63, 64, 66, 58, 45, 67,
	68, 1, 63, 45, 58, 1, 63, 69,
	58, 63, 70, 63, 58, 63, 58, 63,
	63, 58, 45, 71, 1, 63, 45, 69,
	58, 45, 1, 63, 45, 69, 58, 45,
	69, 68, 1, 63, 45, 58, 72, 7,
	73, 74, 4, 1, 63, 58, 7, 73,
	74, 4, 1, 63, 58, 4, 1, 63,
	73, 58, 75, 4, 1, 63, 76, 58,
	67, 77, 1, 63, 58, 67, 1, 63,
	69, 58, 69, 77, 1, 63, 58, 59,
	60, 4, 1, 63, 64, 66, 58, 59,
	60, 61, 66, 4, 1, 63, 9, 64,
	58, 79, 80, 81, 82, 13, 14, 83,
	20, 84, 78, 85, 80, 86, 82, 13,
	14, 83, 84, 78, 80, 86, 82, 13,
	14, 83, 84, 78, 87, 88, 89, 14,
	83, 87, 78, 90, 80, 91, 92, 13,
	14, 83, 19, 84, 78, 93, 80, 13,
	14, 83, 84, 86, 78, 80, 13, 14,
	83, 84, 86, 78, 87, 94, 89, 14,
	83, 87, 78, 83, 95, 83, 78, 83,
	78, 83, 83, 78, 87, 96, 14, 83,
	87, 94, 78, 87, 14, 83, 87, 94,
	78, 97, 17, 98, 99, 13, 14, 83,
	78, 17, 98, 99, 13, 14, 83, 78,
	13, 14, 83, 98, 78, 100, 13, 14,
	83, 101, 78, 88, 102, 14, 83, 78,
	14, 83, 94, 78, 88, 14, 83, 94,
	78, 94, 102, 14, 83, 78, 90, 80,
	13, 14, 83, 84, 86, 78, 90, 80,
	91, 86, 13, 14, 83, 19, 84, 78,
	11, 13, 14, 12, 78, 79, 80, 86,
	82, 13, 14, 83, 84, 78, 104, 48,
	24, 25, 51, 55, 105, 103, 48, 24,
	25, 51, 55, 105, 103, 106, 107, 108,
	25, 51, 106, 103, 47, 48, 109, 110,
	24, 25, 51, 26, 55, 103, 106, 111,
	108, 25, 51, 106, 103, 51, 112, 51,
	103, 51, 103, 51, 51, 103, 106, 113,
	25, 51, 106, 111, 103, 106, 25, 51,
	106, 111, 103, 114, 30, 115, 116, 24,
	25, 51, 103, 30, 115, 116, 24, 25,
	51, 103, 24, 25, 51, 115, 103, 47,
	48, 24, 25, 51, 55, 105, 103, 117,
	24, 25, 51, 118, 103, 107, 119, 25,
	51, 103, 25, 51, 111, 103, 107, 25,
	51, 111, 103, 111, 119, 25, 51, 103,
	47, 48, 109, 105, 24, 25, 51, 26,
	55, 103, 22, 24, 25, 26, 23, 120,
	22, 24, 25, 23, 120, 122, 123, 124,
	125, 35, 36, 126, 37, 127, 121, 128,
	123, 35, 36, 126, 127, 125, 121, 123,
	35, 36, 126, 127, 125, 121, 129, 130,
	131, 36, 126, 129, 121, 122, 123, 124,
	52, 35, 36, 126, 37, 127, 121, 129,
	132, 131, 36, 126, 129, 121, 126, 133,
	126, 121, 126, 121, 126, 126, 121, 129,
	134, 36, 126, 129, 132, 121, 129, 36,
	126, 129, 132, 121, 135, 40, 136, 137,
	35, 36, 126, 121, 40, 136, 137, 35,
	36, 126, 121, 35, 36, 126, 136, 121,
	122, 123, 35, 36, 126, 127, 125, 121,
	138, 35, 36, 126, 139, 121, 130, 140,
	36, 126, 121, 36, 126, 132, 121, 130,
	36, 126, 132, 121, 132, 140, 36, 126,
	121, 45, 46, 47, 48, 109, 105, 24,
	25, 51, 26, 45, 55, 52, 103, 59,
	141, 61, 62, 4, 1, 63, 9, 64,
	58, 45, 46, 47, 48, 142, 143, 24,
	144, 145, 52, 26, 45, 55, 58, 22,
	24, 144, 63, 26, 146, 58, 145, 147,
	145, 58, 145, 58, 145, 145, 58, 45,
	71, 22, 24, 144, 63, 45, 146, 58,
	149, 43, 151, 150, 148, 43, 151, 150,
	148, 151, 152, 151, 148, 151, 148, 151,
	151, 148, 45, 52, 45, 120,
}

var _indicSyllableMachine_trans_targs []byte = []byte{
	39, 45, 50, 2, 51, 5, 6, 53,
	57, 58, 39, 67, 11, 73, 68, 14,
	15, 75, 80, 81, 84, 39, 89, 21,
	95, 90, 98, 39, 24, 25, 97, 103,
	39, 112, 30, 118, 113, 121, 33, 34,
	120, 126, 39, 137, 39, 40, 60, 85,
	87, 105, 106, 91, 107, 127, 128, 99,
	135, 140, 39, 41, 43, 8, 59, 46,
	54, 42, 1, 44, 48, 0, 47, 49,
	52, 3, 4, 55, 7, 56, 39, 61,
	63, 18, 83, 69, 76, 62, 9, 64,
	78, 71, 65, 17, 82, 66, 10, 70,
	72, 74, 12, 13, 77, 16, 79, 39,
	86, 26, 88, 101, 93, 19, 104, 20,
	92, 94, 96, 22, 23, 100, 27, 102,
	39, 39, 108, 110, 28, 35, 114, 122,
	109, 111, 124, 116, 29, 115, 117, 119,
	31, 32, 123, 36, 125, 129, 130, 134,
	131, 132, 37, 133, 39, 136, 38, 138,
	139,
}

var _indicSyllableMachine_trans_actions []byte = []byte{
	21, 0, 5, 0, 5, 0, 0, 5,
	5, 5, 23, 5, 0, 5, 0, 0,
	0, 5, 5, 5, 5, 29, 5, 0,
	36, 0, 36, 31, 0, 0, 36, 5,
	25, 5, 0, 5, 0, 5, 0, 0,
	5, 5, 27, 0, 7, 5, 5, 36,
	0, 39, 39, 0, 5, 36, 5, 36,
	5, 0, 9, 5, 0, 0, 5, 0,
	5, 5, 0, 5, 5, 0, 0, 5,
	5, 0, 0, 0, 0, 5, 11, 5,
	0, 0, 5, 0, 5, 5, 0, 5,
	5, 5, 5, 0, 5, 5, 0, 0,
	5, 5, 0, 0, 0, 0, 5, 17,
	36, 0, 36, 5, 5, 0, 36, 0,
	0, 5, 36, 0, 0, 0, 0, 5,
	19, 13, 5, 0, 0, 0, 0, 5,
	5, 5, 5, 5, 0, 0, 5, 5,
	0, 0, 0, 0, 5, 0, 33, 33,
	0, 0, 0, 0, 15, 5, 0, 0,
	0,
}

var _indicSyllableMachine_to_state_actions []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 1,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

var _indicSyllableMachine_from_state_actions []byte = []byte{
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 3,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0,
}

var _indicSyllableMachine_eof_trans []int16 = []int16{
	1, 1, 1, 1, 1, 1, 1, 1,
	1, 11, 11, 11, 11, 11, 11, 11,
	11, 11, 11, 22, 22, 28, 22, 22,
	22, 22, 22, 22, 33, 33, 33, 33,
	33, 33, 33, 33, 33, 1, 43, 0,
	59, 59, 59, 59, 59, 59, 59, 59,
	59, 59, 59, 59, 59, 59, 59, 59,
	59, 59, 59, 59, 79, 79, 79, 79,
	79, 79, 79, 79, 79, 79, 79, 79,
	79, 79, 79, 79, 79, 79, 79, 79,
	79, 79, 79, 79, 79, 104, 104, 104,
	104, 104, 104, 104, 104, 104, 104, 104,
	104, 104, 104, 104, 104, 104, 104, 104,
	104, 121, 121, 122, 122, 122, 122, 122,
	122, 122, 122, 122, 122, 122, 122, 122,
	122, 122, 122, 122, 122, 122, 122, 104,
	59, 59, 59, 59, 59, 59, 59, 149,
	149, 149, 149, 149, 121,
}

const indicSyllableMachine_start int = 39
const indicSyllableMachine_first_final int = 39
const indicSyllableMachine_error int = -1

const indicSyllableMachine_en_main int = 39

func findSyllablesIndic(buffer *Buffer) {
	var p, ts, te, act, cs int
	info := buffer.Info

	{
		cs = indicSyllableMachine_start
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
		_acts = int(_indicSyllableMachine_from_state_actions[cs])
		_nacts = uint(_indicSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _indicSyllableMachine_actions[_acts-1] {
			case 1:
				ts = p

			}
		}

		_keys = int(_indicSyllableMachine_key_offsets[cs])
		_trans = int(_indicSyllableMachine_index_offsets[cs])

		_klen = int(_indicSyllableMachine_single_lengths[cs])
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
				case (info[p].complexCategory) < _indicSyllableMachine_trans_keys[_mid]:
					_upper = _mid - 1
				case (info[p].complexCategory) > _indicSyllableMachine_trans_keys[_mid]:
					_lower = _mid + 1
				default:
					_trans += int(_mid - int(_keys))
					goto _match
				}
			}
			_keys += _klen
			_trans += _klen
		}

		_klen = int(_indicSyllableMachine_range_lengths[cs])
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
				case (info[p].complexCategory) < _indicSyllableMachine_trans_keys[_mid]:
					_upper = _mid - 2
				case (info[p].complexCategory) > _indicSyllableMachine_trans_keys[_mid+1]:
					_lower = _mid + 2
				default:
					_trans += int((_mid - int(_keys)) >> 1)
					goto _match
				}
			}
			_trans += _klen
		}

	_match:
		_trans = int(_indicSyllableMachine_indicies[_trans])
	_eof_trans:
		cs = int(_indicSyllableMachine_trans_targs[_trans])

		if _indicSyllableMachine_trans_actions[_trans] == 0 {
			goto _again
		}

		_acts = int(_indicSyllableMachine_trans_actions[_trans])
		_nacts = uint(_indicSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _indicSyllableMachine_actions[_acts-1] {
			case 2:
				te = p + 1

			case 3:
				act = 1
			case 4:
				act = 5
			case 5:
				act = 6
			case 6:
				te = p + 1
				{
					foundSyllableIndic(indicNonIndicCluster, ts, te, info, &syllableSerial)
				}
			case 7:
				te = p
				p--
				{
					foundSyllableIndic(indicConsonantSyllable, ts, te, info, &syllableSerial)
				}
			case 8:
				te = p
				p--
				{
					foundSyllableIndic(indicVowelSyllable, ts, te, info, &syllableSerial)
				}
			case 9:
				te = p
				p--
				{
					foundSyllableIndic(indicStandaloneCluster, ts, te, info, &syllableSerial)
				}
			case 10:
				te = p
				p--
				{
					foundSyllableIndic(indicSymbolCluster, ts, te, info, &syllableSerial)
				}
			case 11:
				te = p
				p--
				{
					foundSyllableIndic(indicBrokenCluster, ts, te, info, &syllableSerial)
				}
			case 12:
				te = p
				p--
				{
					foundSyllableIndic(indicNonIndicCluster, ts, te, info, &syllableSerial)
				}
			case 13:
				p = (te) - 1
				{
					foundSyllableIndic(indicConsonantSyllable, ts, te, info, &syllableSerial)
				}
			case 14:
				p = (te) - 1
				{
					foundSyllableIndic(indicVowelSyllable, ts, te, info, &syllableSerial)
				}
			case 15:
				p = (te) - 1
				{
					foundSyllableIndic(indicStandaloneCluster, ts, te, info, &syllableSerial)
				}
			case 16:
				p = (te) - 1
				{
					foundSyllableIndic(indicSymbolCluster, ts, te, info, &syllableSerial)
				}
			case 17:
				p = (te) - 1
				{
					foundSyllableIndic(indicBrokenCluster, ts, te, info, &syllableSerial)
				}
			case 18:
				switch act {
				case 1:
					{
						p = (te) - 1
						foundSyllableIndic(indicConsonantSyllable, ts, te, info, &syllableSerial)
					}
				case 5:
					{
						p = (te) - 1
						foundSyllableIndic(indicBrokenCluster, ts, te, info, &syllableSerial)
					}
				case 6:
					{
						p = (te) - 1
						foundSyllableIndic(indicNonIndicCluster, ts, te, info, &syllableSerial)
					}
				}

			}
		}

	_again:
		_acts = int(_indicSyllableMachine_to_state_actions[cs])
		_nacts = uint(_indicSyllableMachine_actions[_acts])
		_acts++
		for ; _nacts > 0; _nacts-- {
			_acts++
			switch _indicSyllableMachine_actions[_acts-1] {
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
			if _indicSyllableMachine_eof_trans[cs] > 0 {
				_trans = int(_indicSyllableMachine_eof_trans[cs] - 1)
				goto _eof_trans
			}
		}

	}

	_ = act // needed by Ragel, but unused
}
