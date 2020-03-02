# Changelog

This file lists the main changes with each version of the Fyne toolkit.
More detailed release notes can be found on the [releases page](https://github.com/fyne-io/fyne/releases). 

## 1.2.3 - 2 March 2020

### Added

 * Add media and volume icons to default themes (#649)
 * Add Canvas.PixelCoordinateForPosition to find pixel locations if required
 * Add ProgressInfinite dialog

### Changed

 * Warn if -executable or -sourceDir flags are used for package on mobile (#652)
 * Update scale based on device for mobile apps
 * Windows without a title will now be named "Fyne Application"
 * Revert fix to quit mobile apps - this is not allowed in guidelines

### Fixed

 * App.UniqueID() did not return current app ID
 * Fyne package ignored -name flag for ios and android builds (#657)
 * Possible crash when appending tabs to TabContainer
 * FixedSize windows not rescaling when dragged between monitors (#654)
 * Fix issues where older Android devices may not background or rotate (#677)
 * Crash when setting theme before window content set (#688)
 * Correct form extend behaviour (#694)
 * Select drop-down width is wrong if the drop-down is too tall for the window (#706)


## 1.2.2 - 29 January 2020

### Added

* Add SelectedText() function to Entry widget
* New mobile.Device interface exposing ShowVirtualKeyboard() (and Hide...)

### Changed

* Scale calculations are now relative to system scale - the default "1" matches the system
* Update scale on Linux to be "auto" by default (and numbers are relative to 96DPI standard) (#595)
* When auto scaling check the monitor in the middle of the window, not top left
* bundled files now have a standard header to optimise some tools like go report card
* Shortcuts are now handled by the event queue - fixed possible deadlock

### Fixed

* Scroll horizontally when holding shift key (#579)
* Updating text and calling refresh for widget doesn't work (#607)
* Corrected visual behaviour of extended widgets including Entry, Select, Check, Radio and Icon (#615)
* Entries and Selects that are extended would crash on right click.
* PasswordEntry created from Entry with Password = true has no revealer
* Dialog width not always sufficient for title
* Pasting unicode characters could panic (#597)
* Setting theme before application start panics on macOS (#626)
* MenuItem type conflicts with other projects (#632)


## 1.2.1 - 24 December 2019

### Added

* Add TouchDown, TouchUp and TouchCancel API in driver/mobile for device specific events
* Add support for adding and removing tabs from a tab container (#444)

### Fixed

* Issues when settings changes may not be monitored (#576)
* Layout of hidden tab container contents on mobile (#578)
* Mobile apps would not quit when Quit() was called (#580)
* Shadows disappeared when theme changes (#589)
* iOS apps could stop rendering after many refreshes (#584)
* Fyne package could fail on Windows (#586)
* Horizontal only scroll container may not refresh using scroll wheel


## 1.2 - 12 December 2019

### Added

* Mobile support - iOS and Android, including "fyne package" command
* Support for OpenGL ES and embedded linux
* New BaseWidget for building custom widgets
* Support for diagonal gradients
* Global settings are now saved and can be set using the new fyne_settings app
* Support rendering in Go playground using playground.Render() helpers
* "fyne install" command to package and install apps on the local computer
* Add horizontal scrolling to ScrollContainer
* Add preferences API
* Add show/hide password icon when created from NewPasswordEntry
* Add NewGridLayoutWithRows to specify a grid layout with a set number of rows
* Add NewAdaptiveGridLayout which uses a column grid layout when horizontal and rows in vertical


### Changed

* New Logo! Thanks to Storm for his work on this :)
* Applications no longer have a default (Fyne logo) icon
* Input events now execute one at a time to maintain the correct order
* Button and other widget callbacks no longer launch new goroutines
* FYNE_THEME and FYNE_SCALE are now overrides to the global configuration
* The first opened window no longer exits the app when closed (unless none others are open or Window.SetMaster() is called)
* "fyne package" now defaults icon to "Icon.png" so the parameter is optional
* Calling ExtendBaseWidget() sets up the renderer for extended widgets
* Entry widget now has a visible Disabled state, ReadOnly has been deprecated
* Bundled images optimised to save space
* Optimise rendering to reduce refresh on TabContainer and ScrollContainer


### Fixed

* Correct the colour of Entry widget cursor if theme changes
* Error where widgets created before main() function could crash (#490)
* App.Run panics if called without a window (#527)
* Support context menu for disabled entry widgets (#488)
* Fix issue where images using fyne.ImageFillOriginal may not show initially (#558)


## 1.1.2 - 12 October 2019

### Added

### Changed

* Default scale value for canvases is now 1.0 instead of Auto (DPI based)

### Fixed

* Correct icon name in linux packages
* Fullscreen before showing a window works again
* Incorrect MinSize of FixedGrid layout in some situations
* Update text size on theme change
* Text handling crashes (#411, #484, #485)
* Layout of image only buttons
* TabItem.Content changes are reflected when refreshing TabContainer (#456)

## 1.1.1 - 17 August 2019

### Added

* Add support for custom Windows manifest files in fyne package

### Changed

* Dismiss non-modal popovers on secondary tap
* Only measure visible objects in layouts and minSize calculations (#343)
* Don't propagate show/hide in the model - allowing children of tabs to remain hidden
* Disable cut/copy for password fields
* Correctly calculate grid layout minsize as width changes
* Select text at end of line when double tapping beyond width

### Fixed

* Scale could be too large on macOS Retina screens
* Window with fixed size changes size when un-minimized on Windows (#300)
* Setting text on a label could crash if it was not yet shown (#381)
* Multiple Entry widgets could have selections simultaneously (#341)
* Hover effect of radio widget too low (#383)
* Missing shadow on Select widget
* Incorrect rendering of subimages within Image object
* Size calculation caches could be skipped causing degraded performance


## 1.1 - 1 July 2019

### Added

* Menubar and PopUpMenu (#41)
* PopUp widgets (regular and modal) and canvas overlay support (#242)
* Add gradient (linear and radial) to canvas
* Add shadow support for overlays, buttons and scrollcontainer
* Text can now be selected (#67)
* Support moving through inputs with Tab / Shift-Tab (#82)
* canvas.Capture() to save the content of a canvas
* Horizontal layout for widget.Radio
* Select widget (#21)
* Add support for disabling widgets (#234)
* Support for changing icon color (#246)
* Button hover effect
* Pointer drag event to main API
* support for desktop mouse move events
* Add a new "hints" build tag that can suggest UI improvements

### Changed

* TabContainer tab location can now be set with SetTabLocation()
* Dialog windows now appear as modal popups within a window
* Don't add a button bar to a form if it has no buttons
* Moved driver/gl package to internal/driver/gl
* Clicking/Tapping in an entry will position the cursor
* A container with no layout will not change the position or size of it's content
* Update the fyne_demo app to reflect the expanding feature set

### Fixed

* Allow scrollbars to be dragged (#133)
* Unicode char input with Option key on macOS (#247)
* Resizng fixed size windows (#248)
* Fixed various bugs in window sizing and padding
* Button icons do not center align if label is empty (#284)


## 1.0.1 - 20 April 2019

### Added

* Support for go modules
* Transparent backgrounds for widgets
* Entry.OnCursorChanged()
* Radio.Append() and Radio.SetSelected() (#229)

### Changed

* Clicking outside a focused element will unfocus it
* Handle key repeat for non-runes (#165)

### Fixed

* Remove duplicate options from a Radio widget (#230)
* Issue where paste shortcut is not called for Ctrl-V keyboard combination
* Cursor position when clearing text in Entry (#214)
* Antialias of lines and circles (fyne-io/examples#14)
* Crash on centering of windows (#220)
* Possible crash when closing secondary windows
* Possible crash when showing dialog
* Initial visibility of scroll bar in ScrollContainer
* Setting window icon when different from app icon.
* Possible panic on app.Quit() (#175)
* Various caches and race condition issues (#194, #217, #209).


## 1.0 - 19 March 2019

The first major release of the Fyne toolkit delivers a stable release of the
main functionality required to build basic GUI applications across multiple
platforms.

### Features

* Canvas API (rect, line, circle, text, image)
* Widget API (box, button, check, entry, form, group, hyperlink, icon, label, progress bar, radio, scroller, tabs and toolbar)
* Light and dark themes
* Pointer, key and shortcut APIs (generic and desktop extension)
* OpenGL driver for Linux, macOS and Windows
* Tools for embedding data and packaging releases

