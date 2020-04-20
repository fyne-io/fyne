package format

import "github.com/lucor/goinfo"

// report reprents a report
type report struct {
	Summary string      `json:"summary"`
	Info    goinfo.Info `json:"info,omitempty"`
}

// makeReport makes a report from a reporter
func makeReport(reporter goinfo.Reporter) (report, error) {
	info, err := reporter.Info()
	if err != nil {
		return report{}, err
	}

	return report{
		Summary: reporter.Summary(),
		Info:    info,
	}, nil
}
