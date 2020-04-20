// +build android

package gomobile

/*
#include <android/log.h>
#include <jni.h>
#include <stdlib.h>

#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "GoLog", __VA_ARGS__)

static jclass find_class(JNIEnv *env, const char *class_name) {
	jclass clazz = (*env)->FindClass(env, class_name);
	if (clazz == NULL) {
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find %s", class_name);
		return NULL;
	}
	return clazz;
}

static jmethodID find_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
	jmethodID m = (*env)->GetMethodID(env, clazz, name, sig);
	if (m == 0) {
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find method %s %s", name, sig);
		return 0;
	}
	return m;
}

static jmethodID find_static_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
	jmethodID m = (*env)->GetStaticMethodID(env, clazz, name, sig);
	if (m == 0) {
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find method %s %s", name, sig);
		return 0;
	}
	return m;
}

jobject getContentResolver(uintptr_t jni_env, uintptr_t ctx) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jclass ctxClass = (*env)->GetObjectClass(env, ctx);
	jmethodID getContentResolver = find_method(env, ctxClass, "getContentResolver", "()Landroid/content/ContentResolver;");

	return (jobject)(*env)->CallObjectMethod(env, ctx, getContentResolver);
}

void* openStream(uintptr_t jni_env, uintptr_t ctx, char* uriCstr) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jobject resolver = getContentResolver(jni_env, ctx);

	jclass resolverClass = (*env)->GetObjectClass(env, resolver);
	jmethodID openInputStream = find_method(env, resolverClass, "openInputStream", "(Landroid/net/Uri;)Ljava/io/InputStream;");

	jstring uriStr = (*env)->NewStringUTF(env, uriCstr);
	jclass uriClass = find_class(env, "android/net/Uri");
	jmethodID parse = find_static_method(env, uriClass, "parse", "(Ljava/lang/String;)Landroid/net/Uri;");
	jobject uri = (jobject)(*env)->CallStaticObjectMethod(env, uriClass, parse, uriStr);

	jobject stream = (jobject)(*env)->CallObjectMethod(env, resolver, openInputStream, uri);
	return (*env)->NewGlobalRef(env, stream);
}

char* readStream(uintptr_t jni_env, uintptr_t ctx, void* stream, int len, int* total) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jclass streamClass = (*env)->GetObjectClass(env, stream);
	jmethodID read = find_method(env, streamClass, "read", "([BII)I");

	jbyteArray data = (*env)->NewByteArray(env, len);
	int count = (int)(*env)->CallIntMethod(env, stream, read, data, 0, len);
	*total = count;

	if (count == -1) {
		return NULL;
	}

	char* bytes = malloc(sizeof(char)*count);
	(*env)->GetByteArrayRegion(env, data, 0, count, bytes);
	return bytes;
}

void closeStream(uintptr_t jni_env, uintptr_t ctx, void* stream) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jclass streamClass = (*env)->GetObjectClass(env, stream);
	jmethodID close = find_method(env, streamClass, "close", "()V");
	(*env)->CallVoidMethod(env, stream, close);

	(*env)->DeleteGlobalRef(env, stream);
}
*/
import "C"
import (
	"errors"
	"io"
	"unsafe"

	"github.com/fyne-io/mobile/app"
)

type javaStream struct {
	stream unsafe.Pointer // java.io.InputStream
}

// Declare conformity to ReadCloser interface
var _ io.ReadCloser = (*javaStream)(nil)

func (s *javaStream) Read(p []byte) (int, error) {
	count := 0
	err := app.RunOnJVM(func(_, env, ctx uintptr) error {
		cCount := C.int(0)
		cBytes := unsafe.Pointer(C.readStream(C.uintptr_t(env), C.uintptr_t(ctx), s.stream, C.int(len(p)), &cCount))
		if cCount == -1 {
			return io.EOF
		}
		defer C.free(cBytes)
		count = int(cCount) // avoid sending -1 instead of 0 on completion

		bytes := C.GoBytes(cBytes, cCount)
		for i := 0; i < int(count); i++ {
			p[i] = bytes[i]
		}
		return nil
	})

	return int(count), err
}

func (s *javaStream) Close() error {
	app.RunOnJVM(func(_, env, ctx uintptr) error {
		C.closeStream(C.uintptr_t(env), C.uintptr_t(ctx), s.stream)

		return nil
	})

	return nil
}

func openStream(uri string) unsafe.Pointer {
	uriStr := C.CString(uri)
	defer C.free(unsafe.Pointer(uriStr))

	var stream unsafe.Pointer
	app.RunOnJVM(func(_, env, ctx uintptr) error {
		stream = unsafe.Pointer(C.openStream(C.uintptr_t(env), C.uintptr_t(ctx), uriStr))

		return nil
	})
	return stream
}

func nativeFileOpen(f *fileOpen) (io.ReadCloser, error) {
	ret := openStream(f.uri)
	if ret == nil {
		return nil, errors.New("resource not found at URI")
	}

	stream := &javaStream{}
	stream.stream = ret
	return stream, nil
}

func nativeFileSave(f *fileSave) (io.WriteCloser, error) {
	panic("Please implement me")
}
