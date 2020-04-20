package goinfo

import (
	"io"
)

// Info reprents the collected info
type Info map[string]interface{}

// Reporter is the interface that wraps the Summary and Info method methods
// along with the Errors interface
type Reporter interface {
	// Summary returns the summary's report
	Summary() string
	// Info returns the collected info
	Info() (Info, error)
}

// Formatter is the interface that wraps the Write method
type Formatter interface {
	// Write writes the reports collected by reporters to io.Writer
	Write(io.Writer, []Reporter) error
}

// Write the collected info by reporters to w using the specified format
func Write(w io.Writer, reporters []Reporter, format Formatter) error {
	return format.Write(w, reporters)
}
