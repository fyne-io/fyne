package canvas

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAdvanceShaderTime(t *testing.T) {
	var elapsed time.Duration
	var lastTick time.Time
	base := time.Unix(100, 0)

	// the first tick only establishes the clock, accumulating no time
	elapsed, lastTick = advanceShaderTime(elapsed, lastTick, base)
	assert.Equal(t, time.Duration(0), elapsed)

	// a normal tick advances elapsed by its own duration
	elapsed, lastTick = advanceShaderTime(elapsed, lastTick, base.Add(16*time.Millisecond))
	assert.Equal(t, 16*time.Millisecond, elapsed)

	// an over-cap gap is treated as a pause/resume and contributes nothing, so
	// the shader neither jumps forward nor counts time spent stopped
	elapsed, lastTick = advanceShaderTime(elapsed, lastTick, base.Add(16*time.Millisecond+5*time.Second))
	assert.Equal(t, 16*time.Millisecond, elapsed)

	// once resumed, a normal tick continues to accumulate from where it left off
	elapsed, _ = advanceShaderTime(elapsed, lastTick, base.Add(32*time.Millisecond+5*time.Second))
	assert.Equal(t, 32*time.Millisecond, elapsed)
}
