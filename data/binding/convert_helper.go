package binding

import (
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func stripFormatPrecision(in string) string {
	// quick exit if certainly not float
	if !strings.ContainsAny(in, "f") {
		return in
	}

	start := -1
	end := -1
	runes := []rune(in)
	for i, r := range runes {
		switch r {
		case '%':
			if i > 0 && start == i-1 { // ignore %%
				start = -1
			} else {
				start = i
			}
		case 'f':
			if start == -1 { // not part of format
				continue
			}
			end = i
		}

		if end > -1 {
			break
		}
	}
	if end == start+1 { // no width/precision
		return in
	}

	sizeRunes := runes[start+1 : end]
	width, err := parseFloat(string(sizeRunes))
	if err != nil {
		return string(runes[:start+1]) + string(runes[:end])
	}

	if sizeRunes[0] == '.' { // formats like %.2f
		return string(runes[:start+1]) + string(runes[end:])
	}
	return string(runes[:start+1]) + strconv.Itoa(int(width)) + string(runes[end:])
}

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
