//========================================================================
// GLFW 3.2 POSIX - www.glfw.org
//------------------------------------------------------------------------
// Copyright (c) 2002-2006 Marcus Geelnard
// Copyright (c) 2006-2016 Camilla Berglund <elmindreda@glfw.org>
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


//////////////////////////////////////////////////////////////////////////
//////                       GLFW internal API                      //////
//////////////////////////////////////////////////////////////////////////

GLFWbool _glfwInitThreadLocalStoragePOSIX(void)
{
    if (pthread_key_create(&_glfw.posix_tls.context, NULL) != 0)
    {
        _glfwInputError(GLFW_PLATFORM_ERROR,
                        "POSIX: Failed to create context TLS");
        return GLFW_FALSE;
    }

    _glfw.posix_tls.allocated = GLFW_TRUE;
    return GLFW_TRUE;
}

void _glfwTerminateThreadLocalStoragePOSIX(void)
{
    if (_glfw.posix_tls.allocated)
        pthread_key_delete(_glfw.posix_tls.context);
}


//////////////////////////////////////////////////////////////////////////
//////                       GLFW platform API                      //////
//////////////////////////////////////////////////////////////////////////

void _glfwPlatformSetCurrentContext(_GLFWwindow* context)
{
    pthread_setspecific(_glfw.posix_tls.context, context);
}

_GLFWwindow* _glfwPlatformGetCurrentContext(void)
{
    return pthread_getspecific(_glfw.posix_tls.context);
}

