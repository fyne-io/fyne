//go:build !ci && !android && !ios && !wasm && !test_web_driver && !tinygo

package app

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	internalapp "fyne.io/fyne/v2/internal/app"
	"fyne.io/fyne/v2/internal/scheduler"
)

const notificationTemplate = `$title = %q
$content = %q
$iconPath = "file:///%s"
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastImageAndText02)
$toastXml = [xml] $template.GetXml()
$toastXml.GetElementsByTagName("text")[0].AppendChild($toastXml.CreateTextNode($title)) > $null
$toastXml.GetElementsByTagName("text")[1].AppendChild($toastXml.CreateTextNode($content)) > $null
$toastXml.GetElementsByTagName("image")[0].SetAttribute("src", $iconPath) > $null
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($toastXml.OuterXml)
$toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("%s").Show($toast);`

const scheduledNotificationTemplate = `$title = %q
$content = %q
$iconPath = "file:///%s"
$id = %q
$delivery = [DateTimeOffset]::Parse(%q)
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$template = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastImageAndText02)
$toastXml = [xml] $template.GetXml()
$toastXml.GetElementsByTagName("text")[0].AppendChild($toastXml.CreateTextNode($title)) > $null
$toastXml.GetElementsByTagName("text")[1].AppendChild($toastXml.CreateTextNode($content)) > $null
$toastXml.GetElementsByTagName("image")[0].SetAttribute("src", $iconPath) > $null
$xml = New-Object Windows.Data.Xml.Dom.XmlDocument
$xml.LoadXml($toastXml.OuterXml)
$scheduled = [Windows.UI.Notifications.ScheduledToastNotification]::new($xml, $delivery)
$scheduled.Tag = $id
[Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("%s").AddToSchedule($scheduled);`

const cancelScheduledNotificationTemplate = `$id = "%s"
[Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime] > $null
$notifier = [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier("%s")
foreach ($s in $notifier.GetScheduledToastNotifications()) {
    if ($s.Tag -eq $id) { $notifier.RemoveFromSchedule($s) }
}`

func (a *fyneApp) OpenURL(url *url.URL) error {
	cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

var scriptNum = 0

func (a *fyneApp) SendNotification(n *fyne.Notification) {
	title := n.Title
	content := n.Content
	iconFilePath := a.cachedIconPath()
	appID := a.notificationAppID()

	script := fmt.Sprintf(notificationTemplate, title, content, iconFilePath, appID)
	go runScript("notify", script)
}

func (a *fyneApp) ScheduleNotification(n *fyne.Notification, when time.Time) (*fyne.ScheduledNotification, error) {
	if !when.After(time.Now()) {
		return nil, errors.New("scheduled delivery time must be in the future")
	}

	id, err := scheduler.NewID()
	if err != nil {
		return nil, err
	}

	title := n.Title
	content := n.Content
	iconFilePath := a.cachedIconPath()
	delivery := when.UTC().Format(time.RFC3339)
	appID := a.notificationAppID()

	script := fmt.Sprintf(scheduledNotificationTemplate, title, content, iconFilePath, id, delivery, appID)
	go runScript("schedule", script)
	return fyne.NewScheduledNotification(id, n, when), nil
}

func (a *fyneApp) CancelScheduledNotification(id string) error {
	appID := a.notificationAppID()
	script := fmt.Sprintf(cancelScheduledNotificationTemplate, id, appID)
	go runScript("cancel", script)
	return nil
}

func (a *fyneApp) notificationAppID() string {
	appID := a.UniqueID()
	if appID == "" || strings.Index(appID, "missing-id") == 0 {
		appID = a.Metadata().Name
	}
	return appID
}

// SetSystemTrayMenu creates a system tray item and attaches the specified menu.
// By default, this will use the application icon.
func (a *fyneApp) SetSystemTrayMenu(menu *fyne.Menu) {
	a.Driver().(systrayDriver).SetSystemTrayMenu(menu)
}

// SetSystemTrayIcon sets a custom image for the system tray icon.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayIcon(icon fyne.Resource) {
	a.Driver().(systrayDriver).SetSystemTrayIcon(icon)
}

// SetSystemTrayWindow assigns a window to be shown with the system tray menu is tapped.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayWindow(w fyne.Window) {
	a.Driver().(systrayDriver).SetSystemTrayWindow(w)
}

func runScript(name, script string) {
	scriptNum++
	appID := fyne.CurrentApp().UniqueID()
	fileName := fmt.Sprintf("fyne-%s-%s-%d.ps1", appID, name, scriptNum)

	tmpFilePath := filepath.Join(os.TempDir(), fileName)
	err := os.WriteFile(tmpFilePath, []byte(script), 0o600)
	if err != nil {
		fyne.LogError("Could not write script to show notification", err)
		return
	}
	defer os.Remove(tmpFilePath)

	launch := "(Get-Content -Encoding UTF8 -Path " + tmpFilePath + " -Raw) | Invoke-Expression"
	cmd := exec.Command("PowerShell", "-ExecutionPolicy", "Bypass", launch)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = cmd.Run()
	if err != nil {
		fyne.LogError("Failed to launch windows notify script", err)
	}
}

func watchTheme(s *settings) {
	go internalapp.WatchTheme(func() {
		fyne.Do(s.setupTheme)
	})
}

func (a *fyneApp) registerRepositories() {
	// no-op
}
