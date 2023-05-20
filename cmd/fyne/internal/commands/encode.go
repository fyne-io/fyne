package commands

import (
	"html"
)

func encodeXMLString(in string) string {
	return html.EscapeString(in)
}
