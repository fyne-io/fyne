#version 110

uniform sampler2D tex;

varying vec2 fragTexCoord;

void main() {
    vec4 texColor = texture2D(tex, fragTexCoord);
    if(texColor.a < 0.01)
        discard;
    gl_FragColor = texColor;
}
