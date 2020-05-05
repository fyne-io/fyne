// +build !ci

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

jobject getSystemService(uintptr_t jni_env, uintptr_t ctx, char *service) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jstring serviceStr = (*env)->NewStringUTF(env, service);

	jclass ctxClass = (*env)->GetObjectClass(env, ctx);
	jmethodID getSystemService = find_method(env, ctxClass, "getSystemService", "(Ljava/lang/String;)Ljava/lang/Object;");

	return (jobject)(*env)->CallObjectMethod(env, ctx, getSystemService, serviceStr);
}

int nextId = 1;

void sendNotification(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx, char *title, char *body) {
	JNIEnv *env = (JNIEnv*)jni_env;
	jstring titleStr = (*env)->NewStringUTF(env, title);
	jstring bodyStr = (*env)->NewStringUTF(env, body);

	jclass cls = find_class(env, "android/app/Notification$Builder");
	jmethodID constructor = find_method(env, cls, "<init>", "(Landroid/content/Context;)V");
	jobject builder = (*env)->NewObject(env, cls, constructor, ctx);

	jmethodID setContentTitle = find_method(env, cls, "setContentTitle", "(Ljava/lang/CharSequence;)Landroid/app/Notification$Builder;");
	(*env)->CallObjectMethod(env, builder, setContentTitle, titleStr);

	jmethodID setContentText = find_method(env, cls, "setContentText", "(Ljava/lang/CharSequence;)Landroid/app/Notification$Builder;");
	(*env)->CallObjectMethod(env, builder, setContentText, bodyStr);

	int iconID = 17629184; // constant of "unknown app icon"
	jmethodID setSmallIcon = find_method(env, cls, "setSmallIcon", "(I)Landroid/app/Notification$Builder;");
	(*env)->CallObjectMethod(env, builder, setSmallIcon, iconID);

	jmethodID build = find_method(env, cls, "build", "()Landroid/app/Notification;");
	jobject notif = (*env)->CallObjectMethod(env, builder, build);

	jclass mgrCls = find_class(env, "android/app/NotificationManager");
	jobject mgr = getSystemService(env, ctx, "notification");

	jmethodID notify = find_method(env, mgrCls, "notify", "(ILandroid/app/Notification;)V");
	(*env)->CallVoidMethod(env, mgr, notify, nextId, notif);
	nextId++;
}