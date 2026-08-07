// An alpha-modulated texture with optional rounded corners. Covers images,
// rasters, gradients and text.
Texture2D tex : register(t0);
SamplerState samp : register(s0);

float4 main(PSIn input) : SV_TARGET
{
    float alpha = texParams.x;
    float cornerRadius = texParams.y;
    float2 size = texParams.zw;

    float fragAlpha = 1.0;
    if (cornerRadius > 0.5) {
        float2 normalizedCoord = (input.uv - inset.xy) / (inset.zw - inset.xy);
        float2 p = normalizedCoord * size;
        float2 halfSize = size * 0.5;
        float dist = length(max(abs(p - halfSize) - halfSize + cornerRadius, 0.0)) - cornerRadius;
        fragAlpha = 1.0 - smoothstep(-1.0, 1.0, dist);
    }

    float4 texColor = tex.Sample(samp, input.uv) * alpha * fragAlpha;
    if (texColor.a < 0.01) {
        discard;
    }
    return texColor;
}
