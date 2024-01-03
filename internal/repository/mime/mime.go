package mime

import "strings"

// Split spits the mimetype into its main type and subtype.
func Split(mimeTypeFull string) (mimeType, mimeSubType string) {
	mimeType, mimeSubType, ok := strings.Cut(mimeTypeFull, "/")
	if !ok || mimeSubType == "" {
		return "", ""
	}

	return mimeType, mimeSubType
}
