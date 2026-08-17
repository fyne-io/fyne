//go:build glstats

// Package-local GL call instrumentation, compiled in only with -tags glstats.
// It exists to produce before/after numbers for text rendering changes and is
// not part of a normal build.
//
// Environment:
//
//	FYNE_GL_STATS_OUT    file to write a JSON report to. When set, the run also
//	                     ends by itself once enough frames have been recorded.
//	                     Leave it unset on a phone, where there is nowhere
//	                     convenient to write and no way to set it: the summary
//	                     still reaches the log, and the caller decides when the
//	                     run is over by calling ReportStats.
//	FYNE_GL_STATS_FRAMES number of frames to record before writing (default 600)
//	FYNE_GL_STATS_WARMUP frames to discard at the start of each phase (default 60)

package gl

import (
	"encoding/json"
	"fmt"
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

// phaseStats is one labelled stretch of a run, so a single session can measure
// scrolling and then typing without needing to be started twice. That matters on
// a device, where launching the app twice with different settings is far more
// awkward than on a desktop.
type phaseStats struct {
	Name   string       `json:"name"`
	Frames []frameStats `json:"frames"`
}

type statsContext struct {
	context

	cur     frameStats
	phases  []*phaseStats
	started time.Time
	last    time.Time
	inFrame bool

	limit    int
	warmup   int
	out      string
	autoExit bool
	written  bool
}

// active is the statsContext of the current painter, so that the harness can
// label phases and ask for a report without threading a reference to it.
var active *statsContext

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
	s := &statsContext{
		context:  c,
		limit:    envInt("FYNE_GL_STATS_FRAMES", 600),
		warmup:   envInt("FYNE_GL_STATS_WARMUP", 60),
		out:      out,
		autoExit: out != "",
		phases:   []*phaseStats{{Name: "default"}},
	}
	active = s
	return s
}

// MarkPhase ends the stretch being measured and starts another under a new
// name. Frames recorded from here on are reported separately.
func MarkPhase(name string) {
	if active == nil {
		return
	}
	active.endFrame()
	active.phases = append(active.phases, &phaseStats{Name: name})
}

// ReportStats writes the summary and ends the process. Call it when the
// workload is done; on a phone this is the only thing that ends the run.
func ReportStats() {
	if active == nil {
		return
	}
	active.endFrame()
	active.write()
}

func (c *statsContext) phase() *phaseStats {
	return c.phases[len(c.phases)-1]
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
	p := c.phase()
	p.Frames = append(p.Frames, c.cur)

	if c.autoExit && len(p.Frames) >= c.limit+c.warmup && !c.written {
		c.write()
	}
}

// measured drops the warmup frames, which are dominated by first-time
// rasterisation and texture allocation rather than the steady state.
func (p *phaseStats) measured(warmup int) []frameStats {
	if len(p.Frames) > warmup {
		return p.Frames[warmup:]
	}
	return nil
}

func (c *statsContext) write() {
	c.written = true

	// The summary goes to stderr because that is the one channel available
	// everywhere: on Android it is piped into logcat under the tag "Fyne".
	for _, p := range c.phases {
		frames := p.measured(c.warmup)
		if len(frames) == 0 {
			continue
		}
		var bufCalls, drawCalls, bufFloats, texBytes, subCalls float64
		var submit float64
		for _, f := range frames {
			bufCalls += float64(f.BufferDataCalls)
			drawCalls += float64(f.DrawArraysCalls)
			bufFloats += float64(f.BufferDataFloats)
			texBytes += float64(f.TexImage2DBytes + f.TexSubImage2DBytes)
			subCalls += float64(f.TexSubImage2DCalls)
			submit += float64(f.SubmitNanos)
		}
		n := float64(len(frames))
		vertKiB := bufFloats * 4 / n / 1024
		texKiB := texBytes / n / 1024
		fmt.Fprintf(os.Stderr,
			"glstats phase=%s frames=%d bufCalls=%.1f drawCalls=%.1f subImgCalls=%.2f "+
				"vertKiB=%.2f texKiB=%.2f totalKiB=%.2f submitMs=%.2f\n",
			p.Name, len(frames), bufCalls/n, drawCalls/n, subCalls/n,
			vertKiB, texKiB, vertKiB+texKiB, submit/n/1e6)
	}

	if c.out != "" {
		if data, err := json.MarshalIndent(struct {
			Phases []*phaseStats `json:"phases"`
			Frames []frameStats  `json:"frames"`
		}{c.phases, c.phases[len(c.phases)-1].measured(c.warmup)}, "", "  "); err == nil {
			_ = os.WriteFile(c.out, data, 0o644)
			os.Stderr.WriteString("glstats: wrote " + strconv.Itoa(len(c.phases)) + " phase(s) to " + c.out + "\n")
		}
	}
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
