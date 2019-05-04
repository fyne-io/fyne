package fyne

import (
	"bytes"
	"fmt"
)

// GoString converts a Resource object to Go code.
// This is useful if serialising to a Go file for compilation into a binary.
func (r *StaticResource) GoString() string {
	var buffer bytes.Buffer

	buffer.WriteString("&fyne.StaticResource{\n")
	buffer.WriteString("\tStaticName: \"" + r.StaticName + "\",\n")
	buffer.WriteString("\tStaticContent: []byte{\n\t\t")
	for i, v := range r.StaticContent {
		if i > 0 {
			buffer.WriteString(", ")
		}

		buffer.WriteString(fmt.Sprint(v))
	}
	buffer.WriteString("}}")

	return buffer.String()
}
