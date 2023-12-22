package mime

import "strings"

// Split spits the mimetype into its main type and subtype.
func Split(mimeTypeFull string) (mimeType, mimeSubType string) {
	// Replace with strings.Cut() when Go 1.18 is our new base version.
	separatorIndex := strings.IndexByte(mimeTypeFull, '/')
	if separatorIndex == -1 || mimeTypeFull[separatorIndex+1:] == "" {
		return "", "" // Empty or only one part. Ignore.
	}

	mimeType = mimeTypeFull[:separatorIndex]
	mimeSubType = mimeTypeFull[separatorIndex+1:]
	return
}
