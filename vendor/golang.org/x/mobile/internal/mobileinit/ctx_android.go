// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mobileinit

/*
#include <jni.h>
#include <stdlib.h>

static char* lockJNI(JavaVM *vm, uintptr_t* envp, int* attachedp) {
	JNIEnv* env;

	if (vm == NULL) {
		return "no current JVM";
	}

	*attachedp = 0;
	switch ((*vm)->GetEnv(vm, (void**)&env, JNI_VERSION_1_6)) {
	case JNI_OK:
		break;
	case JNI_EDETACHED:
		if ((*vm)->AttachCurrentThread(vm, &env, 0) != 0) {
			return "cannot attach to JVM";
		}
		*attachedp = 1;
		break;
	case JNI_EVERSION:
		return "bad JNI version";
	default:
		return "unknown JNI error from GetEnv";
	}

	*envp = (uintptr_t)env;
	return NULL;
}

static char* checkException(uintptr_t jnienv) {
	jthrowable exc;
	JNIEnv* env = (JNIEnv*)jnienv;

	if (!(*env)->ExceptionCheck(env)) {
		return NULL;
	}

	exc = (*env)->ExceptionOccurred(env);
	(*env)->ExceptionClear(env);

	jclass clazz = (*env)->FindClass(env, "java/lang/Throwable");
	jmethodID toString = (*env)->GetMethodID(env, clazz, "toString", "()Ljava/lang/String;");
	jobject msgStr = (*env)->CallObjectMethod(env, exc, toString);
	return (char*)(*env)->GetStringUTFChars(env, msgStr, 0);
}

static void unlockJNI(JavaVM *vm) {
	(*vm)->DetachCurrentThread(vm);
}
*/
import "C"

import (
	"errors"
	"runtime"
	"unsafe"
)

// currentVM is stored to initialize other cgo packages.
//
// As all the Go packages in a program form a single shared library,
// there can only be one JNI_OnLoad function for initialization. In
// OpenJDK there is JNI_GetCreatedJavaVMs, but this is not available
// on android.
var currentVM *C.JavaVM

// currentCtx is Android's android.context.Context. May be NULL.
var currentCtx C.jobject

// SetCurrentContext populates the global Context object with the specified
// current JavaVM instance (vm) and android.context.Context object (ctx).
// The android.context.Context object must be a global reference.
func SetCurrentContext(vm unsafe.Pointer, ctx uintptr) {
	currentVM = (*C.JavaVM)(vm)
	currentCtx = (C.jobject)(ctx)
}

// RunOnJVM runs fn on a new goroutine locked to an OS thread with a JNIEnv.
//
// RunOnJVM blocks until the call to fn is complete. Any Java
// exception or failure to attach to the JVM is returned as an error.
//
// The function fn takes vm, the current JavaVM*,
// env, the current JNIEnv*, and
// ctx, a jobject representing the global android.context.Context.
func RunOnJVM(fn func(vm, env, ctx uintptr) error) error {
	errch := make(chan error)
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		env := C.uintptr_t(0)
		attached := C.int(0)
		if errStr := C.lockJNI(currentVM, &env, &attached); errStr != nil {
			errch <- errors.New(C.GoString(errStr))
			return
		}
		if attached != 0 {
			defer C.unlockJNI(currentVM)
		}

		vm := uintptr(unsafe.Pointer(currentVM))
		if err := fn(vm, uintptr(env), uintptr(currentCtx)); err != nil {
			errch <- err
			return
		}

		if exc := C.checkException(env); exc != nil {
			errch <- errors.New(C.GoString(exc))
			C.free(unsafe.Pointer(exc))
			return
		}
		errch <- nil
	}()
	return <-errch
}
