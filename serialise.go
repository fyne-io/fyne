package fyne

import "bytes"
import "fmt"
import "io/ioutil"

func fromFile(name string) *Resource {
	data, err := ioutil.ReadFile(cachePath(name))

	if err != nil {
		return nil
	}

	return NewResource(name, data)
}

func toFile(res *Resource) {
	ioutil.WriteFile(cachePath(res.Name), res.Content, 0644)
}

// ToGo converts a Resource object to Go code.
// This is useful if serialising to a go file for compilation into a binary
func ToGo(res *Resource) string {
	var buffer bytes.Buffer

	buffer.WriteString("&fyne.Resource{\n")
	buffer.WriteString("\tName:    \"" + res.Name + "\",\n")
	buffer.WriteString("\tContent: []byte{")
	for i, v := range res.Content {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(fmt.Sprint(v))
	}
	buffer.WriteString("}}")

	return buffer.String()
}
