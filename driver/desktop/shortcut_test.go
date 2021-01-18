package desktop

import (
	"testing"

	"fyne.io/fyne/v2"
)

func TestCustomShortcut_Shortcut(t *testing.T) {
	type fields struct {
		KeyName  fyne.KeyName
		Modifier Modifier
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
				Modifier: ControlModifier,
			},
			want: "CustomDesktop:Control+C",
		},
		{
			name: "Ctrl+Alt+Esc",
			fields: fields{
				KeyName:  fyne.KeyEscape,
				Modifier: ControlModifier + AltModifier,
			},
			want: "CustomDesktop:Control+Alt+Escape",
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
		mods Modifier
		want string
	}{
		{
			name: "None",
			mods: 0,
			want: "",
		},
		{
			name: "Ctrl",
			mods: ControlModifier,
			want: "Control",
		},
		{
			name: "Shift+Ctrl",
			mods: ShiftModifier + ControlModifier,
			want: "Shift+Control",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := modifierToString(tt.mods); got != tt.want {
				t.Errorf("modifierToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
