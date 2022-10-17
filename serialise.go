package fyne

import (
	"fmt"
	"strings"
)

// GoString converts a Resource object to Go code.
// This is useful if serialising to a Go file for compilation into a binary.
func (r *StaticResource) GoString() string {
	buffer := strings.Builder{}

	buffer.WriteString("&fyne.StaticResource{\n\tStaticName: \"")
	buffer.WriteString(r.StaticName)
	buffer.WriteString("\",\n\tStaticContent: []byte{\n\t\t")
	for i, v := range r.StaticContent {
		if i > 0 {
			buffer.WriteString(", ")
		}

		fmt.Fprint(&buffer, v)
	}
	buffer.WriteString("}}")

	return buffer.String()
}
