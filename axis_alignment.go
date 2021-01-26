package fyne

type (
	// Axis represents the two cardinal directions in two dimensions.
	Axis int
	// MainAxisAlignment represents how the children should be placed along the main axis in a flex layout.
	MainAxisAlignment int
	// CrossAxisAlignment represents how the children should be placed along the cross axis in a flex layout.
	CrossAxisAlignment int
	// AxisAlignment represents both Main and Cross AxisAlignment in one object.
	AxisAlignment struct {
		MainAxisAlignment  MainAxisAlignment
		CrossAxisAlignment CrossAxisAlignment
	}
)

// Axis options
const (
	AxisHorizontal Axis = iota
	AxisVertical
)

// MainAxisAlignment options
const (
	// Place the children as close to the start of the main axis as possible.
	MainAxisAlignmentStart MainAxisAlignment = iota
	// Place the children as close to the end of the main axis as possible.
	MainAxisAlignmentEnd
	// Place the children as close to the middle of the main axis as possible.
	MainAxisAlignmentCenter
	// Place the free space evenly between the children.
	MainAxisAlignmentSpaceBetween
	// Place the free space evenly between the children as well as half of that
	// space before and after the first and last child.
	MainAxisAlignmentSpaceAround
	// Place the free space evenly between the children as well as before and
	// after the first and last child.
	MainAxisAlignmentSpaceEvenly
)

// CrossAxisAlignment options
const (
	// Place the children so that their centers align with the middle of the
	// cross axis.
	CrossAxisAlignmentCenter CrossAxisAlignment = iota
	// Place the children with their start edge aligned with the start side of
	// the cross axis.
	CrossAxisAlignmentStart
	// Place the children as close to the end of the cross axis as possible.
	CrossAxisAlignmentEnd
	// Place the children along the cross axis such that their baselines match.
	CrossAxisAlignmentBaseline
)
