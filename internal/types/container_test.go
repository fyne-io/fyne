package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"fyne.io/fyne/v2"
	container2 "fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/internal/types"
)

func TestIsContainer(t *testing.T) {
	assert.True(t, types.IsContainer(&fyne.Container{}))
	assert.True(t, types.IsContainer(container2.NewWithoutLayout()))
	assert.True(t, types.IsContainer(&extendedContainer{}))
}

type extendedContainer struct {
	fyne.Container
}
