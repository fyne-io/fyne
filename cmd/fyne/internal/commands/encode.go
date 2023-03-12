package commands

import "strings"

func encodeXMLString(in string) string {
	amped := strings.ReplaceAll(in, "&", "&amp;")
	return strings.ReplaceAll(strings.ReplaceAll(amped, "<", "&lt;"), ">", "&gt;")
}
