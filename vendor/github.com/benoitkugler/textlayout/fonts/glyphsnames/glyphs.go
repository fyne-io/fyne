// copied from https://git.maze.io/go/unipdf/src/branch/master/internal/textencoding
package glyphsnames

import (
	"regexp"
	"strconv"
	"strings"
)

// GlyphToRune returns the rune corresponding to glyph `glyph` if there is one.
func GlyphToRune(glyph string) (rune, bool) {
	// We treat glyph "eight.lf" the same as glyph "eight".
	if strings.Contains(glyph, ".") {
		groups := rePrefix.FindStringSubmatch(glyph)
		if groups != nil {
			glyph = groups[1]
		}
	}
	// First lookup the glyph in all the tables.
	if alias := glyphToAlias(glyph); alias != "" {
		glyph = alias
	}
	if r := glyphToRune(glyph); r != 0 {
		return r, true
	}
	if r := ligatureToCodePoint(glyph); r != 0 {
		return r, true
	}

	// Next try all the glyph naming conventions.
	if groups := reUniEncoding.FindStringSubmatch(glyph); groups != nil {
		n, err := strconv.ParseInt(groups[1], 16, 32)
		if err == nil {
			return rune(n), true
		}
	}

	if groups := reEncoding.FindStringSubmatch(glyph); groups != nil {
		n, err := strconv.Atoi(groups[1])
		if err == nil {
			return rune(n), true
		}
	}

	return 0, false
}

var (
	reEncoding    = regexp.MustCompile(`^[A-Za-z](\d{1,5})$`) // C211
	reUniEncoding = regexp.MustCompile(`^uni([\dA-F]{4})$`)   // uniFB03
	rePrefix      = regexp.MustCompile(`^(\w+)\.\w+$`)        // eight.pnum => eight
)

// ligatureToCodePoint maps ligatures without corresponding unicode code points. We use the Unicode private
// use area (https://en.wikipedia.org/wiki/Private_Use_Areas) to store them.
// These runes are mapped to strings in RuneToString which uses the reverse mappings in
// ligatureToString.
func ligatureToCodePoint(ligature string) rune {
	switch ligature {
	case "f_t":
		return 0xe000
	case "f_j":
		return 0xe001
	case "f_b":
		return 0xe002
	case "f_h":
		return 0xe003
	case "f_k":
		return 0xe004
	case "t_t":
		return 0xe005
	case "t_f":
		return 0xe006
	case "f_f_j":
		return 0xe007
	case "f_f_b":
		return 0xe008
	case "f_f_h":
		return 0xe009
	case "f_f_k":
		return 0xe00a
	case "T_h":
		return 0xe00b
	}

	return 0
}

func glyphToAlias(glyph string) string {
	// 2462 entries
	switch glyph {
	case "f_f":
		return "ff"
	case "f_f_i":
		return "ffi"
	case "f_f_l":
		return "ffl"
	case "f_i":
		return "fi"
	case "f_l":
		return "fl"
	case "ascriptturn":
		return "AEmacron"
	case "mturndescend":
		return "Adblgrave"
	case "aturn":
		return "Adotmacron"
	case "nlftlfthook":
		return "Ainvertedbreve"
	case "upAlpha":
		return "Alpha"
	case "Ismallcap":
		return "Aringacute"
	case "Cbb":
		return "BbbC"
	case "Cdblstruck":
		return "BbbC"
	case "Hbb":
		return "BbbH"
	case "Hdblstruck":
		return "BbbH"
	case "Nbb":
		return "BbbN"
	case "Ndblstruck":
		return "BbbN"
	case "Pbb":
		return "BbbP"
	case "Pdblstruck":
		return "BbbP"
	case "Qbb":
		return "BbbQ"
	case "Qdblstruck":
		return "BbbQ"
	case "Rbb":
		return "BbbR"
	case "Rdblstruck":
		return "BbbR"
	case "Zbb":
		return "BbbZ"
	case "Zdblstruck":
		return "BbbZ"
	case "upBeta":
		return "Beta"
	case "OI":
		return "Btopbar"
	case "Hmacron":
		return "Cacute"
	case "Cdot":
		return "Cdotaccent"
	case "Che":
		return "Checyrillic"
	case "afii10041":
		return "Checyrillic"
	case "lcircumflex":
		return "Chedescendercyrillic"
	case "upChi":
		return "Chi"
	case "yusbig":
		return "Chi"
	case "gcursive":
		return "DZ"
	case "Gbar":
		return "DZcaron"
	case "Dslash":
		return "Dcroat"
	case "De":
		return "Decyrillic"
	case "afii10021":
		return "Decyrillic"
	case "Khartdes":
		return "Deicoptic"
	case "increment":
		return "Delta"
	case "upDelta":
		return "Deltagreek"
	case "eshlooprev":
		return "Dhook"
	case "mbfdigamma":
		return "Digamma"
	case "GeKarev":
		return "Digammagreek"
	case "upDigamma":
		return "Digammagreek"
	case "Gsmallcap":
		return "Dz"
	case "gbar":
		return "Dzcaron"
	case "Dzhe":
		return "Dzhecyrillic"
	case "afii10145":
		return "Dzhecyrillic"
	case "Ecyril":
		return "Ecyrillic"
	case "afii10053":
		return "Ecyrillic"
	case "Nsmallcap":
		return "Edblgrave"
	case "Edot":
		return "Edotaccent"
	case "OEsmallcap":
		return "Einvertedbreve"
	case "El":
		return "Elcyrillic"
	case "afii10029":
		return "Elcyrillic"
	case "Em":
		return "Emcyrillic"
	case "afii10030":
		return "Emcyrillic"
	case "Ng":
		return "Eng"
	case "kra":
		return "Eogonek"
	case "upEpsilon":
		return "Epsilon"
	case "strictequivalence":
		return "Equiv"
	case "Trthook":
		return "Ereversed"
	case "Ecyrilrev":
		return "Ereversedcyrillic"
	case "afii10047":
		return "Ereversedcyrillic"
	case "upEta":
		return "Eta"
	case "Euler":
		return "Eulerconst"
	case "euro":
		return "Euro"
	case "epsilon1revclosed":
		return "Ezhcaron"
	case "Ohook":
		return "Feicoptic"
	case "Upsilon2":
		return "Fhook"
	case "Fturn":
		return "Finv"
	case "FFIsmall":
		return "Fsmall"
	case "FFLsmall":
		return "Fsmall"
	case "FFsmall":
		return "Fsmall"
	case "FIsmall":
		return "Fsmall"
	case "FLsmall":
		return "Fsmall"
	case "babygamma":
		return "Gacute"
	case "upGamma":
		return "Gamma"
	case "Ustrt":
		return "Gangiacoptic"
	case "drthook":
		return "Gcaron"
	case "Gcedilla":
		return "Gcommaaccent"
	case "Gdot":
		return "Gdotaccent"
	case "Ge":
		return "Gecyrillic"
	case "afii10020":
		return "Gecyrillic"
	case "Geupturn":
		return "Gheupturncyrillic"
	case "afii10050":
		return "Gheupturncyrillic"
	case "Game":
		return "Gmir"
	case "ogoneknosp":
		return "Gsmallhook"
	case "cturn":
		return "Gstroke"
	case "whitesquare":
		return "H22073"
	case "box":
		return "H22073"
	case "mdlgwhtsquare":
		return "H22073"
	case "square":
		return "H22073"
	case "Tertdes":
		return "Horicoptic"
	case "Inodot":
		return "I"
	case "yoghhacek":
		return "Icaron"
	case "Idotaccent":
		return "Idot"
	case "Ie":
		return "Iecyrillic"
	case "afii10022":
		return "Iecyrillic"
	case "Iblackletter":
		return "Ifraktur"
	case "Ifractur":
		return "Ifraktur"
	case "Im":
		return "Ifraktur"
	case "Ii":
		return "Iicyrillic"
	case "afii10026":
		return "Iicyrillic"
	case "rturnascend":
		return "Iinvertedbreve"
	case "Io":
		return "Iocyrillic"
	case "afii10023":
		return "Iocyrillic"
	case "upIota":
		return "Iota"
	case "zbar":
		return "Iotaafrican"
	case "Yogh":
		return "Istroke"
	case "upKappa":
		return "Kappa"
	case "erev":
		return "Kcaron"
	case "Kcommaaccent":
		return "Kcedilla"
	case "Kha":
		return "Khacyrillic"
	case "afii10039":
		return "Khacyrillic"
	case "Escedilla":
		return "Kheicoptic"
	case "Yoghrev":
		return "Khook"
	case "Kje":
		return "Kjecyrillic"
	case "afii10061":
		return "Kjecyrillic"
	case "Enrtdes":
		return "Koppagreek"
	case "upKoppa":
		return "Koppagreek"
	case "ghacek":
		return "LJ"
	case "upLambda":
		return "Lambda"
	case "Lcommaaccent":
		return "Lcedilla"
	case "gcedilla1":
		return "Lcedilla1"
	case "Ldot":
		return "Ldotaccent"
	case "Khacek":
		return "Lj"
	case "Lje":
		return "Ljecyrillic"
	case "afii10058":
		return "Ljecyrillic"
	case "upMu":
		return "Mu"
	case "tmacron":
		return "Ncaron"
	case "Ncedilla":
		return "Ncommaaccent"
	case "tquoteright":
		return "Ncommaaccent"
	case "arrowdblne":
		return "Nearrow"
	case "upNu":
		return "Nu"
	case "arrowdblnw":
		return "Nwarrow"
	case "Ocyril":
		return "Ocyrillic"
	case "afii10032":
		return "Ocyrillic"
	case "Ohungarumlaut":
		return "Odblacute"
	case "rdescend":
		return "Odblgrave"
	case "pipe":
		return "Ohorn"
	case "pipedblbar":
		return "Oi"
	case "Ohm":
		return "Omega"
	case "ohm":
		return "Omega"
	case "upOmega":
		return "Omegagreek"
	case "mho":
		return "Omegainv"
	case "ohminverted":
		return "Omegainv"
	case "upOmicron":
		return "Omicron"
	case "yat":
		return "Omicron"
	case "epsilon1rev":
		return "Oogonekmacron"
	case "YR":
		return "Oopen"
	case "Ostrokeacute":
		return "Oslashacute"
	case "lyogh":
		return "Oslashacute"
	case "Yusbig":
		return "Phi"
	case "upPhi":
		return "Phi"
	case "DZhacek":
		return "Phook"
	case "upPi":
		return "Pi"
	case "planck":
		return "Planckconst"
	case "upPsi":
		return "Psi"
	case "endofproof":
		return "QED"
	case "eop":
		return "QED"
	case "Rcommaaccent":
		return "Rcedilla"
	case "Rsmallcap":
		return "Rdblgrave"
	case "Rblackletter":
		return "Rfraktur"
	case "Re":
		return "Rfraktur"
	case "Rfractur":
		return "Rfraktur"
	case "upRho":
		return "Rho"
	case "srthook":
		return "Rinvertedbreve"
	case "linevertdblnosp":
		return "Rsmallinverted"
	case "Germandbls":
		return "S"
	case "SS":
		return "S"
	case "250c":
		return "SF010000"
	case "253c":
		return "SF050000"
	case "252c":
		return "SF060000"
	case "251c":
		return "SF080000"
	case "255d":
		return "SF260000"
	case "255c":
		return "SF270000"
	case "255b":
		return "SF280000"
	case "255e":
		return "SF360000"
	case "255f":
		return "SF370000"
	case "255a":
		return "SF380000"
	case "256c":
		return "SF440000"
	case "256b":
		return "SF530000"
	case "256a":
		return "SF540000"
	case "EnGe":
		return "Sampigreek"
	case "upSampi":
		return "Sampigreek"
	case "bbar":
		return "Scaron"
	case "circleS":
		return "Scircle"
	case "trthook":
		return "Scommaaccent"
	case "arrowdblse":
		return "Searrow"
	case "Sha":
		return "Shacyrillic"
	case "afii10042":
		return "Shacyrillic"
	case "Pehook":
		return "Sheicoptic"
	case "Ustrtbar":
		return "Shimacoptic"
	case "upSigma":
		return "Sigma"
	case "Germandblssmall":
		return "Ssmall"
	case "SSsmall":
		return "Ssmall"
	case "Kabar":
		return "Stigmagreek"
	case "upStigma":
		return "Stigmagreek"
	case "arrowdblsw":
		return "Swarrow"
	case "upTau":
		return "Tau"
	case "Kcedilla1":
		return "Tcedilla1"
	case "Tcedilla":
		return "Tcommaaccent"
	case "upTheta":
		return "Theta"
	case "ahacek":
		return "Tretroflexhook"
	case "Tse":
		return "Tsecyrillic"
	case "afii10040":
		return "Tsecyrillic"
	case "Tshe":
		return "Tshecyrillic"
	case "afii10060":
		return "Tshecyrillic"
	case "Ucyril":
		return "Ucyrillic"
	case "afii10037":
		return "Ucyrillic"
	case "jhookdblbar":
		return "Udblgrave"
	case "aacutering":
		return "Udieresisgrave"
	case "Ihacek":
		return "Uhorn"
	case "Epsilon1":
		return "Uhungarumlaut"
	case "Udblacute":
		return "Uhungarumlaut"
	case "fscript":
		return "Uogonek"
	case "upUpsilon":
		return "Upsilon"
	case "Upsilonhooksymbol":
		return "Upsilon1"
	case "Zhertdes":
		return "Upsilon1"
	case "zhertdes":
		return "Upsilonacutehooksymbolgreek"
	case "Ohacek":
		return "Upsilonafrican"
	case "Zecedilla":
		return "Upsilondieresishooksymbolgreek"
	case "Eturn":
		return "Uring"
	case "Ucyrilbreve":
		return "Ushortcyrillic"
	case "afii10062":
		return "Ushortcyrillic"
	case "forceextr":
		return "VDash"
	case "ohacek":
		return "Vhook"
	case "Gamma1":
		return "Wcircumflex"
	case "Yat":
		return "Xi"
	case "upXi":
		return "Xi"
	case "Iota1":
		return "Ycircumflex"
	case "Uhacek":
		return "Yhook"
	case "Yi":
		return "Yicyrillic"
	case "afii10056":
		return "Yicyrillic"
	case "Nhook":
		return "Zcaron"
	case "Zdot":
		return "Zdotaccent"
	case "lambdabar":
		return "Zdotaccent"
	case "upZeta":
		return "Zeta"
	case "telephoneblack":
		return "a4"
	case "maltese":
		return "a9"
	case "maltesecross":
		return "a9"
	case "pointingindexrightwhite":
		return "a12"
	case "checkmark":
		return "a19"
	case "bigstar":
		return "a35"
	case "blackstar":
		return "a35"
	case "circledstar":
		return "a37"
	case "varstar":
		return "a49"
	case "dingasterisk":
		return "a56"
	case "circlesolid":
		return "a71"
	case "mdlgblkcircle":
		return "a71"
	case "bulletaltone":
		return "a71"
	case "blackcircle":
		return "a71"
	case "H18533":
		return "a71"
	case "filledbox":
		return "a73"
	case "squaresolid":
		return "a73"
	case "mdlgblksquare":
		return "a73"
	case "blacksquare":
		return "a73"
	case "trianglesolid":
		return "a76"
	case "blackuppointingtriangle":
		return "a76"
	case "bigblacktriangleup":
		return "a76"
	case "triagup":
		return "a76"
	case "blackdownpointingtriangle":
		return "a77"
	case "triangledownsld":
		return "a77"
	case "triagdn":
		return "a77"
	case "bigblacktriangledown":
		return "a77"
	case "diamondrhombsolid":
		return "a78"
	case "blackdiamond":
		return "a78"
	case "mdlgblkdiamond":
		return "a78"
	case "semicirclelertsld":
		return "a81"
	case "blackrighthalfcircle":
		return "a81"
	case "onecircle":
		return "a120"
	case "twocircle":
		return "a121"
	case "threecircle":
		return "a122"
	case "fourcircle":
		return "a123"
	case "fivecircle":
		return "a124"
	case "sixcircle":
		return "a125"
	case "sevencircle":
		return "a126"
	case "eightcircle":
		return "a127"
	case "ninecircle":
		return "a128"
	case "tencircle":
		return "a129"
	case "onecircleinversesansserif":
		return "a150"
	case "twocircleinversesansserif":
		return "a151"
	case "threecircleinversesansserif":
		return "a152"
	case "fourcircleinversesansserif":
		return "a153"
	case "fivecircleinversesansserif":
		return "a154"
	case "sixcircleinversesansserif":
		return "a155"
	case "sevencircleinversesansserif":
		return "a156"
	case "eightcircleinversesansserif":
		return "a157"
	case "ninecircleinversesansserif":
		return "a158"
	case "updownarrow":
		return "a164"
	case "arrowbothv":
		return "a164"
	case "arrowupdn":
		return "a164"
	case "draftingarrow":
		return "a166"
	case "arrowrightheavy":
		return "a169"
	case "Yoghhacek":
		return "acaron"
	case "acutecmb":
		return "acutecomb"
	case "arrowanticlockw":
		return "acwopencirclearrow"
	case "upslopeellipsis":
		return "adots"
	case "lrthook":
		return "aeacute"
	case "lefttoright":
		return "afii299"
	case "righttoleft":
		return "afii300"
	case "zerojoin":
		return "afii301"
	case "Acyril":
		return "afii10017"
	case "Acyrillic":
		return "afii10017"
	case "Be":
		return "afii10018"
	case "Becyrillic":
		return "afii10018"
	case "Vecyrillic":
		return "afii10019"
	case "Ve":
		return "afii10019"
	case "Zhe":
		return "afii10024"
	case "Zhecyrillic":
		return "afii10024"
	case "Zecyrillic":
		return "afii10025"
	case "Ze":
		return "afii10025"
	case "Iibreve":
		return "afii10027"
	case "Iishortcyrillic":
		return "afii10027"
	case "Kacyrillic":
		return "afii10028"
	case "Ka":
		return "afii10028"
	case "En":
		return "afii10031"
	case "Encyrillic":
		return "afii10031"
	case "Pecyril":
		return "afii10033"
	case "Pecyrillic":
		return "afii10033"
	case "Ercyrillic":
		return "afii10034"
	case "Er":
		return "afii10034"
	case "Es":
		return "afii10035"
	case "Escyrillic":
		return "afii10035"
	case "Tecyrillic":
		return "afii10036"
	case "Te":
		return "afii10036"
	case "Efcyrillic":
		return "afii10038"
	case "Ef":
		return "afii10038"
	case "Shchacyrillic":
		return "afii10043"
	case "Shcha":
		return "afii10043"
	case "Hard":
		return "afii10044"
	case "Hardsigncyrillic":
		return "afii10044"
	case "Yericyrillic":
		return "afii10045"
	case "Yeri":
		return "afii10045"
	case "Soft":
		return "afii10046"
	case "Softsigncyrillic":
		return "afii10046"
	case "Iu":
		return "afii10048"
	case "IUcyrillic":
		return "afii10048"
	case "Ia":
		return "afii10049"
	case "IAcyrillic":
		return "afii10049"
	case "Dje":
		return "afii10051"
	case "Djecyrillic":
		return "afii10051"
	case "Gje":
		return "afii10052"
	case "Gjecyrillic":
		return "afii10052"
	case "Dze":
		return "afii10054"
	case "Dzecyrillic":
		return "afii10054"
	case "Icyril":
		return "afii10055"
	case "Icyrillic":
		return "afii10055"
	case "Je":
		return "afii10057"
	case "Jecyrillic":
		return "afii10057"
	case "Nje":
		return "afii10059"
	case "Njecyrillic":
		return "afii10059"
	case "acyrillic":
		return "afii10065"
	case "acyril":
		return "afii10065"
	case "vecyrillic":
		return "afii10067"
	case "ve":
		return "afii10067"
	case "gecyrillic":
		return "afii10068"
	case "ge":
		return "afii10068"
	case "decyrillic":
		return "afii10069"
	case "de":
		return "afii10069"
	case "io":
		return "afii10071"
	case "iocyrillic":
		return "afii10071"
	case "ze":
		return "afii10073"
	case "zecyrillic":
		return "afii10073"
	case "iibreve":
		return "afii10075"
	case "iishortcyrillic":
		return "afii10075"
	case "en":
		return "afii10079"
	case "encyrillic":
		return "afii10079"
	case "te":
		return "afii10084"
	case "tecyrillic":
		return "afii10084"
	case "ucyrillic":
		return "afii10085"
	case "ucyril":
		return "afii10085"
	case "efcyrillic":
		return "afii10086"
	case "ef":
		return "afii10086"
	case "kha":
		return "afii10087"
	case "khacyrillic":
		return "afii10087"
	case "shacyrillic":
		return "afii10090"
	case "sha":
		return "afii10090"
	case "shchacyrillic":
		return "afii10091"
	case "shcha":
		return "afii10091"
	case "iu":
		return "afii10096"
	case "iucyrillic":
		return "afii10096"
	case "iacyrillic":
		return "afii10097"
	case "ia":
		return "afii10097"
	case "dzecyrillic":
		return "afii10102"
	case "dze":
		return "afii10102"
	case "icyrillic":
		return "afii10103"
	case "icyril":
		return "afii10103"
	case "je":
		return "afii10105"
	case "jecyrillic":
		return "afii10105"
	case "njecyrillic":
		return "afii10107"
	case "nje":
		return "afii10107"
	case "kjecyrillic":
		return "afii10109"
	case "kje":
		return "afii10109"
	case "ushortcyrillic":
		return "afii10110"
	case "ucyrilbreve":
		return "afii10110"
	case "Yatcyrillic":
		return "afii10146"
	case "Fitacyrillic":
		return "afii10147"
	case "Izhitsacyrillic":
		return "afii10148"
	case "fitacyrillic":
		return "afii10195"
	case "izhitsacyrillic":
		return "afii10196"
	case "afii10190":
		return "afii10196"
	case "arabiccomma":
		return "afii57388"
	case "commaarabic":
		return "afii57388"
	case "threearabic":
		return "afii57395"
	case "threehackarabic":
		return "afii57395"
	case "arabicindicdigitthree":
		return "afii57395"
	case "sixhackarabic":
		return "afii57398"
	case "arabicindicdigitsix":
		return "afii57398"
	case "sixarabic":
		return "afii57398"
	case "sevenhackarabic":
		return "afii57399"
	case "arabicindicdigitseven":
		return "afii57399"
	case "sevenarabic":
		return "afii57399"
	case "arabicsemicolon":
		return "afii57403"
	case "semicolonarabic":
		return "afii57403"
	case "questionarabic":
		return "afii57407"
	case "arabicquestionmark":
		return "afii57407"
	case "alefmaddaabovearabic":
		return "afii57410"
	case "alefwithmaddaabove":
		return "afii57410"
	case "alefhamzaabovearabic":
		return "afii57411"
	case "alefwithhamzaabove":
		return "afii57411"
	case "wawwithhamzaabove":
		return "afii57412"
	case "wawhamzaabovearabic":
		return "afii57412"
	case "teh":
		return "afii57418"
	case "teharabic":
		return "afii57418"
	case "hah":
		return "afii57421"
	case "haharabic":
		return "afii57421"
	case "khaharabic":
		return "afii57422"
	case "khah":
		return "afii57422"
	case "dalarabic":
		return "afii57423"
	case "dal":
		return "afii57423"
	case "seenarabic":
		return "afii57427"
	case "seen":
		return "afii57427"
	case "sheenarabic":
		return "afii57428"
	case "sheen":
		return "afii57428"
	case "sadarabic":
		return "afii57429"
	case "sad":
		return "afii57429"
	case "dad":
		return "afii57430"
	case "dadarabic":
		return "afii57430"
	case "ainarabic":
		return "afii57433"
	case "ain":
		return "afii57433"
	case "feharabic":
		return "afii57441"
	case "feh":
		return "afii57441"
	case "qaf":
		return "afii57442"
	case "qafarabic":
		return "afii57442"
	case "arabickaf":
		return "afii57443"
	case "kafarabic":
		return "afii57443"
	case "lam":
		return "afii57444"
	case "lamarabic":
		return "afii57444"
	case "meem":
		return "afii57445"
	case "meemarabic":
		return "afii57445"
	case "fathatanarabic":
		return "afii57451"
	case "fathatan":
		return "afii57451"
	case "dammatan":
		return "afii57452"
	case "dammatanarabic":
		return "afii57452"
	case "dammatanaltonearabic":
		return "afii57452"
	case "kasraarabic":
		return "afii57456"
	case "kasra":
		return "afii57456"
	case "jeh":
		return "afii57508"
	case "jeharabic":
		return "afii57508"
	case "tteharabic":
		return "afii57511"
	case "ddalarabic":
		return "afii57512"
	case "noonghunnaarabic":
		return "afii57514"
	case "arabicae":
		return "afii57534"
	case "sheqel":
		return "afii57636"
	case "sheqelhebrew":
		return "afii57636"
	case "newsheqelsign":
		return "afii57636"
	case "newsheqel":
		return "afii57636"
	case "maqaf":
		return "afii57645"
	case "maqafhebrew":
		return "afii57645"
	case "gimelhebrew":
		return "afii57666"
	case "gimel":
		return "afii57666"
	case "hehebrew":
		return "afii57668"
	case "he":
		return "afii57668"
	case "zayin":
		return "afii57670"
	case "zayinhebrew":
		return "afii57670"
	case "hethebrew":
		return "afii57671"
	case "het":
		return "afii57671"
	case "yodhebrew":
		return "afii57673"
	case "yod":
		return "afii57673"
	case "finalkafshevahebrew":
		return "afii57674"
	case "finalkaf":
		return "afii57674"
	case "finalkafsheva":
		return "afii57674"
	case "finalkafqamatshebrew":
		return "afii57674"
	case "finalkafhebrew":
		return "afii57674"
	case "finalkafqamats":
		return "afii57674"
	case "kaffinal":
		return "afii57674"
	case "finalnunhebrew":
		return "afii57679"
	case "nunfinal":
		return "afii57679"
	case "finalnun":
		return "afii57679"
	case "pehebrew":
		return "afii57684"
	case "pe":
		return "afii57684"
	case "tsadi":
		return "afii57686"
	case "tsadihebrew":
		return "afii57686"
	case "shinwithsindot":
		return "afii57695"
	case "shinsindothebrew":
		return "afii57695"
	case "shinsindot":
		return "afii57695"
	case "vavvavhebrew":
		return "afii57716"
	case "vavdbl":
		return "afii57716"
	case "vavyodhebrew":
		return "afii57717"
	case "vavyod":
		return "afii57717"
	case "qamatsquarterhebrew":
		return "afii57797"
	case "qamatsqatanquarterhebrew":
		return "afii57797"
	case "qamats1a":
		return "afii57797"
	case "qamatshebrew":
		return "afii57797"
	case "qamatsqatannarrowhebrew":
		return "afii57797"
	case "unitseparator":
		return "afii57797"
	case "qamatswidehebrew":
		return "afii57797"
	case "qamats":
		return "afii57797"
	case "qamats27":
		return "afii57797"
	case "qamatsqatanhebrew":
		return "afii57797"
	case "qamatsqatanwidehebrew":
		return "afii57797"
	case "qamatsde":
		return "afii57797"
	case "qamats1c":
		return "afii57797"
	case "qamats29":
		return "afii57797"
	case "qamatsnarrowhebrew":
		return "afii57797"
	case "qamats10":
		return "afii57797"
	case "qamats33":
		return "afii57797"
	case "sheva2e":
		return "afii57799"
	case "shevaquarterhebrew":
		return "afii57799"
	case "sheva15":
		return "afii57799"
	case "sheva115":
		return "afii57799"
	case "shevahebrew":
		return "afii57799"
	case "sheva":
		return "afii57799"
	case "endtransblock":
		return "afii57799"
	case "sheva22":
		return "afii57799"
	case "shevawidehebrew":
		return "afii57799"
	case "shevanarrowhebrew":
		return "afii57799"
	case "sindothebrew":
		return "afii57803"
	case "sindot":
		return "afii57803"
	case "rafehebrew":
		return "afii57841"
	case "rafe":
		return "afii57841"
	case "paseq":
		return "afii57842"
	case "paseqhebrew":
		return "afii57842"
	case "lscript":
		return "afii61289"
	case "lsquare":
		return "afii61289"
	case "liter":
		return "afii61289"
	case "ell":
		return "afii61289"
	case "pdf":
		return "afii61573"
	case "lro":
		return "afii61574"
	case "rlo":
		return "afii61575"
	case "zerowidthnonjoiner":
		return "afii61664"
	case "cwm":
		return "afii61664"
	case "zeronojoin":
		return "afii61664"
	case "compwordmark":
		return "afii61664"
	case "arabicfivepointedstar":
		return "afii63167"
	case "asteriskaltonearabic":
		return "afii63167"
	case "asteriskarabic":
		return "afii63167"
	case "commareversedmod":
		return "afii64937"
	case "numeralgreek":
		return "afii64937"
	case "ainfinal":
		return "ainfinalarabic"
	case "aininitial":
		return "aininitialarabic"
	case "ainmedial":
		return "ainmedialarabic"
	case "nrthook":
		return "ainvertedbreve"
	case "afii57664":
		return "alef"
	case "alefhebrew":
		return "alef"
	case "afii57415":
		return "alefarabic"
	case "arabicalef":
		return "alefarabic"
	case "alefwithmapiq":
		return "alefdageshhebrew"
	case "aleffinal":
		return "aleffinalarabic"
	case "alefwithhamzaabovefinal":
		return "alefhamzaabovefinalarabic"
	case "afii57413":
		return "alefhamzabelowarabic"
	case "alefwithhamzabelow":
		return "alefhamzabelowarabic"
	case "alefwithhamzabelowfinal":
		return "alefhamzabelowfinalarabic"
	case "aleflamed":
		return "aleflamedhebrew"
	case "alefwithmaddaabovefinal":
		return "alefmaddaabovefinalarabic"
	case "afii57449":
		return "alefmaksuraarabic"
	case "alefmaksura":
		return "alefmaksuraarabic"
	case "alefmaksurafinal":
		return "alefmaksurafinalarabic"
	case "yehmedial":
		return "alefmaksuramedialarabic"
	case "yehmedialarabic":
		return "alefmaksuramedialarabic"
	case "alefwithpatah":
		return "alefpatahhebrew"
	case "alefwithqamats":
		return "alefqamatshebrew"
	case "alephmath":
		return "aleph"
	case "backcong":
		return "allequal"
	case "upalpha":
		return "alpha"
	case "c158":
		return "amacron"
	case "langle":
		return "angbracketleft"
	case "rangle":
		return "angbracketright"
	case "afii59770":
		return "angkhankhuthai"
	case "angbracketleftBig":
		return "angleleft"
	case "angbracketleftBigg":
		return "angleleft"
	case "angbracketleftbig":
		return "angleleft"
	case "angbracketleftbigg":
		return "angleleft"
	case "angbracketrightBig":
		return "angleright"
	case "angbracketrightBigg":
		return "angleright"
	case "angbracketrightbig":
		return "angleright"
	case "angbracketrightbigg":
		return "angleright"
	case "Angstrom":
		return "angstrom"
	case "acwgapcirclearrow":
		return "anticlockwise"
	case "afii57929":
		return "apostrophemod"
	case "approachlimit":
		return "approaches"
	case "doteq":
		return "approaches"
	case "almostequal":
		return "approxequal"
	case "approx":
		return "approxequal"
	case "equaldotleftright":
		return "approxequalorimage"
	case "fallingdotseq":
		return "approxequalorimage"
	case "tildetrpl":
		return "approxident"
	case "almostorequal":
		return "approxorequal"
	case "approxeq":
		return "approxorequal"
	case "profline":
		return "arc"
	case "corresponds":
		return "arceq"
	case "arrowsemanticlockw":
		return "archleftdown"
	case "curvearrowleft":
		return "archleftdown"
	case "arrowsemclockw":
		return "archrightdown"
	case "curvearrowright":
		return "archrightdown"
	case "lmidtilde":
		return "aringacute"
	case "a163":
		return "arrowboth"
	case "leftrightarrow":
		return "arrowboth"
	case "downdasharrow":
		return "arrowdashdown"
	case "leftdasharrow":
		return "arrowdashleft"
	case "rightdasharrow":
		return "arrowdashright"
	case "updasharrow":
		return "arrowdashup"
	case "Leftrightarrow":
		return "arrowdblboth"
	case "arrowdbllongboth":
		return "arrowdblboth"
	case "dblarrowleft":
		return "arrowdblboth"
	case "Updownarrow":
		return "arrowdblbothv"
	case "arrowdbllongbothv":
		return "arrowdblbothv"
	case "Downarrow":
		return "arrowdbldown"
	case "Leftarrow":
		return "arrowdblleft"
	case "arrowleftdbl":
		return "arrowdblleft"
	case "Rightarrow":
		return "arrowdblright"
	case "dblarrowright":
		return "arrowdblright"
	case "Uparrow":
		return "arrowdblup"
	case "downarrow":
		return "arrowdown"
	case "swarrow":
		return "arrowdownleft"
	case "searrow":
		return "arrowdownright"
	case "arrowopendown":
		return "arrowdownwhite"
	case "downwhitearrow":
		return "arrowdownwhite"
	case "iotasub":
		return "arrowheadrightmod"
	case "hookrightarrow":
		return "arrowhookleft"
	case "hookleftarrow":
		return "arrowhookright"
	case "leftarrow":
		return "arrowleft"
	case "leftharpoondown":
		return "arrowleftbothalf"
	case "arrowdblleftnot":
		return "arrowleftdblstroke"
	case "nLeftarrow":
		return "arrowleftdblstroke"
	case "notdblarrowleft":
		return "arrowleftdblstroke"
	case "arrowparrleftright":
		return "arrowleftoverright"
	case "leftrightarrows":
		return "arrowleftoverright"
	case "arrowopenleft":
		return "arrowleftwhite"
	case "leftwhitearrow":
		return "arrowleftwhite"
	case "a161":
		return "arrowright"
	case "rightarrow":
		return "arrowright"
	case "rightharpoondown":
		return "arrowrightbothalf"
	case "arrowdblrightnot":
		return "arrowrightdblstroke"
	case "nRightarrow":
		return "arrowrightdblstroke"
	case "notdblarrowright":
		return "arrowrightdblstroke"
	case "arrowparrrightleft":
		return "arrowrightoverleft"
	case "rightleftarrows":
		return "arrowrightoverleft"
	case "arrowopenright":
		return "arrowrightwhite"
	case "rightwhitearrow":
		return "arrowrightwhite"
	case "barleftarrow":
		return "arrowtableft"
	case "rightarrowbar":
		return "arrowtabright"
	case "leftarrowtail":
		return "arrowtailleft"
	case "rightarrowtail":
		return "arrowtailright"
	case "Lleftarrow":
		return "arrowtripleleft"
	case "Rrightarrow":
		return "arrowtripleright"
	case "uparrow":
		return "arrowup"
	case "arrowupdnbse":
		return "arrowupdownbase"
	case "updownarrowbar":
		return "arrowupdownbase"
	case "nwarrow":
		return "arrowupleft"
	case "dblarrowupdown":
		return "arrowupleftofdown"
	case "updownarrows":
		return "arrowupleftofdown"
	case "nearrow":
		return "arrowupright"
	case "arrowopenup":
		return "arrowupwhite"
	case "upwhitearrow":
		return "arrowupwhite"
	case "linevert":
		return "ascript"
	case "macron1":
		return "ascriptturned"
	case "overscore1":
		return "ascriptturned"
	case "assertion":
		return "assert"
	case "ast":
		return "asteriskmath"
	case "asteriskcentered":
		return "asteriskmath"
	case "approxequalalt":
		return "asymptoticallyequal"
	case "asymptequal":
		return "asymptoticallyequal"
	case "simeq":
		return "asymptoticallyequal"
	case "similarequal":
		return "asymptoticallyequal"
	case "atsign":
		return "at"
	case "alternativeayin":
		return "ayinaltonehebrew"
	case "afii57682":
		return "ayinhebrew"
	case "ayin":
		return "ayinhebrew"
	case "primedblrev":
		return "backdprime"
	case "primedblrev1":
		return "backdprime"
	case "secondrev":
		return "backdprime"
	case "primetriplerev":
		return "backtrprime"
	case "primetriplerev1":
		return "backtrprime"
	case "afii59743":
		return "bahtthai"
	case "vert":
		return "bar"
	case "verticalbar":
		return "bar"
	case "tableftright":
		return "barleftarrowrightarrowba"
	case "home":
		return "barovernorthwestarrow"
	case "nor":
		return "barvee"
	case "afii10066":
		return "becyrillic"
	case "be":
		return "becyrillic"
	case "afii57416":
		return "beharabic"
	case "beh":
		return "beharabic"
	case "behfinal":
		return "behfinalarabic"
	case "behinitial":
		return "behinitialarabic"
	case "behmedial":
		return "behmedialarabic"
	case "behwithmeeminitial":
		return "behmeeminitialarabic"
	case "behwithmeemisolated":
		return "behmeemisolatedarabic"
	case "behwithnoonfinal":
		return "behnoonfinalarabic"
	case "upbeta":
		return "beta"
	case "Gehook":
		return "betasymbolgreek"
	case "upvarbeta":
		return "betasymbolgreek"
	case "betdagesh":
		return "betdageshhebrew"
	case "betwithdagesh":
		return "betdageshhebrew"
	case "bethmath":
		return "beth"
	case "afii57665":
		return "bethebrew"
	case "bet":
		return "bethebrew"
	case "betwithrafe":
		return "betrafehebrew"
	case "acute1":
		return "bhook"
	case "narylogicalor":
		return "bigvee"
	case "narylogicaland":
		return "bigwedge"
	case "ringsubnosp":
		return "bilabialclick"
	case "circlenwopen":
		return "blackcircleulquadwhite"
	case "semicircleleftsld":
		return "blacklefthalfcircle"
	case "blackpointerleft":
		return "blackleftpointingpointer"
	case "triaglf":
		return "blackleftpointingpointer"
	case "blacktriangleleft":
		return "blackleftpointingtriangle"
	case "triangleleftsld1":
		return "blackleftpointingtriangle"
	case "llblacktriangle":
		return "blacklowerlefttriangle"
	case "triangleswsld":
		return "blacklowerlefttriangle"
	case "lrblacktriangle":
		return "blacklowerrighttriangle"
	case "trianglesesld":
		return "blacklowerrighttriangle"
	case "filledrect":
		return "blackrectangle"
	case "hrectangleblack":
		return "blackrectangle"
	case "blackpointerright":
		return "blackrightpointingpointer"
	case "triagrt":
		return "blackrightpointingpointer"
	case "blacktriangleright":
		return "blackrightpointingtriangle"
	case "trianglerightsld1":
		return "blackrightpointingtriangle"
	case "H18543":
		return "blacksmallsquare"
	case "smallboxfilled":
		return "blacksmallsquare"
	case "smblksquare":
		return "blacksmallsquare"
	case "blacksmiley":
		return "blacksmilingface"
	case "invsmileface":
		return "blacksmilingface"
	case "smalltriangleinvsld":
		return "blacktriangledown"
	case "tranglenwsld":
		return "blackupperlefttriangle"
	case "ulblacktriangle":
		return "blackupperlefttriangle"
	case "trianglenesld":
		return "blackupperrighttriangle"
	case "urblacktriangle":
		return "blackupperrighttriangle"
	case "blacktriangle":
		return "blackuppointingsmalltriangle"
	case "smalltrianglesld":
		return "blackuppointingsmalltriangle"
	case "visiblespace":
		return "blank"
	case "visualspace":
		return "blank"
	case "blockfull":
		return "block"
	case "afii59706":
		return "bobaimaithai"
	case "bottomarc":
		return "botsemicircle"
	case "squarevertbisect":
		return "boxbar"
	case "braceleftBig":
		return "braceleft"
	case "braceleftBigg":
		return "braceleft"
	case "braceleftbig":
		return "braceleft"
	case "braceleftbigg":
		return "braceleft"
	case "lbrace":
		return "braceleft"
	case "bracehtipdownleft":
		return "braceleftvertical"
	case "bracehtipdownright":
		return "braceleftvertical"
	case "bracerightBig":
		return "braceright"
	case "bracerightBigg":
		return "braceright"
	case "bracerightbig":
		return "braceright"
	case "bracerightbigg":
		return "braceright"
	case "rbrace":
		return "braceright"
	case "appleopen":
		return "bracerightbt"
	case "enter":
		return "bracerightmid"
	case "carriagereturnleft":
		return "bracerighttp"
	case "bracehtipupleft":
		return "bracerightvertical"
	case "bracehtipupright":
		return "bracerightvertical"
	case "bracketleftBig":
		return "bracketleft"
	case "bracketleftBigg":
		return "bracketleft"
	case "bracketleftbig":
		return "bracketleft"
	case "bracketleftbigg":
		return "bracketleft"
	case "lbrack":
		return "bracketleft"
	case "bracketrightBig":
		return "bracketright"
	case "bracketrightBigg":
		return "bracketright"
	case "bracketrightbig":
		return "bracketright"
	case "bracketrightbigg":
		return "bracketright"
	case "rbrack":
		return "bracketright"
	case "contextmenu":
		return "bracketrightbt"
	case "power":
		return "bracketrighttp"
	case "rho1":
		return "bridgeinvertedbelowcmb"
	case "smblkcircle":
		return "bullet"
	case "bulletmath":
		return "bulletoperator"
	case "productdot":
		return "bulletoperator"
	case "vysmblkcircle":
		return "bulletoperator"
	case "bullseye1":
		return "bullseye"
	case "ct":
		return "c"
	case "overstore":
		return "c143"
	case "hmacron":
		return "cacute"
	case "candra":
		return "candrabinducmb"
	case "whitearrowupfrombar":
		return "capslock"
	case "afii61248":
		return "careof"
	case "caret":
		return "caretinsert"
	case "check":
		return "caroncmb"
	case "carriagerreturn":
		return "carriagereturn"
	case "linevertsub":
		return "ccurl"
	case "cdotaccent":
		return "cdot"
	case "Koppa":
		return "cedillacmb"
	case "ceilingleftBig":
		return "ceilingleft"
	case "ceilingleftBigg":
		return "ceilingleft"
	case "ceilingleftbig":
		return "ceilingleft"
	case "ceilingleftbigg":
		return "ceilingleft"
	case "lceil":
		return "ceilingleft"
	case "ceilingrightBig":
		return "ceilingright"
	case "ceilingrightBigg":
		return "ceilingright"
	case "ceilingrightbig":
		return "ceilingright"
	case "ceilingrightbigg":
		return "ceilingright"
	case "rceil":
		return "ceilingright"
	case "celsius":
		return "centigrade"
	case "degreecentigrade":
		return "centigrade"
	case "CL":
		return "centreline"
	case "afii10089":
		return "checyrillic"
	case "che":
		return "checyrillic"
	case "upchi":
		return "chi"
	case "afii59690":
		return "chochangthai"
	case "afii59688":
		return "chochanthai"
	case "afii59689":
		return "chochingthai"
	case "afii59692":
		return "chochoethai"
	case "ringequal":
		return "circeq"
	case "circledast":
		return "circleasterisk"
	case "circlebottomsld":
		return "circlebottomhalfblack"
	case "enclosecircle":
		return "circlecopyrt"
	case "circleminus1":
		return "circleddash"
	case "circledequal":
		return "circleequal"
	case "circlemultiplydisplay":
		return "circlemultiply"
	case "circlemultiplytext":
		return "circlemultiply"
	case "otimes":
		return "circlemultiply"
	case "timescircle":
		return "circlemultiply"
	case "circledot":
		return "circleot"
	case "circledotdisplay":
		return "circleot"
	case "circledottext":
		return "circleot"
	case "odot":
		return "circleot"
	case "circleplusdisplay":
		return "circleplus"
	case "circleplustext":
		return "circleplus"
	case "oplus":
		return "circleplus"
	case "pluscircle":
		return "circleplus"
	case "circledcirc":
		return "circlering"
	case "circletopsld":
		return "circletophalfblack"
	case "circlenesld":
		return "circleurquadblack"
	case "circleverthatch":
		return "circlevertfill"
	case "circlelefthalfblack":
		return "circlewithlefthalfblack"
	case "circleleftsld":
		return "circlewithlefthalfblack"
	case "circlerighthalfblack":
		return "circlewithrighthalfblack"
	case "circlerightsld":
		return "circlewithrighthalfblack"
	case "hat":
		return "circumflexcmb"
	case "hatwide":
		return "circumflexcmb"
	case "hatwider":
		return "circumflexcmb"
	case "hatwidest":
		return "circumflexcmb"
	case "cwgapcirclearrow":
		return "clockwise"
	case "a112":
		return "club"
	case "clubsuit":
		return "club"
	case "clubsuitblack":
		return "club"
	case "varclubsuit":
		return "clubsuitwhite"
	case "arrowsoutheast":
		return "coarmenian"
	case "mathcolon":
		return "colon"
	case "colonequal":
		return "coloneq"
	case "Colonmonetary":
		return "colonmonetary"
	case "coloncur":
		return "colonmonetary"
	case "coloncurrency":
		return "colonmonetary"
	case "colonsign":
		return "colonmonetary"
	case "iotadiaeresis":
		return "commaabovecmb"
	case "ocommatopright":
		return "commaaboverightcmb"
	case "upsilondiaeresis":
		return "commareversedabovecmb"
	case "oturnedcomma":
		return "commaturnedabovecmb"
	case "approximatelyequal":
		return "congruent"
	case "cong":
		return "congruent"
	case "contintegral":
		return "contourintegral"
	case "contintegraldisplay":
		return "contourintegral"
	case "contintegraltext":
		return "contourintegral"
	case "oint":
		return "contourintegral"
	case "ACK":
		return "controlACK"
	case "BEL":
		return "controlBEL"
	case "BS":
		return "controlBS"
	case "CAN":
		return "controlCAN"
	case "CR":
		return "controlCR"
	case "nonmarkingreturn":
		return "controlCR"
	case "XON":
		return "controlDC1"
	case "DC1":
		return "controlDC1"
	case "DC2":
		return "controlDC2"
	case "XOF":
		return "controlDC3"
	case "DC3":
		return "controlDC3"
	case "DC4":
		return "controlDC4"
	case "DEL":
		return "controlDEL"
	case "DC0":
		return "controlDLE"
	case "DLE":
		return "controlDLE"
	case "EM":
		return "controlEM"
	case "ENQ":
		return "controlENQ"
	case "EOT":
		return "controlEOT"
	case "ESC":
		return "controlESC"
	case "ETB":
		return "controlETB"
	case "ETX":
		return "controlETX"
	case "FF":
		return "controlFF"
	case "FS":
		return "controlFS"
	case "IFS":
		return "controlFS"
	case "GS":
		return "controlGS"
	case "IGS":
		return "controlGS"
	case "HT":
		return "controlHT"
	case "LF":
		return "controlLF"
	case "NAK":
		return "controlNAK"
	case ".null":
		return "controlNULL"
	case "NUL":
		return "controlNULL"
	case "IRS":
		return "controlRS"
	case "RS":
		return "controlRS"
	case "SI":
		return "controlSI"
	case "SO":
		return "controlSO"
	case "STX":
		return "controlSOT"
	case "SOH":
		return "controlSTX"
	case "EOF":
		return "controlSUB"
	case "SUB":
		return "controlSUB"
	case "SYN":
		return "controlSYN"
	case "IUS":
		return "controlUS"
	case "US":
		return "controlUS"
	case "VT":
		return "controlVT"
	case "amalg":
		return "coproduct"
	case "coprod":
		return "coproductdisplay"
	case "coproducttext":
		return "coproductdisplay"
	case "dotdblsubnosp":
		return "cstretched"
	case "multiplymultiset":
		return "cupdot"
	case "multiset":
		return "cupleftarrow"
	case "curland":
		return "curlyand"
	case "curlywedge":
		return "curlyand"
	case "uprise":
		return "curlyand"
	case "looparrowleft":
		return "curlyleft"
	case "curlor":
		return "curlyor"
	case "curlyvee":
		return "curlyor"
	case "downfall":
		return "curlyor"
	case "looparrowright":
		return "curlyright"
	case "arrowclockw":
		return "cwopencirclearrow"
	case "dadfinal":
		return "dadfinalarabic"
	case "dadinitial":
		return "dadinitialarabic"
	case "dadmedial":
		return "dadmedialarabic"
	case "afii57807":
		return "dagesh"
	case "dageshhebrew":
		return "dagesh"
	case "spaceopenbox":
		return "dagesh"
	case "ddagger":
		return "daggerdbl"
	case "daletdageshhebrew":
		return "daletdagesh"
	case "daletwithdagesh":
		return "daletdagesh"
	case "dalethmath":
		return "daleth"
	case "afii57667":
		return "daletqamatshebrew"
	case "dalet":
		return "daletqamatshebrew"
	case "dalethatafpatah":
		return "daletqamatshebrew"
	case "dalethatafpatahhebrew":
		return "daletqamatshebrew"
	case "dalethatafsegol":
		return "daletqamatshebrew"
	case "dalethatafsegolhebrew":
		return "daletqamatshebrew"
	case "dalethebrew":
		return "daletqamatshebrew"
	case "dalethiriq":
		return "daletqamatshebrew"
	case "dalethiriqhebrew":
		return "daletqamatshebrew"
	case "daletholam":
		return "daletqamatshebrew"
	case "daletholamhebrew":
		return "daletqamatshebrew"
	case "daletpatah":
		return "daletqamatshebrew"
	case "daletpatahhebrew":
		return "daletqamatshebrew"
	case "daletqamats":
		return "daletqamatshebrew"
	case "daletqubuts":
		return "daletqamatshebrew"
	case "daletqubutshebrew":
		return "daletqamatshebrew"
	case "daletsegol":
		return "daletqamatshebrew"
	case "daletsegolhebrew":
		return "daletqamatshebrew"
	case "daletsheva":
		return "daletqamatshebrew"
	case "daletshevahebrew":
		return "daletqamatshebrew"
	case "dalettsere":
		return "daletqamatshebrew"
	case "dalettserehebrew":
		return "daletqamatshebrew"
	case "dalfinal":
		return "dalfinalarabic"
	case "afii57455":
		return "dammaarabic"
	case "damma":
		return "dammaarabic"
	case "dammalowarabic":
		return "dammaarabic"
	case "dammahontatweel":
		return "dammamedial"
	case "dargahebrew":
		return "dargalefthebrew"
	case "shiftout":
		return "dargalefthebrew"
	case "excess":
		return "dashcolon"
	case "dblarrowdown":
		return "dblarrowdwn"
	case "downdownarrows":
		return "dblarrowdwn"
	case "twoheadleftarrow":
		return "dblarrowheadleft"
	case "twoheadrightarrow":
		return "dblarrowheadright"
	case "upuparrows":
		return "dblarrowup"
	case "lBrack":
		return "dblbracketleft"
	case "rBrack":
		return "dblbracketright"
	case "doubleintegral":
		return "dblintegral"
	case "iint":
		return "dblintegral"
	case "integraldbl":
		return "dblintegral"
	case "Vert":
		return "dblverticalbar"
	case "bardbl":
		return "dblverticalbar"
	case "verticalbardbl":
		return "dblverticalbar"
	case "vertlinedbl":
		return "dblverticalbar"
	case "downslopeellipsis":
		return "ddots"
	case "decimalseparatorarabic":
		return "decimalseparatorpersian"
	case "deltaequal":
		return "defines"
	case "triangleq":
		return "defines"
	case "kelvin":
		return "degreekelvin"
	case "devcon4":
		return "dehihebrew"
	case "khartdes":
		return "deicoptic"
	case "updelta":
		return "delta"
	case "macronsubnosp":
		return "dezh"
	case "gravesub":
		return "dhook"
	case "a111":
		return "diamond"
	case "diamondsolid":
		return "diamond"
	case "vardiamondsuit":
		return "diamond"
	case "smwhtdiamond":
		return "diamondmath"
	case "diamondsuit":
		return "diamondsuitwhite"
	case "ddot":
		return "dieresiscmb"
	case "dialytikatonos":
		return "dieresistonos"
	case "bumpeq":
		return "difference"
	case "c144":
		return "divide"
	case "div":
		return "divide"
	case "divideonmultiply":
		return "dividemultiply"
	case "divideontimes":
		return "dividemultiply"
	case "bar1":
		return "divides"
	case "mid":
		return "divides"
	case "vextendsingle":
		return "divides"
	case "divslash":
		return "divisionslash"
	case "slashmath":
		return "divisionslash"
	case "afii10099":
		return "djecyrillic"
	case "dje":
		return "djecyrillic"
	case "blockthreeqtrshaded":
		return "dkshade"
	case "shadedark":
		return "dkshade"
	case "dcroat":
		return "dmacron"
	case "dslash":
		return "dmacron"
	case "blocklowhalf":
		return "dnblock"
	case "afii59694":
		return "dochadathai"
	case "afii59700":
		return "dodekthai"
	case "escudo":
		return "dollar"
	case "mathdollar":
		return "dollar"
	case "milreis":
		return "dollar"
	case "iotadiaeresistonos":
		return "dotaccent"
	case "dot":
		return "dotaccentcmb"
	case "Stigma":
		return "dotbelowcomb"
	case "dotbelowcmb":
		return "dotbelowcomb"
	case "breveinvnosp":
		return "dotlessjstrokehook"
	case "geomproportion":
		return "dotsminusdots"
	case "proportiongeom":
		return "dotsminusdots"
	case "circledash":
		return "dottedcircle"
	case "xbsol":
		return "downslope"
	case "macronsub":
		return "dtail"
	case "gamma1":
		return "dz"
	case "tildesubnosp":
		return "dzaltone"
	case "Ghacek":
		return "dzcaron"
	case "underscorenosp":
		return "dzcurl"
	case "afii10193":
		return "dzhecyrillic"
	case "dzhe":
		return "dzhecyrillic"
	case "afii10101":
		return "ecyrillic"
	case "ecyril":
		return "ecyrillic"
	case "edot":
		return "edotaccent"
	case "afii57400":
		return "eighthackarabic"
	case "arabicindicdigiteight":
		return "eighthackarabic"
	case "eightarabic":
		return "eighthackarabic"
	case "musicalnotedbl":
		return "eighthnotebeamed"
	case "twonotes":
		return "eighthnotebeamed"
	case "eightsub":
		return "eightinferior"
	case "extendedarabicindicdigiteight":
		return "eightpersian"
	case "afii59768":
		return "eightthai"
	case "omegaclosed":
		return "einvertedbreve"
	case "afii10077":
		return "elcyrillic"
	case "el":
		return "elcyrillic"
	case "in":
		return "element"
	case "elipsis":
		return "ellipsis"
	case "unicodeellipsis":
		return "ellipsis"
	case "vdots":
		return "ellipsisvertical"
	case "vertellipsis":
		return "ellipsisvertical"
	case "afii10078":
		return "emcyrillic"
	case "em":
		return "emcyrillic"
	case "punctdash":
		return "emdash"
	case "varnothing":
		return "emptyset"
	case "rangedash":
		return "endash"
	case "ng":
		return "eng"
	case "ringrighthalfcenter":
		return "eopen"
	case "cedillanosp":
		return "eopenclosed"
	case "ringlefthalfsup":
		return "eopenreversed"
	case "tackdownmid":
		return "eopenreversedclosed"
	case "tackupmid":
		return "eopenreversedhook"
	case "upepsilon":
		return "epsilon"
	case "upvarepsilon":
		return "epsilon1"
	case "chevertbar":
		return "epsilon1"
	case "Hcyril":
		return "epsiloninv"
	case "upbackepsilon":
		return "epsiloninv"
	case "equalcolon":
		return "eqcolon"
	case "definequal":
		return "eqdef"
	case "equalgreater":
		return "eqgtr"
	case "equalless":
		return "eqless"
	case "curlyeqsucc":
		return "equalorfollows"
	case "equalfollows1":
		return "equalorfollows"
	case "eqslantgtr":
		return "equalorgreater"
	case "eqslantless":
		return "equalorless"
	case "curlyeqprec":
		return "equalorprecedes"
	case "equalprecedes1":
		return "equalorprecedes"
	case "eqsim":
		return "equalorsimilar"
	case "minustilde":
		return "equalorsimilar"
	case "equalinferior":
		return "equalsub"
	case "equiv":
		return "equivalence"
	case "asymp":
		return "equivasymptotic"
	case "afii10082":
		return "ercyrillic"
	case "er":
		return "ercyrillic"
	case "acutesub":
		return "ereversed"
	case "afii10095":
		return "ereversedcyrillic"
	case "ecyrilrev":
		return "ereversedcyrillic"
	case "afii10083":
		return "escyrillic"
	case "es":
		return "escyrillic"
	case "candrabindunosp":
		return "esh"
	case "apostrophesupnosp":
		return "eshcurl"
	case "commaturnsupnosp":
		return "eshsquatreversed"
	case "upeta":
		return "eta"
	case "Dbar":
		return "eth"
	case "Dmacron":
		return "eth"
	case "matheth":
		return "eth"
	case "arrowbothvbase":
		return "etnahtalefthebrew"
	case "etnahtafoukhhebrew":
		return "etnahtalefthebrew"
	case "etnahtafoukhlefthebrew":
		return "etnahtalefthebrew"
	case "etnahtahebrew":
		return "etnahtalefthebrew"
	case "Exclam":
		return "exclamdbl"
	case "exists":
		return "existential"
	case "thereexists":
		return "existential"
	case "plussubnosp":
		return "ezh"
	case "jdotlessbar":
		return "ezhcaron"
	case "minussubnosp":
		return "ezhcurl"
	case "Udieresishacek":
		return "ezhreversed"
	case "udieresishacek":
		return "ezhtail"
	case "degreefahrenheit":
		return "fahrenheit"
	case "degreefarenheit":
		return "fahrenheit"
	case "farenheit":
		return "fahrenheit"
	case "fathamedial":
		return "fathahontatweel"
	case "afii57454":
		return "fathalowarabic"
	case "fatha":
		return "fathalowarabic"
	case "fathaarabic":
		return "fathalowarabic"
	case "arrowwaveright":
		return "feharmenian"
	case "fehfinal":
		return "fehfinalarabic"
	case "fehinitial":
		return "fehinitialarabic"
	case "fehmedial":
		return "fehmedialarabic"
	case "ohook":
		return "feicoptic"
	case "venus":
		return "female"
	case "finalkafdagesh":
		return "finalkafdageshhebrew"
	case "finalkafwithdagesh":
		return "finalkafdageshhebrew"
	case "afii57677":
		return "finalmemhebrew"
	case "finalmem":
		return "finalmemhebrew"
	case "memfinal":
		return "finalmemhebrew"
	case "afii57683":
		return "finalpehebrew"
	case "finalpe":
		return "finalpehebrew"
	case "pefinal":
		return "finalpehebrew"
	case "afii57685":
		return "finaltsadi"
	case "finaltsadihebrew":
		return "finaltsadi"
	case "tsadifinal":
		return "finaltsadi"
	case "afii57397":
		return "fivearabic"
	case "arabicindicdigitfive":
		return "fivearabic"
	case "fivehackarabic":
		return "fivearabic"
	case "fivesub":
		return "fiveinferior"
	case "extendedarabicindicdigitfive":
		return "fivepersian"
	case "afii59765":
		return "fivethai"
	case "floorleftBig":
		return "floorleft"
	case "floorleftBigg":
		return "floorleft"
	case "floorleftbig":
		return "floorleft"
	case "floorleftbigg":
		return "floorleft"
	case "lfloor":
		return "floorleft"
	case "floorrightBig":
		return "floorright"
	case "floorrightBigg":
		return "floorright"
	case "floorrightbig":
		return "floorright"
	case "floorrightbigg":
		return "floorright"
	case "rfloor":
		return "floorright"
	case "Vcursive":
		return "florin"
	case "afii59711":
		return "fofanthai"
	case "afii59709":
		return "fofathai"
	case "succnapprox":
		return "follownotdbleqv"
	case "succneqq":
		return "follownotslnteql"
	case "followsnotequivlnt":
		return "followornoteqvlnt"
	case "succnsim":
		return "followornoteqvlnt"
	case "notfollowsoreql":
		return "followsequal"
	case "succeq":
		return "followsequal"
	case "followsequal1":
		return "followsorcurly"
	case "succcurlyeq":
		return "followsorcurly"
	case "followsequivlnt":
		return "followsorequal"
	case "succsim":
		return "followsorequal"
	case "afii59759":
		return "fongmanthai"
	case "Vdash":
		return "forces"
	case "force":
		return "forces"
	case "Vvdash":
		return "forcesbar"
	case "tacktrpl":
		return "forcesbar"
	case "pitchfork":
		return "fork"
	case "afii57396":
		return "fourarabic"
	case "arabicindicdigitfour":
		return "fourarabic"
	case "fourhackarabic":
		return "fourarabic"
	case "foursub":
		return "fourinferior"
	case "extendedarabicindicdigitfour":
		return "fourpersian"
	case "afii59764":
		return "fourthai"
	case "fracslash":
		return "fraction"
	case "fraction1":
		return "fraction"
	case "hturn":
		return "gacute"
	case "afii57509":
		return "gafarabic"
	case "gaf":
		return "gafarabic"
	case "gaffinal":
		return "gaffinalarabic"
	case "gafinitial":
		return "gafinitialarabic"
	case "gafmedial":
		return "gafmedialarabic"
	case "upgamma":
		return "gamma"
	case "ustrt":
		return "gangiacoptic"
	case "gcommaaccent":
		return "gcedilla"
	case "gdotaccent":
		return "gdot"
	case "Bumpeq":
		return "geomequivalent"
	case "Doteq":
		return "geometricallyequal"
	case "equalsdots":
		return "geometricallyequal"
	case "geomequal":
		return "geometricallyequal"
	case "endtext":
		return "gereshaccenthebrew"
	case "geresh":
		return "gereshhebrew"
	case "endtrans":
		return "gereshmuqdamhebrew"
	case "enquiry":
		return "gershayimaccenthebrew"
	case "gershayim":
		return "gershayimhebrew"
	case "verymuchgreater":
		return "ggg"
	case "afii57434":
		return "ghainarabic"
	case "ghain":
		return "ghainarabic"
	case "ghainfinal":
		return "ghainfinalarabic"
	case "ghaininitial":
		return "ghaininitialarabic"
	case "ghainmedial":
		return "ghainmedialarabic"
	case "afii10098":
		return "gheupturncyrillic"
	case "geupturn":
		return "gheupturncyrillic"
	case "gimelmath":
		return "gimel"
	case "gimeldagesh":
		return "gimeldageshhebrew"
	case "gimelwithdagesh":
		return "gimeldageshhebrew"
	case "afii10100":
		return "gjecyrillic"
	case "gje":
		return "gjecyrillic"
	case "hooksubpalatnosp":
		return "glottalstop"
	case "dotsubnosp":
		return "glottalstopinverted"
	case "hooksubretronosp":
		return "glottalstopreversed"
	case "brevesubnosp":
		return "glottalstopstroke"
	case "breveinvsubnosp":
		return "glottalstopstrokereversed"
	case "greaternotequivlnt":
		return "gnsim"
	case "nabla":
		return "gradient"
	case "gravecomb":
		return "gravecmb"
	case "diaeresistonos":
		return "gravelowmod"
	case "gtreqqless":
		return "greaterdbleqlless"
	case "gtrdot":
		return "greaterdot"
	case "geq":
		return "greaterequal"
	case "greaterequalless":
		return "greaterequalorless"
	case "greaterlessequal":
		return "greaterequalorless"
	case "gtreqless":
		return "greaterequalorless"
	case "gnapprox":
		return "greaternotdblequal"
	case "gneq":
		return "greaternotequal"
	case "gtrapprox":
		return "greaterorapproxeql"
	case "greaterequivlnt":
		return "greaterorequivalent"
	case "greaterorsimilar":
		return "greaterorequivalent"
	case "gtrsim":
		return "greaterorequivalent"
	case "gtrless":
		return "greaterorless"
	case "gneqq":
		return "greaterornotdbleql"
	case "greaterornotequal":
		return "greaterornotdbleql"
	case "geqq":
		return "greateroverequal"
	case "greaterdblequal":
		return "greateroverequal"
	case "notgreaterdblequal":
		return "greateroverequal"
	case "hehaltonearabic":
		return "haaltonearabic"
	case "hahfinal":
		return "hahfinalarabic"
	case "hahinitial":
		return "hahinitialarabic"
	case "hahmedial":
		return "hahmedialarabic"
	case "afii57409":
		return "hamzadammaarabic"
	case "hamza":
		return "hamzadammaarabic"
	case "hamzaarabic":
		return "hamzadammaarabic"
	case "hamzadammatanarabic":
		return "hamzadammaarabic"
	case "hamzafathaarabic":
		return "hamzadammaarabic"
	case "hamzafathatanarabic":
		return "hamzadammaarabic"
	case "hamzalowarabic":
		return "hamzadammaarabic"
	case "hamzalowkasraarabic":
		return "hamzadammaarabic"
	case "hamzalowkasratanarabic":
		return "hamzadammaarabic"
	case "hamzasukunarabic":
		return "hamzadammaarabic"
	case "afii10092":
		return "hardsigncyrillic"
	case "hard":
		return "hardsigncyrillic"
	case "downharpoonleft":
		return "harpoondownleft"
	case "downharpoonright":
		return "harpoondownright"
	case "arrowlefttophalf":
		return "harpoonleftbarbup"
	case "leftharpoonup":
		return "harpoonleftbarbup"
	case "rightleftharpoons":
		return "harpoonleftright"
	case "arrowrighttophalf":
		return "harpoonrightbarbup"
	case "rightharpoonup":
		return "harpoonrightbarbup"
	case "leftrightharpoons":
		return "harpoonrightleft"
	case "upharpoonleft":
		return "harpoonupleft"
	case "upharpoonright":
		return "harpoonupright"
	case "hatafpatahwidehebrew":
		return "hatafpatah16"
	case "hatafpatahquarterhebrew":
		return "hatafpatah16"
	case "hatafpatahhebrew":
		return "hatafpatah16"
	case "hatafpatahnarrowhebrew":
		return "hatafpatah16"
	case "hatafpatah2f":
		return "hatafpatah16"
	case "afii57800":
		return "hatafpatah16"
	case "endmedium":
		return "hatafpatah16"
	case "hatafpatah23":
		return "hatafpatah16"
	case "hatafpatah":
		return "hatafpatah16"
	case "hatafqamatshebrew":
		return "hatafqamats28"
	case "afii57802":
		return "hatafqamats28"
	case "substitute":
		return "hatafqamats28"
	case "hatafqamats34":
		return "hatafqamats28"
	case "hatafqamatswidehebrew":
		return "hatafqamats28"
	case "hatafqamatsnarrowhebrew":
		return "hatafqamats28"
	case "hatafqamatsquarterhebrew":
		return "hatafqamats28"
	case "hatafqamats1b":
		return "hatafqamats28"
	case "hatafqamats":
		return "hatafqamats28"
	case "endoffile":
		return "hatafqamats28"
	case "afii57801":
		return "hatafsegolwidehebrew"
	case "cancel":
		return "hatafsegolwidehebrew"
	case "hatafsegol":
		return "hatafsegolwidehebrew"
	case "hatafsegol17":
		return "hatafsegolwidehebrew"
	case "hatafsegol24":
		return "hatafsegolwidehebrew"
	case "hatafsegol30":
		return "hatafsegolwidehebrew"
	case "hatafsegolhebrew":
		return "hatafsegolwidehebrew"
	case "hatafsegolnarrowhebrew":
		return "hatafsegolwidehebrew"
	case "hatafsegolquarterhebrew":
		return "hatafsegolwidehebrew"
	case "a110":
		return "heart"
	case "heartsuitblack":
		return "heart"
	case "varheartsuit":
		return "heart"
	case "heartsuit":
		return "heartsuitwhite"
	case "hedagesh":
		return "hedageshhebrew"
	case "hewithmapiq":
		return "hedageshhebrew"
	case "afii57470":
		return "heharabic"
	case "heh":
		return "heharabic"
	case "hehfinal":
		return "hehfinalarabic"
	case "hehfinalalttwoarabic":
		return "hehfinalarabic"
	case "hehinitial":
		return "hehinitialarabic"
	case "hehmedial":
		return "hehmedialarabic"
	case "rhotichook":
		return "henghook"
	case "hermitconjmatrix":
		return "hermitmatrix"
	case "tildevertsupnosp":
		return "hhooksuperior"
	case "hiriq":
		return "hiriq14"
	case "hiriq2d":
		return "hiriq14"
	case "afii57793":
		return "hiriq14"
	case "hiriqhebrew":
		return "hiriq14"
	case "escape":
		return "hiriq14"
	case "hiriqnarrowhebrew":
		return "hiriq14"
	case "hiriqquarterhebrew":
		return "hiriq14"
	case "hiriq21":
		return "hiriq14"
	case "hiriqwidehebrew":
		return "hiriq14"
	case "afii59723":
		return "hohipthai"
	case "afii57806":
		return "holamquarterhebrew"
	case "holam":
		return "holamquarterhebrew"
	case "holam19":
		return "holamquarterhebrew"
	case "holam26":
		return "holamquarterhebrew"
	case "holam32":
		return "holamquarterhebrew"
	case "holamhebrew":
		return "holamquarterhebrew"
	case "holamnarrowhebrew":
		return "holamquarterhebrew"
	case "holamwidehebrew":
		return "holamquarterhebrew"
	case "spaceliteral":
		return "holamquarterhebrew"
	case "afii59726":
		return "honokhukthai"
	case "hookabovecomb":
		return "hookcmb"
	case "ovhook":
		return "hookcmb"
	case "tertdes":
		return "horicoptic"
	case "afii00208":
		return "horizontalbar"
	case "horizbar":
		return "horizontalbar"
	case "longdash":
		return "horizontalbar"
	case "quotedash":
		return "horizontalbar"
	case "rectangle":
		return "hrectangle"
	case "xsupnosp":
		return "hsuperior"
	case "SD190100":
		return "hturned"
	case "Zbar":
		return "hv"
	case "hyphen-minus":
		return "hyphen"
	case "hyphenchar":
		return "hyphen"
	case "hyphenminus":
		return "hyphen"
	case "hyphen1":
		return "hyphentwo"
	case "jhacek":
		return "icaron"
	case "rturn":
		return "idblgrave"
	case "dquoteright":
		return "idieresis"
	case "afii10070":
		return "iecyrillic"
	case "ie":
		return "iecyrillic"
	case "afii10074":
		return "iicyrillic"
	case "ii":
		return "iicyrillic"
	case "integraltrpl":
		return "iiint"
	case "tripleintegral":
		return "iiint"
	case "rturnhook":
		return "iinvertedbreve"
	case "rturnrthook":
		return "iinvertedbreve"
	case "auxiliaryoff":
		return "iluyhebrew"
	case "devcon3":
		return "iluyhebrew"
	case "image":
		return "imageof"
	case "equaldotrightleft":
		return "imageorapproximatelyequal"
	case "imageorapproxequal":
		return "imageorapproximatelyequal"
	case "risingdotseq":
		return "imageorapproximatelyequal"
	case "infty":
		return "infinity"
	case "clwintegral":
		return "intclockwise"
	case "backslashBig":
		return "integerdivide"
	case "backslashBigg":
		return "integerdivide"
	case "backslashbig":
		return "integerdivide"
	case "backslashbigg":
		return "integerdivide"
	case "backslashmath":
		return "integerdivide"
	case "smallsetminus":
		return "integerdivide"
	case "int":
		return "integral"
	case "integraldisplay":
		return "integral"
	case "integraltext":
		return "integral"
	case "intbottom":
		return "integralbt"
	case "integralbottom":
		return "integralbt"
	case "integraltop":
		return "integraltp"
	case "inttop":
		return "integraltp"
	case "cap":
		return "intersection"
	case "Cap":
		return "intersectiondbl"
	case "bigcap":
		return "intersectiondisplay"
	case "intersectiontext":
		return "intersectiondisplay"
	case "naryintersection":
		return "intersectiondisplay"
	case "sqcap":
		return "intersectionsq"
	case "bulletinverse":
		return "invbullet"
	case "inversebullet":
		return "invbullet"
	case "inversewhitecircle":
		return "invcircle"
	case "whitecircleinverse":
		return "invcircle"
	case "Sinvlazy":
		return "invlazys"
	case "lazysinv":
		return "invlazys"
	case "invsemicircledn":
		return "invwhitelowerhalfcircle"
	case "invsemicircleup":
		return "invwhiteupperhalfcircle"
	case "upiota":
		return "iota"
	case "gammasuper":
		return "iotalatin"
	case "highcomman":
		return "itilde"
	case "bridgesubnosp":
		return "jcrossedtail"
	case "afii57420":
		return "jeemarabic"
	case "jeem":
		return "jeemarabic"
	case "jeemfinal":
		return "jeemfinalarabic"
	case "jeeminitial":
		return "jeeminitialarabic"
	case "jeemmedial":
		return "jeemmedialarabic"
	case "jehfinal":
		return "jehfinalarabic"
	case "overscoredblnosp":
		return "jsuperior"
	case "afii10076":
		return "kacyrillic"
	case "ka":
		return "kacyrillic"
	case "afii57675":
		return "kaf"
	case "kafhebrew":
		return "kaf"
	case "kafdageshhebrew":
		return "kafdagesh"
	case "kafwithdagesh":
		return "kafdagesh"
	case "arabickaffinal":
		return "kaffinalarabic"
	case "kafinitial":
		return "kafinitialarabic"
	case "kafmedial":
		return "kafmedialarabic"
	case "kafwithrafe":
		return "kafrafehebrew"
	case "upkappa":
		return "kappa"
	case "TeTse":
		return "kappasymbolgreek"
	case "upvarkappa":
		return "kappasymbolgreek"
	case "afii57440":
		return "kashidaautonosidebearingarabic"
	case "kashidaautoarabic":
		return "kashidaautonosidebearingarabic"
	case "tatweel":
		return "kashidaautonosidebearingarabic"
	case "tatweelarabic":
		return "kashidaautonosidebearingarabic"
	case "kasrahontatweel":
		return "kasramedial"
	case "afii57453":
		return "kasratanarabic"
	case "kasratan":
		return "kasratanarabic"
	case "kcedilla":
		return "kcommaaccent"
	case "arrowrightnot":
		return "keharmenian"
	case "homothetic":
		return "kernelcontraction"
	case "khahfinal":
		return "khahfinalarabic"
	case "khahinitial":
		return "khahinitialarabic"
	case "khahmedial":
		return "khahmedialarabic"
	case "escedilla":
		return "kheicoptic"
	case "afii59682":
		return "khokhaithai"
	case "afii59685":
		return "khokhonthai"
	case "afii59683":
		return "khokhuatthai"
	case "afii59684":
		return "khokhwaithai"
	case "afii59771":
		return "khomutthai"
	case "yoghrev":
		return "khook"
	case "afii59686":
		return "khorakhangthai"
	case "afii59681":
		return "kokaithai"
	case "archdblsubnosp":
		return "kturned"
	case "afii59749":
		return "lakkhangyaothai"
	case "lamwithaleffinal":
		return "lamaleffinalarabic"
	case "lamwithalefhamzaabovefinal":
		return "lamalefhamzaabovefinalarabic"
	case "lamwithalefhamzaaboveisolatedd":
		return "lamalefhamzaaboveisolatedarabic"
	case "lamwithalefhamzabelowfinal":
		return "lamalefhamzabelowfinalarabic"
	case "lamwithalefhamzabelowisolated":
		return "lamalefhamzabelowisolatedarabic"
	case "lamwithalefisolated":
		return "lamalefisolatedarabic"
	case "lamwithalefmaddaabovefinal":
		return "lamalefmaddaabovefinalarabic"
	case "lamwithalefmaddaaboveisolatedd":
		return "lamalefmaddaaboveisolatedarabic"
	case "uplambda":
		return "lambda"
	case "2bar":
		return "lambdastroke"
	case "lameddageshhebrew":
		return "lameddagesh"
	case "lamedwithdagesh":
		return "lameddagesh"
	case "afii57676":
		return "lamedholamhebrew"
	case "lamed":
		return "lamedholamhebrew"
	case "lamedhebrew":
		return "lamedholamhebrew"
	case "lamedholam":
		return "lamedholamhebrew"
	case "lamedholamdagesh":
		return "lamedholamhebrew"
	case "lamedholamdageshhebrew":
		return "lamedholamhebrew"
	case "lamfinal":
		return "lamfinalarabic"
	case "lamwithhahinitial":
		return "lamhahinitialarabic"
	case "laminitial":
		return "laminitialarabic"
	case "lammeemjeeminitialarabic":
		return "laminitialarabic"
	case "lammeemkhahinitialarabic":
		return "laminitialarabic"
	case "lamwithjeeminitial":
		return "lamjeeminitialarabic"
	case "lamwithkhahinitial":
		return "lamkhahinitialarabic"
	case "allahisolated":
		return "lamlamhehisolatedarabic"
	case "lammedial":
		return "lammedialarabic"
	case "lamwithmeemwithhahinitial":
		return "lammeemhahinitialarabic"
	case "lamwithmeeminitial":
		return "lammeeminitialarabic"
	case "lgwhtcircle":
		return "largecircle"
	case "yoghtail":
		return "lbar"
	case "xsuper":
		return "lbelt"
	case "lcedilla":
		return "lcommaaccent"
	case "ldot":
		return "ldotaccent"
	case "droang":
		return "leftangleabovecmb"
	case "arrowsquiggleleft":
		return "leftsquigarrow"
	case "lesseqqgtr":
		return "lessdbleqlgreater"
	case "leq":
		return "lessequal"
	case "lesseqgtr":
		return "lessequalorgreater"
	case "lessequalgreater":
		return "lessequalorgreater"
	case "lnapprox":
		return "lessnotdblequal"
	case "lneq":
		return "lessnotequal"
	case "lessapprox":
		return "lessorapproxeql"
	case "leqslant":
		return "lessorequalslant"
	case "notlessorslnteql":
		return "lessorequalslant"
	case "lessequivlnt":
		return "lessorequivalent"
	case "lessorsimilar":
		return "lessorequivalent"
	case "lesssim":
		return "lessorequivalent"
	case "lessgtr":
		return "lessorgreater"
	case "lessornotdbleql":
		return "lessornotequal"
	case "lneqq":
		return "lessornotequal"
	case "leqq":
		return "lessoverequal"
	case "lessdblequal":
		return "lessoverequal"
	case "notlessdblequal":
		return "lessoverequal"
	case "toneextrahigh":
		return "lezh"
	case "blocklefthalf":
		return "lfblock"
	case "glottalrevsuper":
		return "lhookretroflex"
	case "arrowrightdown":
		return "linefeed"
	case "afii08941":
		return "lira"
	case "khacek":
		return "lj"
	case "afii10106":
		return "ljecyrillic"
	case "lje":
		return "ljecyrillic"
	case "swquadarc":
		return "llarc"
	case "verymuchless":
		return "lll"
	case "ssuper":
		return "lmiddletilde"
	case "lessnotequivlnt":
		return "lnsim"
	case "afii59724":
		return "lochulathai"
	case "logicalanddisplay":
		return "logicaland"
	case "logicalandtext":
		return "logicaland"
	case "wedge":
		return "logicaland"
	case "neg":
		return "logicalnot"
	case "logicalordisplay":
		return "logicalor"
	case "logicalortext":
		return "logicalor"
	case "vee":
		return "logicalor"
	case "afii59717":
		return "lolingthai"
	case "Obar":
		return "longs"
	case "longdbls":
		return "longs"
	case "longsh":
		return "longs"
	case "longsi":
		return "longs"
	case "longsl":
		return "longs"
	case "slong":
		return "longs"
	case "slongt":
		return "longst"
	case "mdlgwhtlozenge":
		return "lozenge"
	case "sequadarc":
		return "lrarc"
	case "afii59718":
		return "luthai"
	case "overscore":
		return "macron"
	case "underbar":
		return "macronbelowcmb"
	case "mahapakhlefthebrew":
		return "mahapakhhebrew"
	case "verttab":
		return "mahapakhhebrew"
	case "afii59755":
		return "maichattawathai"
	case "afii59752":
		return "maiekthai"
	case "afii59728":
		return "maihanakatthai"
	case "afii59751":
		return "maitaikhuthai"
	case "afii59753":
		return "maithothai"
	case "afii59754":
		return "maitrithai"
	case "afii59750":
		return "maiyamokthai"
	case "male":
		return "mars"
	case "synch":
		return "masoracirclehebrew"
	case "measurequal":
		return "measeq"
	case "rightanglearc":
		return "measuredrightangle"
	case "meemfinal":
		return "meemfinalarabic"
	case "meeminitial":
		return "meeminitialarabic"
	case "meemmedial":
		return "meemmedialarabic"
	case "meemwithmeeminitial":
		return "meemmeeminitialarabic"
	case "afii57678":
		return "mem"
	case "memhebrew":
		return "mem"
	case "memdagesh":
		return "memdageshhebrew"
	case "memwithdagesh":
		return "memdageshhebrew"
	case "formfeed":
		return "merkhahebrew"
	case "merkhalefthebrew":
		return "merkhahebrew"
	case "merkhakefulalefthebrew":
		return "merkhakefulahebrew"
	case "Cblackletter":
		return "mfrakC"
	case "Cfractur":
		return "mfrakC"
	case "Cfraktur":
		return "mfrakC"
	case "Hblackletter":
		return "mfrakH"
	case "Hfractur":
		return "mfrakH"
	case "Hfraktur":
		return "mfrakH"
	case "Zblackletter":
		return "mfrakZ"
	case "Zfractur":
		return "mfrakZ"
	case "Zfraktur":
		return "mfrakZ"
	case "tonelow":
		return "mhook"
	case "circleminus":
		return "minuscircle"
	case "ominus":
		return "minuscircle"
	case "minussub":
		return "minusinferior"
	case "mp":
		return "minusplus"
	case "prime":
		return "minute"
	case "prime1":
		return "minute"
	case "tonemid":
		return "mlonglegturned"
	case "truestate":
		return "models"
	case "afii59713":
		return "momathai"
	case "Bscript":
		return "mscrB"
	case "Escript":
		return "mscrE"
	case "Fscript":
		return "mscrF"
	case "Hscript":
		return "mscrH"
	case "Iscript":
		return "mscrI"
	case "Lscript":
		return "mscrL"
	case "Mscript":
		return "mscrM"
	case "Rscript":
		return "mscrR"
	case "escript":
		return "mscre"
	case "gscriptmath":
		return "mscrg"
	case "0script":
		return "mscro"
	case "oscript":
		return "mscro"
	case "tonehigh":
		return "mturned"
	case "mu1":
		return "mu"
	case "gg":
		return "muchgreater"
	case "greatermuch":
		return "muchgreater"
	case "lessmuch":
		return "muchless"
	case "upmu":
		return "mugreek"
	case "ltimes":
		return "multicloseleft"
	case "rtimes":
		return "multicloseright"
	case "leftthreetimes":
		return "multiopenleft"
	case "rightthreetimes":
		return "multiopenright"
	case "times":
		return "multiply"
	case "munahhebrew":
		return "munahlefthebrew"
	case "eighthnote":
		return "musicalnote"
	case "flat":
		return "musicflatsign"
	case "sharp":
		return "musicsharpsign"
	case "barwedge":
		return "nand"
	case "notalmostequal":
		return "napprox"
	case "notequivasymptotic":
		return "nasymp"
	case "hyphennobreak":
		return "nbhyphen"
	case "Tmacron":
		return "ncedilla"
	case "ncommaaccent":
		return "ncedilla"
	case "afii59687":
		return "ngonguthai"
	case "notgreaterequivlnt":
		return "ngtrsim"
	case "toneextralow":
		return "nhookleft"
	case "gravenosp":
		return "nhookretroflex"
	case "afii59757":
		return "nikhahitthai"
	case "afii57401":
		return "ninehackarabic"
	case "arabicindicdigitnine":
		return "ninehackarabic"
	case "ninearabic":
		return "ninehackarabic"
	case "ninesub":
		return "nineinferior"
	case "extendedarabicindicdigitnine":
		return "ninepersian"
	case "afii59769":
		return "ninethai"
	case "glottalstopbarinv":
		return "nlegrightlong"
	case "notlessgreater":
		return "nlessgtr"
	case "notlessequivlnt":
		return "nlesssim"
	case "nbspace":
		return "nonbreakingspace"
	case "afii59699":
		return "nonenthai"
	case "afii59705":
		return "nonuthai"
	case "afii57446":
		return "noonarabic"
	case "noon":
		return "noonarabic"
	case "noonfinal":
		return "noonfinalarabic"
	case "nooninitial":
		return "noonhehinitialarabic"
	case "nooninitialarabic":
		return "noonhehinitialarabic"
	case "noonwithjeeminitial":
		return "noonjeeminitialarabic"
	case "noonmedial":
		return "noonmedialarabic"
	case "noonwithmeeminitial":
		return "noonmeeminitialarabic"
	case "noonwithmeemisolated":
		return "noonmeemisolatedarabic"
	case "ncong":
		return "notapproxequal"
	case "nleftrightarrow":
		return "notarrowboth"
	case "nleftarrow":
		return "notarrowleft"
	case "nrightarrow":
		return "notarrowright"
	case "nmid":
		return "notbar"
	case "notdivides":
		return "notbar"
	case "nni":
		return "notcontains"
	case "notowner":
		return "notcontains"
	case "notsuchthat":
		return "notcontains"
	case "arrowdbllongbothnot":
		return "notdblarrowboth"
	case "nLeftrightarrow":
		return "notdblarrowboth"
	case "notelementof":
		return "notelement"
	case "notin":
		return "notelement"
	case "ne":
		return "notequal"
	case "nexists":
		return "notexistential"
	case "nVdash":
		return "notforces"
	case "notforce":
		return "notforces"
	case "nVDash":
		return "notforcesextra"
	case "notforceextr":
		return "notforcesextra"
	case "ngtr":
		return "notgreater"
	case "ngeq":
		return "notgreaternorequal"
	case "notgreaterequal":
		return "notgreaternorequal"
	case "notgreaterequal1":
		return "notgreaternorequal"
	case "ngtrless":
		return "notgreaternorless"
	case "notgreaterless":
		return "notgreaternorless"
	case "geqslant":
		return "notgreaterorslnteql"
	case "greaterorequalslant":
		return "notgreaterorslnteql"
	case "nequiv":
		return "notidentical"
	case "notequivalence":
		return "notidentical"
	case "nless":
		return "notless"
	case "nleq":
		return "notlessnorequal"
	case "notlessequal":
		return "notlessnorequal"
	case "notlessequal1":
		return "notlessnorequal"
	case "notbardbl":
		return "notparallel"
	case "nparallel":
		return "notparallel"
	case "notpreceeds":
		return "notprecedes"
	case "nprec":
		return "notprecedes"
	case "notsatisfy":
		return "notsatisfies"
	case "nvDash":
		return "notsatisfies"
	case "nsim":
		return "notsimilar"
	case "notpropersubset":
		return "notsubset"
	case "nsubset":
		return "notsubset"
	case "notreflexsubset":
		return "notsubseteql"
	case "nsubseteq":
		return "notsubseteql"
	case "notfollows":
		return "notsucceeds"
	case "nsucc":
		return "notsucceeds"
	case "notpropersuperset":
		return "notsuperset"
	case "nsupset":
		return "notsuperset"
	case "notreflexsuperset":
		return "notsuperseteql"
	case "nsupseteq":
		return "notsuperseteql"
	case "nottriangleleftequal":
		return "nottriangeqlleft"
	case "ntrianglelefteq":
		return "nottriangeqlleft"
	case "nottrianglerightequal":
		return "nottriangeqlright"
	case "ntrianglerighteq":
		return "nottriangeqlright"
	case "ntriangleleft":
		return "nottriangleleft"
	case "ntriangleright":
		return "nottriangleright"
	case "notturnstileleft":
		return "notturnstile"
	case "nvdash":
		return "notturnstile"
	case "preceedsnotequal":
		return "npreccurlyeq"
	case "notasymptequal":
		return "nsime"
	case "notsubsetsqequal":
		return "nsqsubseteq"
	case "notsupersetsqequal":
		return "nsqsupseteq"
	case "followsnotequal":
		return "nsucccurlyeq"
	case "dbar":
		return "ntilde"
	case "upnu":
		return "nu"
	case "octothorpe":
		return "numbersign"
	case "afii61352":
		return "numero"
	case "afii57680":
		return "nun"
	case "nunhebrew":
		return "nun"
	case "nundageshhebrew":
		return "nundagesh"
	case "nunwithdagesh":
		return "nundagesh"
	case "afii59725":
		return "oangthai"
	case "circumflexnosp":
		return "obarred"
	case "afii10080":
		return "ocyrillic"
	case "ocyril":
		return "ocyrillic"
	case "rrthook":
		return "odblgrave"
	case "arrowwaveleft":
		return "oharmenian"
	case "pipedbl":
		return "ohorn"
	case "odblacute":
		return "ohungarumlaut"
	case "exclam1":
		return "oi"
	case "volintegral":
		return "oiiint"
	case "volumeintegral":
		return "oiiint"
	case "surfaceintegral":
		return "oiint"
	case "surfintegral":
		return "oiint"
	case "cclwcontintegral":
		return "ointctrclockwise"
	case "rfishhookrev":
		return "oinvertedbreve"
	case "devcon2":
		return "olehebrew"
	case "upomega":
		return "omega"
	case "upvarpi":
		return "omega1"
	case "pisymbolgreek":
		return "omega1"
	case "Kartdes":
		return "omega1"
	case "macronnosp":
		return "omegalatinclosed"
	case "Gebar":
		return "omegatonos"
	case "upomicron":
		return "omicron"
	case "onedotlead":
		return "onedotenleader"
	case "onedotleader":
		return "onedotenleader"
	case "afii57393":
		return "onehackarabic"
	case "arabicindicdigitone":
		return "onehackarabic"
	case "onearabic":
		return "onehackarabic"
	case "onesub":
		return "oneinferior"
	case "extendedarabicindicdigitone":
		return "onepersian"
	case "afii59761":
		return "onethai"
	case "epsilon1revhook":
		return "oogonekmacron"
	case "grave1":
		return "oopen"
	case "original":
		return "origof"
	case "rightangle":
		return "orthogonal"
	case "veebar":
		return "orunderscore"
	case "xor":
		return "orunderscore"
	case "mturn":
		return "oslashacute"
	case "ostrokeacute":
		return "oslashacute"
	case "overbar":
		return "overlinecmb"
	case "nHdownarrow":
		return "pagedown"
	case "nHuparrow":
		return "pageup"
	case "afii59727":
		return "paiyannoithai"
	case "bardbl2":
		return "parallel"
	case "vextenddouble":
		return "parallel"
	case "filledparallelogram":
		return "parallelogramblack"
	case "parenleftBig":
		return "parenleft"
	case "parenleftBigg":
		return "parenleft"
	case "parenleftbig":
		return "parenleft"
	case "parenleftbigg":
		return "parenleft"
	case "ornateleftparenthesis":
		return "parenleftaltonearabic"
	case "parenleftsub":
		return "parenleftinferior"
	case "parenrightBig":
		return "parenright"
	case "parenrightBigg":
		return "parenright"
	case "parenrightbig":
		return "parenright"
	case "parenrightbigg":
		return "parenright"
	case "ornaterightparenthesis":
		return "parenrightaltonearabic"
	case "parenrightsub":
		return "parenrightinferior"
	case "help":
		return "parenrighttp"
	case "partial":
		return "partialdiff"
	case "null":
		return "pashtahebrew"
	case "patahquarterhebrew":
		return "patah11"
	case "afii57798":
		return "patah11"
	case "patah2a":
		return "patah11"
	case "patahwidehebrew":
		return "patah11"
	case "patahhebrew":
		return "patah11"
	case "patah1d":
		return "patah11"
	case "patah":
		return "patah11"
	case "recordseparator":
		return "patah11"
	case "patahnarrowhebrew":
		return "patah11"
	case "backspace":
		return "pazerhebrew"
	case "afii10081":
		return "pecyrillic"
	case "pecyril":
		return "pecyrillic"
	case "pedageshhebrew":
		return "pedagesh"
	case "pewithdagesh":
		return "pedagesh"
	case "finalpewithdagesh":
		return "pefinaldageshhebrew"
	case "afii57506":
		return "peharabic"
	case "peh":
		return "peharabic"
	case "pehfinal":
		return "pehfinalarabic"
	case "pehinitial":
		return "pehinitialarabic"
	case "pehmedial":
		return "pehmedialarabic"
	case "pewithrafe":
		return "perafehebrew"
	case "afii57381":
		return "percentarabic"
	case "arabicpercentsign":
		return "percentarabic"
	case "cdotp":
		return "periodcentered"
	case "middot":
		return "periodcentered"
	case "doublebarwedge":
		return "perpcorrespond"
	case "bot":
		return "perpendicular"
	case "Pts":
		return "peseta"
	case "pesetas":
		return "peseta"
	case "lcedilla1":
		return "peso1"
	case "pesoph":
		return "peso1"
	case "upvarphi":
		return "phi"
	case "phisymbolgreek":
		return "phi1"
	case "upphi":
		return "phi1"
	case "zecedilla":
		return "phi1"
	case "overscorenosp":
		return "philatin"
	case "afii59738":
		return "phinthuthai"
	case "Dzhacek":
		return "phook"
	case "afii59710":
		return "phophanthai"
	case "afii59708":
		return "phophungthai"
	case "afii59712":
		return "phosamphaothai"
	case "uppi":
		return "pi"
	case "arrowleftnot":
		return "piwrarmenian"
	case "planckover2pi1":
		return "planckover2pi"
	case "hslash":
		return "planckover2pi"
	case "pi1":
		return "plusbelowcmb"
	case "plussub":
		return "plusinferior"
	case "pm":
		return "plusminus"
	case "afii59707":
		return "poplathai"
	case "precnapprox":
		return "precedenotdbleqv"
	case "precneqq":
		return "precedenotslnteql"
	case "preceedsnotsimilar":
		return "precedeornoteqvlnt"
	case "precnsim":
		return "precedeornoteqvlnt"
	case "prec":
		return "precedes"
	case "notprecedesoreql":
		return "precedesequal"
	case "preceq":
		return "precedesequal"
	case "preccurlyeq":
		return "precedesorcurly"
	case "precedesequal1":
		return "precedesorcurly"
	case "precedequivlnt":
		return "precedesorequal"
	case "precsim":
		return "precedesorequal"
	case "Rx":
		return "prescription"
	case "backprime":
		return "primereversed"
	case "minuterev":
		return "primereversed"
	case "primerev":
		return "primereversed"
	case "primerev1":
		return "primereversed"
	case "primereverse":
		return "primereversed"
	case "prod":
		return "product"
	case "productdisplay":
		return "product"
	case "producttext":
		return "product"
	case "varbarwedge":
		return "projective"
	case "subset":
		return "propersubset"
	case "superset":
		return "propersuperset"
	case "supset":
		return "propersuperset"
	case "Colon":
		return "proportion"
	case "propto":
		return "proportional"
	case "lowerrank":
		return "prurel"
	case "uppsi":
		return "psi"
	case "shiftin":
		return "qadmahebrew"
	case "qaffinal":
		return "qaffinalarabic"
	case "qafinitial":
		return "qafinitialarabic"
	case "qafmedial":
		return "qafmedialarabic"
	case "acknowledge":
		return "qarneyparahebrew"
	case "circumflexsubnosp":
		return "qhook"
	case "qofdageshhebrew":
		return "qofdagesh"
	case "qofwithdagesh":
		return "qofdagesh"
	case "afii57687":
		return "qofqubutshebrew"
	case "qof":
		return "qofqubutshebrew"
	case "qofhatafpatah":
		return "qofqubutshebrew"
	case "qofhatafpatahhebrew":
		return "qofqubutshebrew"
	case "qofhatafsegol":
		return "qofqubutshebrew"
	case "qofhatafsegolhebrew":
		return "qofqubutshebrew"
	case "qofhebrew":
		return "qofqubutshebrew"
	case "qofhiriq":
		return "qofqubutshebrew"
	case "qofhiriqhebrew":
		return "qofqubutshebrew"
	case "qofholam":
		return "qofqubutshebrew"
	case "qofholamhebrew":
		return "qofqubutshebrew"
	case "qofpatah":
		return "qofqubutshebrew"
	case "qofpatahhebrew":
		return "qofqubutshebrew"
	case "qofqamats":
		return "qofqubutshebrew"
	case "qofqamatshebrew":
		return "qofqubutshebrew"
	case "qofqubuts":
		return "qofqubutshebrew"
	case "qofsegol":
		return "qofqubutshebrew"
	case "qofsegolhebrew":
		return "qofqubutshebrew"
	case "qofsheva":
		return "qofqubutshebrew"
	case "qofshevahebrew":
		return "qofqubutshebrew"
	case "qoftsere":
		return "qofqubutshebrew"
	case "qoftserehebrew":
		return "qofqubutshebrew"
	case "afii57796":
		return "qubutswidehebrew"
	case "blankb":
		return "qubutswidehebrew"
	case "qibuts":
		return "qubutswidehebrew"
	case "qubuts":
		return "qubutswidehebrew"
	case "qubuts18":
		return "qubutswidehebrew"
	case "qubuts25":
		return "qubutswidehebrew"
	case "qubuts31":
		return "qubutswidehebrew"
	case "qubutshebrew":
		return "qubutswidehebrew"
	case "qubutsnarrowhebrew":
		return "qubutswidehebrew"
	case "qubutsquarterhebrew":
		return "qubutswidehebrew"
	case "questionequal":
		return "questeq"
	case "quotesinglleft":
		return "quoteleft"
	case "quoteleftreversed":
		return "quotereversed"
	case "quotesinglrev":
		return "quotereversed"
	case "quotesinglright":
		return "quoteright"
	case "napostrophe":
		return "quoterightn"
	case "radicalBig":
		return "radical"
	case "radicalBigg":
		return "radical"
	case "radicalbig":
		return "radical"
	case "radicalbigg":
		return "radical"
	case "radicalbt":
		return "radical"
	case "radicaltp":
		return "radical"
	case "radicalvertex":
		return "radical"
	case "sqrt":
		return "radical"
	case "squareroot":
		return "radical"
	case "mathratio":
		return "ratio"
	case "rcommaaccent":
		return "rcedilla"
	case "Rsmallcapinv":
		return "rdblgrave"
	case "soundcopyright":
		return "recordright"
	case "refmark":
		return "referencemark"
	case "subseteq":
		return "reflexsubset"
	case "subsetorequal":
		return "reflexsubset"
	case "supersetorequal":
		return "reflexsuperset"
	case "supseteq":
		return "reflexsuperset"
	case "circleR":
		return "registered"
	case "afii57425":
		return "reharabic"
	case "reh":
		return "reharabic"
	case "rehyehaleflamarabic":
		return "reharabic"
	case "arrownortheast":
		return "reharmenian"
	case "rehfinal":
		return "rehfinalarabic"
	case "reshwithdagesh":
		return "reshdageshhebrew"
	case "afii57688":
		return "reshhiriq"
	case "resh":
		return "reshhiriq"
	case "reshhatafpatah":
		return "reshhiriq"
	case "reshhatafpatahhebrew":
		return "reshhiriq"
	case "reshhatafsegol":
		return "reshhiriq"
	case "reshhatafsegolhebrew":
		return "reshhiriq"
	case "reshhebrew":
		return "reshhiriq"
	case "reshhiriqhebrew":
		return "reshhiriq"
	case "reshholam":
		return "reshhiriq"
	case "reshholamhebrew":
		return "reshhiriq"
	case "reshpatah":
		return "reshhiriq"
	case "reshpatahhebrew":
		return "reshhiriq"
	case "reshqamats":
		return "reshhiriq"
	case "reshqamatshebrew":
		return "reshhiriq"
	case "reshqubuts":
		return "reshhiriq"
	case "reshqubutshebrew":
		return "reshhiriq"
	case "reshsegol":
		return "reshhiriq"
	case "reshsegolhebrew":
		return "reshhiriq"
	case "reshsheva":
		return "reshhiriq"
	case "reshshevahebrew":
		return "reshhiriq"
	case "reshtsere":
		return "reshhiriq"
	case "reshtserehebrew":
		return "reshhiriq"
	case "backsimeq":
		return "revasymptequal"
	case "backsim":
		return "reversedtilde"
	case "revsimilar":
		return "reversedtilde"
	case "tildereversed":
		return "reversedtilde"
	case "arrowlongbothnot":
		return "reviamugrashhebrew"
	case "reviahebrew":
		return "reviamugrashhebrew"
	case "invnot":
		return "revlogicalnot"
	case "logicalnotreversed":
		return "revlogicalnot"
	case "acutedblnosp":
		return "rfishhook"
	case "haceknosp":
		return "rfishhookreversed"
	case "uprho":
		return "rho"
	case "ringnosp":
		return "rhook"
	case "dieresisnosp":
		return "rhookturned"
	case "tetse":
		return "rhosymbolgreek"
	case "upvarrho":
		return "rhosymbolgreek"
	case "urcorner":
		return "rightanglene"
	case "ulcorner":
		return "rightanglenw"
	case "lrcorner":
		return "rightanglese"
	case "llcorner":
		return "rightanglesw"
	case "beta1":
		return "righttackbelowcmb"
	case "varlrtriangle":
		return "righttriangle"
	case "ocirc":
		return "ringcmb"
	case "Upsilon1tonos":
		return "ringhalfleftbelowcmb"
	case "numeralgreeksub":
		return "ringhalfright"
	case "kappa1":
		return "ringhalfrightbelowcmb"
	case "eqcirc":
		return "ringinequal"
	case "hooksupnosp":
		return "rlongleg"
	case "dotnosp":
		return "rlonglegturned"
	case "afii59715":
		return "roruathai"
	case "afii57513":
		return "rreharabic"
	case "blockrighthalf":
		return "rtblock"
	case "brevenosp":
		return "rturned"
	case "acuterightnosp":
		return "rturnedsuperior"
	case "rturnhooksuper":
		return "rturnrthooksuper"
	case "rupees":
		return "rupee"
	case "afii59716":
		return "ruthai"
	case "sadfinal":
		return "sadfinalarabic"
	case "sadinitial":
		return "sadinitialarabic"
	case "sadmedial":
		return "sadmedialarabic"
	case "afii57681":
		return "samekh"
	case "samekhhebrew":
		return "samekh"
	case "samekhdagesh":
		return "samekhdageshhebrew"
	case "samekhwithdagesh":
		return "samekhdageshhebrew"
	case "afii59730":
		return "saraaathai"
	case "afii59745":
		return "saraaethai"
	case "afii59748":
		return "saraaimaimalaithai"
	case "afii59747":
		return "saraaimaimuanthai"
	case "afii59731":
		return "saraamthai"
	case "afii59729":
		return "saraathai"
	case "afii59744":
		return "saraethai"
	case "afii59733":
		return "saraiithai"
	case "afii59732":
		return "saraithai"
	case "afii59746":
		return "saraothai"
	case "afii59735":
		return "saraueethai"
	case "afii59734":
		return "sarauethai"
	case "afii59736":
		return "sarauthai"
	case "afii59737":
		return "sarauuthai"
	case "satisfy":
		return "satisfies"
	case "vDash":
		return "satisfies"
	case "length":
		return "schwa"
	case "afii10846":
		return "schwacyrillic"
	case "halflength":
		return "schwahook"
	case "higherrank":
		return "scurel"
	case "dprime":
		return "second"
	case "primedbl":
		return "second"
	case "primedbl1":
		return "second"
	case "seenfinal":
		return "seenfinalarabic"
	case "seeninitial":
		return "seeninitialarabic"
	case "seenmedial":
		return "seenmedialarabic"
	case "afii57795":
		return "segolhebrew"
	case "groupseparator":
		return "segolhebrew"
	case "segol":
		return "segolhebrew"
	case "segol1f":
		return "segolhebrew"
	case "segol2c":
		return "segolhebrew"
	case "segol13":
		return "segolhebrew"
	case "segolnarrowhebrew":
		return "segolhebrew"
	case "segolquarterhebrew":
		return "segolhebrew"
	case "segolwidehebrew":
		return "segolhebrew"
	case "arrowlongboth":
		return "seharmenian"
	case "sevensub":
		return "seveninferior"
	case "extendedarabicindicdigitseven":
		return "sevenpersian"
	case "afii59767":
		return "seventhai"
	case "afii57457":
		return "shaddaarabic"
	case "shadda":
		return "shaddaarabic"
	case "shaddafathatanarabic":
		return "shaddaarabic"
	case "shaddawithdammaisolated":
		return "shaddadammaarabic"
	case "shaddawithdammatanisolated":
		return "shaddadammatanarabic"
	case "shaddawithfathaisolated":
		return "shaddafathaarabic"
	case "shaddamedial":
		return "shaddahontatweel"
	case "shaddawithkasraisolated":
		return "shaddakasraarabic"
	case "shaddawithkasratanisolated":
		return "shaddakasratanarabic"
	case "shaddawithdammalow":
		return "shaddawithdammaisolatedlow"
	case "shaddawithdammatanlow":
		return "shaddawithdammatanisolatedlow"
	case "shaddawithfathaisolatedlow":
		return "shaddawithfathalow"
	case "shaddawithfathatanisolatedlow":
		return "shaddawithfathatanlow"
	case "shaddawithkasralow":
		return "shaddawithkasraisolatedlow"
	case "shaddawithkasratanlow":
		return "shaddawithkasratanisolatedlow"
	case "blockhalfshaded":
		return "shade"
	case "shademedium":
		return "shade"
	case "blockqtrshaded":
		return "shadelight"
	case "ltshade":
		return "shadelight"
	case "sheenfinal":
		return "sheenfinalarabic"
	case "sheeninitial":
		return "sheeninitialarabic"
	case "sheenmedial":
		return "sheenmedialarabic"
	case "pehook":
		return "sheicoptic"
	case "Lsh":
		return "shiftleft"
	case "Rsh":
		return "shiftright"
	case "ustrtbar":
		return "shimacoptic"
	case "afii57689":
		return "shin"
	case "shinhebrew":
		return "shin"
	case "shindageshhebrew":
		return "shindagesh"
	case "shinwithdagesh":
		return "shindagesh"
	case "shindageshshindothebrew":
		return "shindageshshindot"
	case "shinwithdageshandshindot":
		return "shindageshshindot"
	case "shindageshsindot":
		return "shindageshsindothebrew"
	case "shinwithdageshandsindot":
		return "shindageshsindothebrew"
	case "afii57804":
		return "shindothebrew"
	case "shindot":
		return "shindothebrew"
	case "afii57694":
		return "shinshindot"
	case "shinshindothebrew":
		return "shinshindot"
	case "shinwithshindot":
		return "shinshindot"
	case "gravedblnosp":
		return "shook"
	case "upsigma":
		return "sigma"
	case "upvarsigma":
		return "sigma1"
	case "sigmafinal":
		return "sigma1"
	case "Chertdes":
		return "sigmalunatesymbolgreek"
	case "afii57839":
		return "siluqlefthebrew"
	case "meteg":
		return "siluqlefthebrew"
	case "newline":
		return "siluqlefthebrew"
	case "siluqhebrew":
		return "siluqlefthebrew"
	case "sim":
		return "similar"
	case "tildemath":
		return "similar"
	case "tildeoperator":
		return "similar"
	case "approxnotequal":
		return "simneqq"
	case "sine":
		return "sinewave"
	case "sixsub":
		return "sixinferior"
	case "extendedarabicindicdigitsix":
		return "sixpersian"
	case "afii59766":
		return "sixthai"
	case "mathslash":
		return "slash"
	case "slashBig":
		return "slash"
	case "slashBigg":
		return "slash"
	case "slashbig":
		return "slash"
	case "slashbigg":
		return "slash"
	case "frown":
		return "slurabove"
	case "smalltriangleleftsld":
		return "smallblacktriangleleft"
	case "smalltrianglerightsld":
		return "smallblacktriangleright"
	case "elementsmall":
		return "smallin"
	case "smallelement":
		return "smallin"
	case "ownersmall":
		return "smallni"
	case "smallcontains":
		return "smallni"
	case "slurbelow":
		return "smile"
	case "whitesmilingface":
		return "smileface"
	case "afii57658":
		return "sofpasuqhebrew"
	case "sofpasuq":
		return "sofpasuqhebrew"
	case "sfthyphen":
		return "softhyphen"
	case "afii10094":
		return "softsigncyrillic"
	case "soft":
		return "softsigncyrillic"
	case "dei":
		return "soliduslongoverlaycmb"
	case "negationslash":
		return "soliduslongoverlaycmb"
	case "not":
		return "soliduslongoverlaycmb"
	case "Dei":
		return "solidusshortoverlaycmb"
	case "afii59721":
		return "sorusithai"
	case "afii59720":
		return "sosalathai"
	case "afii59691":
		return "sosothai"
	case "afii59722":
		return "sosuathai"
	case "spacehackarabic":
		return "space"
	case "a109":
		return "spade"
	case "spadesuit":
		return "spade"
	case "spadesuitblack":
		return "spade"
	case "varspadesuit":
		return "spadesuitwhite"
	case "sqimageornotequal":
		return "sqsubsetneq"
	case "sqoriginornotequal":
		return "sqsupsetneq"
	case "sigmalunate":
		return "squarebelowcmb"
	case "boxcrossdiaghatch":
		return "squarediagonalcrosshatchfill"
	case "squarecrossfill":
		return "squarediagonalcrosshatchfill"
	case "boxdot":
		return "squaredot"
	case "boxhorizhatch":
		return "squarehorizontalfill"
	case "squarehfill":
		return "squarehorizontalfill"
	case "sqsubset":
		return "squareimage"
	case "squareleftsld":
		return "squareleftblack"
	case "squaresesld":
		return "squarelrblack"
	case "boxminus":
		return "squareminus"
	case "boxtimes":
		return "squaremultiply"
	case "sqsupset":
		return "squareoriginal"
	case "boxcrosshatch":
		return "squareorthogonalcrosshatchfill"
	case "squarehvfill":
		return "squareorthogonalcrosshatchfill"
	case "boxplus":
		return "squareplus"
	case "squarerightsld":
		return "squarerightblack"
	case "squarenwsld":
		return "squareulblack"
	case "boxleftdiaghatch":
		return "squareupperlefttolowerrightfill"
	case "squarenwsefill":
		return "squareupperlefttolowerrightfill"
	case "boxrtdiaghatch":
		return "squareupperrighttolowerleftfill"
	case "squareneswfill":
		return "squareupperrighttolowerleftfill"
	case "boxverthatch":
		return "squareverticalfill"
	case "squarevfill":
		return "squareverticalfill"
	case "blackinwhitesquare":
		return "squarewhitewithsmallblack"
	case "boxnested":
		return "squarewhitewithsmallblack"
	case "leftrightsquigarrow":
		return "squiggleleftright"
	case "arrowsquiggleright":
		return "squiggleright"
	case "rightsquigarrow":
		return "squiggleright"
	case "boxrounded":
		return "squoval"
	case "starequal":
		return "stareq"
	case "Subset":
		return "subsetdbl"
	case "notsubsetordbleql":
		return "subsetdblequal"
	case "subseteqq":
		return "subsetdblequal"
	case "notsubsetoreql":
		return "subsetnotequal"
	case "subsetneq":
		return "subsetnotequal"
	case "subsetnoteql":
		return "subsetnotequal"
	case "subsetneqq":
		return "subsetornotdbleql"
	case "sqsubseteq":
		return "subsetsqequal"
	case "follows":
		return "succeeds"
	case "succ":
		return "succeeds"
	case "contains":
		return "suchthat"
	case "ni":
		return "suchthat"
	case "owner":
		return "suchthat"
	case "afii57458":
		return "sukunarabic"
	case "sukun":
		return "sukunarabic"
	case "sukunontatweel":
		return "sukunmedial"
	case "sum":
		return "summation"
	case "summationdisplay":
		return "summation"
	case "summationtext":
		return "summation"
	case "compass":
		return "sun"
	case "Supset":
		return "supersetdbl"
	case "notsupersetordbleql":
		return "supersetdblequal"
	case "supseteqq":
		return "supersetdblequal"
	case "notsupersetoreql":
		return "supersetnotequal"
	case "supersetnoteql":
		return "supersetnotequal"
	case "supsetneq":
		return "supersetnotequal"
	case "supsetneqq":
		return "supersetornotdbleql"
	case "sqsupseteq":
		return "supersetsqequal"
	case "latticetop":
		return "tackdown"
	case "top":
		return "tackdown"
	case "dashv":
		return "tackleft"
	case "turnstileright":
		return "tackleft"
	case "afii57431":
		return "taharabic"
	case "tah":
		return "taharabic"
	case "tahfinal":
		return "tahfinalarabic"
	case "tahinitial":
		return "tahinitialarabic"
	case "tahmedial":
		return "tahmedialarabic"
	case "fathatanontatweel":
		return "tatweelwithfathatanabove"
	case "uptau":
		return "tau"
	case "tavdages":
		return "tavdagesh"
	case "tavdageshhebrew":
		return "tavdagesh"
	case "tavwithdagesh":
		return "tavdagesh"
	case "afii57690":
		return "tavhebrew"
	case "tav":
		return "tavhebrew"
	case "tcaronaltone":
		return "tcaron1"
	case "barmidshortnosp":
		return "tccurl"
	case "tcommaaccent":
		return "tcedilla"
	case "kcedilla1":
		return "tcedilla1"
	case "afii57507":
		return "tcheharabic"
	case "tcheh":
		return "tcheharabic"
	case "tchehfinal":
		return "tchehfinalarabic"
	case "tchehinitial":
		return "tchehinitialarabic"
	case "tchehmeeminitialarabic":
		return "tchehinitialarabic"
	case "tchehmedial":
		return "tchehmedialarabic"
	case "tehfinal":
		return "tehfinalarabic"
	case "tehwithhahinitial":
		return "tehhahinitialarabic"
	case "tehinitial":
		return "tehinitialarabic"
	case "tehwithjeeminitial":
		return "tehjeeminitialarabic"
	case "afii57417":
		return "tehmarbutaarabic"
	case "tehmarbuta":
		return "tehmarbutaarabic"
	case "tehmarbutafinal":
		return "tehmarbutafinalarabic"
	case "tehmedial":
		return "tehmedialarabic"
	case "tehwithmeeminitial":
		return "tehmeeminitialarabic"
	case "tehwithmeemisolated":
		return "tehmeemisolatedarabic"
	case "tehwithnoonfinal":
		return "tehnoonfinalarabic"
	case "tel":
		return "telephone"
	case "bell":
		return "telishagedolahebrew"
	case "datalinkescape":
		return "telishaqetanahebrew"
	case "devcon0":
		return "telishaqetanahebrew"
	case "tildemidnosp":
		return "tesh"
	case "tetdageshhebrew":
		return "tetdagesh"
	case "tetwithdagesh":
		return "tetdagesh"
	case "afii57672":
		return "tethebrew"
	case "tet":
		return "tethebrew"
	case "Lcircumflex":
		return "tetsecyrillic"
	case "starttext":
		return "tevirhebrew"
	case "tevirlefthebrew":
		return "tevirhebrew"
	case "afii57424":
		return "thalarabic"
	case "thal":
		return "thalarabic"
	case "thalfinal":
		return "thalfinalarabic"
	case "afii59756":
		return "thanthakhatthai"
	case "afii57419":
		return "theharabic"
	case "theh":
		return "theharabic"
	case "thehfinal":
		return "thehfinalarabic"
	case "thehinitial":
		return "thehinitialarabic"
	case "thehmedial":
		return "thehmedialarabic"
	case "uptheta":
		return "theta"
	case "gehook":
		return "theta1"
	case "upvartheta":
		return "theta1"
	case "thetasymbolgreek":
		return "theta1"
	case "afii59697":
		return "thonangmonthothai"
	case "Ahacek":
		return "thook"
	case "afii59698":
		return "thophuthaothai"
	case "afii59703":
		return "thothahanthai"
	case "afii59696":
		return "thothanthai"
	case "afii59704":
		return "thothongthai"
	case "afii59702":
		return "thothungthai"
	case "thousandsseparatorpersian":
		return "thousandsseparatorarabic"
	case "threesub":
		return "threeinferior"
	case "extendedarabicindicdigitthree":
		return "threepersian"
	case "afii59763":
		return "threethai"
	case "tie":
		return "tieconcat"
	case "tie1":
		return "tieconcat"
	case "ilde":
		return "tilde"
	case "tildewide":
		return "tilde"
	case "tildewider":
		return "tilde"
	case "tildewidest":
		return "tilde"
	case "wideutilde":
		return "tildebelowcmb"
	case "tildecomb":
		return "tildecmb"
	case "arrowwaveboth":
		return "tipehahebrew"
	case "tipehalefthebrew":
		return "tipehahebrew"
	case "arrownorthwest":
		return "tiwnarmenian"
	case "eturn":
		return "tonefive"
	case "afii59695":
		return "topatakthai"
	case "toparc":
		return "topsemicircle"
	case "afii59701":
		return "totaothai"
	case "commasuprightnosp":
		return "tretroflexhook"
	case "triangledot":
		return "trianglecdot"
	case "triangleleftsld":
		return "triangleleftblack"
	case "triangleftequal":
		return "triangleleftequal"
	case "trianglelefteq":
		return "triangleleftequal"
	case "trianglerightsld":
		return "trianglerightblack"
	case "trianglerighteq":
		return "trianglerightequal"
	case "triangrightequal":
		return "trianglerightequal"
	case "primetripl":
		return "trprime"
	case "primetripl1":
		return "trprime"
	case "underscoredblnosp":
		return "ts"
	case "tsadidageshhebrew":
		return "tsadidagesh"
	case "tsadiwithdagesh":
		return "tsadidagesh"
	case "afii10088":
		return "tsecyrillic"
	case "tse":
		return "tsecyrillic"
	case "afii57794":
		return "tsere12"
	case "tserenarrowhebrew":
		return "tsere12"
	case "tserehebrew":
		return "tsere12"
	case "tsere1e":
		return "tsere12"
	case "tsere":
		return "tsere12"
	case "tserewidehebrew":
		return "tsere12"
	case "fileseparator":
		return "tsere12"
	case "tsere2b":
		return "tsere12"
	case "tserequarterhebrew":
		return "tsere12"
	case "afii10108":
		return "tshecyrillic"
	case "tshe":
		return "tshecyrillic"
	case "commasuprevnosp":
		return "tturned"
	case "iotaturn":
		return "turnediota"
	case "vdash":
		return "turnstileleft"
	case "afii57394":
		return "twoarabic"
	case "arabicindicdigittwo":
		return "twoarabic"
	case "twohackarabic":
		return "twoarabic"
	case "enleadertwodots":
		return "twodotleader"
	case "twodotenleader":
		return "twodotleader"
	case "twodotlead":
		return "twodotleader"
	case "twosub":
		return "twoinferior"
	case "extendedarabicindicdigittwo":
		return "twopersian"
	case "afii59762":
		return "twothai"
	case "gravesubnosp":
		return "ubar"
	case "deltaturn":
		return "ubreve"
	case "uhungarumlaut":
		return "udblacute"
	case "eshshortrev":
		return "udblgrave"
	case "Aacutering":
		return "udieresiscaron"
	case "ihacek":
		return "uhorn"
	case "tturn":
		return "uinvertedbreve"
	case "nwquadarc":
		return "ularc"
	case "dbllowline":
		return "underscoredbl"
	case "twolowline":
		return "underscoredbl"
	case "midhorizellipsis":
		return "unicodecdots"
	case "cup":
		return "union"
	case "Cup":
		return "uniondbl"
	case "unionmultidisplay":
		return "unionmulti"
	case "unionmultitext":
		return "unionmulti"
	case "uplus":
		return "unionmulti"
	case "sqcup":
		return "unionsq"
	case "unionsqdisplay":
		return "unionsq"
	case "unionsqtext":
		return "unionsq"
	case "bigcup":
		return "uniontext"
	case "naryunion":
		return "uniontext"
	case "uniondisplay":
		return "uniontext"
	case "forall":
		return "universal"
	case "blockuphalf":
		return "upblock"
	case "gekarev":
		return "updigamma"
	case "enrtdes":
		return "upkoppa"
	case "Kavertbar":
		return "upoldKoppa"
	case "kavertbar":
		return "upoldkoppa"
	case "enge":
		return "upsampi"
	case "upupsilon":
		return "upsilon"
	case "acutesubnosp":
		return "upsilonlatin"
	case "xsol":
		return "upslope"
	case "kabar":
		return "upstigma"
	case "Upsilon1dieresis":
		return "uptackbelowcmb"
	case "Upsilon1diaeresis":
		return "uptackbelowcmb"
	case "Chevertbar":
		return "upvarTheta"
	case "nequadarc":
		return "urarc"
	case "Dbar1":
		return "utilde"
	case "perspcorrespond":
		return "vardoublebarwedge"
	case "clwcontintegral":
		return "varointclockwise"
	case "triangleright":
		return "vartriangleleft"
	case "triangleleft":
		return "vartriangleright"
	case "afii57669":
		return "vav"
	case "vavhebrew":
		return "vav"
	case "afii57723":
		return "vavdageshhebrew"
	case "vavdagesh":
		return "vavdageshhebrew"
	case "vavdagesh65":
		return "vavdageshhebrew"
	case "vavwithdagesh":
		return "vavdageshhebrew"
	case "afii57700":
		return "vavholam"
	case "vavholamhebrew":
		return "vavholam"
	case "vavwithholam":
		return "vavholam"
	case "vec":
		return "vector"
	case "equiangular":
		return "veeeq"
	case "afii57505":
		return "veharabic"
	case "veh":
		return "veharabic"
	case "vehfinal":
		return "vehfinalarabic"
	case "vehinitial":
		return "vehinitialarabic"
	case "vehmedial":
		return "vehmedialarabic"
	case "Sampi":
		return "verticallinebelowcmb"
	case "arrowlongbothv":
		return "vewarmenian"
	case "tackleftsubnosp":
		return "vhook"
	case "vertrectangle":
		return "vrectangle"
	case "filledvertrect":
		return "vrectangleblack"
	case "tackrightsubnosp":
		return "vturned"
	case "openbullet1":
		return "vysmwhtcircle"
	case "ringmath":
		return "vysmwhtcircle"
	case "afii57448":
		return "wawarabic"
	case "waw":
		return "wawarabic"
	case "wawfinal":
		return "wawfinalarabic"
	case "wawwithhamzaabovefinal":
		return "wawhamzaabovefinalarabic"
	case "estimates":
		return "wedgeq"
	case "Pscript":
		return "weierstrass"
	case "wp":
		return "weierstrass"
	case "openbullet":
		return "whitebullet"
	case "smwhtcircle":
		return "whitebullet"
	case "circle":
		return "whitecircle"
	case "mdlgwhtcircle":
		return "whitecircle"
	case "diamondrhomb":
		return "whitediamond"
	case "mdlgwhtdiamond":
		return "whitediamond"
	case "blackinwhitediamond":
		return "whitediamondcontainingblacksmalldiamond"
	case "diamondrhombnested":
		return "whitediamondcontainingblacksmalldiamond"
	case "smalltriangleinv":
		return "whitedownpointingsmalltriangle"
	case "triangledown":
		return "whitedownpointingsmalltriangle"
	case "bigtriangledown":
		return "whitedownpointingtriangle"
	case "triangleinv":
		return "whitedownpointingtriangle"
	case "smalltriangleleft":
		return "whiteleftpointingsmalltriangle"
	case "triangleleft1":
		return "whiteleftpointingtriangle"
	case "triaglfopen":
		return "whitepointerleft"
	case "triagrtopen":
		return "whitepointerright"
	case "smalltriangleright":
		return "whiterightpointingsmalltriangle"
	case "triangleright1":
		return "whiterightpointingtriangle"
	case "H18551":
		return "whitesmallsquare"
	case "smallbox":
		return "whitesmallsquare"
	case "smwhtsquare":
		return "whitesmallsquare"
	case "bigwhitestar":
		return "whitestar"
	case "smalltriangle":
		return "whiteuppointingsmalltriangle"
	case "vartriangle":
		return "whiteuppointingsmalltriangle"
	case "bigtriangleup":
		return "whiteuppointingtriangle"
	case "triangle":
		return "whiteuppointingtriangle"
	case "afii59719":
		return "wowaenthai"
	case "wr":
		return "wreathproduct"
	case "diaeresistonosnosp":
		return "wsuperior"
	case "anglesupnosp":
		return "wturned"
	case "upxi":
		return "xi"
	case "afii59758":
		return "yamakkanthai"
	case "afii10194":
		return "yatcyrillic"
	case "Ibar":
		return "ycircumflex"
	case "afii57450":
		return "yeharabic"
	case "yeh":
		return "yeharabic"
	case "afii57519":
		return "yehbarreearabic"
	case "yehfinal":
		return "yehfinalarabic"
	case "afii57414":
		return "yehhamzaabovearabic"
	case "yehwithhamzaabove":
		return "yehhamzaabovearabic"
	case "yehwithhamzaabovefinal":
		return "yehhamzaabovefinalarabic"
	case "yehwithhamzaaboveinitial":
		return "yehhamzaaboveinitialarabic"
	case "yehwithhamzaabovemedial":
		return "yehhamzaabovemedialarabic"
	case "alefmaksurainitialarabic":
		return "yehinitialarabic"
	case "yehinitial":
		return "yehinitialarabic"
	case "yehwithmeeminitial":
		return "yehmeeminitialarabic"
	case "yehwithmeemisolated":
		return "yehmeemisolatedarabic"
	case "yehwithnoonfinal":
		return "yehnoonfinalarabic"
	case "Yen":
		return "yen"
	case "auxiliaryon":
		return "yerahbenyomohebrew"
	case "devcon1":
		return "yerahbenyomohebrew"
	case "yerahbenyomolefthebrew":
		return "yerahbenyomohebrew"
	case "afii10093":
		return "yericyrillic"
	case "yeri":
		return "yericyrillic"
	case "startofhead":
		return "yetivhebrew"
	case "uhacek":
		return "yhook"
	case "afii10104":
		return "yicyrillic"
	case "yi":
		return "yicyrillic"
	case "arrowsouthwest":
		return "yiwnarmenian"
	case "yoddagesh":
		return "yoddageshhebrew"
	case "yodwithdagesh":
		return "yoddageshhebrew"
	case "afii57718":
		return "yodyodhebrew"
	case "yoddbl":
		return "yodyodhebrew"
	case "afii57705":
		return "yodyodpatahhebrew"
	case "doubleyodpatah":
		return "yodyodpatahhebrew"
	case "doubleyodpatahhebrew":
		return "yodyodpatahhebrew"
	case "chertdes":
		return "yotgreek"
	case "afii59714":
		return "yoyakthai"
	case "afii59693":
		return "yoyingthai"
	case "dzhacek":
		return "yr"
	case "iotasubnosp":
		return "ysuperior"
	case "hornnosp":
		return "yturned"
	case "afii57432":
		return "zaharabic"
	case "zah":
		return "zaharabic"
	case "zahfinal":
		return "zahfinalarabic"
	case "zahinitial":
		return "zahinitialarabic"
	case "zahmedial":
		return "zahmedialarabic"
	case "afii57426":
		return "zainarabic"
	case "zain":
		return "zainarabic"
	case "zainfinal":
		return "zainfinalarabic"
	case "arrowloopright":
		return "zaqefgadolhebrew"
	case "arrowloopleft":
		return "zaqefqatanhebrew"
	case "arrowzigzag":
		return "zarqahebrew"
	case "zayindagesh":
		return "zayindageshhebrew"
	case "zayinwithdagesh":
		return "zayindageshhebrew"
	case "nleg":
		return "zcaron"
	case "tackdownsubnosp":
		return "zcurl"
	case "mcapturn":
		return "zdotaccent"
	case "zdot":
		return "zdotaccent"
	case "zerodot":
		return "zero"
	case "zeroslash":
		return "zero"
	case "afii57392":
		return "zerohackarabic"
	case "arabicindicdigitzero":
		return "zerohackarabic"
	case "zeroarabic":
		return "zerohackarabic"
	case "zerosub":
		return "zeroinferior"
	case "extendedarabicindicdigitzero":
		return "zeropersian"
	case "afii59760":
		return "zerothai"
	case "bom":
		return "zerowidthjoiner"
	case "zerowidthnobreakspace":
		return "zerowidthjoiner"
	case "zerospace":
		return "zerowidthspace"
	case "upzeta":
		return "zeta"
	case "afii10072":
		return "zhecyrillic"
	case "zhe":
		return "zhecyrillic"
	case "negacknowledge":
		return "zinorhebrew"
	case "tackupsubnosp":
		return "zretroflexhook"
	}

	return ""
}

func glyphToRune(glyph string) rune {
	// 6339 entries
	switch glyph {
	case ".notdef":
		return 0xfffd // � '\ufffd'
	case "250a":
		return 0x250a // ┊ '\u250a'
	case "250b":
		return 0x250b // ┋ '\u250b'
	case "250d":
		return 0x250d // ┍ '\u250d'
	case "250e":
		return 0x250e // ┎ '\u250e'
	case "250f":
		return 0x250f // ┏ '\u250f'
	case "251a":
		return 0x251a // ┚ '\u251a'
	case "251b":
		return 0x251b // ┛ '\u251b'
	case "251d":
		return 0x251d // ┝ '\u251d'
	case "251e":
		return 0x251e // ┞ '\u251e'
	case "251f":
		return 0x251f // ┟ '\u251f'
	case "252a":
		return 0x252a // ┪ '\u252a'
	case "252b":
		return 0x252b // ┫ '\u252b'
	case "252d":
		return 0x252d // ┭ '\u252d'
	case "252e":
		return 0x252e // ┮ '\u252e'
	case "252f":
		return 0x252f // ┯ '\u252f'
	case "253a":
		return 0x253a // ┺ '\u253a'
	case "253b":
		return 0x253b // ┻ '\u253b'
	case "253d":
		return 0x253d // ┽ '\u253d'
	case "253e":
		return 0x253e // ┾ '\u253e'
	case "253f":
		return 0x253f // ┿ '\u253f'
	case "254a":
		return 0x254a // ╊ '\u254a'
	case "254b":
		return 0x254b // ╋ '\u254b'
	case "254c":
		return 0x254c // ╌ '\u254c'
	case "254d":
		return 0x254d // ╍ '\u254d'
	case "254e":
		return 0x254e // ╎ '\u254e'
	case "254f":
		return 0x254f // ╏ '\u254f'
	case "256d":
		return 0x256d // ╭ '\u256d'
	case "256e":
		return 0x256e // ╮ '\u256e'
	case "256f":
		return 0x256f // ╯ '\u256f'
	case "257a":
		return 0x257a // ╺ '\u257a'
	case "257b":
		return 0x257b // ╻ '\u257b'
	case "257c":
		return 0x257c // ╼ '\u257c'
	case "257d":
		return 0x257d // ╽ '\u257d'
	case "257e":
		return 0x257e // ╾ '\u257e'
	case "257f":
		return 0x257f // ╿ '\u257f'
	case "A":
		return 0x0041 // A 'A'
	case "AE":
		return 0x00c6 // Æ '\u00c6'
	case "AEacute":
		return 0x01fc // Ǽ '\u01fc'
	case "AEmacron":
		return 0x01e2 // Ǣ '\u01e2'
	case "AEsmall":
		return 0xf7e6 //  '\uf7e6'
	case "APLboxquestion":
		return 0x2370 // ⍰ '\u2370'
	case "APLboxupcaret":
		return 0x2353 // ⍓ '\u2353'
	case "APLnotbackslash":
		return 0x2340 // ⍀ '\u2340'
	case "APLnotslash":
		return 0x233f // ⌿ '\u233f'
	case "Aacute":
		return 0x00c1 // Á '\u00c1'
	case "Aacutesmall":
		return 0xf7e1 //  '\uf7e1'
	case "Abreve":
		return 0x0102 // Ă '\u0102'
	case "Abreveacute":
		return 0x1eae // Ắ '\u1eae'
	case "Abrevecyrillic":
		return 0x04d0 // Ӑ '\u04d0'
	case "Abrevedotbelow":
		return 0x1eb6 // Ặ '\u1eb6'
	case "Abrevegrave":
		return 0x1eb0 // Ằ '\u1eb0'
	case "Abrevehookabove":
		return 0x1eb2 // Ẳ '\u1eb2'
	case "Abrevetilde":
		return 0x1eb4 // Ẵ '\u1eb4'
	case "Acaron":
		return 0x01cd // Ǎ '\u01cd'
	case "Acircle":
		return 0x24b6 // Ⓐ '\u24b6'
	case "Acircumflex":
		return 0x00c2 // Â '\u00c2'
	case "Acircumflexacute":
		return 0x1ea4 // Ấ '\u1ea4'
	case "Acircumflexdotbelow":
		return 0x1eac // Ậ '\u1eac'
	case "Acircumflexgrave":
		return 0x1ea6 // Ầ '\u1ea6'
	case "Acircumflexhookabove":
		return 0x1ea8 // Ẩ '\u1ea8'
	case "Acircumflexsmall":
		return 0xf7e2 //  '\uf7e2'
	case "Acircumflextilde":
		return 0x1eaa // Ẫ '\u1eaa'
	case "Acute":
		return 0xf6c9 //  '\uf6c9'
	case "Acutesmall":
		return 0xf7b4 //  '\uf7b4'
	case "Adblgrave":
		return 0x0200 // Ȁ '\u0200'
	case "Adieresis":
		return 0x00c4 // Ä '\u00c4'
	case "Adieresiscyrillic":
		return 0x04d2 // Ӓ '\u04d2'
	case "Adieresismacron":
		return 0x01de // Ǟ '\u01de'
	case "Adieresissmall":
		return 0xf7e4 //  '\uf7e4'
	case "Adotbelow":
		return 0x1ea0 // Ạ '\u1ea0'
	case "Adotmacron":
		return 0x01e0 // Ǡ '\u01e0'
	case "Agrave":
		return 0x00c0 // À '\u00c0'
	case "Agravesmall":
		return 0xf7e0 //  '\uf7e0'
	case "Ahookabove":
		return 0x1ea2 // Ả '\u1ea2'
	case "Aiecyrillic":
		return 0x04d4 // Ӕ '\u04d4'
	case "Ainvertedbreve":
		return 0x0202 // Ȃ '\u0202'
	case "Alpha":
		return 0x0391 // Α '\u0391'
	case "Alphatonos":
		return 0x0386 // Ά '\u0386'
	case "Amacron":
		return 0x0100 // Ā '\u0100'
	case "Amonospace":
		return 0xff21 // Ａ '\uff21'
	case "Aogonek":
		return 0x0104 // Ą '\u0104'
	case "Aring":
		return 0x00c5 // Å '\u00c5'
	case "Aringacute":
		return 0x01fa // Ǻ '\u01fa'
	case "Aringbelow":
		return 0x1e00 // Ḁ '\u1e00'
	case "Aringsmall":
		return 0xf7e5 //  '\uf7e5'
	case "Asmall":
		return 0xf761 //  '\uf761'
	case "Atilde":
		return 0x00c3 // Ã '\u00c3'
	case "Atildesmall":
		return 0xf7e3 //  '\uf7e3'
	case "Aybarmenian":
		return 0x0531 // Ա '\u0531'
	case "B":
		return 0x0042 // B 'B'
	case "Barv":
		return 0x2ae7 // ⫧ '\u2ae7'
	case "BbbA":
		return 0x1d538 // 𝔸 '\U0001d538'
	case "BbbB":
		return 0x1d539 // 𝔹 '\U0001d539'
	case "BbbC":
		return 0x2102 // ℂ '\u2102'
	case "BbbD":
		return 0x1d53b // 𝔻 '\U0001d53b'
	case "BbbE":
		return 0x1d53c // 𝔼 '\U0001d53c'
	case "BbbF":
		return 0x1d53d // 𝔽 '\U0001d53d'
	case "BbbG":
		return 0x1d53e // 𝔾 '\U0001d53e'
	case "BbbGamma":
		return 0x213e // ℾ '\u213e'
	case "BbbH":
		return 0x210d // ℍ '\u210d'
	case "BbbI":
		return 0x1d540 // 𝕀 '\U0001d540'
	case "BbbJ":
		return 0x1d541 // 𝕁 '\U0001d541'
	case "BbbK":
		return 0x1d542 // 𝕂 '\U0001d542'
	case "BbbL":
		return 0x1d543 // 𝕃 '\U0001d543'
	case "BbbM":
		return 0x1d544 // 𝕄 '\U0001d544'
	case "BbbN":
		return 0x2115 // ℕ '\u2115'
	case "BbbO":
		return 0x1d546 // 𝕆 '\U0001d546'
	case "BbbP":
		return 0x2119 // ℙ '\u2119'
	case "BbbPi":
		return 0x213f // ℿ '\u213f'
	case "BbbQ":
		return 0x211a // ℚ '\u211a'
	case "BbbR":
		return 0x211d // ℝ '\u211d'
	case "BbbS":
		return 0x1d54a // 𝕊 '\U0001d54a'
	case "BbbT":
		return 0x1d54b // 𝕋 '\U0001d54b'
	case "BbbU":
		return 0x1d54c // 𝕌 '\U0001d54c'
	case "BbbV":
		return 0x1d54d // 𝕍 '\U0001d54d'
	case "BbbW":
		return 0x1d54e // 𝕎 '\U0001d54e'
	case "BbbX":
		return 0x1d54f // 𝕏 '\U0001d54f'
	case "BbbY":
		return 0x1d550 // 𝕐 '\U0001d550'
	case "BbbZ":
		return 0x2124 // ℤ '\u2124'
	case "Bbba":
		return 0x1d552 // 𝕒 '\U0001d552'
	case "Bbbb":
		return 0x1d553 // 𝕓 '\U0001d553'
	case "Bbbc":
		return 0x1d554 // 𝕔 '\U0001d554'
	case "Bbbd":
		return 0x1d555 // 𝕕 '\U0001d555'
	case "Bbbe":
		return 0x1d556 // 𝕖 '\U0001d556'
	case "Bbbeight":
		return 0x1d7e0 // 𝟠 '\U0001d7e0'
	case "Bbbf":
		return 0x1d557 // 𝕗 '\U0001d557'
	case "Bbbfive":
		return 0x1d7dd // 𝟝 '\U0001d7dd'
	case "Bbbfour":
		return 0x1d7dc // 𝟜 '\U0001d7dc'
	case "Bbbg":
		return 0x1d558 // 𝕘 '\U0001d558'
	case "Bbbgamma":
		return 0x213d // ℽ '\u213d'
	case "Bbbh":
		return 0x1d559 // 𝕙 '\U0001d559'
	case "Bbbi":
		return 0x1d55a // 𝕚 '\U0001d55a'
	case "Bbbj":
		return 0x1d55b // 𝕛 '\U0001d55b'
	case "Bbbk":
		return 0x1d55c // 𝕜 '\U0001d55c'
	case "Bbbl":
		return 0x1d55d // 𝕝 '\U0001d55d'
	case "Bbbm":
		return 0x1d55e // 𝕞 '\U0001d55e'
	case "Bbbn":
		return 0x1d55f // 𝕟 '\U0001d55f'
	case "Bbbnine":
		return 0x1d7e1 // 𝟡 '\U0001d7e1'
	case "Bbbo":
		return 0x1d560 // 𝕠 '\U0001d560'
	case "Bbbone":
		return 0x1d7d9 // 𝟙 '\U0001d7d9'
	case "Bbbp":
		return 0x1d561 // 𝕡 '\U0001d561'
	case "Bbbpi":
		return 0x213c // ℼ '\u213c'
	case "Bbbq":
		return 0x1d562 // 𝕢 '\U0001d562'
	case "Bbbr":
		return 0x1d563 // 𝕣 '\U0001d563'
	case "Bbbs":
		return 0x1d564 // 𝕤 '\U0001d564'
	case "Bbbseven":
		return 0x1d7df // 𝟟 '\U0001d7df'
	case "Bbbsix":
		return 0x1d7de // 𝟞 '\U0001d7de'
	case "Bbbsum":
		return 0x2140 // ⅀ '\u2140'
	case "Bbbt":
		return 0x1d565 // 𝕥 '\U0001d565'
	case "Bbbthree":
		return 0x1d7db // 𝟛 '\U0001d7db'
	case "Bbbtwo":
		return 0x1d7da // 𝟚 '\U0001d7da'
	case "Bbbu":
		return 0x1d566 // 𝕦 '\U0001d566'
	case "Bbbv":
		return 0x1d567 // 𝕧 '\U0001d567'
	case "Bbbw":
		return 0x1d568 // 𝕨 '\U0001d568'
	case "Bbbx":
		return 0x1d569 // 𝕩 '\U0001d569'
	case "Bbby":
		return 0x1d56a // 𝕪 '\U0001d56a'
	case "Bbbz":
		return 0x1d56b // 𝕫 '\U0001d56b'
	case "Bbbzero":
		return 0x1d7d8 // 𝟘 '\U0001d7d8'
	case "Bcircle":
		return 0x24b7 // Ⓑ '\u24b7'
	case "Bdotaccent":
		return 0x1e02 // Ḃ '\u1e02'
	case "Bdotbelow":
		return 0x1e04 // Ḅ '\u1e04'
	case "Benarmenian":
		return 0x0532 // Բ '\u0532'
	case "Beta":
		return 0x0392 // Β '\u0392'
	case "Bhook":
		return 0x0181 // Ɓ '\u0181'
	case "Blinebelow":
		return 0x1e06 // Ḇ '\u1e06'
	case "Bmonospace":
		return 0xff22 // Ｂ '\uff22'
	case "Brevesmall":
		return 0xf6f4 //  '\uf6f4'
	case "Bsmall":
		return 0xf762 //  '\uf762'
	case "Bsmallcap":
		return 0x0229 // ȩ '\u0229'
	case "Btopbar":
		return 0x0182 // Ƃ '\u0182'
	case "C":
		return 0x0043 // C 'C'
	case "Caarmenian":
		return 0x053e // Ծ '\u053e'
	case "Cacute":
		return 0x0106 // Ć '\u0106'
	case "Caron":
		return 0xf6ca //  '\uf6ca'
	case "Caronsmall":
		return 0xf6f5 //  '\uf6f5'
	case "Ccaron":
		return 0x010c // Č '\u010c'
	case "Ccedilla":
		return 0x00c7 // Ç '\u00c7'
	case "Ccedillaacute":
		return 0x1e08 // Ḉ '\u1e08'
	case "Ccedillasmall":
		return 0xf7e7 //  '\uf7e7'
	case "Ccircle":
		return 0x24b8 // Ⓒ '\u24b8'
	case "Ccircumflex":
		return 0x0108 // Ĉ '\u0108'
	case "Cdotaccent":
		return 0x010a // Ċ '\u010a'
	case "Cedillasmall":
		return 0xf7b8 //  '\uf7b8'
	case "Chaarmenian":
		return 0x0549 // Չ '\u0549'
	case "Cheabkhasiancyrillic":
		return 0x04bc // Ҽ '\u04bc'
	case "Checyrillic":
		return 0x0427 // Ч '\u0427'
	case "Chedescenderabkhasiancyrillic":
		return 0x04be // Ҿ '\u04be'
	case "Chedescendercyrillic":
		return 0x04b6 // Ҷ '\u04b6'
	case "Chedieresiscyrillic":
		return 0x04f4 // Ӵ '\u04f4'
	case "Cheharmenian":
		return 0x0543 // Ճ '\u0543'
	case "Chekhakassiancyrillic":
		return 0x04cb // Ӌ '\u04cb'
	case "Cheverticalstrokecyrillic":
		return 0x04b8 // Ҹ '\u04b8'
	case "Chi":
		return 0x03a7 // Χ '\u03a7'
	case "Chook":
		return 0x0187 // Ƈ '\u0187'
	case "Circumflexsmall":
		return 0xf6f6 //  '\uf6f6'
	case "Cmonospace":
		return 0xff23 // Ｃ '\uff23'
	case "Coarmenian":
		return 0x0551 // Ց '\u0551'
	case "Coloneq":
		return 0x2a74 // ⩴ '\u2a74'
	case "Csmall":
		return 0xf763 //  '\uf763'
	case "D":
		return 0x0044 // D 'D'
	case "DDownarrow":
		return 0x27f1 // ⟱ '\u27f1'
	case "DZ":
		return 0x01f1 // Ǳ '\u01f1'
	case "DZcaron":
		return 0x01c4 // Ǆ '\u01c4'
	case "Daarmenian":
		return 0x0534 // Դ '\u0534'
	case "Dafrican":
		return 0x0189 // Ɖ '\u0189'
	case "DashV":
		return 0x2ae5 // ⫥ '\u2ae5'
	case "DashVDash":
		return 0x27da // ⟚ '\u27da'
	case "Dashv":
		return 0x2ae4 // ⫤ '\u2ae4'
	case "Dcaron":
		return 0x010e // Ď '\u010e'
	case "Dcaron1":
		return 0xf810 //  '\uf810'
	case "Dcedilla":
		return 0x1e10 // Ḑ '\u1e10'
	case "Dcircle":
		return 0x24b9 // Ⓓ '\u24b9'
	case "Dcircumflexbelow":
		return 0x1e12 // Ḓ '\u1e12'
	case "Dcroat":
		return 0x0110 // Đ '\u0110'
	case "Ddotaccent":
		return 0x1e0a // Ḋ '\u1e0a'
	case "Ddotbelow":
		return 0x1e0c // Ḍ '\u1e0c'
	case "Ddownarrow":
		return 0x290b // ⤋ '\u290b'
	case "Decyrillic":
		return 0x0414 // Д '\u0414'
	case "Deicoptic":
		return 0x03ee // Ϯ '\u03ee'
	case "Delta":
		return 0x2206 // ∆ '\u2206'
	case "Deltagreek":
		return 0x0394 // Δ '\u0394'
	case "Dhook":
		return 0x018a // Ɗ '\u018a'
	case "Dieresis":
		return 0xf6cb //  '\uf6cb'
	case "DieresisAcute":
		return 0xf6cc //  '\uf6cc'
	case "DieresisGrave":
		return 0xf6cd //  '\uf6cd'
	case "Dieresissmall":
		return 0xf7a8 //  '\uf7a8'
	case "Digamma":
		return 0x1d7cb // 𝟋 '\U0001d7cb'
	case "Digammagreek":
		return 0x03dc // Ϝ '\u03dc'
	case "Dlinebelow":
		return 0x1e0e // Ḏ '\u1e0e'
	case "Dmonospace":
		return 0xff24 // Ｄ '\uff24'
	case "Dotaccentsmall":
		return 0xf6f7 //  '\uf6f7'
	case "Dsmall":
		return 0xf764 //  '\uf764'
	case "Dtopbar":
		return 0x018b // Ƌ '\u018b'
	case "Dz":
		return 0x01f2 // ǲ '\u01f2'
	case "Dzcaron":
		return 0x01c5 // ǅ '\u01c5'
	case "Dzeabkhasiancyrillic":
		return 0x04e0 // Ӡ '\u04e0'
	case "Dzhecyrillic":
		return 0x040f // Џ '\u040f'
	case "E":
		return 0x0045 // E 'E'
	case "Eacute":
		return 0x00c9 // É '\u00c9'
	case "Eacutesmall":
		return 0xf7e9 //  '\uf7e9'
	case "Ebreve":
		return 0x0114 // Ĕ '\u0114'
	case "Ecaron":
		return 0x011a // Ě '\u011a'
	case "Ecedillabreve":
		return 0x1e1c // Ḝ '\u1e1c'
	case "Echarmenian":
		return 0x0535 // Ե '\u0535'
	case "Ecircle":
		return 0x24ba // Ⓔ '\u24ba'
	case "Ecircumflex":
		return 0x00ca // Ê '\u00ca'
	case "Ecircumflexacute":
		return 0x1ebe // Ế '\u1ebe'
	case "Ecircumflexbelow":
		return 0x1e18 // Ḙ '\u1e18'
	case "Ecircumflexdotbelow":
		return 0x1ec6 // Ệ '\u1ec6'
	case "Ecircumflexgrave":
		return 0x1ec0 // Ề '\u1ec0'
	case "Ecircumflexhookabove":
		return 0x1ec2 // Ể '\u1ec2'
	case "Ecircumflexsmall":
		return 0xf7ea //  '\uf7ea'
	case "Ecircumflextilde":
		return 0x1ec4 // Ễ '\u1ec4'
	case "Ecyrillic":
		return 0x0404 // Є '\u0404'
	case "Edblgrave":
		return 0x0204 // Ȅ '\u0204'
	case "Edieresis":
		return 0x00cb // Ë '\u00cb'
	case "Edieresissmall":
		return 0xf7eb //  '\uf7eb'
	case "Edotaccent":
		return 0x0116 // Ė '\u0116'
	case "Edotbelow":
		return 0x1eb8 // Ẹ '\u1eb8'
	case "Egrave":
		return 0x00c8 // È '\u00c8'
	case "Egravesmall":
		return 0xf7e8 //  '\uf7e8'
	case "Eharmenian":
		return 0x0537 // Է '\u0537'
	case "Ehookabove":
		return 0x1eba // Ẻ '\u1eba'
	case "Eightroman":
		return 0x2167 // Ⅷ '\u2167'
	case "Einvertedbreve":
		return 0x0206 // Ȇ '\u0206'
	case "Eiotifiedcyrillic":
		return 0x0464 // Ѥ '\u0464'
	case "Elcyrillic":
		return 0x041b // Л '\u041b'
	case "Elevenroman":
		return 0x216a // Ⅺ '\u216a'
	case "Emacron":
		return 0x0112 // Ē '\u0112'
	case "Emacronacute":
		return 0x1e16 // Ḗ '\u1e16'
	case "Emacrongrave":
		return 0x1e14 // Ḕ '\u1e14'
	case "Emcyrillic":
		return 0x041c // М '\u041c'
	case "Emonospace":
		return 0xff25 // Ｅ '\uff25'
	case "Endescendercyrillic":
		return 0x04a2 // Ң '\u04a2'
	case "Eng":
		return 0x014a // Ŋ '\u014a'
	case "Enghecyrillic":
		return 0x04a4 // Ҥ '\u04a4'
	case "Enhookcyrillic":
		return 0x04c7 // Ӈ '\u04c7'
	case "Eogonek":
		return 0x0118 // Ę '\u0118'
	case "Eopen":
		return 0x0190 // Ɛ '\u0190'
	case "Epsilon":
		return 0x0395 // Ε '\u0395'
	case "Epsilontonos":
		return 0x0388 // Έ '\u0388'
	case "Equiv":
		return 0x2263 // ≣ '\u2263'
	case "Ereversed":
		return 0x018e // Ǝ '\u018e'
	case "Ereversedcyrillic":
		return 0x042d // Э '\u042d'
	case "Esdescendercyrillic":
		return 0x04aa // Ҫ '\u04aa'
	case "Esh":
		return 0x01a9 // Ʃ '\u01a9'
	case "Esmall":
		return 0xf765 //  '\uf765'
	case "Eta":
		return 0x0397 // Η '\u0397'
	case "Etarmenian":
		return 0x0538 // Ը '\u0538'
	case "Etatonos":
		return 0x0389 // Ή '\u0389'
	case "Eth":
		return 0x00d0 // Ð '\u00d0'
	case "Ethsmall":
		return 0xf7f0 //  '\uf7f0'
	case "Etilde":
		return 0x1ebc // Ẽ '\u1ebc'
	case "Etildebelow":
		return 0x1e1a // Ḛ '\u1e1a'
	case "Eulerconst":
		return 0x2107 // ℇ '\u2107'
	case "Euro":
		return 0x20ac // € '\u20ac'
	case "Ezh":
		return 0x01b7 // Ʒ '\u01b7'
	case "Ezhcaron":
		return 0x01ee // Ǯ '\u01ee'
	case "Ezhreversed":
		return 0x01b8 // Ƹ '\u01b8'
	case "F":
		return 0x0046 // F 'F'
	case "Fcircle":
		return 0x24bb // Ⓕ '\u24bb'
	case "Fdotaccent":
		return 0x1e1e // Ḟ '\u1e1e'
	case "Feharmenian":
		return 0x0556 // Ֆ '\u0556'
	case "Feicoptic":
		return 0x03e4 // Ϥ '\u03e4'
	case "Fhook":
		return 0x0191 // Ƒ '\u0191'
	case "Finv":
		return 0x2132 // Ⅎ '\u2132'
	case "Fiveroman":
		return 0x2164 // Ⅴ '\u2164'
	case "Fmonospace":
		return 0xff26 // Ｆ '\uff26'
	case "Fourroman":
		return 0x2163 // Ⅳ '\u2163'
	case "Fsmall":
		return 0xf766 //  '\uf766'
	case "G":
		return 0x0047 // G 'G'
	case "GBsquare":
		return 0x3387 // ㎇ '\u3387'
	case "Gacute":
		return 0x01f4 // Ǵ '\u01f4'
	case "Gamma":
		return 0x0393 // Γ '\u0393'
	case "Gammaafrican":
		return 0x0194 // Ɣ '\u0194'
	case "Gangiacoptic":
		return 0x03ea // Ϫ '\u03ea'
	case "Gbreve":
		return 0x011e // Ğ '\u011e'
	case "Gcaron":
		return 0x01e6 // Ǧ '\u01e6'
	case "Gcircle":
		return 0x24bc // Ⓖ '\u24bc'
	case "Gcircumflex":
		return 0x011c // Ĝ '\u011c'
	case "Gcommaaccent":
		return 0x0122 // Ģ '\u0122'
	case "Gdotaccent":
		return 0x0120 // Ġ '\u0120'
	case "Gecyrillic":
		return 0x0413 // Г '\u0413'
	case "Ghadarmenian":
		return 0x0542 // Ղ '\u0542'
	case "Ghemiddlehookcyrillic":
		return 0x0494 // Ҕ '\u0494'
	case "Ghestrokecyrillic":
		return 0x0492 // Ғ '\u0492'
	case "Gheupturncyrillic":
		return 0x0490 // Ґ '\u0490'
	case "Ghook":
		return 0x0193 // Ɠ '\u0193'
	case "Gimarmenian":
		return 0x0533 // Գ '\u0533'
	case "Gmacron":
		return 0x1e20 // Ḡ '\u1e20'
	case "Gmir":
		return 0x2141 // ⅁ '\u2141'
	case "Gmonospace":
		return 0xff27 // Ｇ '\uff27'
	case "Grave":
		return 0xf6ce //  '\uf6ce'
	case "Gravesmall":
		return 0xf760 //  '\uf760'
	case "Gsmall":
		return 0xf767 //  '\uf767'
	case "Gsmallcaphook":
		return 0x022b // ȫ '\u022b'
	case "Gsmallhook":
		return 0x029b // ʛ '\u029b'
	case "Gstroke":
		return 0x01e4 // Ǥ '\u01e4'
	case "Gt":
		return 0x2aa2 // ⪢ '\u2aa2'
	case "H":
		return 0x0048 // H 'H'
	case "H22073":
		return 0x25a1 // □ '\u25a1'
	case "HPsquare":
		return 0x33cb // ㏋ '\u33cb'
	case "Haabkhasiancyrillic":
		return 0x04a8 // Ҩ '\u04a8'
	case "Hadescendercyrillic":
		return 0x04b2 // Ҳ '\u04b2'
	case "Hbar":
		return 0x0126 // Ħ '\u0126'
	case "Hbrevebelow":
		return 0x1e2a // Ḫ '\u1e2a'
	case "Hcedilla":
		return 0x1e28 // Ḩ '\u1e28'
	case "Hcircle":
		return 0x24bd // Ⓗ '\u24bd'
	case "Hcircumflex":
		return 0x0124 // Ĥ '\u0124'
	case "Hdieresis":
		return 0x1e26 // Ḧ '\u1e26'
	case "Hdotaccent":
		return 0x1e22 // Ḣ '\u1e22'
	case "Hdotbelow":
		return 0x1e24 // Ḥ '\u1e24'
	case "Hermaphrodite":
		return 0x26a5 // ⚥ '\u26a5'
	case "Hmonospace":
		return 0xff28 // Ｈ '\uff28'
	case "Hoarmenian":
		return 0x0540 // Հ '\u0540'
	case "Horicoptic":
		return 0x03e8 // Ϩ '\u03e8'
	case "Hsmall":
		return 0xf768 //  '\uf768'
	case "Hsmallcap":
		return 0x022c // Ȭ '\u022c'
	case "Hungarumlaut":
		return 0xf6cf //  '\uf6cf'
	case "Hungarumlautsmall":
		return 0xf6f8 //  '\uf6f8'
	case "Hzsquare":
		return 0x3390 // ㎐ '\u3390'
	case "I":
		return 0x0049 // I 'I'
	case "IJ":
		return 0x0132 // Ĳ '\u0132'
	case "Iacute":
		return 0x00cd // Í '\u00cd'
	case "Iacutesmall":
		return 0xf7ed //  '\uf7ed'
	case "Ibreve":
		return 0x012c // Ĭ '\u012c'
	case "Icaron":
		return 0x01cf // Ǐ '\u01cf'
	case "Icircle":
		return 0x24be // Ⓘ '\u24be'
	case "Icircumflex":
		return 0x00ce // Î '\u00ce'
	case "Icircumflexsmall":
		return 0xf7ee //  '\uf7ee'
	case "Icyril1":
		return 0x03fc // ϼ '\u03fc'
	case "Idblgrave":
		return 0x0208 // Ȉ '\u0208'
	case "Idieresis":
		return 0x00cf // Ï '\u00cf'
	case "Idieresisacute":
		return 0x1e2e // Ḯ '\u1e2e'
	case "Idieresiscyrillic":
		return 0x04e4 // Ӥ '\u04e4'
	case "Idieresissmall":
		return 0xf7ef //  '\uf7ef'
	case "Idot":
		return 0x0130 // İ '\u0130'
	case "Idotbelow":
		return 0x1eca // Ị '\u1eca'
	case "Iebrevecyrillic":
		return 0x04d6 // Ӗ '\u04d6'
	case "Iecyrillic":
		return 0x0415 // Е '\u0415'
	case "Iehook":
		return 0x03f8 // ϸ '\u03f8'
	case "Iehookogonek":
		return 0x03fa // Ϻ '\u03fa'
	case "Ifraktur":
		return 0x2111 // ℑ '\u2111'
	case "Igrave":
		return 0x00cc // Ì '\u00cc'
	case "Igravesmall":
		return 0xf7ec //  '\uf7ec'
	case "Ihookabove":
		return 0x1ec8 // Ỉ '\u1ec8'
	case "Iicyrillic":
		return 0x0418 // И '\u0418'
	case "Iinvertedbreve":
		return 0x020a // Ȋ '\u020a'
	case "Imacron":
		return 0x012a // Ī '\u012a'
	case "Imacroncyrillic":
		return 0x04e2 // Ӣ '\u04e2'
	case "Imonospace":
		return 0xff29 // Ｉ '\uff29'
	case "Iniarmenian":
		return 0x053b // Ի '\u053b'
	case "Iocyrillic":
		return 0x0401 // Ё '\u0401'
	case "Iogonek":
		return 0x012e // Į '\u012e'
	case "Iota":
		return 0x0399 // Ι '\u0399'
	case "Iotaafrican":
		return 0x0196 // Ɩ '\u0196'
	case "Iotadiaeresis":
		return 0x02f3 // ˳ '\u02f3'
	case "Iotadieresis":
		return 0x03aa // Ϊ '\u03aa'
	case "Iotatonos":
		return 0x038a // Ί '\u038a'
	case "Ismall":
		return 0xf769 //  '\uf769'
	case "Istroke":
		return 0x0197 // Ɨ '\u0197'
	case "Itilde":
		return 0x0128 // Ĩ '\u0128'
	case "Itildebelow":
		return 0x1e2c // Ḭ '\u1e2c'
	case "Izhitsadblgravecyrillic":
		return 0x0476 // Ѷ '\u0476'
	case "J":
		return 0x004a // J 'J'
	case "Jaarmenian":
		return 0x0541 // Ձ '\u0541'
	case "Jcircle":
		return 0x24bf // Ⓙ '\u24bf'
	case "Jcircumflex":
		return 0x0134 // Ĵ '\u0134'
	case "Jheharmenian":
		return 0x054b // Ջ '\u054b'
	case "Jmonospace":
		return 0xff2a // Ｊ '\uff2a'
	case "Join":
		return 0x2a1d // ⨝ '\u2a1d'
	case "Jsmall":
		return 0xf76a //  '\uf76a'
	case "K":
		return 0x004b // K 'K'
	case "KBsquare":
		return 0x3385 // ㎅ '\u3385'
	case "KKsquare":
		return 0x33cd // ㏍ '\u33cd'
	case "Kabashkircyrillic":
		return 0x04a0 // Ҡ '\u04a0'
	case "Kacute":
		return 0x1e30 // Ḱ '\u1e30'
	case "Kadescendercyrillic":
		return 0x049a // Қ '\u049a'
	case "Kahook":
		return 0x03ff // Ͽ '\u03ff'
	case "Kahookcyrillic":
		return 0x04c3 // Ӄ '\u04c3'
	case "Kappa":
		return 0x039a // Κ '\u039a'
	case "Kastrokecyrillic":
		return 0x049e // Ҟ '\u049e'
	case "Kaverticalstrokecyrillic":
		return 0x049c // Ҝ '\u049c'
	case "Kcaron":
		return 0x01e8 // Ǩ '\u01e8'
	case "Kcedilla":
		return 0x0136 // Ķ '\u0136'
	case "Kcircle":
		return 0x24c0 // Ⓚ '\u24c0'
	case "Kdotbelow":
		return 0x1e32 // Ḳ '\u1e32'
	case "Keharmenian":
		return 0x0554 // Ք '\u0554'
	case "Kenarmenian":
		return 0x053f // Կ '\u053f'
	case "Khacyrillic":
		return 0x0425 // Х '\u0425'
	case "Kheicoptic":
		return 0x03e6 // Ϧ '\u03e6'
	case "Khook":
		return 0x0198 // Ƙ '\u0198'
	case "Kjecyrillic":
		return 0x040c // Ќ '\u040c'
	case "Klinebelow":
		return 0x1e34 // Ḵ '\u1e34'
	case "Kmonospace":
		return 0xff2b // Ｋ '\uff2b'
	case "Koppacyrillic":
		return 0x0480 // Ҁ '\u0480'
	case "Koppagreek":
		return 0x03de // Ϟ '\u03de'
	case "Ksicyrillic":
		return 0x046e // Ѯ '\u046e'
	case "Ksmall":
		return 0xf76b //  '\uf76b'
	case "L":
		return 0x004c // L 'L'
	case "LJ":
		return 0x01c7 // Ǉ '\u01c7'
	case "LL":
		return 0xf6bf //  '\uf6bf'
	case "LLeftarrow":
		return 0x2b45 // ⭅ '\u2b45'
	case "Lacute":
		return 0x0139 // Ĺ '\u0139'
	case "Lambda":
		return 0x039b // Λ '\u039b'
	case "Lbrbrak":
		return 0x27ec // ⟬ '\u27ec'
	case "Lcaron":
		return 0x013d // Ľ '\u013d'
	case "Lcaron1":
		return 0xf812 //  '\uf812'
	case "Lcedilla":
		return 0x013b // Ļ '\u013b'
	case "Lcedilla1":
		return 0xf81a //  '\uf81a'
	case "Lcircle":
		return 0x24c1 // Ⓛ '\u24c1'
	case "Lcircumflexbelow":
		return 0x1e3c // Ḽ '\u1e3c'
	case "Ldotaccent":
		return 0x013f // Ŀ '\u013f'
	case "Ldotbelow":
		return 0x1e36 // Ḷ '\u1e36'
	case "Ldotbelowmacron":
		return 0x1e38 // Ḹ '\u1e38'
	case "Ldsh":
		return 0x21b2 // ↲ '\u21b2'
	case "Liwnarmenian":
		return 0x053c // Լ '\u053c'
	case "Lj":
		return 0x01c8 // ǈ '\u01c8'
	case "Ljecyrillic":
		return 0x0409 // Љ '\u0409'
	case "Llinebelow":
		return 0x1e3a // Ḻ '\u1e3a'
	case "Lmonospace":
		return 0xff2c // Ｌ '\uff2c'
	case "Longleftarrow":
		return 0x27f8 // ⟸ '\u27f8'
	case "Longleftrightarrow":
		return 0x27fa // ⟺ '\u27fa'
	case "Longmapsfrom":
		return 0x27fd // ⟽ '\u27fd'
	case "Longmapsto":
		return 0x27fe // ⟾ '\u27fe'
	case "Longrightarrow":
		return 0x27f9 // ⟹ '\u27f9'
	case "Lparengtr":
		return 0x2995 // ⦕ '\u2995'
	case "Lslash":
		return 0x0141 // Ł '\u0141'
	case "Lslashsmall":
		return 0xf6f9 //  '\uf6f9'
	case "Lsmall":
		return 0xf76c //  '\uf76c'
	case "Lsmallcap":
		return 0x022f // ȯ '\u022f'
	case "Lt":
		return 0x2aa1 // ⪡ '\u2aa1'
	case "Lvzigzag":
		return 0x29da // ⧚ '\u29da'
	case "M":
		return 0x004d // M 'M'
	case "MBsquare":
		return 0x3386 // ㎆ '\u3386'
	case "Macron":
		return 0xf6d0 //  '\uf6d0'
	case "Macronsmall":
		return 0xf7af //  '\uf7af'
	case "Macute":
		return 0x1e3e // Ḿ '\u1e3e'
	case "Mapsfrom":
		return 0x2906 // ⤆ '\u2906'
	case "Mapsto":
		return 0x2907 // ⤇ '\u2907'
	case "Mcircle":
		return 0x24c2 // Ⓜ '\u24c2'
	case "Mdotaccent":
		return 0x1e40 // Ṁ '\u1e40'
	case "Mdotbelow":
		return 0x1e42 // Ṃ '\u1e42'
	case "Menarmenian":
		return 0x0544 // Մ '\u0544'
	case "Mmonospace":
		return 0xff2d // Ｍ '\uff2d'
	case "Msmall":
		return 0xf76d //  '\uf76d'
	case "Mturned":
		return 0x019c // Ɯ '\u019c'
	case "Mu":
		return 0x039c // Μ '\u039c'
	case "N":
		return 0x004e // N 'N'
	case "NJ":
		return 0x01ca // Ǌ '\u01ca'
	case "Nacute":
		return 0x0143 // Ń '\u0143'
	case "Ncaron":
		return 0x0147 // Ň '\u0147'
	case "Ncedilla1":
		return 0xf81c //  '\uf81c'
	case "Ncircle":
		return 0x24c3 // Ⓝ '\u24c3'
	case "Ncircumflexbelow":
		return 0x1e4a // Ṋ '\u1e4a'
	case "Ncommaaccent":
		return 0x0145 // Ņ '\u0145'
	case "Ndotaccent":
		return 0x1e44 // Ṅ '\u1e44'
	case "Ndotbelow":
		return 0x1e46 // Ṇ '\u1e46'
	case "Nearrow":
		return 0x21d7 // ⇗ '\u21d7'
	case "Nhookleft":
		return 0x019d // Ɲ '\u019d'
	case "Nineroman":
		return 0x2168 // Ⅸ '\u2168'
	case "Nj":
		return 0x01cb // ǋ '\u01cb'
	case "Nlinebelow":
		return 0x1e48 // Ṉ '\u1e48'
	case "Nmonospace":
		return 0xff2e // Ｎ '\uff2e'
	case "Not":
		return 0x2aec // ⫬ '\u2aec'
	case "Nowarmenian":
		return 0x0546 // Ն '\u0546'
	case "Nsmall":
		return 0xf76e //  '\uf76e'
	case "Ntilde":
		return 0x00d1 // Ñ '\u00d1'
	case "Ntildesmall":
		return 0xf7f1 //  '\uf7f1'
	case "Nu":
		return 0x039d // Ν '\u039d'
	case "Nwarrow":
		return 0x21d6 // ⇖ '\u21d6'
	case "O":
		return 0x004f // O 'O'
	case "OE":
		return 0x0152 // Œ '\u0152'
	case "OEsmall":
		return 0xf6fa //  '\uf6fa'
	case "Oacute":
		return 0x00d3 // Ó '\u00d3'
	case "Oacutesmall":
		return 0xf7f3 //  '\uf7f3'
	case "Obarredcyrillic":
		return 0x04e8 // Ө '\u04e8'
	case "Obarreddieresiscyrillic":
		return 0x04ea // Ӫ '\u04ea'
	case "Obreve":
		return 0x014e // Ŏ '\u014e'
	case "Ocaron":
		return 0x01d1 // Ǒ '\u01d1'
	case "Ocenteredtilde":
		return 0x019f // Ɵ '\u019f'
	case "Ocircle":
		return 0x24c4 // Ⓞ '\u24c4'
	case "Ocircumflex":
		return 0x00d4 // Ô '\u00d4'
	case "Ocircumflexacute":
		return 0x1ed0 // Ố '\u1ed0'
	case "Ocircumflexdotbelow":
		return 0x1ed8 // Ộ '\u1ed8'
	case "Ocircumflexgrave":
		return 0x1ed2 // Ồ '\u1ed2'
	case "Ocircumflexhookabove":
		return 0x1ed4 // Ổ '\u1ed4'
	case "Ocircumflexsmall":
		return 0xf7f4 //  '\uf7f4'
	case "Ocircumflextilde":
		return 0x1ed6 // Ỗ '\u1ed6'
	case "Ocyrillic":
		return 0x041e // О '\u041e'
	case "Odblacute":
		return 0x0150 // Ő '\u0150'
	case "Odblgrave":
		return 0x020c // Ȍ '\u020c'
	case "Odieresis":
		return 0x00d6 // Ö '\u00d6'
	case "Odieresiscyrillic":
		return 0x04e6 // Ӧ '\u04e6'
	case "Odieresissmall":
		return 0xf7f6 //  '\uf7f6'
	case "Odotbelow":
		return 0x1ecc // Ọ '\u1ecc'
	case "Ogoneksmall":
		return 0xf6fb //  '\uf6fb'
	case "Ograve":
		return 0x00d2 // Ò '\u00d2'
	case "Ogravesmall":
		return 0xf7f2 //  '\uf7f2'
	case "Oharmenian":
		return 0x0555 // Օ '\u0555'
	case "Ohookabove":
		return 0x1ece // Ỏ '\u1ece'
	case "Ohorn":
		return 0x01a0 // Ơ '\u01a0'
	case "Ohornacute":
		return 0x1eda // Ớ '\u1eda'
	case "Ohorndotbelow":
		return 0x1ee2 // Ợ '\u1ee2'
	case "Ohorngrave":
		return 0x1edc // Ờ '\u1edc'
	case "Ohornhookabove":
		return 0x1ede // Ở '\u1ede'
	case "Ohorntilde":
		return 0x1ee0 // Ỡ '\u1ee0'
	case "Oi":
		return 0x01a2 // Ƣ '\u01a2'
	case "Oinvertedbreve":
		return 0x020e // Ȏ '\u020e'
	case "Omacron":
		return 0x014c // Ō '\u014c'
	case "Omacronacute":
		return 0x1e52 // Ṓ '\u1e52'
	case "Omacrongrave":
		return 0x1e50 // Ṑ '\u1e50'
	case "Omega":
		return 0x2126 // Ω '\u2126'
	case "Omegacyrillic":
		return 0x0460 // Ѡ '\u0460'
	case "Omegagreek":
		return 0x03a9 // Ω '\u03a9'
	case "Omegainv":
		return 0x2127 // ℧ '\u2127'
	case "Omegaroundcyrillic":
		return 0x047a // Ѻ '\u047a'
	case "Omegatitlocyrillic":
		return 0x047c // Ѽ '\u047c'
	case "Omegatonos":
		return 0x038f // Ώ '\u038f'
	case "Omicron":
		return 0x039f // Ο '\u039f'
	case "Omicrontonos":
		return 0x038c // Ό '\u038c'
	case "Omonospace":
		return 0xff2f // Ｏ '\uff2f'
	case "Oneroman":
		return 0x2160 // Ⅰ '\u2160'
	case "Oogonek":
		return 0x01ea // Ǫ '\u01ea'
	case "Oogonekmacron":
		return 0x01ec // Ǭ '\u01ec'
	case "Oopen":
		return 0x0186 // Ɔ '\u0186'
	case "Oslash":
		return 0x00d8 // Ø '\u00d8'
	case "Oslashacute":
		return 0x01fe // Ǿ '\u01fe'
	case "Oslashsmall":
		return 0xf7f8 //  '\uf7f8'
	case "Osmall":
		return 0xf76f //  '\uf76f'
	case "Otcyrillic":
		return 0x047e // Ѿ '\u047e'
	case "Otilde":
		return 0x00d5 // Õ '\u00d5'
	case "Otildeacute":
		return 0x1e4c // Ṍ '\u1e4c'
	case "Otildedieresis":
		return 0x1e4e // Ṏ '\u1e4e'
	case "Otildesmall":
		return 0xf7f5 //  '\uf7f5'
	case "Otimes":
		return 0x2a37 // ⨷ '\u2a37'
	case "P":
		return 0x0050 // P 'P'
	case "Pacute":
		return 0x1e54 // Ṕ '\u1e54'
	case "Pcircle":
		return 0x24c5 // Ⓟ '\u24c5'
	case "Pdotaccent":
		return 0x1e56 // Ṗ '\u1e56'
	case "Peharmenian":
		return 0x054a // Պ '\u054a'
	case "Pemiddlehookcyrillic":
		return 0x04a6 // Ҧ '\u04a6'
	case "Phi":
		return 0x03a6 // Φ '\u03a6'
	case "Phook":
		return 0x01a4 // Ƥ '\u01a4'
	case "Pi":
		return 0x03a0 // Π '\u03a0'
	case "Piwrarmenian":
		return 0x0553 // Փ '\u0553'
	case "Planckconst":
		return 0x210e // ℎ '\u210e'
	case "Pmonospace":
		return 0xff30 // Ｐ '\uff30'
	case "Prec":
		return 0x2abb // ⪻ '\u2abb'
	case "PropertyLine":
		return 0x214a // ⅊ '\u214a'
	case "Psi":
		return 0x03a8 // Ψ '\u03a8'
	case "Psicyrillic":
		return 0x0470 // Ѱ '\u0470'
	case "Psmall":
		return 0xf770 //  '\uf770'
	case "Q":
		return 0x0051 // Q 'Q'
	case "QED":
		return 0x220e // ∎ '\u220e'
	case "Qcircle":
		return 0x24c6 // Ⓠ '\u24c6'
	case "Qmonospace":
		return 0xff31 // Ｑ '\uff31'
	case "Qsmall":
		return 0xf771 //  '\uf771'
	case "Question":
		return 0x2047 // ⁇ '\u2047'
	case "R":
		return 0x0052 // R 'R'
	case "RRightarrow":
		return 0x2b46 // ⭆ '\u2b46'
	case "Raarmenian":
		return 0x054c // Ռ '\u054c'
	case "Racute":
		return 0x0154 // Ŕ '\u0154'
	case "Rbrbrak":
		return 0x27ed // ⟭ '\u27ed'
	case "Rcaron":
		return 0x0158 // Ř '\u0158'
	case "Rcedilla":
		return 0x0156 // Ŗ '\u0156'
	case "Rcedilla1":
		return 0xf81e //  '\uf81e'
	case "Rcircle":
		return 0x24c7 // Ⓡ '\u24c7'
	case "Rcircumflex":
		return 0xf831 //  '\uf831'
	case "Rdblgrave":
		return 0x0210 // Ȑ '\u0210'
	case "Rdotaccent":
		return 0x1e58 // Ṙ '\u1e58'
	case "Rdotbelow":
		return 0x1e5a // Ṛ '\u1e5a'
	case "Rdotbelowmacron":
		return 0x1e5c // Ṝ '\u1e5c'
	case "Rdsh":
		return 0x21b3 // ↳ '\u21b3'
	case "Reharmenian":
		return 0x0550 // Ր '\u0550'
	case "Rfraktur":
		return 0x211c // ℜ '\u211c'
	case "Rho":
		return 0x03a1 // Ρ '\u03a1'
	case "Ringsmall":
		return 0xf6fc //  '\uf6fc'
	case "Rinvertedbreve":
		return 0x0212 // Ȓ '\u0212'
	case "Rlinebelow":
		return 0x1e5e // Ṟ '\u1e5e'
	case "Rmonospace":
		return 0xff32 // Ｒ '\uff32'
	case "Rparenless":
		return 0x2996 // ⦖ '\u2996'
	case "Rsmall":
		return 0xf772 //  '\uf772'
	case "Rsmallinverted":
		return 0x0281 // ʁ '\u0281'
	case "Rsmallinvertedsuperior":
		return 0x02b6 // ʶ '\u02b6'
	case "Rturnsuper":
		return 0x023f // ȿ '\u023f'
	case "Rvzigzag":
		return 0x29db // ⧛ '\u29db'
	case "S":
		return 0x0053 // S 'S'
	case "SD150100":
		return 0x024f // ɏ '\u024f'
	case "SF010000":
		return 0x250c // ┌ '\u250c'
	case "SF020000":
		return 0x2514 // └ '\u2514'
	case "SF030000":
		return 0x2510 // ┐ '\u2510'
	case "SF040000":
		return 0x2518 // ┘ '\u2518'
	case "SF050000":
		return 0x253c // ┼ '\u253c'
	case "SF060000":
		return 0x252c // ┬ '\u252c'
	case "SF070000":
		return 0x2534 // ┴ '\u2534'
	case "SF080000":
		return 0x251c // ├ '\u251c'
	case "SF090000":
		return 0x2524 // ┤ '\u2524'
	case "SF100000":
		return 0x2500 // ─ '\u2500'
	case "SF110000":
		return 0x2502 // │ '\u2502'
	case "SF190000":
		return 0x2561 // ╡ '\u2561'
	case "SF200000":
		return 0x2562 // ╢ '\u2562'
	case "SF210000":
		return 0x2556 // ╖ '\u2556'
	case "SF220000":
		return 0x2555 // ╕ '\u2555'
	case "SF230000":
		return 0x2563 // ╣ '\u2563'
	case "SF240000":
		return 0x2551 // ║ '\u2551'
	case "SF250000":
		return 0x2557 // ╗ '\u2557'
	case "SF260000":
		return 0x255d // ╝ '\u255d'
	case "SF270000":
		return 0x255c // ╜ '\u255c'
	case "SF280000":
		return 0x255b // ╛ '\u255b'
	case "SF360000":
		return 0x255e // ╞ '\u255e'
	case "SF370000":
		return 0x255f // ╟ '\u255f'
	case "SF380000":
		return 0x255a // ╚ '\u255a'
	case "SF390000":
		return 0x2554 // ╔ '\u2554'
	case "SF400000":
		return 0x2569 // ╩ '\u2569'
	case "SF410000":
		return 0x2566 // ╦ '\u2566'
	case "SF420000":
		return 0x2560 // ╠ '\u2560'
	case "SF430000":
		return 0x2550 // ═ '\u2550'
	case "SF440000":
		return 0x256c // ╬ '\u256c'
	case "SF450000":
		return 0x2567 // ╧ '\u2567'
	case "SF460000":
		return 0x2568 // ╨ '\u2568'
	case "SF470000":
		return 0x2564 // ╤ '\u2564'
	case "SF480000":
		return 0x2565 // ╥ '\u2565'
	case "SF490000":
		return 0x2559 // ╙ '\u2559'
	case "SF500000":
		return 0x2558 // ╘ '\u2558'
	case "SF510000":
		return 0x2552 // ╒ '\u2552'
	case "SF520000":
		return 0x2553 // ╓ '\u2553'
	case "SF530000":
		return 0x256b // ╫ '\u256b'
	case "SF540000":
		return 0x256a // ╪ '\u256a'
	case "Sacute":
		return 0x015a // Ś '\u015a'
	case "Sacutedotaccent":
		return 0x1e64 // Ṥ '\u1e64'
	case "Sampigreek":
		return 0x03e0 // Ϡ '\u03e0'
	case "Scaron":
		return 0x0160 // Š '\u0160'
	case "Scarondotaccent":
		return 0x1e66 // Ṧ '\u1e66'
	case "Scaronsmall":
		return 0xf6fd //  '\uf6fd'
	case "Scedilla":
		return 0x015e // Ş '\u015e'
	case "Scedilla1":
		return 0xf816 //  '\uf816'
	case "Schwa":
		return 0x018f // Ə '\u018f'
	case "Schwacyrillic":
		return 0x04d8 // Ә '\u04d8'
	case "Schwadieresiscyrillic":
		return 0x04da // Ӛ '\u04da'
	case "Scircle":
		return 0x24c8 // Ⓢ '\u24c8'
	case "Scircumflex":
		return 0x015c // Ŝ '\u015c'
	case "Scommaaccent":
		return 0x0218 // Ș '\u0218'
	case "Sdotaccent":
		return 0x1e60 // Ṡ '\u1e60'
	case "Sdotbelow":
		return 0x1e62 // Ṣ '\u1e62'
	case "Sdotbelowdotaccent":
		return 0x1e68 // Ṩ '\u1e68'
	case "Searrow":
		return 0x21d8 // ⇘ '\u21d8'
	case "Seharmenian":
		return 0x054d // Ս '\u054d'
	case "Sevenroman":
		return 0x2166 // Ⅶ '\u2166'
	case "Shaarmenian":
		return 0x0547 // Շ '\u0547'
	case "Shacyrillic":
		return 0x0428 // Ш '\u0428'
	case "Sheicoptic":
		return 0x03e2 // Ϣ '\u03e2'
	case "Shhacyrillic":
		return 0x04ba // Һ '\u04ba'
	case "Shimacoptic":
		return 0x03ec // Ϭ '\u03ec'
	case "Sigma":
		return 0x03a3 // Σ '\u03a3'
	case "Sixroman":
		return 0x2165 // Ⅵ '\u2165'
	case "Smonospace":
		return 0xff33 // Ｓ '\uff33'
	case "Sqcap":
		return 0x2a4e // ⩎ '\u2a4e'
	case "Sqcup":
		return 0x2a4f // ⩏ '\u2a4f'
	case "Ssmall":
		return 0xf773 //  '\uf773'
	case "Stigmagreek":
		return 0x03da // Ϛ '\u03da'
	case "Succ":
		return 0x2abc // ⪼ '\u2abc'
	case "Swarrow":
		return 0x21d9 // ⇙ '\u21d9'
	case "T":
		return 0x0054 // T 'T'
	case "Tau":
		return 0x03a4 // Τ '\u03a4'
	case "Tbar":
		return 0x0166 // Ŧ '\u0166'
	case "Tcaron":
		return 0x0164 // Ť '\u0164'
	case "Tcaron1":
		return 0xf814 //  '\uf814'
	case "Tcedilla1":
		return 0xf818 //  '\uf818'
	case "Tcircle":
		return 0x24c9 // Ⓣ '\u24c9'
	case "Tcircumflexbelow":
		return 0x1e70 // Ṱ '\u1e70'
	case "Tcommaaccent":
		return 0x0162 // Ţ '\u0162'
	case "Tdotaccent":
		return 0x1e6a // Ṫ '\u1e6a'
	case "Tdotbelow":
		return 0x1e6c // Ṭ '\u1e6c'
	case "Tedescendercyrillic":
		return 0x04ac // Ҭ '\u04ac'
	case "Tenroman":
		return 0x2169 // Ⅹ '\u2169'
	case "Tetsecyrillic":
		return 0x04b4 // Ҵ '\u04b4'
	case "Theta":
		return 0x0398 // Θ '\u0398'
	case "Thook":
		return 0x01ac // Ƭ '\u01ac'
	case "Thorn":
		return 0x00de // Þ '\u00de'
	case "Thornsmall":
		return 0xf7fe //  '\uf7fe'
	case "Threeroman":
		return 0x2162 // Ⅲ '\u2162'
	case "Tildesmall":
		return 0xf6fe //  '\uf6fe'
	case "Tiwnarmenian":
		return 0x054f // Տ '\u054f'
	case "Tlinebelow":
		return 0x1e6e // Ṯ '\u1e6e'
	case "Tmonospace":
		return 0xff34 // Ｔ '\uff34'
	case "Toarmenian":
		return 0x0539 // Թ '\u0539'
	case "Tonefive":
		return 0x01bc // Ƽ '\u01bc'
	case "Tonesix":
		return 0x0184 // Ƅ '\u0184'
	case "Tonetwo":
		return 0x01a7 // Ƨ '\u01a7'
	case "Tretroflexhook":
		return 0x01ae // Ʈ '\u01ae'
	case "Tsecyrillic":
		return 0x0426 // Ц '\u0426'
	case "Tshecyrillic":
		return 0x040b // Ћ '\u040b'
	case "Tsmall":
		return 0xf774 //  '\uf774'
	case "Twelveroman":
		return 0x216b // Ⅻ '\u216b'
	case "Tworoman":
		return 0x2161 // Ⅱ '\u2161'
	case "U":
		return 0x0055 // U 'U'
	case "UUparrow":
		return 0x27f0 // ⟰ '\u27f0'
	case "Uacute":
		return 0x00da // Ú '\u00da'
	case "Uacutesmall":
		return 0xf7fa //  '\uf7fa'
	case "Ubreve":
		return 0x016c // Ŭ '\u016c'
	case "Ucaron":
		return 0x01d3 // Ǔ '\u01d3'
	case "Ucedilla":
		return 0xf833 //  '\uf833'
	case "Ucircle":
		return 0x24ca // Ⓤ '\u24ca'
	case "Ucircumflex":
		return 0x00db // Û '\u00db'
	case "Ucircumflexbelow":
		return 0x1e76 // Ṷ '\u1e76'
	case "Ucircumflexsmall":
		return 0xf7fb //  '\uf7fb'
	case "Ucyrillic":
		return 0x0423 // У '\u0423'
	case "Udblgrave":
		return 0x0214 // Ȕ '\u0214'
	case "Udieresis":
		return 0x00dc // Ü '\u00dc'
	case "Udieresisacute":
		return 0x01d7 // Ǘ '\u01d7'
	case "Udieresisbelow":
		return 0x1e72 // Ṳ '\u1e72'
	case "Udieresiscaron":
		return 0x01d9 // Ǚ '\u01d9'
	case "Udieresiscyrillic":
		return 0x04f0 // Ӱ '\u04f0'
	case "Udieresisgrave":
		return 0x01db // Ǜ '\u01db'
	case "Udieresismacron":
		return 0x01d5 // Ǖ '\u01d5'
	case "Udieresissmall":
		return 0xf7fc //  '\uf7fc'
	case "Udotbelow":
		return 0x1ee4 // Ụ '\u1ee4'
	case "Ugrave":
		return 0x00d9 // Ù '\u00d9'
	case "Ugravesmall":
		return 0xf7f9 //  '\uf7f9'
	case "Uhookabove":
		return 0x1ee6 // Ủ '\u1ee6'
	case "Uhorn":
		return 0x01af // Ư '\u01af'
	case "Uhornacute":
		return 0x1ee8 // Ứ '\u1ee8'
	case "Uhorndotbelow":
		return 0x1ef0 // Ự '\u1ef0'
	case "Uhorngrave":
		return 0x1eea // Ừ '\u1eea'
	case "Uhornhookabove":
		return 0x1eec // Ử '\u1eec'
	case "Uhorntilde":
		return 0x1eee // Ữ '\u1eee'
	case "Uhungarumlaut":
		return 0x0170 // Ű '\u0170'
	case "Uhungarumlautcyrillic":
		return 0x04f2 // Ӳ '\u04f2'
	case "Uinvertedbreve":
		return 0x0216 // Ȗ '\u0216'
	case "Ukcyrillic":
		return 0x0478 // Ѹ '\u0478'
	case "Umacron":
		return 0x016a // Ū '\u016a'
	case "Umacroncyrillic":
		return 0x04ee // Ӯ '\u04ee'
	case "Umacrondieresis":
		return 0x1e7a // Ṻ '\u1e7a'
	case "Umonospace":
		return 0xff35 // Ｕ '\uff35'
	case "Uogonek":
		return 0x0172 // Ų '\u0172'
	case "Upsilon":
		return 0x03a5 // Υ '\u03a5'
	case "Upsilon1":
		return 0x03d2 // ϒ '\u03d2'
	case "Upsilonacutehooksymbolgreek":
		return 0x03d3 // ϓ '\u03d3'
	case "Upsilonafrican":
		return 0x01b1 // Ʊ '\u01b1'
	case "Upsilondiaeresis":
		return 0x02f4 // ˴ '\u02f4'
	case "Upsilondieresis":
		return 0x03ab // Ϋ '\u03ab'
	case "Upsilondieresishooksymbolgreek":
		return 0x03d4 // ϔ '\u03d4'
	case "Upsilontonos":
		return 0x038e // Ύ '\u038e'
	case "Uring":
		return 0x016e // Ů '\u016e'
	case "Ushortcyrillic":
		return 0x040e // Ў '\u040e'
	case "Usmall":
		return 0xf775 //  '\uf775'
	case "Ustraightcyrillic":
		return 0x04ae // Ү '\u04ae'
	case "Ustraightstrokecyrillic":
		return 0x04b0 // Ұ '\u04b0'
	case "Utilde":
		return 0x0168 // Ũ '\u0168'
	case "Utildeacute":
		return 0x1e78 // Ṹ '\u1e78'
	case "Utildebelow":
		return 0x1e74 // Ṵ '\u1e74'
	case "Uuparrow":
		return 0x290a // ⤊ '\u290a'
	case "V":
		return 0x0056 // V 'V'
	case "VDash":
		return 0x22ab // ⊫ '\u22ab'
	case "Vbar":
		return 0x2aeb // ⫫ '\u2aeb'
	case "Vcircle":
		return 0x24cb // Ⓥ '\u24cb'
	case "Vdotbelow":
		return 0x1e7e // Ṿ '\u1e7e'
	case "Vee":
		return 0x2a54 // ⩔ '\u2a54'
	case "Vewarmenian":
		return 0x054e // Վ '\u054e'
	case "Vhook":
		return 0x01b2 // Ʋ '\u01b2'
	case "Vmonospace":
		return 0xff36 // Ｖ '\uff36'
	case "Voarmenian":
		return 0x0548 // Ո '\u0548'
	case "Vsmall":
		return 0xf776 //  '\uf776'
	case "Vtilde":
		return 0x1e7c // Ṽ '\u1e7c'
	case "Vvert":
		return 0x2980 // ⦀ '\u2980'
	case "W":
		return 0x0057 // W 'W'
	case "Wacute":
		return 0x1e82 // Ẃ '\u1e82'
	case "Wcircle":
		return 0x24cc // Ⓦ '\u24cc'
	case "Wcircumflex":
		return 0x0174 // Ŵ '\u0174'
	case "Wdieresis":
		return 0x1e84 // Ẅ '\u1e84'
	case "Wdotaccent":
		return 0x1e86 // Ẇ '\u1e86'
	case "Wdotbelow":
		return 0x1e88 // Ẉ '\u1e88'
	case "Wedge":
		return 0x2a53 // ⩓ '\u2a53'
	case "Wgrave":
		return 0x1e80 // Ẁ '\u1e80'
	case "Wmonospace":
		return 0xff37 // Ｗ '\uff37'
	case "Wsmall":
		return 0xf777 //  '\uf777'
	case "X":
		return 0x0058 // X 'X'
	case "Xcircle":
		return 0x24cd // Ⓧ '\u24cd'
	case "Xdieresis":
		return 0x1e8c // Ẍ '\u1e8c'
	case "Xdotaccent":
		return 0x1e8a // Ẋ '\u1e8a'
	case "Xeharmenian":
		return 0x053d // Խ '\u053d'
	case "Xi":
		return 0x039e // Ξ '\u039e'
	case "Xmonospace":
		return 0xff38 // Ｘ '\uff38'
	case "Xsmall":
		return 0xf778 //  '\uf778'
	case "Y":
		return 0x0059 // Y 'Y'
	case "Yacute":
		return 0x00dd // Ý '\u00dd'
	case "Yacutesmall":
		return 0xf7fd //  '\uf7fd'
	case "Ycircle":
		return 0x24ce // Ⓨ '\u24ce'
	case "Ycircumflex":
		return 0x0176 // Ŷ '\u0176'
	case "Ydieresis":
		return 0x0178 // Ÿ '\u0178'
	case "Ydieresissmall":
		return 0xf7ff //  '\uf7ff'
	case "Ydotaccent":
		return 0x1e8e // Ẏ '\u1e8e'
	case "Ydotbelow":
		return 0x1ef4 // Ỵ '\u1ef4'
	case "Yerudieresiscyrillic":
		return 0x04f8 // Ӹ '\u04f8'
	case "Ygrave":
		return 0x1ef2 // Ỳ '\u1ef2'
	case "Yhook":
		return 0x01b3 // Ƴ '\u01b3'
	case "Yhookabove":
		return 0x1ef6 // Ỷ '\u1ef6'
	case "Yiarmenian":
		return 0x0545 // Յ '\u0545'
	case "Yicyrillic":
		return 0x0407 // Ї '\u0407'
	case "Yiwnarmenian":
		return 0x0552 // Ւ '\u0552'
	case "Ymonospace":
		return 0xff39 // Ｙ '\uff39'
	case "Ysmall":
		return 0xf779 //  '\uf779'
	case "Ysmallcap":
		return 0x021f // ȟ '\u021f'
	case "Ytilde":
		return 0x1ef8 // Ỹ '\u1ef8'
	case "Yup":
		return 0x2144 // ⅄ '\u2144'
	case "Yusbigcyrillic":
		return 0x046a // Ѫ '\u046a'
	case "Yusbigiotifiedcyrillic":
		return 0x046c // Ѭ '\u046c'
	case "Yuslittlecyrillic":
		return 0x0466 // Ѧ '\u0466'
	case "Yuslittleiotifiedcyrillic":
		return 0x0468 // Ѩ '\u0468'
	case "Z":
		return 0x005a // Z 'Z'
	case "Zaarmenian":
		return 0x0536 // Զ '\u0536'
	case "Zacute":
		return 0x0179 // Ź '\u0179'
	case "Zcaron":
		return 0x017d // Ž '\u017d'
	case "Zcaronsmall":
		return 0xf6ff //  '\uf6ff'
	case "Zcircle":
		return 0x24cf // Ⓩ '\u24cf'
	case "Zcircumflex":
		return 0x1e90 // Ẑ '\u1e90'
	case "Zdotaccent":
		return 0x017b // Ż '\u017b'
	case "Zdotbelow":
		return 0x1e92 // Ẓ '\u1e92'
	case "Zedescendercyrillic":
		return 0x0498 // Ҙ '\u0498'
	case "Zedieresiscyrillic":
		return 0x04de // Ӟ '\u04de'
	case "Zeta":
		return 0x0396 // Ζ '\u0396'
	case "Zhearmenian":
		return 0x053a // Ժ '\u053a'
	case "Zhebreve":
		return 0x03fd // Ͻ '\u03fd'
	case "Zhebrevecyrillic":
		return 0x04c1 // Ӂ '\u04c1'
	case "Zhedescendercyrillic":
		return 0x0496 // Җ '\u0496'
	case "Zhedieresiscyrillic":
		return 0x04dc // Ӝ '\u04dc'
	case "Zlinebelow":
		return 0x1e94 // Ẕ '\u1e94'
	case "Zmonospace":
		return 0xff3a // Ｚ '\uff3a'
	case "Zsmall":
		return 0xf77a //  '\uf77a'
	case "Zstroke":
		return 0x01b5 // Ƶ '\u01b5'
	case "a":
		return 0x0061 // a 'a'
	case "a1":
		return 0x2701 // ✁ '\u2701'
	case "a2":
		return 0x2702 // ✂ '\u2702'
	case "a3":
		return 0x2704 // ✄ '\u2704'
	case "a4":
		return 0x260e // ☎ '\u260e'
	case "a5":
		return 0x2706 // ✆ '\u2706'
	case "a6":
		return 0x271d // ✝ '\u271d'
	case "a7":
		return 0x271e // ✞ '\u271e'
	case "a8":
		return 0x271f // ✟ '\u271f'
	case "a9":
		return 0x2720 // ✠ '\u2720'
	case "a10":
		return 0x2721 // ✡ '\u2721'
	case "a11":
		return 0x261b // ☛ '\u261b'
	case "a12":
		return 0x261e // ☞ '\u261e'
	case "a13":
		return 0x270c // ✌ '\u270c'
	case "a14":
		return 0x270d // ✍ '\u270d'
	case "a15":
		return 0x270e // ✎ '\u270e'
	case "a16":
		return 0x270f // ✏ '\u270f'
	case "a17":
		return 0x2711 // ✑ '\u2711'
	case "a18":
		return 0x2712 // ✒ '\u2712'
	case "a19":
		return 0x2713 // ✓ '\u2713'
	case "a20":
		return 0x2714 // ✔ '\u2714'
	case "a21":
		return 0x2715 // ✕ '\u2715'
	case "a22":
		return 0x2716 // ✖ '\u2716'
	case "a23":
		return 0x2717 // ✗ '\u2717'
	case "a24":
		return 0x2718 // ✘ '\u2718'
	case "a25":
		return 0x2719 // ✙ '\u2719'
	case "a26":
		return 0x271a // ✚ '\u271a'
	case "a27":
		return 0x271b // ✛ '\u271b'
	case "a28":
		return 0x271c // ✜ '\u271c'
	case "a29":
		return 0x2722 // ✢ '\u2722'
	case "a30":
		return 0x2723 // ✣ '\u2723'
	case "a31":
		return 0x2724 // ✤ '\u2724'
	case "a32":
		return 0x2725 // ✥ '\u2725'
	case "a33":
		return 0x2726 // ✦ '\u2726'
	case "a34":
		return 0x2727 // ✧ '\u2727'
	case "a35":
		return 0x2605 // ★ '\u2605'
	case "a36":
		return 0x2729 // ✩ '\u2729'
	case "a37":
		return 0x272a // ✪ '\u272a'
	case "a38":
		return 0x272b // ✫ '\u272b'
	case "a39":
		return 0x272c // ✬ '\u272c'
	case "a40":
		return 0x272d // ✭ '\u272d'
	case "a41":
		return 0x272e // ✮ '\u272e'
	case "a42":
		return 0x272f // ✯ '\u272f'
	case "a43":
		return 0x2730 // ✰ '\u2730'
	case "a44":
		return 0x2731 // ✱ '\u2731'
	case "a45":
		return 0x2732 // ✲ '\u2732'
	case "a46":
		return 0x2733 // ✳ '\u2733'
	case "a47":
		return 0x2734 // ✴ '\u2734'
	case "a48":
		return 0x2735 // ✵ '\u2735'
	case "a49":
		return 0x2736 // ✶ '\u2736'
	case "a50":
		return 0x2737 // ✷ '\u2737'
	case "a51":
		return 0x2738 // ✸ '\u2738'
	case "a52":
		return 0x2739 // ✹ '\u2739'
	case "a53":
		return 0x273a // ✺ '\u273a'
	case "a54":
		return 0x273b // ✻ '\u273b'
	case "a55":
		return 0x273c // ✼ '\u273c'
	case "a56":
		return 0x273d // ✽ '\u273d'
	case "a57":
		return 0x273e // ✾ '\u273e'
	case "a58":
		return 0x273f // ✿ '\u273f'
	case "a59":
		return 0x2740 // ❀ '\u2740'
	case "a60":
		return 0x2741 // ❁ '\u2741'
	case "a61":
		return 0x2742 // ❂ '\u2742'
	case "a62":
		return 0x2743 // ❃ '\u2743'
	case "a63":
		return 0x2744 // ❄ '\u2744'
	case "a64":
		return 0x2745 // ❅ '\u2745'
	case "a65":
		return 0x2746 // ❆ '\u2746'
	case "a66":
		return 0x2747 // ❇ '\u2747'
	case "a67":
		return 0x2748 // ❈ '\u2748'
	case "a68":
		return 0x2749 // ❉ '\u2749'
	case "a69":
		return 0x274a // ❊ '\u274a'
	case "a70":
		return 0x274b // ❋ '\u274b'
	case "a71":
		return 0x25cf // ● '\u25cf'
	case "a72":
		return 0x274d // ❍ '\u274d'
	case "a73":
		return 0x25a0 // ■ '\u25a0'
	case "a74":
		return 0x274f // ❏ '\u274f'
	case "a75":
		return 0x2751 // ❑ '\u2751'
	case "a76":
		return 0x25b2 // ▲ '\u25b2'
	case "a77":
		return 0x25bc // ▼ '\u25bc'
	case "a78":
		return 0x25c6 // ◆ '\u25c6'
	case "a79":
		return 0x2756 // ❖ '\u2756'
	case "a81":
		return 0x25d7 // ◗ '\u25d7'
	case "a82":
		return 0x2758 // ❘ '\u2758'
	case "a83":
		return 0x2759 // ❙ '\u2759'
	case "a84":
		return 0x275a // ❚ '\u275a'
	case "a85":
		return 0xf8de //  '\uf8de'
	case "a86":
		return 0xf8e0 //  '\uf8e0'
	case "a87":
		return 0xf8e1 //  '\uf8e1'
	case "a88":
		return 0xf8e2 //  '\uf8e2'
	case "a89":
		return 0xf8d7 //  '\uf8d7'
	case "a90":
		return 0xf8d8 //  '\uf8d8'
	case "a91":
		return 0xf8db //  '\uf8db'
	case "a92":
		return 0xf8dc //  '\uf8dc'
	case "a93":
		return 0xf8d9 //  '\uf8d9'
	case "a94":
		return 0xf8da //  '\uf8da'
	case "a95":
		return 0xf8e3 //  '\uf8e3'
	case "a96":
		return 0xf8e4 //  '\uf8e4'
	case "a97":
		return 0x275b // ❛ '\u275b'
	case "a98":
		return 0x275c // ❜ '\u275c'
	case "a99":
		return 0x275d // ❝ '\u275d'
	case "a100":
		return 0x275e // ❞ '\u275e'
	case "a101":
		return 0x2761 // ❡ '\u2761'
	case "a102":
		return 0x2762 // ❢ '\u2762'
	case "a103":
		return 0x2763 // ❣ '\u2763'
	case "a104":
		return 0x2764 // ❤ '\u2764'
	case "a105":
		return 0x2710 // ✐ '\u2710'
	case "a106":
		return 0x2765 // ❥ '\u2765'
	case "a107":
		return 0x2766 // ❦ '\u2766'
	case "a108":
		return 0x2767 // ❧ '\u2767'
	case "a117":
		return 0x2709 // ✉ '\u2709'
	case "a118":
		return 0x2708 // ✈ '\u2708'
	case "a119":
		return 0x2707 // ✇ '\u2707'
	case "a120":
		return 0x2460 // ① '\u2460'
	case "a121":
		return 0x2461 // ② '\u2461'
	case "a122":
		return 0x2462 // ③ '\u2462'
	case "a123":
		return 0x2463 // ④ '\u2463'
	case "a124":
		return 0x2464 // ⑤ '\u2464'
	case "a125":
		return 0x2465 // ⑥ '\u2465'
	case "a126":
		return 0x2466 // ⑦ '\u2466'
	case "a127":
		return 0x2467 // ⑧ '\u2467'
	case "a128":
		return 0x2468 // ⑨ '\u2468'
	case "a129":
		return 0x2469 // ⑩ '\u2469'
	case "a130":
		return 0x2776 // ❶ '\u2776'
	case "a131":
		return 0x2777 // ❷ '\u2777'
	case "a132":
		return 0x2778 // ❸ '\u2778'
	case "a133":
		return 0x2779 // ❹ '\u2779'
	case "a134":
		return 0x277a // ❺ '\u277a'
	case "a135":
		return 0x277b // ❻ '\u277b'
	case "a136":
		return 0x277c // ❼ '\u277c'
	case "a137":
		return 0x277d // ❽ '\u277d'
	case "a138":
		return 0x277e // ❾ '\u277e'
	case "a139":
		return 0x277f // ❿ '\u277f'
	case "a140":
		return 0x2780 // ➀ '\u2780'
	case "a141":
		return 0x2781 // ➁ '\u2781'
	case "a142":
		return 0x2782 // ➂ '\u2782'
	case "a143":
		return 0x2783 // ➃ '\u2783'
	case "a144":
		return 0x2784 // ➄ '\u2784'
	case "a145":
		return 0x2785 // ➅ '\u2785'
	case "a146":
		return 0x2786 // ➆ '\u2786'
	case "a147":
		return 0x2787 // ➇ '\u2787'
	case "a148":
		return 0x2788 // ➈ '\u2788'
	case "a149":
		return 0x2789 // ➉ '\u2789'
	case "a150":
		return 0x278a // ➊ '\u278a'
	case "a151":
		return 0x278b // ➋ '\u278b'
	case "a152":
		return 0x278c // ➌ '\u278c'
	case "a153":
		return 0x278d // ➍ '\u278d'
	case "a154":
		return 0x278e // ➎ '\u278e'
	case "a155":
		return 0x278f // ➏ '\u278f'
	case "a156":
		return 0x2790 // ➐ '\u2790'
	case "a157":
		return 0x2791 // ➑ '\u2791'
	case "a158":
		return 0x2792 // ➒ '\u2792'
	case "a159":
		return 0x2793 // ➓ '\u2793'
	case "a160":
		return 0x2794 // ➔ '\u2794'
	case "a162":
		return 0x27a3 // ➣ '\u27a3'
	case "a164":
		return 0x2195 // ↕ '\u2195'
	case "a165":
		return 0x2799 // ➙ '\u2799'
	case "a166":
		return 0x279b // ➛ '\u279b'
	case "a167":
		return 0x279c // ➜ '\u279c'
	case "a168":
		return 0x279d // ➝ '\u279d'
	case "a169":
		return 0x279e // ➞ '\u279e'
	case "a170":
		return 0x279f // ➟ '\u279f'
	case "a171":
		return 0x27a0 // ➠ '\u27a0'
	case "a172":
		return 0x27a1 // ➡ '\u27a1'
	case "a173":
		return 0x27a2 // ➢ '\u27a2'
	case "a174":
		return 0x27a4 // ➤ '\u27a4'
	case "a175":
		return 0x27a5 // ➥ '\u27a5'
	case "a176":
		return 0x27a6 // ➦ '\u27a6'
	case "a177":
		return 0x27a7 // ➧ '\u27a7'
	case "a178":
		return 0x27a8 // ➨ '\u27a8'
	case "a179":
		return 0x27a9 // ➩ '\u27a9'
	case "a180":
		return 0x27ab // ➫ '\u27ab'
	case "a181":
		return 0x27ad // ➭ '\u27ad'
	case "a182":
		return 0x27af // ➯ '\u27af'
	case "a183":
		return 0x27b2 // ➲ '\u27b2'
	case "a184":
		return 0x27b3 // ➳ '\u27b3'
	case "a185":
		return 0x27b5 // ➵ '\u27b5'
	case "a186":
		return 0x27b8 // ➸ '\u27b8'
	case "a187":
		return 0x27ba // ➺ '\u27ba'
	case "a188":
		return 0x27bb // ➻ '\u27bb'
	case "a189":
		return 0x27bc // ➼ '\u27bc'
	case "a190":
		return 0x27bd // ➽ '\u27bd'
	case "a191":
		return 0x27be // ➾ '\u27be'
	case "a192":
		return 0x279a // ➚ '\u279a'
	case "a193":
		return 0x27aa // ➪ '\u27aa'
	case "a194":
		return 0x27b6 // ➶ '\u27b6'
	case "a195":
		return 0x27b9 // ➹ '\u27b9'
	case "a196":
		return 0x2798 // ➘ '\u2798'
	case "a197":
		return 0x27b4 // ➴ '\u27b4'
	case "a198":
		return 0x27b7 // ➷ '\u27b7'
	case "a199":
		return 0x27ac // ➬ '\u27ac'
	case "a200":
		return 0x27ae // ➮ '\u27ae'
	case "a201":
		return 0x27b1 // ➱ '\u27b1'
	case "a202":
		return 0x2703 // ✃ '\u2703'
	case "a203":
		return 0x2750 // ❐ '\u2750'
	case "a204":
		return 0x2752 // ❒ '\u2752'
	case "a205":
		return 0xf8dd //  '\uf8dd'
	case "a206":
		return 0xf8df //  '\uf8df'
	case "aabengali":
		return 0x0986 // আ '\u0986'
	case "aacute":
		return 0x00e1 // á '\u00e1'
	case "aadeva":
		return 0x0906 // आ '\u0906'
	case "aagujarati":
		return 0x0a86 // આ '\u0a86'
	case "aagurmukhi":
		return 0x0a06 // ਆ '\u0a06'
	case "aamatragurmukhi":
		return 0x0a3e // ਾ '\u0a3e'
	case "aarusquare":
		return 0x3303 // ㌃ '\u3303'
	case "aavowelsignbengali":
		return 0x09be // া '\u09be'
	case "aavowelsigndeva":
		return 0x093e // ा '\u093e'
	case "aavowelsigngujarati":
		return 0x0abe // ા '\u0abe'
	case "abbreviationmarkarmenian":
		return 0x055f // ՟ '\u055f'
	case "abbreviationsigndeva":
		return 0x0970 // ॰ '\u0970'
	case "abengali":
		return 0x0985 // অ '\u0985'
	case "abopomofo":
		return 0x311a // ㄚ '\u311a'
	case "abreve":
		return 0x0103 // ă '\u0103'
	case "abreveacute":
		return 0x1eaf // ắ '\u1eaf'
	case "abrevecyrillic":
		return 0x04d1 // ӑ '\u04d1'
	case "abrevedotbelow":
		return 0x1eb7 // ặ '\u1eb7'
	case "abrevegrave":
		return 0x1eb1 // ằ '\u1eb1'
	case "abrevehookabove":
		return 0x1eb3 // ẳ '\u1eb3'
	case "abrevetilde":
		return 0x1eb5 // ẵ '\u1eb5'
	case "acaron":
		return 0x01ce // ǎ '\u01ce'
	case "accountof":
		return 0x2100 // ℀ '\u2100'
	case "accurrent":
		return 0x23e6 // ⏦ '\u23e6'
	case "acidfree":
		return 0x267e // ♾ '\u267e'
	case "acircle":
		return 0x24d0 // ⓐ '\u24d0'
	case "acircumflex":
		return 0x00e2 // â '\u00e2'
	case "acircumflexacute":
		return 0x1ea5 // ấ '\u1ea5'
	case "acircumflexdotbelow":
		return 0x1ead // ậ '\u1ead'
	case "acircumflexgrave":
		return 0x1ea7 // ầ '\u1ea7'
	case "acircumflexhookabove":
		return 0x1ea9 // ẩ '\u1ea9'
	case "acircumflextilde":
		return 0x1eab // ẫ '\u1eab'
	case "acute":
		return 0x00b4 // ´ '\u00b4'
	case "acutebelowcmb":
		return 0x0317 // ̗ '\u0317'
	case "acutecomb":
		return 0x0301 // ́ '\u0301'
	case "acutedeva":
		return 0x0954 // ॔ '\u0954'
	case "acutelowmod":
		return 0x02cf // ˏ '\u02cf'
	case "acutenosp":
		return 0x0274 // ɴ '\u0274'
	case "acutetonecmb":
		return 0x0341 // ́ '\u0341'
	case "acwcirclearrow":
		return 0x2940 // ⥀ '\u2940'
	case "acwleftarcarrow":
		return 0x2939 // ⤹ '\u2939'
	case "acwopencirclearrow":
		return 0x21ba // ↺ '\u21ba'
	case "acwoverarcarrow":
		return 0x293a // ⤺ '\u293a'
	case "acwunderarcarrow":
		return 0x293b // ⤻ '\u293b'
	case "adblgrave":
		return 0x0201 // ȁ '\u0201'
	case "addakgurmukhi":
		return 0x0a71 // ੱ '\u0a71'
	case "addresssubject":
		return 0x2101 // ℁ '\u2101'
	case "adeva":
		return 0x0905 // अ '\u0905'
	case "adieresis":
		return 0x00e4 // ä '\u00e4'
	case "adieresiscyrillic":
		return 0x04d3 // ӓ '\u04d3'
	case "adieresismacron":
		return 0x01df // ǟ '\u01df'
	case "adotbelow":
		return 0x1ea1 // ạ '\u1ea1'
	case "adotmacron":
		return 0x01e1 // ǡ '\u01e1'
	case "adots":
		return 0x22f0 // ⋰ '\u22f0'
	case "ae":
		return 0x00e6 // æ '\u00e6'
	case "aeacute":
		return 0x01fd // ǽ '\u01fd'
	case "aekorean":
		return 0x3150 // ㅐ '\u3150'
	case "aemacron":
		return 0x01e3 // ǣ '\u01e3'
	case "afii299":
		return 0x200e //  '\u200e'
	case "afii300":
		return 0x200f //  '\u200f'
	case "afii301":
		return 0x200d //  '\u200d'
	case "afii10017":
		return 0x0410 // А '\u0410'
	case "afii10018":
		return 0x0411 // Б '\u0411'
	case "afii10019":
		return 0x0412 // В '\u0412'
	case "afii10024":
		return 0x0416 // Ж '\u0416'
	case "afii10025":
		return 0x0417 // З '\u0417'
	case "afii10027":
		return 0x0419 // Й '\u0419'
	case "afii10028":
		return 0x041a // К '\u041a'
	case "afii10031":
		return 0x041d // Н '\u041d'
	case "afii10033":
		return 0x041f // П '\u041f'
	case "afii10034":
		return 0x0420 // Р '\u0420'
	case "afii10035":
		return 0x0421 // С '\u0421'
	case "afii10036":
		return 0x0422 // Т '\u0422'
	case "afii10038":
		return 0x0424 // Ф '\u0424'
	case "afii10043":
		return 0x0429 // Щ '\u0429'
	case "afii10044":
		return 0x042a // Ъ '\u042a'
	case "afii10045":
		return 0x042b // Ы '\u042b'
	case "afii10046":
		return 0x042c // Ь '\u042c'
	case "afii10048":
		return 0x042e // Ю '\u042e'
	case "afii10049":
		return 0x042f // Я '\u042f'
	case "afii10051":
		return 0x0402 // Ђ '\u0402'
	case "afii10052":
		return 0x0403 // Ѓ '\u0403'
	case "afii10054":
		return 0x0405 // Ѕ '\u0405'
	case "afii10055":
		return 0x0406 // І '\u0406'
	case "afii10057":
		return 0x0408 // Ј '\u0408'
	case "afii10059":
		return 0x040a // Њ '\u040a'
	case "afii10063":
		return 0xf6c4 //  '\uf6c4'
	case "afii10064":
		return 0xf6c5 //  '\uf6c5'
	case "afii10065":
		return 0x0430 // а '\u0430'
	case "afii10067":
		return 0x0432 // в '\u0432'
	case "afii10068":
		return 0x0433 // г '\u0433'
	case "afii10069":
		return 0x0434 // д '\u0434'
	case "afii10071":
		return 0x0451 // ё '\u0451'
	case "afii10073":
		return 0x0437 // з '\u0437'
	case "afii10075":
		return 0x0439 // й '\u0439'
	case "afii10079":
		return 0x043d // н '\u043d'
	case "afii10084":
		return 0x0442 // т '\u0442'
	case "afii10085":
		return 0x0443 // у '\u0443'
	case "afii10086":
		return 0x0444 // ф '\u0444'
	case "afii10087":
		return 0x0445 // х '\u0445'
	case "afii10090":
		return 0x0448 // ш '\u0448'
	case "afii10091":
		return 0x0449 // щ '\u0449'
	case "afii10096":
		return 0x044e // ю '\u044e'
	case "afii10097":
		return 0x044f // я '\u044f'
	case "afii10102":
		return 0x0455 // ѕ '\u0455'
	case "afii10103":
		return 0x0456 // і '\u0456'
	case "afii10105":
		return 0x0458 // ј '\u0458'
	case "afii10107":
		return 0x045a // њ '\u045a'
	case "afii10109":
		return 0x045c // ќ '\u045c'
	case "afii10110":
		return 0x045e // ў '\u045e'
	case "afii10146":
		return 0x0462 // Ѣ '\u0462'
	case "afii10147":
		return 0x0472 // Ѳ '\u0472'
	case "afii10148":
		return 0x0474 // Ѵ '\u0474'
	case "afii10192":
		return 0xf6c6 //  '\uf6c6'
	case "afii10195":
		return 0x0473 // ѳ '\u0473'
	case "afii10196":
		return 0x0475 // ѵ '\u0475'
	case "afii10831":
		return 0xf6c7 //  '\uf6c7'
	case "afii10832":
		return 0xf6c8 //  '\uf6c8'
	case "afii57388":
		return 0x060c // ، '\u060c'
	case "afii57395":
		return 0x0663 // ٣ '\u0663'
	case "afii57398":
		return 0x0666 // ٦ '\u0666'
	case "afii57399":
		return 0x0667 // ٧ '\u0667'
	case "afii57403":
		return 0x061b // ؛ '\u061b'
	case "afii57407":
		return 0x061f // ؟ '\u061f'
	case "afii57410":
		return 0x0622 // آ '\u0622'
	case "afii57411":
		return 0x0623 // أ '\u0623'
	case "afii57412":
		return 0x0624 // ؤ '\u0624'
	case "afii57418":
		return 0x062a // ت '\u062a'
	case "afii57421":
		return 0x062d // ح '\u062d'
	case "afii57422":
		return 0x062e // خ '\u062e'
	case "afii57423":
		return 0x062f // د '\u062f'
	case "afii57427":
		return 0x0633 // س '\u0633'
	case "afii57428":
		return 0x0634 // ش '\u0634'
	case "afii57429":
		return 0x0635 // ص '\u0635'
	case "afii57430":
		return 0x0636 // ض '\u0636'
	case "afii57433":
		return 0x0639 // ع '\u0639'
	case "afii57441":
		return 0x0641 // ف '\u0641'
	case "afii57442":
		return 0x0642 // ق '\u0642'
	case "afii57443":
		return 0x0643 // ك '\u0643'
	case "afii57444":
		return 0x0644 // ل '\u0644'
	case "afii57445":
		return 0x0645 // م '\u0645'
	case "afii57451":
		return 0x064b // ً '\u064b'
	case "afii57452":
		return 0x064c // ٌ '\u064c'
	case "afii57456":
		return 0x0650 // ِ '\u0650'
	case "afii57508":
		return 0x0698 // ژ '\u0698'
	case "afii57511":
		return 0x0679 // ٹ '\u0679'
	case "afii57512":
		return 0x0688 // ڈ '\u0688'
	case "afii57514":
		return 0x06ba // ں '\u06ba'
	case "afii57534":
		return 0x06d5 // ە '\u06d5'
	case "afii57636":
		return 0x20aa // ₪ '\u20aa'
	case "afii57645":
		return 0x05be // ־ '\u05be'
	case "afii57666":
		return 0x05d2 // ג '\u05d2'
	case "afii57668":
		return 0x05d4 // ה '\u05d4'
	case "afii57670":
		return 0x05d6 // ז '\u05d6'
	case "afii57671":
		return 0x05d7 // ח '\u05d7'
	case "afii57673":
		return 0x05d9 // י '\u05d9'
	case "afii57674":
		return 0x05da // ך '\u05da'
	case "afii57679":
		return 0x05df // ן '\u05df'
	case "afii57684":
		return 0x05e4 // פ '\u05e4'
	case "afii57686":
		return 0x05e6 // צ '\u05e6'
	case "afii57695":
		return 0xfb2b // שׂ '\ufb2b'
	case "afii57716":
		return 0x05f0 // װ '\u05f0'
	case "afii57717":
		return 0x05f1 // ױ '\u05f1'
	case "afii57797":
		return 0x05b8 // ָ '\u05b8'
	case "afii57799":
		return 0x05b0 // ְ '\u05b0'
	case "afii57803":
		return 0x05c2 // ׂ '\u05c2'
	case "afii57841":
		return 0x05bf // ֿ '\u05bf'
	case "afii57842":
		return 0x05c0 // ׀ '\u05c0'
	case "afii61289":
		return 0x2113 // ℓ '\u2113'
	case "afii61573":
		return 0x202c //  '\u202c'
	case "afii61574":
		return 0x202d //  '\u202d'
	case "afii61575":
		return 0x202e //  '\u202e'
	case "afii61664":
		return 0x200c //  '\u200c'
	case "afii63167":
		return 0x066d // ٭ '\u066d'
	case "afii64937":
		return 0x02bd // ʽ '\u02bd'
	case "agrave":
		return 0x00e0 // à '\u00e0'
	case "agujarati":
		return 0x0a85 // અ '\u0a85'
	case "agurmukhi":
		return 0x0a05 // ਅ '\u0a05'
	case "ahiragana":
		return 0x3042 // あ '\u3042'
	case "ahookabove":
		return 0x1ea3 // ả '\u1ea3'
	case "aibengali":
		return 0x0990 // ঐ '\u0990'
	case "aibopomofo":
		return 0x311e // ㄞ '\u311e'
	case "aideva":
		return 0x0910 // ऐ '\u0910'
	case "aiecyrillic":
		return 0x04d5 // ӕ '\u04d5'
	case "aigujarati":
		return 0x0a90 // ઐ '\u0a90'
	case "aigurmukhi":
		return 0x0a10 // ਐ '\u0a10'
	case "aimatragurmukhi":
		return 0x0a48 // ੈ '\u0a48'
	case "ainfinalarabic":
		return 0xfeca // ﻊ '\ufeca'
	case "aininitialarabic":
		return 0xfecb // ﻋ '\ufecb'
	case "ainisolated":
		return 0xfec9 // ﻉ '\ufec9'
	case "ainmedialarabic":
		return 0xfecc // ﻌ '\ufecc'
	case "ainvertedbreve":
		return 0x0203 // ȃ '\u0203'
	case "aivowelsignbengali":
		return 0x09c8 // ৈ '\u09c8'
	case "aivowelsigndeva":
		return 0x0948 // ै '\u0948'
	case "aivowelsigngujarati":
		return 0x0ac8 // ૈ '\u0ac8'
	case "akatakana":
		return 0x30a2 // ア '\u30a2'
	case "akatakanahalfwidth":
		return 0xff71 // ｱ '\uff71'
	case "akorean":
		return 0x314f // ㅏ '\u314f'
	case "alef":
		return 0x05d0 // א '\u05d0'
	case "alefarabic":
		return 0x0627 // ا '\u0627'
	case "alefdageshhebrew":
		return 0xfb30 // אּ '\ufb30'
	case "aleffinalarabic":
		return 0xfe8e // ﺎ '\ufe8e'
	case "alefhamzaabovefinalarabic":
		return 0xfe84 // ﺄ '\ufe84'
	case "alefhamzabelowarabic":
		return 0x0625 // إ '\u0625'
	case "alefhamzabelowfinalarabic":
		return 0xfe88 // ﺈ '\ufe88'
	case "alefisolated":
		return 0xfe8d // ﺍ '\ufe8d'
	case "aleflamedhebrew":
		return 0xfb4f // ﭏ '\ufb4f'
	case "alefmaddaabovefinalarabic":
		return 0xfe82 // ﺂ '\ufe82'
	case "alefmaksuraarabic":
		return 0x0649 // ى '\u0649'
	case "alefmaksurafinalarabic":
		return 0xfef0 // ﻰ '\ufef0'
	case "alefmaksuraisolated":
		return 0xfeef // ﻯ '\ufeef'
	case "alefmaksuramedialarabic":
		return 0xfef4 // ﻴ '\ufef4'
	case "alefpatahhebrew":
		return 0xfb2e // אַ '\ufb2e'
	case "alefqamatshebrew":
		return 0xfb2f // אָ '\ufb2f'
	case "alefwasla":
		return 0x0671 // ٱ '\u0671'
	case "alefwaslafinal":
		return 0xfb51 // ﭑ '\ufb51'
	case "alefwaslaisolated":
		return 0xfb50 // ﭐ '\ufb50'
	case "alefwithfathatanfinal":
		return 0xfd3c // ﴼ '\ufd3c'
	case "alefwithfathatanisolated":
		return 0xfd3d // ﴽ '\ufd3d'
	case "alefwithhamzaaboveisolated":
		return 0xfe83 // ﺃ '\ufe83'
	case "alefwithhamzabelowisolated":
		return 0xfe87 // ﺇ '\ufe87'
	case "alefwithmaddaaboveisolated":
		return 0xfe81 // ﺁ '\ufe81'
	case "aleph":
		return 0x2135 // ℵ '\u2135'
	case "allequal":
		return 0x224c // ≌ '\u224c'
	case "alpha":
		return 0x03b1 // α '\u03b1'
	case "alphatonos":
		return 0x03ac // ά '\u03ac'
	case "altselector":
		return 0xd802 //  '\ufffd'
	case "amacron":
		return 0x0101 // ā '\u0101'
	case "amonospace":
		return 0xff41 // ａ '\uff41'
	case "ampersand":
		return 0x0026 // & '&'
	case "ampersandmonospace":
		return 0xff06 // ＆ '\uff06'
	case "ampersandsmall":
		return 0xf726 //  '\uf726'
	case "amsquare":
		return 0x33c2 // ㏂ '\u33c2'
	case "anbopomofo":
		return 0x3122 // ㄢ '\u3122'
	case "angbopomofo":
		return 0x3124 // ㄤ '\u3124'
	case "angbracketleft":
		return 0x27e8 // ⟨ '\u27e8'
	case "angbracketright":
		return 0x27e9 // ⟩ '\u27e9'
	case "angdnr":
		return 0x299f // ⦟ '\u299f'
	case "angkhankhuthai":
		return 0x0e5a // ๚ '\u0e5a'
	case "angle":
		return 0x2220 // ∠ '\u2220'
	case "anglebracketleft":
		return 0x3008 // 〈 '\u3008'
	case "anglebracketleftvertical":
		return 0xfe3f // ︿ '\ufe3f'
	case "anglebracketright":
		return 0x3009 // 〉 '\u3009'
	case "anglebracketrightvertical":
		return 0xfe40 // ﹀ '\ufe40'
	case "angleleft":
		return 0x2329 // 〈 '\u2329'
	case "angleright":
		return 0x232a // 〉 '\u232a'
	case "angles":
		return 0x299e // ⦞ '\u299e'
	case "angleubar":
		return 0x29a4 // ⦤ '\u29a4'
	case "angstrom":
		return 0x212b // Å '\u212b'
	case "annuity":
		return 0x20e7 // ⃧ '\u20e7'
	case "anoteleia":
		return 0x0387 // · '\u0387'
	case "anticlockwise":
		return 0x27f2 // ⟲ '\u27f2'
	case "anudattadeva":
		return 0x0952 // ॒ '\u0952'
	case "anusvarabengali":
		return 0x0982 // ং '\u0982'
	case "anusvaradeva":
		return 0x0902 // ं '\u0902'
	case "anusvaragujarati":
		return 0x0a82 // ં '\u0a82'
	case "aogonek":
		return 0x0105 // ą '\u0105'
	case "apaatosquare":
		return 0x3300 // ㌀ '\u3300'
	case "aparen":
		return 0x249c // ⒜ '\u249c'
	case "apostrophe":
		return 0x0245 // Ʌ '\u0245'
	case "apostrophearmenian":
		return 0x055a // ՚ '\u055a'
	case "apostrophemod":
		return 0x02bc // ʼ '\u02bc'
	case "apostropherev":
		return 0x0246 // Ɇ '\u0246'
	case "apple":
		return 0xf8ff //  '\uf8ff'
	case "approaches":
		return 0x2250 // ≐ '\u2250'
	case "approxeqq":
		return 0x2a70 // ⩰ '\u2a70'
	case "approxequal":
		return 0x2248 // ≈ '\u2248'
	case "approxequalorimage":
		return 0x2252 // ≒ '\u2252'
	case "approxident":
		return 0x224b // ≋ '\u224b'
	case "approxorequal":
		return 0x224a // ≊ '\u224a'
	case "araeaekorean":
		return 0x318e // ㆎ '\u318e'
	case "araeakorean":
		return 0x318d // ㆍ '\u318d'
	case "arc":
		return 0x2312 // ⌒ '\u2312'
	case "arceq":
		return 0x2258 // ≘ '\u2258'
	case "archleftdown":
		return 0x21b6 // ↶ '\u21b6'
	case "archrightdown":
		return 0x21b7 // ↷ '\u21b7'
	case "arighthalfring":
		return 0x1e9a // ẚ '\u1e9a'
	case "aring":
		return 0x00e5 // å '\u00e5'
	case "aringacute":
		return 0x01fb // ǻ '\u01fb'
	case "aringbelow":
		return 0x1e01 // ḁ '\u1e01'
	case "arrowbardown":
		return 0x0590 //  '\u0590'
	case "arrowbarleft":
		return 0x058d // ֍ '\u058d'
	case "arrowbarright":
		return 0x058f // ֏ '\u058f'
	case "arrowbarup":
		return 0x058e // ֎ '\u058e'
	case "arrowboth":
		return 0x2194 // ↔ '\u2194'
	case "arrowdashdown":
		return 0x21e3 // ⇣ '\u21e3'
	case "arrowdashleft":
		return 0x21e0 // ⇠ '\u21e0'
	case "arrowdashright":
		return 0x21e2 // ⇢ '\u21e2'
	case "arrowdashup":
		return 0x21e1 // ⇡ '\u21e1'
	case "arrowdblboth":
		return 0x21d4 // ⇔ '\u21d4'
	case "arrowdblbothv":
		return 0x21d5 // ⇕ '\u21d5'
	case "arrowdbldown":
		return 0x21d3 // ⇓ '\u21d3'
	case "arrowdblleft":
		return 0x21d0 // ⇐ '\u21d0'
	case "arrowdblright":
		return 0x21d2 // ⇒ '\u21d2'
	case "arrowdblup":
		return 0x21d1 // ⇑ '\u21d1'
	case "arrowdown":
		return 0x2193 // ↓ '\u2193'
	case "arrowdownleft":
		return 0x2199 // ↙ '\u2199'
	case "arrowdownright":
		return 0x2198 // ↘ '\u2198'
	case "arrowdownwhite":
		return 0x21e9 // ⇩ '\u21e9'
	case "arrowheaddownmod":
		return 0x02c5 // ˅ '\u02c5'
	case "arrowheadleftmod":
		return 0x02c2 // ˂ '\u02c2'
	case "arrowheadrightmod":
		return 0x02c3 // ˃ '\u02c3'
	case "arrowheadupmod":
		return 0x02c4 // ˄ '\u02c4'
	case "arrowhookleft":
		return 0x21aa // ↪ '\u21aa'
	case "arrowhookright":
		return 0x21a9 // ↩ '\u21a9'
	case "arrowhorizex":
		return 0xf8e7 //  '\uf8e7'
	case "arrowleft":
		return 0x2190 // ← '\u2190'
	case "arrowleftbothalf":
		return 0x21bd // ↽ '\u21bd'
	case "arrowleftdblstroke":
		return 0x21cd // ⇍ '\u21cd'
	case "arrowleftoverright":
		return 0x21c6 // ⇆ '\u21c6'
	case "arrowleftwhite":
		return 0x21e6 // ⇦ '\u21e6'
	case "arrowright":
		return 0x2192 // → '\u2192'
	case "arrowrightbothalf":
		return 0x21c1 // ⇁ '\u21c1'
	case "arrowrightdblstroke":
		return 0x21cf // ⇏ '\u21cf'
	case "arrowrightoverleft":
		return 0x21c4 // ⇄ '\u21c4'
	case "arrowrightwhite":
		return 0x21e8 // ⇨ '\u21e8'
	case "arrowtableft":
		return 0x21e4 // ⇤ '\u21e4'
	case "arrowtabright":
		return 0x21e5 // ⇥ '\u21e5'
	case "arrowtailleft":
		return 0x21a2 // ↢ '\u21a2'
	case "arrowtailright":
		return 0x21a3 // ↣ '\u21a3'
	case "arrowtripleleft":
		return 0x21da // ⇚ '\u21da'
	case "arrowtripleright":
		return 0x21db // ⇛ '\u21db'
	case "arrowup":
		return 0x2191 // ↑ '\u2191'
	case "arrowupdownbase":
		return 0x21a8 // ↨ '\u21a8'
	case "arrowupleft":
		return 0x2196 // ↖ '\u2196'
	case "arrowupleftofdown":
		return 0x21c5 // ⇅ '\u21c5'
	case "arrowupright":
		return 0x2197 // ↗ '\u2197'
	case "arrowupwhite":
		return 0x21e7 // ⇧ '\u21e7'
	case "arrowvertex":
		return 0xf8e6 //  '\uf8e6'
	case "ascendercompwordmark":
		return 0xd80a //  '\ufffd'
	case "asciicircum":
		return 0x005e // ^ '^'
	case "asciicircummonospace":
		return 0xff3e // ＾ '\uff3e'
	case "asciitilde":
		return 0x007e // ~ '~'
	case "asciitildemonospace":
		return 0xff5e // ～ '\uff5e'
	case "ascript":
		return 0x0251 // ɑ '\u0251'
	case "ascriptturned":
		return 0x0252 // ɒ '\u0252'
	case "asmallhiragana":
		return 0x3041 // ぁ '\u3041'
	case "asmallkatakana":
		return 0x30a1 // ァ '\u30a1'
	case "asmallkatakanahalfwidth":
		return 0xff67 // ｧ '\uff67'
	case "assert":
		return 0x22a6 // ⊦ '\u22a6'
	case "asteq":
		return 0x2a6e // ⩮ '\u2a6e'
	case "asteraccent":
		return 0x20f0 // ⃰ '\u20f0'
	case "asterisk":
		return 0x002a // * '*'
	case "asteriskmath":
		return 0x2217 // ∗ '\u2217'
	case "asteriskmonospace":
		return 0xff0a // ＊ '\uff0a'
	case "asterisksmall":
		return 0xfe61 // ﹡ '\ufe61'
	case "asterism":
		return 0x2042 // ⁂ '\u2042'
	case "astrosun":
		return 0x2609 // ☉ '\u2609'
	case "asuperior":
		return 0xf6e9 //  '\uf6e9'
	case "asymptoticallyequal":
		return 0x2243 // ≃ '\u2243'
	case "at":
		return 0x0040 // @ '@'
	case "atilde":
		return 0x00e3 // ã '\u00e3'
	case "atmonospace":
		return 0xff20 // ＠ '\uff20'
	case "atsmall":
		return 0xfe6b // ﹫ '\ufe6b'
	case "aturned":
		return 0x0250 // ɐ '\u0250'
	case "aubengali":
		return 0x0994 // ঔ '\u0994'
	case "aubopomofo":
		return 0x3120 // ㄠ '\u3120'
	case "audeva":
		return 0x0914 // औ '\u0914'
	case "augujarati":
		return 0x0a94 // ઔ '\u0a94'
	case "augurmukhi":
		return 0x0a14 // ਔ '\u0a14'
	case "aulengthmarkbengali":
		return 0x09d7 // ৗ '\u09d7'
	case "aumatragurmukhi":
		return 0x0a4c // ੌ '\u0a4c'
	case "auvowelsignbengali":
		return 0x09cc // ৌ '\u09cc'
	case "auvowelsigndeva":
		return 0x094c // ौ '\u094c'
	case "auvowelsigngujarati":
		return 0x0acc // ૌ '\u0acc'
	case "avagrahadeva":
		return 0x093d // ऽ '\u093d'
	case "awint":
		return 0x2a11 // ⨑ '\u2a11'
	case "aybarmenian":
		return 0x0561 // ա '\u0561'
	case "ayinaltonehebrew":
		return 0xfb20 // ﬠ '\ufb20'
	case "ayinhebrew":
		return 0x05e2 // ע '\u05e2'
	case "b":
		return 0x0062 // b 'b'
	case "bNot":
		return 0x2aed // ⫭ '\u2aed'
	case "babengali":
		return 0x09ac // ব '\u09ac'
	case "backdprime":
		return 0x2036 // ‶ '\u2036'
	case "backed":
		return 0x024c // Ɍ '\u024c'
	case "backslash":
		return 0x005c // \\ '\\'
	case "backslashmonospace":
		return 0xff3c // ＼ '\uff3c'
	case "backtrprime":
		return 0x2037 // ‷ '\u2037'
	case "badeva":
		return 0x092c // ब '\u092c'
	case "bagmember":
		return 0x22ff // ⋿ '\u22ff'
	case "bagujarati":
		return 0x0aac // બ '\u0aac'
	case "bagurmukhi":
		return 0x0a2c // ਬ '\u0a2c'
	case "bahiragana":
		return 0x3070 // ば '\u3070'
	case "bahtthai":
		return 0x0e3f // ฿ '\u0e3f'
	case "bakatakana":
		return 0x30d0 // バ '\u30d0'
	case "bar":
		return 0x007c // | '|'
	case "barV":
		return 0x2aea // ⫪ '\u2aea'
	case "barcap":
		return 0x2a43 // ⩃ '\u2a43'
	case "barcup":
		return 0x2a42 // ⩂ '\u2a42'
	case "bardownharpoonleft":
		return 0x2961 // ⥡ '\u2961'
	case "bardownharpoonright":
		return 0x295d // ⥝ '\u295d'
	case "barleftarrowrightarrowba":
		return 0x21b9 // ↹ '\u21b9'
	case "barleftharpoondown":
		return 0x2956 // ⥖ '\u2956'
	case "barleftharpoonup":
		return 0x2952 // ⥒ '\u2952'
	case "barmidlongnosp":
		return 0x02a9 // ʩ '\u02a9'
	case "barmonospace":
		return 0xff5c // ｜ '\uff5c'
	case "barovernorthwestarrow":
		return 0x21b8 // ↸ '\u21b8'
	case "barrightarrowdiamond":
		return 0x2920 // ⤠ '\u2920'
	case "barrightharpoondown":
		return 0x295f // ⥟ '\u295f'
	case "barrightharpoonup":
		return 0x295b // ⥛ '\u295b'
	case "baruparrow":
		return 0x2912 // ⤒ '\u2912'
	case "barupharpoonleft":
		return 0x2958 // ⥘ '\u2958'
	case "barupharpoonright":
		return 0x2954 // ⥔ '\u2954'
	case "barvee":
		return 0x22bd // ⊽ '\u22bd'
	case "bbopomofo":
		return 0x3105 // ㄅ '\u3105'
	case "bbrktbrk":
		return 0x23b6 // ⎶ '\u23b6'
	case "bcircle":
		return 0x24d1 // ⓑ '\u24d1'
	case "bdotaccent":
		return 0x1e03 // ḃ '\u1e03'
	case "bdotbelow":
		return 0x1e05 // ḅ '\u1e05'
	case "bdtriplevdash":
		return 0x2506 // ┆ '\u2506'
	case "beamedsixteenthnotes":
		return 0x266c // ♬ '\u266c'
	case "because":
		return 0x2235 // ∵ '\u2235'
	case "becyrillic":
		return 0x0431 // б '\u0431'
	case "beharabic":
		return 0x0628 // ب '\u0628'
	case "behfinalarabic":
		return 0xfe90 // ﺐ '\ufe90'
	case "behinitialarabic":
		return 0xfe91 // ﺑ '\ufe91'
	case "behiragana":
		return 0x3079 // べ '\u3079'
	case "behisolated":
		return 0xfe8f // ﺏ '\ufe8f'
	case "behmedialarabic":
		return 0xfe92 // ﺒ '\ufe92'
	case "behmeeminitialarabic":
		return 0xfc9f // ﲟ '\ufc9f'
	case "behmeemisolatedarabic":
		return 0xfc08 // ﰈ '\ufc08'
	case "behnoonfinalarabic":
		return 0xfc6d // ﱭ '\ufc6d'
	case "behwithalefmaksurafinal":
		return 0xfc6e // ﱮ '\ufc6e'
	case "behwithalefmaksuraisolated":
		return 0xfc09 // ﰉ '\ufc09'
	case "behwithhahinitial":
		return 0xfc9d // ﲝ '\ufc9d'
	case "behwithhehinitial":
		return 0xe812 //  '\ue812'
	case "behwithjeeminitial":
		return 0xfc9c // ﲜ '\ufc9c'
	case "behwithkhahinitial":
		return 0xfc9e // ﲞ '\ufc9e'
	case "behwithrehfinal":
		return 0xfc6a // ﱪ '\ufc6a'
	case "behwithyehfinal":
		return 0xfc6f // ﱯ '\ufc6f'
	case "behwithyehisolated":
		return 0xfc0a // ﰊ '\ufc0a'
	case "bekatakana":
		return 0x30d9 // ベ '\u30d9'
	case "benarmenian":
		return 0x0562 // բ '\u0562'
	case "benzenr":
		return 0x23e3 // ⏣ '\u23e3'
	case "beta":
		return 0x03b2 // β '\u03b2'
	case "betasymbolgreek":
		return 0x03d0 // ϐ '\u03d0'
	case "betdageshhebrew":
		return 0xfb31 // בּ '\ufb31'
	case "beth":
		return 0x2136 // ℶ '\u2136'
	case "bethebrew":
		return 0x05d1 // ב '\u05d1'
	case "betrafehebrew":
		return 0xfb4c // בֿ '\ufb4c'
	case "between":
		return 0x226c // ≬ '\u226c'
	case "bhabengali":
		return 0x09ad // ভ '\u09ad'
	case "bhadeva":
		return 0x092d // भ '\u092d'
	case "bhagujarati":
		return 0x0aad // ભ '\u0aad'
	case "bhagurmukhi":
		return 0x0a2d // ਭ '\u0a2d'
	case "bhook":
		return 0x0253 // ɓ '\u0253'
	case "bigbot":
		return 0x27d8 // ⟘ '\u27d8'
	case "bigcupdot":
		return 0x2a03 // ⨃ '\u2a03'
	case "biginterleave":
		return 0x2afc // ⫼ '\u2afc'
	case "bigodot":
		return 0x2a00 // ⨀ '\u2a00'
	case "bigoplus":
		return 0x2a01 // ⨁ '\u2a01'
	case "bigotimes":
		return 0x2a02 // ⨂ '\u2a02'
	case "bigslopedvee":
		return 0x2a57 // ⩗ '\u2a57'
	case "bigslopedwedge":
		return 0x2a58 // ⩘ '\u2a58'
	case "bigsqcap":
		return 0x2a05 // ⨅ '\u2a05'
	case "bigsqcup":
		return 0x2a06 // ⨆ '\u2a06'
	case "bigtalloblong":
		return 0x2aff // ⫿ '\u2aff'
	case "bigtimes":
		return 0x2a09 // ⨉ '\u2a09'
	case "bigtop":
		return 0x27d9 // ⟙ '\u27d9'
	case "bigtriangleleft":
		return 0x2a1e // ⨞ '\u2a1e'
	case "biguplus":
		return 0x2a04 // ⨄ '\u2a04'
	case "bigvee":
		return 0x22c1 // ⋁ '\u22c1'
	case "bigwedge":
		return 0x22c0 // ⋀ '\u22c0'
	case "bihiragana":
		return 0x3073 // び '\u3073'
	case "bikatakana":
		return 0x30d3 // ビ '\u30d3'
	case "bilabialclick":
		return 0x0298 // ʘ '\u0298'
	case "bindigurmukhi":
		return 0x0a02 // ਂ '\u0a02'
	case "birusquare":
		return 0x3331 // ㌱ '\u3331'
	case "blackcircledownarrow":
		return 0x29ed // ⧭ '\u29ed'
	case "blackcircledrightdot":
		return 0x2688 // ⚈ '\u2688'
	case "blackcircledtwodots":
		return 0x2689 // ⚉ '\u2689'
	case "blackcircleulquadwhite":
		return 0x25d5 // ◕ '\u25d5'
	case "blackdiamonddownarrow":
		return 0x29ea // ⧪ '\u29ea'
	case "blackhourglass":
		return 0x29d7 // ⧗ '\u29d7'
	case "blacklefthalfcircle":
		return 0x25d6 // ◖ '\u25d6'
	case "blackleftpointingpointer":
		return 0x25c4 // ◄ '\u25c4'
	case "blackleftpointingtriangle":
		return 0x25c0 // ◀ '\u25c0'
	case "blacklenticularbracketleft":
		return 0x3010 // 【 '\u3010'
	case "blacklenticularbracketleftvertical":
		return 0xfe3b // ︻ '\ufe3b'
	case "blacklenticularbracketright":
		return 0x3011 // 】 '\u3011'
	case "blacklenticularbracketrightvertical":
		return 0xfe3c // ︼ '\ufe3c'
	case "blacklowerlefttriangle":
		return 0x25e3 // ◣ '\u25e3'
	case "blacklowerrighttriangle":
		return 0x25e2 // ◢ '\u25e2'
	case "blackrectangle":
		return 0x25ac // ▬ '\u25ac'
	case "blackrightpointingpointer":
		return 0x25ba // ► '\u25ba'
	case "blackrightpointingtriangle":
		return 0x25b6 // ▶ '\u25b6'
	case "blacksmallsquare":
		return 0x25aa // ▪ '\u25aa'
	case "blacksmilingface":
		return 0x263b // ☻ '\u263b'
	case "blacktriangledown":
		return 0x25be // ▾ '\u25be'
	case "blackupperlefttriangle":
		return 0x25e4 // ◤ '\u25e4'
	case "blackupperrighttriangle":
		return 0x25e5 // ◥ '\u25e5'
	case "blackuppointingsmalltriangle":
		return 0x25b4 // ▴ '\u25b4'
	case "blank":
		return 0x2423 // ␣ '\u2423'
	case "blinebelow":
		return 0x1e07 // ḇ '\u1e07'
	case "blkhorzoval":
		return 0x2b2c // ⬬ '\u2b2c'
	case "blkvertoval":
		return 0x2b2e // ⬮ '\u2b2e'
	case "block":
		return 0x2588 // █ '\u2588'
	case "bmonospace":
		return 0xff42 // ｂ '\uff42'
	case "bobaimaithai":
		return 0x0e1a // บ '\u0e1a'
	case "bohiragana":
		return 0x307c // ぼ '\u307c'
	case "bokatakana":
		return 0x30dc // ボ '\u30dc'
	case "botsemicircle":
		return 0x25e1 // ◡ '\u25e1'
	case "bowtie":
		return 0x22c8 // ⋈ '\u22c8'
	case "boxast":
		return 0x29c6 // ⧆ '\u29c6'
	case "boxbar":
		return 0x25eb // ◫ '\u25eb'
	case "boxbox":
		return 0x29c8 // ⧈ '\u29c8'
	case "boxbslash":
		return 0x29c5 // ⧅ '\u29c5'
	case "boxcircle":
		return 0x29c7 // ⧇ '\u29c7'
	case "boxdiag":
		return 0x29c4 // ⧄ '\u29c4'
	case "boxonbox":
		return 0x29c9 // ⧉ '\u29c9'
	case "bparen":
		return 0x249d // ⒝ '\u249d'
	case "bqsquare":
		return 0x33c3 // ㏃ '\u33c3'
	case "braceex":
		return 0xf8f4 //  '\uf8f4'
	case "braceleft":
		return 0x007b // { '{'
	case "braceleftbt":
		return 0xf8f3 //  '\uf8f3'
	case "braceleftmid":
		return 0xf8f2 //  '\uf8f2'
	case "braceleftmonospace":
		return 0xff5b // ｛ '\uff5b'
	case "braceleftsmall":
		return 0xfe5b // ﹛ '\ufe5b'
	case "bracelefttp":
		return 0xf8f1 //  '\uf8f1'
	case "braceleftvertical":
		return 0xfe37 // ︷ '\ufe37'
	case "braceright":
		return 0x007d // } '}'
	case "bracerightbt":
		return 0xf8fe //  '\uf8fe'
	case "bracerightmid":
		return 0xf8fd //  '\uf8fd'
	case "bracerightmonospace":
		return 0xff5d // ｝ '\uff5d'
	case "bracerightsmall":
		return 0xfe5c // ﹜ '\ufe5c'
	case "bracerighttp":
		return 0xf8fc //  '\uf8fc'
	case "bracerightvertical":
		return 0xfe38 // ︸ '\ufe38'
	case "bracketleft":
		return 0x005b // [ '['
	case "bracketleftbt":
		return 0xf8f0 //  '\uf8f0'
	case "bracketleftex":
		return 0xf8ef //  '\uf8ef'
	case "bracketleftmonospace":
		return 0xff3b // ［ '\uff3b'
	case "bracketleftquill":
		return 0x2045 // ⁅ '\u2045'
	case "bracketlefttp":
		return 0xf8ee //  '\uf8ee'
	case "bracketright":
		return 0x005d // ] ']'
	case "bracketrightbt":
		return 0xf8fb //  '\uf8fb'
	case "bracketrightex":
		return 0xf8fa //  '\uf8fa'
	case "bracketrightmonospace":
		return 0xff3d // ］ '\uff3d'
	case "bracketrightquill":
		return 0x2046 // ⁆ '\u2046'
	case "bracketrighttp":
		return 0xf8f9 //  '\uf8f9'
	case "breve":
		return 0x02d8 // ˘ '\u02d8'
	case "breve1":
		return 0xf006 //  '\uf006'
	case "brevebelowcmb":
		return 0x032e // ̮ '\u032e'
	case "brevecmb":
		return 0x0306 // ̆ '\u0306'
	case "breveinvertedbelowcmb":
		return 0x032f // ̯ '\u032f'
	case "breveinvertedcmb":
		return 0x0311 // ̑ '\u0311'
	case "breveinverteddoublecmb":
		return 0x0361 // ͡ '\u0361'
	case "bridgebelowcmb":
		return 0x032a // ̪ '\u032a'
	case "bridgeinvertedbelowcmb":
		return 0x033a // ̺ '\u033a'
	case "bridgeinvsubnosp":
		return 0x02ad // ʭ '\u02ad'
	case "brokenbar":
		return 0x00a6 // ¦ '\u00a6'
	case "bsimilarleftarrow":
		return 0x2b41 // ⭁ '\u2b41'
	case "bsimilarrightarrow":
		return 0x2b47 // ⭇ '\u2b47'
	case "bsolhsub":
		return 0x27c8 // ⟈ '\u27c8'
	case "bstroke":
		return 0x0180 // ƀ '\u0180'
	case "bsuperior":
		return 0xf6ea //  '\uf6ea'
	case "btimes":
		return 0x2a32 // ⨲ '\u2a32'
	case "btopbar":
		return 0x0183 // ƃ '\u0183'
	case "buhiragana":
		return 0x3076 // ぶ '\u3076'
	case "bukatakana":
		return 0x30d6 // ブ '\u30d6'
	case "bullet":
		return 0x2022 // • '\u2022'
	case "bulletoperator":
		return 0x2219 // ∙ '\u2219'
	case "bullseye":
		return 0x25ce // ◎ '\u25ce'
	case "bumpeqq":
		return 0x2aae // ⪮ '\u2aae'
	case "c":
		return 0x0063 // c 'c'
	case "c128":
		return 0x0080 //  '\u0080'
	case "c129":
		return 0x0081 //  '\u0081'
	case "c141":
		return 0x008d //  '\u008d'
	case "c142":
		return 0x008e //  '\u008e'
	case "c143":
		return 0x008f //  '\u008f'
	case "caarmenian":
		return 0x056e // ծ '\u056e'
	case "cabengali":
		return 0x099a // চ '\u099a'
	case "cacute":
		return 0x0107 // ć '\u0107'
	case "cadauna":
		return 0x2106 // ℆ '\u2106'
	case "cadeva":
		return 0x091a // च '\u091a'
	case "cagujarati":
		return 0x0a9a // ચ '\u0a9a'
	case "cagurmukhi":
		return 0x0a1a // ਚ '\u0a1a'
	case "calsquare":
		return 0x3388 // ㎈ '\u3388'
	case "candrabindubengali":
		return 0x0981 // ঁ '\u0981'
	case "candrabinducmb":
		return 0x0310 // ̐ '\u0310'
	case "candrabindudeva":
		return 0x0901 // ँ '\u0901'
	case "candrabindugujarati":
		return 0x0a81 // ઁ '\u0a81'
	case "capbarcup":
		return 0x2a49 // ⩉ '\u2a49'
	case "capdot":
		return 0x2a40 // ⩀ '\u2a40'
	case "capitalcompwordmark":
		return 0xd809 //  '\ufffd'
	case "capovercup":
		return 0x2a47 // ⩇ '\u2a47'
	case "capslock":
		return 0x21ea // ⇪ '\u21ea'
	case "capwedge":
		return 0x2a44 // ⩄ '\u2a44'
	case "careof":
		return 0x2105 // ℅ '\u2105'
	case "caretinsert":
		return 0x2038 // ‸ '\u2038'
	case "caron":
		return 0x02c7 // ˇ '\u02c7'
	case "caron1":
		return 0xf00a //  '\uf00a'
	case "caronbelowcmb":
		return 0x032c // ̬ '\u032c'
	case "caroncmb":
		return 0x030c // ̌ '\u030c'
	case "carriagereturn":
		return 0x21b5 // ↵ '\u21b5'
	case "cbopomofo":
		return 0x3118 // ㄘ '\u3118'
	case "ccaron":
		return 0x010d // č '\u010d'
	case "ccedilla":
		return 0x00e7 // ç '\u00e7'
	case "ccedillaacute":
		return 0x1e09 // ḉ '\u1e09'
	case "ccircle":
		return 0x24d2 // ⓒ '\u24d2'
	case "ccircumflex":
		return 0x0109 // ĉ '\u0109'
	case "ccurl":
		return 0x0255 // ɕ '\u0255'
	case "ccwundercurvearrow":
		return 0x293f // ⤿ '\u293f'
	case "cdot":
		return 0x010b // ċ '\u010b'
	case "cdsquare":
		return 0x33c5 // ㏅ '\u33c5'
	case "cedilla":
		return 0x00b8 // ¸ '\u00b8'
	case "cedilla1":
		return 0xf008 //  '\uf008'
	case "cedilla2":
		return 0xf00d //  '\uf00d'
	case "cedillacmb":
		return 0x0327 // ̧ '\u0327'
	case "ceilingleft":
		return 0x2308 // ⌈ '\u2308'
	case "ceilingright":
		return 0x2309 // ⌉ '\u2309'
	case "cent":
		return 0x00a2 // ¢ '\u00a2'
	case "centigrade":
		return 0x2103 // ℃ '\u2103'
	case "centinferior":
		return 0xf6df //  '\uf6df'
	case "centmonospace":
		return 0xffe0 // ￠ '\uffe0'
	case "centoldstyle":
		return 0xf7a2 //  '\uf7a2'
	case "centreline":
		return 0x2104 // ℄ '\u2104'
	case "centsuperior":
		return 0xf6e0 //  '\uf6e0'
	case "chaarmenian":
		return 0x0579 // չ '\u0579'
	case "chabengali":
		return 0x099b // ছ '\u099b'
	case "chadeva":
		return 0x091b // छ '\u091b'
	case "chagujarati":
		return 0x0a9b // છ '\u0a9b'
	case "chagurmukhi":
		return 0x0a1b // ਛ '\u0a1b'
	case "chbopomofo":
		return 0x3114 // ㄔ '\u3114'
	case "cheabkhasiancyrillic":
		return 0x04bd // ҽ '\u04bd'
	case "checyrillic":
		return 0x0447 // ч '\u0447'
	case "chedescenderabkhasiancyrillic":
		return 0x04bf // ҿ '\u04bf'
	case "chedescendercyrillic":
		return 0x04b7 // ҷ '\u04b7'
	case "chedieresiscyrillic":
		return 0x04f5 // ӵ '\u04f5'
	case "cheharmenian":
		return 0x0573 // ճ '\u0573'
	case "chekhakassiancyrillic":
		return 0x04cc // ӌ '\u04cc'
	case "cheverticalstrokecyrillic":
		return 0x04b9 // ҹ '\u04b9'
	case "chi":
		return 0x03c7 // χ '\u03c7'
	case "chieuchacirclekorean":
		return 0x3277 // ㉷ '\u3277'
	case "chieuchaparenkorean":
		return 0x3217 // ㈗ '\u3217'
	case "chieuchcirclekorean":
		return 0x3269 // ㉩ '\u3269'
	case "chieuchkorean":
		return 0x314a // ㅊ '\u314a'
	case "chieuchparenkorean":
		return 0x3209 // ㈉ '\u3209'
	case "chochangthai":
		return 0x0e0a // ช '\u0e0a'
	case "chochanthai":
		return 0x0e08 // จ '\u0e08'
	case "chochingthai":
		return 0x0e09 // ฉ '\u0e09'
	case "chochoethai":
		return 0x0e0c // ฌ '\u0e0c'
	case "chook":
		return 0x0188 // ƈ '\u0188'
	case "cieucacirclekorean":
		return 0x3276 // ㉶ '\u3276'
	case "cieucaparenkorean":
		return 0x3216 // ㈖ '\u3216'
	case "cieuccirclekorean":
		return 0x3268 // ㉨ '\u3268'
	case "cieuckorean":
		return 0x3148 // ㅈ '\u3148'
	case "cieucparenkorean":
		return 0x3208 // ㈈ '\u3208'
	case "cieucuparenkorean":
		return 0x321c // ㈜ '\u321c'
	case "cirE":
		return 0x29c3 // ⧃ '\u29c3'
	case "cirbot":
		return 0x27df // ⟟ '\u27df'
	case "circeq":
		return 0x2257 // ≗ '\u2257'
	case "circleasterisk":
		return 0x229b // ⊛ '\u229b'
	case "circlebottomhalfblack":
		return 0x25d2 // ◒ '\u25d2'
	case "circlecopyrt":
		return 0x20dd // ⃝ '\u20dd'
	case "circledbullet":
		return 0x29bf // ⦿ '\u29bf'
	case "circleddash":
		return 0x229d // ⊝ '\u229d'
	case "circledivide":
		return 0x2298 // ⊘ '\u2298'
	case "circledownarrow":
		return 0x29ec // ⧬ '\u29ec'
	case "circledparallel":
		return 0x29b7 // ⦷ '\u29b7'
	case "circledrightdot":
		return 0x2686 // ⚆ '\u2686'
	case "circledtwodots":
		return 0x2687 // ⚇ '\u2687'
	case "circledvert":
		return 0x29b6 // ⦶ '\u29b6'
	case "circledwhitebullet":
		return 0x29be // ⦾ '\u29be'
	case "circleequal":
		return 0x229c // ⊜ '\u229c'
	case "circlehbar":
		return 0x29b5 // ⦵ '\u29b5'
	case "circlellquad":
		return 0x25f5 // ◵ '\u25f5'
	case "circlelrquad":
		return 0x25f6 // ◶ '\u25f6'
	case "circlemultiply":
		return 0x2297 // ⊗ '\u2297'
	case "circleonleftarrow":
		return 0x2b30 // ⬰ '\u2b30'
	case "circleonrightarrow":
		return 0x21f4 // ⇴ '\u21f4'
	case "circleot":
		return 0x2299 // ⊙ '\u2299'
	case "circleplus":
		return 0x2295 // ⊕ '\u2295'
	case "circlepostalmark":
		return 0x3036 // 〶 '\u3036'
	case "circlering":
		return 0x229a // ⊚ '\u229a'
	case "circletophalfblack":
		return 0x25d3 // ◓ '\u25d3'
	case "circleulquad":
		return 0x25f4 // ◴ '\u25f4'
	case "circleurquad":
		return 0x25f7 // ◷ '\u25f7'
	case "circleurquadblack":
		return 0x25d4 // ◔ '\u25d4'
	case "circlevertfill":
		return 0x25cd // ◍ '\u25cd'
	case "circlewithlefthalfblack":
		return 0x25d0 // ◐ '\u25d0'
	case "circlewithrighthalfblack":
		return 0x25d1 // ◑ '\u25d1'
	case "circumflex":
		return 0x02c6 // ˆ '\u02c6'
	case "circumflex1":
		return 0xf003 //  '\uf003'
	case "circumflexbelowcmb":
		return 0x032d // ̭ '\u032d'
	case "circumflexcmb":
		return 0x0302 // ̂ '\u0302'
	case "cirfnint":
		return 0x2a10 // ⨐ '\u2a10'
	case "cirmid":
		return 0x2aef // ⫯ '\u2aef'
	case "cirscir":
		return 0x29c2 // ⧂ '\u29c2'
	case "clear":
		return 0x2327 // ⌧ '\u2327'
	case "clickalveolar":
		return 0x01c2 // ǂ '\u01c2'
	case "clickdental":
		return 0x01c0 // ǀ '\u01c0'
	case "clicklateral":
		return 0x01c1 // ǁ '\u01c1'
	case "clickretroflex":
		return 0x01c3 // ǃ '\u01c3'
	case "clockwise":
		return 0x27f3 // ⟳ '\u27f3'
	case "closedvarcap":
		return 0x2a4d // ⩍ '\u2a4d'
	case "closedvarcup":
		return 0x2a4c // ⩌ '\u2a4c'
	case "closedvarcupsmashprod":
		return 0x2a50 // ⩐ '\u2a50'
	case "closure":
		return 0x2050 // ⁐ '\u2050'
	case "club":
		return 0x2663 // ♣ '\u2663'
	case "clubsuitwhite":
		return 0x2667 // ♧ '\u2667'
	case "cmcubedsquare":
		return 0x33a4 // ㎤ '\u33a4'
	case "cmonospace":
		return 0xff43 // ｃ '\uff43'
	case "cmsquaredsquare":
		return 0x33a0 // ㎠ '\u33a0'
	case "coarmenian":
		return 0x0581 // ց '\u0581'
	case "colon":
		return 0x003a // : ':'
	case "coloneq":
		return 0x2254 // ≔ '\u2254'
	case "colonmonetary":
		return 0x20a1 // ₡ '\u20a1'
	case "colonmonospace":
		return 0xff1a // ： '\uff1a'
	case "colonsmall":
		return 0xfe55 // ﹕ '\ufe55'
	case "colontriangularhalfmod":
		return 0x02d1 // ˑ '\u02d1'
	case "colontriangularmod":
		return 0x02d0 // ː '\u02d0'
	case "comma":
		return 0x002c // , ','
	case "commaabovecmb":
		return 0x0313 // ̓ '\u0313'
	case "commaaboverightcmb":
		return 0x0315 // ̕ '\u0315'
	case "commaaccent":
		return 0xf6c3 //  '\uf6c3'
	case "commaarmenian":
		return 0x055d // ՝ '\u055d'
	case "commainferior":
		return 0xf6e1 //  '\uf6e1'
	case "commaminus":
		return 0x2a29 // ⨩ '\u2a29'
	case "commamonospace":
		return 0xff0c // ， '\uff0c'
	case "commareversedabovecmb":
		return 0x0314 // ̔ '\u0314'
	case "commasmall":
		return 0xfe50 // ﹐ '\ufe50'
	case "commasubnosp":
		return 0x0299 // ʙ '\u0299'
	case "commasuperior":
		return 0xf6e2 //  '\uf6e2'
	case "commaturnedabovecmb":
		return 0x0312 // ̒ '\u0312'
	case "commaturnedmod":
		return 0x02bb // ʻ '\u02bb'
	case "complement":
		return 0x2201 // ∁ '\u2201'
	case "concavediamond":
		return 0x27e1 // ⟡ '\u27e1'
	case "concavediamondtickleft":
		return 0x27e2 // ⟢ '\u27e2'
	case "concavediamondtickright":
		return 0x27e3 // ⟣ '\u27e3'
	case "congdot":
		return 0x2a6d // ⩭ '\u2a6d'
	case "congruent":
		return 0x2245 // ≅ '\u2245'
	case "conictaper":
		return 0x2332 // ⌲ '\u2332'
	case "conjquant":
		return 0x2a07 // ⨇ '\u2a07'
	case "contourintegral":
		return 0x222e // ∮ '\u222e'
	case "control":
		return 0x2303 // ⌃ '\u2303'
	case "controlACK":
		return 0x0006 //  '\x06'
	case "controlBEL":
		return 0x0007 //  '\a'
	case "controlBS":
		return 0x0008 //  '\b'
	case "controlCAN":
		return 0x0018 //  '\x18'
	case "controlCR":
		return 0x000d //  '\r'
	case "controlDC1":
		return 0x0011 //  '\x11'
	case "controlDC2":
		return 0x0012 //  '\x12'
	case "controlDC3":
		return 0x0013 //  '\x13'
	case "controlDC4":
		return 0x0014 //  '\x14'
	case "controlDEL":
		return 0x007f //  '\u007f'
	case "controlDLE":
		return 0x0010 //  '\x10'
	case "controlEM":
		return 0x0019 //  '\x19'
	case "controlENQ":
		return 0x0005 //  '\x05'
	case "controlEOT":
		return 0x0004 //  '\x04'
	case "controlESC":
		return 0x001b //  '\x1b'
	case "controlETB":
		return 0x0017 //  '\x17'
	case "controlETX":
		return 0x0003 //  '\x03'
	case "controlFF":
		return 0x000c //  '\f'
	case "controlFS":
		return 0x001c //  '\x1c'
	case "controlGS":
		return 0x001d //  '\x1d'
	case "controlHT":
		return 0x0009 //  '\t'
	case "controlLF":
		return 0x000a //  '\n'
	case "controlNAK":
		return 0x0015 //  '\x15'
	case "controlNULL":
		return 0x0000 //  '\x00'
	case "controlRS":
		return 0x001e //  '\x1e'
	case "controlSI":
		return 0x000f //  '\x0f'
	case "controlSO":
		return 0x000e //  '\x0e'
	case "controlSOT":
		return 0x0002 //  '\x02'
	case "controlSTX":
		return 0x0001 //  '\x01'
	case "controlSUB":
		return 0x001a //  '\x1a'
	case "controlSYN":
		return 0x0016 //  '\x16'
	case "controlUS":
		return 0x001f //  '\x1f'
	case "controlVT":
		return 0x000b //  '\v'
	case "coproduct":
		return 0x2a3f // ⨿ '\u2a3f'
	case "coproductdisplay":
		return 0x2210 // ∐ '\u2210'
	case "copyright":
		return 0x00a9 // © '\u00a9'
	case "copyrightsans":
		return 0xf8e9 //  '\uf8e9'
	case "copyrightserif":
		return 0xf6d9 //  '\uf6d9'
	case "cornerbracketleft":
		return 0x300c // 「 '\u300c'
	case "cornerbracketlefthalfwidth":
		return 0xff62 // ｢ '\uff62'
	case "cornerbracketleftvertical":
		return 0xfe41 // ﹁ '\ufe41'
	case "cornerbracketright":
		return 0x300d // 」 '\u300d'
	case "cornerbracketrighthalfwidth":
		return 0xff63 // ｣ '\uff63'
	case "cornerbracketrightvertical":
		return 0xfe42 // ﹂ '\ufe42'
	case "corporationsquare":
		return 0x337f // ㍿ '\u337f'
	case "cosquare":
		return 0x33c7 // ㏇ '\u33c7'
	case "coverkgsquare":
		return 0x33c6 // ㏆ '\u33c6'
	case "cparen":
		return 0x249e // ⒞ '\u249e'
	case "cruzeiro":
		return 0x20a2 // ₢ '\u20a2'
	case "cstretch":
		return 0x0227 // ȧ '\u0227'
	case "cstretched":
		return 0x0297 // ʗ '\u0297'
	case "csub":
		return 0x2acf // ⫏ '\u2acf'
	case "csube":
		return 0x2ad1 // ⫑ '\u2ad1'
	case "csup":
		return 0x2ad0 // ⫐ '\u2ad0'
	case "csupe":
		return 0x2ad2 // ⫒ '\u2ad2'
	case "cuberoot":
		return 0x221b // ∛ '\u221b'
	case "cupbarcap":
		return 0x2a48 // ⩈ '\u2a48'
	case "cupdot":
		return 0x228d // ⊍ '\u228d'
	case "cupleftarrow":
		return 0x228c // ⊌ '\u228c'
	case "cupovercap":
		return 0x2a46 // ⩆ '\u2a46'
	case "cupvee":
		return 0x2a45 // ⩅ '\u2a45'
	case "curlyand":
		return 0x22cf // ⋏ '\u22cf'
	case "curlyleft":
		return 0x21ab // ↫ '\u21ab'
	case "curlyor":
		return 0x22ce // ⋎ '\u22ce'
	case "curlyright":
		return 0x21ac // ↬ '\u21ac'
	case "currency":
		return 0x00a4 // ¤ '\u00a4'
	case "curvearrowleftplus":
		return 0x293d // ⤽ '\u293d'
	case "curvearrowrightminus":
		return 0x293c // ⤼ '\u293c'
	case "cwcirclearrow":
		return 0x2941 // ⥁ '\u2941'
	case "cwopencirclearrow":
		return 0x21bb // ↻ '\u21bb'
	case "cwrightarcarrow":
		return 0x2938 // ⤸ '\u2938'
	case "cwundercurvearrow":
		return 0x293e // ⤾ '\u293e'
	case "cyrBreve":
		return 0xf6d1 //  '\uf6d1'
	case "cyrFlex":
		return 0xf6d2 //  '\uf6d2'
	case "cyrbreve":
		return 0xf6d4 //  '\uf6d4'
	case "cyrflex":
		return 0xf6d5 //  '\uf6d5'
	case "d":
		return 0x0064 // d 'd'
	case "daarmenian":
		return 0x0564 // դ '\u0564'
	case "dabengali":
		return 0x09a6 // দ '\u09a6'
	case "dadeva":
		return 0x0926 // द '\u0926'
	case "dadfinalarabic":
		return 0xfebe // ﺾ '\ufebe'
	case "dadinitialarabic":
		return 0xfebf // ﺿ '\ufebf'
	case "dadisolated":
		return 0xfebd // ﺽ '\ufebd'
	case "dadmedialarabic":
		return 0xfec0 // ﻀ '\ufec0'
	case "dagesh":
		return 0x05bc // ּ '\u05bc'
	case "dagger":
		return 0x2020 // † '\u2020'
	case "daggerdbl":
		return 0x2021 // ‡ '\u2021'
	case "dagujarati":
		return 0x0aa6 // દ '\u0aa6'
	case "dagurmukhi":
		return 0x0a26 // ਦ '\u0a26'
	case "dahiragana":
		return 0x3060 // だ '\u3060'
	case "dakatakana":
		return 0x30c0 // ダ '\u30c0'
	case "daletdagesh":
		return 0xfb33 // דּ '\ufb33'
	case "daleth":
		return 0x2138 // ℸ '\u2138'
	case "daletqamatshebrew":
		return 0x05d3 // ד '\u05d3'
	case "dalfinalarabic":
		return 0xfeaa // ﺪ '\ufeaa'
	case "dalisolated":
		return 0xfea9 // ﺩ '\ufea9'
	case "dammaarabic":
		return 0x064f // ُ '\u064f'
	case "dammaisolated":
		return 0xfe78 // ﹸ '\ufe78'
	case "dammalow":
		return 0xe821 //  '\ue821'
	case "dammamedial":
		return 0xfe79 // ﹹ '\ufe79'
	case "dammaonhamza":
		return 0xe835 //  '\ue835'
	case "dammatanisolated":
		return 0xfe72 // ﹲ '\ufe72'
	case "dammatanlow":
		return 0xe824 //  '\ue824'
	case "dammatanonhamza":
		return 0xe836 //  '\ue836'
	case "danda":
		return 0x0964 // । '\u0964'
	case "danger":
		return 0x2621 // ☡ '\u2621'
	case "dargalefthebrew":
		return 0x05a7 // ֧ '\u05a7'
	case "dashV":
		return 0x2ae3 // ⫣ '\u2ae3'
	case "dashVdash":
		return 0x27db // ⟛ '\u27db'
	case "dashcolon":
		return 0x2239 // ∹ '\u2239'
	case "dashleftharpoondown":
		return 0x296b // ⥫ '\u296b'
	case "dashrightharpoondown":
		return 0x296d // ⥭ '\u296d'
	case "dasiapneumatacyrilliccmb":
		return 0x0485 // ҅ '\u0485'
	case "dbkarow":
		return 0x290f // ⤏ '\u290f'
	case "dblGrave":
		return 0xf6d3 //  '\uf6d3'
	case "dblanglebracketleft":
		return 0x300a // 《 '\u300a'
	case "dblanglebracketleftvertical":
		return 0xfe3d // ︽ '\ufe3d'
	case "dblanglebracketright":
		return 0x300b // 》 '\u300b'
	case "dblanglebracketrightvertical":
		return 0xfe3e // ︾ '\ufe3e'
	case "dblarchinvertedbelowcmb":
		return 0x032b // ̫ '\u032b'
	case "dblarrowdwn":
		return 0x21ca // ⇊ '\u21ca'
	case "dblarrowheaddown":
		return 0x058a // ֊ '\u058a'
	case "dblarrowheadleft":
		return 0x219e // ↞ '\u219e'
	case "dblarrowheadright":
		return 0x21a0 // ↠ '\u21a0'
	case "dblarrowheadup":
		return 0x0588 //  '\u0588'
	case "dblarrowup":
		return 0x21c8 // ⇈ '\u21c8'
	case "dblbracketleft":
		return 0x27e6 // ⟦ '\u27e6'
	case "dblbracketright":
		return 0x27e7 // ⟧ '\u27e7'
	case "dbldanda":
		return 0x0965 // ॥ '\u0965'
	case "dblgrave":
		return 0xf6d6 //  '\uf6d6'
	case "dblgravecmb":
		return 0x030f // ̏ '\u030f'
	case "dblintegral":
		return 0x222c // ∬ '\u222c'
	case "dbllowlinecmb":
		return 0x0333 // ̳ '\u0333'
	case "dbloverlinecmb":
		return 0x033f // ̿ '\u033f'
	case "dblprimemod":
		return 0x02ba // ʺ '\u02ba'
	case "dblverticalbar":
		return 0x2016 // ‖ '\u2016'
	case "dblverticallineabovecmb":
		return 0x030e // ̎ '\u030e'
	case "dbopomofo":
		return 0x3109 // ㄉ '\u3109'
	case "dbsquare":
		return 0x33c8 // ㏈ '\u33c8'
	case "dcaron":
		return 0x010f // ď '\u010f'
	case "dcaron1":
		return 0xf811 //  '\uf811'
	case "dcedilla":
		return 0x1e11 // ḑ '\u1e11'
	case "dcircle":
		return 0x24d3 // ⓓ '\u24d3'
	case "dcircumflexbelow":
		return 0x1e13 // ḓ '\u1e13'
	case "ddabengali":
		return 0x09a1 // ড '\u09a1'
	case "ddadeva":
		return 0x0921 // ड '\u0921'
	case "ddagujarati":
		return 0x0aa1 // ડ '\u0aa1'
	case "ddagurmukhi":
		return 0x0a21 // ਡ '\u0a21'
	case "ddalfinalarabic":
		return 0xfb89 // ﮉ '\ufb89'
	case "ddddot":
		return 0x20dc // ⃜ '\u20dc'
	case "dddhadeva":
		return 0x095c // ड़ '\u095c'
	case "dddot":
		return 0x20db // ⃛ '\u20db'
	case "ddhabengali":
		return 0x09a2 // ঢ '\u09a2'
	case "ddhadeva":
		return 0x0922 // ढ '\u0922'
	case "ddhagujarati":
		return 0x0aa2 // ઢ '\u0aa2'
	case "ddhagurmukhi":
		return 0x0a22 // ਢ '\u0a22'
	case "ddotaccent":
		return 0x1e0b // ḋ '\u1e0b'
	case "ddotbelow":
		return 0x1e0d // ḍ '\u1e0d'
	case "ddots":
		return 0x22f1 // ⋱ '\u22f1'
	case "ddotseq":
		return 0x2a77 // ⩷ '\u2a77'
	case "decimalseparatorpersian":
		return 0x066b // ٫ '\u066b'
	case "defines":
		return 0x225c // ≜ '\u225c'
	case "degree":
		return 0x00b0 // ° '\u00b0'
	case "degreekelvin":
		return 0x212a // K '\u212a'
	case "dehihebrew":
		return 0x05ad // ֭ '\u05ad'
	case "dehiragana":
		return 0x3067 // で '\u3067'
	case "deicoptic":
		return 0x03ef // ϯ '\u03ef'
	case "dekatakana":
		return 0x30c7 // デ '\u30c7'
	case "delete":
		return 0x05ba // ֺ '\u05ba'
	case "deleteleft":
		return 0x232b // ⌫ '\u232b'
	case "deleteright":
		return 0x2326 // ⌦ '\u2326'
	case "delta":
		return 0x03b4 // δ '\u03b4'
	case "deltaturned":
		return 0x018d // ƍ '\u018d'
	case "denominatorminusonenumeratorbengali":
		return 0x09f8 // ৸ '\u09f8'
	case "dezh":
		return 0x02a4 // ʤ '\u02a4'
	case "dhabengali":
		return 0x09a7 // ধ '\u09a7'
	case "dhadeva":
		return 0x0927 // ध '\u0927'
	case "dhagujarati":
		return 0x0aa7 // ધ '\u0aa7'
	case "dhagurmukhi":
		return 0x0a27 // ਧ '\u0a27'
	case "dhook":
		return 0x0257 // ɗ '\u0257'
	case "diaeresis":
		return 0x0088 //  '\u0088'
	case "dialytikatonoscmb":
		return 0x0344 // ̈́ '\u0344'
	case "diameter":
		return 0x2300 // ⌀ '\u2300'
	case "diamond":
		return 0x2666 // ♦ '\u2666'
	case "diamondbotblack":
		return 0x2b19 // ⬙ '\u2b19'
	case "diamondcdot":
		return 0x27d0 // ⟐ '\u27d0'
	case "diamondleftarrow":
		return 0x291d // ⤝ '\u291d'
	case "diamondleftarrowbar":
		return 0x291f // ⤟ '\u291f'
	case "diamondleftblack":
		return 0x2b16 // ⬖ '\u2b16'
	case "diamondmath":
		return 0x22c4 // ⋄ '\u22c4'
	case "diamondrightblack":
		return 0x2b17 // ⬗ '\u2b17'
	case "diamondsuitwhite":
		return 0x2662 // ♢ '\u2662'
	case "diamondtopblack":
		return 0x2b18 // ⬘ '\u2b18'
	case "dicei":
		return 0x2680 // ⚀ '\u2680'
	case "diceii":
		return 0x2681 // ⚁ '\u2681'
	case "diceiii":
		return 0x2682 // ⚂ '\u2682'
	case "diceiv":
		return 0x2683 // ⚃ '\u2683'
	case "dicev":
		return 0x2684 // ⚄ '\u2684'
	case "dicevi":
		return 0x2685 // ⚅ '\u2685'
	case "dieresis":
		return 0x00a8 // ¨ '\u00a8'
	case "dieresis1":
		return 0xf005 //  '\uf005'
	case "dieresisacute":
		return 0xf6d7 //  '\uf6d7'
	case "dieresisbelowcmb":
		return 0x0324 // ̤ '\u0324'
	case "dieresiscmb":
		return 0x0308 // ̈ '\u0308'
	case "dieresisgrave":
		return 0xf6d8 //  '\uf6d8'
	case "dieresistonos":
		return 0x0385 // ΅ '\u0385'
	case "difference":
		return 0x224f // ≏ '\u224f'
	case "dihiragana":
		return 0x3062 // ぢ '\u3062'
	case "dikatakana":
		return 0x30c2 // ヂ '\u30c2'
	case "disin":
		return 0x22f2 // ⋲ '\u22f2'
	case "disjquant":
		return 0x2a08 // ⨈ '\u2a08'
	case "dittomark":
		return 0x3003 // 〃 '\u3003'
	case "divide":
		return 0x00f7 // ÷ '\u00f7'
	case "dividemultiply":
		return 0x22c7 // ⋇ '\u22c7'
	case "divides":
		return 0x2223 // ∣ '\u2223'
	case "divisionslash":
		return 0x2215 // ∕ '\u2215'
	case "djecyrillic":
		return 0x0452 // ђ '\u0452'
	case "dkshade":
		return 0x2593 // ▓ '\u2593'
	case "dkshade1":
		return 0xf823 //  '\uf823'
	case "dlinebelow":
		return 0x1e0f // ḏ '\u1e0f'
	case "dlsquare":
		return 0x3397 // ㎗ '\u3397'
	case "dmacron":
		return 0x0111 // đ '\u0111'
	case "dmonospace":
		return 0xff44 // ｄ '\uff44'
	case "dnblock":
		return 0x2584 // ▄ '\u2584'
	case "dneightblock":
		return 0x2581 // ▁ '\u2581'
	case "dnfiveeighthblock":
		return 0x2585 // ▅ '\u2585'
	case "dnquarterblock":
		return 0x2582 // ▂ '\u2582'
	case "dnseveneighthblock":
		return 0x2587 // ▇ '\u2587'
	case "dnthreeeighthblock":
		return 0x2583 // ▃ '\u2583'
	case "dnthreequarterblock":
		return 0x2586 // ▆ '\u2586'
	case "dochadathai":
		return 0x0e0e // ฎ '\u0e0e'
	case "dodekthai":
		return 0x0e14 // ด '\u0e14'
	case "dohiragana":
		return 0x3069 // ど '\u3069'
	case "dokatakana":
		return 0x30c9 // ド '\u30c9'
	case "dollar":
		return 0x0024 // $ '$'
	case "dollarinferior":
		return 0xf6e3 //  '\uf6e3'
	case "dollarmonospace":
		return 0xff04 // ＄ '\uff04'
	case "dollaroldstyle":
		return 0xf724 //  '\uf724'
	case "dollarsmall":
		return 0xfe69 // ﹩ '\ufe69'
	case "dollarsuperior":
		return 0xf6e4 //  '\uf6e4'
	case "dong":
		return 0x20ab // ₫ '\u20ab'
	case "dorusquare":
		return 0x3326 // ㌦ '\u3326'
	case "dotaccent":
		return 0x02d9 // ˙ '\u02d9'
	case "dotaccentcmb":
		return 0x0307 // ̇ '\u0307'
	case "dotbelowcomb":
		return 0x0323 // ̣ '\u0323'
	case "dotcircle1":
		return 0xf820 //  '\uf820'
	case "dotequiv":
		return 0x2a67 // ⩧ '\u2a67'
	case "dotkatakana":
		return 0x30fb // ・ '\u30fb'
	case "dotlessi":
		return 0x0131 // ı '\u0131'
	case "dotlessj":
		return 0xf6be //  '\uf6be'
	case "dotlessjstrokehook":
		return 0x0284 // ʄ '\u0284'
	case "dotmath":
		return 0x22c5 // ⋅ '\u22c5'
	case "dotminus":
		return 0x2238 // ∸ '\u2238'
	case "dotplus":
		return 0x2214 // ∔ '\u2214'
	case "dotsim":
		return 0x2a6a // ⩪ '\u2a6a'
	case "dotsminusdots":
		return 0x223a // ∺ '\u223a'
	case "dottedcircle":
		return 0x25cc // ◌ '\u25cc'
	case "dottedsquare":
		return 0x2b1a // ⬚ '\u2b1a'
	case "dottimes":
		return 0x2a30 // ⨰ '\u2a30'
	case "doublebarvee":
		return 0x2a62 // ⩢ '\u2a62'
	case "doubleplus":
		return 0x29fa // ⧺ '\u29fa'
	case "downarrowbar":
		return 0x2913 // ⤓ '\u2913'
	case "downarrowbarred":
		return 0x2908 // ⤈ '\u2908'
	case "downfishtail":
		return 0x297f // ⥿ '\u297f'
	case "downharpoonleftbar":
		return 0x2959 // ⥙ '\u2959'
	case "downharpoonrightbar":
		return 0x2955 // ⥕ '\u2955'
	case "downharpoonsleftright":
		return 0x2965 // ⥥ '\u2965'
	case "downrightcurvedarrow":
		return 0x2935 // ⤵ '\u2935'
	case "downslope":
		return 0x29f9 // ⧹ '\u29f9'
	case "downtackbelowcmb":
		return 0x031e // ̞ '\u031e'
	case "downtackmod":
		return 0x02d5 // ˕ '\u02d5'
	case "downtriangleleftblack":
		return 0x29e8 // ⧨ '\u29e8'
	case "downtrianglerightblack":
		return 0x29e9 // ⧩ '\u29e9'
	case "downuparrows":
		return 0x21f5 // ⇵ '\u21f5'
	case "downupharpoonsleftright":
		return 0x296f // ⥯ '\u296f'
	case "downzigzagarrow":
		return 0x21af // ↯ '\u21af'
	case "dparen":
		return 0x249f // ⒟ '\u249f'
	case "drbkarow":
		return 0x2910 // ⤐ '\u2910'
	case "dsol":
		return 0x29f6 // ⧶ '\u29f6'
	case "dsub":
		return 0x2a64 // ⩤ '\u2a64'
	case "dsuperior":
		return 0xf6eb //  '\uf6eb'
	case "dtail":
		return 0x0256 // ɖ '\u0256'
	case "dtopbar":
		return 0x018c // ƌ '\u018c'
	case "dualmap":
		return 0x29df // ⧟ '\u29df'
	case "duhiragana":
		return 0x3065 // づ '\u3065'
	case "dukatakana":
		return 0x30c5 // ヅ '\u30c5'
	case "dyogh":
		return 0x0234 // ȴ '\u0234'
	case "dz":
		return 0x01f3 // ǳ '\u01f3'
	case "dzaltone":
		return 0x02a3 // ʣ '\u02a3'
	case "dzcaron":
		return 0x01c6 // ǆ '\u01c6'
	case "dzcurl":
		return 0x02a5 // ʥ '\u02a5'
	case "dzeabkhasiancyrillic":
		return 0x04e1 // ӡ '\u04e1'
	case "dzhecyrillic":
		return 0x045f // џ '\u045f'
	case "e":
		return 0x0065 // e 'e'
	case "eacute":
		return 0x00e9 // é '\u00e9'
	case "earth":
		return 0x2641 // ♁ '\u2641'
	case "ebengali":
		return 0x098f // এ '\u098f'
	case "ebopomofo":
		return 0x311c // ㄜ '\u311c'
	case "ebreve":
		return 0x0115 // ĕ '\u0115'
	case "ecandradeva":
		return 0x090d // ऍ '\u090d'
	case "ecandragujarati":
		return 0x0a8d // ઍ '\u0a8d'
	case "ecandravowelsigndeva":
		return 0x0945 // ॅ '\u0945'
	case "ecandravowelsigngujarati":
		return 0x0ac5 // ૅ '\u0ac5'
	case "ecaron":
		return 0x011b // ě '\u011b'
	case "ecedillabreve":
		return 0x1e1d // ḝ '\u1e1d'
	case "echarmenian":
		return 0x0565 // ե '\u0565'
	case "echyiwnarmenian":
		return 0x0587 // և '\u0587'
	case "ecircle":
		return 0x24d4 // ⓔ '\u24d4'
	case "ecircumflex":
		return 0x00ea // ê '\u00ea'
	case "ecircumflexacute":
		return 0x1ebf // ế '\u1ebf'
	case "ecircumflexbelow":
		return 0x1e19 // ḙ '\u1e19'
	case "ecircumflexdotbelow":
		return 0x1ec7 // ệ '\u1ec7'
	case "ecircumflexgrave":
		return 0x1ec1 // ề '\u1ec1'
	case "ecircumflexhookabove":
		return 0x1ec3 // ể '\u1ec3'
	case "ecircumflextilde":
		return 0x1ec5 // ễ '\u1ec5'
	case "ecyrillic":
		return 0x0454 // є '\u0454'
	case "edblgrave":
		return 0x0205 // ȅ '\u0205'
	case "edeva":
		return 0x090f // ए '\u090f'
	case "edieresis":
		return 0x00eb // ë '\u00eb'
	case "edotaccent":
		return 0x0117 // ė '\u0117'
	case "edotbelow":
		return 0x1eb9 // ẹ '\u1eb9'
	case "eegurmukhi":
		return 0x0a0f // ਏ '\u0a0f'
	case "eematragurmukhi":
		return 0x0a47 // ੇ '\u0a47'
	case "egrave":
		return 0x00e8 // è '\u00e8'
	case "egsdot":
		return 0x2a98 // ⪘ '\u2a98'
	case "egujarati":
		return 0x0a8f // એ '\u0a8f'
	case "eharmenian":
		return 0x0567 // է '\u0567'
	case "ehbopomofo":
		return 0x311d // ㄝ '\u311d'
	case "ehiragana":
		return 0x3048 // え '\u3048'
	case "ehookabove":
		return 0x1ebb // ẻ '\u1ebb'
	case "eibopomofo":
		return 0x311f // ㄟ '\u311f'
	case "eight":
		return 0x0038 // 8 '8'
	case "eightbengali":
		return 0x09ee // ৮ '\u09ee'
	case "eightdeva":
		return 0x096e // ८ '\u096e'
	case "eighteencircle":
		return 0x2471 // ⑱ '\u2471'
	case "eighteenparen":
		return 0x2485 // ⒅ '\u2485'
	case "eighteenperiod":
		return 0x2499 // ⒙ '\u2499'
	case "eightgujarati":
		return 0x0aee // ૮ '\u0aee'
	case "eightgurmukhi":
		return 0x0a6e // ੮ '\u0a6e'
	case "eighthackarabic":
		return 0x0668 // ٨ '\u0668'
	case "eighthangzhou":
		return 0x3028 // 〨 '\u3028'
	case "eighthnotebeamed":
		return 0x266b // ♫ '\u266b'
	case "eightideographicparen":
		return 0x3227 // ㈧ '\u3227'
	case "eightinferior":
		return 0x2088 // ₈ '\u2088'
	case "eightmonospace":
		return 0xff18 // ８ '\uff18'
	case "eightoldstyle":
		return 0xf738 //  '\uf738'
	case "eightparen":
		return 0x247b // ⑻ '\u247b'
	case "eightperiod":
		return 0x248f // ⒏ '\u248f'
	case "eightpersian":
		return 0x06f8 // ۸ '\u06f8'
	case "eightroman":
		return 0x2177 // ⅷ '\u2177'
	case "eightsuperior":
		return 0x2078 // ⁸ '\u2078'
	case "eightthai":
		return 0x0e58 // ๘ '\u0e58'
	case "einvertedbreve":
		return 0x0207 // ȇ '\u0207'
	case "eiotifiedcyrillic":
		return 0x0465 // ѥ '\u0465'
	case "ekatakana":
		return 0x30a8 // エ '\u30a8'
	case "ekatakanahalfwidth":
		return 0xff74 // ｴ '\uff74'
	case "ekonkargurmukhi":
		return 0x0a74 // ੴ '\u0a74'
	case "ekorean":
		return 0x3154 // ㅔ '\u3154'
	case "elcyrillic":
		return 0x043b // л '\u043b'
	case "element":
		return 0x2208 // ∈ '\u2208'
	case "elevencircle":
		return 0x246a // ⑪ '\u246a'
	case "elevenparen":
		return 0x247e // ⑾ '\u247e'
	case "elevenperiod":
		return 0x2492 // ⒒ '\u2492'
	case "elevenroman":
		return 0x217a // ⅺ '\u217a'
	case "elinters":
		return 0x23e7 // ⏧ '\u23e7'
	case "ellipsis":
		return 0x2026 // … '\u2026'
	case "ellipsisvertical":
		return 0x22ee // ⋮ '\u22ee'
	case "elsdot":
		return 0x2a97 // ⪗ '\u2a97'
	case "emacron":
		return 0x0113 // ē '\u0113'
	case "emacronacute":
		return 0x1e17 // ḗ '\u1e17'
	case "emacrongrave":
		return 0x1e15 // ḕ '\u1e15'
	case "emcyrillic":
		return 0x043c // м '\u043c'
	case "emdash":
		return 0x2014 // — '\u2014'
	case "emdashvertical":
		return 0xfe31 // ︱ '\ufe31'
	case "emonospace":
		return 0xff45 // ｅ '\uff45'
	case "emphasismarkarmenian":
		return 0x055b // ՛ '\u055b'
	case "emptyset":
		return 0x2205 // ∅ '\u2205'
	case "emptysetoarr":
		return 0x29b3 // ⦳ '\u29b3'
	case "emptysetoarrl":
		return 0x29b4 // ⦴ '\u29b4'
	case "emptysetobar":
		return 0x29b1 // ⦱ '\u29b1'
	case "emptysetocirc":
		return 0x29b2 // ⦲ '\u29b2'
	case "emptyslot":
		return 0xd801 //  '\ufffd'
	case "emquad":
		return 0x2001 //  '\u2001'
	case "emspace":
		return 0x2003 //  '\u2003'
	case "enbopomofo":
		return 0x3123 // ㄣ '\u3123'
	case "enclosediamond":
		return 0x20df // ⃟ '\u20df'
	case "enclosesquare":
		return 0x20de // ⃞ '\u20de'
	case "enclosetriangle":
		return 0x20e4 // ⃤ '\u20e4'
	case "endash":
		return 0x2013 // – '\u2013'
	case "endashvertical":
		return 0xfe32 // ︲ '\ufe32'
	case "endescendercyrillic":
		return 0x04a3 // ң '\u04a3'
	case "eng":
		return 0x014b // ŋ '\u014b'
	case "engbopomofo":
		return 0x3125 // ㄥ '\u3125'
	case "enghecyrillic":
		return 0x04a5 // ҥ '\u04a5'
	case "enhookcyrillic":
		return 0x04c8 // ӈ '\u04c8'
	case "enquad":
		return 0x2000 //  '\u2000'
	case "enspace":
		return 0x2002 //  '\u2002'
	case "eogonek":
		return 0x0119 // ę '\u0119'
	case "eokorean":
		return 0x3153 // ㅓ '\u3153'
	case "eopen":
		return 0x025b // ɛ '\u025b'
	case "eopenclosed":
		return 0x029a // ʚ '\u029a'
	case "eopenreversed":
		return 0x025c // ɜ '\u025c'
	case "eopenreversedclosed":
		return 0x025e // ɞ '\u025e'
	case "eopenreversedhook":
		return 0x025d // ɝ '\u025d'
	case "eparen":
		return 0x24a0 // ⒠ '\u24a0'
	case "eparsl":
		return 0x29e3 // ⧣ '\u29e3'
	case "epsilon":
		return 0x03b5 // ε '\u03b5'
	case "epsilon1":
		return 0x03f5 // ϵ '\u03f5'
	case "epsilonclosed":
		return 0x022a // Ȫ '\u022a'
	case "epsiloninv":
		return 0x03f6 // ϶ '\u03f6'
	case "epsilontonos":
		return 0x03ad // έ '\u03ad'
	case "eqcolon":
		return 0x2255 // ≕ '\u2255'
	case "eqdef":
		return 0x225d // ≝ '\u225d'
	case "eqdot":
		return 0x2a66 // ⩦ '\u2a66'
	case "eqeq":
		return 0x2a75 // ⩵ '\u2a75'
	case "eqeqeq":
		return 0x2a76 // ⩶ '\u2a76'
	case "eqgtr":
		return 0x22dd // ⋝ '\u22dd'
	case "eqless":
		return 0x22dc // ⋜ '\u22dc'
	case "eqqgtr":
		return 0x2a9a // ⪚ '\u2a9a'
	case "eqqless":
		return 0x2a99 // ⪙ '\u2a99'
	case "eqqplus":
		return 0x2a71 // ⩱ '\u2a71'
	case "eqqsim":
		return 0x2a73 // ⩳ '\u2a73'
	case "eqqslantgtr":
		return 0x2a9c // ⪜ '\u2a9c'
	case "eqqslantless":
		return 0x2a9b // ⪛ '\u2a9b'
	case "equal":
		return 0x003d // = '='
	case "equalleftarrow":
		return 0x2b40 // ⭀ '\u2b40'
	case "equalmonospace":
		return 0xff1d // ＝ '\uff1d'
	case "equalorfollows":
		return 0x22df // ⋟ '\u22df'
	case "equalorgreater":
		return 0x2a96 // ⪖ '\u2a96'
	case "equalorless":
		return 0x2a95 // ⪕ '\u2a95'
	case "equalorprecedes":
		return 0x22de // ⋞ '\u22de'
	case "equalorsimilar":
		return 0x2242 // ≂ '\u2242'
	case "equalparallel":
		return 0x22d5 // ⋕ '\u22d5'
	case "equalrightarrow":
		return 0x2971 // ⥱ '\u2971'
	case "equalsmall":
		return 0xfe66 // ﹦ '\ufe66'
	case "equalsub":
		return 0x208c // ₌ '\u208c'
	case "equalsuperior":
		return 0x207c // ⁼ '\u207c'
	case "equivDD":
		return 0x2a78 // ⩸ '\u2a78'
	case "equivVert":
		return 0x2a68 // ⩨ '\u2a68'
	case "equivVvert":
		return 0x2a69 // ⩩ '\u2a69'
	case "equivalence":
		return 0x2261 // ≡ '\u2261'
	case "equivasymptotic":
		return 0x224d // ≍ '\u224d'
	case "eqvparsl":
		return 0x29e5 // ⧥ '\u29e5'
	case "erbopomofo":
		return 0x3126 // ㄦ '\u3126'
	case "ercyrillic":
		return 0x0440 // р '\u0440'
	case "ereversed":
		return 0x0258 // ɘ '\u0258'
	case "ereversedcyrillic":
		return 0x044d // э '\u044d'
	case "errbarblackcircle":
		return 0x29f3 // ⧳ '\u29f3'
	case "errbarblackdiamond":
		return 0x29f1 // ⧱ '\u29f1'
	case "errbarblacksquare":
		return 0x29ef // ⧯ '\u29ef'
	case "errbarcircle":
		return 0x29f2 // ⧲ '\u29f2'
	case "errbardiamond":
		return 0x29f0 // ⧰ '\u29f0'
	case "errbarsquare":
		return 0x29ee // ⧮ '\u29ee'
	case "escyrillic":
		return 0x0441 // с '\u0441'
	case "esdescendercyrillic":
		return 0x04ab // ҫ '\u04ab'
	case "esh":
		return 0x0283 // ʃ '\u0283'
	case "eshcurl":
		return 0x0286 // ʆ '\u0286'
	case "eshortdeva":
		return 0x090e // ऎ '\u090e'
	case "eshortvowelsigndeva":
		return 0x0946 // ॆ '\u0946'
	case "eshreversedloop":
		return 0x01aa // ƪ '\u01aa'
	case "eshsquatreversed":
		return 0x0285 // ʅ '\u0285'
	case "esmallhiragana":
		return 0x3047 // ぇ '\u3047'
	case "esmallkatakana":
		return 0x30a7 // ェ '\u30a7'
	case "esmallkatakanahalfwidth":
		return 0xff6a // ｪ '\uff6a'
	case "estimated":
		return 0x212e // ℮ '\u212e'
	case "esuperior":
		return 0xf6ec //  '\uf6ec'
	case "eta":
		return 0x03b7 // η '\u03b7'
	case "etarmenian":
		return 0x0568 // ը '\u0568'
	case "etatonos":
		return 0x03ae // ή '\u03ae'
	case "eth":
		return 0x00f0 // ð '\u00f0'
	case "etilde":
		return 0x1ebd // ẽ '\u1ebd'
	case "etildebelow":
		return 0x1e1b // ḛ '\u1e1b'
	case "etnahtalefthebrew":
		return 0x0591 // ֑ '\u0591'
	case "eturned":
		return 0x01dd // ǝ '\u01dd'
	case "eukorean":
		return 0x3161 // ㅡ '\u3161'
	case "eurocurrency":
		return 0x20a0 // ₠ '\u20a0'
	case "evowelsignbengali":
		return 0x09c7 // ে '\u09c7'
	case "evowelsigndeva":
		return 0x0947 // े '\u0947'
	case "evowelsigngujarati":
		return 0x0ac7 // ે '\u0ac7'
	case "exclam":
		return 0x0021 // ! '!'
	case "exclamarmenian":
		return 0x055c // ՜ '\u055c'
	case "exclamdbl":
		return 0x203c // ‼ '\u203c'
	case "exclamdown":
		return 0x00a1 // ¡ '\u00a1'
	case "exclamdownsmall":
		return 0xf7a1 //  '\uf7a1'
	case "exclammonospace":
		return 0xff01 // ！ '\uff01'
	case "exclamsmall":
		return 0xf721 //  '\uf721'
	case "existential":
		return 0x2203 // ∃ '\u2203'
	case "ezh":
		return 0x0292 // ʒ '\u0292'
	case "ezhcaron":
		return 0x01ef // ǯ '\u01ef'
	case "ezhcurl":
		return 0x0293 // ʓ '\u0293'
	case "ezhreversed":
		return 0x01b9 // ƹ '\u01b9'
	case "ezhtail":
		return 0x01ba // ƺ '\u01ba'
	case "f":
		return 0x0066 // f 'f'
	case "f70e":
		return 0xf70e //  '\uf70e'
	case "f70a":
		return 0xf70a //  '\uf70a'
	case "f70c":
		return 0xf70c //  '\uf70c'
	case "f70d":
		return 0xf70d //  '\uf70d'
	case "f70b":
		return 0xf70b //  '\uf70b'
	case "f70f":
		return 0xf70f //  '\uf70f'
	case "f71c":
		return 0xf71c //  '\uf71c'
	case "f71a":
		return 0xf71a //  '\uf71a'
	case "f71d":
		return 0xf71d //  '\uf71d'
	case "f700":
		return 0xf700 //  '\uf700'
	case "f701":
		return 0xf701 //  '\uf701'
	case "f702":
		return 0xf702 //  '\uf702'
	case "f703":
		return 0xf703 //  '\uf703'
	case "f704":
		return 0xf704 //  '\uf704'
	case "f705":
		return 0xf705 //  '\uf705'
	case "f706":
		return 0xf706 //  '\uf706'
	case "f707":
		return 0xf707 //  '\uf707'
	case "f708":
		return 0xf708 //  '\uf708'
	case "f709":
		return 0xf709 //  '\uf709'
	case "f710":
		return 0xf710 //  '\uf710'
	case "f711":
		return 0xf711 //  '\uf711'
	case "f712":
		return 0xf712 //  '\uf712'
	case "f713":
		return 0xf713 //  '\uf713'
	case "f714":
		return 0xf714 //  '\uf714'
	case "f715":
		return 0xf715 //  '\uf715'
	case "f716":
		return 0xf716 //  '\uf716'
	case "f717":
		return 0xf717 //  '\uf717'
	case "f718":
		return 0xf718 //  '\uf718'
	case "f719":
		return 0xf719 //  '\uf719'
	case "fadeva":
		return 0x095e // फ़ '\u095e'
	case "fagurmukhi":
		return 0x0a5e // ਫ਼ '\u0a5e'
	case "fahrenheit":
		return 0x2109 // ℉ '\u2109'
	case "farsiyeh":
		return 0x06cc // ی '\u06cc'
	case "farsiyehfinal":
		return 0xfbfd // ﯽ '\ufbfd'
	case "farsiyehisolated":
		return 0xfbfc // ﯼ '\ufbfc'
	case "fathahontatweel":
		return 0xfe77 // ﹷ '\ufe77'
	case "fathaisolated":
		return 0xfe76 // ﹶ '\ufe76'
	case "fathalow":
		return 0xe820 //  '\ue820'
	case "fathalowarabic":
		return 0x064e // َ '\u064e'
	case "fathaonhamza":
		return 0xe832 //  '\ue832'
	case "fathatanisolated":
		return 0xfe70 // ﹰ '\ufe70'
	case "fathatanlow":
		return 0xe823 //  '\ue823'
	case "fathatanonhamza":
		return 0xe833 //  '\ue833'
	case "fbopomofo":
		return 0x3108 // ㄈ '\u3108'
	case "fbowtie":
		return 0x29d3 // ⧓ '\u29d3'
	case "fcircle":
		return 0x24d5 // ⓕ '\u24d5'
	case "fcmp":
		return 0x2a3e // ⨾ '\u2a3e'
	case "fdiagovnearrow":
		return 0x292f // ⤯ '\u292f'
	case "fdiagovrdiag":
		return 0x292c // ⤬ '\u292c'
	case "fdotaccent":
		return 0x1e1f // ḟ '\u1e1f'
	case "feharmenian":
		return 0x0586 // ֆ '\u0586'
	case "fehfinalarabic":
		return 0xfed2 // ﻒ '\ufed2'
	case "fehinitialarabic":
		return 0xfed3 // ﻓ '\ufed3'
	case "fehisolated":
		return 0xfed1 // ﻑ '\ufed1'
	case "fehmedialarabic":
		return 0xfed4 // ﻔ '\ufed4'
	case "fehwithalefmaksuraisolated":
		return 0xfc31 // ﰱ '\ufc31'
	case "fehwithyehisolated":
		return 0xfc32 // ﰲ '\ufc32'
	case "feicoptic":
		return 0x03e5 // ϥ '\u03e5'
	case "female":
		return 0x2640 // ♀ '\u2640'
	case "ff":
		return 0xfb00 // ﬀ '\ufb00'
	case "ffi":
		return 0xfb03 // ﬃ '\ufb03'
	case "ffl":
		return 0xfb04 // ﬄ '\ufb04'
	case "fi":
		return 0xfb01 // ﬁ '\ufb01'
	case "fifteencircle":
		return 0x246e // ⑮ '\u246e'
	case "fifteenparen":
		return 0x2482 // ⒂ '\u2482'
	case "fifteenperiod":
		return 0x2496 // ⒖ '\u2496'
	case "figuredash":
		return 0x2012 // ‒ '\u2012'
	case "figurespace":
		return 0x2007 //  '\u2007'
	case "finalkafdageshhebrew":
		return 0xfb3a // ךּ '\ufb3a'
	case "finalkafwithqamats":
		return 0xe803 //  '\ue803'
	case "finalkafwithsheva":
		return 0xe802 //  '\ue802'
	case "finalmemhebrew":
		return 0x05dd // ם '\u05dd'
	case "finalpehebrew":
		return 0x05e3 // ף '\u05e3'
	case "finaltsadi":
		return 0x05e5 // ץ '\u05e5'
	case "fint":
		return 0x2a0f // ⨏ '\u2a0f'
	case "firsttonechinese":
		return 0x02c9 // ˉ '\u02c9'
	case "fisheye":
		return 0x25c9 // ◉ '\u25c9'
	case "five":
		return 0x0035 // 5 '5'
	case "fivearabic":
		return 0x0665 // ٥ '\u0665'
	case "fivebengali":
		return 0x09eb // ৫ '\u09eb'
	case "fivedeva":
		return 0x096b // ५ '\u096b'
	case "fiveeighths":
		return 0x215d // ⅝ '\u215d'
	case "fivegujarati":
		return 0x0aeb // ૫ '\u0aeb'
	case "fivegurmukhi":
		return 0x0a6b // ੫ '\u0a6b'
	case "fivehangzhou":
		return 0x3025 // 〥 '\u3025'
	case "fiveideographicparen":
		return 0x3224 // ㈤ '\u3224'
	case "fiveinferior":
		return 0x2085 // ₅ '\u2085'
	case "fivemonospace":
		return 0xff15 // ５ '\uff15'
	case "fiveoldstyle":
		return 0xf735 //  '\uf735'
	case "fiveparen":
		return 0x2478 // ⑸ '\u2478'
	case "fiveperiod":
		return 0x248c // ⒌ '\u248c'
	case "fivepersian":
		return 0x06f5 // ۵ '\u06f5'
	case "fiveroman":
		return 0x2174 // ⅴ '\u2174'
	case "fivesixth":
		return 0x215a // ⅚ '\u215a'
	case "fivesuperior":
		return 0x2075 // ⁵ '\u2075'
	case "fivethai":
		return 0x0e55 // ๕ '\u0e55'
	case "fl":
		return 0xfb02 // ﬂ '\ufb02'
	case "floorleft":
		return 0x230a // ⌊ '\u230a'
	case "floorright":
		return 0x230b // ⌋ '\u230b'
	case "florin":
		return 0x0192 // ƒ '\u0192'
	case "fltns":
		return 0x23e5 // ⏥ '\u23e5'
	case "fmonospace":
		return 0xff46 // ｆ '\uff46'
	case "fmsquare":
		return 0x3399 // ㎙ '\u3399'
	case "fofanthai":
		return 0x0e1f // ฟ '\u0e1f'
	case "fofathai":
		return 0x0e1d // ฝ '\u0e1d'
	case "follownotdbleqv":
		return 0x2aba // ⪺ '\u2aba'
	case "follownotslnteql":
		return 0x2ab6 // ⪶ '\u2ab6'
	case "followornoteqvlnt":
		return 0x22e9 // ⋩ '\u22e9'
	case "followsequal":
		return 0x2ab0 // ⪰ '\u2ab0'
	case "followsorcurly":
		return 0x227d // ≽ '\u227d'
	case "followsorequal":
		return 0x227f // ≿ '\u227f'
	case "fongmanthai":
		return 0x0e4f // ๏ '\u0e4f'
	case "forces":
		return 0x22a9 // ⊩ '\u22a9'
	case "forcesbar":
		return 0x22aa // ⊪ '\u22aa'
	case "fork":
		return 0x22d4 // ⋔ '\u22d4'
	case "forks":
		return 0x2adc // ⫝̸ '\u2adc'
	case "forksnot":
		return 0x2add // ⫝ '\u2add'
	case "forkv":
		return 0x2ad9 // ⫙ '\u2ad9'
	case "four":
		return 0x0034 // 4 '4'
	case "fourarabic":
		return 0x0664 // ٤ '\u0664'
	case "fourbengali":
		return 0x09ea // ৪ '\u09ea'
	case "fourdeva":
		return 0x096a // ४ '\u096a'
	case "fourfifths":
		return 0x2158 // ⅘ '\u2158'
	case "fourgujarati":
		return 0x0aea // ૪ '\u0aea'
	case "fourgurmukhi":
		return 0x0a6a // ੪ '\u0a6a'
	case "fourhangzhou":
		return 0x3024 // 〤 '\u3024'
	case "fourideographicparen":
		return 0x3223 // ㈣ '\u3223'
	case "fourinferior":
		return 0x2084 // ₄ '\u2084'
	case "fourmonospace":
		return 0xff14 // ４ '\uff14'
	case "fournumeratorbengali":
		return 0x09f7 // ৷ '\u09f7'
	case "fouroldstyle":
		return 0xf734 //  '\uf734'
	case "fourparen":
		return 0x2477 // ⑷ '\u2477'
	case "fourperemspace":
		return 0x2005 //  '\u2005'
	case "fourperiod":
		return 0x248b // ⒋ '\u248b'
	case "fourpersian":
		return 0x06f4 // ۴ '\u06f4'
	case "fourroman":
		return 0x2173 // ⅳ '\u2173'
	case "foursuperior":
		return 0x2074 // ⁴ '\u2074'
	case "fourteencircle":
		return 0x246d // ⑭ '\u246d'
	case "fourteenparen":
		return 0x2481 // ⒁ '\u2481'
	case "fourteenperiod":
		return 0x2495 // ⒕ '\u2495'
	case "fourthai":
		return 0x0e54 // ๔ '\u0e54'
	case "fourthroot":
		return 0x221c // ∜ '\u221c'
	case "fourthtonechinese":
		return 0x02cb // ˋ '\u02cb'
	case "fourvdots":
		return 0x2999 // ⦙ '\u2999'
	case "fparen":
		return 0x24a1 // ⒡ '\u24a1'
	case "fraction":
		return 0x2044 // ⁄ '\u2044'
	case "franc":
		return 0x20a3 // ₣ '\u20a3'
	case "fronted":
		return 0x024b // ɋ '\u024b'
	case "fullouterjoin":
		return 0x27d7 // ⟗ '\u27d7'
	case "g":
		return 0x0067 // g 'g'
	case "gabengali":
		return 0x0997 // গ '\u0997'
	case "gacute":
		return 0x01f5 // ǵ '\u01f5'
	case "gadeva":
		return 0x0917 // ग '\u0917'
	case "gafarabic":
		return 0x06af // گ '\u06af'
	case "gaffinalarabic":
		return 0xfb93 // ﮓ '\ufb93'
	case "gafinitialarabic":
		return 0xfb94 // ﮔ '\ufb94'
	case "gafisolated":
		return 0xfb92 // ﮒ '\ufb92'
	case "gafmedialarabic":
		return 0xfb95 // ﮕ '\ufb95'
	case "gagujarati":
		return 0x0a97 // ગ '\u0a97'
	case "gagurmukhi":
		return 0x0a17 // ਗ '\u0a17'
	case "gahiragana":
		return 0x304c // が '\u304c'
	case "gakatakana":
		return 0x30ac // ガ '\u30ac'
	case "gamma":
		return 0x03b3 // γ '\u03b3'
	case "gammalatinsmall":
		return 0x0263 // ɣ '\u0263'
	case "gammasuperior":
		return 0x02e0 // ˠ '\u02e0'
	case "gangiacoptic":
		return 0x03eb // ϫ '\u03eb'
	case "gbopomofo":
		return 0x310d // ㄍ '\u310d'
	case "gbreve":
		return 0x011f // ğ '\u011f'
	case "gcaron":
		return 0x01e7 // ǧ '\u01e7'
	case "gcedilla":
		return 0x0123 // ģ '\u0123'
	case "gcircle":
		return 0x24d6 // ⓖ '\u24d6'
	case "gcircumflex":
		return 0x011d // ĝ '\u011d'
	case "gdot":
		return 0x0121 // ġ '\u0121'
	case "gebar":
		return 0x03cf // Ϗ '\u03cf'
	case "gehiragana":
		return 0x3052 // げ '\u3052'
	case "gekatakana":
		return 0x30b2 // ゲ '\u30b2'
	case "geomequivalent":
		return 0x224e // ≎ '\u224e'
	case "geometricallyequal":
		return 0x2251 // ≑ '\u2251'
	case "geqqslant":
		return 0x2afa // ⫺ '\u2afa'
	case "gereshaccenthebrew":
		return 0x059c // ֜ '\u059c'
	case "gereshhebrew":
		return 0x05f3 // ׳ '\u05f3'
	case "gereshmuqdamhebrew":
		return 0x059d // ֝ '\u059d'
	case "germandbls":
		return 0x00df // ß '\u00df'
	case "gershayimaccenthebrew":
		return 0x059e // ֞ '\u059e'
	case "gershayimhebrew":
		return 0x05f4 // ״ '\u05f4'
	case "gescc":
		return 0x2aa9 // ⪩ '\u2aa9'
	case "gesdot":
		return 0x2a80 // ⪀ '\u2a80'
	case "gesdoto":
		return 0x2a82 // ⪂ '\u2a82'
	case "gesdotol":
		return 0x2a84 // ⪄ '\u2a84'
	case "gesles":
		return 0x2a94 // ⪔ '\u2a94'
	case "getamark":
		return 0x3013 // 〓 '\u3013'
	case "ggg":
		return 0x22d9 // ⋙ '\u22d9'
	case "gggnest":
		return 0x2af8 // ⫸ '\u2af8'
	case "ghabengali":
		return 0x0998 // ঘ '\u0998'
	case "ghadarmenian":
		return 0x0572 // ղ '\u0572'
	case "ghadeva":
		return 0x0918 // घ '\u0918'
	case "ghagujarati":
		return 0x0a98 // ઘ '\u0a98'
	case "ghagurmukhi":
		return 0x0a18 // ਘ '\u0a18'
	case "ghainarabic":
		return 0x063a // غ '\u063a'
	case "ghainfinalarabic":
		return 0xfece // ﻎ '\ufece'
	case "ghaininitialarabic":
		return 0xfecf // ﻏ '\ufecf'
	case "ghainisolated":
		return 0xfecd // ﻍ '\ufecd'
	case "ghainmedialarabic":
		return 0xfed0 // ﻐ '\ufed0'
	case "ghemiddlehookcyrillic":
		return 0x0495 // ҕ '\u0495'
	case "ghestrokecyrillic":
		return 0x0493 // ғ '\u0493'
	case "gheupturncyrillic":
		return 0x0491 // ґ '\u0491'
	case "ghhadeva":
		return 0x095a // ग़ '\u095a'
	case "ghhagurmukhi":
		return 0x0a5a // ਗ਼ '\u0a5a'
	case "ghook":
		return 0x0260 // ɠ '\u0260'
	case "ghzsquare":
		return 0x3393 // ㎓ '\u3393'
	case "gihiragana":
		return 0x304e // ぎ '\u304e'
	case "gikatakana":
		return 0x30ae // ギ '\u30ae'
	case "gimarmenian":
		return 0x0563 // գ '\u0563'
	case "gimel":
		return 0x2137 // ℷ '\u2137'
	case "gimeldageshhebrew":
		return 0xfb32 // גּ '\ufb32'
	case "gjecyrillic":
		return 0x0453 // ѓ '\u0453'
	case "glE":
		return 0x2a92 // ⪒ '\u2a92'
	case "gla":
		return 0x2aa5 // ⪥ '\u2aa5'
	case "gleichstark":
		return 0x29e6 // ⧦ '\u29e6'
	case "glj":
		return 0x2aa4 // ⪤ '\u2aa4'
	case "glottal":
		return 0x0249 // ɉ '\u0249'
	case "glottalinvertedstroke":
		return 0x01be // ƾ '\u01be'
	case "glottalrev":
		return 0x024a // Ɋ '\u024a'
	case "glottalstop":
		return 0x0294 // ʔ '\u0294'
	case "glottalstopbar":
		return 0x0231 // ȱ '\u0231'
	case "glottalstopbarrev":
		return 0x0232 // Ȳ '\u0232'
	case "glottalstopinv":
		return 0x0226 // Ȧ '\u0226'
	case "glottalstopinverted":
		return 0x0296 // ʖ '\u0296'
	case "glottalstopmod":
		return 0x02c0 // ˀ '\u02c0'
	case "glottalstopreversed":
		return 0x0295 // ʕ '\u0295'
	case "glottalstopreversedmod":
		return 0x02c1 // ˁ '\u02c1'
	case "glottalstopreversedsuperior":
		return 0x02e4 // ˤ '\u02e4'
	case "glottalstoprevinv":
		return 0x0225 // ȥ '\u0225'
	case "glottalstopstroke":
		return 0x02a1 // ʡ '\u02a1'
	case "glottalstopstrokereversed":
		return 0x02a2 // ʢ '\u02a2'
	case "gmacron":
		return 0x1e21 // ḡ '\u1e21'
	case "gmonospace":
		return 0xff47 // ｇ '\uff47'
	case "gnsim":
		return 0x22e7 // ⋧ '\u22e7'
	case "gohiragana":
		return 0x3054 // ご '\u3054'
	case "gokatakana":
		return 0x30b4 // ゴ '\u30b4'
	case "gparen":
		return 0x24a2 // ⒢ '\u24a2'
	case "gpasquare":
		return 0x33ac // ㎬ '\u33ac'
	case "gradient":
		return 0x2207 // ∇ '\u2207'
	case "grave":
		return 0x0060 // ` '`'
	case "gravebelowcmb":
		return 0x0316 // ̖ '\u0316'
	case "gravecmb":
		return 0x0300 // ̀ '\u0300'
	case "gravedeva":
		return 0x0953 // ॓ '\u0953'
	case "graveleftnosp":
		return 0x02b3 // ʳ '\u02b3'
	case "gravelowmod":
		return 0x02ce // ˎ '\u02ce'
	case "gravemonospace":
		return 0xff40 // ｀ '\uff40'
	case "gravetonecmb":
		return 0x0340 // ̀ '\u0340'
	case "greater":
		return 0x003e // > '>'
	case "greaterdbleqlless":
		return 0x2a8c // ⪌ '\u2a8c'
	case "greaterdot":
		return 0x22d7 // ⋗ '\u22d7'
	case "greaterequal":
		return 0x2265 // ≥ '\u2265'
	case "greaterequalorless":
		return 0x22db // ⋛ '\u22db'
	case "greatermonospace":
		return 0xff1e // ＞ '\uff1e'
	case "greaternotdblequal":
		return 0x2a8a // ⪊ '\u2a8a'
	case "greaternotequal":
		return 0x2a88 // ⪈ '\u2a88'
	case "greaterorapproxeql":
		return 0x2a86 // ⪆ '\u2a86'
	case "greaterorequivalent":
		return 0x2273 // ≳ '\u2273'
	case "greaterorless":
		return 0x2277 // ≷ '\u2277'
	case "greaterornotdbleql":
		return 0x2269 // ≩ '\u2269'
	case "greateroverequal":
		return 0x2267 // ≧ '\u2267'
	case "greatersmall":
		return 0xfe65 // ﹥ '\ufe65'
	case "gscript":
		return 0x0261 // ɡ '\u0261'
	case "gsime":
		return 0x2a8e // ⪎ '\u2a8e'
	case "gsiml":
		return 0x2a90 // ⪐ '\u2a90'
	case "gstroke":
		return 0x01e5 // ǥ '\u01e5'
	case "gtcc":
		return 0x2aa7 // ⪧ '\u2aa7'
	case "gtcir":
		return 0x2a7a // ⩺ '\u2a7a'
	case "gtlpar":
		return 0x29a0 // ⦠ '\u29a0'
	case "gtquest":
		return 0x2a7c // ⩼ '\u2a7c'
	case "gtrarr":
		return 0x2978 // ⥸ '\u2978'
	case "guhiragana":
		return 0x3050 // ぐ '\u3050'
	case "guillemotleft":
		return 0x00ab // « '\u00ab'
	case "guillemotright":
		return 0x00bb // » '\u00bb'
	case "guilsinglleft":
		return 0x2039 // ‹ '\u2039'
	case "guilsinglright":
		return 0x203a // › '\u203a'
	case "gukatakana":
		return 0x30b0 // グ '\u30b0'
	case "guramusquare":
		return 0x3318 // ㌘ '\u3318'
	case "gysquare":
		return 0x33c9 // ㏉ '\u33c9'
	case "h":
		return 0x0068 // h 'h'
	case "haabkhasiancyrillic":
		return 0x04a9 // ҩ '\u04a9'
	case "haaltonearabic":
		return 0x06c1 // ہ '\u06c1'
	case "habengali":
		return 0x09b9 // হ '\u09b9'
	case "haceksubnosp":
		return 0x029f // ʟ '\u029f'
	case "hadescendercyrillic":
		return 0x04b3 // ҳ '\u04b3'
	case "hadeva":
		return 0x0939 // ह '\u0939'
	case "hagujarati":
		return 0x0ab9 // હ '\u0ab9'
	case "hagurmukhi":
		return 0x0a39 // ਹ '\u0a39'
	case "hahfinalarabic":
		return 0xfea2 // ﺢ '\ufea2'
	case "hahinitialarabic":
		return 0xfea3 // ﺣ '\ufea3'
	case "hahiragana":
		return 0x306f // は '\u306f'
	case "hahisolated":
		return 0xfea1 // ﺡ '\ufea1'
	case "hahmedialarabic":
		return 0xfea4 // ﺤ '\ufea4'
	case "hahwithmeeminitial":
		return 0xfcaa // ﲪ '\ufcaa'
	case "hairspace":
		return 0x200a //  '\u200a'
	case "haitusquare":
		return 0x332a // ㌪ '\u332a'
	case "hakatakana":
		return 0x30cf // ハ '\u30cf'
	case "hakatakanahalfwidth":
		return 0xff8a // ﾊ '\uff8a'
	case "halantgurmukhi":
		return 0x0a4d // ੍ '\u0a4d'
	case "hamzadammaarabic":
		return 0x0621 // ء '\u0621'
	case "hamzaisolated":
		return 0xfe80 // ﺀ '\ufe80'
	case "hangulfiller":
		return 0x3164 // ㅤ '\u3164'
	case "hardsigncyrillic":
		return 0x044a // ъ '\u044a'
	case "harpoondownleft":
		return 0x21c3 // ⇃ '\u21c3'
	case "harpoondownright":
		return 0x21c2 // ⇂ '\u21c2'
	case "harpoonleftbarbup":
		return 0x21bc // ↼ '\u21bc'
	case "harpoonleftright":
		return 0x21cc // ⇌ '\u21cc'
	case "harpoonrightbarbup":
		return 0x21c0 // ⇀ '\u21c0'
	case "harpoonrightleft":
		return 0x21cb // ⇋ '\u21cb'
	case "harpoonupleft":
		return 0x21bf // ↿ '\u21bf'
	case "harpoonupright":
		return 0x21be // ↾ '\u21be'
	case "harrowextender":
		return 0x23af // ⎯ '\u23af'
	case "hasquare":
		return 0x33ca // ㏊ '\u33ca'
	case "hatafpatah16":
		return 0x05b2 // ֲ '\u05b2'
	case "hatafqamats28":
		return 0x05b3 // ֳ '\u05b3'
	case "hatafsegolwidehebrew":
		return 0x05b1 // ֱ '\u05b1'
	case "hatapprox":
		return 0x2a6f // ⩯ '\u2a6f'
	case "hbar":
		return 0x0127 // ħ '\u0127'
	case "hbopomofo":
		return 0x310f // ㄏ '\u310f'
	case "hbrevebelow":
		return 0x1e2b // ḫ '\u1e2b'
	case "hcedilla":
		return 0x1e29 // ḩ '\u1e29'
	case "hcircle":
		return 0x24d7 // ⓗ '\u24d7'
	case "hcircumflex":
		return 0x0125 // ĥ '\u0125'
	case "hcyril":
		return 0x03f7 // Ϸ '\u03f7'
	case "hdieresis":
		return 0x1e27 // ḧ '\u1e27'
	case "hdotaccent":
		return 0x1e23 // ḣ '\u1e23'
	case "hdotbelow":
		return 0x1e25 // ḥ '\u1e25'
	case "heart":
		return 0x2665 // ♥ '\u2665'
	case "heartsuitwhite":
		return 0x2661 // ♡ '\u2661'
	case "hedageshhebrew":
		return 0xfb34 // הּ '\ufb34'
	case "heharabic":
		return 0x0647 // ه '\u0647'
	case "hehfinalaltonearabic":
		return 0xfba7 // ﮧ '\ufba7'
	case "hehfinalarabic":
		return 0xfeea // ﻪ '\ufeea'
	case "hehhamzaabovefinalarabic":
		return 0xfba5 // ﮥ '\ufba5'
	case "hehhamzaaboveisolatedarabic":
		return 0xfba4 // ﮤ '\ufba4'
	case "hehinitialaltonearabic":
		return 0xfba8 // ﮨ '\ufba8'
	case "hehinitialarabic":
		return 0xfeeb // ﻫ '\ufeeb'
	case "hehiragana":
		return 0x3078 // へ '\u3078'
	case "hehisolated":
		return 0xfee9 // ﻩ '\ufee9'
	case "hehmedialaltonearabic":
		return 0xfba9 // ﮩ '\ufba9'
	case "hehmedialarabic":
		return 0xfeec // ﻬ '\ufeec'
	case "hehwithmeeminitial":
		return 0xfcd8 // ﳘ '\ufcd8'
	case "heiseierasquare":
		return 0x337b // ㍻ '\u337b'
	case "hekatakana":
		return 0x30d8 // ヘ '\u30d8'
	case "hekatakanahalfwidth":
		return 0xff8d // ﾍ '\uff8d'
	case "hekutaarusquare":
		return 0x3336 // ㌶ '\u3336'
	case "henghook":
		return 0x0267 // ɧ '\u0267'
	case "hermitmatrix":
		return 0x22b9 // ⊹ '\u22b9'
	case "herutusquare":
		return 0x3339 // ㌹ '\u3339'
	case "hexagon":
		return 0x2394 // ⎔ '\u2394'
	case "hexagonblack":
		return 0x2b23 // ⬣ '\u2b23'
	case "hhook":
		return 0x0266 // ɦ '\u0266'
	case "hhooksuper":
		return 0x023a // Ⱥ '\u023a'
	case "hhooksuperior":
		return 0x02b1 // ʱ '\u02b1'
	case "hieuhacirclekorean":
		return 0x327b // ㉻ '\u327b'
	case "hieuhaparenkorean":
		return 0x321b // ㈛ '\u321b'
	case "hieuhcirclekorean":
		return 0x326d // ㉭ '\u326d'
	case "hieuhkorean":
		return 0x314e // ㅎ '\u314e'
	case "hieuhparenkorean":
		return 0x320d // ㈍ '\u320d'
	case "highhamza":
		return 0x0674 // ٴ '\u0674'
	case "hihiragana":
		return 0x3072 // ひ '\u3072'
	case "hikatakana":
		return 0x30d2 // ヒ '\u30d2'
	case "hikatakanahalfwidth":
		return 0xff8b // ﾋ '\uff8b'
	case "hiriq14":
		return 0x05b4 // ִ '\u05b4'
	case "hknearrow":
		return 0x2924 // ⤤ '\u2924'
	case "hknwarrow":
		return 0x2923 // ⤣ '\u2923'
	case "hksearow":
		return 0x2925 // ⤥ '\u2925'
	case "hkswarow":
		return 0x2926 // ⤦ '\u2926'
	case "hlinebelow":
		return 0x1e96 // ẖ '\u1e96'
	case "hmonospace":
		return 0xff48 // ｈ '\uff48'
	case "hoarmenian":
		return 0x0570 // հ '\u0570'
	case "hohipthai":
		return 0x0e2b // ห '\u0e2b'
	case "hohiragana":
		return 0x307b // ほ '\u307b'
	case "hokatakana":
		return 0x30db // ホ '\u30db'
	case "hokatakanahalfwidth":
		return 0xff8e // ﾎ '\uff8e'
	case "holamquarterhebrew":
		return 0x05b9 // ֹ '\u05b9'
	case "honokhukthai":
		return 0x0e2e // ฮ '\u0e2e'
	case "hookcmb":
		return 0x0309 // ̉ '\u0309'
	case "hookpalatalizedbelowcmb":
		return 0x0321 // ̡ '\u0321'
	case "hookretroflexbelowcmb":
		return 0x0322 // ̢ '\u0322'
	case "hoonsquare":
		return 0x3342 // ㍂ '\u3342'
	case "horicoptic":
		return 0x03e9 // ϩ '\u03e9'
	case "horizontalbar":
		return 0x2015 // ― '\u2015'
	case "horiztab":
		return 0x05a2 // ֢ '\u05a2'
	case "horncmb":
		return 0x031b // ̛ '\u031b'
	case "hotsprings":
		return 0x2668 // ♨ '\u2668'
	case "hourglass":
		return 0x29d6 // ⧖ '\u29d6'
	case "house":
		return 0x2302 // ⌂ '\u2302'
	case "hparen":
		return 0x24a3 // ⒣ '\u24a3'
	case "hrectangle":
		return 0x25ad // ▭ '\u25ad'
	case "hsuper":
		return 0x0239 // ȹ '\u0239'
	case "hsuperior":
		return 0x02b0 // ʰ '\u02b0'
	case "hturned":
		return 0x0265 // ɥ '\u0265'
	case "huhiragana":
		return 0x3075 // ふ '\u3075'
	case "huiitosquare":
		return 0x3333 // ㌳ '\u3333'
	case "hukatakana":
		return 0x30d5 // フ '\u30d5'
	case "hukatakanahalfwidth":
		return 0xff8c // ﾌ '\uff8c'
	case "hungarumlaut":
		return 0x02dd // ˝ '\u02dd'
	case "hungarumlaut1":
		return 0xf009 //  '\uf009'
	case "hungarumlautcmb":
		return 0x030b // ̋ '\u030b'
	case "hv":
		return 0x0195 // ƕ '\u0195'
	case "hyphen":
		return 0x002d // - '-'
	case "hyphenbullet":
		return 0x2043 // ⁃ '\u2043'
	case "hyphendot":
		return 0x2027 // ‧ '\u2027'
	case "hypheninferior":
		return 0xf6e5 //  '\uf6e5'
	case "hyphenmonospace":
		return 0xff0d // － '\uff0d'
	case "hyphensmall":
		return 0xfe63 // ﹣ '\ufe63'
	case "hyphensuperior":
		return 0xf6e6 //  '\uf6e6'
	case "hyphentwo":
		return 0x2010 // ‐ '\u2010'
	case "hzigzag":
		return 0x3030 // 〰 '\u3030'
	case "i":
		return 0x0069 // i 'i'
	case "iacute":
		return 0x00ed // í '\u00ed'
	case "ibar":
		return 0x01f8 // Ǹ '\u01f8'
	case "ibengali":
		return 0x0987 // ই '\u0987'
	case "ibopomofo":
		return 0x3127 // ㄧ '\u3127'
	case "ibreve":
		return 0x012d // ĭ '\u012d'
	case "icaron":
		return 0x01d0 // ǐ '\u01d0'
	case "icircle":
		return 0x24d8 // ⓘ '\u24d8'
	case "icircumflex":
		return 0x00ee // î '\u00ee'
	case "idblgrave":
		return 0x0209 // ȉ '\u0209'
	case "ideographearthcircle":
		return 0x328f // ㊏ '\u328f'
	case "ideographfirecircle":
		return 0x328b // ㊋ '\u328b'
	case "ideographicallianceparen":
		return 0x323f // ㈿ '\u323f'
	case "ideographiccallparen":
		return 0x323a // ㈺ '\u323a'
	case "ideographiccentrecircle":
		return 0x32a5 // ㊥ '\u32a5'
	case "ideographicclose":
		return 0x3006 // 〆 '\u3006'
	case "ideographiccomma":
		return 0x3001 // 、 '\u3001'
	case "ideographiccommaleft":
		return 0xff64 // ､ '\uff64'
	case "ideographiccongratulationparen":
		return 0x3237 // ㈷ '\u3237'
	case "ideographiccorrectcircle":
		return 0x32a3 // ㊣ '\u32a3'
	case "ideographicearthparen":
		return 0x322f // ㈯ '\u322f'
	case "ideographicenterpriseparen":
		return 0x323d // ㈽ '\u323d'
	case "ideographicexcellentcircle":
		return 0x329d // ㊝ '\u329d'
	case "ideographicfestivalparen":
		return 0x3240 // ㉀ '\u3240'
	case "ideographicfinancialcircle":
		return 0x3296 // ㊖ '\u3296'
	case "ideographicfinancialparen":
		return 0x3236 // ㈶ '\u3236'
	case "ideographicfireparen":
		return 0x322b // ㈫ '\u322b'
	case "ideographichaveparen":
		return 0x3232 // ㈲ '\u3232'
	case "ideographichighcircle":
		return 0x32a4 // ㊤ '\u32a4'
	case "ideographiciterationmark":
		return 0x3005 // 々 '\u3005'
	case "ideographiclaborcircle":
		return 0x3298 // ㊘ '\u3298'
	case "ideographiclaborparen":
		return 0x3238 // ㈸ '\u3238'
	case "ideographicleftcircle":
		return 0x32a7 // ㊧ '\u32a7'
	case "ideographiclowcircle":
		return 0x32a6 // ㊦ '\u32a6'
	case "ideographicmedicinecircle":
		return 0x32a9 // ㊩ '\u32a9'
	case "ideographicmetalparen":
		return 0x322e // ㈮ '\u322e'
	case "ideographicmoonparen":
		return 0x322a // ㈪ '\u322a'
	case "ideographicnameparen":
		return 0x3234 // ㈴ '\u3234'
	case "ideographicperiod":
		return 0x3002 // 。 '\u3002'
	case "ideographicprintcircle":
		return 0x329e // ㊞ '\u329e'
	case "ideographicreachparen":
		return 0x3243 // ㉃ '\u3243'
	case "ideographicrepresentparen":
		return 0x3239 // ㈹ '\u3239'
	case "ideographicresourceparen":
		return 0x323e // ㈾ '\u323e'
	case "ideographicrightcircle":
		return 0x32a8 // ㊨ '\u32a8'
	case "ideographicsecretcircle":
		return 0x3299 // ㊙ '\u3299'
	case "ideographicselfparen":
		return 0x3242 // ㉂ '\u3242'
	case "ideographicsocietyparen":
		return 0x3233 // ㈳ '\u3233'
	case "ideographicspace":
		return 0x3000 //  '\u3000'
	case "ideographicspecialparen":
		return 0x3235 // ㈵ '\u3235'
	case "ideographicstockparen":
		return 0x3231 // ㈱ '\u3231'
	case "ideographicstudyparen":
		return 0x323b // ㈻ '\u323b'
	case "ideographicsunparen":
		return 0x3230 // ㈰ '\u3230'
	case "ideographicsuperviseparen":
		return 0x323c // ㈼ '\u323c'
	case "ideographicwaterparen":
		return 0x322c // ㈬ '\u322c'
	case "ideographicwoodparen":
		return 0x322d // ㈭ '\u322d'
	case "ideographiczero":
		return 0x3007 // 〇 '\u3007'
	case "ideographmetalcircle":
		return 0x328e // ㊎ '\u328e'
	case "ideographmooncircle":
		return 0x328a // ㊊ '\u328a'
	case "ideographnamecircle":
		return 0x3294 // ㊔ '\u3294'
	case "ideographsuncircle":
		return 0x3290 // ㊐ '\u3290'
	case "ideographwatercircle":
		return 0x328c // ㊌ '\u328c'
	case "ideographwoodcircle":
		return 0x328d // ㊍ '\u328d'
	case "ideva":
		return 0x0907 // इ '\u0907'
	case "idieresis":
		return 0x00ef // ï '\u00ef'
	case "idieresisacute":
		return 0x1e2f // ḯ '\u1e2f'
	case "idieresiscyrillic":
		return 0x04e5 // ӥ '\u04e5'
	case "idotbelow":
		return 0x1ecb // ị '\u1ecb'
	case "iebrevecyrillic":
		return 0x04d7 // ӗ '\u04d7'
	case "iecyrillic":
		return 0x0435 // е '\u0435'
	case "iehook":
		return 0x03f9 // Ϲ '\u03f9'
	case "iehookogonek":
		return 0x03fb // ϻ '\u03fb'
	case "ieungacirclekorean":
		return 0x3275 // ㉵ '\u3275'
	case "ieungaparenkorean":
		return 0x3215 // ㈕ '\u3215'
	case "ieungcirclekorean":
		return 0x3267 // ㉧ '\u3267'
	case "ieungkorean":
		return 0x3147 // ㅇ '\u3147'
	case "ieungparenkorean":
		return 0x3207 // ㈇ '\u3207'
	case "igrave":
		return 0x00ec // ì '\u00ec'
	case "igujarati":
		return 0x0a87 // ઇ '\u0a87'
	case "igurmukhi":
		return 0x0a07 // ਇ '\u0a07'
	case "ihiragana":
		return 0x3044 // い '\u3044'
	case "ihookabove":
		return 0x1ec9 // ỉ '\u1ec9'
	case "iibengali":
		return 0x0988 // ঈ '\u0988'
	case "iicyrillic":
		return 0x0438 // и '\u0438'
	case "iideva":
		return 0x0908 // ई '\u0908'
	case "iigujarati":
		return 0x0a88 // ઈ '\u0a88'
	case "iigurmukhi":
		return 0x0a08 // ਈ '\u0a08'
	case "iiiint":
		return 0x2a0c // ⨌ '\u2a0c'
	case "iiint":
		return 0x222d // ∭ '\u222d'
	case "iimatragurmukhi":
		return 0x0a40 // ੀ '\u0a40'
	case "iinfin":
		return 0x29dc // ⧜ '\u29dc'
	case "iinvertedbreve":
		return 0x020b // ȋ '\u020b'
	case "iivowelsignbengali":
		return 0x09c0 // ী '\u09c0'
	case "iivowelsigndeva":
		return 0x0940 // ी '\u0940'
	case "iivowelsigngujarati":
		return 0x0ac0 // ી '\u0ac0'
	case "ij":
		return 0x0133 // ĳ '\u0133'
	case "ikatakana":
		return 0x30a4 // イ '\u30a4'
	case "ikatakanahalfwidth":
		return 0xff72 // ｲ '\uff72'
	case "ikorean":
		return 0x3163 // ㅣ '\u3163'
	case "iluyhebrew":
		return 0x05ac // ֬ '\u05ac'
	case "imacron":
		return 0x012b // ī '\u012b'
	case "imacroncyrillic":
		return 0x04e3 // ӣ '\u04e3'
	case "imageof":
		return 0x22b7 // ⊷ '\u22b7'
	case "imageorapproximatelyequal":
		return 0x2253 // ≓ '\u2253'
	case "imath":
		return 0x1d6a4 // 𝚤 '\U0001d6a4'
	case "imatragurmukhi":
		return 0x0a3f // ਿ '\u0a3f'
	case "imonospace":
		return 0xff49 // ｉ '\uff49'
	case "infinity":
		return 0x221e // ∞ '\u221e'
	case "iniarmenian":
		return 0x056b // ի '\u056b'
	case "intBar":
		return 0x2a0e // ⨎ '\u2a0e'
	case "intbar":
		return 0x2a0d // ⨍ '\u2a0d'
	case "intcap":
		return 0x2a19 // ⨙ '\u2a19'
	case "intclockwise":
		return 0x2231 // ∱ '\u2231'
	case "intcup":
		return 0x2a1a // ⨚ '\u2a1a'
	case "integerdivide":
		return 0x2216 // ∖ '\u2216'
	case "integral":
		return 0x222b // ∫ '\u222b'
	case "integralbt":
		return 0x2321 // ⌡ '\u2321'
	case "integralex":
		return 0xf8f5 //  '\uf8f5'
	case "integraltp":
		return 0x2320 // ⌠ '\u2320'
	case "intercal":
		return 0x22ba // ⊺ '\u22ba'
	case "interleave":
		return 0x2af4 // ⫴ '\u2af4'
	case "interrobang":
		return 0x203d // ‽ '\u203d'
	case "interrobangdown":
		return 0x2e18 // ⸘ '\u2e18'
	case "intersection":
		return 0x2229 // ∩ '\u2229'
	case "intersectiondbl":
		return 0x22d2 // ⋒ '\u22d2'
	case "intersectiondisplay":
		return 0x22c2 // ⋂ '\u22c2'
	case "intersectionsq":
		return 0x2293 // ⊓ '\u2293'
	case "intextender":
		return 0x23ae // ⎮ '\u23ae'
	case "intisquare":
		return 0x3305 // ㌅ '\u3305'
	case "intlarhk":
		return 0x2a17 // ⨗ '\u2a17'
	case "intprod":
		return 0x2a3c // ⨼ '\u2a3c'
	case "intprodr":
		return 0x2a3d // ⨽ '\u2a3d'
	case "intx":
		return 0x2a18 // ⨘ '\u2a18'
	case "invbullet":
		return 0x25d8 // ◘ '\u25d8'
	case "invcircle":
		return 0x25d9 // ◙ '\u25d9'
	case "invlazys":
		return 0x223e // ∾ '\u223e'
	case "invwhitelowerhalfcircle":
		return 0x25db // ◛ '\u25db'
	case "invwhiteupperhalfcircle":
		return 0x25da // ◚ '\u25da'
	case "iogonek":
		return 0x012f // į '\u012f'
	case "iota":
		return 0x03b9 // ι '\u03b9'
	case "iota1":
		return 0x01f9 // ǹ '\u01f9'
	case "iotadieresis":
		return 0x03ca // ϊ '\u03ca'
	case "iotadieresistonos":
		return 0x0390 // ΐ '\u0390'
	case "iotalatin":
		return 0x0269 // ɩ '\u0269'
	case "iotatonos":
		return 0x03af // ί '\u03af'
	case "iparen":
		return 0x24a4 // ⒤ '\u24a4'
	case "irigurmukhi":
		return 0x0a72 // ੲ '\u0a72'
	case "isinE":
		return 0x22f9 // ⋹ '\u22f9'
	case "isindot":
		return 0x22f5 // ⋵ '\u22f5'
	case "isinobar":
		return 0x22f7 // ⋷ '\u22f7'
	case "isins":
		return 0x22f4 // ⋴ '\u22f4'
	case "isinvb":
		return 0x22f8 // ⋸ '\u22f8'
	case "ismallhiragana":
		return 0x3043 // ぃ '\u3043'
	case "ismallkatakana":
		return 0x30a3 // ィ '\u30a3'
	case "ismallkatakanahalfwidth":
		return 0xff68 // ｨ '\uff68'
	case "issharbengali":
		return 0x09fa // ৺ '\u09fa'
	case "istroke":
		return 0x0268 // ɨ '\u0268'
	case "isuperior":
		return 0xf6ed //  '\uf6ed'
	case "iterationhiragana":
		return 0x309d // ゝ '\u309d'
	case "iterationkatakana":
		return 0x30fd // ヽ '\u30fd'
	case "itilde":
		return 0x0129 // ĩ '\u0129'
	case "itildebelow":
		return 0x1e2d // ḭ '\u1e2d'
	case "iubopomofo":
		return 0x3129 // ㄩ '\u3129'
	case "ivowelsignbengali":
		return 0x09bf // ি '\u09bf'
	case "ivowelsigndeva":
		return 0x093f // ि '\u093f'
	case "ivowelsigngujarati":
		return 0x0abf // િ '\u0abf'
	case "izhitsadblgravecyrillic":
		return 0x0477 // ѷ '\u0477'
	case "j":
		return 0x006a // j 'j'
	case "jaarmenian":
		return 0x0571 // ձ '\u0571'
	case "jabengali":
		return 0x099c // জ '\u099c'
	case "jadeva":
		return 0x091c // ज '\u091c'
	case "jagujarati":
		return 0x0a9c // જ '\u0a9c'
	case "jagurmukhi":
		return 0x0a1c // ਜ '\u0a1c'
	case "jbopomofo":
		return 0x3110 // ㄐ '\u3110'
	case "jcaron":
		return 0x01f0 // ǰ '\u01f0'
	case "jcircle":
		return 0x24d9 // ⓙ '\u24d9'
	case "jcircumflex":
		return 0x0135 // ĵ '\u0135'
	case "jcrossedtail":
		return 0x029d // ʝ '\u029d'
	case "jcrosstail":
		return 0x022d // ȭ '\u022d'
	case "jdotlessstroke":
		return 0x025f // ɟ '\u025f'
	case "jeemarabic":
		return 0x062c // ج '\u062c'
	case "jeemfinalarabic":
		return 0xfe9e // ﺞ '\ufe9e'
	case "jeeminitialarabic":
		return 0xfe9f // ﺟ '\ufe9f'
	case "jeemisolated":
		return 0xfe9d // ﺝ '\ufe9d'
	case "jeemmedialarabic":
		return 0xfea0 // ﺠ '\ufea0'
	case "jeemwithmeeminitial":
		return 0xfca8 // ﲨ '\ufca8'
	case "jehfinalarabic":
		return 0xfb8b // ﮋ '\ufb8b'
	case "jehisolated":
		return 0xfb8a // ﮊ '\ufb8a'
	case "jhabengali":
		return 0x099d // ঝ '\u099d'
	case "jhadeva":
		return 0x091d // झ '\u091d'
	case "jhagujarati":
		return 0x0a9d // ઝ '\u0a9d'
	case "jhagurmukhi":
		return 0x0a1d // ਝ '\u0a1d'
	case "jheharmenian":
		return 0x057b // ջ '\u057b'
	case "jis":
		return 0x3004 // 〄 '\u3004'
	case "jmath":
		return 0x1d6a5 // 𝚥 '\U0001d6a5'
	case "jmonospace":
		return 0xff4a // ｊ '\uff4a'
	case "jparen":
		return 0x24a5 // ⒥ '\u24a5'
	case "jsuper":
		return 0x023b // Ȼ '\u023b'
	case "jsuperior":
		return 0x02b2 // ʲ '\u02b2'
	case "k":
		return 0x006b // k 'k'
	case "kabashkircyrillic":
		return 0x04a1 // ҡ '\u04a1'
	case "kabengali":
		return 0x0995 // ক '\u0995'
	case "kacute":
		return 0x1e31 // ḱ '\u1e31'
	case "kacyrillic":
		return 0x043a // к '\u043a'
	case "kadescendercyrillic":
		return 0x049b // қ '\u049b'
	case "kadeva":
		return 0x0915 // क '\u0915'
	case "kaf":
		return 0x05db // כ '\u05db'
	case "kafdagesh":
		return 0xfb3b // כּ '\ufb3b'
	case "kaffinalarabic":
		return 0xfeda // ﻚ '\ufeda'
	case "kafinitialarabic":
		return 0xfedb // ﻛ '\ufedb'
	case "kafisolated":
		return 0xfed9 // ﻙ '\ufed9'
	case "kafmedialarabic":
		return 0xfedc // ﻜ '\ufedc'
	case "kafrafehebrew":
		return 0xfb4d // כֿ '\ufb4d'
	case "kagujarati":
		return 0x0a95 // ક '\u0a95'
	case "kagurmukhi":
		return 0x0a15 // ਕ '\u0a15'
	case "kahiragana":
		return 0x304b // か '\u304b'
	case "kahook":
		return 0x0400 // Ѐ '\u0400'
	case "kahookcyrillic":
		return 0x04c4 // ӄ '\u04c4'
	case "kakatakana":
		return 0x30ab // カ '\u30ab'
	case "kakatakanahalfwidth":
		return 0xff76 // ｶ '\uff76'
	case "kappa":
		return 0x03ba // κ '\u03ba'
	case "kappasymbolgreek":
		return 0x03f0 // ϰ '\u03f0'
	case "kapyeounmieumkorean":
		return 0x3171 // ㅱ '\u3171'
	case "kapyeounphieuphkorean":
		return 0x3184 // ㆄ '\u3184'
	case "kapyeounpieupkorean":
		return 0x3178 // ㅸ '\u3178'
	case "kapyeounssangpieupkorean":
		return 0x3179 // ㅹ '\u3179'
	case "karoriisquare":
		return 0x330d // ㌍ '\u330d'
	case "kartdes":
		return 0x03d7 // ϗ '\u03d7'
	case "kashidaautonosidebearingarabic":
		return 0x0640 // ـ '\u0640'
	case "kasmallkatakana":
		return 0x30f5 // ヵ '\u30f5'
	case "kasquare":
		return 0x3384 // ㎄ '\u3384'
	case "kasraisolated":
		return 0xfe7a // ﹺ '\ufe7a'
	case "kasralow":
		return 0xe826 //  '\ue826'
	case "kasramedial":
		return 0xfe7b // ﹻ '\ufe7b'
	case "kasratanarabic":
		return 0x064d // ٍ '\u064d'
	case "kasratanisolated":
		return 0xfe74 // ﹴ '\ufe74'
	case "kasratanlow":
		return 0xe827 //  '\ue827'
	case "kastrokecyrillic":
		return 0x049f // ҟ '\u049f'
	case "katahiraprolongmarkhalfwidth":
		return 0xff70 // ｰ '\uff70'
	case "kaverticalstrokecyrillic":
		return 0x049d // ҝ '\u049d'
	case "kbopomofo":
		return 0x310e // ㄎ '\u310e'
	case "kcalsquare":
		return 0x3389 // ㎉ '\u3389'
	case "kcaron":
		return 0x01e9 // ǩ '\u01e9'
	case "kcircle":
		return 0x24da // ⓚ '\u24da'
	case "kcommaaccent":
		return 0x0137 // ķ '\u0137'
	case "kdotbelow":
		return 0x1e33 // ḳ '\u1e33'
	case "keharmenian":
		return 0x0584 // ք '\u0584'
	case "keheh":
		return 0x06a9 // ک '\u06a9'
	case "kehehfinal":
		return 0xfb8f // ﮏ '\ufb8f'
	case "kehehinitial":
		return 0xfb90 // ﮐ '\ufb90'
	case "kehehisolated":
		return 0xfb8e // ﮎ '\ufb8e'
	case "kehehmedial":
		return 0xfb91 // ﮑ '\ufb91'
	case "kehiragana":
		return 0x3051 // け '\u3051'
	case "kekatakana":
		return 0x30b1 // ケ '\u30b1'
	case "kekatakanahalfwidth":
		return 0xff79 // ｹ '\uff79'
	case "kenarmenian":
		return 0x056f // կ '\u056f'
	case "kernelcontraction":
		return 0x223b // ∻ '\u223b'
	case "kesmallkatakana":
		return 0x30f6 // ヶ '\u30f6'
	case "kgreenlandic":
		return 0x0138 // ĸ '\u0138'
	case "khabengali":
		return 0x0996 // খ '\u0996'
	case "khadeva":
		return 0x0916 // ख '\u0916'
	case "khagujarati":
		return 0x0a96 // ખ '\u0a96'
	case "khagurmukhi":
		return 0x0a16 // ਖ '\u0a16'
	case "khahfinalarabic":
		return 0xfea6 // ﺦ '\ufea6'
	case "khahinitialarabic":
		return 0xfea7 // ﺧ '\ufea7'
	case "khahisolated":
		return 0xfea5 // ﺥ '\ufea5'
	case "khahmedialarabic":
		return 0xfea8 // ﺨ '\ufea8'
	case "khahwithmeeminitial":
		return 0xfcac // ﲬ '\ufcac'
	case "kheicoptic":
		return 0x03e7 // ϧ '\u03e7'
	case "khhadeva":
		return 0x0959 // ख़ '\u0959'
	case "khhagurmukhi":
		return 0x0a59 // ਖ਼ '\u0a59'
	case "khieukhacirclekorean":
		return 0x3278 // ㉸ '\u3278'
	case "khieukhaparenkorean":
		return 0x3218 // ㈘ '\u3218'
	case "khieukhcirclekorean":
		return 0x326a // ㉪ '\u326a'
	case "khieukhkorean":
		return 0x314b // ㅋ '\u314b'
	case "khieukhparenkorean":
		return 0x320a // ㈊ '\u320a'
	case "khokhaithai":
		return 0x0e02 // ข '\u0e02'
	case "khokhonthai":
		return 0x0e05 // ฅ '\u0e05'
	case "khokhuatthai":
		return 0x0e03 // ฃ '\u0e03'
	case "khokhwaithai":
		return 0x0e04 // ค '\u0e04'
	case "khomutthai":
		return 0x0e5b // ๛ '\u0e5b'
	case "khook":
		return 0x0199 // ƙ '\u0199'
	case "khorakhangthai":
		return 0x0e06 // ฆ '\u0e06'
	case "khzsquare":
		return 0x3391 // ㎑ '\u3391'
	case "kihiragana":
		return 0x304d // き '\u304d'
	case "kikatakana":
		return 0x30ad // キ '\u30ad'
	case "kikatakanahalfwidth":
		return 0xff77 // ｷ '\uff77'
	case "kiroguramusquare":
		return 0x3315 // ㌕ '\u3315'
	case "kiromeetorusquare":
		return 0x3316 // ㌖ '\u3316'
	case "kirosquare":
		return 0x3314 // ㌔ '\u3314'
	case "kiyeokacirclekorean":
		return 0x326e // ㉮ '\u326e'
	case "kiyeokaparenkorean":
		return 0x320e // ㈎ '\u320e'
	case "kiyeokcirclekorean":
		return 0x3260 // ㉠ '\u3260'
	case "kiyeokkorean":
		return 0x3131 // ㄱ '\u3131'
	case "kiyeokparenkorean":
		return 0x3200 // ㈀ '\u3200'
	case "kiyeoksioskorean":
		return 0x3133 // ㄳ '\u3133'
	case "klinebelow":
		return 0x1e35 // ḵ '\u1e35'
	case "klsquare":
		return 0x3398 // ㎘ '\u3398'
	case "kmcubedsquare":
		return 0x33a6 // ㎦ '\u33a6'
	case "kmonospace":
		return 0xff4b // ｋ '\uff4b'
	case "kmsquaredsquare":
		return 0x33a2 // ㎢ '\u33a2'
	case "kohiragana":
		return 0x3053 // こ '\u3053'
	case "kohmsquare":
		return 0x33c0 // ㏀ '\u33c0'
	case "kokaithai":
		return 0x0e01 // ก '\u0e01'
	case "kokatakana":
		return 0x30b3 // コ '\u30b3'
	case "kokatakanahalfwidth":
		return 0xff7a // ｺ '\uff7a'
	case "kooposquare":
		return 0x331e // ㌞ '\u331e'
	case "koppacyrillic":
		return 0x0481 // ҁ '\u0481'
	case "koreanstandardsymbol":
		return 0x327f // ㉿ '\u327f'
	case "koroniscmb":
		return 0x0343 // ̓ '\u0343'
	case "kparen":
		return 0x24a6 // ⒦ '\u24a6'
	case "kpasquare":
		return 0x33aa // ㎪ '\u33aa'
	case "ksicyrillic":
		return 0x046f // ѯ '\u046f'
	case "ktsquare":
		return 0x33cf // ㏏ '\u33cf'
	case "kturn":
		return 0x022e // Ȯ '\u022e'
	case "kturned":
		return 0x029e // ʞ '\u029e'
	case "kuhiragana":
		return 0x304f // く '\u304f'
	case "kukatakana":
		return 0x30af // ク '\u30af'
	case "kukatakanahalfwidth":
		return 0xff78 // ｸ '\uff78'
	case "kvsquare":
		return 0x33b8 // ㎸ '\u33b8'
	case "kwsquare":
		return 0x33be // ㎾ '\u33be'
	case "l":
		return 0x006c // l 'l'
	case "lAngle":
		return 0x27ea // ⟪ '\u27ea'
	case "lBrace":
		return 0x2983 // ⦃ '\u2983'
	case "lParen":
		return 0x2985 // ⦅ '\u2985'
	case "labengali":
		return 0x09b2 // ল '\u09b2'
	case "lacute":
		return 0x013a // ĺ '\u013a'
	case "ladeva":
		return 0x0932 // ल '\u0932'
	case "lagujarati":
		return 0x0ab2 // લ '\u0ab2'
	case "lagurmukhi":
		return 0x0a32 // ਲ '\u0a32'
	case "lakkhangyaothai":
		return 0x0e45 // ๅ '\u0e45'
	case "lamaleffinalarabic":
		return 0xfefc // ﻼ '\ufefc'
	case "lamalefhamzaabovefinalarabic":
		return 0xfef8 // ﻸ '\ufef8'
	case "lamalefhamzaaboveisolatedarabic":
		return 0xfef7 // ﻷ '\ufef7'
	case "lamalefhamzabelowfinalarabic":
		return 0xfefa // ﻺ '\ufefa'
	case "lamalefhamzabelowisolatedarabic":
		return 0xfef9 // ﻹ '\ufef9'
	case "lamalefisolatedarabic":
		return 0xfefb // ﻻ '\ufefb'
	case "lamalefmaddaabovefinalarabic":
		return 0xfef6 // ﻶ '\ufef6'
	case "lamalefmaddaaboveisolatedarabic":
		return 0xfef5 // ﻵ '\ufef5'
	case "lambda":
		return 0x03bb // λ '\u03bb'
	case "lambdastroke":
		return 0x019b // ƛ '\u019b'
	case "lameddagesh":
		return 0xfb3c // לּ '\ufb3c'
	case "lamedholamhebrew":
		return 0x05dc // ל '\u05dc'
	case "lamedwithdageshandholam":
		return 0xe805 //  '\ue805'
	case "lamedwithholam":
		return 0xe804 //  '\ue804'
	case "lamfinalarabic":
		return 0xfede // ﻞ '\ufede'
	case "lamhahinitialarabic":
		return 0xfcca // ﳊ '\ufcca'
	case "laminitialarabic":
		return 0xfedf // ﻟ '\ufedf'
	case "lamisolated":
		return 0xfedd // ﻝ '\ufedd'
	case "lamjeeminitialarabic":
		return 0xfcc9 // ﳉ '\ufcc9'
	case "lamkhahinitialarabic":
		return 0xfccb // ﳋ '\ufccb'
	case "lamlamhehisolatedarabic":
		return 0xfdf2 // ﷲ '\ufdf2'
	case "lammedialarabic":
		return 0xfee0 // ﻠ '\ufee0'
	case "lammeemhahinitialarabic":
		return 0xfd88 // ﶈ '\ufd88'
	case "lammeeminitialarabic":
		return 0xfccc // ﳌ '\ufccc'
	case "lamwithalefmaksuraisolated":
		return 0xfc43 // ﱃ '\ufc43'
	case "lamwithhahisolated":
		return 0xfc40 // ﱀ '\ufc40'
	case "lamwithhehinitial":
		return 0xfccd // ﳍ '\ufccd'
	case "lamwithjeemisolated":
		return 0xfc3f // ﰿ '\ufc3f'
	case "lamwithkhahisolated":
		return 0xfc41 // ﱁ '\ufc41'
	case "lamwithmeemisolated":
		return 0xfc42 // ﱂ '\ufc42'
	case "lamwithmeemwithjeeminitial":
		return 0xe811 //  '\ue811'
	case "lamwithyehisolated":
		return 0xfc44 // ﱄ '\ufc44'
	case "langledot":
		return 0x2991 // ⦑ '\u2991'
	case "laplac":
		return 0x29e0 // ⧠ '\u29e0'
	case "largecircle":
		return 0x25ef // ◯ '\u25ef'
	case "lat":
		return 0x2aab // ⪫ '\u2aab'
	case "late":
		return 0x2aad // ⪭ '\u2aad'
	case "lbag":
		return 0x27c5 // ⟅ '\u27c5'
	case "lbar":
		return 0x019a // ƚ '\u019a'
	case "lbbar":
		return 0x2114 // ℔ '\u2114'
	case "lbelt":
		return 0x026c // ɬ '\u026c'
	case "lblkbrbrak":
		return 0x2997 // ⦗ '\u2997'
	case "lbopomofo":
		return 0x310c // ㄌ '\u310c'
	case "lbracelend":
		return 0x23a9 // ⎩ '\u23a9'
	case "lbracemid":
		return 0x23a8 // ⎨ '\u23a8'
	case "lbraceuend":
		return 0x23a7 // ⎧ '\u23a7'
	case "lbrackextender":
		return 0x23a2 // ⎢ '\u23a2'
	case "lbracklend":
		return 0x23a3 // ⎣ '\u23a3'
	case "lbracklltick":
		return 0x298f // ⦏ '\u298f'
	case "lbrackubar":
		return 0x298b // ⦋ '\u298b'
	case "lbrackuend":
		return 0x23a1 // ⎡ '\u23a1'
	case "lbrackultick":
		return 0x298d // ⦍ '\u298d'
	case "lbrbrak":
		return 0x2772 // ❲ '\u2772'
	case "lcaron":
		return 0x013e // ľ '\u013e'
	case "lcaron1":
		return 0xf813 //  '\uf813'
	case "lcircle":
		return 0x24db // ⓛ '\u24db'
	case "lcircumflexbelow":
		return 0x1e3d // ḽ '\u1e3d'
	case "lcommaaccent":
		return 0x013c // ļ '\u013c'
	case "lcurvyangle":
		return 0x29fc // ⧼ '\u29fc'
	case "ldotaccent":
		return 0x0140 // ŀ '\u0140'
	case "ldotbelow":
		return 0x1e37 // ḷ '\u1e37'
	case "ldotbelowmacron":
		return 0x1e39 // ḹ '\u1e39'
	case "leftangleabovecmb":
		return 0x031a // ̚ '\u031a'
	case "leftarrowapprox":
		return 0x2b4a // ⭊ '\u2b4a'
	case "leftarrowbackapprox":
		return 0x2b42 // ⭂ '\u2b42'
	case "leftarrowbsimilar":
		return 0x2b4b // ⭋ '\u2b4b'
	case "leftarrowless":
		return 0x2977 // ⥷ '\u2977'
	case "leftarrowonoplus":
		return 0x2b32 // ⬲ '\u2b32'
	case "leftarrowplus":
		return 0x2946 // ⥆ '\u2946'
	case "leftarrowshortrightarrow":
		return 0x2943 // ⥃ '\u2943'
	case "leftarrowsimilar":
		return 0x2973 // ⥳ '\u2973'
	case "leftarrowsubset":
		return 0x297a // ⥺ '\u297a'
	case "leftarrowtriangle":
		return 0x21fd // ⇽ '\u21fd'
	case "leftarrowx":
		return 0x2b3e // ⬾ '\u2b3e'
	case "leftbkarrow":
		return 0x290c // ⤌ '\u290c'
	case "leftcurvedarrow":
		return 0x2b3f // ⬿ '\u2b3f'
	case "leftdbkarrow":
		return 0x290e // ⤎ '\u290e'
	case "leftdbltail":
		return 0x291b // ⤛ '\u291b'
	case "leftdotarrow":
		return 0x2b38 // ⬸ '\u2b38'
	case "leftdowncurvedarrow":
		return 0x2936 // ⤶ '\u2936'
	case "leftfishtail":
		return 0x297c // ⥼ '\u297c'
	case "leftharpoonaccent":
		return 0x20d0 // ⃐ '\u20d0'
	case "leftharpoondownbar":
		return 0x295e // ⥞ '\u295e'
	case "leftharpoonsupdown":
		return 0x2962 // ⥢ '\u2962'
	case "leftharpoonupbar":
		return 0x295a // ⥚ '\u295a'
	case "leftharpoonupdash":
		return 0x296a // ⥪ '\u296a'
	case "leftleftarrows":
		return 0x21c7 // ⇇ '\u21c7'
	case "leftmoon":
		return 0x263e // ☾ '\u263e'
	case "leftouterjoin":
		return 0x27d5 // ⟕ '\u27d5'
	case "leftrightarrowcircle":
		return 0x2948 // ⥈ '\u2948'
	case "leftrightarrowtriangle":
		return 0x21ff // ⇿ '\u21ff'
	case "leftrightharpoondowndown":
		return 0x2950 // ⥐ '\u2950'
	case "leftrightharpoondownup":
		return 0x294b // ⥋ '\u294b'
	case "leftrightharpoonsdown":
		return 0x2967 // ⥧ '\u2967'
	case "leftrightharpoonsup":
		return 0x2966 // ⥦ '\u2966'
	case "leftrightharpoonupdown":
		return 0x294a // ⥊ '\u294a'
	case "leftrightharpoonupup":
		return 0x294e // ⥎ '\u294e'
	case "leftsquigarrow":
		return 0x21dc // ⇜ '\u21dc'
	case "lefttackbelowcmb":
		return 0x0318 // ̘ '\u0318'
	case "lefttail":
		return 0x2919 // ⤙ '\u2919'
	case "leftthreearrows":
		return 0x2b31 // ⬱ '\u2b31'
	case "leftwavearrow":
		return 0x219c // ↜ '\u219c'
	case "leqqslant":
		return 0x2af9 // ⫹ '\u2af9'
	case "lescc":
		return 0x2aa8 // ⪨ '\u2aa8'
	case "lesdot":
		return 0x2a7f // ⩿ '\u2a7f'
	case "lesdoto":
		return 0x2a81 // ⪁ '\u2a81'
	case "lesdotor":
		return 0x2a83 // ⪃ '\u2a83'
	case "lesges":
		return 0x2a93 // ⪓ '\u2a93'
	case "less":
		return 0x003c // < '<'
	case "lessdbleqlgreater":
		return 0x2a8b // ⪋ '\u2a8b'
	case "lessdot":
		return 0x22d6 // ⋖ '\u22d6'
	case "lessequal":
		return 0x2264 // ≤ '\u2264'
	case "lessequalorgreater":
		return 0x22da // ⋚ '\u22da'
	case "lessmonospace":
		return 0xff1c // ＜ '\uff1c'
	case "lessnotdblequal":
		return 0x2a89 // ⪉ '\u2a89'
	case "lessnotequal":
		return 0x2a87 // ⪇ '\u2a87'
	case "lessorapproxeql":
		return 0x2a85 // ⪅ '\u2a85'
	case "lessorequalslant":
		return 0x2a7d // ⩽ '\u2a7d'
	case "lessorequivalent":
		return 0x2272 // ≲ '\u2272'
	case "lessorgreater":
		return 0x2276 // ≶ '\u2276'
	case "lessornotequal":
		return 0x2268 // ≨ '\u2268'
	case "lessoverequal":
		return 0x2266 // ≦ '\u2266'
	case "lesssmall":
		return 0xfe64 // ﹤ '\ufe64'
	case "lezh":
		return 0x026e // ɮ '\u026e'
	case "lfblock":
		return 0x258c // ▌ '\u258c'
	case "lfbowtie":
		return 0x29d1 // ⧑ '\u29d1'
	case "lfeighthblock":
		return 0x258f // ▏ '\u258f'
	case "lffiveeighthblock":
		return 0x258b // ▋ '\u258b'
	case "lfquarterblock":
		return 0x258e // ▎ '\u258e'
	case "lfseveneighthblock":
		return 0x2589 // ▉ '\u2589'
	case "lfthreeeighthblock":
		return 0x258d // ▍ '\u258d'
	case "lfthreequarterblock":
		return 0x258a // ▊ '\u258a'
	case "lftimes":
		return 0x29d4 // ⧔ '\u29d4'
	case "lgE":
		return 0x2a91 // ⪑ '\u2a91'
	case "lgblkcircle":
		return 0x2b24 // ⬤ '\u2b24'
	case "lgblksquare":
		return 0x2b1b // ⬛ '\u2b1b'
	case "lgwhtsquare":
		return 0x2b1c // ⬜ '\u2b1c'
	case "lhookretroflex":
		return 0x026d // ɭ '\u026d'
	case "linefeed":
		return 0x21b4 // ↴ '\u21b4'
	case "lineseparator":
		return 0x2028 //  '\u2028'
	case "linevertnosp":
		return 0x0280 // ʀ '\u0280'
	case "linevertsubnosp":
		return 0x029c // ʜ '\u029c'
	case "lira":
		return 0x20a4 // ₤ '\u20a4'
	case "liwnarmenian":
		return 0x056c // լ '\u056c'
	case "lj":
		return 0x01c9 // ǉ '\u01c9'
	case "ljecyrillic":
		return 0x0459 // љ '\u0459'
	case "ll":
		return 0xf6c0 //  '\uf6c0'
	case "lladeva":
		return 0x0933 // ळ '\u0933'
	case "llagujarati":
		return 0x0ab3 // ળ '\u0ab3'
	case "llangle":
		return 0x2989 // ⦉ '\u2989'
	case "llarc":
		return 0x25df // ◟ '\u25df'
	case "llinebelow":
		return 0x1e3b // ḻ '\u1e3b'
	case "lll":
		return 0x22d8 // ⋘ '\u22d8'
	case "llladeva":
		return 0x0934 // ऴ '\u0934'
	case "lllnest":
		return 0x2af7 // ⫷ '\u2af7'
	case "llparenthesis":
		return 0x2987 // ⦇ '\u2987'
	case "lltriangle":
		return 0x25fa // ◺ '\u25fa'
	case "llvocalicbengali":
		return 0x09e1 // ৡ '\u09e1'
	case "llvocalicdeva":
		return 0x0961 // ॡ '\u0961'
	case "llvocalicvowelsignbengali":
		return 0x09e3 // ৣ '\u09e3'
	case "llvocalicvowelsigndeva":
		return 0x0963 // ॣ '\u0963'
	case "lmiddletilde":
		return 0x026b // ɫ '\u026b'
	case "lmonospace":
		return 0xff4c // ｌ '\uff4c'
	case "lmoustache":
		return 0x23b0 // ⎰ '\u23b0'
	case "lmsquare":
		return 0x33d0 // ㏐ '\u33d0'
	case "lnsim":
		return 0x22e6 // ⋦ '\u22e6'
	case "lochulathai":
		return 0x0e2c // ฬ '\u0e2c'
	case "logicaland":
		return 0x2227 // ∧ '\u2227'
	case "logicalnot":
		return 0x00ac // ¬ '\u00ac'
	case "logicalor":
		return 0x2228 // ∨ '\u2228'
	case "logonek":
		return 0xf830 //  '\uf830'
	case "lolingthai":
		return 0x0e25 // ล '\u0e25'
	case "longdashv":
		return 0x27de // ⟞ '\u27de'
	case "longdivision":
		return 0x27cc // ⟌ '\u27cc'
	case "longleftarrow":
		return 0x27f5 // ⟵ '\u27f5'
	case "longleftrightarrow":
		return 0x27f7 // ⟷ '\u27f7'
	case "longleftsquigarrow":
		return 0x2b33 // ⬳ '\u2b33'
	case "longmapsfrom":
		return 0x27fb // ⟻ '\u27fb'
	case "longmapsto":
		return 0x27fc // ⟼ '\u27fc'
	case "longrightarrow":
		return 0x27f6 // ⟶ '\u27f6'
	case "longrightsquigarrow":
		return 0x27ff // ⟿ '\u27ff'
	case "longs":
		return 0x017f // ſ '\u017f'
	case "longst":
		return 0xfb05 // ﬅ '\ufb05'
	case "lowered":
		return 0x024e // Ɏ '\u024e'
	case "lowint":
		return 0x2a1c // ⨜ '\u2a1c'
	case "lowlinecenterline":
		return 0xfe4e // ﹎ '\ufe4e'
	case "lowlinecmb":
		return 0x0332 // ̲ '\u0332'
	case "lowlinedashed":
		return 0xfe4d // ﹍ '\ufe4d'
	case "lozenge":
		return 0x25ca // ◊ '\u25ca'
	case "lozengeminus":
		return 0x27e0 // ⟠ '\u27e0'
	case "lparen":
		return 0x24a7 // ⒧ '\u24a7'
	case "lparenextender":
		return 0x239c // ⎜ '\u239c'
	case "lparenlend":
		return 0x239d // ⎝ '\u239d'
	case "lparenless":
		return 0x2993 // ⦓ '\u2993'
	case "lparenuend":
		return 0x239b // ⎛ '\u239b'
	case "lrarc":
		return 0x25de // ◞ '\u25de'
	case "lre":
		return 0x202a //  '\u202a'
	case "lrtriangle":
		return 0x25ff // ◿ '\u25ff'
	case "lrtriangleeq":
		return 0x29e1 // ⧡ '\u29e1'
	case "lsime":
		return 0x2a8d // ⪍ '\u2a8d'
	case "lsimg":
		return 0x2a8f // ⪏ '\u2a8f'
	case "lslash":
		return 0x0142 // ł '\u0142'
	case "lsqhook":
		return 0x2acd // ⫍ '\u2acd'
	case "lsuper":
		return 0x026a // ɪ '\u026a'
	case "lsuperior":
		return 0xf6ee //  '\uf6ee'
	case "ltcc":
		return 0x2aa6 // ⪦ '\u2aa6'
	case "ltcir":
		return 0x2a79 // ⩹ '\u2a79'
	case "ltlarr":
		return 0x2976 // ⥶ '\u2976'
	case "ltquest":
		return 0x2a7b // ⩻ '\u2a7b'
	case "ltrivb":
		return 0x29cf // ⧏ '\u29cf'
	case "ltshade1":
		return 0xf821 //  '\uf821'
	case "luthai":
		return 0x0e26 // ฦ '\u0e26'
	case "lvboxline":
		return 0x23b8 // ⎸ '\u23b8'
	case "lvocalicbengali":
		return 0x098c // ঌ '\u098c'
	case "lvocalicdeva":
		return 0x090c // ऌ '\u090c'
	case "lvocalicvowelsignbengali":
		return 0x09e2 // ৢ '\u09e2'
	case "lvocalicvowelsigndeva":
		return 0x0962 // ॢ '\u0962'
	case "lvzigzag":
		return 0x29d8 // ⧘ '\u29d8'
	case "lxsquare":
		return 0x33d3 // ㏓ '\u33d3'
	case "m":
		return 0x006d // m 'm'
	case "mabengali":
		return 0x09ae // ম '\u09ae'
	case "macron":
		return 0x00af // ¯ '\u00af'
	case "macronbelowcmb":
		return 0x0331 // ̱ '\u0331'
	case "macroncmb":
		return 0x0304 // ̄ '\u0304'
	case "macronlowmod":
		return 0x02cd // ˍ '\u02cd'
	case "macronmonospace":
		return 0xffe3 // ￣ '\uffe3'
	case "macute":
		return 0x1e3f // ḿ '\u1e3f'
	case "madeva":
		return 0x092e // म '\u092e'
	case "magujarati":
		return 0x0aae // મ '\u0aae'
	case "magurmukhi":
		return 0x0a2e // ਮ '\u0a2e'
	case "mahapakhhebrew":
		return 0x05a4 // ֤ '\u05a4'
	case "mahiragana":
		return 0x307e // ま '\u307e'
	case "maichattawalowleftthai":
		return 0xf895 //  '\uf895'
	case "maichattawalowrightthai":
		return 0xf894 //  '\uf894'
	case "maichattawathai":
		return 0x0e4b // ๋ '\u0e4b'
	case "maichattawaupperleftthai":
		return 0xf893 //  '\uf893'
	case "maieklowleftthai":
		return 0xf88c //  '\uf88c'
	case "maieklowrightthai":
		return 0xf88b //  '\uf88b'
	case "maiekthai":
		return 0x0e48 // ่ '\u0e48'
	case "maiekupperleftthai":
		return 0xf88a //  '\uf88a'
	case "maihanakatleftthai":
		return 0xf884 //  '\uf884'
	case "maihanakatthai":
		return 0x0e31 // ั '\u0e31'
	case "maitaikhuleftthai":
		return 0xf889 //  '\uf889'
	case "maitaikhuthai":
		return 0x0e47 // ็ '\u0e47'
	case "maitholowleftthai":
		return 0xf88f //  '\uf88f'
	case "maitholowrightthai":
		return 0xf88e //  '\uf88e'
	case "maithothai":
		return 0x0e49 // ้ '\u0e49'
	case "maithoupperleftthai":
		return 0xf88d //  '\uf88d'
	case "maitrilowleftthai":
		return 0xf892 //  '\uf892'
	case "maitrilowrightthai":
		return 0xf891 //  '\uf891'
	case "maitrithai":
		return 0x0e4a // ๊ '\u0e4a'
	case "maitriupperleftthai":
		return 0xf890 //  '\uf890'
	case "maiyamokthai":
		return 0x0e46 // ๆ '\u0e46'
	case "makatakana":
		return 0x30de // マ '\u30de'
	case "makatakanahalfwidth":
		return 0xff8f // ﾏ '\uff8f'
	case "mansyonsquare":
		return 0x3347 // ㍇ '\u3347'
	case "mapsdown":
		return 0x21a7 // ↧ '\u21a7'
	case "mapsfrom":
		return 0x21a4 // ↤ '\u21a4'
	case "mapsto":
		return 0x21a6 // ↦ '\u21a6'
	case "mapsup":
		return 0x21a5 // ↥ '\u21a5'
	case "mars":
		return 0x2642 // ♂ '\u2642'
	case "masoracirclehebrew":
		return 0x05af // ֯ '\u05af'
	case "masquare":
		return 0x3383 // ㎃ '\u3383'
	case "mbfA":
		return 0x1d400 // 𝐀 '\U0001d400'
	case "mbfAlpha":
		return 0x1d6a8 // 𝚨 '\U0001d6a8'
	case "mbfB":
		return 0x1d401 // 𝐁 '\U0001d401'
	case "mbfBeta":
		return 0x1d6a9 // 𝚩 '\U0001d6a9'
	case "mbfC":
		return 0x1d402 // 𝐂 '\U0001d402'
	case "mbfChi":
		return 0x1d6be // 𝚾 '\U0001d6be'
	case "mbfD":
		return 0x1d403 // 𝐃 '\U0001d403'
	case "mbfDelta":
		return 0x1d6ab // 𝚫 '\U0001d6ab'
	case "mbfDigamma":
		return 0x1d7ca // 𝟊 '\U0001d7ca'
	case "mbfE":
		return 0x1d404 // 𝐄 '\U0001d404'
	case "mbfEpsilon":
		return 0x1d6ac // 𝚬 '\U0001d6ac'
	case "mbfEta":
		return 0x1d6ae // 𝚮 '\U0001d6ae'
	case "mbfF":
		return 0x1d405 // 𝐅 '\U0001d405'
	case "mbfG":
		return 0x1d406 // 𝐆 '\U0001d406'
	case "mbfGamma":
		return 0x1d6aa // 𝚪 '\U0001d6aa'
	case "mbfH":
		return 0x1d407 // 𝐇 '\U0001d407'
	case "mbfI":
		return 0x1d408 // 𝐈 '\U0001d408'
	case "mbfIota":
		return 0x1d6b0 // 𝚰 '\U0001d6b0'
	case "mbfJ":
		return 0x1d409 // 𝐉 '\U0001d409'
	case "mbfK":
		return 0x1d40a // 𝐊 '\U0001d40a'
	case "mbfKappa":
		return 0x1d6b1 // 𝚱 '\U0001d6b1'
	case "mbfL":
		return 0x1d40b // 𝐋 '\U0001d40b'
	case "mbfLambda":
		return 0x1d6b2 // 𝚲 '\U0001d6b2'
	case "mbfM":
		return 0x1d40c // 𝐌 '\U0001d40c'
	case "mbfMu":
		return 0x1d6b3 // 𝚳 '\U0001d6b3'
	case "mbfN":
		return 0x1d40d // 𝐍 '\U0001d40d'
	case "mbfNu":
		return 0x1d6b4 // 𝚴 '\U0001d6b4'
	case "mbfO":
		return 0x1d40e // 𝐎 '\U0001d40e'
	case "mbfOmega":
		return 0x1d6c0 // 𝛀 '\U0001d6c0'
	case "mbfOmicron":
		return 0x1d6b6 // 𝚶 '\U0001d6b6'
	case "mbfP":
		return 0x1d40f // 𝐏 '\U0001d40f'
	case "mbfPhi":
		return 0x1d6bd // 𝚽 '\U0001d6bd'
	case "mbfPi":
		return 0x1d6b7 // 𝚷 '\U0001d6b7'
	case "mbfPsi":
		return 0x1d6bf // 𝚿 '\U0001d6bf'
	case "mbfQ":
		return 0x1d410 // 𝐐 '\U0001d410'
	case "mbfR":
		return 0x1d411 // 𝐑 '\U0001d411'
	case "mbfRho":
		return 0x1d6b8 // 𝚸 '\U0001d6b8'
	case "mbfS":
		return 0x1d412 // 𝐒 '\U0001d412'
	case "mbfSigma":
		return 0x1d6ba // 𝚺 '\U0001d6ba'
	case "mbfT":
		return 0x1d413 // 𝐓 '\U0001d413'
	case "mbfTau":
		return 0x1d6bb // 𝚻 '\U0001d6bb'
	case "mbfTheta":
		return 0x1d6af // 𝚯 '\U0001d6af'
	case "mbfU":
		return 0x1d414 // 𝐔 '\U0001d414'
	case "mbfUpsilon":
		return 0x1d6bc // 𝚼 '\U0001d6bc'
	case "mbfV":
		return 0x1d415 // 𝐕 '\U0001d415'
	case "mbfW":
		return 0x1d416 // 𝐖 '\U0001d416'
	case "mbfX":
		return 0x1d417 // 𝐗 '\U0001d417'
	case "mbfXi":
		return 0x1d6b5 // 𝚵 '\U0001d6b5'
	case "mbfY":
		return 0x1d418 // 𝐘 '\U0001d418'
	case "mbfZ":
		return 0x1d419 // 𝐙 '\U0001d419'
	case "mbfZeta":
		return 0x1d6ad // 𝚭 '\U0001d6ad'
	case "mbfa":
		return 0x1d41a // 𝐚 '\U0001d41a'
	case "mbfalpha":
		return 0x1d6c2 // 𝛂 '\U0001d6c2'
	case "mbfb":
		return 0x1d41b // 𝐛 '\U0001d41b'
	case "mbfbeta":
		return 0x1d6c3 // 𝛃 '\U0001d6c3'
	case "mbfc":
		return 0x1d41c // 𝐜 '\U0001d41c'
	case "mbfchi":
		return 0x1d6d8 // 𝛘 '\U0001d6d8'
	case "mbfd":
		return 0x1d41d // 𝐝 '\U0001d41d'
	case "mbfdelta":
		return 0x1d6c5 // 𝛅 '\U0001d6c5'
	case "mbfe":
		return 0x1d41e // 𝐞 '\U0001d41e'
	case "mbfepsilon":
		return 0x1d6c6 // 𝛆 '\U0001d6c6'
	case "mbfeta":
		return 0x1d6c8 // 𝛈 '\U0001d6c8'
	case "mbff":
		return 0x1d41f // 𝐟 '\U0001d41f'
	case "mbffrakA":
		return 0x1d56c // 𝕬 '\U0001d56c'
	case "mbffrakB":
		return 0x1d56d // 𝕭 '\U0001d56d'
	case "mbffrakC":
		return 0x1d56e // 𝕮 '\U0001d56e'
	case "mbffrakD":
		return 0x1d56f // 𝕯 '\U0001d56f'
	case "mbffrakE":
		return 0x1d570 // 𝕰 '\U0001d570'
	case "mbffrakF":
		return 0x1d571 // 𝕱 '\U0001d571'
	case "mbffrakG":
		return 0x1d572 // 𝕲 '\U0001d572'
	case "mbffrakH":
		return 0x1d573 // 𝕳 '\U0001d573'
	case "mbffrakI":
		return 0x1d574 // 𝕴 '\U0001d574'
	case "mbffrakJ":
		return 0x1d575 // 𝕵 '\U0001d575'
	case "mbffrakK":
		return 0x1d576 // 𝕶 '\U0001d576'
	case "mbffrakL":
		return 0x1d577 // 𝕷 '\U0001d577'
	case "mbffrakM":
		return 0x1d578 // 𝕸 '\U0001d578'
	case "mbffrakN":
		return 0x1d579 // 𝕹 '\U0001d579'
	case "mbffrakO":
		return 0x1d57a // 𝕺 '\U0001d57a'
	case "mbffrakP":
		return 0x1d57b // 𝕻 '\U0001d57b'
	case "mbffrakQ":
		return 0x1d57c // 𝕼 '\U0001d57c'
	case "mbffrakR":
		return 0x1d57d // 𝕽 '\U0001d57d'
	case "mbffrakS":
		return 0x1d57e // 𝕾 '\U0001d57e'
	case "mbffrakT":
		return 0x1d57f // 𝕿 '\U0001d57f'
	case "mbffrakU":
		return 0x1d580 // 𝖀 '\U0001d580'
	case "mbffrakV":
		return 0x1d581 // 𝖁 '\U0001d581'
	case "mbffrakW":
		return 0x1d582 // 𝖂 '\U0001d582'
	case "mbffrakX":
		return 0x1d583 // 𝖃 '\U0001d583'
	case "mbffrakY":
		return 0x1d584 // 𝖄 '\U0001d584'
	case "mbffrakZ":
		return 0x1d585 // 𝖅 '\U0001d585'
	case "mbffraka":
		return 0x1d586 // 𝖆 '\U0001d586'
	case "mbffrakb":
		return 0x1d587 // 𝖇 '\U0001d587'
	case "mbffrakc":
		return 0x1d588 // 𝖈 '\U0001d588'
	case "mbffrakd":
		return 0x1d589 // 𝖉 '\U0001d589'
	case "mbffrake":
		return 0x1d58a // 𝖊 '\U0001d58a'
	case "mbffrakf":
		return 0x1d58b // 𝖋 '\U0001d58b'
	case "mbffrakg":
		return 0x1d58c // 𝖌 '\U0001d58c'
	case "mbffrakh":
		return 0x1d58d // 𝖍 '\U0001d58d'
	case "mbffraki":
		return 0x1d58e // 𝖎 '\U0001d58e'
	case "mbffrakj":
		return 0x1d58f // 𝖏 '\U0001d58f'
	case "mbffrakk":
		return 0x1d590 // 𝖐 '\U0001d590'
	case "mbffrakl":
		return 0x1d591 // 𝖑 '\U0001d591'
	case "mbffrakm":
		return 0x1d592 // 𝖒 '\U0001d592'
	case "mbffrakn":
		return 0x1d593 // 𝖓 '\U0001d593'
	case "mbffrako":
		return 0x1d594 // 𝖔 '\U0001d594'
	case "mbffrakp":
		return 0x1d595 // 𝖕 '\U0001d595'
	case "mbffrakq":
		return 0x1d596 // 𝖖 '\U0001d596'
	case "mbffrakr":
		return 0x1d597 // 𝖗 '\U0001d597'
	case "mbffraks":
		return 0x1d598 // 𝖘 '\U0001d598'
	case "mbffrakt":
		return 0x1d599 // 𝖙 '\U0001d599'
	case "mbffraku":
		return 0x1d59a // 𝖚 '\U0001d59a'
	case "mbffrakv":
		return 0x1d59b // 𝖛 '\U0001d59b'
	case "mbffrakw":
		return 0x1d59c // 𝖜 '\U0001d59c'
	case "mbffrakx":
		return 0x1d59d // 𝖝 '\U0001d59d'
	case "mbffraky":
		return 0x1d59e // 𝖞 '\U0001d59e'
	case "mbffrakz":
		return 0x1d59f // 𝖟 '\U0001d59f'
	case "mbfg":
		return 0x1d420 // 𝐠 '\U0001d420'
	case "mbfgamma":
		return 0x1d6c4 // 𝛄 '\U0001d6c4'
	case "mbfh":
		return 0x1d421 // 𝐡 '\U0001d421'
	case "mbfi":
		return 0x1d422 // 𝐢 '\U0001d422'
	case "mbfiota":
		return 0x1d6ca // 𝛊 '\U0001d6ca'
	case "mbfitA":
		return 0x1d468 // 𝑨 '\U0001d468'
	case "mbfitAlpha":
		return 0x1d71c // 𝜜 '\U0001d71c'
	case "mbfitB":
		return 0x1d469 // 𝑩 '\U0001d469'
	case "mbfitBeta":
		return 0x1d71d // 𝜝 '\U0001d71d'
	case "mbfitC":
		return 0x1d46a // 𝑪 '\U0001d46a'
	case "mbfitChi":
		return 0x1d732 // 𝜲 '\U0001d732'
	case "mbfitD":
		return 0x1d46b // 𝑫 '\U0001d46b'
	case "mbfitDelta":
		return 0x1d71f // 𝜟 '\U0001d71f'
	case "mbfitE":
		return 0x1d46c // 𝑬 '\U0001d46c'
	case "mbfitEpsilon":
		return 0x1d720 // 𝜠 '\U0001d720'
	case "mbfitEta":
		return 0x1d722 // 𝜢 '\U0001d722'
	case "mbfitF":
		return 0x1d46d // 𝑭 '\U0001d46d'
	case "mbfitG":
		return 0x1d46e // 𝑮 '\U0001d46e'
	case "mbfitGamma":
		return 0x1d71e // 𝜞 '\U0001d71e'
	case "mbfitH":
		return 0x1d46f // 𝑯 '\U0001d46f'
	case "mbfitI":
		return 0x1d470 // 𝑰 '\U0001d470'
	case "mbfitIota":
		return 0x1d724 // 𝜤 '\U0001d724'
	case "mbfitJ":
		return 0x1d471 // 𝑱 '\U0001d471'
	case "mbfitK":
		return 0x1d472 // 𝑲 '\U0001d472'
	case "mbfitKappa":
		return 0x1d725 // 𝜥 '\U0001d725'
	case "mbfitL":
		return 0x1d473 // 𝑳 '\U0001d473'
	case "mbfitLambda":
		return 0x1d726 // 𝜦 '\U0001d726'
	case "mbfitM":
		return 0x1d474 // 𝑴 '\U0001d474'
	case "mbfitMu":
		return 0x1d727 // 𝜧 '\U0001d727'
	case "mbfitN":
		return 0x1d475 // 𝑵 '\U0001d475'
	case "mbfitNu":
		return 0x1d728 // 𝜨 '\U0001d728'
	case "mbfitO":
		return 0x1d476 // 𝑶 '\U0001d476'
	case "mbfitOmega":
		return 0x1d734 // 𝜴 '\U0001d734'
	case "mbfitOmicron":
		return 0x1d72a // 𝜪 '\U0001d72a'
	case "mbfitP":
		return 0x1d477 // 𝑷 '\U0001d477'
	case "mbfitPhi":
		return 0x1d731 // 𝜱 '\U0001d731'
	case "mbfitPi":
		return 0x1d72b // 𝜫 '\U0001d72b'
	case "mbfitPsi":
		return 0x1d733 // 𝜳 '\U0001d733'
	case "mbfitQ":
		return 0x1d478 // 𝑸 '\U0001d478'
	case "mbfitR":
		return 0x1d479 // 𝑹 '\U0001d479'
	case "mbfitRho":
		return 0x1d72c // 𝜬 '\U0001d72c'
	case "mbfitS":
		return 0x1d47a // 𝑺 '\U0001d47a'
	case "mbfitSigma":
		return 0x1d72e // 𝜮 '\U0001d72e'
	case "mbfitT":
		return 0x1d47b // 𝑻 '\U0001d47b'
	case "mbfitTau":
		return 0x1d72f // 𝜯 '\U0001d72f'
	case "mbfitTheta":
		return 0x1d723 // 𝜣 '\U0001d723'
	case "mbfitU":
		return 0x1d47c // 𝑼 '\U0001d47c'
	case "mbfitUpsilon":
		return 0x1d730 // 𝜰 '\U0001d730'
	case "mbfitV":
		return 0x1d47d // 𝑽 '\U0001d47d'
	case "mbfitW":
		return 0x1d47e // 𝑾 '\U0001d47e'
	case "mbfitX":
		return 0x1d47f // 𝑿 '\U0001d47f'
	case "mbfitXi":
		return 0x1d729 // 𝜩 '\U0001d729'
	case "mbfitY":
		return 0x1d480 // 𝒀 '\U0001d480'
	case "mbfitZ":
		return 0x1d481 // 𝒁 '\U0001d481'
	case "mbfitZeta":
		return 0x1d721 // 𝜡 '\U0001d721'
	case "mbfita":
		return 0x1d482 // 𝒂 '\U0001d482'
	case "mbfitalpha":
		return 0x1d736 // 𝜶 '\U0001d736'
	case "mbfitb":
		return 0x1d483 // 𝒃 '\U0001d483'
	case "mbfitbeta":
		return 0x1d737 // 𝜷 '\U0001d737'
	case "mbfitc":
		return 0x1d484 // 𝒄 '\U0001d484'
	case "mbfitchi":
		return 0x1d74c // 𝝌 '\U0001d74c'
	case "mbfitd":
		return 0x1d485 // 𝒅 '\U0001d485'
	case "mbfitdelta":
		return 0x1d739 // 𝜹 '\U0001d739'
	case "mbfite":
		return 0x1d486 // 𝒆 '\U0001d486'
	case "mbfitepsilon":
		return 0x1d73a // 𝜺 '\U0001d73a'
	case "mbfiteta":
		return 0x1d73c // 𝜼 '\U0001d73c'
	case "mbfitf":
		return 0x1d487 // 𝒇 '\U0001d487'
	case "mbfitg":
		return 0x1d488 // 𝒈 '\U0001d488'
	case "mbfitgamma":
		return 0x1d738 // 𝜸 '\U0001d738'
	case "mbfith":
		return 0x1d489 // 𝒉 '\U0001d489'
	case "mbfiti":
		return 0x1d48a // 𝒊 '\U0001d48a'
	case "mbfitiota":
		return 0x1d73e // 𝜾 '\U0001d73e'
	case "mbfitj":
		return 0x1d48b // 𝒋 '\U0001d48b'
	case "mbfitk":
		return 0x1d48c // 𝒌 '\U0001d48c'
	case "mbfitkappa":
		return 0x1d73f // 𝜿 '\U0001d73f'
	case "mbfitl":
		return 0x1d48d // 𝒍 '\U0001d48d'
	case "mbfitlambda":
		return 0x1d740 // 𝝀 '\U0001d740'
	case "mbfitm":
		return 0x1d48e // 𝒎 '\U0001d48e'
	case "mbfitmu":
		return 0x1d741 // 𝝁 '\U0001d741'
	case "mbfitn":
		return 0x1d48f // 𝒏 '\U0001d48f'
	case "mbfitnabla":
		return 0x1d735 // 𝜵 '\U0001d735'
	case "mbfitnu":
		return 0x1d742 // 𝝂 '\U0001d742'
	case "mbfito":
		return 0x1d490 // 𝒐 '\U0001d490'
	case "mbfitomega":
		return 0x1d74e // 𝝎 '\U0001d74e'
	case "mbfitomicron":
		return 0x1d744 // 𝝄 '\U0001d744'
	case "mbfitp":
		return 0x1d491 // 𝒑 '\U0001d491'
	case "mbfitpartial":
		return 0x1d74f // 𝝏 '\U0001d74f'
	case "mbfitphi":
		return 0x1d74b // 𝝋 '\U0001d74b'
	case "mbfitpi":
		return 0x1d745 // 𝝅 '\U0001d745'
	case "mbfitpsi":
		return 0x1d74d // 𝝍 '\U0001d74d'
	case "mbfitq":
		return 0x1d492 // 𝒒 '\U0001d492'
	case "mbfitr":
		return 0x1d493 // 𝒓 '\U0001d493'
	case "mbfitrho":
		return 0x1d746 // 𝝆 '\U0001d746'
	case "mbfits":
		return 0x1d494 // 𝒔 '\U0001d494'
	case "mbfitsansA":
		return 0x1d63c // 𝘼 '\U0001d63c'
	case "mbfitsansAlpha":
		return 0x1d790 // 𝞐 '\U0001d790'
	case "mbfitsansB":
		return 0x1d63d // 𝘽 '\U0001d63d'
	case "mbfitsansBeta":
		return 0x1d791 // 𝞑 '\U0001d791'
	case "mbfitsansC":
		return 0x1d63e // 𝘾 '\U0001d63e'
	case "mbfitsansChi":
		return 0x1d7a6 // 𝞦 '\U0001d7a6'
	case "mbfitsansD":
		return 0x1d63f // 𝘿 '\U0001d63f'
	case "mbfitsansDelta":
		return 0x1d793 // 𝞓 '\U0001d793'
	case "mbfitsansE":
		return 0x1d640 // 𝙀 '\U0001d640'
	case "mbfitsansEpsilon":
		return 0x1d794 // 𝞔 '\U0001d794'
	case "mbfitsansEta":
		return 0x1d796 // 𝞖 '\U0001d796'
	case "mbfitsansF":
		return 0x1d641 // 𝙁 '\U0001d641'
	case "mbfitsansG":
		return 0x1d642 // 𝙂 '\U0001d642'
	case "mbfitsansGamma":
		return 0x1d792 // 𝞒 '\U0001d792'
	case "mbfitsansH":
		return 0x1d643 // 𝙃 '\U0001d643'
	case "mbfitsansI":
		return 0x1d644 // 𝙄 '\U0001d644'
	case "mbfitsansIota":
		return 0x1d798 // 𝞘 '\U0001d798'
	case "mbfitsansJ":
		return 0x1d645 // 𝙅 '\U0001d645'
	case "mbfitsansK":
		return 0x1d646 // 𝙆 '\U0001d646'
	case "mbfitsansKappa":
		return 0x1d799 // 𝞙 '\U0001d799'
	case "mbfitsansL":
		return 0x1d647 // 𝙇 '\U0001d647'
	case "mbfitsansLambda":
		return 0x1d79a // 𝞚 '\U0001d79a'
	case "mbfitsansM":
		return 0x1d648 // 𝙈 '\U0001d648'
	case "mbfitsansMu":
		return 0x1d79b // 𝞛 '\U0001d79b'
	case "mbfitsansN":
		return 0x1d649 // 𝙉 '\U0001d649'
	case "mbfitsansNu":
		return 0x1d79c // 𝞜 '\U0001d79c'
	case "mbfitsansO":
		return 0x1d64a // 𝙊 '\U0001d64a'
	case "mbfitsansOmega":
		return 0x1d7a8 // 𝞨 '\U0001d7a8'
	case "mbfitsansOmicron":
		return 0x1d79e // 𝞞 '\U0001d79e'
	case "mbfitsansP":
		return 0x1d64b // 𝙋 '\U0001d64b'
	case "mbfitsansPhi":
		return 0x1d7a5 // 𝞥 '\U0001d7a5'
	case "mbfitsansPi":
		return 0x1d79f // 𝞟 '\U0001d79f'
	case "mbfitsansPsi":
		return 0x1d7a7 // 𝞧 '\U0001d7a7'
	case "mbfitsansQ":
		return 0x1d64c // 𝙌 '\U0001d64c'
	case "mbfitsansR":
		return 0x1d64d // 𝙍 '\U0001d64d'
	case "mbfitsansRho":
		return 0x1d7a0 // 𝞠 '\U0001d7a0'
	case "mbfitsansS":
		return 0x1d64e // 𝙎 '\U0001d64e'
	case "mbfitsansSigma":
		return 0x1d7a2 // 𝞢 '\U0001d7a2'
	case "mbfitsansT":
		return 0x1d64f // 𝙏 '\U0001d64f'
	case "mbfitsansTau":
		return 0x1d7a3 // 𝞣 '\U0001d7a3'
	case "mbfitsansTheta":
		return 0x1d797 // 𝞗 '\U0001d797'
	case "mbfitsansU":
		return 0x1d650 // 𝙐 '\U0001d650'
	case "mbfitsansUpsilon":
		return 0x1d7a4 // 𝞤 '\U0001d7a4'
	case "mbfitsansV":
		return 0x1d651 // 𝙑 '\U0001d651'
	case "mbfitsansW":
		return 0x1d652 // 𝙒 '\U0001d652'
	case "mbfitsansX":
		return 0x1d653 // 𝙓 '\U0001d653'
	case "mbfitsansXi":
		return 0x1d79d // 𝞝 '\U0001d79d'
	case "mbfitsansY":
		return 0x1d654 // 𝙔 '\U0001d654'
	case "mbfitsansZ":
		return 0x1d655 // 𝙕 '\U0001d655'
	case "mbfitsansZeta":
		return 0x1d795 // 𝞕 '\U0001d795'
	case "mbfitsansa":
		return 0x1d656 // 𝙖 '\U0001d656'
	case "mbfitsansalpha":
		return 0x1d7aa // 𝞪 '\U0001d7aa'
	case "mbfitsansb":
		return 0x1d657 // 𝙗 '\U0001d657'
	case "mbfitsansbeta":
		return 0x1d7ab // 𝞫 '\U0001d7ab'
	case "mbfitsansc":
		return 0x1d658 // 𝙘 '\U0001d658'
	case "mbfitsanschi":
		return 0x1d7c0 // 𝟀 '\U0001d7c0'
	case "mbfitsansd":
		return 0x1d659 // 𝙙 '\U0001d659'
	case "mbfitsansdelta":
		return 0x1d7ad // 𝞭 '\U0001d7ad'
	case "mbfitsanse":
		return 0x1d65a // 𝙚 '\U0001d65a'
	case "mbfitsansepsilon":
		return 0x1d7ae // 𝞮 '\U0001d7ae'
	case "mbfitsanseta":
		return 0x1d7b0 // 𝞰 '\U0001d7b0'
	case "mbfitsansf":
		return 0x1d65b // 𝙛 '\U0001d65b'
	case "mbfitsansg":
		return 0x1d65c // 𝙜 '\U0001d65c'
	case "mbfitsansgamma":
		return 0x1d7ac // 𝞬 '\U0001d7ac'
	case "mbfitsansh":
		return 0x1d65d // 𝙝 '\U0001d65d'
	case "mbfitsansi":
		return 0x1d65e // 𝙞 '\U0001d65e'
	case "mbfitsansiota":
		return 0x1d7b2 // 𝞲 '\U0001d7b2'
	case "mbfitsansj":
		return 0x1d65f // 𝙟 '\U0001d65f'
	case "mbfitsansk":
		return 0x1d660 // 𝙠 '\U0001d660'
	case "mbfitsanskappa":
		return 0x1d7b3 // 𝞳 '\U0001d7b3'
	case "mbfitsansl":
		return 0x1d661 // 𝙡 '\U0001d661'
	case "mbfitsanslambda":
		return 0x1d7b4 // 𝞴 '\U0001d7b4'
	case "mbfitsansm":
		return 0x1d662 // 𝙢 '\U0001d662'
	case "mbfitsansmu":
		return 0x1d7b5 // 𝞵 '\U0001d7b5'
	case "mbfitsansn":
		return 0x1d663 // 𝙣 '\U0001d663'
	case "mbfitsansnabla":
		return 0x1d7a9 // 𝞩 '\U0001d7a9'
	case "mbfitsansnu":
		return 0x1d7b6 // 𝞶 '\U0001d7b6'
	case "mbfitsanso":
		return 0x1d664 // 𝙤 '\U0001d664'
	case "mbfitsansomega":
		return 0x1d7c2 // 𝟂 '\U0001d7c2'
	case "mbfitsansomicron":
		return 0x1d7b8 // 𝞸 '\U0001d7b8'
	case "mbfitsansp":
		return 0x1d665 // 𝙥 '\U0001d665'
	case "mbfitsanspartial":
		return 0x1d7c3 // 𝟃 '\U0001d7c3'
	case "mbfitsansphi":
		return 0x1d7bf // 𝞿 '\U0001d7bf'
	case "mbfitsanspi":
		return 0x1d7b9 // 𝞹 '\U0001d7b9'
	case "mbfitsanspsi":
		return 0x1d7c1 // 𝟁 '\U0001d7c1'
	case "mbfitsansq":
		return 0x1d666 // 𝙦 '\U0001d666'
	case "mbfitsansr":
		return 0x1d667 // 𝙧 '\U0001d667'
	case "mbfitsansrho":
		return 0x1d7ba // 𝞺 '\U0001d7ba'
	case "mbfitsanss":
		return 0x1d668 // 𝙨 '\U0001d668'
	case "mbfitsanssigma":
		return 0x1d7bc // 𝞼 '\U0001d7bc'
	case "mbfitsanst":
		return 0x1d669 // 𝙩 '\U0001d669'
	case "mbfitsanstau":
		return 0x1d7bd // 𝞽 '\U0001d7bd'
	case "mbfitsanstheta":
		return 0x1d7b1 // 𝞱 '\U0001d7b1'
	case "mbfitsansu":
		return 0x1d66a // 𝙪 '\U0001d66a'
	case "mbfitsansupsilon":
		return 0x1d7be // 𝞾 '\U0001d7be'
	case "mbfitsansv":
		return 0x1d66b // 𝙫 '\U0001d66b'
	case "mbfitsansvarTheta":
		return 0x1d7a1 // 𝞡 '\U0001d7a1'
	case "mbfitsansvarepsilon":
		return 0x1d7c4 // 𝟄 '\U0001d7c4'
	case "mbfitsansvarkappa":
		return 0x1d7c6 // 𝟆 '\U0001d7c6'
	case "mbfitsansvarphi":
		return 0x1d7c7 // 𝟇 '\U0001d7c7'
	case "mbfitsansvarpi":
		return 0x1d7c9 // 𝟉 '\U0001d7c9'
	case "mbfitsansvarrho":
		return 0x1d7c8 // 𝟈 '\U0001d7c8'
	case "mbfitsansvarsigma":
		return 0x1d7bb // 𝞻 '\U0001d7bb'
	case "mbfitsansvartheta":
		return 0x1d7c5 // 𝟅 '\U0001d7c5'
	case "mbfitsansw":
		return 0x1d66c // 𝙬 '\U0001d66c'
	case "mbfitsansx":
		return 0x1d66d // 𝙭 '\U0001d66d'
	case "mbfitsansxi":
		return 0x1d7b7 // 𝞷 '\U0001d7b7'
	case "mbfitsansy":
		return 0x1d66e // 𝙮 '\U0001d66e'
	case "mbfitsansz":
		return 0x1d66f // 𝙯 '\U0001d66f'
	case "mbfitsanszeta":
		return 0x1d7af // 𝞯 '\U0001d7af'
	case "mbfitsigma":
		return 0x1d748 // 𝝈 '\U0001d748'
	case "mbfitt":
		return 0x1d495 // 𝒕 '\U0001d495'
	case "mbfittau":
		return 0x1d749 // 𝝉 '\U0001d749'
	case "mbfittheta":
		return 0x1d73d // 𝜽 '\U0001d73d'
	case "mbfitu":
		return 0x1d496 // 𝒖 '\U0001d496'
	case "mbfitupsilon":
		return 0x1d74a // 𝝊 '\U0001d74a'
	case "mbfitv":
		return 0x1d497 // 𝒗 '\U0001d497'
	case "mbfitvarTheta":
		return 0x1d72d // 𝜭 '\U0001d72d'
	case "mbfitvarepsilon":
		return 0x1d750 // 𝝐 '\U0001d750'
	case "mbfitvarkappa":
		return 0x1d752 // 𝝒 '\U0001d752'
	case "mbfitvarphi":
		return 0x1d753 // 𝝓 '\U0001d753'
	case "mbfitvarpi":
		return 0x1d755 // 𝝕 '\U0001d755'
	case "mbfitvarrho":
		return 0x1d754 // 𝝔 '\U0001d754'
	case "mbfitvarsigma":
		return 0x1d747 // 𝝇 '\U0001d747'
	case "mbfitvartheta":
		return 0x1d751 // 𝝑 '\U0001d751'
	case "mbfitw":
		return 0x1d498 // 𝒘 '\U0001d498'
	case "mbfitx":
		return 0x1d499 // 𝒙 '\U0001d499'
	case "mbfitxi":
		return 0x1d743 // 𝝃 '\U0001d743'
	case "mbfity":
		return 0x1d49a // 𝒚 '\U0001d49a'
	case "mbfitz":
		return 0x1d49b // 𝒛 '\U0001d49b'
	case "mbfitzeta":
		return 0x1d73b // 𝜻 '\U0001d73b'
	case "mbfj":
		return 0x1d423 // 𝐣 '\U0001d423'
	case "mbfk":
		return 0x1d424 // 𝐤 '\U0001d424'
	case "mbfkappa":
		return 0x1d6cb // 𝛋 '\U0001d6cb'
	case "mbfl":
		return 0x1d425 // 𝐥 '\U0001d425'
	case "mbflambda":
		return 0x1d6cc // 𝛌 '\U0001d6cc'
	case "mbfm":
		return 0x1d426 // 𝐦 '\U0001d426'
	case "mbfmu":
		return 0x1d6cd // 𝛍 '\U0001d6cd'
	case "mbfn":
		return 0x1d427 // 𝐧 '\U0001d427'
	case "mbfnabla":
		return 0x1d6c1 // 𝛁 '\U0001d6c1'
	case "mbfnu":
		return 0x1d6ce // 𝛎 '\U0001d6ce'
	case "mbfo":
		return 0x1d428 // 𝐨 '\U0001d428'
	case "mbfomega":
		return 0x1d6da // 𝛚 '\U0001d6da'
	case "mbfomicron":
		return 0x1d6d0 // 𝛐 '\U0001d6d0'
	case "mbfp":
		return 0x1d429 // 𝐩 '\U0001d429'
	case "mbfpartial":
		return 0x1d6db // 𝛛 '\U0001d6db'
	case "mbfphi":
		return 0x1d6df // 𝛟 '\U0001d6df'
	case "mbfpi":
		return 0x1d6d1 // 𝛑 '\U0001d6d1'
	case "mbfpsi":
		return 0x1d6d9 // 𝛙 '\U0001d6d9'
	case "mbfq":
		return 0x1d42a // 𝐪 '\U0001d42a'
	case "mbfr":
		return 0x1d42b // 𝐫 '\U0001d42b'
	case "mbfrho":
		return 0x1d6d2 // 𝛒 '\U0001d6d2'
	case "mbfs":
		return 0x1d42c // 𝐬 '\U0001d42c'
	case "mbfsansA":
		return 0x1d5d4 // 𝗔 '\U0001d5d4'
	case "mbfsansAlpha":
		return 0x1d756 // 𝝖 '\U0001d756'
	case "mbfsansB":
		return 0x1d5d5 // 𝗕 '\U0001d5d5'
	case "mbfsansBeta":
		return 0x1d757 // 𝝗 '\U0001d757'
	case "mbfsansC":
		return 0x1d5d6 // 𝗖 '\U0001d5d6'
	case "mbfsansChi":
		return 0x1d76c // 𝝬 '\U0001d76c'
	case "mbfsansD":
		return 0x1d5d7 // 𝗗 '\U0001d5d7'
	case "mbfsansDelta":
		return 0x1d759 // 𝝙 '\U0001d759'
	case "mbfsansE":
		return 0x1d5d8 // 𝗘 '\U0001d5d8'
	case "mbfsansEpsilon":
		return 0x1d75a // 𝝚 '\U0001d75a'
	case "mbfsansEta":
		return 0x1d75c // 𝝜 '\U0001d75c'
	case "mbfsansF":
		return 0x1d5d9 // 𝗙 '\U0001d5d9'
	case "mbfsansG":
		return 0x1d5da // 𝗚 '\U0001d5da'
	case "mbfsansGamma":
		return 0x1d758 // 𝝘 '\U0001d758'
	case "mbfsansH":
		return 0x1d5db // 𝗛 '\U0001d5db'
	case "mbfsansI":
		return 0x1d5dc // 𝗜 '\U0001d5dc'
	case "mbfsansIota":
		return 0x1d75e // 𝝞 '\U0001d75e'
	case "mbfsansJ":
		return 0x1d5dd // 𝗝 '\U0001d5dd'
	case "mbfsansK":
		return 0x1d5de // 𝗞 '\U0001d5de'
	case "mbfsansKappa":
		return 0x1d75f // 𝝟 '\U0001d75f'
	case "mbfsansL":
		return 0x1d5df // 𝗟 '\U0001d5df'
	case "mbfsansLambda":
		return 0x1d760 // 𝝠 '\U0001d760'
	case "mbfsansM":
		return 0x1d5e0 // 𝗠 '\U0001d5e0'
	case "mbfsansMu":
		return 0x1d761 // 𝝡 '\U0001d761'
	case "mbfsansN":
		return 0x1d5e1 // 𝗡 '\U0001d5e1'
	case "mbfsansNu":
		return 0x1d762 // 𝝢 '\U0001d762'
	case "mbfsansO":
		return 0x1d5e2 // 𝗢 '\U0001d5e2'
	case "mbfsansOmega":
		return 0x1d76e // 𝝮 '\U0001d76e'
	case "mbfsansOmicron":
		return 0x1d764 // 𝝤 '\U0001d764'
	case "mbfsansP":
		return 0x1d5e3 // 𝗣 '\U0001d5e3'
	case "mbfsansPhi":
		return 0x1d76b // 𝝫 '\U0001d76b'
	case "mbfsansPi":
		return 0x1d765 // 𝝥 '\U0001d765'
	case "mbfsansPsi":
		return 0x1d76d // 𝝭 '\U0001d76d'
	case "mbfsansQ":
		return 0x1d5e4 // 𝗤 '\U0001d5e4'
	case "mbfsansR":
		return 0x1d5e5 // 𝗥 '\U0001d5e5'
	case "mbfsansRho":
		return 0x1d766 // 𝝦 '\U0001d766'
	case "mbfsansS":
		return 0x1d5e6 // 𝗦 '\U0001d5e6'
	case "mbfsansSigma":
		return 0x1d768 // 𝝨 '\U0001d768'
	case "mbfsansT":
		return 0x1d5e7 // 𝗧 '\U0001d5e7'
	case "mbfsansTau":
		return 0x1d769 // 𝝩 '\U0001d769'
	case "mbfsansTheta":
		return 0x1d75d // 𝝝 '\U0001d75d'
	case "mbfsansU":
		return 0x1d5e8 // 𝗨 '\U0001d5e8'
	case "mbfsansUpsilon":
		return 0x1d76a // 𝝪 '\U0001d76a'
	case "mbfsansV":
		return 0x1d5e9 // 𝗩 '\U0001d5e9'
	case "mbfsansW":
		return 0x1d5ea // 𝗪 '\U0001d5ea'
	case "mbfsansX":
		return 0x1d5eb // 𝗫 '\U0001d5eb'
	case "mbfsansXi":
		return 0x1d763 // 𝝣 '\U0001d763'
	case "mbfsansY":
		return 0x1d5ec // 𝗬 '\U0001d5ec'
	case "mbfsansZ":
		return 0x1d5ed // 𝗭 '\U0001d5ed'
	case "mbfsansZeta":
		return 0x1d75b // 𝝛 '\U0001d75b'
	case "mbfsansa":
		return 0x1d5ee // 𝗮 '\U0001d5ee'
	case "mbfsansalpha":
		return 0x1d770 // 𝝰 '\U0001d770'
	case "mbfsansb":
		return 0x1d5ef // 𝗯 '\U0001d5ef'
	case "mbfsansbeta":
		return 0x1d771 // 𝝱 '\U0001d771'
	case "mbfsansc":
		return 0x1d5f0 // 𝗰 '\U0001d5f0'
	case "mbfsanschi":
		return 0x1d786 // 𝞆 '\U0001d786'
	case "mbfsansd":
		return 0x1d5f1 // 𝗱 '\U0001d5f1'
	case "mbfsansdelta":
		return 0x1d773 // 𝝳 '\U0001d773'
	case "mbfsanse":
		return 0x1d5f2 // 𝗲 '\U0001d5f2'
	case "mbfsanseight":
		return 0x1d7f4 // 𝟴 '\U0001d7f4'
	case "mbfsansepsilon":
		return 0x1d774 // 𝝴 '\U0001d774'
	case "mbfsanseta":
		return 0x1d776 // 𝝶 '\U0001d776'
	case "mbfsansf":
		return 0x1d5f3 // 𝗳 '\U0001d5f3'
	case "mbfsansfive":
		return 0x1d7f1 // 𝟱 '\U0001d7f1'
	case "mbfsansfour":
		return 0x1d7f0 // 𝟰 '\U0001d7f0'
	case "mbfsansg":
		return 0x1d5f4 // 𝗴 '\U0001d5f4'
	case "mbfsansgamma":
		return 0x1d772 // 𝝲 '\U0001d772'
	case "mbfsansh":
		return 0x1d5f5 // 𝗵 '\U0001d5f5'
	case "mbfsansi":
		return 0x1d5f6 // 𝗶 '\U0001d5f6'
	case "mbfsansiota":
		return 0x1d778 // 𝝸 '\U0001d778'
	case "mbfsansj":
		return 0x1d5f7 // 𝗷 '\U0001d5f7'
	case "mbfsansk":
		return 0x1d5f8 // 𝗸 '\U0001d5f8'
	case "mbfsanskappa":
		return 0x1d779 // 𝝹 '\U0001d779'
	case "mbfsansl":
		return 0x1d5f9 // 𝗹 '\U0001d5f9'
	case "mbfsanslambda":
		return 0x1d77a // 𝝺 '\U0001d77a'
	case "mbfsansm":
		return 0x1d5fa // 𝗺 '\U0001d5fa'
	case "mbfsansmu":
		return 0x1d77b // 𝝻 '\U0001d77b'
	case "mbfsansn":
		return 0x1d5fb // 𝗻 '\U0001d5fb'
	case "mbfsansnabla":
		return 0x1d76f // 𝝯 '\U0001d76f'
	case "mbfsansnine":
		return 0x1d7f5 // 𝟵 '\U0001d7f5'
	case "mbfsansnu":
		return 0x1d77c // 𝝼 '\U0001d77c'
	case "mbfsanso":
		return 0x1d5fc // 𝗼 '\U0001d5fc'
	case "mbfsansomega":
		return 0x1d788 // 𝞈 '\U0001d788'
	case "mbfsansomicron":
		return 0x1d77e // 𝝾 '\U0001d77e'
	case "mbfsansone":
		return 0x1d7ed // 𝟭 '\U0001d7ed'
	case "mbfsansp":
		return 0x1d5fd // 𝗽 '\U0001d5fd'
	case "mbfsanspartial":
		return 0x1d789 // 𝞉 '\U0001d789'
	case "mbfsansphi":
		return 0x1d785 // 𝞅 '\U0001d785'
	case "mbfsanspi":
		return 0x1d77f // 𝝿 '\U0001d77f'
	case "mbfsanspsi":
		return 0x1d787 // 𝞇 '\U0001d787'
	case "mbfsansq":
		return 0x1d5fe // 𝗾 '\U0001d5fe'
	case "mbfsansr":
		return 0x1d5ff // 𝗿 '\U0001d5ff'
	case "mbfsansrho":
		return 0x1d780 // 𝞀 '\U0001d780'
	case "mbfsanss":
		return 0x1d600 // 𝘀 '\U0001d600'
	case "mbfsansseven":
		return 0x1d7f3 // 𝟳 '\U0001d7f3'
	case "mbfsanssigma":
		return 0x1d782 // 𝞂 '\U0001d782'
	case "mbfsanssix":
		return 0x1d7f2 // 𝟲 '\U0001d7f2'
	case "mbfsanst":
		return 0x1d601 // 𝘁 '\U0001d601'
	case "mbfsanstau":
		return 0x1d783 // 𝞃 '\U0001d783'
	case "mbfsanstheta":
		return 0x1d777 // 𝝷 '\U0001d777'
	case "mbfsansthree":
		return 0x1d7ef // 𝟯 '\U0001d7ef'
	case "mbfsanstwo":
		return 0x1d7ee // 𝟮 '\U0001d7ee'
	case "mbfsansu":
		return 0x1d602 // 𝘂 '\U0001d602'
	case "mbfsansupsilon":
		return 0x1d784 // 𝞄 '\U0001d784'
	case "mbfsansv":
		return 0x1d603 // 𝘃 '\U0001d603'
	case "mbfsansvarTheta":
		return 0x1d767 // 𝝧 '\U0001d767'
	case "mbfsansvarepsilon":
		return 0x1d78a // 𝞊 '\U0001d78a'
	case "mbfsansvarkappa":
		return 0x1d78c // 𝞌 '\U0001d78c'
	case "mbfsansvarphi":
		return 0x1d78d // 𝞍 '\U0001d78d'
	case "mbfsansvarpi":
		return 0x1d78f // 𝞏 '\U0001d78f'
	case "mbfsansvarrho":
		return 0x1d78e // 𝞎 '\U0001d78e'
	case "mbfsansvarsigma":
		return 0x1d781 // 𝞁 '\U0001d781'
	case "mbfsansvartheta":
		return 0x1d78b // 𝞋 '\U0001d78b'
	case "mbfsansw":
		return 0x1d604 // 𝘄 '\U0001d604'
	case "mbfsansx":
		return 0x1d605 // 𝘅 '\U0001d605'
	case "mbfsansxi":
		return 0x1d77d // 𝝽 '\U0001d77d'
	case "mbfsansy":
		return 0x1d606 // 𝘆 '\U0001d606'
	case "mbfsansz":
		return 0x1d607 // 𝘇 '\U0001d607'
	case "mbfsanszero":
		return 0x1d7ec // 𝟬 '\U0001d7ec'
	case "mbfsanszeta":
		return 0x1d775 // 𝝵 '\U0001d775'
	case "mbfscrA":
		return 0x1d4d0 // 𝓐 '\U0001d4d0'
	case "mbfscrB":
		return 0x1d4d1 // 𝓑 '\U0001d4d1'
	case "mbfscrC":
		return 0x1d4d2 // 𝓒 '\U0001d4d2'
	case "mbfscrD":
		return 0x1d4d3 // 𝓓 '\U0001d4d3'
	case "mbfscrE":
		return 0x1d4d4 // 𝓔 '\U0001d4d4'
	case "mbfscrF":
		return 0x1d4d5 // 𝓕 '\U0001d4d5'
	case "mbfscrG":
		return 0x1d4d6 // 𝓖 '\U0001d4d6'
	case "mbfscrH":
		return 0x1d4d7 // 𝓗 '\U0001d4d7'
	case "mbfscrI":
		return 0x1d4d8 // 𝓘 '\U0001d4d8'
	case "mbfscrJ":
		return 0x1d4d9 // 𝓙 '\U0001d4d9'
	case "mbfscrK":
		return 0x1d4da // 𝓚 '\U0001d4da'
	case "mbfscrL":
		return 0x1d4db // 𝓛 '\U0001d4db'
	case "mbfscrM":
		return 0x1d4dc // 𝓜 '\U0001d4dc'
	case "mbfscrN":
		return 0x1d4dd // 𝓝 '\U0001d4dd'
	case "mbfscrO":
		return 0x1d4de // 𝓞 '\U0001d4de'
	case "mbfscrP":
		return 0x1d4df // 𝓟 '\U0001d4df'
	case "mbfscrQ":
		return 0x1d4e0 // 𝓠 '\U0001d4e0'
	case "mbfscrR":
		return 0x1d4e1 // 𝓡 '\U0001d4e1'
	case "mbfscrS":
		return 0x1d4e2 // 𝓢 '\U0001d4e2'
	case "mbfscrT":
		return 0x1d4e3 // 𝓣 '\U0001d4e3'
	case "mbfscrU":
		return 0x1d4e4 // 𝓤 '\U0001d4e4'
	case "mbfscrV":
		return 0x1d4e5 // 𝓥 '\U0001d4e5'
	case "mbfscrW":
		return 0x1d4e6 // 𝓦 '\U0001d4e6'
	case "mbfscrX":
		return 0x1d4e7 // 𝓧 '\U0001d4e7'
	case "mbfscrY":
		return 0x1d4e8 // 𝓨 '\U0001d4e8'
	case "mbfscrZ":
		return 0x1d4e9 // 𝓩 '\U0001d4e9'
	case "mbfscra":
		return 0x1d4ea // 𝓪 '\U0001d4ea'
	case "mbfscrb":
		return 0x1d4eb // 𝓫 '\U0001d4eb'
	case "mbfscrc":
		return 0x1d4ec // 𝓬 '\U0001d4ec'
	case "mbfscrd":
		return 0x1d4ed // 𝓭 '\U0001d4ed'
	case "mbfscre":
		return 0x1d4ee // 𝓮 '\U0001d4ee'
	case "mbfscrf":
		return 0x1d4ef // 𝓯 '\U0001d4ef'
	case "mbfscrg":
		return 0x1d4f0 // 𝓰 '\U0001d4f0'
	case "mbfscrh":
		return 0x1d4f1 // 𝓱 '\U0001d4f1'
	case "mbfscri":
		return 0x1d4f2 // 𝓲 '\U0001d4f2'
	case "mbfscrj":
		return 0x1d4f3 // 𝓳 '\U0001d4f3'
	case "mbfscrk":
		return 0x1d4f4 // 𝓴 '\U0001d4f4'
	case "mbfscrl":
		return 0x1d4f5 // 𝓵 '\U0001d4f5'
	case "mbfscrm":
		return 0x1d4f6 // 𝓶 '\U0001d4f6'
	case "mbfscrn":
		return 0x1d4f7 // 𝓷 '\U0001d4f7'
	case "mbfscro":
		return 0x1d4f8 // 𝓸 '\U0001d4f8'
	case "mbfscrp":
		return 0x1d4f9 // 𝓹 '\U0001d4f9'
	case "mbfscrq":
		return 0x1d4fa // 𝓺 '\U0001d4fa'
	case "mbfscrr":
		return 0x1d4fb // 𝓻 '\U0001d4fb'
	case "mbfscrs":
		return 0x1d4fc // 𝓼 '\U0001d4fc'
	case "mbfscrt":
		return 0x1d4fd // 𝓽 '\U0001d4fd'
	case "mbfscru":
		return 0x1d4fe // 𝓾 '\U0001d4fe'
	case "mbfscrv":
		return 0x1d4ff // 𝓿 '\U0001d4ff'
	case "mbfscrw":
		return 0x1d500 // 𝔀 '\U0001d500'
	case "mbfscrx":
		return 0x1d501 // 𝔁 '\U0001d501'
	case "mbfscry":
		return 0x1d502 // 𝔂 '\U0001d502'
	case "mbfscrz":
		return 0x1d503 // 𝔃 '\U0001d503'
	case "mbfsigma":
		return 0x1d6d4 // 𝛔 '\U0001d6d4'
	case "mbft":
		return 0x1d42d // 𝐭 '\U0001d42d'
	case "mbftau":
		return 0x1d6d5 // 𝛕 '\U0001d6d5'
	case "mbftheta":
		return 0x1d6c9 // 𝛉 '\U0001d6c9'
	case "mbfu":
		return 0x1d42e // 𝐮 '\U0001d42e'
	case "mbfupsilon":
		return 0x1d6d6 // 𝛖 '\U0001d6d6'
	case "mbfv":
		return 0x1d42f // 𝐯 '\U0001d42f'
	case "mbfvarTheta":
		return 0x1d6b9 // 𝚹 '\U0001d6b9'
	case "mbfvarepsilon":
		return 0x1d6dc // 𝛜 '\U0001d6dc'
	case "mbfvarkappa":
		return 0x1d6de // 𝛞 '\U0001d6de'
	case "mbfvarphi":
		return 0x1d6d7 // 𝛗 '\U0001d6d7'
	case "mbfvarpi":
		return 0x1d6e1 // 𝛡 '\U0001d6e1'
	case "mbfvarrho":
		return 0x1d6e0 // 𝛠 '\U0001d6e0'
	case "mbfvarsigma":
		return 0x1d6d3 // 𝛓 '\U0001d6d3'
	case "mbfvartheta":
		return 0x1d6dd // 𝛝 '\U0001d6dd'
	case "mbfw":
		return 0x1d430 // 𝐰 '\U0001d430'
	case "mbfx":
		return 0x1d431 // 𝐱 '\U0001d431'
	case "mbfxi":
		return 0x1d6cf // 𝛏 '\U0001d6cf'
	case "mbfy":
		return 0x1d432 // 𝐲 '\U0001d432'
	case "mbfz":
		return 0x1d433 // 𝐳 '\U0001d433'
	case "mbfzeta":
		return 0x1d6c7 // 𝛇 '\U0001d6c7'
	case "mbopomofo":
		return 0x3107 // ㄇ '\u3107'
	case "mbsquare":
		return 0x33d4 // ㏔ '\u33d4'
	case "mcircle":
		return 0x24dc // ⓜ '\u24dc'
	case "mcubedsquare":
		return 0x33a5 // ㎥ '\u33a5'
	case "mdblkcircle":
		return 0x26ab // ⚫ '\u26ab'
	case "mdblkdiamond":
		return 0x2b25 // ⬥ '\u2b25'
	case "mdblklozenge":
		return 0x2b27 // ⬧ '\u2b27'
	case "mdblksquare":
		return 0x25fc // ◼ '\u25fc'
	case "mdlgblklozenge":
		return 0x29eb // ⧫ '\u29eb'
	case "mdotaccent":
		return 0x1e41 // ṁ '\u1e41'
	case "mdotbelow":
		return 0x1e43 // ṃ '\u1e43'
	case "mdsmblkcircle":
		return 0x2981 // ⦁ '\u2981'
	case "mdsmblksquare":
		return 0x25fe // ◾ '\u25fe'
	case "mdsmwhtcircle":
		return 0x26ac // ⚬ '\u26ac'
	case "mdsmwhtsquare":
		return 0x25fd // ◽ '\u25fd'
	case "mdwhtcircle":
		return 0x26aa // ⚪ '\u26aa'
	case "mdwhtdiamond":
		return 0x2b26 // ⬦ '\u2b26'
	case "mdwhtlozenge":
		return 0x2b28 // ⬨ '\u2b28'
	case "mdwhtsquare":
		return 0x25fb // ◻ '\u25fb'
	case "measangledltosw":
		return 0x29af // ⦯ '\u29af'
	case "measangledrtose":
		return 0x29ae // ⦮ '\u29ae'
	case "measangleldtosw":
		return 0x29ab // ⦫ '\u29ab'
	case "measanglelutonw":
		return 0x29a9 // ⦩ '\u29a9'
	case "measanglerdtose":
		return 0x29aa // ⦪ '\u29aa'
	case "measanglerutone":
		return 0x29a8 // ⦨ '\u29a8'
	case "measangleultonw":
		return 0x29ad // ⦭ '\u29ad'
	case "measangleurtone":
		return 0x29ac // ⦬ '\u29ac'
	case "measeq":
		return 0x225e // ≞ '\u225e'
	case "measuredangle":
		return 0x2221 // ∡ '\u2221'
	case "measuredangleleft":
		return 0x299b // ⦛ '\u299b'
	case "measuredrightangle":
		return 0x22be // ⊾ '\u22be'
	case "medblackstar":
		return 0x2b51 // ⭑ '\u2b51'
	case "medwhitestar":
		return 0x2b50 // ⭐ '\u2b50'
	case "meemfinalarabic":
		return 0xfee2 // ﻢ '\ufee2'
	case "meeminitialarabic":
		return 0xfee3 // ﻣ '\ufee3'
	case "meemisolated":
		return 0xfee1 // ﻡ '\ufee1'
	case "meemmedialarabic":
		return 0xfee4 // ﻤ '\ufee4'
	case "meemmeeminitialarabic":
		return 0xfcd1 // ﳑ '\ufcd1'
	case "meemmeemisolatedarabic":
		return 0xfc48 // ﱈ '\ufc48'
	case "meemwithhahinitial":
		return 0xfccf // ﳏ '\ufccf'
	case "meemwithjeeminitial":
		return 0xfcce // ﳎ '\ufcce'
	case "meemwithkhahinitial":
		return 0xfcd0 // ﳐ '\ufcd0'
	case "meetorusquare":
		return 0x334d // ㍍ '\u334d'
	case "mehiragana":
		return 0x3081 // め '\u3081'
	case "meizierasquare":
		return 0x337e // ㍾ '\u337e'
	case "mekatakana":
		return 0x30e1 // メ '\u30e1'
	case "mekatakanahalfwidth":
		return 0xff92 // ﾒ '\uff92'
	case "mem":
		return 0x05de // מ '\u05de'
	case "memdageshhebrew":
		return 0xfb3e // מּ '\ufb3e'
	case "menarmenian":
		return 0x0574 // մ '\u0574'
	case "merkhahebrew":
		return 0x05a5 // ֥ '\u05a5'
	case "merkhakefulahebrew":
		return 0x05a6 // ֦ '\u05a6'
	case "mfrakA":
		return 0x1d504 // 𝔄 '\U0001d504'
	case "mfrakB":
		return 0x1d505 // 𝔅 '\U0001d505'
	case "mfrakC":
		return 0x212d // ℭ '\u212d'
	case "mfrakD":
		return 0x1d507 // 𝔇 '\U0001d507'
	case "mfrakE":
		return 0x1d508 // 𝔈 '\U0001d508'
	case "mfrakF":
		return 0x1d509 // 𝔉 '\U0001d509'
	case "mfrakG":
		return 0x1d50a // 𝔊 '\U0001d50a'
	case "mfrakH":
		return 0x210c // ℌ '\u210c'
	case "mfrakJ":
		return 0x1d50d // 𝔍 '\U0001d50d'
	case "mfrakK":
		return 0x1d50e // 𝔎 '\U0001d50e'
	case "mfrakL":
		return 0x1d50f // 𝔏 '\U0001d50f'
	case "mfrakM":
		return 0x1d510 // 𝔐 '\U0001d510'
	case "mfrakN":
		return 0x1d511 // 𝔑 '\U0001d511'
	case "mfrakO":
		return 0x1d512 // 𝔒 '\U0001d512'
	case "mfrakP":
		return 0x1d513 // 𝔓 '\U0001d513'
	case "mfrakQ":
		return 0x1d514 // 𝔔 '\U0001d514'
	case "mfrakS":
		return 0x1d516 // 𝔖 '\U0001d516'
	case "mfrakT":
		return 0x1d517 // 𝔗 '\U0001d517'
	case "mfrakU":
		return 0x1d518 // 𝔘 '\U0001d518'
	case "mfrakV":
		return 0x1d519 // 𝔙 '\U0001d519'
	case "mfrakW":
		return 0x1d51a // 𝔚 '\U0001d51a'
	case "mfrakX":
		return 0x1d51b // 𝔛 '\U0001d51b'
	case "mfrakY":
		return 0x1d51c // 𝔜 '\U0001d51c'
	case "mfrakZ":
		return 0x2128 // ℨ '\u2128'
	case "mfraka":
		return 0x1d51e // 𝔞 '\U0001d51e'
	case "mfrakb":
		return 0x1d51f // 𝔟 '\U0001d51f'
	case "mfrakc":
		return 0x1d520 // 𝔠 '\U0001d520'
	case "mfrakd":
		return 0x1d521 // 𝔡 '\U0001d521'
	case "mfrake":
		return 0x1d522 // 𝔢 '\U0001d522'
	case "mfrakf":
		return 0x1d523 // 𝔣 '\U0001d523'
	case "mfrakg":
		return 0x1d524 // 𝔤 '\U0001d524'
	case "mfrakh":
		return 0x1d525 // 𝔥 '\U0001d525'
	case "mfraki":
		return 0x1d526 // 𝔦 '\U0001d526'
	case "mfrakj":
		return 0x1d527 // 𝔧 '\U0001d527'
	case "mfrakk":
		return 0x1d528 // 𝔨 '\U0001d528'
	case "mfrakl":
		return 0x1d529 // 𝔩 '\U0001d529'
	case "mfrakm":
		return 0x1d52a // 𝔪 '\U0001d52a'
	case "mfrakn":
		return 0x1d52b // 𝔫 '\U0001d52b'
	case "mfrako":
		return 0x1d52c // 𝔬 '\U0001d52c'
	case "mfrakp":
		return 0x1d52d // 𝔭 '\U0001d52d'
	case "mfrakq":
		return 0x1d52e // 𝔮 '\U0001d52e'
	case "mfrakr":
		return 0x1d52f // 𝔯 '\U0001d52f'
	case "mfraks":
		return 0x1d530 // 𝔰 '\U0001d530'
	case "mfrakt":
		return 0x1d531 // 𝔱 '\U0001d531'
	case "mfraku":
		return 0x1d532 // 𝔲 '\U0001d532'
	case "mfrakv":
		return 0x1d533 // 𝔳 '\U0001d533'
	case "mfrakw":
		return 0x1d534 // 𝔴 '\U0001d534'
	case "mfrakx":
		return 0x1d535 // 𝔵 '\U0001d535'
	case "mfraky":
		return 0x1d536 // 𝔶 '\U0001d536'
	case "mfrakz":
		return 0x1d537 // 𝔷 '\U0001d537'
	case "mhook":
		return 0x0271 // ɱ '\u0271'
	case "mhzsquare":
		return 0x3392 // ㎒ '\u3392'
	case "micro":
		return 0x0095 //  '\u0095'
	case "midbarvee":
		return 0x2a5d // ⩝ '\u2a5d'
	case "midbarwedge":
		return 0x2a5c // ⩜ '\u2a5c'
	case "midcir":
		return 0x2af0 // ⫰ '\u2af0'
	case "middledotkatakanahalfwidth":
		return 0xff65 // ･ '\uff65'
	case "mieumacirclekorean":
		return 0x3272 // ㉲ '\u3272'
	case "mieumaparenkorean":
		return 0x3212 // ㈒ '\u3212'
	case "mieumcirclekorean":
		return 0x3264 // ㉤ '\u3264'
	case "mieumkorean":
		return 0x3141 // ㅁ '\u3141'
	case "mieumpansioskorean":
		return 0x3170 // ㅰ '\u3170'
	case "mieumparenkorean":
		return 0x3204 // ㈄ '\u3204'
	case "mieumpieupkorean":
		return 0x316e // ㅮ '\u316e'
	case "mieumsioskorean":
		return 0x316f // ㅯ '\u316f'
	case "mihiragana":
		return 0x307f // み '\u307f'
	case "mikatakana":
		return 0x30df // ミ '\u30df'
	case "mikatakanahalfwidth":
		return 0xff90 // ﾐ '\uff90'
	case "mill":
		return 0x20a5 // ₥ '\u20a5'
	case "minus":
		return 0x2212 // − '\u2212'
	case "minusbelowcmb":
		return 0x0320 // ̠ '\u0320'
	case "minuscircle":
		return 0x2296 // ⊖ '\u2296'
	case "minusdot":
		return 0x2a2a // ⨪ '\u2a2a'
	case "minusfdots":
		return 0x2a2b // ⨫ '\u2a2b'
	case "minusinferior":
		return 0x208b // ₋ '\u208b'
	case "minusmod":
		return 0x02d7 // ˗ '\u02d7'
	case "minusplus":
		return 0x2213 // ∓ '\u2213'
	case "minusrdots":
		return 0x2a2c // ⨬ '\u2a2c'
	case "minussuperior":
		return 0x207b // ⁻ '\u207b'
	case "minute":
		return 0x2032 // ′ '\u2032'
	case "miribaarusquare":
		return 0x334a // ㍊ '\u334a'
	case "mirisquare":
		return 0x3349 // ㍉ '\u3349'
	case "mitA":
		return 0x1d434 // 𝐴 '\U0001d434'
	case "mitAlpha":
		return 0x1d6e2 // 𝛢 '\U0001d6e2'
	case "mitB":
		return 0x1d435 // 𝐵 '\U0001d435'
	case "mitBbbD":
		return 0x2145 // ⅅ '\u2145'
	case "mitBbbd":
		return 0x2146 // ⅆ '\u2146'
	case "mitBbbe":
		return 0x2147 // ⅇ '\u2147'
	case "mitBbbi":
		return 0x2148 // ⅈ '\u2148'
	case "mitBbbj":
		return 0x2149 // ⅉ '\u2149'
	case "mitBeta":
		return 0x1d6e3 // 𝛣 '\U0001d6e3'
	case "mitC":
		return 0x1d436 // 𝐶 '\U0001d436'
	case "mitChi":
		return 0x1d6f8 // 𝛸 '\U0001d6f8'
	case "mitD":
		return 0x1d437 // 𝐷 '\U0001d437'
	case "mitDelta":
		return 0x1d6e5 // 𝛥 '\U0001d6e5'
	case "mitE":
		return 0x1d438 // 𝐸 '\U0001d438'
	case "mitEpsilon":
		return 0x1d6e6 // 𝛦 '\U0001d6e6'
	case "mitEta":
		return 0x1d6e8 // 𝛨 '\U0001d6e8'
	case "mitF":
		return 0x1d439 // 𝐹 '\U0001d439'
	case "mitG":
		return 0x1d43a // 𝐺 '\U0001d43a'
	case "mitGamma":
		return 0x1d6e4 // 𝛤 '\U0001d6e4'
	case "mitH":
		return 0x1d43b // 𝐻 '\U0001d43b'
	case "mitI":
		return 0x1d43c // 𝐼 '\U0001d43c'
	case "mitIota":
		return 0x1d6ea // 𝛪 '\U0001d6ea'
	case "mitJ":
		return 0x1d43d // 𝐽 '\U0001d43d'
	case "mitK":
		return 0x1d43e // 𝐾 '\U0001d43e'
	case "mitKappa":
		return 0x1d6eb // 𝛫 '\U0001d6eb'
	case "mitL":
		return 0x1d43f // 𝐿 '\U0001d43f'
	case "mitLambda":
		return 0x1d6ec // 𝛬 '\U0001d6ec'
	case "mitM":
		return 0x1d440 // 𝑀 '\U0001d440'
	case "mitMu":
		return 0x1d6ed // 𝛭 '\U0001d6ed'
	case "mitN":
		return 0x1d441 // 𝑁 '\U0001d441'
	case "mitNu":
		return 0x1d6ee // 𝛮 '\U0001d6ee'
	case "mitO":
		return 0x1d442 // 𝑂 '\U0001d442'
	case "mitOmega":
		return 0x1d6fa // 𝛺 '\U0001d6fa'
	case "mitOmicron":
		return 0x1d6f0 // 𝛰 '\U0001d6f0'
	case "mitP":
		return 0x1d443 // 𝑃 '\U0001d443'
	case "mitPhi":
		return 0x1d6f7 // 𝛷 '\U0001d6f7'
	case "mitPi":
		return 0x1d6f1 // 𝛱 '\U0001d6f1'
	case "mitPsi":
		return 0x1d6f9 // 𝛹 '\U0001d6f9'
	case "mitQ":
		return 0x1d444 // 𝑄 '\U0001d444'
	case "mitR":
		return 0x1d445 // 𝑅 '\U0001d445'
	case "mitRho":
		return 0x1d6f2 // 𝛲 '\U0001d6f2'
	case "mitS":
		return 0x1d446 // 𝑆 '\U0001d446'
	case "mitSigma":
		return 0x1d6f4 // 𝛴 '\U0001d6f4'
	case "mitT":
		return 0x1d447 // 𝑇 '\U0001d447'
	case "mitTau":
		return 0x1d6f5 // 𝛵 '\U0001d6f5'
	case "mitTheta":
		return 0x1d6e9 // 𝛩 '\U0001d6e9'
	case "mitU":
		return 0x1d448 // 𝑈 '\U0001d448'
	case "mitUpsilon":
		return 0x1d6f6 // 𝛶 '\U0001d6f6'
	case "mitV":
		return 0x1d449 // 𝑉 '\U0001d449'
	case "mitW":
		return 0x1d44a // 𝑊 '\U0001d44a'
	case "mitX":
		return 0x1d44b // 𝑋 '\U0001d44b'
	case "mitXi":
		return 0x1d6ef // 𝛯 '\U0001d6ef'
	case "mitY":
		return 0x1d44c // 𝑌 '\U0001d44c'
	case "mitZ":
		return 0x1d44d // 𝑍 '\U0001d44d'
	case "mitZeta":
		return 0x1d6e7 // 𝛧 '\U0001d6e7'
	case "mita":
		return 0x1d44e // 𝑎 '\U0001d44e'
	case "mitalpha":
		return 0x1d6fc // 𝛼 '\U0001d6fc'
	case "mitb":
		return 0x1d44f // 𝑏 '\U0001d44f'
	case "mitbeta":
		return 0x1d6fd // 𝛽 '\U0001d6fd'
	case "mitc":
		return 0x1d450 // 𝑐 '\U0001d450'
	case "mitchi":
		return 0x1d712 // 𝜒 '\U0001d712'
	case "mitd":
		return 0x1d451 // 𝑑 '\U0001d451'
	case "mitdelta":
		return 0x1d6ff // 𝛿 '\U0001d6ff'
	case "mite":
		return 0x1d452 // 𝑒 '\U0001d452'
	case "mitepsilon":
		return 0x1d700 // 𝜀 '\U0001d700'
	case "miteta":
		return 0x1d702 // 𝜂 '\U0001d702'
	case "mitf":
		return 0x1d453 // 𝑓 '\U0001d453'
	case "mitg":
		return 0x1d454 // 𝑔 '\U0001d454'
	case "mitgamma":
		return 0x1d6fe // 𝛾 '\U0001d6fe'
	case "miti":
		return 0x1d456 // 𝑖 '\U0001d456'
	case "mitiota":
		return 0x1d704 // 𝜄 '\U0001d704'
	case "mitj":
		return 0x1d457 // 𝑗 '\U0001d457'
	case "mitk":
		return 0x1d458 // 𝑘 '\U0001d458'
	case "mitkappa":
		return 0x1d705 // 𝜅 '\U0001d705'
	case "mitl":
		return 0x1d459 // 𝑙 '\U0001d459'
	case "mitlambda":
		return 0x1d706 // 𝜆 '\U0001d706'
	case "mitm":
		return 0x1d45a // 𝑚 '\U0001d45a'
	case "mitmu":
		return 0x1d707 // 𝜇 '\U0001d707'
	case "mitn":
		return 0x1d45b // 𝑛 '\U0001d45b'
	case "mitnabla":
		return 0x1d6fb // 𝛻 '\U0001d6fb'
	case "mitnu":
		return 0x1d708 // 𝜈 '\U0001d708'
	case "mito":
		return 0x1d45c // 𝑜 '\U0001d45c'
	case "mitomega":
		return 0x1d714 // 𝜔 '\U0001d714'
	case "mitomicron":
		return 0x1d70a // 𝜊 '\U0001d70a'
	case "mitp":
		return 0x1d45d // 𝑝 '\U0001d45d'
	case "mitpartial":
		return 0x1d715 // 𝜕 '\U0001d715'
	case "mitphi":
		return 0x1d711 // 𝜑 '\U0001d711'
	case "mitpi":
		return 0x1d70b // 𝜋 '\U0001d70b'
	case "mitpsi":
		return 0x1d713 // 𝜓 '\U0001d713'
	case "mitq":
		return 0x1d45e // 𝑞 '\U0001d45e'
	case "mitr":
		return 0x1d45f // 𝑟 '\U0001d45f'
	case "mitrho":
		return 0x1d70c // 𝜌 '\U0001d70c'
	case "mits":
		return 0x1d460 // 𝑠 '\U0001d460'
	case "mitsansA":
		return 0x1d608 // 𝘈 '\U0001d608'
	case "mitsansB":
		return 0x1d609 // 𝘉 '\U0001d609'
	case "mitsansC":
		return 0x1d60a // 𝘊 '\U0001d60a'
	case "mitsansD":
		return 0x1d60b // 𝘋 '\U0001d60b'
	case "mitsansE":
		return 0x1d60c // 𝘌 '\U0001d60c'
	case "mitsansF":
		return 0x1d60d // 𝘍 '\U0001d60d'
	case "mitsansG":
		return 0x1d60e // 𝘎 '\U0001d60e'
	case "mitsansH":
		return 0x1d60f // 𝘏 '\U0001d60f'
	case "mitsansI":
		return 0x1d610 // 𝘐 '\U0001d610'
	case "mitsansJ":
		return 0x1d611 // 𝘑 '\U0001d611'
	case "mitsansK":
		return 0x1d612 // 𝘒 '\U0001d612'
	case "mitsansL":
		return 0x1d613 // 𝘓 '\U0001d613'
	case "mitsansM":
		return 0x1d614 // 𝘔 '\U0001d614'
	case "mitsansN":
		return 0x1d615 // 𝘕 '\U0001d615'
	case "mitsansO":
		return 0x1d616 // 𝘖 '\U0001d616'
	case "mitsansP":
		return 0x1d617 // 𝘗 '\U0001d617'
	case "mitsansQ":
		return 0x1d618 // 𝘘 '\U0001d618'
	case "mitsansR":
		return 0x1d619 // 𝘙 '\U0001d619'
	case "mitsansS":
		return 0x1d61a // 𝘚 '\U0001d61a'
	case "mitsansT":
		return 0x1d61b // 𝘛 '\U0001d61b'
	case "mitsansU":
		return 0x1d61c // 𝘜 '\U0001d61c'
	case "mitsansV":
		return 0x1d61d // 𝘝 '\U0001d61d'
	case "mitsansW":
		return 0x1d61e // 𝘞 '\U0001d61e'
	case "mitsansX":
		return 0x1d61f // 𝘟 '\U0001d61f'
	case "mitsansY":
		return 0x1d620 // 𝘠 '\U0001d620'
	case "mitsansZ":
		return 0x1d621 // 𝘡 '\U0001d621'
	case "mitsansa":
		return 0x1d622 // 𝘢 '\U0001d622'
	case "mitsansb":
		return 0x1d623 // 𝘣 '\U0001d623'
	case "mitsansc":
		return 0x1d624 // 𝘤 '\U0001d624'
	case "mitsansd":
		return 0x1d625 // 𝘥 '\U0001d625'
	case "mitsanse":
		return 0x1d626 // 𝘦 '\U0001d626'
	case "mitsansf":
		return 0x1d627 // 𝘧 '\U0001d627'
	case "mitsansg":
		return 0x1d628 // 𝘨 '\U0001d628'
	case "mitsansh":
		return 0x1d629 // 𝘩 '\U0001d629'
	case "mitsansi":
		return 0x1d62a // 𝘪 '\U0001d62a'
	case "mitsansj":
		return 0x1d62b // 𝘫 '\U0001d62b'
	case "mitsansk":
		return 0x1d62c // 𝘬 '\U0001d62c'
	case "mitsansl":
		return 0x1d62d // 𝘭 '\U0001d62d'
	case "mitsansm":
		return 0x1d62e // 𝘮 '\U0001d62e'
	case "mitsansn":
		return 0x1d62f // 𝘯 '\U0001d62f'
	case "mitsanso":
		return 0x1d630 // 𝘰 '\U0001d630'
	case "mitsansp":
		return 0x1d631 // 𝘱 '\U0001d631'
	case "mitsansq":
		return 0x1d632 // 𝘲 '\U0001d632'
	case "mitsansr":
		return 0x1d633 // 𝘳 '\U0001d633'
	case "mitsanss":
		return 0x1d634 // 𝘴 '\U0001d634'
	case "mitsanst":
		return 0x1d635 // 𝘵 '\U0001d635'
	case "mitsansu":
		return 0x1d636 // 𝘶 '\U0001d636'
	case "mitsansv":
		return 0x1d637 // 𝘷 '\U0001d637'
	case "mitsansw":
		return 0x1d638 // 𝘸 '\U0001d638'
	case "mitsansx":
		return 0x1d639 // 𝘹 '\U0001d639'
	case "mitsansy":
		return 0x1d63a // 𝘺 '\U0001d63a'
	case "mitsansz":
		return 0x1d63b // 𝘻 '\U0001d63b'
	case "mitsigma":
		return 0x1d70e // 𝜎 '\U0001d70e'
	case "mitt":
		return 0x1d461 // 𝑡 '\U0001d461'
	case "mittau":
		return 0x1d70f // 𝜏 '\U0001d70f'
	case "mittheta":
		return 0x1d703 // 𝜃 '\U0001d703'
	case "mitu":
		return 0x1d462 // 𝑢 '\U0001d462'
	case "mitupsilon":
		return 0x1d710 // 𝜐 '\U0001d710'
	case "mitv":
		return 0x1d463 // 𝑣 '\U0001d463'
	case "mitvarTheta":
		return 0x1d6f3 // 𝛳 '\U0001d6f3'
	case "mitvarepsilon":
		return 0x1d716 // 𝜖 '\U0001d716'
	case "mitvarkappa":
		return 0x1d718 // 𝜘 '\U0001d718'
	case "mitvarphi":
		return 0x1d719 // 𝜙 '\U0001d719'
	case "mitvarpi":
		return 0x1d71b // 𝜛 '\U0001d71b'
	case "mitvarrho":
		return 0x1d71a // 𝜚 '\U0001d71a'
	case "mitvarsigma":
		return 0x1d70d // 𝜍 '\U0001d70d'
	case "mitvartheta":
		return 0x1d717 // 𝜗 '\U0001d717'
	case "mitw":
		return 0x1d464 // 𝑤 '\U0001d464'
	case "mitx":
		return 0x1d465 // 𝑥 '\U0001d465'
	case "mitxi":
		return 0x1d709 // 𝜉 '\U0001d709'
	case "mity":
		return 0x1d466 // 𝑦 '\U0001d466'
	case "mitz":
		return 0x1d467 // 𝑧 '\U0001d467'
	case "mitzeta":
		return 0x1d701 // 𝜁 '\U0001d701'
	case "mlcp":
		return 0x2adb // ⫛ '\u2adb'
	case "mlonglegturned":
		return 0x0270 // ɰ '\u0270'
	case "mlsquare":
		return 0x3396 // ㎖ '\u3396'
	case "mmcubedsquare":
		return 0x33a3 // ㎣ '\u33a3'
	case "mmonospace":
		return 0xff4d // ｍ '\uff4d'
	case "mmsquaredsquare":
		return 0x339f // ㎟ '\u339f'
	case "models":
		return 0x22a7 // ⊧ '\u22a7'
	case "modtwosum":
		return 0x2a0a // ⨊ '\u2a0a'
	case "mohiragana":
		return 0x3082 // も '\u3082'
	case "mohmsquare":
		return 0x33c1 // ㏁ '\u33c1'
	case "mokatakana":
		return 0x30e2 // モ '\u30e2'
	case "mokatakanahalfwidth":
		return 0xff93 // ﾓ '\uff93'
	case "molsquare":
		return 0x33d6 // ㏖ '\u33d6'
	case "momathai":
		return 0x0e21 // ม '\u0e21'
	case "moverssquare":
		return 0x33a7 // ㎧ '\u33a7'
	case "moverssquaredsquare":
		return 0x33a8 // ㎨ '\u33a8'
	case "mparen":
		return 0x24a8 // ⒨ '\u24a8'
	case "mpasquare":
		return 0x33ab // ㎫ '\u33ab'
	case "msansA":
		return 0x1d5a0 // 𝖠 '\U0001d5a0'
	case "msansB":
		return 0x1d5a1 // 𝖡 '\U0001d5a1'
	case "msansC":
		return 0x1d5a2 // 𝖢 '\U0001d5a2'
	case "msansD":
		return 0x1d5a3 // 𝖣 '\U0001d5a3'
	case "msansE":
		return 0x1d5a4 // 𝖤 '\U0001d5a4'
	case "msansF":
		return 0x1d5a5 // 𝖥 '\U0001d5a5'
	case "msansG":
		return 0x1d5a6 // 𝖦 '\U0001d5a6'
	case "msansH":
		return 0x1d5a7 // 𝖧 '\U0001d5a7'
	case "msansI":
		return 0x1d5a8 // 𝖨 '\U0001d5a8'
	case "msansJ":
		return 0x1d5a9 // 𝖩 '\U0001d5a9'
	case "msansK":
		return 0x1d5aa // 𝖪 '\U0001d5aa'
	case "msansL":
		return 0x1d5ab // 𝖫 '\U0001d5ab'
	case "msansM":
		return 0x1d5ac // 𝖬 '\U0001d5ac'
	case "msansN":
		return 0x1d5ad // 𝖭 '\U0001d5ad'
	case "msansO":
		return 0x1d5ae // 𝖮 '\U0001d5ae'
	case "msansP":
		return 0x1d5af // 𝖯 '\U0001d5af'
	case "msansQ":
		return 0x1d5b0 // 𝖰 '\U0001d5b0'
	case "msansR":
		return 0x1d5b1 // 𝖱 '\U0001d5b1'
	case "msansS":
		return 0x1d5b2 // 𝖲 '\U0001d5b2'
	case "msansT":
		return 0x1d5b3 // 𝖳 '\U0001d5b3'
	case "msansU":
		return 0x1d5b4 // 𝖴 '\U0001d5b4'
	case "msansV":
		return 0x1d5b5 // 𝖵 '\U0001d5b5'
	case "msansW":
		return 0x1d5b6 // 𝖶 '\U0001d5b6'
	case "msansX":
		return 0x1d5b7 // 𝖷 '\U0001d5b7'
	case "msansY":
		return 0x1d5b8 // 𝖸 '\U0001d5b8'
	case "msansZ":
		return 0x1d5b9 // 𝖹 '\U0001d5b9'
	case "msansa":
		return 0x1d5ba // 𝖺 '\U0001d5ba'
	case "msansb":
		return 0x1d5bb // 𝖻 '\U0001d5bb'
	case "msansc":
		return 0x1d5bc // 𝖼 '\U0001d5bc'
	case "msansd":
		return 0x1d5bd // 𝖽 '\U0001d5bd'
	case "msanse":
		return 0x1d5be // 𝖾 '\U0001d5be'
	case "msanseight":
		return 0x1d7ea // 𝟪 '\U0001d7ea'
	case "msansf":
		return 0x1d5bf // 𝖿 '\U0001d5bf'
	case "msansfive":
		return 0x1d7e7 // 𝟧 '\U0001d7e7'
	case "msansfour":
		return 0x1d7e6 // 𝟦 '\U0001d7e6'
	case "msansg":
		return 0x1d5c0 // 𝗀 '\U0001d5c0'
	case "msansh":
		return 0x1d5c1 // 𝗁 '\U0001d5c1'
	case "msansi":
		return 0x1d5c2 // 𝗂 '\U0001d5c2'
	case "msansj":
		return 0x1d5c3 // 𝗃 '\U0001d5c3'
	case "msansk":
		return 0x1d5c4 // 𝗄 '\U0001d5c4'
	case "msansl":
		return 0x1d5c5 // 𝗅 '\U0001d5c5'
	case "msansm":
		return 0x1d5c6 // 𝗆 '\U0001d5c6'
	case "msansn":
		return 0x1d5c7 // 𝗇 '\U0001d5c7'
	case "msansnine":
		return 0x1d7eb // 𝟫 '\U0001d7eb'
	case "msanso":
		return 0x1d5c8 // 𝗈 '\U0001d5c8'
	case "msansone":
		return 0x1d7e3 // 𝟣 '\U0001d7e3'
	case "msansp":
		return 0x1d5c9 // 𝗉 '\U0001d5c9'
	case "msansq":
		return 0x1d5ca // 𝗊 '\U0001d5ca'
	case "msansr":
		return 0x1d5cb // 𝗋 '\U0001d5cb'
	case "msanss":
		return 0x1d5cc // 𝗌 '\U0001d5cc'
	case "msansseven":
		return 0x1d7e9 // 𝟩 '\U0001d7e9'
	case "msanssix":
		return 0x1d7e8 // 𝟨 '\U0001d7e8'
	case "msanst":
		return 0x1d5cd // 𝗍 '\U0001d5cd'
	case "msansthree":
		return 0x1d7e5 // 𝟥 '\U0001d7e5'
	case "msanstwo":
		return 0x1d7e4 // 𝟤 '\U0001d7e4'
	case "msansu":
		return 0x1d5ce // 𝗎 '\U0001d5ce'
	case "msansv":
		return 0x1d5cf // 𝗏 '\U0001d5cf'
	case "msansw":
		return 0x1d5d0 // 𝗐 '\U0001d5d0'
	case "msansx":
		return 0x1d5d1 // 𝗑 '\U0001d5d1'
	case "msansy":
		return 0x1d5d2 // 𝗒 '\U0001d5d2'
	case "msansz":
		return 0x1d5d3 // 𝗓 '\U0001d5d3'
	case "msanszero":
		return 0x1d7e2 // 𝟢 '\U0001d7e2'
	case "mscrA":
		return 0x1d49c // 𝒜 '\U0001d49c'
	case "mscrB":
		return 0x212c // ℬ '\u212c'
	case "mscrC":
		return 0x1d49e // 𝒞 '\U0001d49e'
	case "mscrD":
		return 0x1d49f // 𝒟 '\U0001d49f'
	case "mscrE":
		return 0x2130 // ℰ '\u2130'
	case "mscrF":
		return 0x2131 // ℱ '\u2131'
	case "mscrG":
		return 0x1d4a2 // 𝒢 '\U0001d4a2'
	case "mscrH":
		return 0x210b // ℋ '\u210b'
	case "mscrI":
		return 0x2110 // ℐ '\u2110'
	case "mscrJ":
		return 0x1d4a5 // 𝒥 '\U0001d4a5'
	case "mscrK":
		return 0x1d4a6 // 𝒦 '\U0001d4a6'
	case "mscrL":
		return 0x2112 // ℒ '\u2112'
	case "mscrM":
		return 0x2133 // ℳ '\u2133'
	case "mscrN":
		return 0x1d4a9 // 𝒩 '\U0001d4a9'
	case "mscrO":
		return 0x1d4aa // 𝒪 '\U0001d4aa'
	case "mscrP":
		return 0x1d4ab // 𝒫 '\U0001d4ab'
	case "mscrQ":
		return 0x1d4ac // 𝒬 '\U0001d4ac'
	case "mscrR":
		return 0x211b // ℛ '\u211b'
	case "mscrS":
		return 0x1d4ae // 𝒮 '\U0001d4ae'
	case "mscrT":
		return 0x1d4af // 𝒯 '\U0001d4af'
	case "mscrU":
		return 0x1d4b0 // 𝒰 '\U0001d4b0'
	case "mscrV":
		return 0x1d4b1 // 𝒱 '\U0001d4b1'
	case "mscrW":
		return 0x1d4b2 // 𝒲 '\U0001d4b2'
	case "mscrX":
		return 0x1d4b3 // 𝒳 '\U0001d4b3'
	case "mscrY":
		return 0x1d4b4 // 𝒴 '\U0001d4b4'
	case "mscrZ":
		return 0x1d4b5 // 𝒵 '\U0001d4b5'
	case "mscra":
		return 0x1d4b6 // 𝒶 '\U0001d4b6'
	case "mscrb":
		return 0x1d4b7 // 𝒷 '\U0001d4b7'
	case "mscrc":
		return 0x1d4b8 // 𝒸 '\U0001d4b8'
	case "mscrd":
		return 0x1d4b9 // 𝒹 '\U0001d4b9'
	case "mscre":
		return 0x212f // ℯ '\u212f'
	case "mscrf":
		return 0x1d4bb // 𝒻 '\U0001d4bb'
	case "mscrg":
		return 0x210a // ℊ '\u210a'
	case "mscrh":
		return 0x1d4bd // 𝒽 '\U0001d4bd'
	case "mscri":
		return 0x1d4be // 𝒾 '\U0001d4be'
	case "mscrj":
		return 0x1d4bf // 𝒿 '\U0001d4bf'
	case "mscrk":
		return 0x1d4c0 // 𝓀 '\U0001d4c0'
	case "mscrl":
		return 0x1d4c1 // 𝓁 '\U0001d4c1'
	case "mscrm":
		return 0x1d4c2 // 𝓂 '\U0001d4c2'
	case "mscrn":
		return 0x1d4c3 // 𝓃 '\U0001d4c3'
	case "mscro":
		return 0x2134 // ℴ '\u2134'
	case "mscrp":
		return 0x1d4c5 // 𝓅 '\U0001d4c5'
	case "mscrq":
		return 0x1d4c6 // 𝓆 '\U0001d4c6'
	case "mscrr":
		return 0x1d4c7 // 𝓇 '\U0001d4c7'
	case "mscrs":
		return 0x1d4c8 // 𝓈 '\U0001d4c8'
	case "mscrt":
		return 0x1d4c9 // 𝓉 '\U0001d4c9'
	case "mscru":
		return 0x1d4ca // 𝓊 '\U0001d4ca'
	case "mscrv":
		return 0x1d4cb // 𝓋 '\U0001d4cb'
	case "mscrw":
		return 0x1d4cc // 𝓌 '\U0001d4cc'
	case "mscrx":
		return 0x1d4cd // 𝓍 '\U0001d4cd'
	case "mscry":
		return 0x1d4ce // 𝓎 '\U0001d4ce'
	case "mscrz":
		return 0x1d4cf // 𝓏 '\U0001d4cf'
	case "mssquare":
		return 0x33b3 // ㎳ '\u33b3'
	case "msuperior":
		return 0xf6ef //  '\uf6ef'
	case "mttA":
		return 0x1d670 // 𝙰 '\U0001d670'
	case "mttB":
		return 0x1d671 // 𝙱 '\U0001d671'
	case "mttC":
		return 0x1d672 // 𝙲 '\U0001d672'
	case "mttD":
		return 0x1d673 // 𝙳 '\U0001d673'
	case "mttE":
		return 0x1d674 // 𝙴 '\U0001d674'
	case "mttF":
		return 0x1d675 // 𝙵 '\U0001d675'
	case "mttG":
		return 0x1d676 // 𝙶 '\U0001d676'
	case "mttH":
		return 0x1d677 // 𝙷 '\U0001d677'
	case "mttI":
		return 0x1d678 // 𝙸 '\U0001d678'
	case "mttJ":
		return 0x1d679 // 𝙹 '\U0001d679'
	case "mttK":
		return 0x1d67a // 𝙺 '\U0001d67a'
	case "mttL":
		return 0x1d67b // 𝙻 '\U0001d67b'
	case "mttM":
		return 0x1d67c // 𝙼 '\U0001d67c'
	case "mttN":
		return 0x1d67d // 𝙽 '\U0001d67d'
	case "mttO":
		return 0x1d67e // 𝙾 '\U0001d67e'
	case "mttP":
		return 0x1d67f // 𝙿 '\U0001d67f'
	case "mttQ":
		return 0x1d680 // 𝚀 '\U0001d680'
	case "mttR":
		return 0x1d681 // 𝚁 '\U0001d681'
	case "mttS":
		return 0x1d682 // 𝚂 '\U0001d682'
	case "mttT":
		return 0x1d683 // 𝚃 '\U0001d683'
	case "mttU":
		return 0x1d684 // 𝚄 '\U0001d684'
	case "mttV":
		return 0x1d685 // 𝚅 '\U0001d685'
	case "mttW":
		return 0x1d686 // 𝚆 '\U0001d686'
	case "mttX":
		return 0x1d687 // 𝚇 '\U0001d687'
	case "mttY":
		return 0x1d688 // 𝚈 '\U0001d688'
	case "mttZ":
		return 0x1d689 // 𝚉 '\U0001d689'
	case "mtta":
		return 0x1d68a // 𝚊 '\U0001d68a'
	case "mttb":
		return 0x1d68b // 𝚋 '\U0001d68b'
	case "mttc":
		return 0x1d68c // 𝚌 '\U0001d68c'
	case "mttd":
		return 0x1d68d // 𝚍 '\U0001d68d'
	case "mtte":
		return 0x1d68e // 𝚎 '\U0001d68e'
	case "mtteight":
		return 0x1d7fe // 𝟾 '\U0001d7fe'
	case "mttf":
		return 0x1d68f // 𝚏 '\U0001d68f'
	case "mttfive":
		return 0x1d7fb // 𝟻 '\U0001d7fb'
	case "mttfour":
		return 0x1d7fa // 𝟺 '\U0001d7fa'
	case "mttg":
		return 0x1d690 // 𝚐 '\U0001d690'
	case "mtth":
		return 0x1d691 // 𝚑 '\U0001d691'
	case "mtti":
		return 0x1d692 // 𝚒 '\U0001d692'
	case "mttj":
		return 0x1d693 // 𝚓 '\U0001d693'
	case "mttk":
		return 0x1d694 // 𝚔 '\U0001d694'
	case "mttl":
		return 0x1d695 // 𝚕 '\U0001d695'
	case "mttm":
		return 0x1d696 // 𝚖 '\U0001d696'
	case "mttn":
		return 0x1d697 // 𝚗 '\U0001d697'
	case "mttnine":
		return 0x1d7ff // 𝟿 '\U0001d7ff'
	case "mtto":
		return 0x1d698 // 𝚘 '\U0001d698'
	case "mttone":
		return 0x1d7f7 // 𝟷 '\U0001d7f7'
	case "mttp":
		return 0x1d699 // 𝚙 '\U0001d699'
	case "mttq":
		return 0x1d69a // 𝚚 '\U0001d69a'
	case "mttr":
		return 0x1d69b // 𝚛 '\U0001d69b'
	case "mtts":
		return 0x1d69c // 𝚜 '\U0001d69c'
	case "mttseven":
		return 0x1d7fd // 𝟽 '\U0001d7fd'
	case "mttsix":
		return 0x1d7fc // 𝟼 '\U0001d7fc'
	case "mttt":
		return 0x1d69d // 𝚝 '\U0001d69d'
	case "mttthree":
		return 0x1d7f9 // 𝟹 '\U0001d7f9'
	case "mtttwo":
		return 0x1d7f8 // 𝟸 '\U0001d7f8'
	case "mttu":
		return 0x1d69e // 𝚞 '\U0001d69e'
	case "mttv":
		return 0x1d69f // 𝚟 '\U0001d69f'
	case "mttw":
		return 0x1d6a0 // 𝚠 '\U0001d6a0'
	case "mttx":
		return 0x1d6a1 // 𝚡 '\U0001d6a1'
	case "mtty":
		return 0x1d6a2 // 𝚢 '\U0001d6a2'
	case "mttz":
		return 0x1d6a3 // 𝚣 '\U0001d6a3'
	case "mttzero":
		return 0x1d7f6 // 𝟶 '\U0001d7f6'
	case "mturned":
		return 0x026f // ɯ '\u026f'
	case "mu":
		return 0x00b5 // µ '\u00b5'
	case "muasquare":
		return 0x3382 // ㎂ '\u3382'
	case "muchgreater":
		return 0x226b // ≫ '\u226b'
	case "muchless":
		return 0x226a // ≪ '\u226a'
	case "mufsquare":
		return 0x338c // ㎌ '\u338c'
	case "mugreek":
		return 0x03bc // μ '\u03bc'
	case "mugsquare":
		return 0x338d // ㎍ '\u338d'
	case "muhiragana":
		return 0x3080 // む '\u3080'
	case "mukatakana":
		return 0x30e0 // ム '\u30e0'
	case "mukatakanahalfwidth":
		return 0xff91 // ﾑ '\uff91'
	case "mulsquare":
		return 0x3395 // ㎕ '\u3395'
	case "multicloseleft":
		return 0x22c9 // ⋉ '\u22c9'
	case "multicloseright":
		return 0x22ca // ⋊ '\u22ca'
	case "multimap":
		return 0x22b8 // ⊸ '\u22b8'
	case "multimapinv":
		return 0x27dc // ⟜ '\u27dc'
	case "multiopenleft":
		return 0x22cb // ⋋ '\u22cb'
	case "multiopenright":
		return 0x22cc // ⋌ '\u22cc'
	case "multiply":
		return 0x00d7 // × '\u00d7'
	case "mumsquare":
		return 0x339b // ㎛ '\u339b'
	case "munahlefthebrew":
		return 0x05a3 // ֣ '\u05a3'
	case "musicalnote":
		return 0x266a // ♪ '\u266a'
	case "musicflatsign":
		return 0x266d // ♭ '\u266d'
	case "musicsharpsign":
		return 0x266f // ♯ '\u266f'
	case "mussquare":
		return 0x33b2 // ㎲ '\u33b2'
	case "muvsquare":
		return 0x33b6 // ㎶ '\u33b6'
	case "muwsquare":
		return 0x33bc // ㎼ '\u33bc'
	case "mvmegasquare":
		return 0x33b9 // ㎹ '\u33b9'
	case "mvsquare":
		return 0x33b7 // ㎷ '\u33b7'
	case "mwmegasquare":
		return 0x33bf // ㎿ '\u33bf'
	case "mwsquare":
		return 0x33bd // ㎽ '\u33bd'
	case "n":
		return 0x006e // n 'n'
	case "nVleftarrow":
		return 0x21fa // ⇺ '\u21fa'
	case "nVleftarrowtail":
		return 0x2b3a // ⬺ '\u2b3a'
	case "nVleftrightarrow":
		return 0x21fc // ⇼ '\u21fc'
	case "nVrightarrow":
		return 0x21fb // ⇻ '\u21fb'
	case "nVrightarrowtail":
		return 0x2915 // ⤕ '\u2915'
	case "nVtwoheadleftarrow":
		return 0x2b35 // ⬵ '\u2b35'
	case "nVtwoheadleftarrowtail":
		return 0x2b3d // ⬽ '\u2b3d'
	case "nVtwoheadrightarrow":
		return 0x2901 // ⤁ '\u2901'
	case "nVtwoheadrightarrowtail":
		return 0x2918 // ⤘ '\u2918'
	case "nabengali":
		return 0x09a8 // ন '\u09a8'
	case "nacute":
		return 0x0144 // ń '\u0144'
	case "nadeva":
		return 0x0928 // न '\u0928'
	case "nagujarati":
		return 0x0aa8 // ન '\u0aa8'
	case "nagurmukhi":
		return 0x0a28 // ਨ '\u0a28'
	case "nahiragana":
		return 0x306a // な '\u306a'
	case "naira":
		return 0x20a6 // ₦ '\u20a6'
	case "nakatakana":
		return 0x30ca // ナ '\u30ca'
	case "nakatakanahalfwidth":
		return 0xff85 // ﾅ '\uff85'
	case "nand":
		return 0x22bc // ⊼ '\u22bc'
	case "napprox":
		return 0x2249 // ≉ '\u2249'
	case "nasquare":
		return 0x3381 // ㎁ '\u3381'
	case "nasymp":
		return 0x226d // ≭ '\u226d'
	case "natural":
		return 0x266e // ♮ '\u266e'
	case "nbhyphen":
		return 0x2011 // ‑ '\u2011'
	case "nbopomofo":
		return 0x310b // ㄋ '\u310b'
	case "ncaron":
		return 0x0148 // ň '\u0148'
	case "ncedilla":
		return 0x0146 // ņ '\u0146'
	case "ncedilla1":
		return 0xf81d //  '\uf81d'
	case "ncircle":
		return 0x24dd // ⓝ '\u24dd'
	case "ncircumflexbelow":
		return 0x1e4b // ṋ '\u1e4b'
	case "ndotaccent":
		return 0x1e45 // ṅ '\u1e45'
	case "ndotbelow":
		return 0x1e47 // ṇ '\u1e47'
	case "nehiragana":
		return 0x306d // ね '\u306d'
	case "nekatakana":
		return 0x30cd // ネ '\u30cd'
	case "nekatakanahalfwidth":
		return 0xff88 // ﾈ '\uff88'
	case "neovnwarrow":
		return 0x2931 // ⤱ '\u2931'
	case "neovsearrow":
		return 0x292e // ⤮ '\u292e'
	case "neswarrow":
		return 0x2922 // ⤢ '\u2922'
	case "neuter":
		return 0x26b2 // ⚲ '\u26b2'
	case "nfsquare":
		return 0x338b // ㎋ '\u338b'
	case "ngabengali":
		return 0x0999 // ঙ '\u0999'
	case "ngadeva":
		return 0x0919 // ङ '\u0919'
	case "ngagujarati":
		return 0x0a99 // ઙ '\u0a99'
	case "ngagurmukhi":
		return 0x0a19 // ਙ '\u0a19'
	case "ngonguthai":
		return 0x0e07 // ง '\u0e07'
	case "ngtrsim":
		return 0x2275 // ≵ '\u2275'
	case "nhVvert":
		return 0x2af5 // ⫵ '\u2af5'
	case "nhiragana":
		return 0x3093 // ん '\u3093'
	case "nhookleft":
		return 0x0272 // ɲ '\u0272'
	case "nhookretroflex":
		return 0x0273 // ɳ '\u0273'
	case "nhpar":
		return 0x2af2 // ⫲ '\u2af2'
	case "nieunacirclekorean":
		return 0x326f // ㉯ '\u326f'
	case "nieunaparenkorean":
		return 0x320f // ㈏ '\u320f'
	case "nieuncieuckorean":
		return 0x3135 // ㄵ '\u3135'
	case "nieuncirclekorean":
		return 0x3261 // ㉡ '\u3261'
	case "nieunhieuhkorean":
		return 0x3136 // ㄶ '\u3136'
	case "nieunkorean":
		return 0x3134 // ㄴ '\u3134'
	case "nieunpansioskorean":
		return 0x3168 // ㅨ '\u3168'
	case "nieunparenkorean":
		return 0x3201 // ㈁ '\u3201'
	case "nieunsioskorean":
		return 0x3167 // ㅧ '\u3167'
	case "nieuntikeutkorean":
		return 0x3166 // ㅦ '\u3166'
	case "nihiragana":
		return 0x306b // に '\u306b'
	case "nikatakana":
		return 0x30cb // ニ '\u30cb'
	case "nikatakanahalfwidth":
		return 0xff86 // ﾆ '\uff86'
	case "nikhahitleftthai":
		return 0xf899 //  '\uf899'
	case "nikhahitthai":
		return 0x0e4d // ํ '\u0e4d'
	case "nine":
		return 0x0039 // 9 '9'
	case "ninebengali":
		return 0x09ef // ৯ '\u09ef'
	case "ninedeva":
		return 0x096f // ९ '\u096f'
	case "ninegujarati":
		return 0x0aef // ૯ '\u0aef'
	case "ninegurmukhi":
		return 0x0a6f // ੯ '\u0a6f'
	case "ninehackarabic":
		return 0x0669 // ٩ '\u0669'
	case "ninehangzhou":
		return 0x3029 // 〩 '\u3029'
	case "nineideographicparen":
		return 0x3228 // ㈨ '\u3228'
	case "nineinferior":
		return 0x2089 // ₉ '\u2089'
	case "ninemonospace":
		return 0xff19 // ９ '\uff19'
	case "nineoldstyle":
		return 0xf739 //  '\uf739'
	case "nineparen":
		return 0x247c // ⑼ '\u247c'
	case "nineperiod":
		return 0x2490 // ⒐ '\u2490'
	case "ninepersian":
		return 0x06f9 // ۹ '\u06f9'
	case "nineroman":
		return 0x2178 // ⅸ '\u2178'
	case "ninesuperior":
		return 0x2079 // ⁹ '\u2079'
	case "nineteencircle":
		return 0x2472 // ⑲ '\u2472'
	case "nineteenparen":
		return 0x2486 // ⒆ '\u2486'
	case "nineteenperiod":
		return 0x249a // ⒚ '\u249a'
	case "ninethai":
		return 0x0e59 // ๙ '\u0e59'
	case "niobar":
		return 0x22fe // ⋾ '\u22fe'
	case "nis":
		return 0x22fc // ⋼ '\u22fc'
	case "nisd":
		return 0x22fa // ⋺ '\u22fa'
	case "nj":
		return 0x01cc // ǌ '\u01cc'
	case "nkatakana":
		return 0x30f3 // ン '\u30f3'
	case "nkatakanahalfwidth":
		return 0xff9d // ﾝ '\uff9d'
	case "nlegrightlong":
		return 0x019e // ƞ '\u019e'
	case "nlessgtr":
		return 0x2278 // ≸ '\u2278'
	case "nlesssim":
		return 0x2274 // ≴ '\u2274'
	case "nlinebelow":
		return 0x1e49 // ṉ '\u1e49'
	case "nmonospace":
		return 0xff4e // ｎ '\uff4e'
	case "nmsquare":
		return 0x339a // ㎚ '\u339a'
	case "nnabengali":
		return 0x09a3 // ণ '\u09a3'
	case "nnadeva":
		return 0x0923 // ण '\u0923'
	case "nnagujarati":
		return 0x0aa3 // ણ '\u0aa3'
	case "nnagurmukhi":
		return 0x0a23 // ਣ '\u0a23'
	case "nnnadeva":
		return 0x0929 // ऩ '\u0929'
	case "nohiragana":
		return 0x306e // の '\u306e'
	case "nokatakana":
		return 0x30ce // ノ '\u30ce'
	case "nokatakanahalfwidth":
		return 0xff89 // ﾉ '\uff89'
	case "nonbreakingspace":
		return 0x00a0 //  '\u00a0'
	case "nonenthai":
		return 0x0e13 // ณ '\u0e13'
	case "nonuthai":
		return 0x0e19 // น '\u0e19'
	case "noonarabic":
		return 0x0646 // ن '\u0646'
	case "noonfinalarabic":
		return 0xfee6 // ﻦ '\ufee6'
	case "noonghunnafinalarabic":
		return 0xfb9f // ﮟ '\ufb9f'
	case "noonhehinitialarabic":
		return 0xfee7 // ﻧ '\ufee7'
	case "noonisolated":
		return 0xfee5 // ﻥ '\ufee5'
	case "noonjeeminitialarabic":
		return 0xfcd2 // ﳒ '\ufcd2'
	case "noonjeemisolatedarabic":
		return 0xfc4b // ﱋ '\ufc4b'
	case "noonmedialarabic":
		return 0xfee8 // ﻨ '\ufee8'
	case "noonmeeminitialarabic":
		return 0xfcd5 // ﳕ '\ufcd5'
	case "noonmeemisolatedarabic":
		return 0xfc4e // ﱎ '\ufc4e'
	case "noonnoonfinalarabic":
		return 0xfc8d // ﲍ '\ufc8d'
	case "noonwithalefmaksurafinal":
		return 0xfc8e // ﲎ '\ufc8e'
	case "noonwithalefmaksuraisolated":
		return 0xfc4f // ﱏ '\ufc4f'
	case "noonwithhahinitial":
		return 0xfcd3 // ﳓ '\ufcd3'
	case "noonwithhehinitial":
		return 0xe815 //  '\ue815'
	case "noonwithkhahinitial":
		return 0xfcd4 // ﳔ '\ufcd4'
	case "noonwithyehfinal":
		return 0xfc8f // ﲏ '\ufc8f'
	case "noonwithyehisolated":
		return 0xfc50 // ﱐ '\ufc50'
	case "noonwithzainfinal":
		return 0xfc70 // ﱰ '\ufc70'
	case "notapproxequal":
		return 0x2247 // ≇ '\u2247'
	case "notarrowboth":
		return 0x21ae // ↮ '\u21ae'
	case "notarrowleft":
		return 0x219a // ↚ '\u219a'
	case "notarrowright":
		return 0x219b // ↛ '\u219b'
	case "notbar":
		return 0x2224 // ∤ '\u2224'
	case "notcontains":
		return 0x220c // ∌ '\u220c'
	case "notdblarrowboth":
		return 0x21ce // ⇎ '\u21ce'
	case "notelement":
		return 0x2209 // ∉ '\u2209'
	case "notequal":
		return 0x2260 // ≠ '\u2260'
	case "notexistential":
		return 0x2204 // ∄ '\u2204'
	case "notforces":
		return 0x22ae // ⊮ '\u22ae'
	case "notforcesextra":
		return 0x22af // ⊯ '\u22af'
	case "notgreater":
		return 0x226f // ≯ '\u226f'
	case "notgreaternorequal":
		return 0x2271 // ≱ '\u2271'
	case "notgreaternorless":
		return 0x2279 // ≹ '\u2279'
	case "notgreaterorslnteql":
		return 0x2a7e // ⩾ '\u2a7e'
	case "notidentical":
		return 0x2262 // ≢ '\u2262'
	case "notless":
		return 0x226e // ≮ '\u226e'
	case "notlessnorequal":
		return 0x2270 // ≰ '\u2270'
	case "notparallel":
		return 0x2226 // ∦ '\u2226'
	case "notprecedes":
		return 0x2280 // ⊀ '\u2280'
	case "notsatisfies":
		return 0x22ad // ⊭ '\u22ad'
	case "notsimilar":
		return 0x2241 // ≁ '\u2241'
	case "notsubset":
		return 0x2284 // ⊄ '\u2284'
	case "notsubseteql":
		return 0x2288 // ⊈ '\u2288'
	case "notsucceeds":
		return 0x2281 // ⊁ '\u2281'
	case "notsuperset":
		return 0x2285 // ⊅ '\u2285'
	case "notsuperseteql":
		return 0x2289 // ⊉ '\u2289'
	case "nottriangeqlleft":
		return 0x22ec // ⋬ '\u22ec'
	case "nottriangeqlright":
		return 0x22ed // ⋭ '\u22ed'
	case "nottriangleleft":
		return 0x22ea // ⋪ '\u22ea'
	case "nottriangleright":
		return 0x22eb // ⋫ '\u22eb'
	case "notturnstile":
		return 0x22ac // ⊬ '\u22ac'
	case "nowarmenian":
		return 0x0576 // ն '\u0576'
	case "nparen":
		return 0x24a9 // ⒩ '\u24a9'
	case "npolint":
		return 0x2a14 // ⨔ '\u2a14'
	case "npreccurlyeq":
		return 0x22e0 // ⋠ '\u22e0'
	case "nsime":
		return 0x2244 // ≄ '\u2244'
	case "nsqsubseteq":
		return 0x22e2 // ⋢ '\u22e2'
	case "nsqsupseteq":
		return 0x22e3 // ⋣ '\u22e3'
	case "nssquare":
		return 0x33b1 // ㎱ '\u33b1'
	case "nsucccurlyeq":
		return 0x22e1 // ⋡ '\u22e1'
	case "nsuperior":
		return 0x207f // ⁿ '\u207f'
	case "ntilde":
		return 0x00f1 // ñ '\u00f1'
	case "nu":
		return 0x03bd // ν '\u03bd'
	case "nuhiragana":
		return 0x306c // ぬ '\u306c'
	case "nukatakana":
		return 0x30cc // ヌ '\u30cc'
	case "nukatakanahalfwidth":
		return 0xff87 // ﾇ '\uff87'
	case "nuktabengali":
		return 0x09bc // ় '\u09bc'
	case "nuktadeva":
		return 0x093c // ़ '\u093c'
	case "nuktagujarati":
		return 0x0abc // ઼ '\u0abc'
	case "nuktagurmukhi":
		return 0x0a3c // ਼ '\u0a3c'
	case "numbersign":
		return 0x0023 // # '#'
	case "numbersignmonospace":
		return 0xff03 // ＃ '\uff03'
	case "numbersignsmall":
		return 0xfe5f // ﹟ '\ufe5f'
	case "numeralsigngreek":
		return 0x0374 // ʹ '\u0374'
	case "numeralsignlowergreek":
		return 0x0375 // ͵ '\u0375'
	case "numero":
		return 0x2116 // № '\u2116'
	case "nun":
		return 0x05e0 // נ '\u05e0'
	case "nundagesh":
		return 0xfb40 // נּ '\ufb40'
	case "nvLeftarrow":
		return 0x2902 // ⤂ '\u2902'
	case "nvLeftrightarrow":
		return 0x2904 // ⤄ '\u2904'
	case "nvRightarrow":
		return 0x2903 // ⤃ '\u2903'
	case "nvinfty":
		return 0x29de // ⧞ '\u29de'
	case "nvleftarrow":
		return 0x21f7 // ⇷ '\u21f7'
	case "nvleftarrowtail":
		return 0x2b39 // ⬹ '\u2b39'
	case "nvleftrightarrow":
		return 0x21f9 // ⇹ '\u21f9'
	case "nvrightarrow":
		return 0x21f8 // ⇸ '\u21f8'
	case "nvrightarrowtail":
		return 0x2914 // ⤔ '\u2914'
	case "nvsquare":
		return 0x33b5 // ㎵ '\u33b5'
	case "nvtwoheadleftarrow":
		return 0x2b34 // ⬴ '\u2b34'
	case "nvtwoheadleftarrowtail":
		return 0x2b3c // ⬼ '\u2b3c'
	case "nvtwoheadrightarrow":
		return 0x2900 // ⤀ '\u2900'
	case "nvtwoheadrightarrowtail":
		return 0x2917 // ⤗ '\u2917'
	case "nwovnearrow":
		return 0x2932 // ⤲ '\u2932'
	case "nwsearrow":
		return 0x2921 // ⤡ '\u2921'
	case "nwsquare":
		return 0x33bb // ㎻ '\u33bb'
	case "nyabengali":
		return 0x099e // ঞ '\u099e'
	case "nyadeva":
		return 0x091e // ञ '\u091e'
	case "nyagujarati":
		return 0x0a9e // ઞ '\u0a9e'
	case "nyagurmukhi":
		return 0x0a1e // ਞ '\u0a1e'
	case "o":
		return 0x006f // o 'o'
	case "oacute":
		return 0x00f3 // ó '\u00f3'
	case "oangthai":
		return 0x0e2d // อ '\u0e2d'
	case "obar":
		return 0x233d // ⌽ '\u233d'
	case "obarred":
		return 0x0275 // ɵ '\u0275'
	case "obarredcyrillic":
		return 0x04e9 // ө '\u04e9'
	case "obarreddieresiscyrillic":
		return 0x04eb // ӫ '\u04eb'
	case "obengali":
		return 0x0993 // ও '\u0993'
	case "obopomofo":
		return 0x311b // ㄛ '\u311b'
	case "obot":
		return 0x29ba // ⦺ '\u29ba'
	case "obrbrak":
		return 0x23e0 // ⏠ '\u23e0'
	case "obreve":
		return 0x014f // ŏ '\u014f'
	case "obslash":
		return 0x29b8 // ⦸ '\u29b8'
	case "ocandradeva":
		return 0x0911 // ऑ '\u0911'
	case "ocandragujarati":
		return 0x0a91 // ઑ '\u0a91'
	case "ocandravowelsigndeva":
		return 0x0949 // ॉ '\u0949'
	case "ocandravowelsigngujarati":
		return 0x0ac9 // ૉ '\u0ac9'
	case "ocaron":
		return 0x01d2 // ǒ '\u01d2'
	case "ocircle":
		return 0x24de // ⓞ '\u24de'
	case "ocircumflex":
		return 0x00f4 // ô '\u00f4'
	case "ocircumflexacute":
		return 0x1ed1 // ố '\u1ed1'
	case "ocircumflexdotbelow":
		return 0x1ed9 // ộ '\u1ed9'
	case "ocircumflexgrave":
		return 0x1ed3 // ồ '\u1ed3'
	case "ocircumflexhookabove":
		return 0x1ed5 // ổ '\u1ed5'
	case "ocircumflextilde":
		return 0x1ed7 // ỗ '\u1ed7'
	case "ocyrillic":
		return 0x043e // о '\u043e'
	case "odblgrave":
		return 0x020d // ȍ '\u020d'
	case "odeva":
		return 0x0913 // ओ '\u0913'
	case "odieresis":
		return 0x00f6 // ö '\u00f6'
	case "odieresiscyrillic":
		return 0x04e7 // ӧ '\u04e7'
	case "odiv":
		return 0x2a38 // ⨸ '\u2a38'
	case "odotbelow":
		return 0x1ecd // ọ '\u1ecd'
	case "odotslashdot":
		return 0x29bc // ⦼ '\u29bc'
	case "oe":
		return 0x0153 // œ '\u0153'
	case "oekorean":
		return 0x315a // ㅚ '\u315a'
	case "ogonek":
		return 0x02db // ˛ '\u02db'
	case "ogonekcmb":
		return 0x0328 // ̨ '\u0328'
	case "ograve":
		return 0x00f2 // ò '\u00f2'
	case "ogreaterthan":
		return 0x29c1 // ⧁ '\u29c1'
	case "ogujarati":
		return 0x0a93 // ઓ '\u0a93'
	case "oharmenian":
		return 0x0585 // օ '\u0585'
	case "ohiragana":
		return 0x304a // お '\u304a'
	case "ohookabove":
		return 0x1ecf // ỏ '\u1ecf'
	case "ohorn":
		return 0x01a1 // ơ '\u01a1'
	case "ohornacute":
		return 0x1edb // ớ '\u1edb'
	case "ohorndotbelow":
		return 0x1ee3 // ợ '\u1ee3'
	case "ohorngrave":
		return 0x1edd // ờ '\u1edd'
	case "ohornhookabove":
		return 0x1edf // ở '\u1edf'
	case "ohorntilde":
		return 0x1ee1 // ỡ '\u1ee1'
	case "ohungarumlaut":
		return 0x0151 // ő '\u0151'
	case "oi":
		return 0x01a3 // ƣ '\u01a3'
	case "oiiint":
		return 0x2230 // ∰ '\u2230'
	case "oiint":
		return 0x222f // ∯ '\u222f'
	case "ointctrclockwise":
		return 0x2233 // ∳ '\u2233'
	case "oinvertedbreve":
		return 0x020f // ȏ '\u020f'
	case "okatakana":
		return 0x30aa // オ '\u30aa'
	case "okatakanahalfwidth":
		return 0xff75 // ｵ '\uff75'
	case "okorean":
		return 0x3157 // ㅗ '\u3157'
	case "olcross":
		return 0x29bb // ⦻ '\u29bb'
	case "olehebrew":
		return 0x05ab // ֫ '\u05ab'
	case "olessthan":
		return 0x29c0 // ⧀ '\u29c0'
	case "omacron":
		return 0x014d // ō '\u014d'
	case "omacronacute":
		return 0x1e53 // ṓ '\u1e53'
	case "omacrongrave":
		return 0x1e51 // ṑ '\u1e51'
	case "omdeva":
		return 0x0950 // ॐ '\u0950'
	case "omega":
		return 0x03c9 // ω '\u03c9'
	case "omega1":
		return 0x03d6 // ϖ '\u03d6'
	case "omegacyrillic":
		return 0x0461 // ѡ '\u0461'
	case "omegalatinclosed":
		return 0x0277 // ɷ '\u0277'
	case "omegaroundcyrillic":
		return 0x047b // ѻ '\u047b'
	case "omegatitlocyrillic":
		return 0x047d // ѽ '\u047d'
	case "omegatonos":
		return 0x03ce // ώ '\u03ce'
	case "omgujarati":
		return 0x0ad0 // ૐ '\u0ad0'
	case "omicron":
		return 0x03bf // ο '\u03bf'
	case "omicrontonos":
		return 0x03cc // ό '\u03cc'
	case "omonospace":
		return 0xff4f // ｏ '\uff4f'
	case "one":
		return 0x0031 // 1 '1'
	case "onebengali":
		return 0x09e7 // ১ '\u09e7'
	case "onedeva":
		return 0x0967 // १ '\u0967'
	case "onedotenleader":
		return 0x2024 // ․ '\u2024'
	case "oneeighth":
		return 0x215b // ⅛ '\u215b'
	case "onefifth":
		return 0x2155 // ⅕ '\u2155'
	case "onefitted":
		return 0xf6dc //  '\uf6dc'
	case "onegujarati":
		return 0x0ae7 // ૧ '\u0ae7'
	case "onegurmukhi":
		return 0x0a67 // ੧ '\u0a67'
	case "onehackarabic":
		return 0x0661 // ١ '\u0661'
	case "onehalf":
		return 0x00bd // ½ '\u00bd'
	case "onehangzhou":
		return 0x3021 // 〡 '\u3021'
	case "oneideographicparen":
		return 0x3220 // ㈠ '\u3220'
	case "oneinferior":
		return 0x2081 // ₁ '\u2081'
	case "onemonospace":
		return 0xff11 // １ '\uff11'
	case "onenumeratorbengali":
		return 0x09f4 // ৴ '\u09f4'
	case "oneoldstyle":
		return 0xf731 //  '\uf731'
	case "oneparen":
		return 0x2474 // ⑴ '\u2474'
	case "oneperiod":
		return 0x2488 // ⒈ '\u2488'
	case "onepersian":
		return 0x06f1 // ۱ '\u06f1'
	case "onequarter":
		return 0x00bc // ¼ '\u00bc'
	case "oneroman":
		return 0x2170 // ⅰ '\u2170'
	case "onesixth":
		return 0x2159 // ⅙ '\u2159'
	case "onesuperior":
		return 0x00b9 // ¹ '\u00b9'
	case "onethai":
		return 0x0e51 // ๑ '\u0e51'
	case "onethird":
		return 0x2153 // ⅓ '\u2153'
	case "oogonek":
		return 0x01eb // ǫ '\u01eb'
	case "oogonekmacron":
		return 0x01ed // ǭ '\u01ed'
	case "oogurmukhi":
		return 0x0a13 // ਓ '\u0a13'
	case "oomatragurmukhi":
		return 0x0a4b // ੋ '\u0a4b'
	case "oopen":
		return 0x0254 // ɔ '\u0254'
	case "oparen":
		return 0x24aa // ⒪ '\u24aa'
	case "operp":
		return 0x29b9 // ⦹ '\u29b9'
	case "opluslhrim":
		return 0x2a2d // ⨭ '\u2a2d'
	case "oplusrhrim":
		return 0x2a2e // ⨮ '\u2a2e'
	case "option":
		return 0x2325 // ⌥ '\u2325'
	case "ordfeminine":
		return 0x00aa // ª '\u00aa'
	case "ordmasculine":
		return 0x00ba // º '\u00ba'
	case "origof":
		return 0x22b6 // ⊶ '\u22b6'
	case "orthogonal":
		return 0x221f // ∟ '\u221f'
	case "orunderscore":
		return 0x22bb // ⊻ '\u22bb'
	case "oshortdeva":
		return 0x0912 // ऒ '\u0912'
	case "oshortvowelsigndeva":
		return 0x094a // ॊ '\u094a'
	case "oslash":
		return 0x00f8 // ø '\u00f8'
	case "oslashacute":
		return 0x01ff // ǿ '\u01ff'
	case "osmallhiragana":
		return 0x3049 // ぉ '\u3049'
	case "osmallkatakana":
		return 0x30a9 // ォ '\u30a9'
	case "osmallkatakanahalfwidth":
		return 0xff6b // ｫ '\uff6b'
	case "osuperior":
		return 0xf6f0 //  '\uf6f0'
	case "otcyrillic":
		return 0x047f // ѿ '\u047f'
	case "otilde":
		return 0x00f5 // õ '\u00f5'
	case "otildeacute":
		return 0x1e4d // ṍ '\u1e4d'
	case "otildedieresis":
		return 0x1e4f // ṏ '\u1e4f'
	case "otimeshat":
		return 0x2a36 // ⨶ '\u2a36'
	case "otimeslhrim":
		return 0x2a34 // ⨴ '\u2a34'
	case "otimesrhrim":
		return 0x2a35 // ⨵ '\u2a35'
	case "oubopomofo":
		return 0x3121 // ㄡ '\u3121'
	case "ounce":
		return 0x2125 // ℥ '\u2125'
	case "overbrace":
		return 0x23de // ⏞ '\u23de'
	case "overbracket":
		return 0x23b4 // ⎴ '\u23b4'
	case "overleftarrow":
		return 0x20d6 // ⃖ '\u20d6'
	case "overleftrightarrow":
		return 0x20e1 // ⃡ '\u20e1'
	case "overline":
		return 0x203e // ‾ '\u203e'
	case "overlinecenterline":
		return 0xfe4a // ﹊ '\ufe4a'
	case "overlinecmb":
		return 0x0305 // ̅ '\u0305'
	case "overlinedashed":
		return 0xfe49 // ﹉ '\ufe49'
	case "overlinedblwavy":
		return 0xfe4c // ﹌ '\ufe4c'
	case "overlinewavy":
		return 0xfe4b // ﹋ '\ufe4b'
	case "overparen":
		return 0x23dc // ⏜ '\u23dc'
	case "ovowelsignbengali":
		return 0x09cb // ো '\u09cb'
	case "ovowelsigndeva":
		return 0x094b // ो '\u094b'
	case "ovowelsigngujarati":
		return 0x0acb // ો '\u0acb'
	case "p":
		return 0x0070 // p 'p'
	case "paampssquare":
		return 0x3380 // ㎀ '\u3380'
	case "paasentosquare":
		return 0x332b // ㌫ '\u332b'
	case "pabengali":
		return 0x09aa // প '\u09aa'
	case "pacute":
		return 0x1e55 // ṕ '\u1e55'
	case "padeva":
		return 0x092a // प '\u092a'
	case "pagedown":
		return 0x21df // ⇟ '\u21df'
	case "pageup":
		return 0x21de // ⇞ '\u21de'
	case "pagujarati":
		return 0x0aaa // પ '\u0aaa'
	case "pagurmukhi":
		return 0x0a2a // ਪ '\u0a2a'
	case "pahiragana":
		return 0x3071 // ぱ '\u3071'
	case "paiyannoithai":
		return 0x0e2f // ฯ '\u0e2f'
	case "pakatakana":
		return 0x30d1 // パ '\u30d1'
	case "palatalizationcyrilliccmb":
		return 0x0484 // ҄ '\u0484'
	case "palochkacyrillic":
		return 0x04c0 // Ӏ '\u04c0'
	case "pansioskorean":
		return 0x317f // ㅿ '\u317f'
	case "paragraph":
		return 0x00b6 // ¶ '\u00b6'
	case "paragraphseparator":
		return 0x2029 //  '\u2029'
	case "parallel":
		return 0x2225 // ∥ '\u2225'
	case "parallelogram":
		return 0x25b1 // ▱ '\u25b1'
	case "parallelogramblack":
		return 0x25b0 // ▰ '\u25b0'
	case "parenleft":
		return 0x0028 // ( '('
	case "parenleftaltonearabic":
		return 0xfd3e // ﴾ '\ufd3e'
	case "parenleftbt":
		return 0xf8ed //  '\uf8ed'
	case "parenleftex":
		return 0xf8ec //  '\uf8ec'
	case "parenleftinferior":
		return 0x208d // ₍ '\u208d'
	case "parenleftmonospace":
		return 0xff08 // （ '\uff08'
	case "parenleftsmall":
		return 0xfe59 // ﹙ '\ufe59'
	case "parenleftsuperior":
		return 0x207d // ⁽ '\u207d'
	case "parenlefttp":
		return 0xf8eb //  '\uf8eb'
	case "parenleftvertical":
		return 0xfe35 // ︵ '\ufe35'
	case "parenright":
		return 0x0029 // ) ')'
	case "parenrightaltonearabic":
		return 0xfd3f // ﴿ '\ufd3f'
	case "parenrightbt":
		return 0xf8f8 //  '\uf8f8'
	case "parenrightex":
		return 0xf8f7 //  '\uf8f7'
	case "parenrightinferior":
		return 0x208e // ₎ '\u208e'
	case "parenrightmonospace":
		return 0xff09 // ） '\uff09'
	case "parenrightsmall":
		return 0xfe5a // ﹚ '\ufe5a'
	case "parenrightsuperior":
		return 0x207e // ⁾ '\u207e'
	case "parenrighttp":
		return 0xf8f6 //  '\uf8f6'
	case "parenrightvertical":
		return 0xfe36 // ︶ '\ufe36'
	case "parsim":
		return 0x2af3 // ⫳ '\u2af3'
	case "partialdiff":
		return 0x2202 // ∂ '\u2202'
	case "partialmeetcontraction":
		return 0x2aa3 // ⪣ '\u2aa3'
	case "pashtahebrew":
		return 0x0599 // ֙ '\u0599'
	case "pasquare":
		return 0x33a9 // ㎩ '\u33a9'
	case "patah11":
		return 0x05b7 // ַ '\u05b7'
	case "pazerhebrew":
		return 0x05a1 // ֡ '\u05a1'
	case "pbopomofo":
		return 0x3106 // ㄆ '\u3106'
	case "pcircle":
		return 0x24df // ⓟ '\u24df'
	case "pdotaccent":
		return 0x1e57 // ṗ '\u1e57'
	case "pecyrillic":
		return 0x043f // п '\u043f'
	case "pedagesh":
		return 0xfb44 // פּ '\ufb44'
	case "peezisquare":
		return 0x333b // ㌻ '\u333b'
	case "pefinaldageshhebrew":
		return 0xfb43 // ףּ '\ufb43'
	case "peharabic":
		return 0x067e // پ '\u067e'
	case "peharmenian":
		return 0x057a // պ '\u057a'
	case "pehfinalarabic":
		return 0xfb57 // ﭗ '\ufb57'
	case "pehinitialarabic":
		return 0xfb58 // ﭘ '\ufb58'
	case "pehiragana":
		return 0x307a // ぺ '\u307a'
	case "pehisolated":
		return 0xfb56 // ﭖ '\ufb56'
	case "pehmedialarabic":
		return 0xfb59 // ﭙ '\ufb59'
	case "pehwithhehinitial":
		return 0xe813 //  '\ue813'
	case "pekatakana":
		return 0x30da // ペ '\u30da'
	case "pemiddlehookcyrillic":
		return 0x04a7 // ҧ '\u04a7'
	case "pentagon":
		return 0x2b20 // ⬠ '\u2b20'
	case "pentagonblack":
		return 0x2b1f // ⬟ '\u2b1f'
	case "perafehebrew":
		return 0xfb4e // פֿ '\ufb4e'
	case "percent":
		return 0x0025 // % '%'
	case "percentarabic":
		return 0x066a // ٪ '\u066a'
	case "percentmonospace":
		return 0xff05 // ％ '\uff05'
	case "percentsmall":
		return 0xfe6a // ﹪ '\ufe6a'
	case "period":
		return 0x002e // . '.'
	case "periodarmenian":
		return 0x0589 // ։ '\u0589'
	case "periodcentered":
		return 0x00b7 // · '\u00b7'
	case "periodcentered.0":
		return 0x0097 //  '\u0097'
	case "periodhalfwidth":
		return 0xff61 // ｡ '\uff61'
	case "periodinferior":
		return 0xf6e7 //  '\uf6e7'
	case "periodmonospace":
		return 0xff0e // ． '\uff0e'
	case "periodsmall":
		return 0xfe52 // ﹒ '\ufe52'
	case "periodsuperior":
		return 0xf6e8 //  '\uf6e8'
	case "perispomenigreekcmb":
		return 0x0342 // ͂ '\u0342'
	case "perp":
		return 0x27c2 // ⟂ '\u27c2'
	case "perpcorrespond":
		return 0x2a5e // ⩞ '\u2a5e'
	case "perpendicular":
		return 0x22a5 // ⊥ '\u22a5'
	case "perps":
		return 0x2ae1 // ⫡ '\u2ae1'
	case "pertenthousand":
		return 0x2031 // ‱ '\u2031'
	case "perthousand":
		return 0x2030 // ‰ '\u2030'
	case "peseta":
		return 0x20a7 // ₧ '\u20a7'
	case "peso1":
		return 0xf81b //  '\uf81b'
	case "pfsquare":
		return 0x338a // ㎊ '\u338a'
	case "phabengali":
		return 0x09ab // ফ '\u09ab'
	case "phadeva":
		return 0x092b // फ '\u092b'
	case "phagujarati":
		return 0x0aab // ફ '\u0aab'
	case "phagurmukhi":
		return 0x0a2b // ਫ '\u0a2b'
	case "phi":
		return 0x03c6 // φ '\u03c6'
	case "phi1":
		return 0x03d5 // ϕ '\u03d5'
	case "phieuphacirclekorean":
		return 0x327a // ㉺ '\u327a'
	case "phieuphaparenkorean":
		return 0x321a // ㈚ '\u321a'
	case "phieuphcirclekorean":
		return 0x326c // ㉬ '\u326c'
	case "phieuphkorean":
		return 0x314d // ㅍ '\u314d'
	case "phieuphparenkorean":
		return 0x320c // ㈌ '\u320c'
	case "philatin":
		return 0x0278 // ɸ '\u0278'
	case "phinthuthai":
		return 0x0e3a // ฺ '\u0e3a'
	case "phook":
		return 0x01a5 // ƥ '\u01a5'
	case "phophanthai":
		return 0x0e1e // พ '\u0e1e'
	case "phophungthai":
		return 0x0e1c // ผ '\u0e1c'
	case "phosamphaothai":
		return 0x0e20 // ภ '\u0e20'
	case "pi":
		return 0x03c0 // π '\u03c0'
	case "pieupacirclekorean":
		return 0x3273 // ㉳ '\u3273'
	case "pieupaparenkorean":
		return 0x3213 // ㈓ '\u3213'
	case "pieupcieuckorean":
		return 0x3176 // ㅶ '\u3176'
	case "pieupcirclekorean":
		return 0x3265 // ㉥ '\u3265'
	case "pieupkiyeokkorean":
		return 0x3172 // ㅲ '\u3172'
	case "pieupkorean":
		return 0x3142 // ㅂ '\u3142'
	case "pieupparenkorean":
		return 0x3205 // ㈅ '\u3205'
	case "pieupsioskiyeokkorean":
		return 0x3174 // ㅴ '\u3174'
	case "pieupsioskorean":
		return 0x3144 // ㅄ '\u3144'
	case "pieupsiostikeutkorean":
		return 0x3175 // ㅵ '\u3175'
	case "pieupthieuthkorean":
		return 0x3177 // ㅷ '\u3177'
	case "pieuptikeutkorean":
		return 0x3173 // ㅳ '\u3173'
	case "pihiragana":
		return 0x3074 // ぴ '\u3074'
	case "pikatakana":
		return 0x30d4 // ピ '\u30d4'
	case "piwrarmenian":
		return 0x0583 // փ '\u0583'
	case "planckover2pi":
		return 0x210f // ℏ '\u210f'
	case "plus":
		return 0x002b // + '+'
	case "plusbelowcmb":
		return 0x031f // ̟ '\u031f'
	case "plusdot":
		return 0x2a25 // ⨥ '\u2a25'
	case "pluseqq":
		return 0x2a72 // ⩲ '\u2a72'
	case "plushat":
		return 0x2a23 // ⨣ '\u2a23'
	case "plusinferior":
		return 0x208a // ₊ '\u208a'
	case "plusminus":
		return 0x00b1 // ± '\u00b1'
	case "plusmod":
		return 0x02d6 // ˖ '\u02d6'
	case "plusmonospace":
		return 0xff0b // ＋ '\uff0b'
	case "plussim":
		return 0x2a26 // ⨦ '\u2a26'
	case "plussmall":
		return 0xfe62 // ﹢ '\ufe62'
	case "plussubtwo":
		return 0x2a27 // ⨧ '\u2a27'
	case "plussuperior":
		return 0x207a // ⁺ '\u207a'
	case "plustrif":
		return 0x2a28 // ⨨ '\u2a28'
	case "pmonospace":
		return 0xff50 // ｐ '\uff50'
	case "pmsquare":
		return 0x33d8 // ㏘ '\u33d8'
	case "pohiragana":
		return 0x307d // ぽ '\u307d'
	case "pointingindexdownwhite":
		return 0x261f // ☟ '\u261f'
	case "pointingindexleftwhite":
		return 0x261c // ☜ '\u261c'
	case "pointingindexupwhite":
		return 0x261d // ☝ '\u261d'
	case "pointint":
		return 0x2a15 // ⨕ '\u2a15'
	case "pokatakana":
		return 0x30dd // ポ '\u30dd'
	case "poplathai":
		return 0x0e1b // ป '\u0e1b'
	case "postalmark":
		return 0x3012 // 〒 '\u3012'
	case "postalmarkface":
		return 0x3020 // 〠 '\u3020'
	case "pparen":
		return 0x24ab // ⒫ '\u24ab'
	case "precapprox":
		return 0x2ab7 // ⪷ '\u2ab7'
	case "precedenotdbleqv":
		return 0x2ab9 // ⪹ '\u2ab9'
	case "precedenotslnteql":
		return 0x2ab5 // ⪵ '\u2ab5'
	case "precedeornoteqvlnt":
		return 0x22e8 // ⋨ '\u22e8'
	case "precedes":
		return 0x227a // ≺ '\u227a'
	case "precedesequal":
		return 0x2aaf // ⪯ '\u2aaf'
	case "precedesorcurly":
		return 0x227c // ≼ '\u227c'
	case "precedesorequal":
		return 0x227e // ≾ '\u227e'
	case "preceqq":
		return 0x2ab3 // ⪳ '\u2ab3'
	case "precneq":
		return 0x2ab1 // ⪱ '\u2ab1'
	case "prescription":
		return 0x211e // ℞ '\u211e'
	case "primedblmod":
		return 0x0243 // Ƀ '\u0243'
	case "primemod":
		return 0x02b9 // ʹ '\u02b9'
	case "primereversed":
		return 0x2035 // ‵ '\u2035'
	case "product":
		return 0x220f // ∏ '\u220f'
	case "profsurf":
		return 0x2313 // ⌓ '\u2313'
	case "projective":
		return 0x2305 // ⌅ '\u2305'
	case "prolongedkana":
		return 0x30fc // ー '\u30fc'
	case "propellor":
		return 0x2318 // ⌘ '\u2318'
	case "propersubset":
		return 0x2282 // ⊂ '\u2282'
	case "propersuperset":
		return 0x2283 // ⊃ '\u2283'
	case "proportion":
		return 0x2237 // ∷ '\u2237'
	case "proportional":
		return 0x221d // ∝ '\u221d'
	case "prurel":
		return 0x22b0 // ⊰ '\u22b0'
	case "psi":
		return 0x03c8 // ψ '\u03c8'
	case "psicyrillic":
		return 0x0471 // ѱ '\u0471'
	case "psilipneumatacyrilliccmb":
		return 0x0486 // ҆ '\u0486'
	case "pssquare":
		return 0x33b0 // ㎰ '\u33b0'
	case "puhiragana":
		return 0x3077 // ぷ '\u3077'
	case "pukatakana":
		return 0x30d7 // プ '\u30d7'
	case "pullback":
		return 0x27d3 // ⟓ '\u27d3'
	case "punctuationspace":
		return 0x2008 //  '\u2008'
	case "pushout":
		return 0x27d4 // ⟔ '\u27d4'
	case "pvsquare":
		return 0x33b4 // ㎴ '\u33b4'
	case "pwsquare":
		return 0x33ba // ㎺ '\u33ba'
	case "q":
		return 0x0071 // q 'q'
	case "qadeva":
		return 0x0958 // क़ '\u0958'
	case "qadmahebrew":
		return 0x05a8 // ֨ '\u05a8'
	case "qaffinalarabic":
		return 0xfed6 // ﻖ '\ufed6'
	case "qafinitialarabic":
		return 0xfed7 // ﻗ '\ufed7'
	case "qafisolated":
		return 0xfed5 // ﻕ '\ufed5'
	case "qafmedialarabic":
		return 0xfed8 // ﻘ '\ufed8'
	case "qarneyparahebrew":
		return 0x059f // ֟ '\u059f'
	case "qbopomofo":
		return 0x3111 // ㄑ '\u3111'
	case "qcircle":
		return 0x24e0 // ⓠ '\u24e0'
	case "qhook":
		return 0x02a0 // ʠ '\u02a0'
	case "qmonospace":
		return 0xff51 // ｑ '\uff51'
	case "qofdagesh":
		return 0xfb47 // קּ '\ufb47'
	case "qofqubutshebrew":
		return 0x05e7 // ק '\u05e7'
	case "qparen":
		return 0x24ac // ⒬ '\u24ac'
	case "qprime":
		return 0x2057 // ⁗ '\u2057'
	case "quarternote":
		return 0x2669 // ♩ '\u2669'
	case "qubutswidehebrew":
		return 0x05bb // ֻ '\u05bb'
	case "questeq":
		return 0x225f // ≟ '\u225f'
	case "question":
		return 0x003f // ? '?'
	case "questionarmenian":
		return 0x055e // ՞ '\u055e'
	case "questiondown":
		return 0x00bf // ¿ '\u00bf'
	case "questiondownsmall":
		return 0xf7bf //  '\uf7bf'
	case "questiongreek":
		return 0x037e // ; '\u037e'
	case "questionmonospace":
		return 0xff1f // ？ '\uff1f'
	case "questionsmall":
		return 0xf73f //  '\uf73f'
	case "quotedbl":
		return 0x0022 // " '"'
	case "quotedblbase":
		return 0x201e // „ '\u201e'
	case "quotedblleft":
		return 0x201c // “ '\u201c'
	case "quotedblmonospace":
		return 0xff02 // ＂ '\uff02'
	case "quotedblprime":
		return 0x301e // 〞 '\u301e'
	case "quotedblprimereversed":
		return 0x301d // 〝 '\u301d'
	case "quotedblrev":
		return 0x201f // ‟ '\u201f'
	case "quotedblright":
		return 0x201d // ” '\u201d'
	case "quoteleft":
		return 0x2018 // ‘ '\u2018'
	case "quoteleftmod":
		return 0x0244 // Ʉ '\u0244'
	case "quotereversed":
		return 0x201b // ‛ '\u201b'
	case "quoteright":
		return 0x2019 // ’ '\u2019'
	case "quoterightn":
		return 0x0149 // ŉ '\u0149'
	case "quotesinglbase":
		return 0x201a // ‚ '\u201a'
	case "quotesingle":
		return 0x0027 // \' '\''
	case "quotesinglemonospace":
		return 0xff07 // ＇ '\uff07'
	case "r":
		return 0x0072 // r 'r'
	case "rAngle":
		return 0x27eb // ⟫ '\u27eb'
	case "rBrace":
		return 0x2984 // ⦄ '\u2984'
	case "rParen":
		return 0x2986 // ⦆ '\u2986'
	case "raarmenian":
		return 0x057c // ռ '\u057c'
	case "rabengali":
		return 0x09b0 // র '\u09b0'
	case "racute":
		return 0x0155 // ŕ '\u0155'
	case "radeva":
		return 0x0930 // र '\u0930'
	case "radical":
		return 0x221a // √ '\u221a'
	case "radicalex":
		return 0xf8e5 //  '\uf8e5'
	case "radoverssquare":
		return 0x33ae // ㎮ '\u33ae'
	case "radoverssquaredsquare":
		return 0x33af // ㎯ '\u33af'
	case "radsquare":
		return 0x33ad // ㎭ '\u33ad'
	case "ragujarati":
		return 0x0ab0 // ર '\u0ab0'
	case "ragurmukhi":
		return 0x0a30 // ਰ '\u0a30'
	case "rahiragana":
		return 0x3089 // ら '\u3089'
	case "raised":
		return 0x024d // ɍ '\u024d'
	case "rakatakana":
		return 0x30e9 // ラ '\u30e9'
	case "rakatakanahalfwidth":
		return 0xff97 // ﾗ '\uff97'
	case "ralowerdiagonalbengali":
		return 0x09f1 // ৱ '\u09f1'
	case "ramiddlediagonalbengali":
		return 0x09f0 // ৰ '\u09f0'
	case "ramshorn":
		return 0x0264 // ɤ '\u0264'
	case "rangledot":
		return 0x2992 // ⦒ '\u2992'
	case "rangledownzigzagarrow":
		return 0x237c // ⍼ '\u237c'
	case "ratio":
		return 0x2236 // ∶ '\u2236'
	case "rayaleflam":
		return 0xe816 //  '\ue816'
	case "rbag":
		return 0x27c6 // ⟆ '\u27c6'
	case "rblkbrbrak":
		return 0x2998 // ⦘ '\u2998'
	case "rbopomofo":
		return 0x3116 // ㄖ '\u3116'
	case "rbracelend":
		return 0x23ad // ⎭ '\u23ad'
	case "rbracemid":
		return 0x23ac // ⎬ '\u23ac'
	case "rbraceuend":
		return 0x23ab // ⎫ '\u23ab'
	case "rbrackextender":
		return 0x23a5 // ⎥ '\u23a5'
	case "rbracklend":
		return 0x23a6 // ⎦ '\u23a6'
	case "rbracklrtick":
		return 0x298e // ⦎ '\u298e'
	case "rbrackubar":
		return 0x298c // ⦌ '\u298c'
	case "rbrackuend":
		return 0x23a4 // ⎤ '\u23a4'
	case "rbrackurtick":
		return 0x2990 // ⦐ '\u2990'
	case "rbrbrak":
		return 0x2773 // ❳ '\u2773'
	case "rcaron":
		return 0x0159 // ř '\u0159'
	case "rcedilla":
		return 0x0157 // ŗ '\u0157'
	case "rcedilla1":
		return 0xf81f //  '\uf81f'
	case "rcircle":
		return 0x24e1 // ⓡ '\u24e1'
	case "rcircumflex":
		return 0xf832 //  '\uf832'
	case "rcurvyangle":
		return 0x29fd // ⧽ '\u29fd'
	case "rdblgrave":
		return 0x0211 // ȑ '\u0211'
	case "rdiagovfdiag":
		return 0x292b // ⤫ '\u292b'
	case "rdiagovsearrow":
		return 0x2930 // ⤰ '\u2930'
	case "rdotaccent":
		return 0x1e59 // ṙ '\u1e59'
	case "rdotbelow":
		return 0x1e5b // ṛ '\u1e5b'
	case "rdotbelowmacron":
		return 0x1e5d // ṝ '\u1e5d'
	case "recordright":
		return 0x2117 // ℗ '\u2117'
	case "referencemark":
		return 0x203b // ※ '\u203b'
	case "reflexsubset":
		return 0x2286 // ⊆ '\u2286'
	case "reflexsuperset":
		return 0x2287 // ⊇ '\u2287'
	case "registered":
		return 0x00ae // ® '\u00ae'
	case "registersans":
		return 0xf8e8 //  '\uf8e8'
	case "registerserif":
		return 0xf6da //  '\uf6da'
	case "reharabic":
		return 0x0631 // ر '\u0631'
	case "reharmenian":
		return 0x0580 // ր '\u0580'
	case "rehfinalarabic":
		return 0xfeae // ﺮ '\ufeae'
	case "rehiragana":
		return 0x308c // れ '\u308c'
	case "rehisolated":
		return 0xfead // ﺭ '\ufead'
	case "rekatakana":
		return 0x30ec // レ '\u30ec'
	case "rekatakanahalfwidth":
		return 0xff9a // ﾚ '\uff9a'
	case "reshdageshhebrew":
		return 0xfb48 // רּ '\ufb48'
	case "reshhiriq":
		return 0x05e8 // ר '\u05e8'
	case "response":
		return 0x211f // ℟ '\u211f'
	case "revangle":
		return 0x29a3 // ⦣ '\u29a3'
	case "revangleubar":
		return 0x29a5 // ⦥ '\u29a5'
	case "revasymptequal":
		return 0x22cd // ⋍ '\u22cd'
	case "revemptyset":
		return 0x29b0 // ⦰ '\u29b0'
	case "reversedtilde":
		return 0x223d // ∽ '\u223d'
	case "reviamugrashhebrew":
		return 0x0597 // ֗ '\u0597'
	case "revlogicalnot":
		return 0x2310 // ⌐ '\u2310'
	case "revnmid":
		return 0x2aee // ⫮ '\u2aee'
	case "rfbowtie":
		return 0x29d2 // ⧒ '\u29d2'
	case "rfishhook":
		return 0x027e // ɾ '\u027e'
	case "rfishhookreversed":
		return 0x027f // ɿ '\u027f'
	case "rftimes":
		return 0x29d5 // ⧕ '\u29d5'
	case "rhabengali":
		return 0x09dd // ঢ় '\u09dd'
	case "rhadeva":
		return 0x095d // ढ़ '\u095d'
	case "rho":
		return 0x03c1 // ρ '\u03c1'
	case "rhook":
		return 0x027d // ɽ '\u027d'
	case "rhookturned":
		return 0x027b // ɻ '\u027b'
	case "rhookturnedsuperior":
		return 0x02b5 // ʵ '\u02b5'
	case "rhosymbolgreek":
		return 0x03f1 // ϱ '\u03f1'
	case "rhotichookmod":
		return 0x02de // ˞ '\u02de'
	case "rieulacirclekorean":
		return 0x3271 // ㉱ '\u3271'
	case "rieulaparenkorean":
		return 0x3211 // ㈑ '\u3211'
	case "rieulcirclekorean":
		return 0x3263 // ㉣ '\u3263'
	case "rieulhieuhkorean":
		return 0x3140 // ㅀ '\u3140'
	case "rieulkiyeokkorean":
		return 0x313a // ㄺ '\u313a'
	case "rieulkiyeoksioskorean":
		return 0x3169 // ㅩ '\u3169'
	case "rieulkorean":
		return 0x3139 // ㄹ '\u3139'
	case "rieulmieumkorean":
		return 0x313b // ㄻ '\u313b'
	case "rieulpansioskorean":
		return 0x316c // ㅬ '\u316c'
	case "rieulparenkorean":
		return 0x3203 // ㈃ '\u3203'
	case "rieulphieuphkorean":
		return 0x313f // ㄿ '\u313f'
	case "rieulpieupkorean":
		return 0x313c // ㄼ '\u313c'
	case "rieulpieupsioskorean":
		return 0x316b // ㅫ '\u316b'
	case "rieulsioskorean":
		return 0x313d // ㄽ '\u313d'
	case "rieulthieuthkorean":
		return 0x313e // ㄾ '\u313e'
	case "rieultikeutkorean":
		return 0x316a // ㅪ '\u316a'
	case "rieulyeorinhieuhkorean":
		return 0x316d // ㅭ '\u316d'
	case "rightanglemdot":
		return 0x299d // ⦝ '\u299d'
	case "rightanglene":
		return 0x231d // ⌝ '\u231d'
	case "rightanglenw":
		return 0x231c // ⌜ '\u231c'
	case "rightanglese":
		return 0x231f // ⌟ '\u231f'
	case "rightanglesqr":
		return 0x299c // ⦜ '\u299c'
	case "rightanglesw":
		return 0x231e // ⌞ '\u231e'
	case "rightarrowapprox":
		return 0x2975 // ⥵ '\u2975'
	case "rightarrowbackapprox":
		return 0x2b48 // ⭈ '\u2b48'
	case "rightarrowbsimilar":
		return 0x2b4c // ⭌ '\u2b4c'
	case "rightarrowdiamond":
		return 0x291e // ⤞ '\u291e'
	case "rightarrowgtr":
		return 0x2b43 // ⭃ '\u2b43'
	case "rightarrowonoplus":
		return 0x27f4 // ⟴ '\u27f4'
	case "rightarrowplus":
		return 0x2945 // ⥅ '\u2945'
	case "rightarrowshortleftarrow":
		return 0x2942 // ⥂ '\u2942'
	case "rightarrowsimilar":
		return 0x2974 // ⥴ '\u2974'
	case "rightarrowsupset":
		return 0x2b44 // ⭄ '\u2b44'
	case "rightarrowtriangle":
		return 0x21fe // ⇾ '\u21fe'
	case "rightarrowx":
		return 0x2947 // ⥇ '\u2947'
	case "rightbkarrow":
		return 0x290d // ⤍ '\u290d'
	case "rightcurvedarrow":
		return 0x2933 // ⤳ '\u2933'
	case "rightdbltail":
		return 0x291c // ⤜ '\u291c'
	case "rightdotarrow":
		return 0x2911 // ⤑ '\u2911'
	case "rightdowncurvedarrow":
		return 0x2937 // ⤷ '\u2937'
	case "rightfishtail":
		return 0x297d // ⥽ '\u297d'
	case "rightharpoonaccent":
		return 0x20d1 // ⃑ '\u20d1'
	case "rightharpoondownbar":
		return 0x2957 // ⥗ '\u2957'
	case "rightharpoonsupdown":
		return 0x2964 // ⥤ '\u2964'
	case "rightharpoonupbar":
		return 0x2953 // ⥓ '\u2953'
	case "rightharpoonupdash":
		return 0x296c // ⥬ '\u296c'
	case "rightimply":
		return 0x2970 // ⥰ '\u2970'
	case "rightleftharpoonsdown":
		return 0x2969 // ⥩ '\u2969'
	case "rightleftharpoonsup":
		return 0x2968 // ⥨ '\u2968'
	case "rightmoon":
		return 0x263d // ☽ '\u263d'
	case "rightouterjoin":
		return 0x27d6 // ⟖ '\u27d6'
	case "rightpentagon":
		return 0x2b54 // ⭔ '\u2b54'
	case "rightpentagonblack":
		return 0x2b53 // ⭓ '\u2b53'
	case "rightrightarrows":
		return 0x21c9 // ⇉ '\u21c9'
	case "righttackbelowcmb":
		return 0x0319 // ̙ '\u0319'
	case "righttail":
		return 0x291a // ⤚ '\u291a'
	case "rightthreearrows":
		return 0x21f6 // ⇶ '\u21f6'
	case "righttriangle":
		return 0x22bf // ⊿ '\u22bf'
	case "rightwavearrow":
		return 0x219d // ↝ '\u219d'
	case "rihiragana":
		return 0x308a // り '\u308a'
	case "rikatakana":
		return 0x30ea // リ '\u30ea'
	case "rikatakanahalfwidth":
		return 0xff98 // ﾘ '\uff98'
	case "ring":
		return 0x02da // ˚ '\u02da'
	case "ring1":
		return 0xf007 //  '\uf007'
	case "ringbelowcmb":
		return 0x0325 // ̥ '\u0325'
	case "ringcmb":
		return 0x030a // ̊ '\u030a'
	case "ringfitted":
		return 0xd80d //  '\ufffd'
	case "ringhalfleft":
		return 0x02bf // ʿ '\u02bf'
	case "ringhalfleftarmenian":
		return 0x0559 // ՙ '\u0559'
	case "ringhalfleftbelowcmb":
		return 0x031c // ̜ '\u031c'
	case "ringhalfleftcentered":
		return 0x02d3 // ˓ '\u02d3'
	case "ringhalfright":
		return 0x02be // ʾ '\u02be'
	case "ringhalfrightbelowcmb":
		return 0x0339 // ̹ '\u0339'
	case "ringhalfrightcentered":
		return 0x02d2 // ˒ '\u02d2'
	case "ringinequal":
		return 0x2256 // ≖ '\u2256'
	case "ringlefthalfsubnosp":
		return 0x028f // ʏ '\u028f'
	case "ringlefthalfsuper":
		return 0x0248 // Ɉ '\u0248'
	case "ringplus":
		return 0x2a22 // ⨢ '\u2a22'
	case "ringrighthalfsubnosp":
		return 0x02ac // ʬ '\u02ac'
	case "ringrighthalfsuper":
		return 0x0247 // ɇ '\u0247'
	case "rinvertedbreve":
		return 0x0213 // ȓ '\u0213'
	case "rittorusquare":
		return 0x3351 // ㍑ '\u3351'
	case "rle":
		return 0x202b //  '\u202b'
	case "rlinebelow":
		return 0x1e5f // ṟ '\u1e5f'
	case "rlongleg":
		return 0x027c // ɼ '\u027c'
	case "rlonglegturned":
		return 0x027a // ɺ '\u027a'
	case "rmonospace":
		return 0xff52 // ｒ '\uff52'
	case "rmoustache":
		return 0x23b1 // ⎱ '\u23b1'
	case "rohiragana":
		return 0x308d // ろ '\u308d'
	case "rokatakana":
		return 0x30ed // ロ '\u30ed'
	case "rokatakanahalfwidth":
		return 0xff9b // ﾛ '\uff9b'
	case "roruathai":
		return 0x0e23 // ร '\u0e23'
	case "rparen":
		return 0x24ad // ⒭ '\u24ad'
	case "rparenextender":
		return 0x239f // ⎟ '\u239f'
	case "rparengtr":
		return 0x2994 // ⦔ '\u2994'
	case "rparenlend":
		return 0x23a0 // ⎠ '\u23a0'
	case "rparenuend":
		return 0x239e // ⎞ '\u239e'
	case "rppolint":
		return 0x2a12 // ⨒ '\u2a12'
	case "rrabengali":
		return 0x09dc // ড় '\u09dc'
	case "rradeva":
		return 0x0931 // ऱ '\u0931'
	case "rragurmukhi":
		return 0x0a5c // ੜ '\u0a5c'
	case "rrangle":
		return 0x298a // ⦊ '\u298a'
	case "rreharabic":
		return 0x0691 // ڑ '\u0691'
	case "rrehfinalarabic":
		return 0xfb8d // ﮍ '\ufb8d'
	case "rrparenthesis":
		return 0x2988 // ⦈ '\u2988'
	case "rrvocalicbengali":
		return 0x09e0 // ৠ '\u09e0'
	case "rrvocalicdeva":
		return 0x0960 // ॠ '\u0960'
	case "rrvocalicgujarati":
		return 0x0ae0 // ૠ '\u0ae0'
	case "rrvocalicvowelsignbengali":
		return 0x09c4 // ৄ '\u09c4'
	case "rrvocalicvowelsigndeva":
		return 0x0944 // ॄ '\u0944'
	case "rrvocalicvowelsigngujarati":
		return 0x0ac4 // ૄ '\u0ac4'
	case "rsolbar":
		return 0x29f7 // ⧷ '\u29f7'
	case "rsqhook":
		return 0x2ace // ⫎ '\u2ace'
	case "rsub":
		return 0x2a65 // ⩥ '\u2a65'
	case "rsuper":
		return 0x023c // ȼ '\u023c'
	case "rsuperior":
		return 0xf6f1 //  '\uf6f1'
	case "rtblock":
		return 0x2590 // ▐ '\u2590'
	case "rteighthblock":
		return 0x2595 // ▕ '\u2595'
	case "rtriltri":
		return 0x29ce // ⧎ '\u29ce'
	case "rturned":
		return 0x0279 // ɹ '\u0279'
	case "rturnedsuperior":
		return 0x02b4 // ʴ '\u02b4'
	case "rturnrthooksuper":
		return 0x023e // Ⱦ '\u023e'
	case "rturnsuper":
		return 0x023d // Ƚ '\u023d'
	case "ruhiragana":
		return 0x308b // る '\u308b'
	case "rukatakana":
		return 0x30eb // ル '\u30eb'
	case "rukatakanahalfwidth":
		return 0xff99 // ﾙ '\uff99'
	case "ruledelayed":
		return 0x29f4 // ⧴ '\u29f4'
	case "rupee":
		return 0x20a8 // ₨ '\u20a8'
	case "rupeemarkbengali":
		return 0x09f2 // ৲ '\u09f2'
	case "rupeesignbengali":
		return 0x09f3 // ৳ '\u09f3'
	case "rupiah":
		return 0xf6dd //  '\uf6dd'
	case "ruthai":
		return 0x0e24 // ฤ '\u0e24'
	case "rvboxline":
		return 0x23b9 // ⎹ '\u23b9'
	case "rvocalicbengali":
		return 0x098b // ঋ '\u098b'
	case "rvocalicdeva":
		return 0x090b // ऋ '\u090b'
	case "rvocalicgujarati":
		return 0x0a8b // ઋ '\u0a8b'
	case "rvocalicvowelsignbengali":
		return 0x09c3 // ৃ '\u09c3'
	case "rvocalicvowelsigndeva":
		return 0x0943 // ृ '\u0943'
	case "rvocalicvowelsigngujarati":
		return 0x0ac3 // ૃ '\u0ac3'
	case "rvzigzag":
		return 0x29d9 // ⧙ '\u29d9'
	case "s":
		return 0x0073 // s 's'
	case "sabengali":
		return 0x09b8 // স '\u09b8'
	case "sacute":
		return 0x015b // ś '\u015b'
	case "sacutedotaccent":
		return 0x1e65 // ṥ '\u1e65'
	case "sadeva":
		return 0x0938 // स '\u0938'
	case "sadfinalarabic":
		return 0xfeba // ﺺ '\ufeba'
	case "sadinitialarabic":
		return 0xfebb // ﺻ '\ufebb'
	case "sadisolated":
		return 0xfeb9 // ﺹ '\ufeb9'
	case "sadmedialarabic":
		return 0xfebc // ﺼ '\ufebc'
	case "sagujarati":
		return 0x0ab8 // સ '\u0ab8'
	case "sagurmukhi":
		return 0x0a38 // ਸ '\u0a38'
	case "sahiragana":
		return 0x3055 // さ '\u3055'
	case "sakatakana":
		return 0x30b5 // サ '\u30b5'
	case "sakatakanahalfwidth":
		return 0xff7b // ｻ '\uff7b'
	case "sallallahoualayhewasallamarabic":
		return 0xfdfa // ﷺ '\ufdfa'
	case "samekh":
		return 0x05e1 // ס '\u05e1'
	case "samekhdageshhebrew":
		return 0xfb41 // סּ '\ufb41'
	case "sansLmirrored":
		return 0x2143 // ⅃ '\u2143'
	case "sansLturned":
		return 0x2142 // ⅂ '\u2142'
	case "saraaathai":
		return 0x0e32 // า '\u0e32'
	case "saraaethai":
		return 0x0e41 // แ '\u0e41'
	case "saraaimaimalaithai":
		return 0x0e44 // ไ '\u0e44'
	case "saraaimaimuanthai":
		return 0x0e43 // ใ '\u0e43'
	case "saraamthai":
		return 0x0e33 // ำ '\u0e33'
	case "saraathai":
		return 0x0e30 // ะ '\u0e30'
	case "saraethai":
		return 0x0e40 // เ '\u0e40'
	case "saraiileftthai":
		return 0xf886 //  '\uf886'
	case "saraiithai":
		return 0x0e35 // ี '\u0e35'
	case "saraileftthai":
		return 0xf885 //  '\uf885'
	case "saraithai":
		return 0x0e34 // ิ '\u0e34'
	case "saraothai":
		return 0x0e42 // โ '\u0e42'
	case "saraueeleftthai":
		return 0xf888 //  '\uf888'
	case "saraueethai":
		return 0x0e37 // ื '\u0e37'
	case "saraueleftthai":
		return 0xf887 //  '\uf887'
	case "sarauethai":
		return 0x0e36 // ึ '\u0e36'
	case "sarauthai":
		return 0x0e38 // ุ '\u0e38'
	case "sarauuthai":
		return 0x0e39 // ู '\u0e39'
	case "satisfies":
		return 0x22a8 // ⊨ '\u22a8'
	case "sbopomofo":
		return 0x3119 // ㄙ '\u3119'
	case "scaron":
		return 0x0161 // š '\u0161'
	case "scarondotaccent":
		return 0x1e67 // ṧ '\u1e67'
	case "scedilla":
		return 0x015f // ş '\u015f'
	case "scedilla1":
		return 0xf817 //  '\uf817'
	case "schwa":
		return 0x0259 // ə '\u0259'
	case "schwacyrillic":
		return 0x04d9 // ә '\u04d9'
	case "schwadieresiscyrillic":
		return 0x04db // ӛ '\u04db'
	case "schwahook":
		return 0x025a // ɚ '\u025a'
	case "scircle":
		return 0x24e2 // ⓢ '\u24e2'
	case "scircumflex":
		return 0x015d // ŝ '\u015d'
	case "scommaaccent":
		return 0x0219 // ș '\u0219'
	case "scpolint":
		return 0x2a13 // ⨓ '\u2a13'
	case "scruple":
		return 0x2108 // ℈ '\u2108'
	case "scurel":
		return 0x22b1 // ⊱ '\u22b1'
	case "sdotaccent":
		return 0x1e61 // ṡ '\u1e61'
	case "sdotbelow":
		return 0x1e63 // ṣ '\u1e63'
	case "sdotbelowdotaccent":
		return 0x1e69 // ṩ '\u1e69'
	case "seagullbelowcmb":
		return 0x033c // ̼ '\u033c'
	case "seagullsubnosp":
		return 0x02af // ʯ '\u02af'
	case "second":
		return 0x2033 // ″ '\u2033'
	case "secondtonechinese":
		return 0x02ca // ˊ '\u02ca'
	case "section":
		return 0x00a7 // § '\u00a7'
	case "seenfinalarabic":
		return 0xfeb2 // ﺲ '\ufeb2'
	case "seeninitialarabic":
		return 0xfeb3 // ﺳ '\ufeb3'
	case "seenisolated":
		return 0xfeb1 // ﺱ '\ufeb1'
	case "seenmedialarabic":
		return 0xfeb4 // ﺴ '\ufeb4'
	case "seenwithmeeminitial":
		return 0xfcb0 // ﲰ '\ufcb0'
	case "segolhebrew":
		return 0x05b6 // ֶ '\u05b6'
	case "segoltahebrew":
		return 0x0592 // ֒ '\u0592'
	case "seharmenian":
		return 0x057d // ս '\u057d'
	case "sehiragana":
		return 0x305b // せ '\u305b'
	case "sekatakana":
		return 0x30bb // セ '\u30bb'
	case "sekatakanahalfwidth":
		return 0xff7e // ｾ '\uff7e'
	case "semicolon":
		return 0x003b // ; ';'
	case "semicolonmonospace":
		return 0xff1b // ； '\uff1b'
	case "semicolonsmall":
		return 0xfe54 // ﹔ '\ufe54'
	case "semivoicedmarkkana":
		return 0x309c // ゜ '\u309c'
	case "semivoicedmarkkanahalfwidth":
		return 0xff9f // ﾟ '\uff9f'
	case "sentisquare":
		return 0x3322 // ㌢ '\u3322'
	case "sentosquare":
		return 0x3323 // ㌣ '\u3323'
	case "seovnearrow":
		return 0x292d // ⤭ '\u292d'
	case "servicemark":
		return 0x2120 // ℠ '\u2120'
	case "setminus":
		return 0x29f5 // ⧵ '\u29f5'
	case "seven":
		return 0x0037 // 7 '7'
	case "sevenbengali":
		return 0x09ed // ৭ '\u09ed'
	case "sevendeva":
		return 0x096d // ७ '\u096d'
	case "seveneighths":
		return 0x215e // ⅞ '\u215e'
	case "sevengujarati":
		return 0x0aed // ૭ '\u0aed'
	case "sevengurmukhi":
		return 0x0a6d // ੭ '\u0a6d'
	case "sevenhangzhou":
		return 0x3027 // 〧 '\u3027'
	case "sevenideographicparen":
		return 0x3226 // ㈦ '\u3226'
	case "seveninferior":
		return 0x2087 // ₇ '\u2087'
	case "sevenmonospace":
		return 0xff17 // ７ '\uff17'
	case "sevenoldstyle":
		return 0xf737 //  '\uf737'
	case "sevenparen":
		return 0x247a // ⑺ '\u247a'
	case "sevenperiod":
		return 0x248e // ⒎ '\u248e'
	case "sevenpersian":
		return 0x06f7 // ۷ '\u06f7'
	case "sevenroman":
		return 0x2176 // ⅶ '\u2176'
	case "sevensuperior":
		return 0x2077 // ⁷ '\u2077'
	case "seventeencircle":
		return 0x2470 // ⑰ '\u2470'
	case "seventeenparen":
		return 0x2484 // ⒄ '\u2484'
	case "seventeenperiod":
		return 0x2498 // ⒘ '\u2498'
	case "seventhai":
		return 0x0e57 // ๗ '\u0e57'
	case "shaarmenian":
		return 0x0577 // շ '\u0577'
	case "shabengali":
		return 0x09b6 // শ '\u09b6'
	case "shaddaarabic":
		return 0x0651 // ّ '\u0651'
	case "shaddadammaarabic":
		return 0xfc61 // ﱡ '\ufc61'
	case "shaddadammatanarabic":
		return 0xfc5e // ﱞ '\ufc5e'
	case "shaddafathaarabic":
		return 0xfc60 // ﱠ '\ufc60'
	case "shaddahontatweel":
		return 0xfe7d // ﹽ '\ufe7d'
	case "shaddaisolated":
		return 0xfe7c // ﹼ '\ufe7c'
	case "shaddakasraarabic":
		return 0xfc62 // ﱢ '\ufc62'
	case "shaddakasratanarabic":
		return 0xfc5f // ﱟ '\ufc5f'
	case "shaddalow":
		return 0xe825 //  '\ue825'
	case "shaddawithdammaisolatedlow":
		return 0xe829 //  '\ue829'
	case "shaddawithdammamedial":
		return 0xfcf3 // ﳳ '\ufcf3'
	case "shaddawithdammatanisolatedlow":
		return 0xe82b //  '\ue82b'
	case "shaddawithfathalow":
		return 0xe828 //  '\ue828'
	case "shaddawithfathamedial":
		return 0xfcf2 // ﳲ '\ufcf2'
	case "shaddawithfathatanisolated":
		return 0xe818 //  '\ue818'
	case "shaddawithfathatanlow":
		return 0xe82a //  '\ue82a'
	case "shaddawithkasraisolatedlow":
		return 0xe82c //  '\ue82c'
	case "shaddawithkasramedial":
		return 0xfcf4 // ﳴ '\ufcf4'
	case "shaddawithkasratanisolatedlow":
		return 0xe82d //  '\ue82d'
	case "shade":
		return 0x2592 // ▒ '\u2592'
	case "shade1":
		return 0xf822 //  '\uf822'
	case "shadelight":
		return 0x2591 // ░ '\u2591'
	case "shadeva":
		return 0x0936 // श '\u0936'
	case "shagujarati":
		return 0x0ab6 // શ '\u0ab6'
	case "shagurmukhi":
		return 0x0a36 // ਸ਼ '\u0a36'
	case "shalshelethebrew":
		return 0x0593 // ֓ '\u0593'
	case "shbopomofo":
		return 0x3115 // ㄕ '\u3115'
	case "sheenfinalarabic":
		return 0xfeb6 // ﺶ '\ufeb6'
	case "sheeninitialarabic":
		return 0xfeb7 // ﺷ '\ufeb7'
	case "sheenisolated":
		return 0xfeb5 // ﺵ '\ufeb5'
	case "sheenmedialarabic":
		return 0xfeb8 // ﺸ '\ufeb8'
	case "sheenwithmeeminitial":
		return 0xfd30 // ﴰ '\ufd30'
	case "sheicoptic":
		return 0x03e3 // ϣ '\u03e3'
	case "shhacyrillic":
		return 0x04bb // һ '\u04bb'
	case "shiftleft":
		return 0x21b0 // ↰ '\u21b0'
	case "shiftright":
		return 0x21b1 // ↱ '\u21b1'
	case "shimacoptic":
		return 0x03ed // ϭ '\u03ed'
	case "shin":
		return 0x05e9 // ש '\u05e9'
	case "shindagesh":
		return 0xfb49 // שּ '\ufb49'
	case "shindageshshindot":
		return 0xfb2c // שּׁ '\ufb2c'
	case "shindageshsindothebrew":
		return 0xfb2d // שּׂ '\ufb2d'
	case "shindothebrew":
		return 0x05c1 // ׁ '\u05c1'
	case "shinshindot":
		return 0xfb2a // שׁ '\ufb2a'
	case "shook":
		return 0x0282 // ʂ '\u0282'
	case "shortdowntack":
		return 0x2adf // ⫟ '\u2adf'
	case "shortlefttack":
		return 0x2ade // ⫞ '\u2ade'
	case "shortrightarrowleftarrow":
		return 0x2944 // ⥄ '\u2944'
	case "shortuptack":
		return 0x2ae0 // ⫠ '\u2ae0'
	case "shuffle":
		return 0x29e2 // ⧢ '\u29e2'
	case "sigma":
		return 0x03c3 // σ '\u03c3'
	case "sigma1":
		return 0x03c2 // ς '\u03c2'
	case "sigmalunatesymbolgreek":
		return 0x03f2 // ϲ '\u03f2'
	case "sihiragana":
		return 0x3057 // し '\u3057'
	case "sikatakana":
		return 0x30b7 // シ '\u30b7'
	case "sikatakanahalfwidth":
		return 0xff7c // ｼ '\uff7c'
	case "siluqlefthebrew":
		return 0x05bd // ֽ '\u05bd'
	case "simgE":
		return 0x2aa0 // ⪠ '\u2aa0'
	case "simgtr":
		return 0x2a9e // ⪞ '\u2a9e'
	case "similar":
		return 0x223c // ∼ '\u223c'
	case "similarleftarrow":
		return 0x2b49 // ⭉ '\u2b49'
	case "similarrightarrow":
		return 0x2972 // ⥲ '\u2972'
	case "simlE":
		return 0x2a9f // ⪟ '\u2a9f'
	case "simless":
		return 0x2a9d // ⪝ '\u2a9d'
	case "simminussim":
		return 0x2a6c // ⩬ '\u2a6c'
	case "simneqq":
		return 0x2246 // ≆ '\u2246'
	case "simplus":
		return 0x2a24 // ⨤ '\u2a24'
	case "simrdots":
		return 0x2a6b // ⩫ '\u2a6b'
	case "sinewave":
		return 0x223f // ∿ '\u223f'
	case "siosacirclekorean":
		return 0x3274 // ㉴ '\u3274'
	case "siosaparenkorean":
		return 0x3214 // ㈔ '\u3214'
	case "sioscieuckorean":
		return 0x317e // ㅾ '\u317e'
	case "sioscirclekorean":
		return 0x3266 // ㉦ '\u3266'
	case "sioskiyeokkorean":
		return 0x317a // ㅺ '\u317a'
	case "sioskorean":
		return 0x3145 // ㅅ '\u3145'
	case "siosnieunkorean":
		return 0x317b // ㅻ '\u317b'
	case "siosparenkorean":
		return 0x3206 // ㈆ '\u3206'
	case "siospieupkorean":
		return 0x317d // ㅽ '\u317d'
	case "siostikeutkorean":
		return 0x317c // ㅼ '\u317c'
	case "six":
		return 0x0036 // 6 '6'
	case "sixbengali":
		return 0x09ec // ৬ '\u09ec'
	case "sixdeva":
		return 0x096c // ६ '\u096c'
	case "sixgujarati":
		return 0x0aec // ૬ '\u0aec'
	case "sixgurmukhi":
		return 0x0a6c // ੬ '\u0a6c'
	case "sixhangzhou":
		return 0x3026 // 〦 '\u3026'
	case "sixideographicparen":
		return 0x3225 // ㈥ '\u3225'
	case "sixinferior":
		return 0x2086 // ₆ '\u2086'
	case "sixmonospace":
		return 0xff16 // ６ '\uff16'
	case "sixoldstyle":
		return 0xf736 //  '\uf736'
	case "sixparen":
		return 0x2479 // ⑹ '\u2479'
	case "sixperemspace":
		return 0x2006 //  '\u2006'
	case "sixperiod":
		return 0x248d // ⒍ '\u248d'
	case "sixpersian":
		return 0x06f6 // ۶ '\u06f6'
	case "sixroman":
		return 0x2175 // ⅵ '\u2175'
	case "sixsuperior":
		return 0x2076 // ⁶ '\u2076'
	case "sixteencircle":
		return 0x246f // ⑯ '\u246f'
	case "sixteencurrencydenominatorbengali":
		return 0x09f9 // ৹ '\u09f9'
	case "sixteenparen":
		return 0x2483 // ⒃ '\u2483'
	case "sixteenperiod":
		return 0x2497 // ⒗ '\u2497'
	case "sixthai":
		return 0x0e56 // ๖ '\u0e56'
	case "slash":
		return 0x002f // / '/'
	case "slashlongnosp":
		return 0x02ab // ʫ '\u02ab'
	case "slashmonospace":
		return 0xff0f // ／ '\uff0f'
	case "slashshortnosp":
		return 0x02aa // ʪ '\u02aa'
	case "slongdotaccent":
		return 0x1e9b // ẛ '\u1e9b'
	case "slurabove":
		return 0x2322 // ⌢ '\u2322'
	case "smallblacktriangleleft":
		return 0x25c2 // ◂ '\u25c2'
	case "smallblacktriangleright":
		return 0x25b8 // ▸ '\u25b8'
	case "smallhighmadda":
		return 0x06e4 // ۤ '\u06e4'
	case "smallin":
		return 0x220a // ∊ '\u220a'
	case "smallni":
		return 0x220d // ∍ '\u220d'
	case "smashtimes":
		return 0x2a33 // ⨳ '\u2a33'
	case "smblkdiamond":
		return 0x2b29 // ⬩ '\u2b29'
	case "smblklozenge":
		return 0x2b2a // ⬪ '\u2b2a'
	case "smeparsl":
		return 0x29e4 // ⧤ '\u29e4'
	case "smile":
		return 0x2323 // ⌣ '\u2323'
	case "smileface":
		return 0x263a // ☺ '\u263a'
	case "smonospace":
		return 0xff53 // ｓ '\uff53'
	case "smt":
		return 0x2aaa // ⪪ '\u2aaa'
	case "smte":
		return 0x2aac // ⪬ '\u2aac'
	case "smwhitestar":
		return 0x2b52 // ⭒ '\u2b52'
	case "smwhtlozenge":
		return 0x2b2b // ⬫ '\u2b2b'
	case "sofpasuqhebrew":
		return 0x05c3 // ׃ '\u05c3'
	case "softhyphen":
		return 0x00ad //  '\u00ad'
	case "softsigncyrillic":
		return 0x044c // ь '\u044c'
	case "sohiragana":
		return 0x305d // そ '\u305d'
	case "sokatakana":
		return 0x30bd // ソ '\u30bd'
	case "sokatakanahalfwidth":
		return 0xff7f // ｿ '\uff7f'
	case "soliduslongoverlaycmb":
		return 0x0338 // ̸ '\u0338'
	case "solidusshortoverlaycmb":
		return 0x0337 // ̷ '\u0337'
	case "sorusithai":
		return 0x0e29 // ษ '\u0e29'
	case "sosalathai":
		return 0x0e28 // ศ '\u0e28'
	case "sosothai":
		return 0x0e0b // ซ '\u0e0b'
	case "sosuathai":
		return 0x0e2a // ส '\u0e2a'
	case "space":
		return 0x0020 //   ' '
	case "spade":
		return 0x2660 // ♠ '\u2660'
	case "spadesuitwhite":
		return 0x2664 // ♤ '\u2664'
	case "sparen":
		return 0x24ae // ⒮ '\u24ae'
	case "sphericalangle":
		return 0x2222 // ∢ '\u2222'
	case "sphericalangleup":
		return 0x29a1 // ⦡ '\u29a1'
	case "sqint":
		return 0x2a16 // ⨖ '\u2a16'
	case "sqlozenge":
		return 0x2311 // ⌑ '\u2311'
	case "sqrtbottom":
		return 0x23b7 // ⎷ '\u23b7'
	case "sqsubsetneq":
		return 0x22e4 // ⋤ '\u22e4'
	case "sqsupsetneq":
		return 0x22e5 // ⋥ '\u22e5'
	case "squarebelowcmb":
		return 0x033b // ̻ '\u033b'
	case "squarebotblack":
		return 0x2b13 // ⬓ '\u2b13'
	case "squarecc":
		return 0x33c4 // ㏄ '\u33c4'
	case "squarecm":
		return 0x339d // ㎝ '\u339d'
	case "squarediagonalcrosshatchfill":
		return 0x25a9 // ▩ '\u25a9'
	case "squaredot":
		return 0x22a1 // ⊡ '\u22a1'
	case "squarehorizontalfill":
		return 0x25a4 // ▤ '\u25a4'
	case "squareimage":
		return 0x228f // ⊏ '\u228f'
	case "squarekg":
		return 0x338f // ㎏ '\u338f'
	case "squarekm":
		return 0x339e // ㎞ '\u339e'
	case "squarekmcapital":
		return 0x33ce // ㏎ '\u33ce'
	case "squareleftblack":
		return 0x25e7 // ◧ '\u25e7'
	case "squarellblack":
		return 0x2b15 // ⬕ '\u2b15'
	case "squarellquad":
		return 0x25f1 // ◱ '\u25f1'
	case "squareln":
		return 0x33d1 // ㏑ '\u33d1'
	case "squarelog":
		return 0x33d2 // ㏒ '\u33d2'
	case "squarelrblack":
		return 0x25ea // ◪ '\u25ea'
	case "squarelrquad":
		return 0x25f2 // ◲ '\u25f2'
	case "squaremg":
		return 0x338e // ㎎ '\u338e'
	case "squaremil":
		return 0x33d5 // ㏕ '\u33d5'
	case "squareminus":
		return 0x229f // ⊟ '\u229f'
	case "squaremm":
		return 0x339c // ㎜ '\u339c'
	case "squaremsquared":
		return 0x33a1 // ㎡ '\u33a1'
	case "squaremultiply":
		return 0x22a0 // ⊠ '\u22a0'
	case "squareoriginal":
		return 0x2290 // ⊐ '\u2290'
	case "squareorthogonalcrosshatchfill":
		return 0x25a6 // ▦ '\u25a6'
	case "squareplus":
		return 0x229e // ⊞ '\u229e'
	case "squarerightblack":
		return 0x25e8 // ◨ '\u25e8'
	case "squaresubnosp":
		return 0x02ae // ʮ '\u02ae'
	case "squaretopblack":
		return 0x2b12 // ⬒ '\u2b12'
	case "squareulblack":
		return 0x25e9 // ◩ '\u25e9'
	case "squareulquad":
		return 0x25f0 // ◰ '\u25f0'
	case "squareupperlefttolowerrightfill":
		return 0x25a7 // ▧ '\u25a7'
	case "squareupperrighttolowerleftfill":
		return 0x25a8 // ▨ '\u25a8'
	case "squareurblack":
		return 0x2b14 // ⬔ '\u2b14'
	case "squareurquad":
		return 0x25f3 // ◳ '\u25f3'
	case "squareverticalfill":
		return 0x25a5 // ▥ '\u25a5'
	case "squarewhitewithsmallblack":
		return 0x25a3 // ▣ '\u25a3'
	case "squiggleleftright":
		return 0x21ad // ↭ '\u21ad'
	case "squiggleright":
		return 0x21dd // ⇝ '\u21dd'
	case "squoval":
		return 0x25a2 // ▢ '\u25a2'
	case "srsquare":
		return 0x33db // ㏛ '\u33db'
	case "ssabengali":
		return 0x09b7 // ষ '\u09b7'
	case "ssadeva":
		return 0x0937 // ष '\u0937'
	case "ssagujarati":
		return 0x0ab7 // ષ '\u0ab7'
	case "ssangcieuckorean":
		return 0x3149 // ㅉ '\u3149'
	case "ssanghieuhkorean":
		return 0x3185 // ㆅ '\u3185'
	case "ssangieungkorean":
		return 0x3180 // ㆀ '\u3180'
	case "ssangkiyeokkorean":
		return 0x3132 // ㄲ '\u3132'
	case "ssangnieunkorean":
		return 0x3165 // ㅥ '\u3165'
	case "ssangpieupkorean":
		return 0x3143 // ㅃ '\u3143'
	case "ssangsioskorean":
		return 0x3146 // ㅆ '\u3146'
	case "ssangtikeutkorean":
		return 0x3138 // ㄸ '\u3138'
	case "sslash":
		return 0x2afd // ⫽ '\u2afd'
	case "ssuperior":
		return 0xf6f2 //  '\uf6f2'
	case "st":
		return 0xfb06 // ﬆ '\ufb06'
	case "star":
		return 0x22c6 // ⋆ '\u22c6'
	case "stareq":
		return 0x225b // ≛ '\u225b'
	case "sterling":
		return 0x00a3 // £ '\u00a3'
	case "sterlingmonospace":
		return 0xffe1 // ￡ '\uffe1'
	case "strns":
		return 0x23e4 // ⏤ '\u23e4'
	case "strokelongoverlaycmb":
		return 0x0336 // ̶ '\u0336'
	case "strokeshortoverlaycmb":
		return 0x0335 // ̵ '\u0335'
	case "subedot":
		return 0x2ac3 // ⫃ '\u2ac3'
	case "submult":
		return 0x2ac1 // ⫁ '\u2ac1'
	case "subrarr":
		return 0x2979 // ⥹ '\u2979'
	case "subsetapprox":
		return 0x2ac9 // ⫉ '\u2ac9'
	case "subsetcirc":
		return 0x27c3 // ⟃ '\u27c3'
	case "subsetdbl":
		return 0x22d0 // ⋐ '\u22d0'
	case "subsetdblequal":
		return 0x2ac5 // ⫅ '\u2ac5'
	case "subsetdot":
		return 0x2abd // ⪽ '\u2abd'
	case "subsetnotequal":
		return 0x228a // ⊊ '\u228a'
	case "subsetornotdbleql":
		return 0x2acb // ⫋ '\u2acb'
	case "subsetplus":
		return 0x2abf // ⪿ '\u2abf'
	case "subsetsqequal":
		return 0x2291 // ⊑ '\u2291'
	case "subsim":
		return 0x2ac7 // ⫇ '\u2ac7'
	case "subsub":
		return 0x2ad5 // ⫕ '\u2ad5'
	case "subsup":
		return 0x2ad3 // ⫓ '\u2ad3'
	case "succapprox":
		return 0x2ab8 // ⪸ '\u2ab8'
	case "succeeds":
		return 0x227b // ≻ '\u227b'
	case "succeqq":
		return 0x2ab4 // ⪴ '\u2ab4'
	case "succneq":
		return 0x2ab2 // ⪲ '\u2ab2'
	case "suchthat":
		return 0x220b // ∋ '\u220b'
	case "suhiragana":
		return 0x3059 // す '\u3059'
	case "sukatakana":
		return 0x30b9 // ス '\u30b9'
	case "sukatakanahalfwidth":
		return 0xff7d // ｽ '\uff7d'
	case "sukunarabic":
		return 0x0652 // ْ '\u0652'
	case "sukunisolated":
		return 0xfe7e // ﹾ '\ufe7e'
	case "sukunlow":
		return 0xe822 //  '\ue822'
	case "sukunmedial":
		return 0xfe7f // ﹿ '\ufe7f'
	case "sukunonhamza":
		return 0xe834 //  '\ue834'
	case "sumbottom":
		return 0x23b3 // ⎳ '\u23b3'
	case "sumint":
		return 0x2a0b // ⨋ '\u2a0b'
	case "summation":
		return 0x2211 // ∑ '\u2211'
	case "sumtop":
		return 0x23b2 // ⎲ '\u23b2'
	case "sun":
		return 0x263c // ☼ '\u263c'
	case "supdsub":
		return 0x2ad8 // ⫘ '\u2ad8'
	case "supedot":
		return 0x2ac4 // ⫄ '\u2ac4'
	case "superscriptalef":
		return 0x0670 // ٰ '\u0670'
	case "supersetdbl":
		return 0x22d1 // ⋑ '\u22d1'
	case "supersetdblequal":
		return 0x2ac6 // ⫆ '\u2ac6'
	case "supersetnotequal":
		return 0x228b // ⊋ '\u228b'
	case "supersetornotdbleql":
		return 0x2acc // ⫌ '\u2acc'
	case "supersetsqequal":
		return 0x2292 // ⊒ '\u2292'
	case "suphsol":
		return 0x27c9 // ⟉ '\u27c9'
	case "suphsub":
		return 0x2ad7 // ⫗ '\u2ad7'
	case "suplarr":
		return 0x297b // ⥻ '\u297b'
	case "supmult":
		return 0x2ac2 // ⫂ '\u2ac2'
	case "supsetapprox":
		return 0x2aca // ⫊ '\u2aca'
	case "supsetcirc":
		return 0x27c4 // ⟄ '\u27c4'
	case "supsetdot":
		return 0x2abe // ⪾ '\u2abe'
	case "supsetplus":
		return 0x2ac0 // ⫀ '\u2ac0'
	case "supsim":
		return 0x2ac8 // ⫈ '\u2ac8'
	case "supsub":
		return 0x2ad4 // ⫔ '\u2ad4'
	case "supsup":
		return 0x2ad6 // ⫖ '\u2ad6'
	case "svsquare":
		return 0x33dc // ㏜ '\u33dc'
	case "syouwaerasquare":
		return 0x337c // ㍼ '\u337c'
	case "t":
		return 0x0074 // t 't'
	case "tabengali":
		return 0x09a4 // ত '\u09a4'
	case "tackdown":
		return 0x22a4 // ⊤ '\u22a4'
	case "tackleft":
		return 0x22a3 // ⊣ '\u22a3'
	case "tadeva":
		return 0x0924 // त '\u0924'
	case "tagujarati":
		return 0x0aa4 // ત '\u0aa4'
	case "tagurmukhi":
		return 0x0a24 // ਤ '\u0a24'
	case "taharabic":
		return 0x0637 // ط '\u0637'
	case "tahfinalarabic":
		return 0xfec2 // ﻂ '\ufec2'
	case "tahinitialarabic":
		return 0xfec3 // ﻃ '\ufec3'
	case "tahiragana":
		return 0x305f // た '\u305f'
	case "tahisolated":
		return 0xfec1 // ﻁ '\ufec1'
	case "tahmedialarabic":
		return 0xfec4 // ﻄ '\ufec4'
	case "taisyouerasquare":
		return 0x337d // ㍽ '\u337d'
	case "takatakana":
		return 0x30bf // タ '\u30bf'
	case "takatakanahalfwidth":
		return 0xff80 // ﾀ '\uff80'
	case "talloblong":
		return 0x2afe // ⫾ '\u2afe'
	case "tatweelwithfathatanabove":
		return 0xfe71 // ﹱ '\ufe71'
	case "tau":
		return 0x03c4 // τ '\u03c4'
	case "tavdagesh":
		return 0xfb4a // תּ '\ufb4a'
	case "tavhebrew":
		return 0x05ea // ת '\u05ea'
	case "tbar":
		return 0x0167 // ŧ '\u0167'
	case "tbopomofo":
		return 0x310a // ㄊ '\u310a'
	case "tcaron":
		return 0x0165 // ť '\u0165'
	case "tcaron1":
		return 0xf815 //  '\uf815'
	case "tccurl":
		return 0x02a8 // ʨ '\u02a8'
	case "tcedilla":
		return 0x0163 // ţ '\u0163'
	case "tcedilla1":
		return 0xf819 //  '\uf819'
	case "tcheharabic":
		return 0x0686 // چ '\u0686'
	case "tchehfinalarabic":
		return 0xfb7b // ﭻ '\ufb7b'
	case "tchehinitialarabic":
		return 0xfb7c // ﭼ '\ufb7c'
	case "tchehisolated":
		return 0xfb7a // ﭺ '\ufb7a'
	case "tchehmedialarabic":
		return 0xfb7d // ﭽ '\ufb7d'
	case "tcircle":
		return 0x24e3 // ⓣ '\u24e3'
	case "tcircumflexbelow":
		return 0x1e71 // ṱ '\u1e71'
	case "tdieresis":
		return 0x1e97 // ẗ '\u1e97'
	case "tdotaccent":
		return 0x1e6b // ṫ '\u1e6b'
	case "tdotbelow":
		return 0x1e6d // ṭ '\u1e6d'
	case "tedescendercyrillic":
		return 0x04ad // ҭ '\u04ad'
	case "tehfinalarabic":
		return 0xfe96 // ﺖ '\ufe96'
	case "tehhahinitialarabic":
		return 0xfca2 // ﲢ '\ufca2'
	case "tehhahisolatedarabic":
		return 0xfc0c // ﰌ '\ufc0c'
	case "tehinitialarabic":
		return 0xfe97 // ﺗ '\ufe97'
	case "tehiragana":
		return 0x3066 // て '\u3066'
	case "tehisolated":
		return 0xfe95 // ﺕ '\ufe95'
	case "tehjeeminitialarabic":
		return 0xfca1 // ﲡ '\ufca1'
	case "tehjeemisolatedarabic":
		return 0xfc0b // ﰋ '\ufc0b'
	case "tehmarbutaarabic":
		return 0x0629 // ة '\u0629'
	case "tehmarbutafinalarabic":
		return 0xfe94 // ﺔ '\ufe94'
	case "tehmarbutaisolated":
		return 0xfe93 // ﺓ '\ufe93'
	case "tehmedialarabic":
		return 0xfe98 // ﺘ '\ufe98'
	case "tehmeeminitialarabic":
		return 0xfca4 // ﲤ '\ufca4'
	case "tehmeemisolatedarabic":
		return 0xfc0e // ﰎ '\ufc0e'
	case "tehnoonfinalarabic":
		return 0xfc73 // ﱳ '\ufc73'
	case "tehwithalefmaksurafinal":
		return 0xfc74 // ﱴ '\ufc74'
	case "tehwithhehinitial":
		return 0xe814 //  '\ue814'
	case "tehwithkhahinitial":
		return 0xfca3 // ﲣ '\ufca3'
	case "tehwithyehfinal":
		return 0xfc75 // ﱵ '\ufc75'
	case "tehwithyehisolated":
		return 0xfc10 // ﰐ '\ufc10'
	case "tekatakana":
		return 0x30c6 // テ '\u30c6'
	case "tekatakanahalfwidth":
		return 0xff83 // ﾃ '\uff83'
	case "telephone":
		return 0x2121 // ℡ '\u2121'
	case "telishagedolahebrew":
		return 0x05a0 // ֠ '\u05a0'
	case "telishaqetanahebrew":
		return 0x05a9 // ֩ '\u05a9'
	case "tenideographicparen":
		return 0x3229 // ㈩ '\u3229'
	case "tenparen":
		return 0x247d // ⑽ '\u247d'
	case "tenperiod":
		return 0x2491 // ⒑ '\u2491'
	case "tenroman":
		return 0x2179 // ⅹ '\u2179'
	case "tesh":
		return 0x02a7 // ʧ '\u02a7'
	case "tetdagesh":
		return 0xfb38 // טּ '\ufb38'
	case "tethebrew":
		return 0x05d8 // ט '\u05d8'
	case "tetsecyrillic":
		return 0x04b5 // ҵ '\u04b5'
	case "tevirhebrew":
		return 0x059b // ֛ '\u059b'
	case "thabengali":
		return 0x09a5 // থ '\u09a5'
	case "thadeva":
		return 0x0925 // थ '\u0925'
	case "thagujarati":
		return 0x0aa5 // થ '\u0aa5'
	case "thagurmukhi":
		return 0x0a25 // ਥ '\u0a25'
	case "thalarabic":
		return 0x0630 // ذ '\u0630'
	case "thalfinalarabic":
		return 0xfeac // ﺬ '\ufeac'
	case "thalisolated":
		return 0xfeab // ﺫ '\ufeab'
	case "thanthakhatlowleftthai":
		return 0xf898 //  '\uf898'
	case "thanthakhatlowrightthai":
		return 0xf897 //  '\uf897'
	case "thanthakhatthai":
		return 0x0e4c // ์ '\u0e4c'
	case "thanthakhatupperleftthai":
		return 0xf896 //  '\uf896'
	case "theharabic":
		return 0x062b // ث '\u062b'
	case "thehfinalarabic":
		return 0xfe9a // ﺚ '\ufe9a'
	case "thehinitialarabic":
		return 0xfe9b // ﺛ '\ufe9b'
	case "thehisolated":
		return 0xfe99 // ﺙ '\ufe99'
	case "thehmedialarabic":
		return 0xfe9c // ﺜ '\ufe9c'
	case "thehwithmeeminitial":
		return 0xfca6 // ﲦ '\ufca6'
	case "thehwithmeemisolated":
		return 0xfc12 // ﰒ '\ufc12'
	case "therefore":
		return 0x2234 // ∴ '\u2234'
	case "thermod":
		return 0x29e7 // ⧧ '\u29e7'
	case "theta":
		return 0x03b8 // θ '\u03b8'
	case "theta1":
		return 0x03d1 // ϑ '\u03d1'
	case "thieuthacirclekorean":
		return 0x3279 // ㉹ '\u3279'
	case "thieuthaparenkorean":
		return 0x3219 // ㈙ '\u3219'
	case "thieuthcirclekorean":
		return 0x326b // ㉫ '\u326b'
	case "thieuthkorean":
		return 0x314c // ㅌ '\u314c'
	case "thieuthparenkorean":
		return 0x320b // ㈋ '\u320b'
	case "thinspace":
		return 0x2009 //  '\u2009'
	case "thirteencircle":
		return 0x246c // ⑬ '\u246c'
	case "thirteenparen":
		return 0x2480 // ⒀ '\u2480'
	case "thirteenperiod":
		return 0x2494 // ⒔ '\u2494'
	case "thonangmonthothai":
		return 0x0e11 // ฑ '\u0e11'
	case "thook":
		return 0x01ad // ƭ '\u01ad'
	case "thophuthaothai":
		return 0x0e12 // ฒ '\u0e12'
	case "thorn":
		return 0x00fe // þ '\u00fe'
	case "thothahanthai":
		return 0x0e17 // ท '\u0e17'
	case "thothanthai":
		return 0x0e10 // ฐ '\u0e10'
	case "thothongthai":
		return 0x0e18 // ธ '\u0e18'
	case "thothungthai":
		return 0x0e16 // ถ '\u0e16'
	case "thousandcyrillic":
		return 0x0482 // ҂ '\u0482'
	case "thousandsseparatorarabic":
		return 0x066c // ٬ '\u066c'
	case "three":
		return 0x0033 // 3 '3'
	case "threebengali":
		return 0x09e9 // ৩ '\u09e9'
	case "threedangle":
		return 0x27c0 // ⟀ '\u27c0'
	case "threedeva":
		return 0x0969 // ३ '\u0969'
	case "threedotcolon":
		return 0x2af6 // ⫶ '\u2af6'
	case "threeeighths":
		return 0x215c // ⅜ '\u215c'
	case "threefifths":
		return 0x2157 // ⅗ '\u2157'
	case "threegujarati":
		return 0x0ae9 // ૩ '\u0ae9'
	case "threegurmukhi":
		return 0x0a69 // ੩ '\u0a69'
	case "threehangzhou":
		return 0x3023 // 〣 '\u3023'
	case "threeideographicparen":
		return 0x3222 // ㈢ '\u3222'
	case "threeinferior":
		return 0x2083 // ₃ '\u2083'
	case "threemonospace":
		return 0xff13 // ３ '\uff13'
	case "threenumeratorbengali":
		return 0x09f6 // ৶ '\u09f6'
	case "threeoldstyle":
		return 0xf733 //  '\uf733'
	case "threeparen":
		return 0x2476 // ⑶ '\u2476'
	case "threeperemspace":
		return 0x2004 //  '\u2004'
	case "threeperiod":
		return 0x248a // ⒊ '\u248a'
	case "threepersian":
		return 0x06f3 // ۳ '\u06f3'
	case "threequarters":
		return 0x00be // ¾ '\u00be'
	case "threequartersemdash":
		return 0xf6de //  '\uf6de'
	case "threeroman":
		return 0x2172 // ⅲ '\u2172'
	case "threesuperior":
		return 0x00b3 // ³ '\u00b3'
	case "threethai":
		return 0x0e53 // ๓ '\u0e53'
	case "threeunderdot":
		return 0x20e8 // ⃨ '\u20e8'
	case "thzsquare":
		return 0x3394 // ㎔ '\u3394'
	case "tieconcat":
		return 0x2040 // ⁀ '\u2040'
	case "tieinfty":
		return 0x29dd // ⧝ '\u29dd'
	case "tihiragana":
		return 0x3061 // ち '\u3061'
	case "tikatakana":
		return 0x30c1 // チ '\u30c1'
	case "tikatakanahalfwidth":
		return 0xff81 // ﾁ '\uff81'
	case "tikeutacirclekorean":
		return 0x3270 // ㉰ '\u3270'
	case "tikeutaparenkorean":
		return 0x3210 // ㈐ '\u3210'
	case "tikeutcirclekorean":
		return 0x3262 // ㉢ '\u3262'
	case "tikeutkorean":
		return 0x3137 // ㄷ '\u3137'
	case "tikeutparenkorean":
		return 0x3202 // ㈂ '\u3202'
	case "tilde":
		return 0x02dc // ˜ '\u02dc'
	case "tilde1":
		return 0xf004 //  '\uf004'
	case "tildebelowcmb":
		return 0x0330 // ̰ '\u0330'
	case "tildecmb":
		return 0x0303 // ̃ '\u0303'
	case "tildedoublecmb":
		return 0x0360 // ͠ '\u0360'
	case "tildenosp":
		return 0x0276 // ɶ '\u0276'
	case "tildeoverlaycmb":
		return 0x0334 // ̴ '\u0334'
	case "tildeverticalcmb":
		return 0x033e // ̾ '\u033e'
	case "timesbar":
		return 0x2a31 // ⨱ '\u2a31'
	case "tipehahebrew":
		return 0x0596 // ֖ '\u0596'
	case "tippigurmukhi":
		return 0x0a70 // ੰ '\u0a70'
	case "titlocyrilliccmb":
		return 0x0483 // ҃ '\u0483'
	case "tiwnarmenian":
		return 0x057f // տ '\u057f'
	case "tlinebelow":
		return 0x1e6f // ṯ '\u1e6f'
	case "tminus":
		return 0x29ff // ⧿ '\u29ff'
	case "tmonospace":
		return 0xff54 // ｔ '\uff54'
	case "toarmenian":
		return 0x0569 // թ '\u0569'
	case "toea":
		return 0x2928 // ⤨ '\u2928'
	case "tohiragana":
		return 0x3068 // と '\u3068'
	case "tokatakana":
		return 0x30c8 // ト '\u30c8'
	case "tokatakanahalfwidth":
		return 0xff84 // ﾄ '\uff84'
	case "tona":
		return 0x2927 // ⤧ '\u2927'
	case "tonebarextrahighmod":
		return 0x02e5 // ˥ '\u02e5'
	case "tonebarextralowmod":
		return 0x02e9 // ˩ '\u02e9'
	case "tonebarhighmod":
		return 0x02e6 // ˦ '\u02e6'
	case "tonebarlowmod":
		return 0x02e8 // ˨ '\u02e8'
	case "tonebarmidmod":
		return 0x02e7 // ˧ '\u02e7'
	case "tonefive":
		return 0x01bd // ƽ '\u01bd'
	case "tonesix":
		return 0x0185 // ƅ '\u0185'
	case "tonetwo":
		return 0x01a8 // ƨ '\u01a8'
	case "tonos":
		return 0x0384 // ΄ '\u0384'
	case "tonsquare":
		return 0x3327 // ㌧ '\u3327'
	case "topatakthai":
		return 0x0e0f // ฏ '\u0e0f'
	case "topbot":
		return 0x2336 // ⌶ '\u2336'
	case "topcir":
		return 0x2af1 // ⫱ '\u2af1'
	case "topfork":
		return 0x2ada // ⫚ '\u2ada'
	case "topsemicircle":
		return 0x25e0 // ◠ '\u25e0'
	case "tortoiseshellbracketleft":
		return 0x3014 // 〔 '\u3014'
	case "tortoiseshellbracketleftsmall":
		return 0xfe5d // ﹝ '\ufe5d'
	case "tortoiseshellbracketleftvertical":
		return 0xfe39 // ︹ '\ufe39'
	case "tortoiseshellbracketright":
		return 0x3015 // 〕 '\u3015'
	case "tortoiseshellbracketrightsmall":
		return 0xfe5e // ﹞ '\ufe5e'
	case "tortoiseshellbracketrightvertical":
		return 0xfe3a // ︺ '\ufe3a'
	case "tosa":
		return 0x2929 // ⤩ '\u2929'
	case "totaothai":
		return 0x0e15 // ต '\u0e15'
	case "towa":
		return 0x292a // ⤪ '\u292a'
	case "tpalatalhook":
		return 0x01ab // ƫ '\u01ab'
	case "tparen":
		return 0x24af // ⒯ '\u24af'
	case "tplus":
		return 0x29fe // ⧾ '\u29fe'
	case "trademark":
		return 0x2122 // ™ '\u2122'
	case "trademarksans":
		return 0xf8ea //  '\uf8ea'
	case "trademarkserif":
		return 0xf6db //  '\uf6db'
	case "trapezium":
		return 0x23e2 // ⏢ '\u23e2'
	case "tretroflexhook":
		return 0x0288 // ʈ '\u0288'
	case "trianglebullet":
		return 0x2023 // ‣ '\u2023'
	case "trianglecdot":
		return 0x25ec // ◬ '\u25ec'
	case "triangleleftblack":
		return 0x25ed // ◭ '\u25ed'
	case "triangleleftequal":
		return 0x22b4 // ⊴ '\u22b4'
	case "triangleminus":
		return 0x2a3a // ⨺ '\u2a3a'
	case "triangleodot":
		return 0x29ca // ⧊ '\u29ca'
	case "triangleplus":
		return 0x2a39 // ⨹ '\u2a39'
	case "trianglerightblack":
		return 0x25ee // ◮ '\u25ee'
	case "trianglerightequal":
		return 0x22b5 // ⊵ '\u22b5'
	case "triangles":
		return 0x29cc // ⧌ '\u29cc'
	case "triangleserifs":
		return 0x29cd // ⧍ '\u29cd'
	case "triangletimes":
		return 0x2a3b // ⨻ '\u2a3b'
	case "triangleubar":
		return 0x29cb // ⧋ '\u29cb'
	case "tripleplus":
		return 0x29fb // ⧻ '\u29fb'
	case "trprime":
		return 0x2034 // ‴ '\u2034'
	case "trslash":
		return 0x2afb // ⫻ '\u2afb'
	case "ts":
		return 0x02a6 // ʦ '\u02a6'
	case "tsadidagesh":
		return 0xfb46 // צּ '\ufb46'
	case "tsecyrillic":
		return 0x0446 // ц '\u0446'
	case "tsere12":
		return 0x05b5 // ֵ '\u05b5'
	case "tshecyrillic":
		return 0x045b // ћ '\u045b'
	case "tsuperior":
		return 0xf6f3 //  '\uf6f3'
	case "ttabengali":
		return 0x099f // ট '\u099f'
	case "ttadeva":
		return 0x091f // ट '\u091f'
	case "ttagujarati":
		return 0x0a9f // ટ '\u0a9f'
	case "ttagurmukhi":
		return 0x0a1f // ਟ '\u0a1f'
	case "ttehfinalarabic":
		return 0xfb67 // ﭧ '\ufb67'
	case "ttehinitialarabic":
		return 0xfb68 // ﭨ '\ufb68'
	case "ttehmedialarabic":
		return 0xfb69 // ﭩ '\ufb69'
	case "tthabengali":
		return 0x09a0 // ঠ '\u09a0'
	case "tthadeva":
		return 0x0920 // ठ '\u0920'
	case "tthagujarati":
		return 0x0aa0 // ઠ '\u0aa0'
	case "tthagurmukhi":
		return 0x0a20 // ਠ '\u0a20'
	case "tturned":
		return 0x0287 // ʇ '\u0287'
	case "tuhiragana":
		return 0x3064 // つ '\u3064'
	case "tukatakana":
		return 0x30c4 // ツ '\u30c4'
	case "tukatakanahalfwidth":
		return 0xff82 // ﾂ '\uff82'
	case "turnangle":
		return 0x29a2 // ⦢ '\u29a2'
	case "turnediota":
		return 0x2129 // ℩ '\u2129'
	case "turnednot":
		return 0x2319 // ⌙ '\u2319'
	case "turnstileleft":
		return 0x22a2 // ⊢ '\u22a2'
	case "tusmallhiragana":
		return 0x3063 // っ '\u3063'
	case "tusmallkatakana":
		return 0x30c3 // ッ '\u30c3'
	case "tusmallkatakanahalfwidth":
		return 0xff6f // ｯ '\uff6f'
	case "twelvecircle":
		return 0x246b // ⑫ '\u246b'
	case "twelveparen":
		return 0x247f // ⑿ '\u247f'
	case "twelveperiod":
		return 0x2493 // ⒓ '\u2493'
	case "twelveroman":
		return 0x217b // ⅻ '\u217b'
	case "twelveudash":
		return 0xd80c //  '\ufffd'
	case "twentycircle":
		return 0x2473 // ⑳ '\u2473'
	case "twentyhangzhou":
		return 0x5344 // 卄 '\u5344'
	case "twentyparen":
		return 0x2487 // ⒇ '\u2487'
	case "twentyperiod":
		return 0x249b // ⒛ '\u249b'
	case "two":
		return 0x0032 // 2 '2'
	case "twoarabic":
		return 0x0662 // ٢ '\u0662'
	case "twobengali":
		return 0x09e8 // ২ '\u09e8'
	case "twocaps":
		return 0x2a4b // ⩋ '\u2a4b'
	case "twocups":
		return 0x2a4a // ⩊ '\u2a4a'
	case "twodeva":
		return 0x0968 // २ '\u0968'
	case "twodotleader":
		return 0x2025 // ‥ '\u2025'
	case "twodotleadervertical":
		return 0xfe30 // ︰ '\ufe30'
	case "twofifths":
		return 0x2156 // ⅖ '\u2156'
	case "twogujarati":
		return 0x0ae8 // ૨ '\u0ae8'
	case "twogurmukhi":
		return 0x0a68 // ੨ '\u0a68'
	case "twohangzhou":
		return 0x3022 // 〢 '\u3022'
	case "twoheaddownarrow":
		return 0x21a1 // ↡ '\u21a1'
	case "twoheadleftarrowtail":
		return 0x2b3b // ⬻ '\u2b3b'
	case "twoheadleftdbkarrow":
		return 0x2b37 // ⬷ '\u2b37'
	case "twoheadmapsfrom":
		return 0x2b36 // ⬶ '\u2b36'
	case "twoheadmapsto":
		return 0x2905 // ⤅ '\u2905'
	case "twoheadrightarrowtail":
		return 0x2916 // ⤖ '\u2916'
	case "twoheaduparrow":
		return 0x219f // ↟ '\u219f'
	case "twoheaduparrowcircle":
		return 0x2949 // ⥉ '\u2949'
	case "twoideographicparen":
		return 0x3221 // ㈡ '\u3221'
	case "twoinferior":
		return 0x2082 // ₂ '\u2082'
	case "twomonospace":
		return 0xff12 // ２ '\uff12'
	case "twonumeratorbengali":
		return 0x09f5 // ৵ '\u09f5'
	case "twooldstyle":
		return 0xf732 //  '\uf732'
	case "twoparen":
		return 0x2475 // ⑵ '\u2475'
	case "twoperiod":
		return 0x2489 // ⒉ '\u2489'
	case "twopersian":
		return 0x06f2 // ۲ '\u06f2'
	case "tworoman":
		return 0x2171 // ⅱ '\u2171'
	case "twostroke":
		return 0x01bb // ƻ '\u01bb'
	case "twosuperior":
		return 0x00b2 // ² '\u00b2'
	case "twothai":
		return 0x0e52 // ๒ '\u0e52'
	case "twothirds":
		return 0x2154 // ⅔ '\u2154'
	case "typecolon":
		return 0x2982 // ⦂ '\u2982'
	case "u":
		return 0x0075 // u 'u'
	case "u2643":
		return 0x2643 // ♃ '\u2643'
	case "uacute":
		return 0x00fa // ú '\u00fa'
	case "ubar":
		return 0x0289 // ʉ '\u0289'
	case "ubengali":
		return 0x0989 // উ '\u0989'
	case "ubopomofo":
		return 0x3128 // ㄨ '\u3128'
	case "ubrbrak":
		return 0x23e1 // ⏡ '\u23e1'
	case "ubreve":
		return 0x016d // ŭ '\u016d'
	case "ucaron":
		return 0x01d4 // ǔ '\u01d4'
	case "ucedilla":
		return 0xf834 //  '\uf834'
	case "ucircle":
		return 0x24e4 // ⓤ '\u24e4'
	case "ucircumflex":
		return 0x00fb // û '\u00fb'
	case "ucircumflexbelow":
		return 0x1e77 // ṷ '\u1e77'
	case "udattadeva":
		return 0x0951 // ॑ '\u0951'
	case "udblacute":
		return 0x0171 // ű '\u0171'
	case "udblgrave":
		return 0x0215 // ȕ '\u0215'
	case "udeva":
		return 0x0909 // उ '\u0909'
	case "udieresis":
		return 0x00fc // ü '\u00fc'
	case "udieresisacute":
		return 0x01d8 // ǘ '\u01d8'
	case "udieresisbelow":
		return 0x1e73 // ṳ '\u1e73'
	case "udieresiscaron":
		return 0x01da // ǚ '\u01da'
	case "udieresiscyrillic":
		return 0x04f1 // ӱ '\u04f1'
	case "udieresisgrave":
		return 0x01dc // ǜ '\u01dc'
	case "udieresismacron":
		return 0x01d6 // ǖ '\u01d6'
	case "udotbelow":
		return 0x1ee5 // ụ '\u1ee5'
	case "ugrave":
		return 0x00f9 // ù '\u00f9'
	case "ugujarati":
		return 0x0a89 // ઉ '\u0a89'
	case "ugurmukhi":
		return 0x0a09 // ਉ '\u0a09'
	case "uhiragana":
		return 0x3046 // う '\u3046'
	case "uhookabove":
		return 0x1ee7 // ủ '\u1ee7'
	case "uhorn":
		return 0x01b0 // ư '\u01b0'
	case "uhornacute":
		return 0x1ee9 // ứ '\u1ee9'
	case "uhorndotbelow":
		return 0x1ef1 // ự '\u1ef1'
	case "uhorngrave":
		return 0x1eeb // ừ '\u1eeb'
	case "uhornhookabove":
		return 0x1eed // ử '\u1eed'
	case "uhorntilde":
		return 0x1eef // ữ '\u1eef'
	case "uhungarumlautcyrillic":
		return 0x04f3 // ӳ '\u04f3'
	case "uinvertedbreve":
		return 0x0217 // ȗ '\u0217'
	case "ukatakana":
		return 0x30a6 // ウ '\u30a6'
	case "ukatakanahalfwidth":
		return 0xff73 // ｳ '\uff73'
	case "ukcyrillic":
		return 0x0479 // ѹ '\u0479'
	case "ukorean":
		return 0x315c // ㅜ '\u315c'
	case "ularc":
		return 0x25dc // ◜ '\u25dc'
	case "ultriangle":
		return 0x25f8 // ◸ '\u25f8'
	case "umacron":
		return 0x016b // ū '\u016b'
	case "umacroncyrillic":
		return 0x04ef // ӯ '\u04ef'
	case "umacrondieresis":
		return 0x1e7b // ṻ '\u1e7b'
	case "umatragurmukhi":
		return 0x0a41 // ੁ '\u0a41'
	case "uminus":
		return 0x2a41 // ⩁ '\u2a41'
	case "umonospace":
		return 0xff55 // ｕ '\uff55'
	case "underbrace":
		return 0x23df // ⏟ '\u23df'
	case "underbracket":
		return 0x23b5 // ⎵ '\u23b5'
	case "underleftarrow":
		return 0x20ee // ⃮ '\u20ee'
	case "underleftharpoondown":
		return 0x20ed // ⃭ '\u20ed'
	case "underparen":
		return 0x23dd // ⏝ '\u23dd'
	case "underrightarrow":
		return 0x20ef // ⃯ '\u20ef'
	case "underrightharpoondown":
		return 0x20ec // ⃬ '\u20ec'
	case "underscore":
		return 0x005f // _ '_'
	case "underscoredbl":
		return 0x2017 // ‗ '\u2017'
	case "underscoremonospace":
		return 0xff3f // ＿ '\uff3f'
	case "underscorevertical":
		return 0xfe33 // ︳ '\ufe33'
	case "underscorewavy":
		return 0xfe4f // ﹏ '\ufe4f'
	case "undertie":
		return 0x203f // ‿ '\u203f'
	case "unicodecdots":
		return 0x22ef // ⋯ '\u22ef'
	case "union":
		return 0x222a // ∪ '\u222a'
	case "uniondbl":
		return 0x22d3 // ⋓ '\u22d3'
	case "unionmulti":
		return 0x228e // ⊎ '\u228e'
	case "unionsq":
		return 0x2294 // ⊔ '\u2294'
	case "uniontext":
		return 0x22c3 // ⋃ '\u22c3'
	case "universal":
		return 0x2200 // ∀ '\u2200'
	case "uogonek":
		return 0x0173 // ų '\u0173'
	case "upand":
		return 0x214b // ⅋ '\u214b'
	case "uparen":
		return 0x24b0 // ⒰ '\u24b0'
	case "uparrowbarred":
		return 0x2909 // ⤉ '\u2909'
	case "uparrowoncircle":
		return 0x29bd // ⦽ '\u29bd'
	case "upblock":
		return 0x2580 // ▀ '\u2580'
	case "updigamma":
		return 0x03dd // ϝ '\u03dd'
	case "updownharpoonleftleft":
		return 0x2951 // ⥑ '\u2951'
	case "updownharpoonleftright":
		return 0x294d // ⥍ '\u294d'
	case "updownharpoonrightleft":
		return 0x294c // ⥌ '\u294c'
	case "updownharpoonrightright":
		return 0x294f // ⥏ '\u294f'
	case "updownharpoonsleftright":
		return 0x296e // ⥮ '\u296e'
	case "upeighthblock":
		return 0x2594 // ▔ '\u2594'
	case "upfishtail":
		return 0x297e // ⥾ '\u297e'
	case "upharpoonleftbar":
		return 0x2960 // ⥠ '\u2960'
	case "upharpoonrightbar":
		return 0x295c // ⥜ '\u295c'
	case "upharpoonsleftright":
		return 0x2963 // ⥣ '\u2963'
	case "upin":
		return 0x27d2 // ⟒ '\u27d2'
	case "upint":
		return 0x2a1b // ⨛ '\u2a1b'
	case "upkoppa":
		return 0x03df // ϟ '\u03df'
	case "upoldKoppa":
		return 0x03d8 // Ϙ '\u03d8'
	case "upoldkoppa":
		return 0x03d9 // ϙ '\u03d9'
	case "upperdothebrew":
		return 0x05c4 // ׄ '\u05c4'
	case "uprightcurvearrow":
		return 0x2934 // ⤴ '\u2934'
	case "upsampi":
		return 0x03e1 // ϡ '\u03e1'
	case "upsilon":
		return 0x03c5 // υ '\u03c5'
	case "upsilondiaeresistonos":
		return 0x02f9 // ˹ '\u02f9'
	case "upsilondieresis":
		return 0x03cb // ϋ '\u03cb'
	case "upsilondieresistonos":
		return 0x03b0 // ΰ '\u03b0'
	case "upsilonlatin":
		return 0x028a // ʊ '\u028a'
	case "upsilontonos":
		return 0x03cd // ύ '\u03cd'
	case "upslope":
		return 0x29f8 // ⧸ '\u29f8'
	case "upstigma":
		return 0x03db // ϛ '\u03db'
	case "uptackbelowcmb":
		return 0x031d // ̝ '\u031d'
	case "uptackmod":
		return 0x02d4 // ˔ '\u02d4'
	case "upvarTheta":
		return 0x03f4 // ϴ '\u03f4'
	case "uragurmukhi":
		return 0x0a73 // ੳ '\u0a73'
	case "urarc":
		return 0x25dd // ◝ '\u25dd'
	case "uring":
		return 0x016f // ů '\u016f'
	case "urtriangle":
		return 0x25f9 // ◹ '\u25f9'
	case "usmallhiragana":
		return 0x3045 // ぅ '\u3045'
	case "usmallkatakana":
		return 0x30a5 // ゥ '\u30a5'
	case "usmallkatakanahalfwidth":
		return 0xff69 // ｩ '\uff69'
	case "ustraightcyrillic":
		return 0x04af // ү '\u04af'
	case "ustraightstrokecyrillic":
		return 0x04b1 // ұ '\u04b1'
	case "utilde":
		return 0x0169 // ũ '\u0169'
	case "utildeacute":
		return 0x1e79 // ṹ '\u1e79'
	case "utildebelow":
		return 0x1e75 // ṵ '\u1e75'
	case "uubengali":
		return 0x098a // ঊ '\u098a'
	case "uudeva":
		return 0x090a // ऊ '\u090a'
	case "uugujarati":
		return 0x0a8a // ઊ '\u0a8a'
	case "uugurmukhi":
		return 0x0a0a // ਊ '\u0a0a'
	case "uumatragurmukhi":
		return 0x0a42 // ੂ '\u0a42'
	case "uuvowelsignbengali":
		return 0x09c2 // ূ '\u09c2'
	case "uuvowelsigndeva":
		return 0x0942 // ू '\u0942'
	case "uuvowelsigngujarati":
		return 0x0ac2 // ૂ '\u0ac2'
	case "uvowelsignbengali":
		return 0x09c1 // ু '\u09c1'
	case "uvowelsigndeva":
		return 0x0941 // ु '\u0941'
	case "uvowelsigngujarati":
		return 0x0ac1 // ુ '\u0ac1'
	case "v":
		return 0x0076 // v 'v'
	case "vBar":
		return 0x2ae8 // ⫨ '\u2ae8'
	case "vBarv":
		return 0x2ae9 // ⫩ '\u2ae9'
	case "vDdash":
		return 0x2ae2 // ⫢ '\u2ae2'
	case "vadeva":
		return 0x0935 // व '\u0935'
	case "vagujarati":
		return 0x0ab5 // વ '\u0ab5'
	case "vagurmukhi":
		return 0x0a35 // ਵ '\u0a35'
	case "vakatakana":
		return 0x30f7 // ヷ '\u30f7'
	case "varVdash":
		return 0x2ae6 // ⫦ '\u2ae6'
	case "varcarriagereturn":
		return 0x23ce // ⏎ '\u23ce'
	case "vardoublebarwedge":
		return 0x2306 // ⌆ '\u2306'
	case "varhexagon":
		return 0x2b21 // ⬡ '\u2b21'
	case "varhexagonblack":
		return 0x2b22 // ⬢ '\u2b22'
	case "varhexagonlrbonds":
		return 0x232c // ⌬ '\u232c'
	case "varika":
		return 0xfb1e // ﬞ '\ufb1e'
	case "varisinobar":
		return 0x22f6 // ⋶ '\u22f6'
	case "varisins":
		return 0x22f3 // ⋳ '\u22f3'
	case "varniobar":
		return 0x22fd // ⋽ '\u22fd'
	case "varnis":
		return 0x22fb // ⋻ '\u22fb'
	case "varointclockwise":
		return 0x2232 // ∲ '\u2232'
	case "vartriangleleft":
		return 0x22b2 // ⊲ '\u22b2'
	case "vartriangleright":
		return 0x22b3 // ⊳ '\u22b3'
	case "varveebar":
		return 0x2a61 // ⩡ '\u2a61'
	case "vav":
		return 0x05d5 // ו '\u05d5'
	case "vavdageshhebrew":
		return 0xfb35 // וּ '\ufb35'
	case "vavholam":
		return 0xfb4b // וֹ '\ufb4b'
	case "vbraceextender":
		return 0x23aa // ⎪ '\u23aa'
	case "vbrtri":
		return 0x29d0 // ⧐ '\u29d0'
	case "vcircle":
		return 0x24e5 // ⓥ '\u24e5'
	case "vdotbelow":
		return 0x1e7f // ṿ '\u1e7f'
	case "vectimes":
		return 0x2a2f // ⨯ '\u2a2f'
	case "vector":
		return 0x20d7 // ⃗ '\u20d7'
	case "veedot":
		return 0x27c7 // ⟇ '\u27c7'
	case "veedoublebar":
		return 0x2a63 // ⩣ '\u2a63'
	case "veeeq":
		return 0x225a // ≚ '\u225a'
	case "veemidvert":
		return 0x2a5b // ⩛ '\u2a5b'
	case "veeodot":
		return 0x2a52 // ⩒ '\u2a52'
	case "veeonvee":
		return 0x2a56 // ⩖ '\u2a56'
	case "veeonwedge":
		return 0x2a59 // ⩙ '\u2a59'
	case "veharabic":
		return 0x06a4 // ڤ '\u06a4'
	case "vehfinalarabic":
		return 0xfb6b // ﭫ '\ufb6b'
	case "vehinitialarabic":
		return 0xfb6c // ﭬ '\ufb6c'
	case "vehisolated":
		return 0xfb6a // ﭪ '\ufb6a'
	case "vehmedialarabic":
		return 0xfb6d // ﭭ '\ufb6d'
	case "vekatakana":
		return 0x30f9 // ヹ '\u30f9'
	case "versicle":
		return 0x2123 // ℣ '\u2123'
	case "verticallineabovecmb":
		return 0x030d // ̍ '\u030d'
	case "verticallinebelowcmb":
		return 0x0329 // ̩ '\u0329'
	case "verticallinelowmod":
		return 0x02cc // ˌ '\u02cc'
	case "verticallinemod":
		return 0x02c8 // ˈ '\u02c8'
	case "vertoverlay":
		return 0x20d2 // ⃒ '\u20d2'
	case "vewarmenian":
		return 0x057e // վ '\u057e'
	case "vhook":
		return 0x028b // ʋ '\u028b'
	case "viewdata":
		return 0x2317 // ⌗ '\u2317'
	case "vikatakana":
		return 0x30f8 // ヸ '\u30f8'
	case "viramabengali":
		return 0x09cd // ্ '\u09cd'
	case "viramadeva":
		return 0x094d // ् '\u094d'
	case "viramagujarati":
		return 0x0acd // ્ '\u0acd'
	case "visargabengali":
		return 0x0983 // ঃ '\u0983'
	case "visargadeva":
		return 0x0903 // ः '\u0903'
	case "visargagujarati":
		return 0x0a83 // ઃ '\u0a83'
	case "vlongdash":
		return 0x27dd // ⟝ '\u27dd'
	case "vmonospace":
		return 0xff56 // ｖ '\uff56'
	case "voarmenian":
		return 0x0578 // ո '\u0578'
	case "voicediterationhiragana":
		return 0x309e // ゞ '\u309e'
	case "voicediterationkatakana":
		return 0x30fe // ヾ '\u30fe'
	case "voicedmarkkana":
		return 0x309b // ゛ '\u309b'
	case "voicedmarkkanahalfwidth":
		return 0xff9e // ﾞ '\uff9e'
	case "vokatakana":
		return 0x30fa // ヺ '\u30fa'
	case "vparen":
		return 0x24b1 // ⒱ '\u24b1'
	case "vrectangle":
		return 0x25af // ▯ '\u25af'
	case "vrectangleblack":
		return 0x25ae // ▮ '\u25ae'
	case "vscript":
		return 0x021b // ț '\u021b'
	case "vtilde":
		return 0x1e7d // ṽ '\u1e7d'
	case "vturn":
		return 0x021c // Ȝ '\u021c'
	case "vturned":
		return 0x028c // ʌ '\u028c'
	case "vuhiragana":
		return 0x3094 // ゔ '\u3094'
	case "vukatakana":
		return 0x30f4 // ヴ '\u30f4'
	case "vysmblksquare":
		return 0x2b1d // ⬝ '\u2b1d'
	case "vysmwhtcircle":
		return 0x2218 // ∘ '\u2218'
	case "vysmwhtsquare":
		return 0x2b1e // ⬞ '\u2b1e'
	case "vzigzag":
		return 0x299a // ⦚ '\u299a'
	case "w":
		return 0x0077 // w 'w'
	case "wacute":
		return 0x1e83 // ẃ '\u1e83'
	case "waekorean":
		return 0x3159 // ㅙ '\u3159'
	case "wahiragana":
		return 0x308f // わ '\u308f'
	case "wakatakana":
		return 0x30ef // ワ '\u30ef'
	case "wakatakanahalfwidth":
		return 0xff9c // ﾜ '\uff9c'
	case "wakorean":
		return 0x3158 // ㅘ '\u3158'
	case "wasmallhiragana":
		return 0x308e // ゎ '\u308e'
	case "wasmallkatakana":
		return 0x30ee // ヮ '\u30ee'
	case "wattosquare":
		return 0x3357 // ㍗ '\u3357'
	case "wavedash":
		return 0x301c // 〜 '\u301c'
	case "wavyunderscorevertical":
		return 0xfe34 // ︴ '\ufe34'
	case "wawarabic":
		return 0x0648 // و '\u0648'
	case "wawfinalarabic":
		return 0xfeee // ﻮ '\ufeee'
	case "wawhamzaabovefinalarabic":
		return 0xfe86 // ﺆ '\ufe86'
	case "wawisolated":
		return 0xfeed // ﻭ '\ufeed'
	case "wawwithhamzaaboveisolated":
		return 0xfe85 // ﺅ '\ufe85'
	case "wbsquare":
		return 0x33dd // ㏝ '\u33dd'
	case "wcircle":
		return 0x24e6 // ⓦ '\u24e6'
	case "wcircumflex":
		return 0x0175 // ŵ '\u0175'
	case "wdieresis":
		return 0x1e85 // ẅ '\u1e85'
	case "wdotaccent":
		return 0x1e87 // ẇ '\u1e87'
	case "wdotbelow":
		return 0x1e89 // ẉ '\u1e89'
	case "wedgebar":
		return 0x2a5f // ⩟ '\u2a5f'
	case "wedgedot":
		return 0x27d1 // ⟑ '\u27d1'
	case "wedgedoublebar":
		return 0x2a60 // ⩠ '\u2a60'
	case "wedgemidvert":
		return 0x2a5a // ⩚ '\u2a5a'
	case "wedgeodot":
		return 0x2a51 // ⩑ '\u2a51'
	case "wedgeonwedge":
		return 0x2a55 // ⩕ '\u2a55'u3016
	case "wedgeq":
		return 0x2259 // ≙ '\u2259'
	case "wehiragana":
		return 0x3091 // ゑ '\u3091'
	case "weierstrass":
		return 0x2118 // ℘ '\u2118'
	case "wekatakana":
		return 0x30f1 // ヱ '\u30f1'
	case "wekorean":
		return 0x315e // ㅞ '\u315e'
	case "weokorean":
		return 0x315d // ㅝ '\u315d'
	case "wgrave":
		return 0x1e81 // ẁ '\u1e81'
	case "whitebullet":
		return 0x25e6 // ◦ '\u25e6'
	case "whitecircle":
		return 0x25cb // ○ '\u25cb'
	case "whitecornerbracketleft":
		return 0x300e // 『 '\u300e'
	case "whitecornerbracketleftvertical":
		return 0xfe43 // ﹃ '\ufe43'
	case "whitecornerbracketright":
		return 0x300f // 』 '\u300f'
	case "whitecornerbracketrightvertical":
		return 0xfe44 // ﹄ '\ufe44'
	case "whitediamond":
		return 0x25c7 // ◇ '\u25c7'
	case "whitediamondcontainingblacksmalldiamond":
		return 0x25c8 // ◈ '\u25c8'
	case "whitedownpointingsmalltriangle":
		return 0x25bf // ▿ '\u25bf'
	case "whitedownpointingtriangle":
		return 0x25bd // ▽ '\u25bd'
	case "whiteinwhitetriangle":
		return 0x27c1 // ⟁ '\u27c1'
	case "whiteleftpointingsmalltriangle":
		return 0x25c3 // ◃ '\u25c3'
	case "whiteleftpointingtriangle":
		return 0x25c1 // ◁ '\u25c1'
	case "whitelenticularbracketleft":
		return 0x3016 // 〖 '\u3016'
	case "whitelenticularbracketright":
		return 0x3017 // 〗 '\u3017'
	case "whitepointerleft":
		return 0x25c5 // ◅ '\u25c5'
	case "whitepointerright":
		return 0x25bb // ▻ '\u25bb'
	case "whiterightpointingsmalltriangle":
		return 0x25b9 // ▹ '\u25b9'
	case "whiterightpointingtriangle":
		return 0x25b7 // ▷ '\u25b7'
	case "whitesmallsquare":
		return 0x25ab // ▫ '\u25ab'
	case "whitesquaretickleft":
		return 0x27e4 // ⟤ '\u27e4'
	case "whitesquaretickright":
		return 0x27e5 // ⟥ '\u27e5'
	case "whitestar":
		return 0x2606 // ☆ '\u2606'
	case "whitetelephone":
		return 0x260f // ☏ '\u260f'
	case "whitetortoiseshellbracketleft":
		return 0x3018 // 〘 '\u3018'
	case "whitetortoiseshellbracketright":
		return 0x3019 // 〙 '\u3019'
	case "whiteuppointingsmalltriangle":
		return 0x25b5 // ▵ '\u25b5'
	case "whiteuppointingtriangle":
		return 0x25b3 // △ '\u25b3'
	case "whthorzoval":
		return 0x2b2d // ⬭ '\u2b2d'
	case "whtvertoval":
		return 0x2b2f // ⬯ '\u2b2f'
	case "wideangledown":
		return 0x29a6 // ⦦ '\u29a6'
	case "wideangleup":
		return 0x29a7 // ⦧ '\u29a7'
	case "widebridgeabove":
		return 0x20e9 // ⃩ '\u20e9'
	case "wihiragana":
		return 0x3090 // ゐ '\u3090'
	case "wikatakana":
		return 0x30f0 // ヰ '\u30f0'
	case "wikorean":
		return 0x315f // ㅟ '\u315f'
	case "wmonospace":
		return 0xff57 // ｗ '\uff57'
	case "wohiragana":
		return 0x3092 // を '\u3092'
	case "wokatakana":
		return 0x30f2 // ヲ '\u30f2'
	case "wokatakanahalfwidth":
		return 0xff66 // ｦ '\uff66'
	case "won":
		return 0x20a9 // ₩ '\u20a9'
	case "wonmonospace":
		return 0xffe6 // ￦ '\uffe6'
	case "wowaenthai":
		return 0x0e27 // ว '\u0e27'
	case "wparen":
		return 0x24b2 // ⒲ '\u24b2'
	case "wreathproduct":
		return 0x2240 // ≀ '\u2240'
	case "wring":
		return 0x1e98 // ẘ '\u1e98'
	case "wsuper":
		return 0x0240 // ɀ '\u0240'
	case "wsuperior":
		return 0x02b7 // ʷ '\u02b7'
	case "wturn":
		return 0x021d // ȝ '\u021d'
	case "wturned":
		return 0x028d // ʍ '\u028d'
	case "wynn":
		return 0x01bf // ƿ '\u01bf'
	case "x":
		return 0x0078 // x 'x'
	case "xabovecmb":
		return 0x033d // ̽ '\u033d'
	case "xbopomofo":
		return 0x3112 // ㄒ '\u3112'
	case "xcircle":
		return 0x24e7 // ⓧ '\u24e7'
	case "xdieresis":
		return 0x1e8d // ẍ '\u1e8d'
	case "xdotaccent":
		return 0x1e8b // ẋ '\u1e8b'
	case "xeharmenian":
		return 0x056d // խ '\u056d'
	case "xi":
		return 0x03be // ξ '\u03be'
	case "xmonospace":
		return 0xff58 // ｘ '\uff58'
	case "xparen":
		return 0x24b3 // ⒳ '\u24b3'
	case "xsuperior":
		return 0x02e3 // ˣ '\u02e3'
	case "y":
		return 0x0079 // y 'y'
	case "yaadosquare":
		return 0x334e // ㍎ '\u334e'
	case "yabengali":
		return 0x09af // য '\u09af'
	case "yacute":
		return 0x00fd // ý '\u00fd'
	case "yadeva":
		return 0x092f // य '\u092f'
	case "yaekorean":
		return 0x3152 // ㅒ '\u3152'
	case "yagujarati":
		return 0x0aaf // ય '\u0aaf'
	case "yagurmukhi":
		return 0x0a2f // ਯ '\u0a2f'
	case "yahiragana":
		return 0x3084 // や '\u3084'
	case "yakatakana":
		return 0x30e4 // ヤ '\u30e4'
	case "yakatakanahalfwidth":
		return 0xff94 // ﾔ '\uff94'
	case "yakorean":
		return 0x3151 // ㅑ '\u3151'
	case "yamakkanthai":
		return 0x0e4e // ๎ '\u0e4e'
	case "yasmallhiragana":
		return 0x3083 // ゃ '\u3083'
	case "yasmallkatakana":
		return 0x30e3 // ャ '\u30e3'
	case "yasmallkatakanahalfwidth":
		return 0xff6c // ｬ '\uff6c'
	case "yatcyrillic":
		return 0x0463 // ѣ '\u0463'
	case "ycircle":
		return 0x24e8 // ⓨ '\u24e8'
	case "ycircumflex":
		return 0x0177 // ŷ '\u0177'
	case "ydieresis":
		return 0x00ff // ÿ '\u00ff'
	case "ydotaccent":
		return 0x1e8f // ẏ '\u1e8f'
	case "ydotbelow":
		return 0x1ef5 // ỵ '\u1ef5'
	case "yeharabic":
		return 0x064a // ي '\u064a'
	case "yehbarreearabic":
		return 0x06d2 // ے '\u06d2'
	case "yehbarreefinalarabic":
		return 0xfbaf // ﮯ '\ufbaf'
	case "yehfinalarabic":
		return 0xfef2 // ﻲ '\ufef2'
	case "yehhamzaabovearabic":
		return 0x0626 // ئ '\u0626'
	case "yehhamzaabovefinalarabic":
		return 0xfe8a // ﺊ '\ufe8a'
	case "yehhamzaaboveinitialarabic":
		return 0xfe8b // ﺋ '\ufe8b'
	case "yehhamzaabovemedialarabic":
		return 0xfe8c // ﺌ '\ufe8c'
	case "yehinitialarabic":
		return 0xfef3 // ﻳ '\ufef3'
	case "yehisolated":
		return 0xfef1 // ﻱ '\ufef1'
	case "yehmeeminitialarabic":
		return 0xfcdd // ﳝ '\ufcdd'
	case "yehmeemisolatedarabic":
		return 0xfc58 // ﱘ '\ufc58'
	case "yehnoonfinalarabic":
		return 0xfc94 // ﲔ '\ufc94'
	case "yehthreedotsbelowarabic":
		return 0x06d1 // ۑ '\u06d1'
	case "yehwithalefmaksurafinal":
		return 0xfc95 // ﲕ '\ufc95'
	case "yehwithalefmaksuraisolated":
		return 0xfc59 // ﱙ '\ufc59'
	case "yehwithhahinitial":
		return 0xfcdb // ﳛ '\ufcdb'
	case "yehwithhamzaaboveisolated":
		return 0xfe89 // ﺉ '\ufe89'
	case "yehwithjeeminitial":
		return 0xfcda // ﳚ '\ufcda'
	case "yehwithkhahinitial":
		return 0xfcdc // ﳜ '\ufcdc'
	case "yehwithrehfinal":
		return 0xfc91 // ﲑ '\ufc91'
	case "yekorean":
		return 0x3156 // ㅖ '\u3156'
	case "yen":
		return 0x00a5 // ¥ '\u00a5'
	case "yenmonospace":
		return 0xffe5 // ￥ '\uffe5'
	case "yeokorean":
		return 0x3155 // ㅕ '\u3155'
	case "yeorinhieuhkorean":
		return 0x3186 // ㆆ '\u3186'
	case "yerahbenyomohebrew":
		return 0x05aa // ֪ '\u05aa'
	case "yericyrillic":
		return 0x044b // ы '\u044b'
	case "yerudieresiscyrillic":
		return 0x04f9 // ӹ '\u04f9'
	case "yesieungkorean":
		return 0x3181 // ㆁ '\u3181'
	case "yesieungpansioskorean":
		return 0x3183 // ㆃ '\u3183'
	case "yesieungsioskorean":
		return 0x3182 // ㆂ '\u3182'
	case "yetivhebrew":
		return 0x059a // ֚ '\u059a'
	case "ygrave":
		return 0x1ef3 // ỳ '\u1ef3'
	case "yhook":
		return 0x01b4 // ƴ '\u01b4'
	case "yhookabove":
		return 0x1ef7 // ỷ '\u1ef7'
	case "yiarmenian":
		return 0x0575 // յ '\u0575'
	case "yicyrillic":
		return 0x0457 // ї '\u0457'
	case "yikorean":
		return 0x3162 // ㅢ '\u3162'
	case "yinyang":
		return 0x262f // ☯ '\u262f'
	case "yiwnarmenian":
		return 0x0582 // ւ '\u0582'
	case "ymonospace":
		return 0xff59 // ｙ '\uff59'
	case "yoddageshhebrew":
		return 0xfb39 // יּ '\ufb39'
	case "yodyodhebrew":
		return 0x05f2 // ײ '\u05f2'
	case "yodyodpatahhebrew":
		return 0xfb1f // ײַ '\ufb1f'
	case "yogh":
		return 0x0222 // Ȣ '\u0222'
	case "yoghcurl":
		return 0x0223 // ȣ '\u0223'
	case "yohiragana":
		return 0x3088 // よ '\u3088'
	case "yoikorean":
		return 0x3189 // ㆉ '\u3189'
	case "yokatakana":
		return 0x30e8 // ヨ '\u30e8'
	case "yokatakanahalfwidth":
		return 0xff96 // ﾖ '\uff96'
	case "yokorean":
		return 0x315b // ㅛ '\u315b'
	case "yosmallhiragana":
		return 0x3087 // ょ '\u3087'
	case "yosmallkatakana":
		return 0x30e7 // ョ '\u30e7'
	case "yosmallkatakanahalfwidth":
		return 0xff6e // ｮ '\uff6e'
	case "yotgreek":
		return 0x03f3 // ϳ '\u03f3'
	case "yoyaekorean":
		return 0x3188 // ㆈ '\u3188'
	case "yoyakorean":
		return 0x3187 // ㆇ '\u3187'
	case "yoyakthai":
		return 0x0e22 // ย '\u0e22'
	case "yoyingthai":
		return 0x0e0d // ญ '\u0e0d'
	case "yparen":
		return 0x24b4 // ⒴ '\u24b4'
	case "ypogegrammeni":
		return 0x037a // ͺ '\u037a'
	case "ypogegrammenigreekcmb":
		return 0x0345 // ͅ '\u0345'
	case "yr":
		return 0x01a6 // Ʀ '\u01a6'
	case "yring":
		return 0x1e99 // ẙ '\u1e99'
	case "ysuper":
		return 0x0241 // Ɂ '\u0241'
	case "ysuperior":
		return 0x02b8 // ʸ '\u02b8'
	case "ytilde":
		return 0x1ef9 // ỹ '\u1ef9'
	case "yturn":
		return 0x021e // Ȟ '\u021e'
	case "yturned":
		return 0x028e // ʎ '\u028e'
	case "yuhiragana":
		return 0x3086 // ゆ '\u3086'
	case "yuikorean":
		return 0x318c // ㆌ '\u318c'
	case "yukatakana":
		return 0x30e6 // ユ '\u30e6'
	case "yukatakanahalfwidth":
		return 0xff95 // ﾕ '\uff95'
	case "yukorean":
		return 0x3160 // ㅠ '\u3160'
	case "yusbigcyrillic":
		return 0x046b // ѫ '\u046b'
	case "yusbigiotifiedcyrillic":
		return 0x046d // ѭ '\u046d'
	case "yuslittlecyrillic":
		return 0x0467 // ѧ '\u0467'
	case "yuslittleiotifiedcyrillic":
		return 0x0469 // ѩ '\u0469'
	case "yusmallhiragana":
		return 0x3085 // ゅ '\u3085'
	case "yusmallkatakana":
		return 0x30e5 // ュ '\u30e5'
	case "yusmallkatakanahalfwidth":
		return 0xff6d // ｭ '\uff6d'
	case "yuyekorean":
		return 0x318b // ㆋ '\u318b'
	case "yuyeokorean":
		return 0x318a // ㆊ '\u318a'
	case "yyabengali":
		return 0x09df // য় '\u09df'
	case "yyadeva":
		return 0x095f // य़ '\u095f'
	case "z":
		return 0x007a // z 'z'
	case "zaarmenian":
		return 0x0566 // զ '\u0566'
	case "zacute":
		return 0x017a // ź '\u017a'
	case "zadeva":
		return 0x095b // ज़ '\u095b'
	case "zagurmukhi":
		return 0x0a5b // ਜ਼ '\u0a5b'
	case "zaharabic":
		return 0x0638 // ظ '\u0638'
	case "zahfinalarabic":
		return 0xfec6 // ﻆ '\ufec6'
	case "zahinitialarabic":
		return 0xfec7 // ﻇ '\ufec7'
	case "zahiragana":
		return 0x3056 // ざ '\u3056'
	case "zahisolated":
		return 0xfec5 // ﻅ '\ufec5'
	case "zahmedialarabic":
		return 0xfec8 // ﻈ '\ufec8'
	case "zainarabic":
		return 0x0632 // ز '\u0632'
	case "zainfinalarabic":
		return 0xfeb0 // ﺰ '\ufeb0'
	case "zainisolated":
		return 0xfeaf // ﺯ '\ufeaf'
	case "zakatakana":
		return 0x30b6 // ザ '\u30b6'
	case "zaqefgadolhebrew":
		return 0x0595 // ֕ '\u0595'
	case "zaqefqatanhebrew":
		return 0x0594 // ֔ '\u0594'
	case "zarqahebrew":
		return 0x0598 // ֘ '\u0598'
	case "zayindageshhebrew":
		return 0xfb36 // זּ '\ufb36'
	case "zbopomofo":
		return 0x3117 // ㄗ '\u3117'
	case "zcaron":
		return 0x017e // ž '\u017e'
	case "zcircle":
		return 0x24e9 // ⓩ '\u24e9'
	case "zcircumflex":
		return 0x1e91 // ẑ '\u1e91'
	case "zcmp":
		return 0x2a1f // ⨟ '\u2a1f'
	case "zcurl":
		return 0x0291 // ʑ '\u0291'
	case "zdotaccent":
		return 0x017c // ż '\u017c'
	case "zdotbelow":
		return 0x1e93 // ẓ '\u1e93'
	case "zedescendercyrillic":
		return 0x0499 // ҙ '\u0499'
	case "zedieresiscyrillic":
		return 0x04df // ӟ '\u04df'
	case "zehiragana":
		return 0x305c // ぜ '\u305c'
	case "zekatakana":
		return 0x30bc // ゼ '\u30bc'
	case "zero":
		return 0x0030 // 0 '0'
	case "zerobengali":
		return 0x09e6 // ০ '\u09e6'
	case "zerodeva":
		return 0x0966 // ० '\u0966'
	case "zerogujarati":
		return 0x0ae6 // ૦ '\u0ae6'
	case "zerogurmukhi":
		return 0x0a66 // ੦ '\u0a66'
	case "zerohackarabic":
		return 0x0660 // ٠ '\u0660'
	case "zeroinferior":
		return 0x2080 // ₀ '\u2080'
	case "zeromonospace":
		return 0xff10 // ０ '\uff10'
	case "zerooldstyle":
		return 0xf730 //  '\uf730'
	case "zeropersian":
		return 0x06f0 // ۰ '\u06f0'
	case "zerosuperior":
		return 0x2070 // ⁰ '\u2070'
	case "zerothai":
		return 0x0e50 // ๐ '\u0e50'
	case "zerowidthjoiner":
		return 0xfeff //  '\ufeff'
	case "zerowidthspace":
		return 0x200b //  '\u200b'
	case "zeta":
		return 0x03b6 // ζ '\u03b6'
	case "zhbopomofo":
		return 0x3113 // ㄓ '\u3113'
	case "zhearmenian":
		return 0x056a // ժ '\u056a'
	case "zhebreve":
		return 0x03fe // Ͼ '\u03fe'
	case "zhebrevecyrillic":
		return 0x04c2 // ӂ '\u04c2'
	case "zhecyrillic":
		return 0x0436 // ж '\u0436'
	case "zhedescendercyrillic":
		return 0x0497 // җ '\u0497'
	case "zhedieresiscyrillic":
		return 0x04dd // ӝ '\u04dd'
	case "zihiragana":
		return 0x3058 // じ '\u3058'
	case "zikatakana":
		return 0x30b8 // ジ '\u30b8'
	case "zinorhebrew":
		return 0x05ae // ֮ '\u05ae'
	case "zlinebelow":
		return 0x1e95 // ẕ '\u1e95'
	case "zmonospace":
		return 0xff5a // ｚ '\uff5a'
	case "zohiragana":
		return 0x305e // ぞ '\u305e'
	case "zokatakana":
		return 0x30be // ゾ '\u30be'
	case "zparen":
		return 0x24b5 // ⒵ '\u24b5'
	case "zpipe":
		return 0x2a20 // ⨠ '\u2a20'
	case "zproject":
		return 0x2a21 // ⨡ '\u2a21'
	case "zretroflexhook":
		return 0x0290 // ʐ '\u0290'
	case "zrthook":
		return 0x0220 // Ƞ '\u0220'
	case "zstroke":
		return 0x01b6 // ƶ '\u01b6'
	case "zuhiragana":
		return 0x305a // ず '\u305a'
	case "zukatakana":
		return 0x30ba // ズ '\u30ba'
	}

	return 0
}
