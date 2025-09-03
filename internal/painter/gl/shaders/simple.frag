#version 110

uniform sampler2D tex;

varying vec2 fragTexCoord;
varying float fragAlpha;

void main() {
    vec4 texColor = texture2D(tex, fragTexCoord);
    texColor.a *= fragAlpha;
    texColor.r *= fragAlpha;
    texColor.g *= fragAlpha;
    texColor.b *= fragAlpha;

    if(texColor.a < 0.01)
        discard;
    gl_FragColor = texColor;
}
