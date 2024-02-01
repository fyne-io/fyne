//go:generate fyne bundle -o metadata_bundled.go -package data -name resourceAuthors ../../../AUTHORS

package data

// Authors contains the list of the toolkit authors, extracted from ../../../AUTHORS.
var Authors = resourceAuthors
