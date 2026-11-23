package painter

import "github.com/go-text/typesetting/font"

func FontCacheLen() (n int) {
	fontCache.Range(func(cacheID, *FontCacheItem) bool { n++; return true })
	return n
}

func ParsedFontsLen() (n int) {
	parsedFonts.Range(func(parsedFontKey, *font.Font) bool { n++; return true })
	return n
}
