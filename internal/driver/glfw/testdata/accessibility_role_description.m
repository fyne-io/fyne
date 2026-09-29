#import <Cocoa/Cocoa.h>
#import "accessibility_darwin.h"
#include <stdio.h>

int main(void) {
    @autoreleasepool {
        for (AccessibilityRole role = AccessibilityRoleWindow; role <= AccessibilityRoleGroup; role++) {
            id element = (id)AccessibilityElementCreate(
                role, "Test", "Test", 0, 0, 100, 30, NULL, NULL, NULL, NULL);
            NSString *expected = role == AccessibilityRoleHeading ? @"heading" :
                NSAccessibilityRoleDescription([element accessibilityRole], [element accessibilitySubrole]);
            fprintf(stderr, "Reading role description for role %d\n", role);
            NSString *modern = [element accessibilityRoleDescription];
            NSString *legacy = [element accessibilityAttributeValue:NSAccessibilityRoleDescriptionAttribute];
            if (modern.length == 0 || ![modern isEqualToString:expected] || ![legacy isEqualToString:modern]) {
                fprintf(stderr, "Incorrect modern or legacy role description for role %d\n", role);
                AccessibilityElementDestroy(element);
                return 1;
            }
            AccessibilityElementDestroy(element);
        }
    }
    return 0;
}
