#version 110

// Batched text and plain rectangles, sampled out of the glyph atlas. The atlas
// holds coverage only, so the colour arrives with each vertex - that is what
// lets differently coloured labels and fills share one draw.
uniform sampler2D tex;

varying vec2 fragTexCoord;
varying vec4 fragColor;

void main() {
    float coverage = texture2D(tex, fragTexCoord).a;
    if (coverage <= 0.0)
        discard;

    // Straight alpha, for the blend function the shape shaders use.
    gl_FragColor = vec4(fragColor.rgb, fragColor.a * coverage);
}
