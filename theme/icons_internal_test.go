package theme

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
)

func helperNewStaticResource() *fyne.StaticResource {
	return &fyne.StaticResource{
		StaticName: "content-remove.svg",
		StaticContent: []byte{
			60, 115, 118, 103, 32, 120, 109, 108, 110, 115, 61, 34, 104, 116, 116, 112, 58, 47, 47, 119, 119, 119, 46, 119, 51, 46, 111, 114, 103, 47, 50, 48, 48, 48, 47, 115, 118, 103, 34, 32, 119, 105, 100, 116, 104, 61, 34, 50, 52, 34, 32, 104, 101, 105, 103, 104, 116, 61, 34, 50, 52, 34, 32, 118, 105, 101, 119, 66, 111, 120, 61, 34, 48, 32, 48, 32, 50, 52, 32, 50, 52, 34, 62, 60, 112, 97, 116, 104, 32, 102, 105, 108, 108, 61, 34, 35, 102, 102, 102, 102, 102, 102, 34, 32, 100, 61, 34, 77, 49, 57, 32, 49, 51, 72, 53, 118, 45, 50, 104, 49, 52, 118, 50, 122, 34, 47, 62, 60, 112, 97, 116, 104, 32, 100, 61, 34, 77, 48, 32, 48, 104, 50, 52, 118, 50, 52, 72, 48, 122, 34, 32, 102, 105, 108, 108, 61, 34, 110, 111, 110, 101, 34, 47, 62, 60, 47, 115, 118, 103, 62},
	}
}

func Test_UnwrapResource(t *testing.T) {
	source := helperNewStaticResource()
	res := NewThemedResource(
		NewErrorThemedResource(
			NewInvertedThemedResource(
				NewDisabledResource(source),
			),
		),
	)

	// check that resource name only has the top-level themed prefix
	assert.Equal(t, "foreground_content-remove.svg", res.Name())

	unwrapped := unwrapResource(res)
	_, ok := unwrapped.(*fyne.StaticResource)
	assert.True(t, ok, "unwrap did not return base resource type")
}
