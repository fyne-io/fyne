package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RingBuffer(t *testing.T) {
	buf := NewRingBuffer[int](8)

	assertDequeue := func(expect int) {
		t.Helper()
		got, valid := buf.Pull()
		assert.True(t, valid, "got invalid pull")
		assert.Equal(t, expect, got, "invalid pull result")
	}

	buf.Push(1)
	buf.Push(2)
	buf.Push(3)
	buf.Push(4)
	buf.Push(5)
	buf.Push(6)
	buf.Push(7)
	buf.Push(8)
	buf.Push(9)
	buf.Push(10)

	assertDequeue(1)
	assertDequeue(2)

	buf.Push(11)
	buf.Push(12)
	buf.Push(13)
	buf.Push(14)
	buf.Push(15)
	buf.Push(16)
	buf.Push(17)
	buf.Push(18)
	buf.Push(19)

	for i := 3; i <= 19; i++ {
		assertDequeue(i)
	}
}
