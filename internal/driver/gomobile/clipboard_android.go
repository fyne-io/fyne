// +build android

package gomobile

/*
#cgo LDFLAGS: -landroid -llog -lEGL -lGLESv2

#include <android/log.h>
#include <jni.h>
#include <stdlib.h>
#include <string.h>

#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "GoLog", __VA_ARGS__)

static jmethodID find_method(JNIEnv *env, jclass clazz, const char *name, const char *sig) {
	jmethodID m = (*env)->GetMethodID(env, clazz, name, sig);
	if (m == 0) {
		(*env)->ExceptionClear(env);
		LOG_FATAL("cannot find method %s %s", name, sig);
		return 0;
	}
	return m;
}

jobject getClipboard(uintptr_t jni_env, uintptr_t ctx) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jclass ctxClass = (*env)->GetObjectClass(env, ctx);
	jmethodID getSystemService = find_method(env, ctxClass, "getSystemService", "(Ljava/lang/String;)Ljava/lang/Object;");

	jstring service = (*env)->NewStringUTF(env, "clipboard");
	return (jobject)(*env)->CallObjectMethod(env, ctx, getSystemService, service);
}

char *getClipboardContent(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jobject mgr = getClipboard(jni_env, ctx);

	jclass mgrClass = (*env)->GetObjectClass(env, mgr);
	jmethodID getText = find_method(env, mgrClass, "getText", "()Ljava/lang/CharSequence;");

	jobject content = (jstring)(*env)->CallObjectMethod(env, mgr, getText);
	if (content == NULL) {
		return "";
	}

	jclass clzCharSequence = (*env)->GetObjectClass(env, content);
	jmethodID toString = (*env)->GetMethodID(env, clzCharSequence, "toString", "()Ljava/lang/String;");
	jobject s = (*env)->CallObjectMethod(env, content, toString);

	const char *chars = (*env)->GetStringUTFChars(env, s, NULL);
	const char *copy = strdup(chars);
	(*env)->ReleaseStringUTFChars(env, s, chars);
	return copy;
}

void setClipboardContent(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *content) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jobject mgr = getClipboard(jni_env, ctx);

	jclass mgrClass = (*env)->GetObjectClass(env, mgr);
	jmethodID setText = find_method(env, mgrClass, "setText", "(Ljava/lang/CharSequence;)V");

	jstring str = (*env)->NewStringUTF(env, content);
	(*env)->CallVoidMethod(env, mgr, setText, str);
}
*/
import "C"
import (
	"unsafe"

	"github.com/fyne-io/mobile/app"
)

// Content returns the clipboard content for Android
func (c *mobileClipboard) Content() string {
	content := ""
	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		chars := C.getClipboardContent(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx))
		content = C.GoString(chars)
		C.free(unsafe.Pointer(chars))
		return nil
	})
	return content
}

// SetContent sets the clipboard content for Android
func (c *mobileClipboard) SetContent(content string) {
	contentStr := C.CString(content)
	defer C.free(unsafe.Pointer(contentStr))

	app.RunOnJVM(func(vm, env, ctx uintptr) error {
		C.setClipboardContent(C.uintptr_t(vm), C.uintptr_t(env), C.uintptr_t(ctx), contentStr)
		return nil
	})
}
