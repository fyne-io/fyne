package desktop

import (
	"strings"
	"testing"

	"fyne.io/fyne/v2"
)

func TestCustomShortcut_Shortcut(t *testing.T) {
	type fields struct {
		KeyName  fyne.KeyName
		Modifier fyne.KeyModifier
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Ctrl+C",
			fields: fields{
				KeyName:  fyne.KeyC,
				Modifier: fyne.KeyModifierControl,
			},
			want: "CustomDesktop:Control+C",
		},
		{
			name: "Ctrl+Alt+Esc",
			fields: fields{
				KeyName:  fyne.KeyEscape,
				Modifier: fyne.KeyModifierControl + fyne.KeyModifierAlt,
			},
			want: "CustomDesktop:Control+Alt+Escape",
		},
		{
			name: "Esc",
			fields: fields{
				KeyName: fyne.KeyEscape,
			},
			want: "CustomDesktop:Escape",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := &CustomShortcut{
				KeyName:  tt.fields.KeyName,
				Modifier: tt.fields.Modifier,
			}
			if got := cs.ShortcutName(); got != tt.want {
				t.Errorf("CustomShortcut.ShortcutName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_modifierToString(t *testing.T) {
	tests := []struct {
		name string
		mods fyne.KeyModifier
		want string
	}{
		{
			name: "None",
			mods: 0,
			want: "",
		},
		{
			name: "Ctrl",
			mods: fyne.KeyModifierControl,
			want: "Control+",
		},
		{
			name: "Shift+Ctrl",
			mods: fyne.KeyModifierShift + fyne.KeyModifierControl,
			want: "Shift+Control+",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &strings.Builder{}
			writeModifiers(w, tt.mods)
			if got := w.String(); got != tt.want {
				t.Errorf("modifierToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
