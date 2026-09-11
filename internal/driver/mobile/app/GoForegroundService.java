package org.golang.app;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.app.Service;
import android.content.Context;
import android.content.Intent;
import android.os.Build;
import android.os.IBinder;

// GoForegroundService is an Android foreground service that does no work of
// its own. It exists only so that the system treats the hosting process as
// foreground, which makes it far less likely to be killed while the app is in
// the background.
//
// Android requires every foreground service to show an ongoing notification,
// so the text for it is supplied by the Go caller. The service is declared in
// AndroidManifest.xml by the Fyne packaging tool, along with the permissions
// and the foreground service type it needs.
public class GoForegroundService extends Service {
	private static final String CHANNEL_ID = "fyne_foreground_service";
	private static final int NOTIFICATION_ID = 1;

	private static final String EXTRA_TITLE = "title";
	private static final String EXTRA_CONTENT = "content";

	// UNKNOWN_APP_ICON is android.R.drawable.sym_def_app_icon, used when the
	// app was packaged without an icon of its own. A notification with no
	// small icon is not displayed at all, which would stop the service.
	private static final int UNKNOWN_APP_ICON = 17629184;

	// start puts the service in the foreground, showing an ongoing
	// notification with the given text. Calling it again updates the text.
	// Called from Go through JNI.
	public static void start(Context ctx, String title, String content) {
		Intent intent = new Intent(ctx, GoForegroundService.class);
		intent.putExtra(EXTRA_TITLE, title);
		intent.putExtra(EXTRA_CONTENT, content);

		if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
			ctx.startForegroundService(intent);
		} else {
			ctx.startService(intent);
		}
	}

	// stop shuts the service down and removes its notification.
	// Called from Go through JNI.
	public static void stop(Context ctx) {
		ctx.stopService(new Intent(ctx, GoForegroundService.class));
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
		if (intent == null) {
			// Re-delivered by the system without our extras, so there is
			// nothing to show. START_NOT_STICKY means this should not
			// happen, but a service that never calls startForeground would
			// crash the app.
			stopSelf();
			return START_NOT_STICKY;
		}

		startForeground(NOTIFICATION_ID, notification(intent.getStringExtra(EXTRA_TITLE),
				intent.getStringExtra(EXTRA_CONTENT)));
		return START_NOT_STICKY;
	}

	@Override
	public IBinder onBind(Intent intent) {
		return null; // not a bound service
	}

	private Notification notification(String title, String content) {
		Notification.Builder builder;
		if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
			builder = new Notification.Builder(this, CHANNEL_ID);
		} else {
			builder = new Notification.Builder(this);
		}

		int icon = getApplicationInfo().icon;
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
