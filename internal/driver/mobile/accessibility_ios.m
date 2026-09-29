//go:build ios

#import <UIKit/UIKit.h>

static NSMutableArray<UIAccessibilityElement*> *pendingElements = nil;

static UIView* getContainerView(void) {
    UIWindow *window = nil;
    NSArray<UIWindow*> *windows = [[UIApplication sharedApplication] windows];
    for (UIWindow *w in windows) {
        if (w.isKeyWindow) {
            window = w;
            break;
        }
    }
    if (window == nil && windows.count > 0) {
        window = windows[0];
    }
    if (window == nil) {
        return nil;
    }

    UIViewController *rootVC = window.rootViewController;
    if (rootVC == nil) {
        return nil;
    }

    return rootVC.view;
}

void clearAccessibilityNodesIOS(void) {
    @autoreleasepool {
        pendingElements = [[NSMutableArray alloc] init];
    }
}

void addAccessibilityNodeIOS(int role, const char *label,
    float x, float y, float width, float height) {
    @autoreleasepool {
        UIView *container = getContainerView();
        if (container == nil) {
            return;
        }

        UIAccessibilityElement *elem = [[UIAccessibilityElement alloc]
            initWithAccessibilityContainer:container];

        elem.accessibilityLabel = label ? [NSString stringWithUTF8String:label] : @"";

        // Map Fyne roles to UIAccessibilityTraits. Role values match the
        // const block in accessibility_ios.go.
        // 0=container, 1=button, 2=checkbox, 3=heading, 4=image, 5=link,
        // 6=list, 7=listItem, 8=progressBar, 9=radio, 10=separator,
        // 11=slider, 12=tab, 13=tabList, 14=table, 15=text, 16=textField,
        // 17=tree, 18=treeItem.
        UIAccessibilityTraits traits = UIAccessibilityTraitNone;
        switch (role) {
            case 1:  traits = UIAccessibilityTraitButton; break;
            case 2:  traits = UIAccessibilityTraitButton; break; // checkbox
            case 3:  traits = UIAccessibilityTraitHeader; break;
            case 4:  traits = UIAccessibilityTraitImage; break;
            case 5:  traits = UIAccessibilityTraitLink; break;
            case 8:  traits = UIAccessibilityTraitUpdatesFrequently; break;
            case 9:  traits = UIAccessibilityTraitButton; break; // radio
            case 11: traits = UIAccessibilityTraitAdjustable; break;
            case 12: traits = UIAccessibilityTraitButton; break; // tab
            case 15: traits = UIAccessibilityTraitStaticText; break;
            case 16: traits = UIAccessibilityTraitNone; break; // text field
            default: traits = UIAccessibilityTraitNone; break;
        }
        elem.accessibilityTraits = traits;

        // Convert from native pixel coordinates to point coordinates for the
        // accessibility frame, which UIKit expects in screen coordinates.
        CGFloat scale = [UIScreen mainScreen].nativeScale;
        CGRect frame = CGRectMake(x / scale, y / scale, width / scale, height / scale);
        elem.accessibilityFrame = UIAccessibilityConvertFrameToScreenCoordinates(frame, container);

        [pendingElements addObject:elem];
    }
}

void commitAccessibilityNodesIOS(void) {
    @autoreleasepool {
        UIView *container = getContainerView();
        if (container == nil) {
            return;
        }

        container.accessibilityElements = [pendingElements copy];
        UIAccessibilityPostNotification(UIAccessibilityLayoutChangedNotification, nil);
    }
}

void setupAccessibilityIOS(void) {
    // No-op: the container view is looked up dynamically each time.
}
