// +build ios

package gl

import "golang.org/x/mobile/gl"

func (p *glPainter) glFreeBuffer(vbo Buffer) {
	ctx := p.glctx()

	ctx.BindBuffer(gl.ARRAY_BUFFER, gl.Buffer(vbo))
	ctx.DeleteBuffer(gl.Buffer(vbo))
}
