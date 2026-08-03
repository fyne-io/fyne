// Port of gl/shaders/line.frag.
//
// A solid colour with an edge feathered from the delta emitted by line.vs.hlsl.
float4 main(PSIn input) : SV_TARGET
{
    float lineWidth = misc.x;
    float feather = shadowOffset.w;
    float4 col = fillColor;

    float dist = length(input.uv);
    if (feather == 0.0 || dist <= lineWidth - feather) {
        return col;
    }
    return float4(col.rgb, lerp(col.a, 0.0, (dist - (lineWidth - feather)) / feather));
}
