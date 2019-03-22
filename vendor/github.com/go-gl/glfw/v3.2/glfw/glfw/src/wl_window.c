//========================================================================
// GLFW 3.2 Wayland - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2014 Jonas Ådahl <jadahl@gmail.com>
//
// This software is provided 'as-is', without any express or implied
// warranty. In no event will the authors be held liable for any damages
// arising from the use of this software.
//
// Permission is granted to anyone to use this software for any purpose,
// including commercial applications, and to alter it and redistribute it
// freely, subject to the following restrictions:
//
// 1. The origin of this software must not be misrepresented; you must not
//    claim that you wrote the original software. If you use this software
//    in a product, an acknowledgment in the product documentation would
//    be appreciated but is not required.
//
// 2. Altered source versions must be plainly marked as such, and must not
//    be misrepresented as being the original software.
//
// 3. This notice may not be removed or altered from any source
//    distribution.
//
//========================================================================

#include "internal.h"

#include <stdio.h>
#include <stdlib.h>
#include <errno.h>
#include <unistd.h>
#include <string.h>
#include <fcntl.h>
#include <sys/mman.h>
#include <poll.h>

#include <wayland-egl.h>
#include <wayland-cursor.h>


static void handlePing(void* data,
                       struct wl_shell_surface* shellSurface,
                       uint32_t serial)
{
    wl_shell_surface_pong(shellSurface, serial);
}

static void handleConfigure(void* data,
                            struct wl_shell_surface* shellSurface,
                            uint32_t edges,
                            int32_t width,
                            int32_t height)
{
    _GLFWwindow* window = data;
    float aspectRatio;
    float targetRatio;

    if (!window->monitor)
    {
        if (window->numer != GLFW_DONT_CARE && window->denom != GLFW_DONT_CARE)
        {
            aspectRatio = (float)width / (float)height;
            targetRatio = (float)window->numer / (float)window->denom;
            if (aspectRatio < targetRatio)
                height = width / targetRatio;
            else if (aspectRatio > targetRatio)
                width = height * targetRatio;
        }

        if (window->minwidth != GLFW_DONT_CARE && width < window->minwidth)
            width = window->minwidth;
        else if (window->maxwidth != GLFW_DONT_CARE && width > window->maxwidth)
            width = window->maxwidth;

        if (window->minheight != GLFW_DONT_CARE && height < window->minheight)
            height = window->minheight;
        else if (window->maxheight != GLFW_DONT_CARE && height > window->maxheight)
            height = window->maxheight;
    }

    _glfwInputWindowSize(window, width, height);
    _glfwPlatformSetWindowSize(window, width, height);
    _glfwInputWindowDamage(window);
}

static void handlePopupDone(void* data,
                            struct wl_shell_surface* shellSurface)
{
}

static const struct wl_shell_surface_listener shellSurfaceListener = {
    handlePing,
    handleConfigure,
    handlePopupDone
};

static void checkScaleChange(_GLFWwindow* window)
{
    int scaledWidth, scaledHeight;
    int scale = 1;
    int i;
    int monitorScale;

    // Check if we will be able to set the buffer scale or not.
    if (_glfw.wl.wl_compositor_version < 3)
        return;

    // Get the scale factor from the highest scale monitor.
    for (i = 0; i < window->wl.monitorsCount; ++i)
    {
        monitorScale = window->wl.monitors[i]->wl.scale;
        if (scale < monitorScale)
            scale = monitorScale;
    }

    // Only change the framebuffer size if the scale changed.
    if (scale != window->wl.scale)
    {
        window->wl.scale = scale;
        scaledWidth = window->wl.width * scale;
        scaledHeight = window->wl.height * scale;
        wl_surface_set_buffer_scale(window->wl.surface, scale);
        wl_egl_window_resize(window->wl.native, scaledWidth, scaledHeight, 0, 0);
        _glfwInputFramebufferSize(window, scaledWidth, scaledHeight);
    }
}

static void handleEnter(void *data,
                        struct wl_surface *surface,
                        struct wl_output *output)
{
    _GLFWwindow* window = data;
    _GLFWmonitor* monitor = wl_output_get_user_data(output);

    if (window->wl.monitorsCount + 1 > window->wl.monitorsSize)
    {
        ++window->wl.monitorsSize;
        window->wl.monitors =
            realloc(window->wl.monitors,
                    window->wl.monitorsSize * sizeof(_GLFWmonitor*));
    }

    window->wl.monitors[window->wl.monitorsCount++] = monitor;

    checkScaleChange(window);
}

static void handleLeave(void *data,
                        struct wl_surface *surface,
                        struct wl_output *output)
{
    _GLFWwindow* window = data;
    _GLFWmonitor* monitor = wl_output_get_user_data(output);
    GLFWbool found;
    int i;

    for (i = 0, found = GLFW_FALSE; i < window->wl.monitorsCount - 1; ++i)
    {
        if (monitor == window->wl.monitors[i])
            found = GLFW_TRUE;
        if (found)
            window->wl.monitors[i] = window->wl.monitors[i + 1];
    }
    window->wl.monitors[--window->wl.monitorsCount] = NULL;

    checkScaleChange(window);
}

static const struct wl_surface_listener surfaceListener = {
    handleEnter,
    handleLeave
};

// Makes the surface considered as XRGB instead of ARGB.
static void setOpaqueRegion(_GLFWwindow* window)
{
    struct wl_region* region;

    region = wl_compositor_create_region(_glfw.wl.compositor);
    if (!region)
        return;

    wl_region_add(region, 0, 0, window->wl.width, window->wl.height);
    wl_surface_set_opaque_region(window->wl.surface, region);
    wl_surface_commit(window->wl.surface);
    wl_region_destroy(region);
}

static GLFWbool createSurface(_GLFWwindow* window,
                              const _GLFWwndconfig* wndconfig)
{
    window->wl.surface = wl_compositor_create_surface(_glfw.wl.compositor);
    if (!window->wl.surface)
        return GLFW_FALSE;

    wl_surface_add_listener(window->wl.surface,
                            &surfaceListener,
                            window);

    wl_surface_set_user_data(window->wl.surface, window);

    window->wl.native = wl_egl_window_create(window->wl.surface,
                                             wndconfig->width,
                                             wndconfig->height);
    if (!window->wl.native)
        return GLFW_FALSE;

    window->wl.width = wndconfig->width;
    window->wl.height = wndconfig->height;
    window->wl.scale = 1;

    // TODO: make this optional once issue #197 is fixed.
    setOpaqueRegion(window);

    return GLFW_TRUE;
}

static GLFWbool createShellSurface(_GLFWwindow* window)
{
    window->wl.shell_surface = wl_shell_get_shell_surface(_glfw.wl.shell,
                                                          window->wl.surface);
    if (!window->wl.shell_surface)
        return GLFW_FALSE;

    wl_shell_surface_add_listener(window->wl.shell_surface,
                                  &shellSurfaceListener,
                                  window);

    if (window->wl.title)
        wl_shell_surface_set_title(window->wl.shell_surface, window->wl.title);

    if (window->monitor)
    {
        wl_shell_surface_set_fullscreen(
            window->wl.shell_surface,
            WL_SHELL_SURFACE_FULLSCREEN_METHOD_DEFAULT,
            0,
            window->monitor->wl.output);
    }
    else if (window->wl.maximized)
    {
        wl_shell_surface_set_maximized(window->wl.shell_surface, NULL);
    }
    else
    {
        wl_shell_surface_set_toplevel(window->wl.shell_surface);
    }

    return GLFW_TRUE;
}

static int
createTmpfileCloexec(char* tmpname)
{
    int fd;

    fd = mkostemp(tmpname, O_CLOEXEC);
    if (fd >= 0)
        unlink(tmpname);

    return fd;
}

static void
handleEvents(int timeout)
{
    struct wl_display* display = _glfw.wl.display;
    struct pollfd fds[] = {
        { wl_display_get_fd(display), POLLIN },
    };

    while (wl_display_prepare_read(display) != 0)
        wl_display_dispatch_pending(display);

    // If an error different from EAGAIN happens, we have likely been
    // disconnected from the Wayland session, try to handle that the best we
    // can.
    if (wl_display_flush(display) < 0 && errno != EAGAIN)
    {
        _GLFWwindow* window = _glfw.windowListHead;
        while (window)
        {
            _glfwInputWindowCloseRequest(window);
            window = window->next;
        }
        wl_display_cancel_read(display);
        return;
    }

    if (poll(fds, 1, timeout) > 0)
    {
        wl_display_read_events(display);
        wl_display_dispatch_pending(display);
    }
    else
    {
        wl_display_cancel_read(display);
    }
}

/*
 * Create a new, unique, anonymous file of the given size, and
 * return the file descriptor for it. The file descriptor is set
 * CLOEXEC. The file is immediately suitable for mmap()'ing
 * the given size at offset zero.
 *
 * The file should not have a permanent backing store like a disk,
 * but may have if XDG_RUNTIME_DIR is not properly implemented in OS.
 *
 * The file name is deleted from the file system.
 *
 * The file is suitable for buffer sharing between processes by
 * transmitting the file descriptor over Unix sockets using the
 * SCM_RIGHTS methods.
 *
 * posix_fallocate() is used to guarantee that disk space is available
 * for the file at the given size. If disk space is insufficent, errno
 * is set to ENOSPC. If posix_fallocate() is not supported, program may
 * receive SIGBUS on accessing mmap()'ed file contents instead.
 */
int
createAnonymousFile(off_t size)
{
    static const char template[] = "/glfw-shared-XXXXXX";
    const char* path;
    char* name;
    int fd;
    int ret;

    path = getenv("XDG_RUNTIME_DIR");
    if (!path)
    {
        errno = ENOENT;
        return -1;
    }

    name = calloc(strlen(path) + sizeof(template), 1);
    strcpy(name, path);
    strcat(name, template);

    fd = createTmpfileCloexec(name);

    free(name);

    if (fd < 0)
        return -1;
    ret = posix_fallocate(fd, 0, size);
    if (ret != 0)
    {
        close(fd);
        errno = ret;
        return -1;
    }
    return fd;
}

// Translates a GLFW standard cursor to a theme cursor name
//
static char *translateCursorShape(int shape)
{
    switch (shape)
    {
        case GLFW_ARROW_CURSOR:
            return "left_ptr";
        case GLFW_IBEAM_CURSOR:
            return "xterm";
        case GLFW_CROSSHAIR_CURSOR:
            return "crosshair";
        case GLFW_HAND_CURSOR:
            return "grabbing";
        case GLFW_HRESIZE_CURSOR:
            return "sb_h_double_arrow";
        case GLFW_VRESIZE_CURSOR:
            return "sb_v_double_arrow";
    }
    return NULL;
}

//////////////////////////////////////////////////////////////////////////
//////                       GLFW platform API                      //////
//////////////////////////////////////////////////////////////////////////

int _glfwPlatformCreateWindow(_GLFWwindow* window,
                              const _GLFWwndconfig* wndconfig,
                              const _GLFWctxconfig* ctxconfig,
                              const _GLFWfbconfig* fbconfig)
{
    if (!createSurface(window, wndconfig))
        return GLFW_FALSE;

    if (ctxconfig->client != GLFW_NO_API)
    {
        if (!_glfwInitEGL())
            return GLFW_FALSE;
        if (!_glfwCreateContextEGL(window, ctxconfig, fbconfig))
            return GLFW_FALSE;
    }

    if (wndconfig->title)
        window->wl.title = strdup(wndconfig->title);

    if (wndconfig->visible)
    {
        if (!createShellSurface(window))
            return GLFW_FALSE;

        window->wl.visible = GLFW_TRUE;
    }
    else
    {
        window->wl.shell_surface = NULL;
        window->wl.visible = GLFW_FALSE;
    }

    window->wl.currentCursor = NULL;

    window->wl.monitors = calloc(1, sizeof(_GLFWmonitor*));
    window->wl.monitorsCount = 0;
    window->wl.monitorsSize = 1;

    return GLFW_TRUE;
}

void _glfwPlatformDestroyWindow(_GLFWwindow* window)
{
    if (window == _glfw.wl.pointerFocus)
    {
        _glfw.wl.pointerFocus = NULL;
        _glfwInputCursorEnter(window, GLFW_FALSE);
    }
    if (window == _glfw.wl.keyboardFocus)
    {
        _glfw.wl.keyboardFocus = NULL;
        _glfwInputWindowFocus(window, GLFW_FALSE);
    }

    if (window->context.destroy)
        window->context.destroy(window);

    if (window->wl.native)
        wl_egl_window_destroy(window->wl.native);

    if (window->wl.shell_surface)
        wl_shell_surface_destroy(window->wl.shell_surface);

    if (window->wl.surface)
        wl_surface_destroy(window->wl.surface);

    free(window->wl.title);
    free(window->wl.monitors);
}

void _glfwPlatformSetWindowTitle(_GLFWwindow* window, const char* title)
{
    if (window->wl.title)
        free(window->wl.title);
    window->wl.title = strdup(title);
    if (window->wl.shell_surface)
        wl_shell_surface_set_title(window->wl.shell_surface, title);
}

void _glfwPlatformSetWindowIcon(_GLFWwindow* window,
                                int count, const GLFWimage* images)
{
    _glfwInputError(GLFW_PLATFORM_ERROR,
                    "Wayland: Setting window icon not supported");
}

void _glfwPlatformGetWindowPos(_GLFWwindow* window, int* xpos, int* ypos)
{
    // A Wayland client is not aware of its position, so just warn and leave it
    // as (0, 0)

    _glfwInputError(GLFW_PLATFORM_ERROR,
                    "Wayland: Window position retrieval not supported");
}

void _glfwPlatformSetWindowPos(_GLFWwindow* window, int xpos, int ypos)
{
    // A Wayland client can not set its position, so just warn

    _glfwInputError(GLFW_PLATFORM_ERROR,
                    "Wayland: Window position setting not supported");
}

void _glfwPlatformGetWindowSize(_GLFWwindow* window, int* width, int* height)
{
    if (width)
        *width = window->wl.width;
    if (height)
        *height = window->wl.height;
}

void _glfwPlatformSetWindowSize(_GLFWwindow* window, int width, int height)
{
    int scaledWidth = width * window->wl.scale;
    int scaledHeight = height * window->wl.scale;
    window->wl.width = width;
    window->wl.height = height;
    wl_egl_window_resize(window->wl.native, scaledWidth, scaledHeight, 0, 0);
    setOpaqueRegion(window);
    _glfwInputFramebufferSize(window, scaledWidth, scaledHeight);
}

void _glfwPlatformSetWindowSizeLimits(_GLFWwindow* window,
                                      int minwidth, int minheight,
                                      int maxwidth, int maxheight)
{
    // TODO: find out how to trigger a resize.
    // The actual limits are checked in the wl_shell_surface::configure handler.
}

void _glfwPlatformSetWindowAspectRatio(_GLFWwindow* window, int numer, int denom)
{
    // TODO: find out how to trigger a resize.
    // The actual limits are checked in the wl_shell_surface::configure handler.
}

void _glfwPlatformGetFramebufferSize(_GLFWwindow* window, int* width, int* height)
{
    _glfwPlatformGetWindowSize(window, width, height);
    *width *= window->wl.scale;
    *height *= window->wl.scale;
}

void _glfwPlatformGetWindowFrameSize(_GLFWwindow* window,
                                     int* left, int* top,
                                     int* right, int* bottom)
{
    // TODO: will need a proper implementation once decorations are
    // implemented, but for now just leave everything as 0.
}

void _glfwPlatformIconifyWindow(_GLFWwindow* window)
{
    // TODO: move to xdg_shell instead of wl_shell.
    _glfwInputError(GLFW_PLATFORM_ERROR,
                    "Wayland: Iconify window not supported");
}

void _glfwPlatformRestoreWindow(_GLFWwindow* window)
{
    // TODO: also do the same for iconified.
    if (window->monitor || window->wl.maximized)
    {
        if (window->wl.shell_surface)
            wl_shell_surface_set_toplevel(window->wl.shell_surface);

        window->wl.maximized = GLFW_FALSE;
    }
}

void _glfwPlatformMaximizeWindow(_GLFWwindow* window)
{
    if (!window->monitor && !window->wl.maximized)
    {
        if (window->wl.shell_surface)
        {
            // Let the compositor select the best output.
            wl_shell_surface_set_maximized(window->wl.shell_surface, NULL);
        }
        window->wl.maximized = GLFW_TRUE;
    }
}

void _glfwPlatformShowWindow(_GLFWwindow* window)
{
    if (!window->monitor)
    {
        if (!window->wl.shell_surface)
            createShellSurface(window);
        window->wl.visible = GLFW_TRUE;
    }
}

void _glfwPlatformHideWindow(_GLFWwindow* window)
{
    if (!window->monitor)
    {
        if (window->wl.shell_surface)
            wl_shell_surface_destroy(window->wl.shell_surface);
        window->wl.visible = GLFW_FALSE;
    }
}

void _glfwPlatformFocusWindow(_GLFWwindow* window)
{
    _glfwInputError(GLFW_PLATFORM_ERROR,
                    "Wayland: Focusing a window requires user interaction");
}

void _glfwPlatformSetWindowMonitor(_GLFWwindow* window,
                                   _GLFWmonitor* monitor,
                                   int xpos, int ypos,
                                   int width, int height,
                                   int refreshRate)
{
    if (monitor)
    {
        wl_shell_surface_set_fullscreen(
            window->wl.shell_surface,
            WL_SHELL_SURFACE_FULLSCREEN_METHOD_DEFAULT,
            refreshRate * 1000, // Convert Hz to mHz.
            monitor->wl.output);
    }
    else
    {
        wl_shell_surface_set_toplevel(window->wl.shell_surface);
    }
    _glfwInputWindowMonitorChange(window, monitor);
}

int _glfwPlatformWindowFocused(_GLFWwindow* window)
{
    return _glfw.wl.keyboardFocus == window;
}

int _glfwPlatformWindowIconified(_GLFWwindow* window)
{
    // TODO: move to xdg_shell, wl_shell doesn't have any iconified concept.
    return GLFW_FALSE;
}

int _glfwPlatformWindowVisible(_GLFWwindow* window)
{
    return window->wl.visible;
}

int _glfwPlatformWindowMaximized(_GLFWwindow* window)
{
    return window->wl.maximized;
}

void _glfwPlatformPollEvents(void)
{
    handleEvents(0);
}

void _glfwPlatformWaitEvents(void)
{
    handleEvents(-1);
}

void _glfwPlatformWaitEventsTimeout(double timeout)
{
    handleEvents((int) (timeout * 1e3));
}

void _glfwPlatformPostEmptyEvent(void)
{
    wl_display_sync(_glfw.wl.display);
}

void _glfwPlatformGetCursorPos(_GLFWwindow* window, double* xpos, double* ypos)
{
    if (xpos)
        *xpos = window->wl.cursorPosX;
    if (ypos)
        *ypos = window->wl.cursorPosY;
}

static GLFWbool isPointerLocked(_GLFWwindow* window);

void _glfwPlatformSetCursorPos(_GLFWwindow* window, double x, double y)
{
    if (isPointerLocked(window))
    {
        zwp_locked_pointer_v1_set_cursor_position_hint(
            window->wl.pointerLock.lockedPointer,
            wl_fixed_from_double(x), wl_fixed_from_double(y));
        wl_surface_commit(window->wl.surface);
    }
}

void _glfwPlatformSetCursorMode(_GLFWwindow* window, int mode)
{
    _glfwPlatformSetCursor(window, window->wl.currentCursor);
}

const char* _glfwPlatformGetKeyName(int key, int scancode)
{
    // TODO
    return NULL;
}

int _glfwPlatformCreateCursor(_GLFWcursor* cursor,
                              const GLFWimage* image,
                              int xhot, int yhot)
{
    struct wl_shm_pool* pool;
    int stride = image->width * 4;
    int length = image->width * image->height * 4;
    void* data;
    int fd, i;

    fd = createAnonymousFile(length);
    if (fd < 0)
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "Wayland: Creating a buffer file for %d B failed: %m",
                        length);
        return GLFW_FALSE;
    }

    data = mmap(NULL, length, PROT_READ | PROT_WRITE, MAP_SHARED, fd, 0);
    if (data == MAP_FAILED)
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "Wayland: Cursor mmap failed: %m");
        close(fd);
        return GLFW_FALSE;
    }

    pool = wl_shm_create_pool(_glfw.wl.shm, fd, length);

    close(fd);
    unsigned char* source = (unsigned char*) image->pixels;
    unsigned char* target = data;
    for (i = 0;  i < image->width * image->height;  i++, source += 4)
    {
        unsigned int alpha = source[3];

        *target++ = (unsigned char) ((source[2] * alpha) / 255);
        *target++ = (unsigned char) ((source[1] * alpha) / 255);
        *target++ = (unsigned char) ((source[0] * alpha) / 255);
        *target++ = (unsigned char) alpha;
    }

    cursor->wl.buffer =
        wl_shm_pool_create_buffer(pool, 0,
                                  image->width,
                                  image->height,
                                  stride, WL_SHM_FORMAT_ARGB8888);
    munmap(data, length);
    wl_shm_pool_destroy(pool);

    cursor->wl.width = image->width;
    cursor->wl.height = image->height;
    cursor->wl.xhot = xhot;
    cursor->wl.yhot = yhot;
    return GLFW_TRUE;
}

int _glfwPlatformCreateStandardCursor(_GLFWcursor* cursor, int shape)
{
    struct wl_cursor* standardCursor;

    standardCursor = wl_cursor_theme_get_cursor(_glfw.wl.cursorTheme,
                                                translateCursorShape(shape));
    if (!standardCursor)
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "Wayland: Standard cursor \"%s\" not found",
                        translateCursorShape(shape));
        return GLFW_FALSE;
    }

    cursor->wl.image = standardCursor->images[0];
    return GLFW_TRUE;
}

void _glfwPlatformDestroyCursor(_GLFWcursor* cursor)
{
    // If it's a standard cursor we don't need to do anything here
    if (cursor->wl.image)
        return;

    if (cursor->wl.buffer)
        wl_buffer_destroy(cursor->wl.buffer);
}

static void handleRelativeMotion(void* data,
                                 struct zwp_relative_pointer_v1* pointer,
                                 uint32_t timeHi,
                                 uint32_t timeLo,
                                 wl_fixed_t dx,
                                 wl_fixed_t dy,
                                 wl_fixed_t dxUnaccel,
                                 wl_fixed_t dyUnaccel)
{
    _GLFWwindow* window = data;

    if (window->cursorMode != GLFW_CURSOR_DISABLED)
        return;

    _glfwInputCursorPos(window,
                        window->virtualCursorPosX + wl_fixed_to_double(dxUnaccel),
                        window->virtualCursorPosY + wl_fixed_to_double(dyUnaccel));
}

static const struct zwp_relative_pointer_v1_listener relativePointerListener = {
    handleRelativeMotion
};

static void handleLocked(void* data,
                         struct zwp_locked_pointer_v1* lockedPointer)
{
}

static void unlockPointer(_GLFWwindow* window)
{
    struct zwp_relative_pointer_v1* relativePointer =
        window->wl.pointerLock.relativePointer;
    struct zwp_locked_pointer_v1* lockedPointer =
        window->wl.pointerLock.lockedPointer;

    zwp_relative_pointer_v1_destroy(relativePointer);
    zwp_locked_pointer_v1_destroy(lockedPointer);

    window->wl.pointerLock.relativePointer = NULL;
    window->wl.pointerLock.lockedPointer = NULL;
}

static void lockPointer(_GLFWwindow* window);

static void handleUnlocked(void* data,
                           struct zwp_locked_pointer_v1* lockedPointer)
{
}

static const struct zwp_locked_pointer_v1_listener lockedPointerListener = {
    handleLocked,
    handleUnlocked
};

static void lockPointer(_GLFWwindow* window)
{
    struct zwp_relative_pointer_v1* relativePointer;
    struct zwp_locked_pointer_v1* lockedPointer;

    if (!_glfw.wl.relativePointerManager)
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "Wayland: no relative pointer manager");
        return;
    }

    relativePointer =
        zwp_relative_pointer_manager_v1_get_relative_pointer(
            _glfw.wl.relativePointerManager,
            _glfw.wl.pointer);
    zwp_relative_pointer_v1_add_listener(relativePointer,
                                         &relativePointerListener,
                                         window);

    lockedPointer =
        zwp_pointer_constraints_v1_lock_pointer(
            _glfw.wl.pointerConstraints,
            window->wl.surface,
            _glfw.wl.pointer,
            NULL,
            ZWP_POINTER_CONSTRAINTS_V1_LIFETIME_PERSISTENT);
    zwp_locked_pointer_v1_add_listener(lockedPointer,
                                       &lockedPointerListener,
                                       window);

    window->wl.pointerLock.relativePointer = relativePointer;
    window->wl.pointerLock.lockedPointer = lockedPointer;

    wl_pointer_set_cursor(_glfw.wl.pointer, _glfw.wl.pointerSerial,
                          NULL, 0, 0);
}

static GLFWbool isPointerLocked(_GLFWwindow* window)
{
    return window->wl.pointerLock.lockedPointer != NULL;
}

void _glfwPlatformSetCursor(_GLFWwindow* window, _GLFWcursor* cursor)
{
    struct wl_buffer* buffer;
    struct wl_cursor* defaultCursor;
    struct wl_cursor_image* image;
    struct wl_surface* surface = _glfw.wl.cursorSurface;

    if (!_glfw.wl.pointer)
        return;

    window->wl.currentCursor = cursor;

    // If we're not in the correct window just save the cursor
    // the next time the pointer enters the window the cursor will change
    if (window != _glfw.wl.pointerFocus)
        return;

    // Unlock possible pointer lock if no longer disabled.
    if (window->cursorMode != GLFW_CURSOR_DISABLED && isPointerLocked(window))
        unlockPointer(window);

    if (window->cursorMode == GLFW_CURSOR_NORMAL)
    {
        if (cursor)
            image = cursor->wl.image;
        else
        {
            defaultCursor = wl_cursor_theme_get_cursor(_glfw.wl.cursorTheme,
                                                       "left_ptr");
            if (!defaultCursor)
            {
                _glfwInputError(GLFW_PLATFORM_ERROR,
                                "Wayland: Standard cursor not found");
                return;
            }
            image = defaultCursor->images[0];
        }

        if (image)
        {
            buffer = wl_cursor_image_get_buffer(image);
            if (!buffer)
                return;
            wl_pointer_set_cursor(_glfw.wl.pointer, _glfw.wl.pointerSerial,
                                  surface,
                                  image->hotspot_x,
                                  image->hotspot_y);
            wl_surface_attach(surface, buffer, 0, 0);
            wl_surface_damage(surface, 0, 0,
                              image->width, image->height);
            wl_surface_commit(surface);
        }
        else
        {
            wl_pointer_set_cursor(_glfw.wl.pointer, _glfw.wl.pointerSerial,
                                  surface,
                                  cursor->wl.xhot,
                                  cursor->wl.yhot);
            wl_surface_attach(surface, cursor->wl.buffer, 0, 0);
            wl_surface_damage(surface, 0, 0,
                              cursor->wl.width, cursor->wl.height);
            wl_surface_commit(surface);
        }
    }
    else if (window->cursorMode == GLFW_CURSOR_DISABLED)
    {
        if (!isPointerLocked(window))
            lockPointer(window);
    }
    else if (window->cursorMode == GLFW_CURSOR_HIDDEN)
    {
        wl_pointer_set_cursor(_glfw.wl.pointer, _glfw.wl.pointerSerial,
                              NULL, 0, 0);
    }
}

void _glfwPlatformSetClipboardString(_GLFWwindow* window, const char* string)
{
    // TODO
    _glfwInputError(GLFW_PLATFORM_ERROR,
                    "Wayland: Clipboard setting not implemented yet");
}

const char* _glfwPlatformGetClipboardString(_GLFWwindow* window)
{
    // TODO
    _glfwInputError(GLFW_PLATFORM_ERROR,
                    "Wayland: Clipboard getting not implemented yet");
    return NULL;
}

char** _glfwPlatformGetRequiredInstanceExtensions(uint32_t* count)
{
    char** extensions;

    *count = 0;

    if (!_glfw.vk.KHR_wayland_surface)
        return NULL;

    extensions = calloc(2, sizeof(char*));
    extensions[0] = strdup("VK_KHR_surface");
    extensions[1] = strdup("VK_KHR_wayland_surface");

    *count = 2;
    return extensions;
}

int _glfwPlatformGetPhysicalDevicePresentationSupport(VkInstance instance,
                                                      VkPhysicalDevice device,
                                                      uint32_t queuefamily)
{
    PFN_vkGetPhysicalDeviceWaylandPresentationSupportKHR vkGetPhysicalDeviceWaylandPresentationSupportKHR =
        (PFN_vkGetPhysicalDeviceWaylandPresentationSupportKHR)
        vkGetInstanceProcAddr(instance, "vkGetPhysicalDeviceWaylandPresentationSupportKHR");
    if (!vkGetPhysicalDeviceWaylandPresentationSupportKHR)
    {
        _glfwInputError(GLFW_API_UNAVAILABLE,
                        "Wayland: Vulkan instance missing VK_KHR_wayland_surface extension");
        return VK_NULL_HANDLE;
    }

    return vkGetPhysicalDeviceWaylandPresentationSupportKHR(device,
                                                            queuefamily,
                                                            _glfw.wl.display);
}

VkResult _glfwPlatformCreateWindowSurface(VkInstance instance,
                                          _GLFWwindow* window,
                                          const VkAllocationCallbacks* allocator,
                                          VkSurfaceKHR* surface)
{
    VkResult err;
    VkWaylandSurfaceCreateInfoKHR sci;
    PFN_vkCreateWaylandSurfaceKHR vkCreateWaylandSurfaceKHR;

    vkCreateWaylandSurfaceKHR = (PFN_vkCreateWaylandSurfaceKHR)
        vkGetInstanceProcAddr(instance, "vkCreateWaylandSurfaceKHR");
    if (!vkCreateWaylandSurfaceKHR)
    {
        _glfwInputError(GLFW_API_UNAVAILABLE,
                        "Wayland: Vulkan instance missing VK_KHR_wayland_surface extension");
        return VK_ERROR_EXTENSION_NOT_PRESENT;
    }

    memset(&sci, 0, sizeof(sci));
    sci.sType = VK_STRUCTURE_TYPE_WAYLAND_SURFACE_CREATE_INFO_KHR;
    sci.display = _glfw.wl.display;
    sci.surface = window->wl.surface;

    err = vkCreateWaylandSurfaceKHR(instance, &sci, allocator, surface);
    if (err)
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "Wayland: Failed to create Vulkan surface: %s",
                        _glfwGetVulkanResultString(err));
    }

    return err;
}


//////////////////////////////////////////////////////////////////////////
//////                        GLFW native API                       //////
//////////////////////////////////////////////////////////////////////////

GLFWAPI struct wl_display* glfwGetWaylandDisplay(void)
{
    _GLFW_REQUIRE_INIT_OR_RETURN(NULL);
    return _glfw.wl.display;
}

GLFWAPI struct wl_surface* glfwGetWaylandWindow(GLFWwindow* handle)
{
    _GLFWwindow* window = (_GLFWwindow*) handle;
    _GLFW_REQUIRE_INIT_OR_RETURN(NULL);
    return window->wl.surface;
}

