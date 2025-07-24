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

    if (intent != NULL) {
        jclass intent_class = (*env)->GetObjectClass(env, intent);
        jmethodID get_data = (*env)->GetMethodID(env, intent_class, "getData", "()Landroid/net/Uri;");
        jobject uri = (*env)->CallObjectMethod(env, intent, get_data);

        if (uri != NULL) {
            jclass uri_class = (*env)->GetObjectClass(env, uri);
            jmethodID to_string = (*env)->GetMethodID(env, uri_class, "toString", "()Ljava/lang/String;");
            jstring uri_string = (*env)->CallObjectMethod(env, uri, to_string);

            const char *utf_str = (*env)->GetStringUTFChars(env, uri_string, NULL);
            receiveURIFromIntent(strdup(utf_str));
            (*env)->ReleaseStringUTFChars(env, uri_string, utf_str);
        }
    }
}
*/
import "C"
import (
	"unsafe"
	log "github.com/schollz/logger"

	"fyne.io/fyne/v2/driver"
)

//export receiveURIFromIntent
func receiveURIFromIntent(uri *C.char) {
	goURI := C.GoString(uri)
	log.Tracef("Received URI: %s", goURI)
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
