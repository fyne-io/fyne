// Replaces the port of gl/shaders/passthrough_2d.vert.
//
// Every quad the painter draws is an axis-aligned rectangle in clip space, so
// the four strip corners are generated from ndcRect and the texture
// coordinates from inset instead of being uploaded per draw - that removes a
// vertex-buffer UpdateSubresource from every draw call. Corner order matches
// what the painter used to upload:
//   0: (x1, y2) uv (minX, maxY)    1: (x1, y1) uv (minX, minY)
//   2: (x2, y2) uv (maxX, maxY)    3: (x2, y1) uv (maxX, minY)
// Shape shaders ignore uv (they work from SV_POSITION and bounds), so inset
// being zero for them is harmless.
PSIn main(uint vid : SV_VertexID)
{
    PSIn output;
    bool right = (vid & 2u) != 0u;
    bool first = (vid & 1u) != 0u;
    float x = right ? ndcRect.z : ndcRect.x;
    float y = first ? ndcRect.y : ndcRect.w;
    float u = right ? inset.z : inset.x;
    float v = first ? inset.y : inset.w;
    output.pos = float4(x, y, 0.0, 1.0);
    output.uv = float2(u, v);
    return output;
}
