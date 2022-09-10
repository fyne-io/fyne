package harfbuzz 

// Code generated with ragel -Z -o ot_myanmar_machine.go ot_myanmar_machine.rl ; sed -i '/^\/\/line/ d' ot_myanmar_machine.go ; goimports -w ot_myanmar_machine.go  DO NOT EDIT.

// ported from harfbuzz/src/hb-ot-shape-complex-myanmar-machine.rl Copyright © 2015 Mozilla Foundation. Google, Inc. Behdad Esfahbod

// myanmar_syllable_type_t
const  (
  myanmarConsonantSyllable = iota
  myanmarPunctuationCluster
  myanmarBrokenCluster
  myanmarNonMyanmarCluster
)

%%{
  machine myanmarSyllableMachine;
  alphtype byte;
  write exports;
  write data;
}%%

%%{

export A    = 10;
export As   = 18;
export C    = 1;
export D    = 32;
export D0   = 20;
export DB   = 3;
export GB   = 11;
export H    = 4;
export IV   = 2;
export MH   = 21;
export ML   = 33;
export MR   = 22;
export MW   = 23;
export MY   = 24;
export PT   = 25;
export V    = 8;
export VAbv = 26;
export VBlw = 27;
export VPre = 28;
export VPst = 29;
export VS   = 30;
export ZWJ  = 6;
export ZWNJ = 5;
export Ra   = 16;
export P    = 31;
export CS   = 19;

j = ZWJ|ZWNJ;			# Joiners
k = (Ra As H);			# Kinzi

c = C|Ra;			# is_consonant

medial_group = MY? As? MR? ((MW MH? ML? | MH ML? | ML) As?)?;
main_vowel_group = (VPre.VS?)* VAbv* VBlw* A* (DB As?)?;
post_vowel_group = VPst MH? ML? As* VAbv* A* (DB As?)?;
pwo_tone_group = PT A* DB? As?;

complex_syllable_tail = As* medial_group main_vowel_group post_vowel_group* pwo_tone_group* V* j?;
syllable_tail = (H (c|IV).VS?)* (H | complex_syllable_tail);

consonant_syllable =	(k|CS)? (c|IV|D|GB).VS? syllable_tail;
punctuation_cluster =	P V;
broken_cluster =	k? VS? syllable_tail;
other =			any;

main := |*
	consonant_syllable	=> { foundSyllableMyanmar (myanmarConsonantSyllable, ts, te, info, &syllableSerial); };
	j			=> { foundSyllableMyanmar (myanmarNonMyanmarCluster, ts, te, info, &syllableSerial); };
	punctuation_cluster	=> { foundSyllableMyanmar (myanmarPunctuationCluster, ts, te, info, &syllableSerial); };
	broken_cluster		=> { foundSyllableMyanmar (myanmarBrokenCluster, ts, te, info, &syllableSerial); };
	other			=> { foundSyllableMyanmar (myanmarNonMyanmarCluster, ts, te, info, &syllableSerial); };
*|;


}%%


func findSyllablesMyanmar (buffer *Buffer){
    var p, ts, te, act, cs int 
    info := buffer.Info;
    %%{
        write init;
        getkey info[p].complexCategory;
    }%%

    pe := len(info)
    eof := pe

    var syllableSerial uint8 = 1;
    %%{
        write exec;
    }%%
    _ = act // needed by Ragel, but unused
}

