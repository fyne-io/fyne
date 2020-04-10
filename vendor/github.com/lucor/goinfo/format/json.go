package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/lucor/goinfo"
)

// JSON writes reports in JSON format
type JSON struct{}

// Write writes the reports collected by reporters to out
func (w *JSON) Write(out io.Writer, reporters []goinfo.Reporter) error {
	reports := []report{}
	for _, reporter := range reporters {
		r, err := makeReport(reporter)
		if err != nil {
			return fmt.Errorf("[%s] could not collect info: %w", reporter.Summary(), err)
		}
		reports = append(reports, r)
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "\t")
	return enc.Encode(reports)
}
