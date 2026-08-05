//go:build !glstats

package gl

// wrapContext is a no-op unless built with -tags glstats, in which case it
// returns a counting wrapper used to gather render statistics.
func wrapContext(c context) context {
	return c
}
