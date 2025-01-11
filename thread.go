package fyne

// Do is used to execute a specified function in the main Fyne runtime context.
// This is required when a background processes wishes to adjust graphical elements of a running app.
// Developers should use this only from within goroutines they have created.
//
// Since: 2.6
func Do(fn func()) {
	CurrentApp().Driver().DoFromGoroutine(fn)
}
