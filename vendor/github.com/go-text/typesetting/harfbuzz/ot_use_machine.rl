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

%%{
  machine useSyllableMachine;
  alphtype byte;
  write exports;
  write data;
}%%

%%{

# Categories used in the Universal Shaping Engine spec:
# https://docs.microsoft.com/en-us/typography/script-development/use

export O	= 0; # OTHER

export B	= 1; # BASE
export N	= 4; # BASE_NUM
export GB	= 5; # BASE_OTHER
export CGJ	= 6; # CGJ
export SUB	= 11; # CONS_SUB
export H	= 12; # HALANT

export HN	= 13; # HALANT_NUM
export ZWNJ	= 14; # Zero width non-joiner
export R	= 18; # REPHA
export CS	= 43; # CONS_WITH_STACKER
export HVM	= 44; # HALANT_OR_VOWEL_MODIFIER
export Sk	= 48; # SAKOT
export G	= 49; # HIEROGLYPH
export J	= 50; # HIEROGLYPH_JOINER
export SB	= 51; # HIEROGLYPH_SEGMENT_BEGIN
export SE	= 52; # HIEROGLYPH_SEGMENT_END

export FAbv	= 24; # CONS_FINAL_ABOVE
export FBlw	= 25; # CONS_FINAL_BELOW
export FPst	= 26; # CONS_FINAL_POST
export MAbv	= 27; # CONS_MED_ABOVE
export MBlw	= 28; # CONS_MED_BELOW
export MPst	= 29; # CONS_MED_POST
export MPre	= 30; # CONS_MED_PRE
export CMAbv	= 31; # CONS_MOD_ABOVE
export CMBlw	= 32; # CONS_MOD_BELOW
export VAbv	= 33; # VOWEL_ABOVE / VOWEL_ABOVE_BELOW / VOWEL_ABOVE_BELOW_POST / VOWEL_ABOVE_POST
export VBlw	= 34; # VOWEL_BELOW / VOWEL_BELOW_POST
export VPst	= 35; # VOWEL_POST	UIPC = Right
export VPre	= 22; # VOWEL_PRE / VOWEL_PRE_ABOVE / VOWEL_PRE_ABOVE_POST / VOWEL_PRE_POST
export VMAbv	= 37; # VOWEL_MOD_ABOVE
export VMBlw	= 38; # VOWEL_MOD_BELOW
export VMPst	= 39; # VOWEL_MOD_POST
export VMPre	= 23; # VOWEL_MOD_PRE
export SMAbv	= 41; # SYM_MOD_ABOVE
export SMBlw	= 42; # SYM_MOD_BELOW
export FMAbv	= 45; # CONS_FINAL_MOD	UIPC = Top
export FMBlw	= 46; # CONS_FINAL_MOD	UIPC = Bottom
export FMPst	= 47; # CONS_FINAL_MOD	UIPC = Not_Applicable

h = H | HVM | Sk;

consonant_modifiers = CMAbv* CMBlw* ((h B | SUB) CMAbv? CMBlw*)*;
medial_consonants = MPre? MAbv? MBlw? MPst?;
dependent_vowels = VPre* VAbv* VBlw* VPst*;
vowel_modifiers = HVM? VMPre* VMAbv* VMBlw* VMPst*;
final_consonants = FAbv* FBlw* FPst*;
final_modifiers = FMAbv* FMBlw* | FMPst?;

complex_syllable_start = (R | CS)? (B | GB);
complex_syllable_middle =
	consonant_modifiers
	medial_consonants
	dependent_vowels
	vowel_modifiers
	(Sk B)*
;
complex_syllable_tail =
	complex_syllable_middle
	final_consonants
	final_modifiers
;
number_joiner_terminated_cluster_tail = (HN N)* HN;
numeral_cluster_tail = (HN N)+;
symbol_cluster_tail = SMAbv+ SMBlw* | SMBlw+;

virama_terminated_cluster_tail =
	consonant_modifiers
	h
;
virama_terminated_cluster =
	complex_syllable_start
	virama_terminated_cluster_tail
;
sakot_terminated_cluster_tail =
	complex_syllable_middle
	Sk
;
sakot_terminated_cluster =
	complex_syllable_start
	sakot_terminated_cluster_tail
;
standard_cluster =
	complex_syllable_start
	complex_syllable_tail
;
broken_cluster =
	R?
	(complex_syllable_tail | number_joiner_terminated_cluster_tail | numeral_cluster_tail | symbol_cluster_tail | virama_terminated_cluster_tail | sakot_terminated_cluster_tail)
;

number_joiner_terminated_cluster = N number_joiner_terminated_cluster_tail;
numeral_cluster = N numeral_cluster_tail?;
symbol_cluster = (O | GB) symbol_cluster_tail?;
hieroglyph_cluster = SB+ | SB* G SE* (J SE* (G SE*)?)*;

other = any;

main := |*
	virama_terminated_cluster		=> { foundSyllableUSE (useViramaTerminatedCluster,data, ts, te, info, &syllableSerial); };
	sakot_terminated_cluster		=> { foundSyllableUSE (useSakotTerminatedCluster,data, ts, te, info, &syllableSerial); };
	standard_cluster			=> { foundSyllableUSE (useStandardCluster,data, ts, te, info, &syllableSerial); };
	number_joiner_terminated_cluster	=> { foundSyllableUSE (useNumberJoinerTerminatedCluster,data, ts, te, info, &syllableSerial); };
	numeral_cluster				=> { foundSyllableUSE (useNumeralCluster,data, ts, te, info, &syllableSerial); };
	symbol_cluster				=> { foundSyllableUSE (useSymbolCluster,data, ts, te, info, &syllableSerial); };
	hieroglyph_cluster			=> { foundSyllableUSE (useHieroglyphCluster,data, ts, te, info, &syllableSerial); };
	broken_cluster				=> { foundSyllableUSE (useBrokenCluster,data, ts, te, info, &syllableSerial); };
	other					=> { foundSyllableUSE (useNonCluster,data, ts, te, info, &syllableSerial); };
*|;

}%%

func findSyllablesUse (buffer * Buffer) {
	info := buffer.Info
	data := preprocessInfoUSE(info)
    p, pe := 0, len(data)
	eof := pe
	var cs, act, ts, te int
	%%{
		write init;
		getkey (data[p]).p.v.complexCategory;
	}%%

	var syllableSerial uint8 = 1;
	%%{
		write exec;
	}%%
	_ = act // needed by Ragel, but unused
}

 
