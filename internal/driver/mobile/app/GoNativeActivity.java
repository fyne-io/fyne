package org.golang.app;

import java.util.Base64;
import java.io.ByteArrayOutputStream;
import java.io.File;
import java.nio.ByteBuffer;
import android.app.Activity;
import android.app.AlarmManager;
import android.app.NativeActivity;
import android.app.NotificationManager;
import android.app.PendingIntent;
import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;
import android.content.pm.ActivityInfo;
import android.content.pm.PackageManager;
import android.content.res.Configuration;
import android.graphics.Rect;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.Bitmap.CompressFormat;
import android.graphics.Matrix;
import android.Manifest;
import android.net.Uri;
import android.os.Build;
import android.os.Bundle;
import android.provider.MediaStore;
import android.text.Editable;
import android.text.InputType;
import android.text.TextWatcher;
import android.text.method.DigitsKeyListener;
import android.util.Log;
import android.view.Gravity;
import android.view.KeyCharacterMap;
import android.view.View;
import android.view.ViewGroup;
import android.view.WindowInsets;
import android.view.accessibility.AccessibilityNodeInfo;
import android.view.WindowInsetsController;
import android.view.inputmethod.EditorInfo;
import android.view.inputmethod.InputMethodManager;
import android.view.KeyEvent;
import android.widget.EditText;
import android.widget.FrameLayout;
import android.widget.TextView;
import android.widget.TextView.OnEditorActionListener;
import java.util.ArrayList;
import java.util.List;
import androidx.camera.core.ImageCapture;
import androidx.camera.core.ImageCaptureException;
import androidx.camera.core.impl.utils.executor.CameraXExecutors;
import java.util.concurrent.Executors;
import androidx.annotation.NonNull;
import androidx.camera.lifecycle.ProcessCameraProvider;
import androidx.lifecycle.LifecycleOwner;
import androidx.lifecycle.Lifecycle;
import androidx.lifecycle.LifecycleService;
import androidx.lifecycle.LifecycleRegistry;
import androidx.lifecycle.ProcessLifecycleOwner;
import com.google.common.util.concurrent.ListenableFuture;
import androidx.camera.core.Camera;
import androidx.camera.core.CameraSelector;
import androidx.camera.core.ImageProxy;
import androidx.camera.camera2.Camera2Config;
import androidx.camera.view.PreviewView;
import androidx.camera.core.Preview;
import androidx.core.content.ContextCompat;
import androidx.core.app.ActivityCompat;

public class GoNativeActivity extends NativeActivity implements LifecycleOwner {
    private static GoNativeActivity goNativeActivity;
    private static final int FILE_OPEN_CODE = 1;
    private static final int FILE_SAVE_CODE = 2;
    private static final int CAMERA_OPEN_CODE = 3;

    private static final int DEFAULT_INPUT_TYPE = InputType.TYPE_TEXT_FLAG_NO_SUGGESTIONS;

    private static final int DEFAULT_KEYBOARD_CODE = 0;
    private static final int SINGLELINE_KEYBOARD_CODE = 1;
    private static final int NUMBER_KEYBOARD_CODE = 2;
    private static final int PASSWORD_KEYBOARD_CODE = 3;

    private native void filePickerReturned(String str);
    private native void capturePhotoReturned(byte[] jpegBytes, int length);
    private native void cameraPreviewFrame(byte[] jpegBytes, int length);
    private native void insetsChanged(int top, int bottom, int left, int right);
    private native void keyboardTyped(String str);
    private native void keyboardDelete();
    private native void backPressed();
    private native void setDarkMode(boolean dark);

    private EditText mTextEdit;
    private boolean ignoreKey = false;
    private boolean keyboardUp = false;

    private final LifecycleRegistry lifecycleRegistry = new LifecycleRegistry(this);
    private int pendingPermissionAction = 0;
    private Preview preview = null;
    private boolean isCameraRunning = false;
    private ProcessCameraProvider globalCameraProvider = null;
    private ImageCapture globalImageCapture = null;
    private BitmapSurfaceProvider globalBitmapSurfaceProvider = null;
    private boolean isPreviewRequested = false;

    // Hoisted out of doShowKeyboard / setupEntry to avoid nested anonymous
    // classes (Runnable -> Listener). javac stores a `MethodParameters`
    // attribute with an empty name on the synthetic `this$1` parameter of a
    // nested anonymous class's constructor; older D8 / R8 versions
    // (e.g. AOSP build-tools 3.3) NPE while reading that attribute.
    private final OnEditorActionListener mEditorActionListener = new OnEditorActionListener() {
        @Override
        public boolean onEditorAction(TextView v, int actionId, KeyEvent event) {
            if (actionId == EditorInfo.IME_ACTION_DONE) {
                keyboardTyped("\n");
            }
            return false;
        }
    };

    private final TextWatcher mTextWatcher = new TextWatcher() {
        @Override
        public void onTextChanged(CharSequence s, int start, int before, int count) {
            if (ignoreKey) {
                return;
            }
            if (count > 0) {
                keyboardTyped(s.subSequence(start, start + count).toString());
            }
        }

        @Override
        public void beforeTextChanged(CharSequence s, int start, int count, int after) {
            if (ignoreKey) {
                return;
            }
            if (count > 0) {
                for (int i = 0; i < count; i++) {
                    keyboardDelete();
                }
            }
        }

        @Override
        public void afterTextChanged(Editable s) {
            // always place one character so all keyboards can send backspace
            if (s.length() < 1) {
                ignoreKey = true;
                mTextEdit.setText(" ");
                mTextEdit.setSelection(mTextEdit.getText().length());
                ignoreKey = false;
            }
        }
    };

    // Accessibility – real-view overlay approach.
    private static FrameLayout mA11yContainer;

    // Staging buffer – written by Go threads, consumed on the UI thread.
    private static final List<int[]>  sStagingData   = new ArrayList<>();
    private static final List<String> sStagingLabels = new ArrayList<>();
    // Signature of the last committed layout; skip UI work when nothing changed.
    private static String sLastCommittedSignature = null;

    public GoNativeActivity() {
        super();
        goNativeActivity = this;
    }

    String getTmpdir() {
        return getCacheDir().getAbsolutePath();
    }

    void updateLayout() {
        try {
            WindowInsets insets = getWindow().getDecorView().getRootWindowInsets();
            if (insets == null) {
                return;
            }

            insetsChanged(insets.getSystemWindowInsetTop(), insets.getSystemWindowInsetBottom(),
                insets.getSystemWindowInsetLeft(), insets.getSystemWindowInsetRight());
        } catch (java.lang.NoSuchMethodError e) {
            Rect insets = new Rect();
            getWindow().getDecorView().getWindowVisibleDisplayFrame(insets);

            View view = findViewById(android.R.id.content).getRootView();
            insetsChanged(insets.top, view.getHeight() - insets.height() - insets.top,
                insets.left, view.getWidth() - insets.width() - insets.left);
        }
    }

    static void showKeyboard(int keyboardType) {
        goNativeActivity.doShowKeyboard(keyboardType);
        goNativeActivity.keyboardUp = true;
    }

    void doShowKeyboard(final int keyboardType) {
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                int imeOptions = EditorInfo.IME_FLAG_NO_ENTER_ACTION;
                int inputType = DEFAULT_INPUT_TYPE;
                String keys = "";
                switch (keyboardType) {
                    case DEFAULT_KEYBOARD_CODE:
                        imeOptions = EditorInfo.IME_FLAG_NO_ENTER_ACTION;
                        break;
                    case SINGLELINE_KEYBOARD_CODE:
                        imeOptions = EditorInfo.IME_ACTION_DONE;
                        break;
                    case NUMBER_KEYBOARD_CODE:
                        imeOptions = EditorInfo.IME_ACTION_DONE;
                        inputType |= InputType.TYPE_CLASS_NUMBER | InputType.TYPE_NUMBER_VARIATION_NORMAL;
                        keys = "0123456789.,-' "; // work around android bug where some number keys are blocked
                        break;
                    case PASSWORD_KEYBOARD_CODE:
                        imeOptions = EditorInfo.IME_ACTION_DONE;
                        inputType |= InputType.TYPE_TEXT_VARIATION_PASSWORD;
                    default:
                        Log.e("Fyne", "unknown keyboard type, use default");
                }
                mTextEdit.setImeOptions(imeOptions|EditorInfo.IME_FLAG_NO_FULLSCREEN);
                mTextEdit.setInputType(inputType);
                if (keys != "") {
                    mTextEdit.setKeyListener(DigitsKeyListener.getInstance(keys));
                }

                mTextEdit.setOnEditorActionListener(mEditorActionListener);

                // always place one character so all keyboards can send backspace
                ignoreKey = true;
                mTextEdit.setText(" ");
                mTextEdit.setSelection(mTextEdit.getText().length());
                ignoreKey = false;

                mTextEdit.setVisibility(View.VISIBLE);
                mTextEdit.bringToFront();
                mTextEdit.requestFocus();

                InputMethodManager m = (InputMethodManager) getSystemService(Context.INPUT_METHOD_SERVICE);
                m.showSoftInput(mTextEdit, 0);
            }
        });
    }

    static void hideKeyboard() {
        goNativeActivity.doHideKeyboard();
        goNativeActivity.keyboardUp = false;
    }

    void doHideKeyboard() {
        InputMethodManager imm = (InputMethodManager) getSystemService(Context.INPUT_METHOD_SERVICE);
        View view = findViewById(android.R.id.content).getRootView();
        imm.hideSoftInputFromWindow(view.getWindowToken(), 0);

        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                mTextEdit.setVisibility(View.GONE);
            }
        });
    }

    static void showFileOpen(String mimes) {
        goNativeActivity.doShowFileOpen(mimes);
    }

    void doShowFileOpen(String mimes) {
        Intent intent = new Intent(Intent.ACTION_OPEN_DOCUMENT);
        if ("application/x-directory".equals(mimes) && Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP) {
            intent = new Intent(Intent.ACTION_OPEN_DOCUMENT_TREE); // ask for a directory picker if OS supports it
            intent.addFlags(Intent.FLAG_GRANT_READ_URI_PERMISSION);
        } else if (mimes.contains("|") && Build.VERSION.SDK_INT >= Build.VERSION_CODES.KITKAT) {
            intent.setType("*/*");
            intent.putExtra(Intent.EXTRA_MIME_TYPES, mimes.split("\\|"));
            intent.addCategory(Intent.CATEGORY_OPENABLE);
        } else {
            intent.setType(mimes);
            intent.addCategory(Intent.CATEGORY_OPENABLE);
        }
        startActivityForResult(Intent.createChooser(intent, "Open File"), FILE_OPEN_CODE);
    }

    @Override
    public void onRequestPermissionsResult(int requestCode, @NonNull String[] permissions, @NonNull int[] grantResults) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults);
    
        if (requestCode == 100) {
            if (grantResults.length > 0 && grantResults[0] == PackageManager.PERMISSION_GRANTED) {
    
                // 4. Read the explicit action intent token
                int actionToExecute = pendingPermissionAction;
                pendingPermissionAction = 0; // Clear immediately to reset the gate
    
                if (actionToExecute == 1) {
                    doStartCameraPreview();
                } else if (actionToExecute == 2) {
                    doCaptureCameraPhoto();
                }
    
            } else {
                pendingPermissionAction = 0; // Clear token if user denied
            }
        }
    }

    private boolean ensureCameraInitialized(Context context, androidx.lifecycle.LifecycleOwner lifecycleOwner, Runnable onReadyTask) {
        if (isCameraRunning && globalCameraProvider != null && globalImageCapture != null && (!isPreviewRequested || preview != null)) {
            if (onReadyTask != null) onReadyTask.run();
            return true;
        }
    
        isCameraRunning = true;
    
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                ListenableFuture<ProcessCameraProvider> providerListenableFuture = ProcessCameraProvider.getInstance(context);
                providerListenableFuture.addListener(new Runnable() {
                    @Override
                    public void run() {
                        try {
                            globalCameraProvider = providerListenableFuture.get();
                            globalCameraProvider.unbindAll(); // Fresh slate
    
                            CameraSelector cameraSelector = new CameraSelector.Builder()
                                    .requireLensFacing(CameraSelector.LENS_FACING_BACK)
                                    .build();
    
                            // 1. Always build the ImageCapture use case
                            globalImageCapture = new ImageCapture.Builder()
                                    .setCaptureMode(ImageCapture.CAPTURE_MODE_MAXIMIZE_QUALITY)
                                    .setTargetRotation(getWindow().getDecorView().getRootView().getDisplay().getRotation())
                                    .build();
    
                            // 2. Build the Preview use case ONLY if requested
                            if (isPreviewRequested) {
                                preview = new Preview.Builder().build();
                                globalBitmapSurfaceProvider = new BitmapSurfaceProvider(context);
                                globalBitmapSurfaceProvider.setFrameCallback(new BitmapSurfaceProvider.FrameProcessor() {
                                    @Override
                                    public void onNewFrame(Bitmap photo) {
                                        ByteArrayOutputStream out = new ByteArrayOutputStream();
                                        photo.compress(Bitmap.CompressFormat.JPEG, 90, out);
                                        byte[] rawBytes = out.toByteArray();
                                        cameraPreviewFrame(rawBytes, rawBytes.length);
                                    }
                                });
                                preview.setSurfaceProvider(globalBitmapSurfaceProvider);
    
                                // Bind both together
                                globalCameraProvider.bindToLifecycle(lifecycleOwner, cameraSelector, preview, globalImageCapture);
                            } else {
                                // Bind just the photo use case to initialize hardware instantly without starting the video thread
                                globalCameraProvider.bindToLifecycle(lifecycleOwner, cameraSelector, globalImageCapture);
                            }
    
                            // 3. Execute any deferred task (like taking the photo) once bound
                            if (onReadyTask != null) {
                                onReadyTask.run();
                            }
    
                        } catch (Exception e) {
                            isCameraRunning = false;
                            globalImageCapture = null;
                            e.printStackTrace();
                        }
                    }
                }, ContextCompat.getMainExecutor(context));
            }
        });
    
        return false;
    }

    private void forceRestartPreviewStream() {
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                // Only attempt to restore if the preview was actually requested and active
                if (isPreviewRequested && isCameraRunning && globalCameraProvider != null) {
                    try {
                        // 1. Force unbind everything to completely reset the frozen hardware pipeline
                        globalCameraProvider.unbindAll();
    
                        CameraSelector cameraSelector = new CameraSelector.Builder()
                                .requireLensFacing(CameraSelector.LENS_FACING_BACK)
                                .build();
    
                        // 2. Clear out any locked up buffers from the old surface provider
                        if (globalBitmapSurfaceProvider != null) {
                            globalBitmapSurfaceProvider.cleanUp();
                        }
    
                        // 3. Instantiate a fresh provider with clean ImageReader slots
                        globalBitmapSurfaceProvider = new BitmapSurfaceProvider(GoNativeActivity.this);
                        globalBitmapSurfaceProvider.setFrameCallback(new BitmapSurfaceProvider.FrameProcessor() {
                            @Override
                            public void onNewFrame(Bitmap photo) {
                                ByteArrayOutputStream out = new ByteArrayOutputStream();
                                photo.compress(Bitmap.CompressFormat.JPEG, 90, out);
                                byte[] rawBytes = out.toByteArray();
                                cameraPreviewFrame(rawBytes, rawBytes.length);
                            }
                        });
    
                        preview = new Preview.Builder().build();
                        preview.setSurfaceProvider(globalBitmapSurfaceProvider);
    
                        // 4. Rebind both use cases cleanly to the lens
                        globalCameraProvider.bindToLifecycle(GoNativeActivity.this, cameraSelector, preview, globalImageCapture);
    
                    } catch (Exception e) {
                        e.printStackTrace();
                    }
                }
            }
        });
    }

    static void startCameraPreview() {
        goNativeActivity.doStartCameraPreview();
    }

    void doStartCameraPreview() {
        // If the preview is already up and running, we can safely skip
        if (isPreviewRequested && isCameraRunning) {
            return;
        }

        if (ContextCompat.checkSelfPermission(this, Manifest.permission.CAMERA) != PackageManager.PERMISSION_GRANTED) {
            pendingPermissionAction = 1; 
            ActivityCompat.requestPermissions((Activity) this, new String[] { Manifest.permission.CAMERA }, 100);
            return;
        }    
    
        isPreviewRequested = true;
        Context context = this;
        androidx.lifecycle.LifecycleOwner lifecycleOwner = this;
    
        // 1. FIXED SCENARIO: Camera is already running (from a prior photo), but preview isn't bound yet
        if (isCameraRunning && globalCameraProvider != null) {
            runOnUiThread(new Runnable() {
                @Override
                public void run() {
                    try {
                        // Unbind everything to cleanly reconfigure the lens options
                        globalCameraProvider.unbindAll();
    
                        CameraSelector cameraSelector = new CameraSelector.Builder()
                                .requireLensFacing(CameraSelector.LENS_FACING_BACK)
                                .build();
    
                        // Rebuild the surface preview use case
                        preview = new Preview.Builder().build();
                        globalBitmapSurfaceProvider = new BitmapSurfaceProvider(context);
                        globalBitmapSurfaceProvider.setFrameCallback(new BitmapSurfaceProvider.FrameProcessor() {
                            @Override
                            public void onNewFrame(Bitmap photo) {
                                ByteArrayOutputStream out = new ByteArrayOutputStream();
                                photo.compress(Bitmap.CompressFormat.JPEG, 90, out);
                                byte[] rawBytes = out.toByteArray();
                                cameraPreviewFrame(rawBytes, rawBytes.length);
                            }
                        });
                        preview.setSurfaceProvider(globalBitmapSurfaceProvider);
    
                        // Rebind BOTH use cases together cleanly
                        globalCameraProvider.bindToLifecycle(lifecycleOwner, cameraSelector, preview, globalImageCapture);
                    } catch (Exception e) {
                        e.printStackTrace();
                    }
                }
            });
        } else {
            // 2. Standard cold start sequence if nothing was running yet
            ensureCameraInitialized(this, this, null);
        }
    }

    static void captureCameraPhoto() {
        goNativeActivity.doCaptureCameraPhoto();
    }

    void doCaptureCameraPhoto() {
        isPreviewRequested = false;
        if (ContextCompat.checkSelfPermission(this, Manifest.permission.CAMERA) != PackageManager.PERMISSION_GRANTED) {
            pendingPermissionAction = 2; 
            ActivityCompat.requestPermissions((Activity) this, new String[] { Manifest.permission.CAMERA }, 100);
            return;
        }
    
        Context context = this;
        androidx.lifecycle.LifecycleOwner lifecycleOwner = this;
    
        // Define the photo-taking logic
        Runnable takePhotoRunnable = new Runnable() {
            @Override
            public void run() {
                if (globalImageCapture == null) return;
    
                globalImageCapture.takePicture(
                    ContextCompat.getMainExecutor(context),
                    new ImageCapture.OnImageCapturedCallback() {
                        @Override
                        public void onCaptureSuccess(@NonNull androidx.camera.core.ImageProxy image) {
                            try {
                                ByteBuffer buffer = image.getPlanes()[0].getBuffer();
                                byte[] bytes = new byte[buffer.remaining()];
                                buffer.get(bytes);
                                
                                Bitmap photo = BitmapFactory.decodeByteArray(bytes, 0, bytes.length);
                                Matrix matrix = new Matrix();
                                matrix.postRotate(image.getImageInfo().getRotationDegrees());
                                Bitmap rotated = Bitmap.createBitmap(photo, 0, 0, photo.getWidth(), photo.getHeight(), matrix, true);
                                
                                ByteArrayOutputStream out = new ByteArrayOutputStream();
                                rotated.compress(Bitmap.CompressFormat.JPEG, 95, out);
                                byte[] rawBytes = out.toByteArray();
                                
                                capturePhotoReturned(rawBytes, rawBytes.length);
                                
                                photo.recycle();
                                rotated.recycle();
                            } catch (Exception e) {
                                e.printStackTrace();
                            } finally {
                                image.close();
                            }
    
                            // Restore preview stream if it was active
                            runOnUiThread(new Runnable() {
                                @Override
                                public void run() {
                                    if (isPreviewRequested && preview != null && globalBitmapSurfaceProvider != null) {
                    forceRestartPreviewStream();
                                    }
                                }
                            });
                        }
    
                        @Override
                        public void onError(@NonNull ImageCaptureException exception) {
                            exception.printStackTrace();
                            runOnUiThread(new Runnable() {
                                @Override
                                public void run() {
                                    if (isPreviewRequested && preview != null && globalBitmapSurfaceProvider != null) {
                    forceRestartPreviewStream();
                                    }
                                }
                            });
                        }
                    }
                );
            }
        };
    
        // Ensure the system is ready. If it's already open, it runs instantly. 
        // If it's closed, it initializes everything first, then fires the runnable.
        boolean initialized = ensureCameraInitialized(context, lifecycleOwner, takePhotoRunnable);
        if (initialized) {
            // Camera was already fully initialized; fire immediately on the main thread
            runOnUiThread(takePhotoRunnable);
        }
    }

    static void stopCameraPreview() {
        goNativeActivity.doStopCameraPreview();
    }

    public void doStopCameraPreview() {
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                // 1. Reset state tracking variables immediately
                isCameraRunning = false;
                isPreviewRequested = false;
                
                // 2. Unbind all active use cases from the Android system lens
                if (globalCameraProvider != null) {
                    try {
                        globalCameraProvider.unbindAll();
                    } catch (Exception e) {
                        e.printStackTrace();
                    }
                    globalCameraProvider = null;
                }
                
                // 3. Terminate background threads and clear the custom ImageReader
                if (globalBitmapSurfaceProvider != null) {
                    try {
                        globalBitmapSurfaceProvider.cleanUp(); // Quits the HandlerThread and drops old frames
                    } catch (Exception e) {
                        e.printStackTrace();
                    }
                    globalBitmapSurfaceProvider = null;
                }
                
                // 4. Nullify remaining references
                globalImageCapture = null;
                preview = null;
            }
        });
    }

    static void showFileSave(String mimes, String filename) {
        goNativeActivity.doShowFileSave(mimes, filename);
    }

    void doShowFileSave(String mimes, String filename) {
        Intent intent = new Intent(Intent.ACTION_CREATE_DOCUMENT);
        if (mimes.contains("|") && Build.VERSION.SDK_INT >= Build.VERSION_CODES.KITKAT) {
            intent.setType("*/*");
            intent.putExtra(Intent.EXTRA_MIME_TYPES, mimes.split("\\|"));
        } else {
            intent.setType(mimes);
        }
        intent.putExtra(Intent.EXTRA_TITLE, filename);
        intent.addCategory(Intent.CATEGORY_OPENABLE);
        startActivityForResult(Intent.createChooser(intent, "Save File"), FILE_SAVE_CODE);
    }

    // -------------------------------------------------------------------------
    // Scheduled notifications via AlarmManager.
    //
    // For delivery to survive the app process being killed, the FyneNotificationReceiver
    // inner class must be declared in AndroidManifest.xml:
    //
    //   <receiver android:name="org.golang.app.FyneNotificationReceiver"
    //             android:exported="false" />
    //
    // The packaging tool is responsible for emitting that line. If the receiver is not
    // registered the schedule call returns false and the Go layer falls back to an
    // in-process scheduler. Cancellation always uses the same identifier so the same
    // PendingIntent can be reproduced and removed from AlarmManager.
    // -------------------------------------------------------------------------

    static boolean scheduleNotification(String id, String title, String body, long deliveryTimeMillis) {
        if (goNativeActivity == null) {
            return false;
        }
        return goNativeActivity.doScheduleNotification(id, title, body, deliveryTimeMillis);
    }

    boolean doScheduleNotification(String id, String title, String body, long deliveryTimeMillis) {
        // If the receiver is not registered in the manifest the scheduled alarm
        // will fire into the void once the app is killed; report failure so the
        // Go layer can fall back to in-process scheduling instead.
        ComponentName receiver = new ComponentName(this, FyneNotificationReceiver.class);
        try {
            getPackageManager().getReceiverInfo(receiver, 0);
        } catch (PackageManager.NameNotFoundException e) {
            return false;
        }

        AlarmManager alarmMgr = (AlarmManager) getSystemService(Context.ALARM_SERVICE);
        if (alarmMgr == null) {
            return false;
        }

        Intent intent = new Intent(this, FyneNotificationReceiver.class);
        intent.putExtra("title", title);
        intent.putExtra("body", body);
        intent.putExtra("notif_id", id.hashCode());

        int flags = PendingIntent.FLAG_UPDATE_CURRENT;
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) {
            flags |= PendingIntent.FLAG_IMMUTABLE;
        }

        PendingIntent pi = PendingIntent.getBroadcast(this, id.hashCode(), intent, flags);
        try {
            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
                alarmMgr.setAndAllowWhileIdle(AlarmManager.RTC_WAKEUP, deliveryTimeMillis, pi);
            } else {
                alarmMgr.set(AlarmManager.RTC_WAKEUP, deliveryTimeMillis, pi);
            }
        } catch (SecurityException e) {
            Log.e("Fyne", "AlarmManager rejected scheduled notification", e);
            return false;
        }
        return true;
    }

    static void cancelScheduledNotification(String id) {
        if (goNativeActivity == null) {
            return;
        }
        goNativeActivity.doCancelScheduledNotification(id);
    }

    void doCancelScheduledNotification(String id) {
        AlarmManager alarmMgr = (AlarmManager) getSystemService(Context.ALARM_SERVICE);
        if (alarmMgr == null) {
            return;
        }

        Intent intent = new Intent(this, FyneNotificationReceiver.class);
        int flags = PendingIntent.FLAG_NO_CREATE;
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.S) {
            flags |= PendingIntent.FLAG_IMMUTABLE;
        }

        PendingIntent pi = PendingIntent.getBroadcast(this, id.hashCode(), intent, flags);
        if (pi != null) {
            alarmMgr.cancel(pi);
            pi.cancel();
        }

        // Also drop any already-posted notification with the same id so cancelling
        // after delivery clears it from the shade.
        NotificationManager mgr = (NotificationManager) getSystemService(Context.NOTIFICATION_SERVICE);
        if (mgr != null) {
            mgr.cancel(id.hashCode());
        }
    }

    static int getRune(int deviceId, int keyCode, int metaState) {
        try {
            int rune = KeyCharacterMap.load(deviceId).get(keyCode, metaState);
            if (rune == 0) {
                return -1;
            }
            return rune;
        } catch (KeyCharacterMap.UnavailableException e) {
            return -1;
        } catch (Exception e) {
            Log.e("Fyne", "exception reading KeyCharacterMap", e);
            return -1;
        }
    }


    private void load() {
        // Interestingly, NativeActivity uses a different method
        // to find native code to execute, avoiding
        // System.loadLibrary. The result is Java methods
        // implemented in C with JNIEXPORT (and JNI_OnLoad) are not
        // available unless an explicit call to System.loadLibrary
        // is done. So we do it here, borrowing the name of the
        // library from the same AndroidManifest.xml metadata used
        // by NativeActivity.
        try {
            ActivityInfo ai = getPackageManager().getActivityInfo(
                    getIntent().getComponent(), PackageManager.GET_META_DATA);
            if (ai.metaData == null) {
                Log.e("Fyne", "loadLibrary: no manifest metadata found");
                return;
            }
            String libName = ai.metaData.getString("android.app.lib_name");
            System.loadLibrary(libName);
        } catch (Exception e) {
            Log.e("Fyne", "loadLibrary failed", e);
        }
    }

    @Override
    public void onCreate(Bundle savedInstanceState) {
        load();
        super.onCreate(savedInstanceState);
        setupEntry();
        updateTheme(getResources().getConfiguration());

        View view = findViewById(android.R.id.content).getRootView();
        view.addOnLayoutChangeListener(new View.OnLayoutChangeListener() {
            public void onLayoutChange (View v, int left, int top, int right, int bottom,
                                        int oldLeft, int oldTop, int oldRight, int oldBottom) {
                GoNativeActivity.this.updateLayout();
            }
        });
        ProcessCameraProvider.configureInstance(Camera2Config.defaultConfig());
        lifecycleRegistry.handleLifecycleEvent(Lifecycle.Event.ON_CREATE);
    }

    private void setupEntry() {
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                mTextEdit = new EditText(goNativeActivity);
                mTextEdit.setVisibility(View.GONE);
                mTextEdit.setInputType(DEFAULT_INPUT_TYPE);

                FrameLayout.LayoutParams mEditTextLayoutParams = new FrameLayout.LayoutParams(
                    FrameLayout.LayoutParams.WRAP_CONTENT, FrameLayout.LayoutParams.WRAP_CONTENT);
                mTextEdit.setLayoutParams(mEditTextLayoutParams);
                addContentView(mTextEdit, mEditTextLayoutParams);

                // always place one character so all keyboards can send backspace
                mTextEdit.setText(" ");
                mTextEdit.setSelection(mTextEdit.getText().length());

                mTextEdit.addTextChangedListener(mTextWatcher);
            }
        });
    }

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        // unhandled request
        if (requestCode != FILE_OPEN_CODE && requestCode != FILE_SAVE_CODE && requestCode != CAMERA_OPEN_CODE) {
            return;
        }

        if (resultCode != Activity.RESULT_OK) {
            // dialog was cancelled
            filePickerReturned("");
            return;
        }

        Uri uri = data.getData();
        filePickerReturned(uri.toString());
    }

    @Override
    public void onBackPressed() {
        if (goNativeActivity.keyboardUp) {
            hideKeyboard();
            return;
        }

        // skip the default behaviour - we can call finishActivity if we want to go back
        backPressed();
    }

    public void finishActivity() {
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                GoNativeActivity.super.onBackPressed();
            }
        });
    }

    @Override
    public void onConfigurationChanged(Configuration config) {
        super.onConfigurationChanged(config);
        updateTheme(config);
    }

    protected void updateTheme(Configuration config) {
        boolean dark = (config.uiMode & Configuration.UI_MODE_NIGHT_MASK) == Configuration.UI_MODE_NIGHT_YES;
        setDarkMode(dark);
        updateSystemBarsAppearance(dark);
    }

    private void updateSystemBarsAppearance(boolean dark) {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.R) {
            WindowInsetsController controller = getWindow().getInsetsController();
            if (controller != null) {
                int lightFlags = WindowInsetsController.APPEARANCE_LIGHT_STATUS_BARS |
                        WindowInsetsController.APPEARANCE_LIGHT_NAVIGATION_BARS;
                controller.setSystemBarsAppearance(dark ? 0 : lightFlags, lightFlags);
            }
        } else if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
            View decorView = getWindow().getDecorView();
            int flags = decorView.getSystemUiVisibility();
            if (dark) {
                flags &= ~View.SYSTEM_UI_FLAG_LIGHT_STATUS_BAR;
                if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                    flags &= ~View.SYSTEM_UI_FLAG_LIGHT_NAVIGATION_BAR;
                }
            } else {
                flags |= View.SYSTEM_UI_FLAG_LIGHT_STATUS_BAR;
                if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                    flags |= View.SYSTEM_UI_FLAG_LIGHT_NAVIGATION_BAR;
                }
            }
            decorView.setSystemUiVisibility(flags);
        }
    }

    // -------------------------------------------------------------------------
    // Android accessibility bridge  (real-view overlay approach)
    // -------------------------------------------------------------------------

    static final int ROLE_BUTTON    = 1;
    static final int ROLE_TEXT      = 2;
    static final int ROLE_LINK      = 3;
    static final int ROLE_CONTAINER = 4;

    // Called once from Go (via JNI).
    static void setupAccessibility() {
        goNativeActivity.doSetupAccessibility();
    }

    void doSetupAccessibility() {
        runOnUiThread(new Runnable() {
            @Override
            public void run() {
                if (mA11yContainer != null) {
                    return; // already set up
                }
                // Full-screen transparent container.  Its children are the real
                // accessibility Views; they intercept no touch events.
                mA11yContainer = new FrameLayout(goNativeActivity);
                mA11yContainer.setClickable(false);
                mA11yContainer.setFocusable(false);

                mA11yContainer.setImportantForAccessibility(
                        View.IMPORTANT_FOR_ACCESSIBILITY_AUTO);

                ViewGroup decorView = (ViewGroup) getWindow().getDecorView();
                decorView.addView(mA11yContainer, new FrameLayout.LayoutParams(
                        FrameLayout.LayoutParams.MATCH_PARENT,
                        FrameLayout.LayoutParams.MATCH_PARENT));

                // If Go already committed nodes before the UI thread ran, apply them now.
                applySnapshot();
            }
        });
    }

    // Called from Go (via JNI) before rebuilding the accessibility tree.
    static synchronized void clearAccessibilityNodes() {
        sStagingData.clear();
        sStagingLabels.clear();
    }

    // Called from Go (via JNI) to register one accessible element.
    static synchronized void addAccessibilityNode(int id, int role, String label,
            int x, int y, int width, int height, int parentID) {
        sStagingData.add(new int[]{id, role, x, y, width, height});
        sStagingLabels.add(label != null ? label : "");
    }

    // Called from Go (via JNI) after all nodes for this frame have been added.
    static synchronized void commitAccessibilityNodes() {
        // Build a cheap signature to avoid redundant UI work.
        StringBuilder sig = new StringBuilder();
        for (int i = 0; i < sStagingData.size(); i++) {
            int[] d = sStagingData.get(i);
            sig.append(d[0]).append(',').append(d[1]).append(',')
               .append(d[2]).append(',').append(d[3]).append(',')
               .append(d[4]).append(',').append(d[5]).append(';')
               .append(sStagingLabels.get(i)).append('|');
        }
        String newSig = sig.toString();
        boolean changed = !newSig.equals(sLastCommittedSignature);
        sLastCommittedSignature = newSig;

        if (!changed) {
            return;
        }

        // Snapshot the staging data for the UI thread (sStagingData is modified
        // on Go threads so we must not access it directly from the UI thread).
        final List<int[]>  snapData   = new ArrayList<>(sStagingData);
        final List<String> snapLabels = new ArrayList<>(sStagingLabels);

        if (mA11yContainer == null) {
            return; // UI not ready yet; doSetupAccessibility will call applySnapshot.
        }
        mA11yContainer.post(new Runnable() {
            @Override
            public void run() {
                rebuildA11yViews(snapData, snapLabels);
            }
        });
    }

    // Applies the most-recently committed snapshot.  Called on the UI thread
    // either from doSetupAccessibility (if setup ran after commit) or from the
    // post() in commitAccessibilityNodes.
    private static void applySnapshot() {
        // Take a consistent copy of the current staging data.
        final List<int[]>  snapData;
        final List<String> snapLabels;
        synchronized (GoNativeActivity.class) {
            snapData   = new ArrayList<>(sStagingData);
            snapLabels = new ArrayList<>(sStagingLabels);
        }
        rebuildA11yViews(snapData, snapLabels);
    }

    // Must be called on the UI thread.
    private static void rebuildA11yViews(List<int[]> data, List<String> labels) {
        if (mA11yContainer == null) {
            return;
        }
        mA11yContainer.removeAllViews();

        for (int i = 0; i < data.size(); i++) {
            final int[] d     = data.get(i);
            final String label = labels.get(i);
            final int role  = d[1];
            final int x = d[2], y = d[3], w = d[4], h = d[5];

            View v = new View(goNativeActivity);
            // Non-interactive so touch events fall through to the GL surface.
            // No background is set, so the view is visually transparent.
            // setAlpha is deliberately NOT called: alpha=0 causes TalkBack to
            // mark the view as not visible to the user and skip it.
            v.setClickable(false);
            v.setFocusable(false);
            v.setContentDescription(label);
            v.setImportantForAccessibility(View.IMPORTANT_FOR_ACCESSIBILITY_YES);

            // Customise the AccessibilityNodeInfo so TalkBack announces the
            // correct role (Button vs label) and actions.
            final int roleCopy = role;
            final String labelCopy = label;
            v.setAccessibilityDelegate(new View.AccessibilityDelegate() {
                @Override
                public void onInitializeAccessibilityNodeInfo(
                        View host, AccessibilityNodeInfo info) {
                    super.onInitializeAccessibilityNodeInfo(host, info);
                    info.setContentDescription(labelCopy);
                    switch (roleCopy) {
                        case ROLE_BUTTON:
                        case ROLE_LINK:
                            info.setClassName("android.widget.Button");
                            info.setClickable(true);
                            info.addAction(AccessibilityNodeInfo.ACTION_CLICK);
                            break;
                        default: // ROLE_TEXT
                            info.setClassName("android.widget.TextView");
                            info.setText(labelCopy);
                            break;
                    }
                }
            });

            // Position the view at the widget's screen-pixel coordinates.
            FrameLayout.LayoutParams lp = new FrameLayout.LayoutParams(w, h);
            lp.leftMargin = x;
            lp.topMargin  = y;
            mA11yContainer.addView(v, lp);
        }
    }

    @Override
    public Lifecycle getLifecycle() {
        return lifecycleRegistry;
    }

    @Override
    protected void onStart() {
        super.onStart();
        lifecycleRegistry.handleLifecycleEvent(Lifecycle.Event.ON_START);
    }

    @Override
    protected void onResume() {
        super.onResume();
        lifecycleRegistry.handleLifecycleEvent(Lifecycle.Event.ON_RESUME);
    }

    @Override
    protected void onPause() {
        lifecycleRegistry.handleLifecycleEvent(Lifecycle.Event.ON_PAUSE);
        super.onPause();
        doStopCameraPreview();
    }

    @Override
    protected void onStop() {
        lifecycleRegistry.handleLifecycleEvent(Lifecycle.Event.ON_STOP);
        super.onStop();
        doStopCameraPreview();
    }

    @Override
    protected void onDestroy() {
        lifecycleRegistry.handleLifecycleEvent(Lifecycle.Event.ON_DESTROY);
        super.onDestroy();
        doStopCameraPreview();
    }
}
