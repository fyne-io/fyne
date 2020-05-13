package theme

import "image/color"

// Brighten takes a color and brightens / dims it by the given
// factor. Each value in the factor represents a 10% increase
// in the intensity of the color.  Negative values allow for 
// dimming the color
func Brighten(c color.Color, factor float64) color.Color {

}