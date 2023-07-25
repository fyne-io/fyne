package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeXMLString(t *testing.T) {
	assert.Equal(t, "", encodeXMLString(""))
	assert.Equal(t, "Hello", encodeXMLString("Hello"))
	assert.Equal(t, "Hi &amp; bye", encodeXMLString("Hi & bye"))
	assert.Equal(t, "Hi &amp;amp; bye", encodeXMLString("Hi &amp; bye"))

	assert.Equal(t, "Hi &lt;&gt; bye", encodeXMLString("Hi <> bye"))
}
