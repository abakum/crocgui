package org.golang.app;

import android.content.Intent;
import android.os.Build;
import android.os.Bundle;
import androidx.activity.OnBackPressedCallback;

public class GoNativeActivity extends org.golang.app.GoNativeActivityBase {

    static {
        System.loadLibrary("croc");
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            getOnBackInvokedDispatcher().registerOnBackInvokedCallback(
                OnBackInvokedDispatcher.PRIORITY_DEFAULT,
                this::handleBackPressed
            );
        } else {
            getOnBackPressedDispatcher().addCallback(this, new OnBackPressedCallback(true) {
                @Override
                public void handleOnBackPressed() {
                    handleBackPressed();
                }
            });
        }

        handleIntent(getIntent());
        // if (!isTaskRoot()) { 
        //     exitApplication(); 
        // }
    }

    @Override
    protected void onNewIntent(Intent intent) {
        handleIntent(Intent);
        super.onNewIntent(intent);
        setIntent(intent);
    }

    @Override
    protected void onResume() {
        super.onResume();
        handleIntent(getIntent());
    }

    private void handleIntent(Intent intent) {
        if (intent == null) {
            return;
        }
        processIntent(intent);
        // String action = intent.getAction();
        // if (Intent.ACTION_VIEW.equals(action) || 
        //     Intent.ACTION_SEND.equals(action) ||
        //     Intent.ACTION_SEND_MULTIPLE.equals(action)) {
        //     processIntent(intent);
        // }
    }

    private void handleBackPressed() {
        exitApplication();
    }

    public void exitApplication() {
        finishAffinity();
        finishAndRemoveTask();
        System.exit(0);
    }

    @Deprecated
    @Override
    public void onBackPressed() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            super.onBackPressed();
            return;
        }
        handleBackPressed();
    }
}