//go:build android

package main

/*
#include <jni.h>

void exitAppJNI(uintptr_t java_vm, uintptr_t jni_env, uintptr_t ctx) {
    JavaVM* vm = (JavaVM*)java_vm;
    JNIEnv* env = (JNIEnv*)jni_env;
    jobject activity = (jobject)ctx;

    jclass cls = (*env)->GetObjectClass(env, activity);
    jmethodID method = (*env)->GetMethodID(env, cls, "exitApplication", "()V");
    (*env)->CallVoidMethod(env, activity, method);
}
*/
import "C"
import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver"
)

func quit(a fyne.App) {
	driver.RunNative(func(ctx interface{}) error {
		ac := ctx.(*driver.AndroidContext)
		C.exitAppJNI(
			C.uintptr_t(ac.VM),
			C.uintptr_t(ac.Env),
			C.uintptr_t(ac.Ctx),
		)
		return nil
	})
}
