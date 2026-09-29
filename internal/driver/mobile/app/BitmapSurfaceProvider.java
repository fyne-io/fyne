package org.golang.app;

import android.content.Context;
import android.graphics.Bitmap;
import android.graphics.BitmapFactory;
import android.graphics.ImageFormat;
import android.graphics.Matrix;
import android.graphics.Rect;
import android.graphics.YuvImage;
import android.media.Image;
import android.media.ImageReader;
import android.os.Handler;
import android.os.HandlerThread;
import android.util.Size;
import android.view.Surface;
import androidx.annotation.NonNull;
import androidx.camera.core.Preview;
import androidx.camera.core.SurfaceRequest;
import androidx.core.content.ContextCompat;
import java.io.ByteArrayOutputStream;
import java.nio.ByteBuffer;
import java.util.concurrent.Executor;

public class BitmapSurfaceProvider implements Preview.SurfaceProvider {
    // Functional interface that receives each newly generated frame
    public interface FrameProcessor {
        void onNewFrame(Bitmap bitmap);
    }

    private final Context context;
    private volatile Bitmap latestBitmap = null;
    private ImageReader imageReader;
    private HandlerThread handlerThread;
    private volatile FrameProcessor frameCallback; // Hook for running logic on every frame
    private volatile int rotationDegrees = 0;
    private final Object lock = new Object();

    public BitmapSurfaceProvider(Context context) {
        this.context = context.getApplicationContext();
    }


    /**
     * Registers a callback to execute on every single incoming camera frame.
     * Note: This runs on the background camera processing thread to avoid blocking the UI.
     */
    public void setFrameCallback(FrameProcessor callback) {
        this.frameCallback = callback;
    }

    public Bitmap getLatestBitmap() {
        synchronized (lock) {
            return latestBitmap;
        }
    }

    @Override
    public void onSurfaceRequested(@NonNull SurfaceRequest request) {
        Size size = request.getResolution();
        
        handlerThread = new HandlerThread("CameraBitmapThread");
        handlerThread.start();
        Handler backgroundHandler = new Handler(handlerThread.getLooper());
        Executor backgroundExecutor = ContextCompat.getMainExecutor(this.context); 

        request.setTransformationInfoListener(backgroundExecutor, new SurfaceRequest.TransformationInfoListener() {
            @Override
            public void onTransformationInfoUpdate(@NonNull SurfaceRequest.TransformationInfo info) {
                rotationDegrees = info.getRotationDegrees();
            }
        });

        imageReader = ImageReader.newInstance(
                size.getWidth(), 
                size.getHeight(), 
                ImageFormat.YUV_420_888, 
                4
        );

        imageReader.setOnImageAvailableListener(new ImageReader.OnImageAvailableListener() {
            @Override
            public void onImageAvailable(ImageReader reader) {
                if (reader == null) return;
                Image image = null;
                try {
                    image = reader.acquireLatestImage();
                    if (image == null) return;

                    Bitmap rawBitmap = yuv420ToBitmapSafely(image);
                    if (rawBitmap != null) {
                        Bitmap correctedBitmap = rotateBitmap(rawBitmap, rotationDegrees);
                        
                        synchronized (lock) {
                            Bitmap oldBitmap = latestBitmap;
                            latestBitmap = correctedBitmap;
                            
                            if (oldBitmap != null && oldBitmap != correctedBitmap) {
                                oldBitmap.recycle();
                            }
                        }

                        // Execute the custom callback for the frame
                        FrameProcessor currentCallback = frameCallback;
                        if (currentCallback != null) {
                            currentCallback.onNewFrame(correctedBitmap);
                        }
                    }
                } catch (Exception e) {
                    e.printStackTrace();
                } finally {
                    if (image != null) {
                        image.close();
                    }
                }
            }
        }, backgroundHandler);

        Surface surface = imageReader.getSurface();
        
        request.provideSurface(surface, backgroundExecutor, new androidx.core.util.Consumer<SurfaceRequest.Result>() {
            @Override
            public void accept(SurfaceRequest.Result result) {
                cleanUp();
            }
        });
    }

    public void cleanUp() {
        synchronized (lock) {
            frameCallback = null; // Clear reference to prevent memory leaks
            if (imageReader != null) {
                imageReader.close();
                imageReader = null;
            }
            if (handlerThread != null) {
                handlerThread.quitSafely();
                handlerThread = null;
            }
            if (latestBitmap != null) {
                latestBitmap.recycle();
                latestBitmap = null;
            }
        }
    }

    private Bitmap rotateBitmap(Bitmap source, int degrees) {
        if (degrees == 0) return source;
        Matrix matrix = new Matrix();
        matrix.postRotate(degrees);
        Bitmap rotated = Bitmap.createBitmap(
                source, 0, 0, source.getWidth(), source.getHeight(), matrix, true
        );
        source.recycle();
        return rotated;
    }

    private Bitmap yuv420ToBitmapSafely(Image image) {
        try {
            int width = image.getWidth();
            int height = image.getHeight();
    
            Image.Plane[] planes = image.getPlanes();
            Image.Plane yPlane = planes[0]; // Y Plane is always index 0
            Image.Plane uPlane = planes[1]; // U Plane is always index 1
            Image.Plane vPlane = planes[2]; // V Plane is always index 2
    
            ByteBuffer yBuffer = yPlane.getBuffer();
            ByteBuffer uBuffer = uPlane.getBuffer();
            ByteBuffer vBuffer = vPlane.getBuffer();
    
            int ySize = yBuffer.remaining();
            int uSize = uBuffer.remaining();
            int vSize = vBuffer.remaining();
    
            // NV21 requires Y size + (UV size)
            byte[] nv21 = new byte[ySize + (width * height / 2)];
    
            // Copy Y data cleanly
            yBuffer.get(nv21, 0, ySize);
    
            // Extract interleaved pixel strides
            int vRowStride = vPlane.getRowStride();
            int vPixelStride = vPlane.getPixelStride();
    
            int nvIndex = ySize;
            byte[] vData = new byte[vSize];
            byte[] uData = new byte[uSize];
    
            vBuffer.get(vData);
            uBuffer.get(uData);
    
            // Correctly stitch the NV21 frame together
            for (int row = 0; row < height / 2; row++) {
                for (int col = 0; col < width / 2; col++) {
                    int vIndex = row * vRowStride + col * vPixelStride;
                    if (vIndex < vData.length && vIndex < uData.length) {
                        nv21[nvIndex++] = vData[vIndex]; // V byte
                        nv21[nvIndex++] = uData[vIndex]; // U byte
                    }
                }
            }
    
            // Compress data to JPEG
            ByteArrayOutputStream out = new ByteArrayOutputStream();
            YuvImage yuvImage = new YuvImage(nv21, ImageFormat.NV21, width, height, null);
            yuvImage.compressToJpeg(new Rect(0, 0, width, height), 80, out);
            byte[] imageBytes = out.toByteArray();
    
            out.close();
            return BitmapFactory.decodeByteArray(imageBytes, 0, imageBytes.length);
    
        } catch (Exception e) {
            // Look at your Logcat console for this tag if it fails!
            android.util.Log.e("BitmapProvider", "Failed to extract YUV frame data", e);
            return null;
        }
    }
}
