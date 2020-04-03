package format

import (
	"fmt"
	"html/template"
	"io"

	"github.com/lucor/goinfo"
)

// HTMLDetails writes reports in HTML Details format
type HTMLDetails struct{}

// Write writes the reports collected by reporters to out
func (w *HTMLDetails) Write(out io.Writer, reporters []goinfo.Reporter) error {
	t := template.Must(template.New("html_details").Parse(htmlDetailsTpl))
	for _, reporter := range reporters {
		r, err := makeReport(reporter)
		if err != nil {
			return fmt.Errorf("[%s] could not collect info: %w", reporter.Summary(), err)
		}

		err = t.Execute(out, r)
		if err != nil {
			return fmt.Errorf("[%s] could not execute the html details template: %w", reporter.Summary(), err)
		}
	}
	return nil
}

const htmlDetailsTpl = `
<details><summary>{{.Summary}}</summary><br><pre>
{{ range $key, $value := .Info -}}
{{ $key }}={{ $value }}
{{ end -}}
</pre></details>
`
