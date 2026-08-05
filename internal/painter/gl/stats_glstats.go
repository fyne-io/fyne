//go:build glstats

// Package-local GL call instrumentation, compiled in only with -tags glstats.
// It exists to produce before/after numbers for text rendering changes and is
// not part of a normal build.
//
// Environment:
//
//	FYNE_GL_STATS_OUT    file to write the JSON report to (default gl_stats.json)
//	FYNE_GL_STATS_FRAMES number of frames to record before writing (default 600)
//	FYNE_GL_STATS_WARMUP frames to discard at the start (default 60)

package gl

import (
	"encoding/json"
	"os"
	"strconv"
	"time"
)

type frameStats struct {
	TexImage2DCalls    int   `json:"texImage2DCalls"`
	TexImage2DBytes    int   `json:"texImage2DBytes"`
	TexSubImage2DCalls int   `json:"texSubImage2DCalls"`
	TexSubImage2DBytes int   `json:"texSubImage2DBytes"`
	DrawArraysCalls    int   `json:"drawArraysCalls"`
	BufferDataCalls    int   `json:"bufferDataCalls"`
	BufferDataFloats   int   `json:"bufferDataFloats"`
	CreateTextureCalls int   `json:"createTextureCalls"`
	DeleteTextureCalls int   `json:"deleteTextureCalls"`
	SubmitNanos        int64 `json:"submitNanos"`
}

type statsContext struct {
	context

	cur     frameStats
	frames  []frameStats
	started time.Time
	last    time.Time
	inFrame bool

	limit   int
	warmup  int
	out     string
	written bool
}

func envInt(name string, fallback int) int {
	if v := os.Getenv(name); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func wrapContext(c context) context {
	out := os.Getenv("FYNE_GL_STATS_OUT")
	if out == "" {
		out = "gl_stats.json"
	}
	return &statsContext{
		context: c,
		limit:   envInt("FYNE_GL_STATS_FRAMES", 600),
		warmup:  envInt("FYNE_GL_STATS_WARMUP", 60),
		out:     out,
	}
}

// Clear marks a frame boundary: the painter calls it once at the start of every
// paint pass, so the previous frame's counters are complete by now.
func (c *statsContext) Clear(mask uint32) {
	c.endFrame()
	c.cur = frameStats{}
	c.started = time.Now()
	c.last = c.started
	c.inFrame = true
	c.context.Clear(mask)
}

// endFrame closes off the in-progress frame. Duration runs from the Clear that
// opened the frame to the last counted GL call within it, so the vsync wait
// inside SwapBuffers is excluded. This measures CPU-side submission cost, not
// GPU execution.
func (c *statsContext) endFrame() {
	if !c.inFrame {
		return
	}
	c.inFrame = false
	c.cur.SubmitNanos = c.last.Sub(c.started).Nanoseconds()
	c.frames = append(c.frames, c.cur)

	if len(c.frames) >= c.limit+c.warmup && !c.written {
		c.write()
	}
}

func (c *statsContext) write() {
	c.written = true
	frames := c.frames
	if len(frames) > c.warmup {
		frames = frames[c.warmup:]
	}

	data, err := json.MarshalIndent(struct {
		Frames []frameStats `json:"frames"`
	}{frames}, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(c.out, data, 0o644)
	os.Stderr.WriteString("glstats: wrote " + strconv.Itoa(len(frames)) + " frames to " + c.out + "\n")
	os.Exit(0)
}

func (c *statsContext) tick() {
	c.last = time.Now()
}

func (c *statsContext) TexImage2D(target uint32, level, width, height int, colorFormat, typ uint32, data []uint8) {
	c.cur.TexImage2DCalls++
	c.cur.TexImage2DBytes += len(data)
	c.context.TexImage2D(target, level, width, height, colorFormat, typ, data)
	c.tick()
}

func (c *statsContext) TexSubImage2D(target uint32, level, xoffset, yoffset, width, height int, colorFormat, typ uint32, data []uint8) {
	c.cur.TexSubImage2DCalls++
	c.cur.TexSubImage2DBytes += len(data)
	c.context.TexSubImage2D(target, level, xoffset, yoffset, width, height, colorFormat, typ, data)
	c.tick()
}

func (c *statsContext) DrawArrays(mode uint32, first, count int) {
	c.cur.DrawArraysCalls++
	c.context.DrawArrays(mode, first, count)
	c.tick()
}

func (c *statsContext) BufferData(target uint32, points []float32, usage uint32) {
	c.cur.BufferDataCalls++
	c.cur.BufferDataFloats += len(points)
	c.context.BufferData(target, points, usage)
	c.tick()
}

func (c *statsContext) BufferSubData(target uint32, points []float32) {
	c.cur.BufferDataCalls++
	c.cur.BufferDataFloats += len(points)
	c.context.BufferSubData(target, points)
	c.tick()
}

func (c *statsContext) CreateTexture() Texture {
	c.cur.CreateTextureCalls++
	return c.context.CreateTexture()
}

func (c *statsContext) DeleteTexture(texture Texture) {
	c.cur.DeleteTextureCalls++
	c.context.DeleteTexture(texture)
}
