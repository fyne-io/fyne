package org.golang.app;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.Service;
import android.content.ComponentName;
import android.content.Context;
import android.content.Intent;
import android.os.Build;
import android.os.IBinder;
import android.os.Looper;
import android.util.Log;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

// GoForegroundService is an Android foreground service that does no work of
// its own. It exists only so that the system treats the hosting process as
// foreground, which makes it far less likely to be killed while the app is in
// the background.
//
// Android requires every foreground service to show an ongoing notification,
// so the text for it is supplied by the Go caller. An app that wants the
// service declares it in its own AndroidManifest.xml, along with the foreground
// service type and the permissions it needs; the manifest that the Fyne
// packaging tool generates does not include it.
public class GoForegroundService extends Service {
	private static final String TAG = "Fyne";
	private static final String CHANNEL_ID = "fyne_foreground_service";
	private static final int NOTIFICATION_ID = 1;

	private static final String EXTRA_TITLE = "title";
	private static final String EXTRA_CONTENT = "content";

	// UNKNOWN_APP_ICON is android.R.drawable.sym_def_app_icon, used when the
	// app was packaged without an icon of its own. A notification with no
	// small icon is not displayed at all, which would stop the service.
	private static final int UNKNOWN_APP_ICON = 17629184;

	// START_TIMEOUT_SECONDS is how long start waits for the system to deliver
	// the start request to the service.
	private static final long START_TIMEOUT_SECONDS = 5;

	// running is true while the service is in the foreground.
	private static volatile boolean running;

	// pendingStart and startError carry the outcome of onStartCommand, which
	// the system runs later on the main thread, back to the start call that
	// is waiting for it.
	private static volatile CountDownLatch pendingStart;
	private static volatile String startError;

	// start puts the service in the foreground, showing an ongoing
	// notification with the given text. Calling it again while the service is
	// running only updates the text. It returns null on success, or a
	// description of why the service could not be started.
	// Called from Go through JNI.
	public static synchronized String start(Context ctx, String title, String content) {
		if (running) {
			// Already in the foreground. Asking for another start is not
			// allowed once the app is in the background (Android 12+), so
			// replace the notification in place instead.
			try {
				NotificationManager manager =
						(NotificationManager) ctx.getSystemService(Context.NOTIFICATION_SERVICE);
				manager.notify(NOTIFICATION_ID, notification(ctx, title, content));
			} catch (RuntimeException e) {
				return failure("cannot update the foreground service notification", e);
			}
			return null;
		}

		Intent intent = new Intent(ctx, GoForegroundService.class);
		intent.putExtra(EXTRA_TITLE, title);
		intent.putExtra(EXTRA_CONTENT, content);

		CountDownLatch started = new CountDownLatch(1);
		startError = null;
		pendingStart = started;

		try {
			if (startService(ctx, intent) == null) {
				pendingStart = null;
				return failure("cannot start the foreground service: " + GoForegroundService.class.getName() +
						" is not declared in AndroidManifest.xml", null);
			}
		} catch (RuntimeException e) {
			pendingStart = null;
			return failure("cannot start the foreground service", e);
		}

		if (Looper.myLooper() == Looper.getMainLooper()) {
			// onStartCommand runs on this thread, so waiting for it here would
			// only ever time out. Failures are still logged by the service.
			return null;
		}

		try {
			if (!started.await(START_TIMEOUT_SECONDS, TimeUnit.SECONDS)) {
				return failure("timed out waiting for the foreground service to start", null);
			}
		} catch (InterruptedException e) {
			Thread.currentThread().interrupt();
			return failure("interrupted while waiting for the foreground service to start", null);
		}
		return startError;
	}

	// stop shuts the service down and removes its notification.
	// Called from Go through JNI.
	public static void stop(Context ctx) {
		ctx.stopService(new Intent(ctx, GoForegroundService.class));
	}

	// startService asks the system to create the service, returning null when
	// the app does not declare it.
	//
	// A plain startService is preferred: unlike startForegroundService it does
	// not oblige the service to reach the foreground, so a failed promotion can
	// be reported and the service stopped without the system killing the app.
	// It is refused while the app is in the background, where
	// startForegroundService is the only option.
	private static ComponentName startService(Context ctx, Intent intent) {
		try {
			return ctx.startService(intent);
		} catch (IllegalStateException e) {
			if (Build.VERSION.SDK_INT < Build.VERSION_CODES.O) {
				throw e;
			}
			return ctx.startForegroundService(intent);
		}
	}

	private static String failure(String message, Throwable cause) {
		if (cause != null) {
			message += ": " + cause;
		}
		Log.e(TAG, message, cause);
		return message;
	}

	@Override
	public void onCreate() {
		super.onCreate();

		if (Build.VERSION.SDK_INT < Build.VERSION_CODES.O) {
			return; // notification channels were added in API 26
		}

		NotificationChannel channel = new NotificationChannel(CHANNEL_ID,
				"Background service", NotificationManager.IMPORTANCE_LOW);
		channel.setDescription("Keeps the app running in the background");
		getSystemService(NotificationManager.class).createNotificationChannel(channel);
	}

	@Override
	public int onStartCommand(Intent intent, int flags, int startId) {
		String error = null;
		if (intent == null) {
			// Re-delivered by the system without our extras, so there is
			// nothing to show. START_NOT_STICKY means this should not happen.
			error = failure("foreground service restarted without a notification to show", null);
		} else {
			try {
				startForeground(NOTIFICATION_ID, notification(this,
						intent.getStringExtra(EXTRA_TITLE), intent.getStringExtra(EXTRA_CONTENT)));
				running = true;
			} catch (RuntimeException e) {
				// Android refuses the promotion when, for example, the
				// permission for the declared foreground service type is
				// missing (14+) or the app is in the background (12+).
				error = failure("cannot move the service to the foreground", e);
			}
		}

		if (error != null) {
			stopSelf();
		}

		startError = error;
		CountDownLatch started = pendingStart;
		if (started != null) {
			started.countDown();
		}
		return START_NOT_STICKY;
	}

	@Override
	public void onDestroy() {
		running = false;
		super.onDestroy();
	}

	@Override
	public IBinder onBind(Intent intent) {
		return null; // not a bound service
	}

	private static Notification notification(Context ctx, String title, String content) {
		Notification.Builder builder;
		if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
			builder = new Notification.Builder(ctx, CHANNEL_ID);
		} else {
			builder = new Notification.Builder(ctx);
		}

		int icon = ctx.getApplicationInfo().icon;
		if (icon == 0) {
			icon = UNKNOWN_APP_ICON;
		}

		return builder.setContentTitle(title == null ? "" : title).
				setContentText(content == null ? "" : content).
				setSmallIcon(icon).
				setOngoing(true).
				build();
	}
}
