// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build (linux && !android) || freebsd || openbsd
// +build linux,!android freebsd openbsd

#include "_cgo_export.h"
#include <EGL/egl.h>
#include <GLES2/gl2.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>
#include <stdio.h>
#include <stdlib.h>

static Atom wm_delete_window;

static Window
new_window(Display *x_dpy, EGLDisplay e_dpy, int w, int h, EGLContext *ctx, EGLSurface *surf) {
	static const EGLint attribs[] = {
		EGL_RENDERABLE_TYPE, EGL_OPENGL_ES2_BIT,
		EGL_SURFACE_TYPE, EGL_WINDOW_BIT,
		EGL_BLUE_SIZE, 8,
		EGL_GREEN_SIZE, 8,
		EGL_RED_SIZE, 8,
		EGL_DEPTH_SIZE, 16,
		EGL_CONFIG_CAVEAT, EGL_NONE,
		EGL_NONE
	};
	EGLConfig config;
	EGLint num_configs;
	if (!eglChooseConfig(e_dpy, attribs, &config, 1, &num_configs)) {
		fprintf(stderr, "eglChooseConfig failed\n");
		exit(1);
	}
	EGLint vid;
	if (!eglGetConfigAttrib(e_dpy, config, EGL_NATIVE_VISUAL_ID, &vid)) {
		fprintf(stderr, "eglGetConfigAttrib failed\n");
		exit(1);
	}

	XVisualInfo visTemplate;
	visTemplate.visualid = vid;
	int num_visuals;
	XVisualInfo *visInfo = XGetVisualInfo(x_dpy, VisualIDMask, &visTemplate, &num_visuals);
	if (!visInfo) {
		fprintf(stderr, "XGetVisualInfo failed\n");
		exit(1);
	}

	Window root = RootWindow(x_dpy, DefaultScreen(x_dpy));
	XSetWindowAttributes attr;

	attr.colormap = XCreateColormap(x_dpy, root, visInfo->visual, AllocNone);
	if (!attr.colormap) {
		fprintf(stderr, "XCreateColormap failed\n");
		exit(1);
	}

	attr.event_mask = StructureNotifyMask | ExposureMask |
		ButtonPressMask | ButtonReleaseMask | ButtonMotionMask;
	Window win = XCreateWindow(
		x_dpy, root, 0, 0, w, h, 0, visInfo->depth, InputOutput,
		visInfo->visual, CWColormap | CWEventMask, &attr);
	XFree(visInfo);

	XSizeHints sizehints;
	sizehints.width  = w;
	sizehints.height = h;
	sizehints.flags = USSize;
	XSetNormalHints(x_dpy, win, &sizehints);
	XSetStandardProperties(x_dpy, win, "App", "App", None, (char **)NULL, 0, &sizehints);

	static const EGLint ctx_attribs[] = {
		EGL_CONTEXT_CLIENT_VERSION, 2,
		EGL_NONE
	};
	*ctx = eglCreateContext(e_dpy, config, EGL_NO_CONTEXT, ctx_attribs);
	if (!*ctx) {
		fprintf(stderr, "eglCreateContext failed\n");
		exit(1);
	}
	*surf = eglCreateWindowSurface(e_dpy, config, win, NULL);
	if (!*surf) {
		fprintf(stderr, "eglCreateWindowSurface failed\n");
		exit(1);
	}
	return win;
}

Display *x_dpy;
EGLDisplay e_dpy;
EGLContext e_ctx;
EGLSurface e_surf;
Window win;

void
createWindow(void) {
	x_dpy = XOpenDisplay(NULL);
	if (!x_dpy) {
		fprintf(stderr, "XOpenDisplay failed\n");
		exit(1);
	}
	e_dpy = eglGetDisplay(x_dpy);
	if (!e_dpy) {
		fprintf(stderr, "eglGetDisplay failed\n");
		exit(1);
	}
	EGLint e_major, e_minor;
	if (!eglInitialize(e_dpy, &e_major, &e_minor)) {
		fprintf(stderr, "eglInitialize failed\n");
		exit(1);
	}
	eglBindAPI(EGL_OPENGL_ES_API);
	win = new_window(x_dpy, e_dpy, 600, 800, &e_ctx, &e_surf);

	wm_delete_window = XInternAtom(x_dpy, "WM_DELETE_WINDOW", True);
	if (wm_delete_window != None) {
		XSetWMProtocols(x_dpy, win, &wm_delete_window, 1);
	}

	XMapWindow(x_dpy, win);
	if (!eglMakeCurrent(e_dpy, e_surf, e_surf, e_ctx)) {
		fprintf(stderr, "eglMakeCurrent failed\n");
		exit(1);
	}

	// Window size and DPI should be initialized before starting app.
	XEvent ev;
	while (1) {
		if (XCheckMaskEvent(x_dpy, StructureNotifyMask, &ev) == False) {
			continue;
		}
		if (ev.type == ConfigureNotify) {
			onResize(ev.xconfigure.width, ev.xconfigure.height);
			break;
		}
	}
}

void
processEvents(void) {
	while (XPending(x_dpy)) {
		XEvent ev;
		XNextEvent(x_dpy, &ev);
		switch (ev.type) {
		case ButtonPress:
			onTouchBegin((float)ev.xbutton.x, (float)ev.xbutton.y);
			break;
		case ButtonRelease:
			onTouchEnd((float)ev.xbutton.x, (float)ev.xbutton.y);
			break;
		case MotionNotify:
			onTouchMove((float)ev.xmotion.x, (float)ev.xmotion.y);
			break;
		case ConfigureNotify:
			onResize(ev.xconfigure.width, ev.xconfigure.height);
			break;
		case ClientMessage:
			if (wm_delete_window != None && (Atom)ev.xclient.data.l[0] == wm_delete_window) {
				onStop();
				return;
			}
			break;
		}
	}
}

void
swapBuffers(void) {
	if (eglSwapBuffers(e_dpy, e_surf) == EGL_FALSE) {
		fprintf(stderr, "eglSwapBuffer failed\n");
		exit(1);
	}
}
