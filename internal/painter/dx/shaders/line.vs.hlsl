// Port of gl/shaders/line.vert.
//
// The vertex carries an edge normal which is scaled by the half line width and
// added to the position; the resulting delta is interpolated so the pixel stage
// can feather the edge.
PSIn main(VSIn input)
{
    PSIn output;
    float2 delta = input.uv * misc.x;
    output.pos = float4(input.pos + delta, 0.0, 1.0);
    output.uv = delta;
    return output;
}
