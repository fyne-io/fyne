// Pixel stage for atlas text.
//
// The atlas stores coverage only, so the glyph's colour arrives per instance
// from the vertex stage rather than being baked into the texture the way the
// whole-run path bakes it. That is what lets differently coloured labels share
// one batch.
Texture2D atlas : register(t0);
SamplerState samp : register(s0);

float4 main(GlyphPSIn input) : SV_TARGET
{
    float coverage = atlas.Sample(samp, input.uv).r;
    if (coverage <= 0.0) {
        discard;
    }
    // Straight alpha, matching the blend state the shape shaders use; the
    // whole-run text path is premultiplied because its texture already is.
    return float4(input.color.rgb, input.color.a * coverage);
}
