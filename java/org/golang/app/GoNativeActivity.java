package org.golang.app;

import android.content.Intent;
import android.os.Bundle;

public class GoNativeActivity extends org.golang.app.GoNativeActivityBase {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        processIntent(getIntent());
    }

    @Override
    protected void onNewIntent(Intent intent) {
        super.onNewIntent(intent);
        processIntent(intent);
    }

    public native void processIntent(Intent intent);
    
    static {
        System.loadLibrary("croc");
    }
}