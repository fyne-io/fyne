// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package geom defines a two-dimensional coordinate system.

The coordinate system is based on an left-handed Cartesian plane.
That is, X increases to the right and Y increases down. For (x,y),

	(0,0) → (1,0)
	  ↓   ↘
	(0,1)   (1,1)

The display window places the origin (0, 0) in the upper-left corner of
the screen. Positions on the plane are measured in typographic points,
1/72 of an inch, which is represented by the Pt type.

Any interface that draws to the screen using types from the geom package
scales the number of pixels to maintain a Pt as 1/72 of an inch.
*/
package geom // import "github.com/fyne-io/mobile/geom"

/*
Notes on the various underlying coordinate systems.

Both Android and iOS (UIKit) use upper-left-origin coordinate systems
with for events, however they have different units.

UIKit measures distance in points. A point is a single-pixel on a
pre-Retina display. UIKit maintains a scale factor that to turn points
into pixels. On current retina devices, the scale factor is 2.0.

A UIKit point does not correspond to a fixed physical distance, as the
iPhone has a 163 DPI/PPI (326 PPI retina) display, and the iPad has a
132 PPI (264 retina) display. Points are 32-bit floats.

Even though point is the official UIKit term, they are commonly called
pixels. Indeed, the units were equivalent until the retina display was
introduced.

N.b. as a UIKit point is unrelated to a typographic point, it is not
related to this packages's Pt and Point types.

More details about iOS drawing:

https://developer.apple.com/library/ios/documentation/2ddrawing/conceptual/drawingprintingios/GraphicsDrawingOverview/GraphicsDrawingOverview.html

Android uses pixels. Sub-pixel precision is possible, so pixels are
represented as 32-bit floats. The ACONFIGURATION_DENSITY enum provides
the screen DPI/PPI, which varies frequently between devices.

It would be tempting to adopt the pixel, given the clear pixel/DPI split
in the core android events API. However, the plot thickens:

http://developer.android.com/training/multiscreen/screendensities.html

Android promotes the notion of a density-independent pixel in many of
their interfaces, often prefixed by "dp". 1dp is a real physical length,
as "independent" means it is assumed to be 1/160th of an inch and is
adjusted for the current screen.

In addition, android has a scale-indepdendent pixel used for expressing
a user's preferred text size. The user text size preference is a useful
notion not yet expressed in the geom package.

For the sake of clarity when working across platforms, the geom package
tries to put distance between it and the word pixel.
*/

import "fmt"

// Pt is a length.
//
// The unit Pt is a typographical point, 1/72 of an inch (0.3527 mm).
//
// It can be be converted to a length in current device pixels by
// multiplying with PixelsPerPt after app initialization is complete.
type Pt float32

// Px converts the length to current device pixels.
func (p Pt) Px(pixelsPerPt float32) float32 { return float32(p) * pixelsPerPt }

// String returns a string representation of p like "3.2pt".
func (p Pt) String() string { return fmt.Sprintf("%.2fpt", p) }

// Point is a point in a two-dimensional plane.
type Point struct {
	X, Y Pt
}

// String returns a string representation of p like "(1.2,3.4)".
func (p Point) String() string { return fmt.Sprintf("(%.2f,%.2f)", p.X, p.Y) }

// A Rectangle is region of points.
// The top-left point is Min, and the bottom-right point is Max.
type Rectangle struct {
	Min, Max Point
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rectangle) String() string { return r.Min.String() + "-" + r.Max.String() }
