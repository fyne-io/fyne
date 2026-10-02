#import <Cocoa/Cocoa.h>
#import "accessibility_darwin.h"
#include <stdio.h>
#include <string.h>

@interface ProbeBaseView : NSView
@end

@implementation ProbeBaseView
- (NSArray*)accessibilityChildren {
    return @[];
}
@end

@interface ProbeInheritedView : ProbeBaseView
@end

@implementation ProbeInheritedView
@end

@interface ProbeOwnView : ProbeBaseView
@end

@implementation ProbeOwnView
- (NSArray*)accessibilityChildren {
    return @[];
}
@end

@interface ProbeWindow : NSObject
@property (nonatomic, retain) NSView* contentView;
@end

@implementation ProbeWindow
- (void)dealloc {
    [_contentView release];
    [super dealloc];
}
@end

int main(int argc, const char* argv[]) {
    @autoreleasepool {
        id element = (id)AccessibilityElementCreate(
            AccessibilityRoleButton, "Button", "Button", 0, 0, 100, 30,
            NULL, NULL, NULL);
        AccessibilityElementSetTitle(element, "Updated button");
        AccessibilityElementSetLabel(element, "Updated label");
        AccessibilityElementSetFrame(element, 20, 40, 120, 50);
        if (![[element accessibilityTitle] isEqualToString:@"Updated button"] ||
            ![[element accessibilityLabel] isEqualToString:@"Updated label"] ||
            !NSEqualRects([[element valueForKey:@"localFrame"] rectValue], NSMakeRect(20, 40, 120, 50)) ||
            !NSEqualRects([element accessibilityFrame], NSMakeRect(20, 40, 120, 50))) {
            fprintf(stderr, "Native attributes were not updated in place\n");
            AccessibilityElementDestroy(element);
            return 1;
        }
        AccessibilityElementDestroy(element);

        id parent = (id)AccessibilityElementCreate(
            AccessibilityRoleGroup, "Group", "Group", 0, 0, 100, 30,
            NULL, NULL, NULL);
        id child = (id)AccessibilityElementCreate(
            AccessibilityRoleButton, "Button", "Button", 0, 0, 100, 30,
            NULL, NULL, NULL);
        [child retain];
        AccessibilityElementAddChild(parent, child);
        AccessibilityElementDestroy(parent);
        if ([child valueForKey:@"parentElement"] != nil) {
            fprintf(stderr, "Retained child refers to a destroyed parent\n");
            [child release];
            return 1;
        }
        [child release];

        if (argc > 1 && strcmp(argv[1], "attributes") != 0) {
            Class viewClass = strcmp(argv[1], "own-method") == 0 ?
                [ProbeOwnView class] : [ProbeInheritedView class];
            ProbeWindow* window = [[ProbeWindow alloc] init];
            NSView* view = [[viewClass alloc] initWithFrame:NSZeroRect];
            window.contentView = view;
            AccessibilitySetTargetWindow(window);
            element = (id)AccessibilityElementCreate(
                AccessibilityRoleButton, "Button", "Button", 0, 0, 100, 30,
                NULL, NULL, NULL);
            AccessibilityAttachToWindow(element);
            NSArray* children = [view accessibilityChildren];
            if ([children count] != 1 || [children firstObject] != element) {
                fprintf(stderr, "Native root was not attached to the content view\n");
                return 1;
            }
            AccessibilityElementDestroy(element);
            if ([[view accessibilityChildren] count] != 0) {
                fprintf(stderr, "Destroyed native root is still registered\n");
                return 1;
            }

            NSView* unrelated = [[ProbeBaseView alloc] initWithFrame:NSZeroRect];
            element = (id)AccessibilityElementCreate(
                AccessibilityRoleButton, "Button", "Button", 0, 0, 100, 30,
                NULL, NULL, NULL);
            AccessibilityAttachToWindow(element);
            if ([[unrelated accessibilityChildren] count] != 0) {
                fprintf(stderr, "Attaching a root changed its content view superclass\n");
                return 1;
            }
            AccessibilityElementDestroy(element);
            [unrelated release];
            [view release];
            [window release];
        }
    }
    return 0;
}
