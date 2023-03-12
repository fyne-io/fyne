package commands

import "strings"

func encodeXMLString(in string) string {
	return strings.ReplaceAll(in, "&", "&amp;")
}
