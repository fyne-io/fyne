package layout

import (
	"log"

	"fyne.io/fyne/v2"
)

type (
	// Axis represents the two cardinal directions in two dimensions.
	Axis int
	// MainAxisAlignment represents how the children should be placed along the main axis in a flex layout.
	MainAxisAlignment int
	// CrossAxisAlignment represents how the children should be placed along the cross axis in a flex layout.
	CrossAxisAlignment int
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

// NewRow creates a new FlexLayout instance with AxisHorizontal.
//
func NewRow(
	mainAxisAlignment MainAxisAlignment,
	crossAxisAlignment CrossAxisAlignment,
) fyne.Layout {
	return &flexLayout{
		Axis: AxisHorizontal, MainAxisAlignment: mainAxisAlignment, CrossAxisAlignment: crossAxisAlignment,
	}
}

// NewColumn creates a new FlexLayout instance with AxisVertical.
//
func NewColumn(
	mainAxisAlignment MainAxisAlignment,
	crossAxisAlignment CrossAxisAlignment,
) fyne.Layout {
	return &flexLayout{
		Axis: AxisVertical, MainAxisAlignment: mainAxisAlignment, CrossAxisAlignment: crossAxisAlignment,
	}
}

// Declare conformity with Layout interface
var _ fyne.Layout = (*flexLayout)(nil)

type flexLayout struct {
	Axis               Axis
	MainAxisAlignment  MainAxisAlignment
	CrossAxisAlignment CrossAxisAlignment
}

// Kind returns what kind of flex layout is this object based on Axis field.
func (l *flexLayout) kind() string {
	if l.Axis == AxisVertical {
		return "ColumnFlexLayout"
	}
	return "RowFlexLayout"
}

func (l *flexLayout) getFlex(obj fyne.CanvasObject) int {
	type flexibleobj interface{ Flex() int }
	flex := 0
	if f, ok := obj.(flexibleobj); ok {
		flex = f.Flex()
	}
	return flex
}

func (l *flexLayout) getCrossSize(obj fyne.CanvasObject, min bool) float32 {
	minSize := obj.MinSize()
	curSize := obj.Size()
	switch l.Axis {
	case AxisHorizontal:
		if min {
			return minSize.Height
		}
		return fyne.Max(curSize.Height, minSize.Height)
	case AxisVertical:
		if min {
			return minSize.Width
		}
		return fyne.Max(curSize.Width, minSize.Width)
	}
	if min {
		return minSize.Height
	}
	return fyne.Max(curSize.Height, minSize.Height)
}

func (l *flexLayout) getMainSize(obj fyne.CanvasObject, min bool) float32 {
	minSize := obj.MinSize()
	curSize := obj.Size()
	switch l.Axis {
	case AxisHorizontal:
		if min {
			return minSize.Width
		}
		return fyne.Max(curSize.Width, minSize.Width)
	case AxisVertical:
		if min {
			return minSize.Height
		}
		return fyne.Max(curSize.Height, minSize.Height)
	}
	if min {
		return minSize.Width
	}
	return fyne.Max(curSize.Width, minSize.Width)
}

func (l *flexLayout) getDistanceToBaseline(obj fyne.CanvasObject) float32 {
	// TODO check
	type baseliner interface{ DistanceToTextBaseline() float32 }
	distance := float32(0)
	if bl, ok := obj.(baseliner); ok {
		distance = bl.DistanceToTextBaseline()
	}
	return distance
}

// Layout is called to pack all child objects into a specified size.
//
func (l *flexLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	// Determine used flex factor, size inflexible items, calculate free space.
	totalFlex := 0
	maxMainSize := size.Width
	if l.Axis == AxisVertical {
		maxMainSize = size.Height
	}
	crossSize := float32(0)
	allocatedSize := float32(0) // Sum of the sizes of the non-flexible children.
	var lastFlexChild fyne.CanvasObject

	// --- Resizing non-flex objects
	for _, obj := range objects {
		flex := l.getFlex(obj)
		if flex > 0 {
			totalFlex += flex
			lastFlexChild = obj
		} else {
			// TODO mainAxisSize min and max??
			if _, ok := obj.(*fyne.Container); ok {
				switch l.Axis {
				case AxisHorizontal:
					obj.Resize(fyne.NewSize(l.getMainSize(obj, false), size.Height))
				case AxisVertical:
					obj.Resize(fyne.NewSize(size.Width, l.getMainSize(obj, false)))
				}
			} else {
				obj.Resize(obj.MinSize())
			}
			allocatedSize += l.getMainSize(obj, false)
			crossSize = fyne.Max(crossSize, l.getCrossSize(obj, false))
		}
	}

	// --- Resizing flex objects
	freeSpace := maxMainSize - allocatedSize
	allocatedFlexSpace := float32(0)
	maxBaselineDistance := float32(0)
	if totalFlex > 0 || l.CrossAxisAlignment == CrossAxisAlignmentBaseline {
		spacePerFlex := freeSpace / float32(totalFlex)
		maxSizeAboveBaseline := float32(0)
		maxSizeBelowBaseline := float32(0)

		for _, obj := range objects {
			flex := l.getFlex(obj)
			if flex > 0 {
				maxChildExtent := spacePerFlex * float32(flex)
				if obj == lastFlexChild {
					maxChildExtent = freeSpace - allocatedFlexSpace
				}
				switch l.Axis {
				case AxisHorizontal:
					obj.Resize(fyne.NewSize(maxChildExtent, l.getCrossSize(obj, false)))
				case AxisVertical:
					obj.Resize(fyne.NewSize(l.getCrossSize(obj, false), maxChildExtent))
				}
				objChildMainSize := l.getMainSize(obj, false)
				if objChildMainSize > maxChildExtent {
					log.Printf("WARNING: Widget %s can not be shrinked to %.1f (current size: %.1f)\n",
						"WidgetName", maxChildExtent, objChildMainSize)
				}
				allocatedSize += objChildMainSize
				allocatedFlexSpace += maxChildExtent
				crossSize = fyne.Max(crossSize, l.getCrossSize(obj, false))
			}
			if l.CrossAxisAlignment == CrossAxisAlignmentBaseline {
				distance := l.getDistanceToBaseline(obj)
				maxBaselineDistance = fyne.Max(maxBaselineDistance, distance)
				maxSizeAboveBaseline = fyne.Max(
					distance,
					maxSizeAboveBaseline,
				)
				maxSizeBelowBaseline = fyne.Max(
					fyne.Max(obj.Size().Height, obj.MinSize().Height)-distance,
					maxSizeBelowBaseline,
				)
				crossSize = fyne.Max(maxSizeAboveBaseline+maxSizeBelowBaseline, crossSize)
			}
		}
	}

	// --- MaxAxisAlignment
	actualSizeDelta := maxMainSize - allocatedSize
	if overflow := fyne.Max(0, -actualSizeDelta); overflow > 0 {
		log.Printf(`[WARNING OVERFLOW]:
 %s overflows by %.1f.
 This warning is commonly triggered when you wrap widgets with a Flexible widget and they cannot
 be shrinked beyond to its minimum size as required for the layout calculations.
 `, l.kind(), overflow)
	}
	remainingSpace := fyne.Max(0, actualSizeDelta)
	leadingSpace := float32(0)
	betweenSpace := float32(0)
	switch l.MainAxisAlignment {
	case MainAxisAlignmentStart:
		leadingSpace = 0
		betweenSpace = 0
	case MainAxisAlignmentEnd:
		leadingSpace = remainingSpace
		betweenSpace = 0
	case MainAxisAlignmentCenter:
		leadingSpace = remainingSpace / 2
		betweenSpace = 0
	case MainAxisAlignmentSpaceBetween:
		leadingSpace = 0
		betweenSpace = 0
		if len(objects) > 1 {
			betweenSpace = remainingSpace / float32(len(objects)-1)
		}
	case MainAxisAlignmentSpaceAround:
		betweenSpace = 0
		if len(objects) > 0 {
			betweenSpace = remainingSpace / float32(len(objects))
		}
		leadingSpace = betweenSpace / 2
	case MainAxisAlignmentSpaceEvenly:
		betweenSpace = 0
		if len(objects) > 0 {
			betweenSpace = remainingSpace / float32(len(objects)+1)
		}
		leadingSpace = betweenSpace
	}

	// --- Position objects
	objMainPos := leadingSpace
	for _, obj := range objects {
		objCrossPos := float32(0)
		switch l.CrossAxisAlignment {
		case CrossAxisAlignmentStart:
			objCrossPos = 0
		case CrossAxisAlignmentEnd:
			objCrossPos = crossSize - l.getCrossSize(obj, false)
		case CrossAxisAlignmentCenter:
			objCrossPos = crossSize/2 - l.getCrossSize(obj, false)/2
		case CrossAxisAlignmentBaseline:
			objCrossPos = 0
			if l.Axis == AxisHorizontal {
				distance := l.getDistanceToBaseline(obj)
				objCrossPos = maxBaselineDistance - distance
			}
		}
		switch l.Axis {
		case AxisHorizontal:
			obj.Move(fyne.NewPos(objMainPos, objCrossPos))
		case AxisVertical:
			obj.Move(fyne.NewPos(objCrossPos, objMainPos))
		}
		objMainPos += l.getMainSize(obj, false) + betweenSpace
	}
}

func (l *flexLayout) MinSize(objects []fyne.CanvasObject) fyne.Size {
	allocatedSize := float32(0)
	crossSize := float32(0)
	maxSizeAboveBaseline := float32(0)
	maxSizeBelowBaseline := float32(0)
	for _, obj := range objects {
		allocatedSize += l.getMainSize(obj, true)
		crossSize = fyne.Max(crossSize, l.getCrossSize(obj, true))
		if l.CrossAxisAlignment == CrossAxisAlignmentBaseline {
			distance := l.getDistanceToBaseline(obj)
			maxSizeAboveBaseline = fyne.Max(
				distance,
				maxSizeAboveBaseline,
			)
			maxSizeBelowBaseline = fyne.Max(
				fyne.Max(obj.Size().Height, obj.MinSize().Height)-distance,
				maxSizeBelowBaseline,
			)
			crossSize = fyne.Max(maxSizeAboveBaseline+maxSizeBelowBaseline, crossSize)
		}
	}
	if l.Axis == AxisHorizontal {
		return fyne.NewSize(allocatedSize, crossSize)
	}
	return fyne.NewSize(crossSize, allocatedSize)
}
