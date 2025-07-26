//go:build android

package main

/*
#include <jni.h>
#include <stdlib.h>
#include <string.h>

extern void receiveURIFromIntent(char* uri);

static void processIntent(JNIEnv* env, jobject activity) {
    jclass activity_class = (*env)->GetObjectClass(env, activity);
    jmethodID get_intent = (*env)->GetMethodID(env, activity_class, "getIntent", "()Landroid/content/Intent;");
    jobject intent = (*env)->CallObjectMethod(env, activity, get_intent);

    if (intent == NULL) return;

    jclass intent_class = (*env)->GetObjectClass(env, intent);

    if (intent_class == NULL) return;

    // Handle ACTION_VIEW (single URI)
    jmethodID get_data = (*env)->GetMethodID(env, intent_class, "getData", "()Landroid/net/Uri;");
    jobject uri = (*env)->CallObjectMethod(env, intent, get_data);

    if (uri != NULL) {
        jclass uri_class = (*env)->GetObjectClass(env, uri);
        jmethodID to_string = (*env)->GetMethodID(env, uri_class, "toString", "()Ljava/lang/String;");
        jstring uri_string = (*env)->CallObjectMethod(env, uri, to_string);

        const char *utf_str = (*env)->GetStringUTFChars(env, uri_string, NULL);
        receiveURIFromIntent(strdup(utf_str));
        (*env)->ReleaseStringUTFChars(env, uri_string, utf_str);
        return;
    }

    // Handle ACTION_SEND (single content)
    jmethodID get_extra = (*env)->GetMethodID(env, intent_class, "getParcelableExtra", "(Ljava/lang/String;)Landroid/os/Parcelable;");
    jstring extra_key = (*env)->NewStringUTF(env, "android.intent.extra.STREAM");
    jobject send_uri = (*env)->CallObjectMethod(env, intent, get_extra, extra_key);
    (*env)->DeleteLocalRef(env, extra_key);

    if (send_uri != NULL) {
        jclass uri_class = (*env)->GetObjectClass(env, send_uri);
        jmethodID to_string = (*env)->GetMethodID(env, uri_class, "toString", "()Ljava/lang/String;");
        jstring uri_string = (*env)->CallObjectMethod(env, send_uri, to_string);

        const char *utf_str = (*env)->GetStringUTFChars(env, uri_string, NULL);
        receiveURIFromIntent(strdup(utf_str));
        (*env)->ReleaseStringUTFChars(env, uri_string, utf_str);
        return;
    }

    // Handle ACTION_SEND_MULTIPLE (multiple content)
    jmethodID get_array = (*env)->GetMethodID(env, intent_class, "getParcelableArrayListExtra", "(Ljava/lang/String;)Ljava/util/ArrayList;");
    jstring array_key = (*env)->NewStringUTF(env, "android.intent.extra.STREAM");
    jobject uri_list = (*env)->CallObjectMethod(env, intent, get_array, array_key);
    (*env)->DeleteLocalRef(env, array_key);

    if (uri_list != NULL) {
        jclass array_list_class = (*env)->GetObjectClass(env, uri_list);
        jmethodID get_size = (*env)->GetMethodID(env, array_list_class, "size", "()I");
        jmethodID get_item = (*env)->GetMethodID(env, array_list_class, "get", "(I)Ljava/lang/Object;");
        jint size = (*env)->CallIntMethod(env, uri_list, get_size);

        // Process all URIs in the list
        for (int i = 0; i < size; i++) {
            jobject current_uri = (*env)->CallObjectMethod(env, uri_list, get_item, i);
            jclass uri_class = (*env)->GetObjectClass(env, current_uri);
            jmethodID to_string = (*env)->GetMethodID(env, uri_class, "toString", "()Ljava/lang/String;");
            jstring uri_string = (*env)->CallObjectMethod(env, current_uri, to_string);

            const char *utf_str = (*env)->GetStringUTFChars(env, uri_string, NULL);
            receiveURIFromIntent(strdup(utf_str));
            (*env)->ReleaseStringUTFChars(env, uri_string, utf_str);

            (*env)->DeleteLocalRef(env, current_uri);
            (*env)->DeleteLocalRef(env, uri_string);
        }
    }
}
*/
import "C"
import (
	"unsafe"

	"fyne.io/fyne/v2/driver"
)

//export receiveURIFromIntent
func receiveURIFromIntent(uri *C.char) {
	uriFromIntent <- C.GoString(uri)
	C.free(unsafe.Pointer(uri))
}

func setupIntentHandler() {
	driver.RunNative(func(ctx interface{}) error {
		ac := ctx.(*driver.AndroidContext)
		C.processIntent(
			(*C.JNIEnv)(unsafe.Pointer(ac.Env)),
			(C.jobject)(unsafe.Pointer(ac.Ctx)),
		)
		return nil
	})
}
