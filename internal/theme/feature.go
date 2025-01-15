package theme

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/internal/cache"
)

type FeatureName string

const FeatureNameDeviceIsMobile = FeatureName("deviceIsMobile")

// FeatureTheme defines the method to look up features that we use internally to apply functional
// differences through a theme override.
type FeatureTheme interface {
	Feature(FeatureName) any
}

// FeatureForWidget looks up the specified feature flag for the requested widget using the current theme.
// This is for internal purposes and will do nothing if the theme has not been overridden with the
// ThemeOverride container.
func FeatureForWidget(name FeatureName, w fyne.Widget) any {
	if custom := cache.WidgetTheme(w); custom != nil {
		if f, ok := custom.(FeatureTheme); ok {
			return f.Feature(name)
		}
	}

	return nil
}
