// +build linux freebsd openbsd netbsd

//Note that you need to have github.com/knightpp/dbus-codegen-go installed from "custom" branch
//go:generate dbus-codegen-go -prefix org.kde -package notifier -output internal/generated/notifier/status_notifier_item.go internal/StatusNotifierItem.xml
//go:generate dbus-codegen-go -prefix com.canonical -package menu -output internal/generated/menu/dbus_menu.go internal/DbusMenu.xml

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

	"fyne.io/systray/internal/generated/menu"
	"fyne.io/systray/internal/generated/notifier"
)

const (
	path     = "/StatusNotifierItem"
	menuPath = "/StatusNotifierMenu"
)

var (
	// to signal quitting the internal main loop
	quitChan = make(chan struct{})

	// instance is the current instance of our DBus tray server
	instance = &tray{menu: &menuLayout{}}
)

// SetTemplateIcon sets the systray icon as a template icon (on macOS), falling back
// to a regular icon on other platforms.
// templateIconBytes and iconBytes should be the content of .ico for windows and
// .ico/.jpg/.png for other platforms.
func SetTemplateIcon(templateIconBytes []byte, regularIconBytes []byte) {
	// TODO handle the templateIconBytes?
	SetIcon(regularIconBytes)
}

// SetIcon sets the systray icon.
// iconBytes should be the content of .ico for windows and .ico/.jpg/.png
// for other platforms.
func SetIcon(iconBytes []byte) {
	instance.iconData = iconBytes

	if instance.props != nil {
		instance.props["org.kde.StatusNotifierItem"]["IconPixmap"].Value = []PX{convertToPixels(iconBytes)}

		if instance.conn != nil {
			notifier.Emit(instance.conn, &notifier.StatusNotifierItem_NewIconSignal{
				Path: path,
				Body: &notifier.StatusNotifierItem_NewIconSignalBody{},
			})
		}
	}
}

// SetTitle sets the systray title, only available on Mac and Linux.
func SetTitle(t string) {
	instance.title = t

	if instance.props != nil {
		instance.props["org.kde.StatusNotifierItem"]["Title"].Value = t

		if instance.conn != nil {
			notifier.Emit(instance.conn, &notifier.StatusNotifierItem_NewTitleSignal{
				Path: path,
				Body: &notifier.StatusNotifierItem_NewTitleSignalBody{},
			})
		}
	}
}

// SetTooltip sets the systray tooltip to display on mouse hover of the tray icon,
// only available on Mac and Windows.
func SetTooltip(tooltip string) {
	instance.tooltipTitle = tooltip
}

// SetTemplateIcon sets the icon of a menu item as a template icon (on macOS). On Windows, it
// falls back to the regular icon bytes and on Linux it does nothing.
// templateIconBytes and regularIconBytes should be the content of .ico for windows and
// .ico/.jpg/.png for other platforms.
func (item *MenuItem) SetTemplateIcon(templateIconBytes []byte, regularIconBytes []byte) {
}

func setInternalLoop(_ bool) {
	// nothing to action on Linux
}

func registerSystray() {
}

func nativeLoop() int {
	nativeStart()
	select {
	case <-quitChan:
		break
	}
	nativeEnd()
	return 0
}

func nativeEnd() {
	systrayExit()
	instance.conn.Close()
}

func quit() {
	close(quitChan)
}

func nativeStart() {
	systrayReady()

	conn, _ := dbus.ConnectSessionBus()
	instance.conn = conn
	notifier.ExportStatusNotifierItem(conn, path, &notifier.UnimplementedStatusNotifierItem{})
	menu.ExportDbusmenu(conn, menuPath, instance)

	name := fmt.Sprintf("org.kde.StatusNotifierItem-%d-1", os.Getpid()) // register id 1 for this process
	_, err := conn.RequestName(name, dbus.NameFlagDoNotQueue)
	if err != nil {
		// fall back to existing name
		name = conn.Names()[0]
	}

	_, err = prop.Export(conn, path, instance.createPropSpec())
	if err != nil {
		log.Printf("Failed to export notifier item properties to bus")
		return
	}
	_, err = prop.Export(conn, menuPath, createMenuPropSpec())
	if err != nil {
		log.Printf("Failed to export notifier menu properties to bus")
		return
	}

	node := introspect.Node{
		Name: path,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			notifier.IntrospectDataStatusNotifierItem,
		},
	}
	err = conn.Export(introspect.NewIntrospectable(&node), path,
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		log.Printf("Failed to export introspection")
		return
	}
	menuNode := introspect.Node{
		Name: menuPath,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			menu.IntrospectDataDbusmenu,
		},
	}
	err = conn.Export(introspect.NewIntrospectable(&menuNode), menuPath,
		"org.freedesktop.DBus.Introspectable")
	if err != nil {
		log.Printf("Failed to export introspection")
		return
	}

	obj := conn.Object("org.kde.StatusNotifierWatcher", "/StatusNotifierWatcher")
	call := obj.Call("org.kde.StatusNotifierWatcher.RegisterStatusNotifierItem", 0, path)
	if call.Err != nil {
		log.Printf("Failed to register our icon with the notifier watcher, maybe no tray running?")
	}
}

// tray is a basic type that handles the dbus functionality
type tray struct {
	// the DBus connection that we will use
	conn *dbus.Conn

	// icon data for the main systray icon
	iconData []byte
	// title and tooltip state
	title, tooltipTitle string

	menu  *menuLayout
	props map[string]map[string]*prop.Prop
}

func (t *tray) createPropSpec() map[string]map[string]*prop.Prop {
	t.props = map[string]map[string]*prop.Prop{
		"org.kde.StatusNotifierItem": {
			"Status": {
				"Active", // Passive, Active or NeedsAttention
				false,
				prop.EmitTrue,
				nil,
			},
			"Title": {
				t.title,
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
				[]PX{convertToPixels(t.iconData)},
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
			"ItemIsMenu": {
				true,
				false,
				prop.EmitTrue,
				nil,
			},
			"Menu": {
				dbus.ObjectPath(menuPath),
				false,
				prop.EmitTrue,
				nil,
			},
			"ToolTip": {
				tooltip{V2: t.tooltipTitle},
				false,
				prop.EmitTrue,
				nil,
			},
		}}
	return t.props
}

type PX struct {
	W, H int
	Pix  []byte
}

// tooltip is our data for a tooltip property.
// Param names need to match the generated code...
type tooltip = struct {
	V0 string // name
	V1 []PX   // icons
	V2 string // title
	V3 string // description
}

func convertToPixels(data []byte) PX {
	if len(data) == 0 {
		return PX{}
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Printf("Failed to read icon format %v", err)
		return PX{}
	}

	return PX{
		img.Bounds().Dx(), img.Bounds().Dy(),
		argbForImage(img),
	}
}

func argbForImage(img image.Image) []byte {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	data := make([]byte, w*h*4)
	i := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			data[i] = byte(a)
			data[i+1] = byte(r)
			data[i+2] = byte(g)
			data[i+3] = byte(b)
			i += 4
		}
	}
	return data
}
