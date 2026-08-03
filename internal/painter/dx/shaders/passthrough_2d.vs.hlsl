// Port of gl/shaders/passthrough_2d.vert.
//
// Emits clip-space vertices unchanged.
PSIn main(VSIn input)
{
    PSIn output;
    output.pos = float4(input.pos, 0.0, 1.0);
    output.uv = input.uv;
    return output;
}
