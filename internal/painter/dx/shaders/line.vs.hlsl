// Replaces the port of gl/shaders/line.vert.
//
// The six vertices of the line quad are generated from the endpoints in
// ndcRect (p1 = xy, p2 = zw) and the edge normal in inset.xy - lines sample no
// texture, so inset is free. The normal is scaled by the half line width in
// misc.x and added to the position; the delta is interpolated so the pixel
// stage can feather the edge. Triangle-list order matches what the painter
// used to upload:
//   0: p1+n   1: p2+n   2: p2-n   3: p2-n   4: p1+n   5: p1-n
PSIn main(uint vid : SV_VertexID)
{
    PSIn output;
    bool second = vid >= 1u && vid <= 3u;
    bool negative = vid == 2u || vid == 3u || vid == 5u;
    float2 pt = second ? ndcRect.zw : ndcRect.xy;
    float2 normal = negative ? -inset.xy : inset.xy;
    float2 delta = normal * misc.x;
    output.pos = float4(pt + delta, 0.0, 1.0);
    output.uv = delta;
    return output;
}
