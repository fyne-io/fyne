#version 110

uniform vec4 color;

uniform sampler2D tex;
varying vec2 fragTexCoord;

void main()
{
    gl_FragColor = vec4(color.r, color.g, color.b, texture2D(tex, fragTexCoord).r*color.a);
}
