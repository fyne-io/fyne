package binding

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func uriFromString(in string) (fyne.URI, error) {
	return storage.ParseURI(in)
}

func uriToString(in fyne.URI) (string, error) {
	if in == nil {
		return "", nil
	}

	return in.String(), nil
}

func parseBool(in string) (bool, error) {
	out, err := strconv.ParseBool(in)
	if err != nil {
		return false, err
	}

	return out, nil
}

func parseFloat(in string) (float64, error) {
	out, err := strconv.ParseFloat(in, 64)
	if err != nil {
		return 0, err
	}

	return out, nil
}

func parseInt(in string) (int, error) {
	out, err := strconv.ParseInt(in, 0, 64)
	if err != nil {
		return 0, err
	}
	return int(out), nil
}

func formatBool(in bool) string {
	return strconv.FormatBool(in)
}

func formatFloat(in float64) string {
	return strconv.FormatFloat(in, 'f', 6, 64)
}

func formatInt(in int) string {
	return strconv.FormatInt(int64(in), 10)
}
