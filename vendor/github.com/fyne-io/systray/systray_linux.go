package systray

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log"
	"os"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/introspect"
	"github.com/godbus/dbus/v5/prop"
)

var (
	// to signal quitting the internal main loop
	quitChan = make(chan struct{})

	// icon data for the main systray icon
	iconData []byte

	// the title of our system tray icon
	title string
)

// SetTemplateIcon sets the systray icon as a template icon (on macOS), falling back
// to a regular icon on other platforms.
// templateIconBytes and iconBytes should be the content of .ico for windows and
// .ico/.jpg/.png for other platforms.
func SetTemplateIcon(templateIconBytes []byte, regularIconBytes []byte) {
	// TODO handle the templateIconBytes?
	iconData = regularIconBytes
}


// SetIcon sets the systray icon.
// iconBytes should be the content of .ico for windows and .ico/.jpg/.png
// for other platforms.
func SetIcon(iconBytes []byte) {
	iconData = iconBytes
}

// SetTitle sets the systray title, only available on Mac and Linux.
func SetTitle(t string) {
	title = t
}

// SetTooltip sets the systray tooltip to display on mouse hover of the tray icon,
// only available on Mac and Windows.
func SetTooltip(tooltip string) {
}

// SetIcon sets the icon of a menu item. Only works on macOS and Windows.
// iconBytes should be the content of .ico/.jpg/.png
func (item *MenuItem) SetIcon(iconBytes []byte) {
}

// SetTemplateIcon sets the icon of a menu item as a template icon (on macOS). On Windows, it
// falls back to the regular icon bytes and on Linux it does nothing.
// templateIconBytes and regularIconBytes should be the content of .ico for windows and
// .ico/.jpg/.png for other platforms.
func (item *MenuItem) SetTemplateIcon(templateIconBytes []byte, regularIconBytes []byte) {
}

func addOrUpdateMenuItem(item *MenuItem) {
}

func setInternalLoop(_ bool) {
	// nothing to action on Linux
}

func registerSystray() {
}

func addSeparator(id uint32) {
}

func hideMenuItem(item *MenuItem) {
}

func showMenuItem(item *MenuItem) {
}

func nativeLoop() int {
	nativeStart()
	select {
	case <- quitChan:
		break
	}
	nativeEnd()
	return 0
}

func nativeEnd() {
	systrayExit()
}

func quit() {
	close(quitChan)
}

func nativeStart() {
	systrayReady()

	intro := introspect.Introspectable(`<node>
<!--  Based on:
 https://invent.kde.org/frameworks/knotifications/-/blob/master/src/org.kde.StatusNotifierItem.xml
 -->
<interface name="org.kde.StatusNotifierItem">
<property name="Category" type="s" access="read"/>
<property name="Id" type="s" access="read"/>
<property name="Title" type="s" access="read"/>
<property name="Status" type="s" access="read"/>
<property name="WindowId" type="i" access="read"/>
<!--  An additional path to add to the theme search path to find the icons specified above.  -->
<property name="IconThemePath" type="s" access="read"/>
<property name="Menu" type="o" access="read"/>
<property name="ItemIsMenu" type="b" access="read"/>
<!--  main icon  -->
<!--  names are preferred over pixmaps  -->
<property name="IconName" type="s" access="read"/>
<!-- struct containing width, height and image data -->
<property name="IconPixmap" type="a(iiay)" access="read">
<annotation name="org.qtproject.QtDBus.QtTypeName" value="KDbusImageVector"/>
</property>
<property name="OverlayIconName" type="s" access="read"/>
<property name="OverlayIconPixmap" type="a(iiay)" access="read">
<annotation name="org.qtproject.QtDBus.QtTypeName" value="KDbusImageVector"/>
</property>
<!--  Requesting attention icon  -->
<property name="AttentionIconName" type="s" access="read"/>
<!-- same definition as image -->
<property name="AttentionIconPixmap" type="a(iiay)" access="read">
<annotation name="org.qtproject.QtDBus.QtTypeName" value="KDbusImageVector"/>
</property>
<property name="AttentionMovieName" type="s" access="read"/>
<!--  tooltip data  -->
<!-- (iiay) is an image -->
<!--  We disable this as we don't support tooltip, so no need to go through it
    <property name="ToolTip" type="(sa(iiay)ss)" access="read">
        <annotation name="org.qtproject.QtDBus.QtTypeName" value="KDbusToolTipStruct"/>
    </property>
     -->
<!--  interaction: the systemtray wants the application to do something  -->
<method name="ContextMenu">
<!--  we're passing the coordinates of the icon, so the app knows where to put the popup window  -->
<arg name="x" type="i" direction="in"/>
<arg name="y" type="i" direction="in"/>
</method>
<method name="Activate">
<arg name="x" type="i" direction="in"/>
<arg name="y" type="i" direction="in"/>
</method>
<method name="SecondaryActivate">
<arg name="x" type="i" direction="in"/>
<arg name="y" type="i" direction="in"/>
</method>
<method name="Scroll">
<arg name="delta" type="i" direction="in"/>
<arg name="orientation" type="s" direction="in"/>
</method>
<!--  Signals: the client wants to change something in the status -->
<signal name="NewTitle"> </signal>
<signal name="NewIcon"> </signal>
<signal name="NewAttentionIcon"> </signal>
<signal name="NewOverlayIcon"> </signal>
<!--  We disable this as we don't support tooltip, so no need to go through it
    <signal name="NewToolTip">
    </signal>
     -->
<signal name="NewStatus">
<arg name="status" type="s"/>
</signal>
<!--  The following items are not supported by specs, but widely used  -->
<signal name="NewIconThemePath">
<arg type="s" name="icon_theme_path" direction="out"/>
</signal>
<signal name="NewMenu"/>
<!--  ayatana labels  -->
<!--  These are commented out because GDBusProxy would otherwise require them,
         but they are not available for KDE indicators
     -->
<!-- <signal name="XAyatanaNewLabel">
        <arg type="s" name="label" direction="out" />
        <arg type="s" name="guide" direction="out" />
    </signal>
    <property name="XAyatanaLabel" type="s" access="read" />
    <property name="XAyatanaLabelGuide" type="s" access="read" /> -->
</interface>`+ introspect.IntrospectDataString+prop.IntrospectDataString+"</node>")

	conn, _ := dbus.ConnectSessionBus()
	conn.Export(introspect.Introspectable(intro), "/StatusNotifierItem",
		"org.freedesktop.DBus.Introspectable")

	_, err := prop.Export(conn, "/StatusNotifierItem", createPropSpec())
	if err != nil {
		log.Printf("Failed to export notifier item properties to bus")
		return
	}
	err = conn.Export(&tray{}, "/StatusNotifierItem", "org.kde.StatusNotifierItem")
	if err != nil {
		log.Printf("Failed to export notifier item to bus")
		return
	}

	name := fmt.Sprintf("org.kde.StatusNotifierItem-%d-1", os.Getpid()) // register id 1 for this process
	_, err = conn.RequestName(name, dbus.NameFlagDoNotQueue)
	if err != nil {
		// fall back to existing name
		name = conn.Names()[0]
	}

	obj := conn.Object("org.kde.StatusNotifierWatcher", dbus.ObjectPath("/StatusNotifierWatcher"))
	call := obj.Call("org.kde.StatusNotifierWatcher.RegisterStatusNotifierItem", 0, name)
	if call.Err != nil {
		log.Printf("Failed to register our icon with the notifier watcher, maybe no tray running?")
	}
}

// tray is a basic type that handles the dbus functionality
type tray struct {
	prop.Properties
}

// ContextMenu method is called when the user has right-clicked on our icon.
func (t *tray) ContextMenu(x, y int) *dbus.Error {
	// not supported for systray lib
	return nil
}

// Activate requests that we perform the primary action, such as showing a menu.
func (t *tray) Activate(x, y int) *dbus.Error {
	// TODO show menu, or have it handled in the dbus?
	return nil
}

// SecondaryActivate is alternative non-context click, such as middle mouse button.
func (t *tray) SecondaryActivate(x, y int) *dbus.Error {
	return nil
}

// Scroll is called when the mouse wheel scrolls over the icon.
func (t *tray) Scroll(delta int, orient string) *dbus.Error {
	return nil
}

type PX struct {
	W, H int
	Pix  []byte
}

func convertToPixels(data []byte) PX {
	if len(iconData) == 0 {
		return PX{}
	}

	p, _, err := image.Decode(bytes.NewReader(iconData))
	if err != nil {
		log.Printf("Failed to read icon format %v", err)
		return PX{}
	}

	w, h := p.Bounds().Dx(), p.Bounds().Dy()
	data2 := make([]byte, w*h*4)
	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := p.At(x, y).RGBA()
			data2[i] = byte(a)
			data2[i+1] = byte(r)
			data2[i+2] = byte(g)
			data2[i+3] = byte(b)
			i += 4
		}
	}
	return PX{
		w, h,
		data2,
	}
}

func createPropSpec() map[string]map[string]*prop.Prop {
	return map[string]map[string]*prop.Prop{
		"org.kde.StatusNotifierItem": {
			"Status": {
				"Active", // Passive, Active or NeedsAttention
				false,
				prop.EmitTrue,
				nil,
			},
			"Title": {
				title,
				false,
				prop.EmitTrue,
				nil,
			},
			"Id": {
				"1",
				false,
				prop.EmitTrue,
				nil,
			},
			"Category": {
				"ApplicationStatus",
				false,
				prop.EmitTrue,
				nil,
			},
			"IconName": {
				"",
				false,
				prop.EmitTrue,
				nil,
			},
			"IconPixmap": {
				[]PX{convertToPixels(iconData)},
				false,
				prop.EmitTrue,
				nil,
			},
			"IconThemePath": {
				"",
				false,
				prop.EmitTrue,
				nil,
			},
		}}
}