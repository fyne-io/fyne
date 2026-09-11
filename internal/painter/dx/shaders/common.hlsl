// Shared prelude prepended to every shader in this directory.
//
// HLSL has no include mechanism unless the caller supplies an ID3DInclude COM
// handler, so the Go side concatenates this file ahead of each shader rather than
// each shader declaring #include "common.hlsl".
//
// The constant buffer is mirrored by the `constants` struct in the directx
// driver's painter.go, and the two must stay in the same order. Every member is a
// float4 to sidestep HLSL's 16-byte packing rules.
cbuffer Constants : register(b0) {
    float4 frame;        // xy: framebuffer size in pixels
    float4 bounds;       // x1, y1, x2, y2 in pixels, top-left origin
    float4 fillColor;
    float4 strokeColor;
    float4 shadowColor;
    float4 radius;       // topRight, bottomRight, topLeft, bottomLeft
    float4 rectHalf;     // xy: half size, z: strokeWidthHalf, w: edgeSoftness
    float4 misc;         // x: strokeWidth, y: addShadow, z: shadowBlurRadius, w: shadowSpread
    float4 shadowOffset; // xy: offset, z: shadowType, w: line feather
    float4 texParams;    // x: alpha, y: cornerRadius, zw: quad size in pixels
                         // psArc reuses xy as start/end angle - it samples no texture
    float4 inset;        // texture coordinate insets (minX, minY, maxX, maxY)
                         // the line shaders reuse xy as the edge normal
    float4 ndcRect;      // quad corners in clip space (x1, y1, x2, y2)
                         // the line shaders reuse it as the two endpoints
};

static const float PI = 3.141592653589793;

struct PSIn {
    float4 pos : SV_POSITION;
    float2 uv  : TEXCOORD0;
};

// Text is drawn a glyph at a time out of one shared coverage atlas, so a run of
// text objects is a single DrawInstanced rather than an UpdateSubresource and a
// Draw for each. Everything a glyph quad needs is here instead of in the shared
// b0 buffer, which is what makes the batch possible. A constant buffer holds the
// batch rather than a StructuredBuffer because these compile as vs_4_0/ps_4_0
// and structured buffers need shader model 5.
//
// Mirrored by `glyphInst` in painter.go, and the array length by
// `glyphBatchMax`; TestGlyphInstMatchesShader pins both.
struct GlyphInst {
    float4 ndc;   // x1, top, x2, bottom in clip space
    float4 uv;    // u1, v1, u2, v2 in the atlas
    float4 color; // rgba, straight alpha
};

cbuffer GlyphInstances : register(b2) {
    GlyphInst gGlyphs[256];
};

// GlyphPSIn carries the per-instance colour to the pixel stage. It cannot ride
// in PSIn because every other vertex shader would then have to write a field it
// has no use for.
struct GlyphPSIn {
    float4 pos                   : SV_POSITION;
    float2 uv                    : TEXCOORD0;
    nointerpolation float4 color : COLOR0;
};

float4 blend_shadow(float4 col, float4 shadow)
{
    float alpha = col.a + shadow.a * (1.0 - col.a);
    if (alpha <= 0.0) return float4(0.0, 0.0, 0.0, 0.0);
    return float4((col.rgb * col.a + shadow.rgb * shadow.a * (1.0 - col.a)) / alpha, alpha);
}

// centredUp returns the fragment position relative to the shape centre with y
// pointing up, the convention the signed-distance shaders compute in.
float2 centredUp(float2 fragPos)
{
    float2 centre = float2((bounds.x + bounds.z) * 0.5, (bounds.y + bounds.w) * 0.5);
    return float2(fragPos.x - centre.x, centre.y - fragPos.y);
}
