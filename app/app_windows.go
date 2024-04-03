//go:build !ci && !js && !android && !ios && !wasm && !test_web_driver
// +build !ci,!js,!android,!ios,!wasm,!test_web_driver

package app

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"golang.org/x/sys/execabs"
	"golang.org/x/sys/windows/registry"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

var once sync.Once

const notificationTemplate = `$title = "%s"
$content = "%s"
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

func isDark() bool {
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Themes\Personalize`, registry.QUERY_VALUE)
	if err != nil { // older version of Windows will not have this key
		return false
	}
	defer k.Close()

	useLight, _, err := k.GetIntegerValue("AppsUseLightTheme")
	if err != nil { // older version of Windows will not have this value
		return false
	}

	return useLight == 0
}

func defaultVariant() fyne.ThemeVariant {
	if isDark() {
		return theme.VariantDark
	}
	return theme.VariantLight
}

func rootConfigDir() string {
	homeDir, _ := os.UserHomeDir()

	desktopConfig := filepath.Join(filepath.Join(homeDir, "AppData"), "Roaming")
	return filepath.Join(desktopConfig, "fyne")
}

func (a *fyneApp) OpenURL(url *url.URL) error {
	cmd := execabs.Command("rundll32", "url.dll,FileProtocolHandler", url.String())
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	return cmd.Run()
}

var scriptNum = 0

func (a *fyneApp) SendNotification(n *fyne.Notification) {
	title := escapeNotificationString(n.Title)
	content := escapeNotificationString(n.Content)
	appID := a.UniqueID()
	iconFilePath := a.cachedNotificationIconPath()
	if appID == "" || strings.Index(appID, "missing-id") == 0 {
		appID = a.Metadata().Name
	}

	script := fmt.Sprintf(notificationTemplate, title, content, iconFilePath, appID)
	go runScript("notify", script)
}

func (a *fyneApp) saveNotificationIconToCache(dirPath, filePath string) error {
	err := os.MkdirAll(dirPath, 0700)
	if err != nil {
		fyne.LogError("Unable to create application cache directory", err)
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		fyne.LogError("Unable to create notification icon file", err)
		return err
	}

	defer file.Close()

	if icon := a.Icon(); icon != nil {
		_, err = file.Write(icon.Content())
		if err != nil {
			fyne.LogError("Unable to write notification icon contents", err)
			return err
		}
	}

	return nil
}

func (a *fyneApp) cachedNotificationIconPath() string {
	if a.Icon() == nil {
		return ""
	}

	dirPath := filepath.Join(rootConfigDir(), a.UniqueID())
	filePath := filepath.Join(dirPath, "icon.png")
	once.Do(func() {
		err := a.saveNotificationIconToCache(dirPath, filePath)
		if err != nil {
			filePath = ""
		}
	})

	return filePath
}

// SetSystemTrayMenu creates a system tray item and attaches the specified menu.
// By default this will use the application icon.
func (a *fyneApp) SetSystemTrayMenu(menu *fyne.Menu) {
	a.Driver().(systrayDriver).SetSystemTrayMenu(menu)
}

// SetSystemTrayIcon sets a custom image for the system tray icon.
// You should have previously called `SetSystemTrayMenu` to initialise the menu icon.
func (a *fyneApp) SetSystemTrayIcon(icon fyne.Resource) {
	a.Driver().(systrayDriver).SetSystemTrayIcon(icon)
}

func escapeNotificationString(in string) string {
	noSlash := strings.ReplaceAll(in, "`", "``")
	return strings.ReplaceAll(noSlash, "\"", "`\"")
}

func runScript(name, script string) {
	scriptNum++
	appID := fyne.CurrentApp().UniqueID()
	fileName := fmt.Sprintf("fyne-%s-%s-%d.ps1", appID, name, scriptNum)

	tmpFilePath := filepath.Join(os.TempDir(), fileName)
	err := os.WriteFile(tmpFilePath, []byte(script), 0600)
	if err != nil {
		fyne.LogError("Could not write script to show notification", err)
		return
	}
	defer os.Remove(tmpFilePath)

	launch := "(Get-Content -Encoding UTF8 -Path " + tmpFilePath + " -Raw) | Invoke-Expression"
	cmd := execabs.Command("PowerShell", "-ExecutionPolicy", "Bypass", launch)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err = cmd.Run()
	if err != nil {
		fyne.LogError("Failed to launch windows notify script", err)
	}
}
func watchTheme() {
	// TODO monitor the Windows theme
}
