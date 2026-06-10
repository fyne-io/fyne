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
    private static final int UNKNOWN_APP_ICON = 17629184; // android.R.drawable.sym_def_app_icon

    @Override
    public void onReceive(Context context, Intent intent) {
        String title = intent.getStringExtra("title");
        String body = intent.getStringExtra("body");
        int notifId = intent.getIntExtra("notif_id", 1);

        NotificationManager mgr = (NotificationManager) context.getSystemService(Context.NOTIFICATION_SERVICE);
        if (mgr == null) {
            return;
        }

        Notification.Builder builder;
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            NotificationChannel channel = new NotificationChannel(CHANNEL_ID,
                    "Fyne Notification", NotificationManager.IMPORTANCE_HIGH);
            mgr.createNotificationChannel(channel);
            builder = new Notification.Builder(context, CHANNEL_ID);
        } else {
            builder = new Notification.Builder(context);
        }

        builder.setContentTitle(title != null ? title : "");
        builder.setContentText(body != null ? body : "");
        builder.setSmallIcon(UNKNOWN_APP_ICON);
        builder.setAutoCancel(true);

        mgr.notify(notifId, builder.build());
    }
}
