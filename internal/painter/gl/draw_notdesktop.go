//go:build !(windows || darwin || linux || openbsd || freebsd)

package gl

func (p *painter) updateBuffer(vbo Buffer, points []float32) {
	p.ctx.BindBuffer(arrayBuffer, vbo)
	p.logError()
	// See the desktop updateBuffer for why this copies.
	p.vertices = append(p.vertices[:0], points...)
	p.ctx.BufferSubData(arrayBuffer, p.vertices)
	p.logError()
}
