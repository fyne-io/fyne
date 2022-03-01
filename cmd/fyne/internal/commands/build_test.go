package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildGenerateMetaLDFlags(t *testing.T) {
	b := &builder{}
	assert.Equal(t, "", b.generateMetaLDFlags())

	b.id = "com.example"
	assert.Equal(t, "-X 'fyne.io/fyne/v2/internal/app.MetaID=com.example'", b.generateMetaLDFlags())

	b.version = "1.2.3"
	assert.Equal(t, "-X 'fyne.io/fyne/v2/internal/app.MetaID=com.example' -X 'fyne.io/fyne/v2/internal/app.MetaVersion=1.2.3'", b.generateMetaLDFlags())
}
