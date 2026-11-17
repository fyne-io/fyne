//go:build windows || darwin || linux || openbsd || freebsd

package gl

func (p *painter) updateBuffer(vbo Buffer, points []float32) {
	p.ctx.BindBuffer(arrayBuffer, vbo)
	p.logError()
	// BufferSubData seems significantly less performant on desktop
	// so use BufferData instead. points is usually a draw call's stack
	// array, and passing it through the context interface would move it to
	// the heap on every draw, so it is uploaded from a buffer the painter owns.
	p.vertices = append(p.vertices[:0], points...)
	p.ctx.BufferData(arrayBuffer, p.vertices, staticDraw)
	p.logError()
}
