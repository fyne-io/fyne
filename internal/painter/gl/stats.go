package gl

import (
	"fmt"
	"log"
	"math"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
)

// glDebug turns on periodic draw statistics: what a frame costs, what it is
// made of, and how much of it the glyph batch carries.
var glDebug = os.Getenv("FYNE_GL_DEBUG") != ""

// debugReportFrames is how many frames each FYNE_GL_DEBUG report averages over.
const debugReportFrames = 120

// drawStats accumulates the FYNE_GL_DEBUG counters between reports.
type drawStats struct {
	since      time.Time // start of the current report window
	frameStart time.Time // the current frame's Clear
	lastDraw   time.Time // end of the current frame's latest draw
	paint      time.Duration
	frames     int

	draws int // DrawArrays calls, counted by countingContext
	// byType counts the draw calls each canvas object type made on its own;
	// batched quads are counted in batches and quads instead.
	byType  map[reflect.Type]int
	batches int
	quads   int
	// textures counts textures created per object type. For text each one is
	// a whole string rasterised and uploaded, which only happens for strings
	// the glyph atlas cannot hold.
	textures map[reflect.Type]int
}

// countingContext counts draw calls for FYNE_GL_DEBUG.
type countingContext struct {
	context
	draws *int
}

func (c *countingContext) DrawArrays(mode uint32, first, count int) {
	*c.draws++
	c.context.DrawArrays(mode, first, count)
}

// startFrame is called from Clear. It closes the books on the previous frame,
// logs a report every debugReportFrames frames, and wraps the context in the
// draw counter on first use (or after Init has replaced it).
func (p *painter) startFrame() {
	s := &p.stats
	if _, ok := p.ctx.(*countingContext); !ok {
		p.ctx = &countingContext{context: p.ctx, draws: &s.draws}
	}

	now := time.Now()
	if s.frameStart.IsZero() {
		s.since = now
	} else {
		if s.lastDraw.After(s.frameStart) {
			s.paint += s.lastDraw.Sub(s.frameStart)
		}
		s.frames++
		if s.frames == debugReportFrames {
			p.report(now)
		}
	}
	s.frameStart = now
}

func (p *painter) report(now time.Time) {
	s := &p.stats
	size := p.canvas.Size()
	log.Printf("gl %.0fx%.0f: %.1f fps, paint %v/frame, %s draws/frame (%s; %s batches carrying %s quads); "+
		"new textures/frame: %s",
		size.Width, size.Height,
		float64(s.frames)/now.Sub(s.since).Seconds(),
		(s.paint / time.Duration(s.frames)).Round(10*time.Microsecond),
		perFrame(s.draws, s.frames), breakdown(s.byType, s.frames),
		perFrame(s.batches, s.frames), perFrame(s.quads, s.frames),
		breakdown(s.textures, s.frames))

	clear(s.byType)
	clear(s.textures)
	*s = drawStats{since: now, byType: s.byType, textures: s.textures}
}

// countDraws attributes the n draw calls obj just made on its own.
func (p *painter) countDraws(obj fyne.CanvasObject, n int) {
	p.stats.lastDraw = time.Now()
	if n > 0 {
		count(&p.stats.byType, obj, n)
	}
}

// countBatch records a flushed batch of quads.
func (p *painter) countBatch(quads int) {
	p.stats.batches++
	p.stats.quads += quads
	p.stats.lastDraw = time.Now()
}

// count adds n to (*m)[type of o], making the map on first use.
func count(m *map[reflect.Type]int, o fyne.CanvasObject, n int) {
	if *m == nil {
		*m = make(map[reflect.Type]int)
	}
	(*m)[reflect.TypeOf(o)] += n
}

// breakdown formats per-frame counts by type, largest first: "Text 402, Line 90".
func breakdown(m map[reflect.Type]int, frames int) string {
	types := make([]reflect.Type, 0, len(m))
	for t := range m {
		types = append(types, t)
	}
	slices.SortFunc(types, func(a, b reflect.Type) int { return m[b] - m[a] })
	parts := make([]string, len(types))
	for i, t := range types {
		parts[i] = fmt.Sprintf("%s %s", strings.TrimPrefix(t.String(), "*canvas."), perFrame(m[t], frames))
	}
	return strings.Join(parts, ", ")
}

// perFrame formats v/frames to one decimal, dropping a trailing ".0".
func perFrame(v, frames int) string {
	return strconv.FormatFloat(math.Round(float64(v)/float64(frames)*10)/10, 'f', -1, 64)
}
