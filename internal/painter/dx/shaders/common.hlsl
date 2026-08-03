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
};

static const float PI = 3.141592653589793;

struct VSIn {
    float2 pos : POSITION;
    float2 uv  : TEXCOORD0;
};

struct PSIn {
    float4 pos : SV_POSITION;
    float2 uv  : TEXCOORD0;
};

float4 blend_shadow(float4 col, float4 shadow)
{
    float alpha = col.a + shadow.a * (1.0 - col.a);
    if (alpha <= 0.0) return float4(0.0, 0.0, 0.0, 0.0);
    return float4((col.rgb * col.a + shadow.rgb * shadow.a * (1.0 - col.a)) / alpha, alpha);
}

// centredGL returns the fragment position relative to the shape centre with y
// pointing up, matching what the GLSL shaders derive from gl_FragCoord.
float2 centredGL(float2 fragPos)
{
    float2 centre = float2((bounds.x + bounds.z) * 0.5, (bounds.y + bounds.w) * 0.5);
    return float2(fragPos.x - centre.x, centre.y - fragPos.y);
}
