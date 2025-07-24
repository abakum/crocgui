package org.golang.app;

import android.content.Intent;
import android.os.Bundle;

public class GoNativeActivity extends org.golang.app.GoNativeActivityBase {
    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        handleIntent(getIntent());  // Обработка при холодном запуске
    }

    @Override
    protected void onNewIntent(Intent intent) {
        super.onNewIntent(intent);
        setIntent(intent);  // Критически важно!
        handleIntent(intent);  // Обработка при горячем запуске
    }

    @Override
    protected void onResume() {
        super.onResume();
        // Двойная проверка для надежности
        if (getIntent() != null && getIntent().getAction() != null) {
            handleIntent(getIntent());
        }
    }

    private void handleIntent(Intent intent) {
        String action = intent.getAction();
        if (Intent.ACTION_VIEW.equals(action) || 
            Intent.ACTION_SEND.equals(action) ||
            Intent.ACTION_SEND_MULTIPLE.equals(action)) {
            processIntent(intent);  // Вызов нативного метода
        }
    }

    static {
        System.loadLibrary("croc");
    }
}

// package org.golang.app;

// import android.content.Intent;
// import android.os.Bundle;
// import android.net.Uri;
// import java.util.ArrayList;

// public class GoNativeActivity extends org.golang.app.GoNativeActivityBase {
//     @Override
//     protected void onCreate(Bundle savedInstanceState) {
//         super.onCreate(savedInstanceState);
//         handleIntent(getIntent());
//     }

//     @Override
//     protected void onNewIntent(Intent intent) {
//         super.onNewIntent(intent);
//         handleIntent(intent);
//     }

//     private void handleIntent(Intent intent) {
//         String action = intent.getAction();
//         String type = intent.getType();

//         if (Intent.ACTION_VIEW.equals(action)) {
//             // Обработка VIEW (Открыть с помощью)
//             Uri uri = intent.getData();
//             if (uri != null) {
//                 processViewUri(uri.toString());
//             }
//         } 
//         else if (Intent.ACTION_SEND.equals(action) && type != null) {
//             // Обработка SEND (Поделиться одним файлом)
//             Uri uri = intent.getParcelableExtra(Intent.EXTRA_STREAM);
//             if (uri != null) {
//                 processSendUri(uri.toString());
//             } else {
//                 String text = intent.getStringExtra(Intent.EXTRA_TEXT);
//                 if (text != null) {
//                     processText(text);
//                 }
//             }
//         } 
//         else if (Intent.ACTION_SEND_MULTIPLE.equals(action) && type != null) {
//             // Обработка SEND_MULTIPLE (Несколько файлов)
//             ArrayList<Uri> uris = intent.getParcelableArrayListExtra(Intent.EXTRA_STREAM);
//             if (uris != null) {
//                 String[] uriStrings = new String[uris.size()];
//                 for (int i = 0; i < uris.size(); i++) {
//                     uriStrings[i] = uris.get(i).toString();
//                 }
//                 processMultipleUris(uriStrings);
//             }
//         }
//     }

//     // Нативные методы для разных типов контента
//     public native void processViewUri(String uri);
//     public native void processSendUri(String uri);
//     public native void processText(String text);
//     public native void processMultipleUris(String[] uris);
    
//     static {
//         System.loadLibrary("croc");
//     }
// }
