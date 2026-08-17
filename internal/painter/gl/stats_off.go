//go:build !glstats

package gl

// wrapContext is a no-op unless built with -tags glstats, in which case it
// returns a counting wrapper used to gather render statistics.
func wrapContext(c context) context {
	return c
}

// MarkPhase labels the stretch of a run being measured. It does nothing unless
// built with -tags glstats.
func MarkPhase(string) {}

// ReportStats writes the gathered render statistics and ends the process. It
// does nothing unless built with -tags glstats.
func ReportStats() {}
