package theme

import (
	_ "fyne.io/fyne/test"
)

// TODO: add tests back in

//func TestIconThemeChangeName(t *testing.T) {
//	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
//	cancel := CancelIcon()
//	name := cancel.Name()
//
//	fyne.CurrentApp().Settings().SetTheme(LightTheme())
//	assert.NotEqual(t, name, cancel.Name())
//}

//func TestIconThemeChangeContent(t *testing.T) {
//	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
//	cancel := ConfirmIcon()
//	content := cancel.Content()
//
//	fyne.CurrentApp().Settings().SetTheme(LightTheme())
//	assert.NotEqual(t, content, cancel.Content())
//}

//func TestNewThemedResource(t *testing.T) {
//	fyne.CurrentApp().Settings().SetTheme(DarkTheme())
//	custom := NewThemedResource(mailcomposeDark, mailcomposeLight)
//	content := custom.Content()
//
//	fyne.CurrentApp().Settings().SetTheme(LightTheme())
//	assert.NotEqual(t, content, custom.Content())
//}
