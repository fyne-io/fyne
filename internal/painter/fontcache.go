//go:build !1.23

package painter

// ClearFontCache is used to remove cached fonts in the case that we wish to re-load Font faces
func ClearFontCache() {
	fontCache.Range(func(key, value interface{}) bool {
		fontCache.Delete(key)
		return true
	})
	fontCustomCache.Range(func(key, value interface{}) bool {
		fontCustomCache.Delete(key)
		return true
	})
}
