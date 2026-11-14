package org.golang.app;

import android.app.Notification;
import android.app.NotificationChannel;
import android.app.NotificationManager;
import android.content.BroadcastReceiver;
import android.content.Context;
import android.content.Intent;
import android.os.Build;

// FyneNotificationReceiver is woken by AlarmManager at the scheduled time and
// posts the notification through NotificationManager. It must be declared in
// the application's AndroidManifest.xml for the alarm to fire after the app
// process has been killed:
//
//   <receiver android:name="org.golang.app.FyneNotificationReceiver"
//             android:exported="false" />
//
// The Fyne packaging tool is responsible for emitting that line.
public class FyneNotificationReceiver extends BroadcastReceiver {
    private static final String CHANNEL_ID = "fyne-notif";
    private static final String SILENT_CHANNEL_ID = "fyne-notif-silent";
    private static final int UNKNOWN_APP_ICON = 17629184; // android.R.drawable.sym_def_app_icon

    @Override
    public void onReceive(Context context, Intent intent) {
        String title = intent.getStringExtra("title");
        String body = intent.getStringExtra("body");
        int notifId = intent.getIntExtra("notif_id", 1);
        boolean silent = intent.getBooleanExtra("silent", false);

        NotificationManager mgr = (NotificationManager) context.getSystemService(Context.NOTIFICATION_SERVICE);
        if (mgr == null) {
            return;
        }

        Notification.Builder builder;
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            NotificationChannel channel;
            if (silent) {
                channel = new NotificationChannel(SILENT_CHANNEL_ID,
                        "Fyne Silent Notification", NotificationManager.IMPORTANCE_HIGH);
                channel.setSound(null, null);
            } else {
                channel = new NotificationChannel(CHANNEL_ID,
                        "Fyne Notification", NotificationManager.IMPORTANCE_HIGH);
            }
            mgr.createNotificationChannel(channel);
            builder = new Notification.Builder(context, channel.getId());
        } else {
            builder = new Notification.Builder(context);
            if (!silent) {
                builder.setDefaults(Notification.DEFAULT_SOUND);
            }
        }

        builder.setContentTitle(title != null ? title : "");
        builder.setContentText(body != null ? body : "");
        builder.setSmallIcon(UNKNOWN_APP_ICON);
        builder.setAutoCancel(true);

        mgr.notify(notifId, builder.build());
    }
}
