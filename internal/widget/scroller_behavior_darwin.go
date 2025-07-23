package widget

/*
int getScrollerPagingBehavior();
*/
import "C"

func getScrollerPagingBehavior() scrollBarTapBehavior {
	if C.getScrollerPagingBehavior() == 0 {
		return scrollBarTapBehaviorScrollOnePage
	}
	return scrollBarTapBehaviorScrollToPosition
}
