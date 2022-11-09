# Freetype for Golang

Package fonts provides a unified API for various font formats.
Its design is vastly inspired by the C freetype2 library, although the logic of the parsers are from diverse origines.

The main purpose of this package is to serve as a font analyser for [fontconfig](github.com/benoitkugler/textlayout/fontconfig), using the `fonts.Loader` interface.
Individuals parsers may also be used separately.
