//go:build darwin

package widget

/*
int getScrollerPagingBehavior();
*/
import "C"

func isScrollerPageOnTap() bool {
	return C.getScrollerPagingBehavior() == 0
}
