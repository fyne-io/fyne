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
func (r *Resource) ToGo() string {
	var buffer bytes.Buffer

	buffer.WriteString("&fyne.Resource{\n")
	buffer.WriteString("\tName: \"" + r.Name + "\",\n")
	buffer.WriteString("\tContent: []byte{")
	for i, v := range r.Content {
		if i == 0 {
			buffer.WriteString("\n")
		} else {
			buffer.WriteString(", ")
		}
		buffer.WriteString(fmt.Sprint(v))
	}
	buffer.WriteString("}}")

	return buffer.String()
}
