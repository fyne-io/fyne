// +build android

#include <android/log.h>
#include <jni.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>

#define LOG_FATAL(...) __android_log_print(ANDROID_LOG_FATAL, "Fyne", __VA_ARGS__)

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

char* getString(uintptr_t jni_env, uintptr_t ctx, jstring str) {
	JNIEnv *env = (JNIEnv*)jni_env;

	const char *chars = (*env)->GetStringUTFChars(env, str, NULL);

	const char *copy = strdup(chars);
	(*env)->ReleaseStringUTFChars(env, str, chars);
	return copy;
}

jobject parseURI(uintptr_t jni_env, uintptr_t ctx, char* uriCstr) {
	JNIEnv *env = (JNIEnv*)jni_env;

	jstring uriStr = (*env)->NewStringUTF(env, uriCstr);
	jclass uriClass = find_class(env, "android/net/Uri");
	jmethodID parse = find_static_method(env, uriClass, "parse", "(Ljava/lang/String;)Landroid/net/Uri;");

	return (jobject)(*env)->CallStaticObjectMethod(env, uriClass, parse, uriStr);
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

	return getString(jni_env, ctx, s);
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

	jobject uri = parseURI(jni_env, ctx, uriCstr);
	jobject stream = (jobject)(*env)->CallObjectMethod(env, resolver, openInputStream, uri);
	jthrowable loadErr = (*env)->ExceptionOccurred(env);

	if (loadErr != NULL) {
		(*env)->ExceptionClear(env);
		return NULL;
	}

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

bool uriCanList(uintptr_t jni_env, uintptr_t ctx, char* uriCstr) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jobject resolver = getContentResolver(jni_env, ctx);
	jstring uriStr = (*env)->NewStringUTF(env, uriCstr);
	jobject uri = parseURI(jni_env, ctx, uriCstr);

	jclass contractClass = find_class(env, "android/provider/DocumentsContract");
	if (contractClass == NULL) { // API 19
		return false;
	}
	jmethodID getDoc = find_static_method(env, contractClass, "getTreeDocumentId", "(Landroid/net/Uri;)Ljava/lang/String;");
	if (getDoc == NULL) { // API 21
		return false;
	}
	jstring docID = (jobject)(*env)->CallStaticObjectMethod(env, contractClass, getDoc, uri);

	jmethodID getTree = find_static_method(env, contractClass, "buildDocumentUriUsingTree", "(Landroid/net/Uri;Ljava/lang/String;)Landroid/net/Uri;");
	jobject treeUri = (jobject)(*env)->CallStaticObjectMethod(env, contractClass, getTree, uri, docID);

	jclass resolverClass = (*env)->GetObjectClass(env, resolver);
	jmethodID getType = find_method(env, resolverClass, "getType", "(Landroid/net/Uri;)Ljava/lang/String;");
	jstring type = (jstring)(*env)->CallObjectMethod(env, resolver, getType, treeUri);

	if (type == NULL) {
		return false;
	}

	char *str = getString(jni_env, ctx, type);
	return strcmp(str, "vnd.android.document/directory") == 0;
}

char* uriList(uintptr_t jni_env, uintptr_t ctx, char* uriCstr) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jobject resolver = getContentResolver(jni_env, ctx);
	jstring uriStr = (*env)->NewStringUTF(env, uriCstr);
	jobject uri = parseURI(jni_env, ctx, uriCstr);

	jclass contractClass = find_class(env, "android/provider/DocumentsContract");
	if (contractClass == NULL) { // API 19
		return "";
	}
	jmethodID getDoc = find_static_method(env, contractClass, "getTreeDocumentId", "(Landroid/net/Uri;)Ljava/lang/String;");
	if (getDoc == NULL) { // API 21
		return "";
	}
	jstring docID = (jobject)(*env)->CallStaticObjectMethod(env, contractClass, getDoc, uri);

	jmethodID getChild = find_static_method(env, contractClass, "buildChildDocumentsUriUsingTree", "(Landroid/net/Uri;Ljava/lang/String;)Landroid/net/Uri;");
	jobject childrenUri = (jobject)(*env)->CallStaticObjectMethod(env, contractClass, getChild, uri, docID);

	jclass stringClass = find_class(env, "java/lang/String");
	jobjectArray project = (*env)->NewObjectArray(env, 1, stringClass, (*env)->NewStringUTF(env, "document_id"));

	jclass resolverClass = (*env)->GetObjectClass(env, resolver);
	jmethodID query = find_method(env, resolverClass, "query", "(Landroid/net/Uri;[Ljava/lang/String;Landroid/os/Bundle;Landroid/os/CancellationSignal;)Landroid/database/Cursor;");
	if (getDoc == NULL) { // API 26
		return "";
	}

	jobject cursor = (jobject)(*env)->CallObjectMethod(env, resolver, query, childrenUri, project, NULL, NULL);
	jclass cursorClass = (*env)->GetObjectClass(env, cursor);
	jmethodID next = find_method(env, cursorClass, "moveToNext", "()Z");
	jmethodID get = find_method(env, cursorClass, "getString", "(I)Ljava/lang/String;");

	char *ret = NULL;
	int len = 0;
	while (((jboolean)(*env)->CallBooleanMethod(env, cursor, next)) == JNI_TRUE) {
		jstring name = (jstring)(*env)->CallObjectMethod(env, cursor, get, 0);
		jobject childUri = (jobject)(*env)->CallStaticObjectMethod(env, contractClass, getChild, uri, name);
		jclass uriClass = (*env)->GetObjectClass(env, childUri);
		jmethodID toString = (*env)->GetMethodID(env, uriClass, "toString", "()Ljava/lang/String;");
		jobject s = (*env)->CallObjectMethod(env, childUri, toString);

		char *uid = getString(jni_env, ctx, name);

		// append
		char *old = ret;
		len = len + strlen(uid) + 1;
		ret = malloc(sizeof(char)*(len+1));
		if (old != NULL) {
			strcpy(ret, old);
			free(old);
		}
		strcat(ret, uid);
		strcat(ret, "|");
	}

	ret[len-1] = '\0';
	return ret;
}