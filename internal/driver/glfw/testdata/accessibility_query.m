#import <Cocoa/Cocoa.h>
#import <ApplicationServices/ApplicationServices.h>
#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <unistd.h>

static double milliseconds(void) {
    struct timespec now;
    clock_gettime(CLOCK_MONOTONIC, &now);
    return now.tv_sec * 1000.0 + now.tv_nsec / 1000000.0;
}

static int compare(const void* left, const void* right) {
    double a = *(const double*)left, b = *(const double*)right;
    return (a > b) - (a < b);
}

static void report(const char* name, double* times, int count) {
    qsort(times, count, sizeof(double), compare);
    printf("%s median_ms=%.3f p95_ms=%.3f max_ms=%.3f\n",
           name, times[count / 2], times[(count - 1) * 95 / 100], times[count - 1]);
}

static CFTypeRef attribute(AXUIElementRef element, CFStringRef name) {
    CFTypeRef value = NULL;
    AXError error = AXUIElementCopyAttributeValue(element, name, &value);
    if (error != kAXErrorSuccess) {
        fprintf(stderr, "AX attribute %s failed: %d\n", [(NSString*)name UTF8String], error);
        exit(1);
    }
    return value;
}

static AXUIElementRef findText(AXUIElementRef element) {
    CFTypeRef role = attribute(element, kAXRoleAttribute);
    BOOL text = CFEqual(role, kAXStaticTextRole);
    BOOL group = CFEqual(role, kAXGroupRole);
    CFRelease(role);
    if (text) {
        return (AXUIElementRef)CFRetain(element);
    }
    CFArrayRef children = NULL;
    AXError error = AXUIElementCopyAttributeValue(element, kAXChildrenAttribute, (CFTypeRef*)&children);
    if (error == kAXErrorAttributeUnsupported || error == kAXErrorNoValue) {
        return NULL;
    }
    if (error != kAXErrorSuccess) {
        fprintf(stderr, "Cannot traverse AX children: %d\n", error);
        exit(1);
    }
    AXUIElementRef found = NULL;
    CFIndex count = CFArrayGetCount(children);
    // Prefer the newest root but the first visible text within a group.
    for (CFIndex i = 0; i < count && !found; i++) {
        found = findText((AXUIElementRef)CFArrayGetValueAtIndex(children, group ? i : count - i - 1));
    }
    CFRelease(children);
    return found;
}

int main(int argc, char** argv) {
    @autoreleasepool {
        if (argc != 3) {
            fprintf(stderr, "usage: accessibility-query PID SAMPLES\n");
            return 1;
        }
        NSDictionary* options = @{(__bridge NSString*)kAXTrustedCheckOptionPrompt: @YES};
        if (!AXIsProcessTrustedWithOptions((__bridge CFDictionaryRef)options)) {
            fprintf(stderr, "Accessibility permission is required for the process running this probe.\n");
            return 2;
        }
        int count = atoi(argv[2]);
        if (count < 2 || count > 10000) {
            fprintf(stderr, "SAMPLES must be between 2 and 10000\n");
            return 1;
        }
        AXUIElementRef app = AXUIElementCreateApplication(atoi(argv[1]));
        CFArrayRef windows = attribute(app, kAXWindowsAttribute);
        if (CFArrayGetCount(windows) == 0) {
            fprintf(stderr, "The application has no accessible windows.\n");
            return 1;
        }
        AXUIElementRef window = (AXUIElementRef)CFArrayGetValueAtIndex(windows, 0);
        CFTypeRef windowRole = attribute(window, kAXRoleAttribute);
        BOOL validWindow = CFEqual(windowRole, kAXWindowRole);
        CFRelease(windowRole);
        if (!validWindow) {
            fprintf(stderr, "macOS returned a non-window AX element. Check the target application before measuring.\n");
            return 1;
        }
        AXUIElementRef retained = findText(window);
        if (!retained) {
            fprintf(stderr, "The window has no accessible text element.\n");
            return 1;
        }
        CFTypeRef position = attribute(retained, kAXPositionAttribute);
        CGPoint point;
        AXValueGetValue(position, kAXValueCGPointType, &point);
        CFRelease(position);
        point.x += 5;
        point.y += 5;
        AXUIElementRef selected = NULL;
        AXError error = AXUIElementCopyElementAtPosition(app, point.x, point.y, &selected);
        if (error != kAXErrorSuccess) {
            fprintf(stderr, "Initial AX hit test failed: %d\n", error);
            return 1;
        }
        AXUIElementRef parent = attribute(retained, kAXParentAttribute);
        AXUIElementRef view = attribute(parent, kAXParentAttribute);
        CFRelease(parent);
        CFArrayRef roots = attribute(view, kAXChildrenAttribute);
        CFIndex initialRoots = CFArrayGetCount(roots);
        CFRelease(roots);
        usleep(1500000);
        CFTypeRef selectedRole = NULL;
        error = AXUIElementCopyAttributeValue(selected, kAXRoleAttribute, &selectedRole);
        if (selectedRole) {
            printf("hit_role=%s\n", [(NSString*)selectedRole UTF8String]);
            CFRelease(selectedRole);
        } else {
            printf("initial_hit_error=%d\n", error);
        }
        double* queries = calloc(count, sizeof(double));
        double* hits = calloc(count, sizeof(double));
        if (!queries || !hits) {
            fprintf(stderr, "Cannot allocate timing samples.\n");
            return 1;
        }
        int changed = 0, invalid = 0;
        for (int i = 0; i < count; i++) {
            double start = milliseconds();
            CFTypeRef role = attribute(window, kAXRoleAttribute);
            queries[i] = milliseconds() - start;
            CFRelease(role);
            CFTypeRef selectedRole = NULL;
            error = AXUIElementCopyAttributeValue(retained, kAXRoleAttribute, &selectedRole);
            if (error != kAXErrorSuccess) {
                invalid++;
            }
            if (selectedRole) {
                CFRelease(selectedRole);
            }
            AXUIElementRef hit = NULL;
            start = milliseconds();
            error = AXUIElementCopyElementAtPosition(app, point.x, point.y, &hit);
            hits[i] = milliseconds() - start;
            if (error != kAXErrorSuccess) {
                fprintf(stderr, "AX hit test failed: %d\n", error);
                return 1;
            }
            changed += !CFEqual(selected, hit);
            CFRelease(hit);
            usleep(1000);
        }
        report("attribute", queries, count);
        report("hit_test", hits, count);
        AXUIElementRef latest = findText(window);
        if (!latest) {
            fprintf(stderr, "The accessible text disappeared during measurement.\n");
            return 1;
        }
        roots = attribute(view, kAXChildrenAttribute);
        printf("hit_changes=%d invalid_selections=%d text_identity_changed=%d roots_before=%ld roots_after=%ld samples=%d\n",
               changed, invalid, !CFEqual(retained, latest), (long)initialRoots, (long)CFArrayGetCount(roots), count);
        CFRelease(roots);
        CFRelease(view);
        CFRelease(latest);
        free(queries);
        free(hits);
        CFRelease(selected);
        CFRelease(retained);
        CFRelease(windows);
        CFRelease(app);
    }
    return 0;
}
