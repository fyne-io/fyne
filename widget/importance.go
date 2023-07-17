package widget

type Importance int

const (
	// MediumImportance applies a standard appearance.
	MediumImportance Importance = iota
	// HighImportance applies a prominent appearance.
	HighImportance
	// LowImportance applies a subtle appearance.
	LowImportance

	// DangerImportance applies an error theme to the widget.
	//
	// Since 2.3
	DangerImportance
	// WarningImportance applies a warning theme to the widget.
	//
	// Since 2.3
	WarningImportance

	// SuccessImportance applies a success theme to the widget.
	SuccessImportance
)
