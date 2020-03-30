package format

import (
	"fmt"
	"io"
	"text/template"

	"github.com/lucor/goinfo"
)

// Text writes reports in text format
type Text struct{}

// Write writes the reports collected by reporters to out
func (w *Text) Write(out io.Writer, reporters []goinfo.Reporter) error {
	t := template.Must(template.New("text").Parse(textTpl))
	for _, reporter := range reporters {
		r, err := makeReport(reporter)
		if err != nil {
			return fmt.Errorf("[%s] could not collect info: %w", reporter.Summary(), err)
		}

		err = t.Execute(out, r)
		if err != nil {
			return fmt.Errorf("[%s] could not execute the text template: %w", reporter.Summary(), err)
		}
	}
	return nil
}

const textTpl = `## {{.Summary}}
{{ range $key, $value := .Info -}}
{{ $key }}="{{ $value }}"
{{ end }}
`
