#version 110

uniform float alpha;
uniform sampler2D tex;
uniform float cornerRadius;  // in pixels
uniform vec2 size;           // in pixels: size of the rendered image quad
uniform vec4 inset;          // texture coordinate insets (minX, minY, maxX, maxY)

varying vec2 fragTexCoord;

void main() {
    float fragAlpha = 1.0;
    if (cornerRadius > 0.5) {
        // normalize texture coords from [insetMin, insetMax] to [0, 1]
        // this makes the rounding calculation work correctly for cropped/covered images
        vec2 normalizedCoord = (fragTexCoord - inset.xy) / (inset.zw - inset.xy);

        vec2 pos = normalizedCoord * size;
        vec2 halfSize = size * 0.5;
        float dist = length(max(abs(pos - halfSize) - halfSize + cornerRadius, 0.0)) - cornerRadius;
        fragAlpha = 1.0 - smoothstep(-1.0, 1.0, dist);
    }

    vec4 texColor = texture2D(tex, fragTexCoord) * alpha * fragAlpha;
    if(texColor.a < 0.01)
        discard;

    gl_FragColor = texColor;
}
