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

// Quads arrive in device pixels relative to the string origin, so the geometry
// can be cached and reused while only the uniforms change.

attribute vec2 vert;          // device pixels from the string origin, y down
attribute vec2 vertTexCoord;

uniform vec2 origin;          // clip-space position of the string origin
uniform vec2 pixelScale;      // clip-space units per device pixel (y negative)

varying vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;

    gl_Position = vec4(origin + vert * pixelScale, 0, 1);
}
