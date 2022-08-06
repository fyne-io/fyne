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

uniform vec4 fill_color;
uniform vec4 stroke_color;

varying vec2 delta;
varying float switch_var;
varying float lineWidth_var;
varying float feather_var;
vec4 color;

void main() {
    if ( switch_var == 1.0 ){
        color = fill_color;
    }; 
    if ( switch_var == 2.0 ){	
        color = stroke_color;
    };
    if ( switch_var > 50.0 ){	
        color = vec4(stroke_color.r, stroke_color.g, stroke_color.b, stroke_color.a * switch_var / 100.0 );
    };
    if ( switch_var >= 10.0 && switch_var <= 50.0 ){	
        color = vec4(fill_color.r, fill_color.g, fill_color.b, fill_color.a * switch_var / 100.0 );
    };
    if ( delta != vec2(0.0, 0.0) ){
        float distance = length(delta);

        if (distance <= lineWidth_var - feather_var) {
            gl_FragColor = color;
        } else {
            gl_FragColor = vec4(color.r, color.g, color.b, mix(color.a, 0.0, (distance - (lineWidth_var - feather_var)) / feather_var));
        }
    } else {
        gl_FragColor = color;
    }
}
