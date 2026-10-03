//go:build !ci && android

#include <jni.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>

static bool callSecretMethod(uintptr_t jni_env, uintptr_t ctx, const char *name,
		const void *in, int inLen, void **out, int *outLen) {
	JNIEnv *env = (JNIEnv*)jni_env;

	jclass cls = (*env)->GetObjectClass(env, (jobject)ctx);
	jmethodID mid = (*env)->GetStaticMethodID(env, cls, name, "([B)[B");
	if (mid == 0) {
		(*env)->ExceptionClear(env);
		return false;
	}

	jbyteArray input = (*env)->NewByteArray(env, inLen);
	if (input == NULL) {
		(*env)->ExceptionClear(env);
		return false;
	}
	if (inLen > 0) {
		(*env)->SetByteArrayRegion(env, input, 0, inLen, (const jbyte *)in);
	}

	jbyteArray result = (jbyteArray)(*env)->CallStaticObjectMethod(env, cls, mid, input);
	(*env)->DeleteLocalRef(env, input);
	if ((*env)->ExceptionCheck(env)) {
		(*env)->ExceptionClear(env);
		return false;
	}
	if (result == NULL) {
		return false;
	}

	jsize len = (*env)->GetArrayLength(env, result);
	void *buf = malloc(len > 0 ? len : 1);
	if (len > 0) {
		(*env)->GetByteArrayRegion(env, result, 0, len, (jbyte *)buf);
	}
	(*env)->DeleteLocalRef(env, result);

	*out = buf;
	*outLen = (int)len;
	return true;
}

bool secretEncrypt(uintptr_t jni_env, uintptr_t ctx, const void *in, int inLen, void **out, int *outLen) {
	return callSecretMethod(jni_env, ctx, "secretEncrypt", in, inLen, out, outLen);
}

bool secretDecrypt(uintptr_t jni_env, uintptr_t ctx, const void *in, int inLen, void **out, int *outLen) {
	return callSecretMethod(jni_env, ctx, "secretDecrypt", in, inLen, out, outLen);
}
