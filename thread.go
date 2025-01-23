package fyne

// DoAndWait is used to execute a specified function in the main Fyne runtime context.
// This is required when a background process wishes to adjust graphical elements of a running app.
// Developers should use this only from within goroutines they have created.
//
// Since: 2.6
func DoAndWait(fn func()) {
	CurrentApp().Driver().DoFromGoroutine(fn, true)
}

// Do is used to execute a specified function in the main Fyne runtime context without waiting.
// This is required when a background process wishes to adjust graphical elements of a running app.
// Developers should use this only from within goroutines they have created and when the result does not have to
// be waited for.
//
// Since: 2.6
func Do(fn func()) {
	CurrentApp().Driver().DoFromGoroutine(fn, false)
}
