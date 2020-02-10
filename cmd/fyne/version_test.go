package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePkgString(t *testing.T) {
	type test struct {
		in string // Test input
		es string // Expected output
		ee bool   // Expecting error?
	}

	ts := []test{
		{in: "/home/user/go/pkg/mod/fyne.io/fyne@v1.2.2/pkg/li", es: "v1.2.2", ee: false},
		{in: "fyne.io/fyne@v1.1", es: "v1.1", ee: false},
		{in: "/home/user/fyne/pkg/linux_amd64/fyne.io/fyne.a", es: "", ee: true},
		{in: "v1.2.2", es: "", ee: true},
		{in: "/home/user/go/pkg/mod/github.com/peterbourgon/diskv@v2.0.1+incompatible/pkg/linux_amd64/github.com/peterbourgon/diskv.a", es: "v2.0.1+incompatible", ee: false},
		{in: "/home/ser/go/pkg/mod/github.com/peterbourgon/diskv/v3@v3.0.0/pkg/linux_amd64/github.com/peterbourgon/diskv/v3.a", es: "v3.0.0", ee: false},
	}

	for _, tc := range ts {
		gs, ge := parsePkgString(tc.in)
		assert.Equal(t, tc.es, gs)
		if tc.ee {
			assert.Error(t, ge)
		}
	}
}
