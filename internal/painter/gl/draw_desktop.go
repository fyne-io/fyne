//go:build windows || darwin || linux || openbsd || freebsd

package gl

func (p *painter) updateBuffer(vbo Buffer, points []float32) {
	p.ctx.BindBuffer(arrayBuffer, vbo)
	p.logError()
	// BufferSubData seems significantly less performant on desktop
	// so use BufferData instead
	p.ctx.BufferData(arrayBuffer, points, staticDraw)
	p.logError()
}
