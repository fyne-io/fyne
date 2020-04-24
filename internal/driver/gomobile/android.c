// +build android

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

// clipboard

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

// file handling

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
