// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build darwin
// +build !ios

#include "_cgo_export.h"
#include <pthread.h>
#include <stdio.h>

#import <Cocoa/Cocoa.h>
#import <Foundation/Foundation.h>
#import <OpenGL/gl3.h>

void makeCurrentContext(GLintptr context) {
	NSOpenGLContext* ctx = (NSOpenGLContext*)context;
	[ctx makeCurrentContext];
}

uint64 threadID() {
	uint64 id;
	if (pthread_threadid_np(pthread_self(), &id)) {
		abort();
	}
	return id;
}

@interface MobileGLView : NSOpenGLView<NSApplicationDelegate, NSWindowDelegate>
{
}
@end

@implementation MobileGLView
- (void)prepareOpenGL {
	[super prepareOpenGL];
	[self setWantsBestResolutionOpenGLSurface:YES];
	GLint swapInt = 1;

#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wdeprecated-declarations"
	[[self openGLContext] setValues:&swapInt forParameter:NSOpenGLCPSwapInterval];
#pragma clang diagnostic pop

	// Using attribute arrays in OpenGL 3.3 requires the use of a VBA.
	// But VBAs don't exist in ES 2. So we bind a default one.
	GLuint vba;
	glGenVertexArrays(1, &vba);
	glBindVertexArray(vba);

	startloop((GLintptr)[self openGLContext]);
}

- (void)reshape {
	[super reshape];

	// Calculate screen PPI.
	//
	// Note that the backingScaleFactor converts from logical
	// pixels to actual pixels, but both of these units vary
	// independently from real world size. E.g.
	//
	// 13" Retina Macbook Pro, 2560x1600, 227ppi, backingScaleFactor=2, scale=3.15
	// 15" Retina Macbook Pro, 2880x1800, 220ppi, backingScaleFactor=2, scale=3.06
	// 27" iMac,               2560x1440, 109ppi, backingScaleFactor=1, scale=1.51
	// 27" Retina iMac,        5120x2880, 218ppi, backingScaleFactor=2, scale=3.03
	NSScreen *screen = [NSScreen mainScreen];
	double screenPixW = [screen frame].size.width * [screen backingScaleFactor];

	CGDirectDisplayID display = (CGDirectDisplayID)[[[screen deviceDescription] valueForKey:@"NSScreenNumber"] intValue];
	CGSize screenSizeMM = CGDisplayScreenSize(display); // in millimeters
	float ppi = 25.4 * screenPixW / screenSizeMM.width;
	float pixelsPerPt = ppi/72.0;

	// The width and height reported to the geom package are the
	// bounds of the OpenGL view. Several steps are necessary.
	// First, [self bounds] gives us the number of logical pixels
	// in the view. Multiplying this by the backingScaleFactor
	// gives us the number of actual pixels.
	NSRect r = [self bounds];
	int w = r.size.width * [screen backingScaleFactor];
	int h = r.size.height * [screen backingScaleFactor];

	setGeom(pixelsPerPt, w, h);
}

- (void)drawRect:(NSRect)theRect {
	// Called during resize. This gets rid of flicker when resizing.
	drawgl();
}

- (void)mouseDown:(NSEvent *)theEvent {
	double scale = [[NSScreen mainScreen] backingScaleFactor];
	NSPoint p = [theEvent locationInWindow];
	eventMouseDown(p.x * scale, p.y * scale);
}

- (void)mouseUp:(NSEvent *)theEvent {
	double scale = [[NSScreen mainScreen] backingScaleFactor];
	NSPoint p = [theEvent locationInWindow];
	eventMouseEnd(p.x * scale, p.y * scale);
}

- (void)mouseDragged:(NSEvent *)theEvent {
	double scale = [[NSScreen mainScreen] backingScaleFactor];
	NSPoint p = [theEvent locationInWindow];
	eventMouseDragged(p.x * scale, p.y * scale);
}

- (void)windowDidBecomeKey:(NSNotification *)notification {
	lifecycleFocused();
}

- (void)windowDidResignKey:(NSNotification *)notification {
	if (![NSApp isHidden]) {
		lifecycleVisible();
	}
}

- (void)applicationDidFinishLaunching:(NSNotification *)aNotification {
	lifecycleAlive();
	[[NSRunningApplication currentApplication] activateWithOptions:(NSApplicationActivateAllWindows | NSApplicationActivateIgnoringOtherApps)];
	[self.window makeKeyAndOrderFront:self];
	lifecycleVisible();
}

- (void)applicationWillTerminate:(NSNotification *)aNotification {
	lifecycleDead();
}

- (void)applicationDidHide:(NSNotification *)aNotification {
	lifecycleAlive();
}

- (void)applicationWillUnhide:(NSNotification *)notification {
	lifecycleVisible();
}

- (void)windowWillClose:(NSNotification *)notification {
	lifecycleAlive();
}

- (BOOL)acceptsFirstResponder {
    return true;
}

- (void)keyDown:(NSEvent *)theEvent {
	[self key:theEvent];
}
- (void)keyUp:(NSEvent *)theEvent {
	[self key:theEvent];
}
- (void)key:(NSEvent *)theEvent {
	NSRange range = [theEvent.characters rangeOfComposedCharacterSequenceAtIndex:0];

	uint8_t buf[4] = {0, 0, 0, 0};
	if (![theEvent.characters getBytes:buf
			maxLength:4
			usedLength:nil
			encoding:NSUTF32LittleEndianStringEncoding
			options:NSStringEncodingConversionAllowLossy
			range:range
			remainingRange:nil]) {
		NSLog(@"failed to read key event %@", theEvent);
		return;
	}

	uint32_t rune = (uint32_t)buf[0]<<0 | (uint32_t)buf[1]<<8 | (uint32_t)buf[2]<<16 | (uint32_t)buf[3]<<24;

	uint8_t direction;
	if ([theEvent isARepeat]) {
		direction = 0;
	} else if (theEvent.type == NSEventTypeKeyDown) {
		direction = 1;
	} else {
		direction = 2;
	}
	eventKey((int32_t)rune, direction, theEvent.keyCode, theEvent.modifierFlags);
}

- (void)flagsChanged:(NSEvent *)theEvent {
	eventFlags(theEvent.modifierFlags);
}
@end

void
runApp(void) {
	[NSAutoreleasePool new];
	[NSApplication sharedApplication];
	[NSApp setActivationPolicy:NSApplicationActivationPolicyRegular];

	id menuBar = [[NSMenu new] autorelease];
	id menuItem = [[NSMenuItem new] autorelease];
	[menuBar addItem:menuItem];
	[NSApp setMainMenu:menuBar];

	id menu = [[NSMenu new] autorelease];
	id name = [[NSProcessInfo processInfo] processName];

	id hideMenuItem = [[[NSMenuItem alloc] initWithTitle:@"Hide"
		action:@selector(hide:) keyEquivalent:@"h"]
		autorelease];
	[menu addItem:hideMenuItem];

	id quitMenuItem = [[[NSMenuItem alloc] initWithTitle:@"Quit"
		action:@selector(terminate:) keyEquivalent:@"q"]
		autorelease];
	[menu addItem:quitMenuItem];
	[menuItem setSubmenu:menu];

	NSRect rect = NSMakeRect(0, 0, 600, 800);

	NSWindow* window = [[[NSWindow alloc] initWithContentRect:rect
			styleMask:NSWindowStyleMaskTitled
			backing:NSBackingStoreBuffered
			defer:NO]
		autorelease];
	window.styleMask |= NSWindowStyleMaskResizable;
	window.styleMask |= NSWindowStyleMaskMiniaturizable;
	window.styleMask |= NSWindowStyleMaskClosable;
	window.title = name;
	[window cascadeTopLeftFromPoint:NSMakePoint(20,20)];

	NSOpenGLPixelFormatAttribute attr[] = {
		NSOpenGLPFAOpenGLProfile, NSOpenGLProfileVersion3_2Core,
		NSOpenGLPFAColorSize,     24,
		NSOpenGLPFAAlphaSize,     8,
		NSOpenGLPFADepthSize,     16,
		NSOpenGLPFAAccelerated,
		NSOpenGLPFADoubleBuffer,
		NSOpenGLPFAAllowOfflineRenderers,
		0
	};
	id pixFormat = [[NSOpenGLPixelFormat alloc] initWithAttributes:attr];
	MobileGLView* view = [[MobileGLView alloc] initWithFrame:rect pixelFormat:pixFormat];
	[window setContentView:view];
	[window setDelegate:view];
	[NSApp setDelegate:view];

	[NSApp run];
}

void stopApp(void) {
	[NSApp terminate:nil];
}
