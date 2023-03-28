package harfbuzz

import (
	"encoding/hex"
	"strings"

	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/opentype/loader"
	"github.com/go-text/typesetting/opentype/tables"
)

// ported from harfbuzz/src/hb-ot-tag.cc Copyright © 2009  Red Hat, Inc. 2011  Google, Inc. Behdad Esfahbod, Roozbeh Pournader

var (
	// OpenType script tag, `DFLT`, for features that are not script-specific.
	tagDefaultScript = loader.NewTag('D', 'F', 'L', 'T')
	// OpenType language tag, `dflt`. Not a valid language tag, but some fonts
	// mistakenly use it.
	tagDefaultLanguage = loader.NewTag('d', 'f', 'l', 't')
)

func oldTagFromScript(script language.Script) tables.Tag {
	/* This seems to be accurate as of end of 2012. */

	switch script {
	case 0:
		return tagDefaultScript

	/* KATAKANA and HIRAGANA both map to 'kana' */
	case language.Hiragana:
		return loader.NewTag('k', 'a', 'n', 'a')

	/* Spaces at the end are preserved, unlike ISO 15924 */
	case language.Lao:
		return loader.NewTag('l', 'a', 'o', ' ')
	case language.Yi:
		return loader.NewTag('y', 'i', ' ', ' ')
	/* Unicode-5.0 additions */
	case language.Nko:
		return loader.NewTag('n', 'k', 'o', ' ')
	/* Unicode-5.1 additions */
	case language.Vai:
		return loader.NewTag('v', 'a', 'i', ' ')
	}

	/* Else, just change first char to lowercase and return */
	return tables.Tag(script | 0x20000000)
}

//  static language.Script
//  hb_ot_old_tag_to_script (hb_tag_t tag)
//  {
//    if (unlikely (tag == HB_OT_TAG_DEFAULT_SCRIPT))
// 	 return HB_SCRIPT_INVALID;

//    /* This side of the conversion is fully algorithmic. */

//    /* Any spaces at the end of the tag are replaced by repeating the last
// 	* letter.  Eg 'nko ' -> 'Nkoo' */
//    if (unlikely ((tag & 0x0000FF00u) == 0x00002000u))
// 	 tag |= (tag >> 8) & 0x0000FF00u; /* Copy second letter to third */
//    if (unlikely ((tag & 0x000000FFu) == 0x00000020u))
// 	 tag |= (tag >> 8) & 0x000000FFu; /* Copy third letter to fourth */

//    /* Change first char to uppercase and return */
//    return (language.Script) (tag & ~0x20000000u);
//  }

func newTagFromScript(script language.Script) tables.Tag {
	switch script {
	case language.Bengali:
		return loader.NewTag('b', 'n', 'g', '2')
	case language.Devanagari:
		return loader.NewTag('d', 'e', 'v', '2')
	case language.Gujarati:
		return loader.NewTag('g', 'j', 'r', '2')
	case language.Gurmukhi:
		return loader.NewTag('g', 'u', 'r', '2')
	case language.Kannada:
		return loader.NewTag('k', 'n', 'd', '2')
	case language.Malayalam:
		return loader.NewTag('m', 'l', 'm', '2')
	case language.Oriya:
		return loader.NewTag('o', 'r', 'y', '2')
	case language.Tamil:
		return loader.NewTag('t', 'm', 'l', '2')
	case language.Telugu:
		return loader.NewTag('t', 'e', 'l', '2')
	case language.Myanmar:
		return loader.NewTag('m', 'y', 'm', '2')
	}

	return tagDefaultScript
}

//  static language.Script
//  hb_ot_new_tag_to_script (hb_tag_t tag)
//  {
//    switch (tag) {
// 	 case newTag('b','n','g','2'):	return HB_SCRIPT_BENGALI;
// 	 case newTag('d','e','v','2'):	return HB_SCRIPT_DEVANAGARI;
// 	 case newTag('g','j','r','2'):	return HB_SCRIPT_GUJARATI;
// 	 case newTag('g','u','r','2'):	return HB_SCRIPT_GURMUKHI;
// 	 case newTag('k','n','d','2'):	return HB_SCRIPT_KANNADA;
// 	 case newTag('m','l','m','2'):	return HB_SCRIPT_MALAYALAM;
// 	 case newTag('o','r','y','2'):	return HB_SCRIPT_ORIYA;
// 	 case newTag('t','m','l','2'):	return HB_SCRIPT_TAMIL;
// 	 case newTag('t','e','l','2'):	return HB_SCRIPT_TELUGU;
// 	 case newTag('m','y','m','2'):	return HB_SCRIPT_MYANMAR;
//    }

//    return HB_SCRIPT_UNKNOWN;
//  }

//  #ifndef HB_DISABLE_DEPRECATED
//  void
//  hb_ot_tags_from_script (language.Script  script,
// 			 hb_tag_t    *script_tag_1,
// 			 hb_tag_t    *script_tag_2)
//  {
//    unsigned int count = 2;
//    hb_tag_t tags[2];
//    otTagsFromScriptAndLanguage (script, HB_LANGUAGE_INVALID, &count, tags, nullptr, nullptr);
//    *script_tag_1 = count > 0 ? tags[0] : HB_OT_TAG_DEFAULT_SCRIPT;
//    *script_tag_2 = count > 1 ? tags[1] : HB_OT_TAG_DEFAULT_SCRIPT;
//  }
//  #endif

//  /*
//   * Complete list at:
//   * https://docs.microsoft.com/en-us/typography/opentype/spec/scripttags
//   *
//   * Most of the script tags are the same as the ISO 15924 tag but lowercased.
//   * So we just do that, and handle the exceptional cases in a switch.
//   */

func allTagsFromScript(script language.Script) []tables.Tag {
	var tags []tables.Tag

	tag := newTagFromScript(script)
	if tag != tagDefaultScript {
		// HB_SCRIPT_MYANMAR maps to 'mym2', but there is no 'mym3'.
		if tag != loader.NewTag('m', 'y', 'm', '2') {
			tags = append(tags, tag|'3')
		}
		tags = append(tags, tag)
	}

	oldTag := oldTagFromScript(script)
	if oldTag != tagDefaultScript {
		tags = append(tags, oldTag)
	}
	return tags
}

//  /**
//   * hb_ot_tag_to_script:
//   * @tag: a script tag
//   *
//   * Converts a script tag to an #language.Script.
//   *
//   * Return value: The #language.Script corresponding to @tag.
//   *
//   **/
//  language.Script
//  hb_ot_tag_to_script (hb_tag_t tag)
//  {
//    unsigned char digit = tag & 0x000000FFu;
//    if (unlikely (digit == '2' || digit == '3'))
// 	 return hb_ot_new_tag_to_script (tag & 0xFFFFFF32);

//    return hb_ot_old_tag_to_script (tag);
//  }

//  /* Language */

//  struct LangTag
//  {
//    char language[4];
//    hb_tag_t tag;

//    int cmp (const char *a) const
//    {
// 	 const char *b = this->language;
// 	 unsigned int da, db;
// 	 const char *p;

// 	 p = strchr (a, '-');
// 	 da = p ? (unsigned int) (p - a) : strlen (a);

// 	 p = strchr (b, '-');
// 	 db = p ? (unsigned int) (p - b) : strlen (b);

// 	 return strncmp (a, b, max (da, db));
//    }
//    int cmp (const LangTag *that) const
//    { return cmp (that->language); }
//  };

//  #include "hb-ot-tag-table.hh"

//  /* The corresponding languages IDs for the following IDs are unclear,
//   * overlap, or are architecturally weird. Needs more research. */

//  /*{"??",	{newTag('B','C','R',' ')}},*/	/* Bible Cree */
//  /*{"zh?",	{newTag('C','H','N',' ')}},*/	/* Chinese (seen in Microsoft fonts) */
//  /*{"ar-Syrc?",	{newTag('G','A','R',' ')}},*/	/* Garshuni */
//  /*{"??",	{newTag('N','G','R',' ')}},*/	/* Nagari */
//  /*{"??",	{newTag('Y','I','C',' ')}},*/	/* Yi Classic */
//  /*{"zh?",	{newTag('Z','H','P',' ')}},*/	/* Chinese Phonetic */

//  #ifndef HB_DISABLE_DEPRECATED
//  hb_tag_t
//  hb_ot_tag_from_language (Language language)
//  {
//    unsigned int count = 1;
//    hb_tag_t tags[1];
//    otTagsFromScriptAndLanguage (HB_SCRIPT_UNKNOWN, language, nullptr, nullptr, &count, tags);
//    return count > 0 ? tags[0] : HB_OT_TAG_DEFAULT_LANGUAGE;
//  }
//  #endif

func otTagsFromLanguage(langStr string, limit int) []tables.Tag {
	// check for matches of multiple subtags.
	if tags := tagsFromComplexLanguage(langStr, limit); len(tags) != 0 {
		return tags
	}

	// find a language matching in the first component.
	s := strings.IndexByte(langStr, '-')
	if s != -1 && limit >= 6 {
		extlangEnd := strings.IndexByte(langStr[s+1:], '-')
		// if there is an extended language tag, use it.
		ref := extlangEnd
		if extlangEnd == -1 {
			ref = len(langStr[s+1:])
		}
		if ref == 3 && isAlpha(langStr[s+1]) {
			langStr = langStr[s+1:]
		}
	}

	if tagIdx := bfindLanguage(langStr); tagIdx != -1 {
		for tagIdx != 0 && otLanguages[tagIdx].language == otLanguages[tagIdx-1].language {
			tagIdx--
		}
		var out []tables.Tag
		for i := 0; tagIdx+i < len(otLanguages) &&
			otLanguages[tagIdx+i].tag != 0 &&
			otLanguages[tagIdx+i].language == otLanguages[tagIdx].language; i++ {
			out = append(out, otLanguages[tagIdx+i].tag)
		}
		return out
	}

	if s == -1 {
		s = len(langStr)
	}
	if s == 3 {
		// assume it's ISO-639-3 and upper-case and use it.
		return []tables.Tag{loader.NewTag(langStr[0], langStr[1], langStr[2], ' ') & ^tables.Tag(0x20202000)}
	}

	return nil
}

// return 0 if no tag
func parsePrivateUseSubtag(privateUseSubtag string, prefix string, normalize func(byte) byte) (tables.Tag, bool) {
	s := strings.Index(privateUseSubtag, prefix)
	if s == -1 {
		return 0, false
	}

	var tag [4]byte
	L := len(privateUseSubtag)
	s += len(prefix)
	if s < L && privateUseSubtag[s] == '-' {
		s += 1
		if L < s+8 {
			return 0, false
		}
		_, err := hex.Decode(tag[:], []byte(privateUseSubtag[s:s+8]))
		if err != nil {
			return 0, false
		}
	} else {
		var i int
		for ; i < 4 && s+i < L && isAlnum(privateUseSubtag[s+i]); i++ {
			tag[i] = normalize(privateUseSubtag[s+i])
		}
		if i == 0 {
			return 0, false
		}

		for ; i < 4; i++ {
			tag[i] = ' '
		}
	}
	out := loader.NewTag(tag[0], tag[1], tag[2], tag[3])
	if (out & 0xDFDFDFDF) == tagDefaultScript {
		out ^= ^tables.Tag(0xDFDFDFDF)
	}
	return out, true
}

// NewOTTagsFromScriptAndLanguage converts an `language.Script` and an `Language`
// to script and language tags.
func NewOTTagsFromScriptAndLanguage(script language.Script, language language.Language) (scriptTags, languageTags []tables.Tag) {
	if language != "" {
		langStr := languageToString(language)
		limit := -1
		privateUseSubtag := ""
		if langStr[0] == 'x' && langStr[1] == '-' {
			privateUseSubtag = langStr
		} else {
			var s int
			for s = 1; s < len(langStr); s++ { // s index in lang_str
				if langStr[s-1] == '-' && langStr[s+1] == '-' {
					if langStr[s] == 'x' {
						privateUseSubtag = langStr[s:]
						if limit == -1 {
							limit = s - 1
						}
						break
					} else if limit == -1 {
						limit = s - 1
					}
				}
			}
			if limit == -1 {
				limit = s
			}
		}

		s, hasScript := parsePrivateUseSubtag(privateUseSubtag, "-hbsc", toLower)
		if hasScript {
			scriptTags = []tables.Tag{s}
		}

		l, hasLanguage := parsePrivateUseSubtag(privateUseSubtag, "-hbot", toUpper)
		if hasLanguage {
			languageTags = append(languageTags, l)
		} else {
			languageTags = otTagsFromLanguage(langStr, limit)
		}
	}

	if len(scriptTags) == 0 {
		scriptTags = allTagsFromScript(script)
	}
	return
}

//  /**
//   * hb_ot_tag_to_language:
//   * @tag: an language tag
//   *
//   * Converts a language tag to an #Language.
//   *
//   * Return value: (transfer none) (nullable):
//   * The #Language corresponding to @tag.
//   *
//   * Since: 0.9.2
//   **/
//  Language
//  hb_ot_tag_to_language (hb_tag_t tag)
//  {
//    unsigned int i;

//    if (tag == HB_OT_TAG_DEFAULT_LANGUAGE)
// 	 return nullptr;

//    {
// 	 Language disambiguated_tag = ambiguousTagToLanguage (tag);
// 	 if (disambiguated_tag != HB_LANGUAGE_INVALID)
// 	   return disambiguated_tag;
//    }

//    for (i = 0; i < ARRAY_LENGTH (ot_languages); i++)
// 	 if (ot_languages[i].tag == tag)
// 	   return hb_language_from_string (ot_languages[i].language, -1);

//    /* Return a custom language in the form of "x-hbot-AABBCCDD".
// 	* If it's three letters long, also guess it's ISO 639-3 and lower-case and
// 	* prepend it (if it's not a registered tag, the private use subtags will
// 	* ensure that calling hb_ot_tag_from_language on the result will still return
// 	* the same tag as the original tag).
// 	*/
//    {
// 	 char buf[20];
// 	 char *str = buf;
// 	 if (isAlpha (tag >> 24)
// 	 && isAlpha ((tag >> 16) & 0xFF)
// 	 && isAlpha ((tag >> 8) & 0xFF)
// 	 && (tag & 0xFF) == ' ')
// 	 {
// 	   buf[0] = TOLOWER (tag >> 24);
// 	   buf[1] = TOLOWER ((tag >> 16) & 0xFF);
// 	   buf[2] = TOLOWER ((tag >> 8) & 0xFF);
// 	   buf[3] = '-';
// 	   str += 4;
// 	 }
// 	 snprintf (str, 16, "x-hbot-%08x", tag);
// 	 return hb_language_from_string (&*buf, -1);
//    }
//  }

//  /**
//   * hb_ot_tags_to_script_and_language:
//   * @script_tag: a script tag
//   * @language_tag: a language tag
//   * @script: (out) (optional): the #language.Script corresponding to @script_tag.
//   * @language: (out) (optional): the #Language corresponding to @script_tag and
//   * @language_tag.
//   *
//   * Converts a script tag and a language tag to an #language.Script and an
//   * #Language.
//   *
//   * Since: 2.0.0
//   **/
//  void
//  hb_ot_tags_to_script_and_language (hb_tag_t       script_tag,
// 					hb_tag_t       language_tag,
// 					language.Script   *script /* OUT */,
// 					Language *language /* OUT */)
//  {
//    language.Script script_out = hb_ot_tag_to_script (script_tag);
//    if (script)
// 	 *script = script_out;
//    if (language)
//    {
// 	 unsigned int script_count = 1;
// 	 hb_tag_t primary_script_tag[1];
// 	 otTagsFromScriptAndLanguage (script_out,
// 					  HB_LANGUAGE_INVALID,
// 					  &script_count,
// 					  primary_script_tag,
// 					  nullptr, nullptr);
// 	 *language = hb_ot_tag_to_language (language_tag);
// 	 if (script_count == 0 || primary_script_tag[0] != script_tag)
// 	 {
// 	   unsigned char *buf;
// 	   const char *lang_str = languageToString (*language);
// 	   size_t len = strlen (lang_str);
// 	   buf = (unsigned char *) malloc (len + 16);
// 	   if (unlikely (!buf))
// 	   {
// 	 *language = nullptr;
// 	   }
// 	   else
// 	   {
// 	 int shift;
// 	 memcpy (buf, lang_str, len);
// 	 if (lang_str[0] != 'x' || lang_str[1] != '-') {
// 	   buf[len++] = '-';
// 	   buf[len++] = 'x';
// 	 }
// 	 buf[len++] = '-';
// 	 buf[len++] = 'h';
// 	 buf[len++] = 'b';
// 	 buf[len++] = 's';
// 	 buf[len++] = 'c';
// 	 buf[len++] = '-';
// 	 for (shift = 28; shift >= 0; shift -= 4)
// 	   buf[len++] = TOHEX (script_tag >> shift);
// 	 *language = hb_language_from_string ((char *) buf, len);
// 	 free (buf);
// 	   }
// 	 }
//    }
//  }
