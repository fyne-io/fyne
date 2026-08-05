#version 100

#ifdef GL_ES
# ifdef GL_FRAGMENT_PRECISION_HIGH
precision highp float;
# else
precision mediump float;
#endif
precision mediump int;
precision lowp sampler2D;
#endif

// Glyph quads arrive in device pixels relative to the origin of the string they
// belong to, rather than in clip space. That keeps the vertex data independent
// of where the string sits and how large the window is, so it can be cached and
// reused for as long as the text itself is unchanged, with only the two
// uniforms below moving it into place.

attribute vec2 vert;          // device pixels from the string origin, y downwards
attribute vec2 vertTexCoord;

uniform vec2 origin;          // clip-space position of the string origin
uniform vec2 pixelScale;      // clip-space units per device pixel (y negative)

varying vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;

    gl_Position = vec4(origin + vert * pixelScale, 0, 1);
}
