// Vertex stage for atlas text. It has no counterpart in gl/shaders: the OpenGL
// painter draws each string as one textured quad from a texture of its own,
// where this draws one quad per glyph out of the shared atlas.
//
// Six vertices per instance as a triangle list, generated from the instance
// rectangle. A strip would need four, but instanced strips leave the question
// of whether the driver restarts between instances, and two triangles cost
// nothing to be sure about. Winding is irrelevant - the rasterizer state culls
// nothing.
//
//   0: (x1, top)  1: (x2, top)  2: (x1, bottom)
//   3: (x1, bottom)  4: (x2, top)  5: (x2, bottom)
GlyphPSIn main(uint vid : SV_VertexID, uint iid : SV_InstanceID)
{
    GlyphInst inst = gGlyphs[iid];

    bool right = vid == 1u || vid == 4u || vid == 5u;
    bool bottom = vid == 2u || vid == 3u || vid == 5u;

    GlyphPSIn output;
    output.pos = float4(right ? inst.ndc.z : inst.ndc.x,
                        bottom ? inst.ndc.w : inst.ndc.y, 0.0, 1.0);
    output.uv = float2(right ? inst.uv.z : inst.uv.x,
                       bottom ? inst.uv.w : inst.uv.y);
    output.color = inst.color;
    return output;
}
