//go:build linux || freebsd || openbsd || netbsd
// +build linux freebsd openbsd netbsd

package systray

import (
	"log"
	"sync/atomic"

	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"

	"fyne.io/systray/internal/generated/menu"
)

// SetIcon sets the icon of a menu item. Only works on macOS and Windows.
// iconBytes should be the content of .ico/.jpg/.png
func (item *MenuItem) SetIcon(iconBytes []byte) {
}

// copyLayout makes full copy of layout
func copyLayout(in *menuLayout, depth int32) *menuLayout {
	out := menuLayout{
		V0: in.V0,
		V1: make(map[string]dbus.Variant, len(in.V1)),
	}
	for k, v := range in.V1 {
		out.V1[k] = v
	}
	if depth != 0 {
		depth--
		out.V2 = make([]dbus.Variant, len(in.V2))
		for i, v := range in.V2 {
			out.V2[i] = dbus.MakeVariant(copyLayout(v.Value().(*menuLayout), depth))
		}
	} else {
		out.V2 = []dbus.Variant{}
	}
	return &out
}

func (t *tray) GetLayout(parentID int32, recursionDepth int32, propertyNames []string) (revision uint32, layout menuLayout, err *dbus.Error) {
	if m, ok := findLayout(parentID); ok {
		instance.menuLock.RLock()
		defer instance.menuLock.RUnlock()
		// return copy of menu layout to prevent panic from cuncurrent access to layout
		return atomic.LoadUint32(&instance.menuVersion), *copyLayout(m, recursionDepth), nil
	}
	return
}

// GetGroupProperties is com.canonical.dbusmenu.GetGroupProperties method.
func (t *tray) GetGroupProperties(ids []int32, propertyNames []string) (properties []struct {
	V0 int32
	V1 map[string]dbus.Variant
}, err *dbus.Error) {
	for _, id := range ids {
		if m, ok := findLayout(id); ok {
			p := struct {
				V0 int32
				V1 map[string]dbus.Variant
			}{
				V0: m.V0,
				V1: make(map[string]dbus.Variant, len(m.V1)),
			}
			instance.menuLock.RLock()
			for k, v := range m.V1 {
				p.V1[k] = v
			}
			instance.menuLock.RUnlock()
			properties = append(properties, p)
		}
	}
	return
}

// GetProperty is com.canonical.dbusmenu.GetProperty method.
func (t *tray) GetProperty(id int32, name string) (value dbus.Variant, err *dbus.Error) {
	if m, ok := findLayout(id); ok {
		instance.menuLock.RLock()
		defer instance.menuLock.RUnlock()
		if p, ok := m.V1[name]; ok {
			return p, nil
		}
	}
	return
}

// Event is com.canonical.dbusmenu.Event method.
func (t *tray) Event(id int32, eventID string, data dbus.Variant, timestamp uint32) (err *dbus.Error) {
	if eventID == "clicked" {
		menuItemsLock.RLock()
		item, ok := menuItems[uint32(id)]
		menuItemsLock.RUnlock()
		if !ok {
			log.Printf("systray error: failed to look up clicked menu item with ID %d\n", id)
			return
		}

		item.ClickedCh <- struct{}{}
	}
	return
}

// EventGroup is com.canonical.dbusmenu.EventGroup method.
func (t *tray) EventGroup(events []struct {
	V0 int32
	V1 string
	V2 dbus.Variant
	V3 uint32
}) (idErrors []int32, err *dbus.Error) {
	for _, event := range events {
		if event.V1 == "clicked" {
			menuItemsLock.RLock()
			item, ok := menuItems[uint32(event.V0)]
			menuItemsLock.RUnlock()
			if !ok {
				log.Printf("systray error: failed to look up clicked menu item with ID %d\n", event.V0)
				return
			}

			item.ClickedCh <- struct{}{}
		}
	}
	return
}

// AboutToShow is com.canonical.dbusmenu.AboutToShow method.
func (t *tray) AboutToShow(id int32) (needUpdate bool, err *dbus.Error) {
	return
}

// AboutToShowGroup is com.canonical.dbusmenu.AboutToShowGroup method.
func (t *tray) AboutToShowGroup(ids []int32) (updatesNeeded []int32, idErrors []int32, err *dbus.Error) {
	return
}

func createMenuPropSpec() map[string]map[string]*prop.Prop {
	return map[string]map[string]*prop.Prop{
		"com.canonical.dbusmenu": {
			"Version": {
				Value:    atomic.LoadUint32(&instance.menuVersion),
				Writable: true,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
			"TextDirection": {
				Value:    "ltr",
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
			"Status": {
				Value:    "normal",
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
			"IconThemePath": {
				Value:    []string{},
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: nil,
			},
		},
	}
}

// menuLayout is a named struct to map into generated bindings. It represents the layout of a menu item
type menuLayout = struct {
	V0 int32                   // the unique ID of this item
	V1 map[string]dbus.Variant // properties for this menu item layout
	V2 []dbus.Variant          // child menu item layouts
}

func addOrUpdateMenuItem(item *MenuItem) {
	var layout *menuLayout
	m, exists := findLayout(int32(item.id))

	if exists {
		layout = m
	} else {
		layout = &menuLayout{
			V0: int32(item.id),
			V1: map[string]dbus.Variant{},
			V2: []dbus.Variant{},
		}

		parent := instance.menu
		if item.parent != nil {
			m, ok := findLayout(int32(item.parent.id))
			if ok {
				parent = m
				instance.menuLock.Lock()
				parent.V1["children-display"] = dbus.MakeVariant("submenu")
				instance.menuLock.Unlock()
			}
		}
		instance.menuLock.Lock()
		parent.V2 = append(parent.V2, dbus.MakeVariant(layout))
		instance.menuLock.Unlock()
	}

	applyItemToLayout(item, layout)
	if exists {
		refresh()
	}
}

func addSeparator(id uint32) {
	layout := &menuLayout{
		V0: int32(id),
		V1: map[string]dbus.Variant{
			"type": dbus.MakeVariant("separator"),
		},
		V2: []dbus.Variant{},
	}
	instance.menuLock.Lock()
	instance.menu.V2 = append(instance.menu.V2, dbus.MakeVariant(layout))
	instance.menuLock.Unlock()
}

func applyItemToLayout(in *MenuItem, out *menuLayout) {
	instance.menuLock.Lock()
	defer instance.menuLock.Unlock()
	out.V1["type"] = dbus.MakeVariant("standard")
	out.V1["enabled"] = dbus.MakeVariant(!in.disabled)
	out.V1["label"] = dbus.MakeVariant(in.title)

	if in.isCheckable {
		out.V1["toggle-type"] = dbus.MakeVariant("checkmark")
		if in.checked {
			out.V1["toggle-state"] = dbus.MakeVariant(1)
		} else {
			out.V1["toggle-state"] = dbus.MakeVariant(0)
		}
	} else {
		out.V1["toggle-type"] = dbus.MakeVariant("")
		out.V1["toggle-state"] = dbus.MakeVariant(0)
	}
}

func findLayout(id int32) (*menuLayout, bool) {
	if id == 0 {
		return instance.menu, true
	}
	instance.menuLock.RLock()
	defer instance.menuLock.RUnlock()
	return findSubLayout(id, instance.menu.V2)
}

func findSubLayout(id int32, vals []dbus.Variant) (*menuLayout, bool) {
	for _, i := range vals {
		item := i.Value().(*menuLayout)
		if item.V0 == id {
			return item, true
		}

		if len(item.V2) > 0 {
			child, ok := findSubLayout(id, item.V2)
			if ok {
				return child, true
			}
		}
	}

	return nil, false
}

func hideMenuItem(item *MenuItem) {
	m, exists := findLayout(int32(item.id))
	if exists {
		instance.menuLock.Lock()
		m.V1["visible"] = dbus.MakeVariant(false)
		instance.menuLock.Unlock()
		refresh()
	}
}

func showMenuItem(item *MenuItem) {
	m, exists := findLayout(int32(item.id))
	if exists {
		instance.menuLock.Lock()
		m.V1["visible"] = dbus.MakeVariant(true)
		instance.menuLock.Unlock()
		refresh()
	}
}

func refresh() {
	instance.lock.RLock()
	conn := instance.conn
	menuProps := instance.menuProps
	instance.lock.RUnlock()
	if conn != nil && menuProps != nil {
		version := atomic.AddUint32(&instance.menuVersion, 1)
		dbusErr := menuProps.Set("com.canonical.dbusmenu", "Version",
			dbus.MakeVariant(version))
		if dbusErr != nil {
			log.Printf("systray error: failed to update menu version: %s\n", dbusErr)
			return
		}

		err := menu.Emit(conn, &menu.Dbusmenu_LayoutUpdatedSignal{
			Path: menuPath,
			Body: &menu.Dbusmenu_LayoutUpdatedSignalBody{
				Revision: version,
			},
		})
		if err != nil {
			log.Printf("systray error: failed to emit layout updated signal: %s\n", err)
		}
	}
}

func resetMenu() {
	instance.menu = &menuLayout{}
	instance.menuVersion++
	refresh()
}
