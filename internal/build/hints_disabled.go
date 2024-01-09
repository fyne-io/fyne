//go:build !hints

package build

// HasHints is false to indicate that hints are not currently switched on.
// To enable please rebuild with "-tags hints" parameters.
const HasHints = false
