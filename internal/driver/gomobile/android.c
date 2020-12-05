// +build android

#include <android/log.h>
#include <dirent.h>
#include <jni.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>

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
	jobject ret = (jobject)(*env)->CallObjectMethod(env, ctx, getSystemService, service);
	jthrowable err = (*env)->ExceptionOccurred(env);

	if (err != NULL) {
		LOG_FATAL("cannot lookup clipboard");
		(*env)->ExceptionClear(env);
		return NULL;
	}
	return ret;
}

char *getClipboardContent(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jobject mgr = getClipboard(jni_env, ctx);
	if (mgr == NULL) {
		return NULL;
	}

	jclass mgrClass = (*env)->GetObjectClass(env, mgr);
	jmethodID getText = find_method(env, mgrClass, "getText", "()Ljava/lang/CharSequence;");

	jobject content = (jstring)(*env)->CallObjectMethod(env, mgr, getText);
	if (content == NULL) {
		return NULL;
	}

	jclass clzCharSequence = (*env)->GetObjectClass(env, content);
	jmethodID toString = (*env)->GetMethodID(env, clzCharSequence, "toString", "()Ljava/lang/String;");
	jobject s = (*env)->CallObjectMethod(env, content, toString);

	return getString(jni_env, ctx, s);
}

void setClipboardContent(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *content) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jobject mgr = getClipboard(jni_env, ctx);
	if (mgr == NULL) {
		return;
	}

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

bool hasPrefix(char* string, char* prefix) {
	size_t lp = strlen(prefix);
	size_t ls = strlen(string);
	if (ls < lp) {
		return false;
	}
	return memcmp(prefix, string, lp) == 0;
}

bool canListContentURI(uintptr_t jni_env, uintptr_t ctx, char* uriCstr) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jobject resolver = getContentResolver(jni_env, ctx);
	jstring uriStr = (*env)->NewStringUTF(env, uriCstr);
	jobject uri = parseURI(jni_env, ctx, uriCstr);
	jthrowable loadErr = (*env)->ExceptionOccurred(env);

	if (loadErr != NULL) {
		(*env)->ExceptionClear(env);
		return false;
	}

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

bool canListFileURI(char* uriCstr) {
	// Get file path from URI
	size_t length = strlen(uriCstr)-7;// -7 for 'file://'
	char* path = malloc(sizeof(char)*(length+1));// +1 for '\0'
	memcpy(path, &uriCstr[7], length);
	path[length] = '\0';

	// Stat path to determine if it points to a directory
	struct stat statbuf;
	int result = stat(path, &statbuf);

	free(path);

	return (result == 0) && S_ISDIR(statbuf.st_mode);
}

bool canListURI(uintptr_t jni_env, uintptr_t ctx, char* uriCstr) {
	if (hasPrefix(uriCstr, "file://")) {
		return canListFileURI(uriCstr);
	} else if (hasPrefix(uriCstr, "content://")) {
		return canListContentURI(jni_env, ctx, uriCstr);
	}
	LOG_FATAL("Unrecognized scheme: %s", uriCstr);
	return false;
}

char* listContentURI(uintptr_t jni_env, uintptr_t ctx, char* uriCstr) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jobject resolver = getContentResolver(jni_env, ctx);
	jstring uriStr = (*env)->NewStringUTF(env, uriCstr);
	jobject uri = parseURI(jni_env, ctx, uriCstr);
	jthrowable loadErr = (*env)->ExceptionOccurred(env);

	if (loadErr != NULL) {
		(*env)->ExceptionClear(env);
		return "";
	}

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
		} else {
			ret[0] = '\0';
		}
		strcat(ret, uid);
		strcat(ret, "|");
	}

	if (ret != NULL) {
		ret[len-1] = '\0';
	}
	return ret;
}

char* listFileURI(char* uriCstr) {

	size_t uriLength = strlen(uriCstr);

	// Get file path from URI
	size_t length = uriLength-7;// -7 for 'file://'
	char* path = malloc(sizeof(char)*(length+1));// +1 for '\0'
	memcpy(path, &uriCstr[7], length);
	path[length] = '\0';

	char *ret = NULL;
	DIR *dfd;
	if ((dfd = opendir(path)) != NULL) {
		struct dirent *dp;
		int len = 0;
		while ((dp = readdir(dfd)) != NULL) {
			if (strcmp(dp->d_name, ".") == 0) {
				// Ignore current directory
				continue;
			}
			if (strcmp(dp->d_name, "..") == 0) {
				// Ignore parent directory
				continue;
			}
			// append
			char *old = ret;
			len = len + uriLength + 1 /* / */ + strlen(dp->d_name) + 1 /* | */;
			ret = malloc(sizeof(char)*(len+1));
			if (old != NULL) {
				strcpy(ret, old);
				free(old);
			} else {
				ret[0] = '\0';
			}
			strcat(ret, uriCstr);
			strcat(ret, "/");
			strcat(ret, dp->d_name);
			strcat(ret, "|");
		}
		if (ret != NULL) {
			ret[len-1] = '\0';
		}
	}

	free(path);

	return ret;
}

char* listURI(uintptr_t jni_env, uintptr_t ctx, char* uriCstr) {
	if (hasPrefix(uriCstr, "file://")) {
		return listFileURI(uriCstr);
	} else if (hasPrefix(uriCstr, "content://")) {
		return listContentURI(jni_env, ctx, uriCstr);
	}
	LOG_FATAL("Unrecognized scheme: %s", uriCstr);
	return "";
}
