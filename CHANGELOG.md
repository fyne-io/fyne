# Changelog

This file lists the main changes with each version of the Fyne toolkit.
More detailed release notes can be found on the [releases page](https://github.com/fyne-io/fyne/releases). 

## 1.4.2 - 9 December 2020

### Added

* [fyne-cli] Add support for passing custom build tags (#1538)

### Changed

* Run validation on content change instead of on each Refresh in widget.Entry

### Fixed

* [fyne-cli] Android: allow to specify an inline password for the keystore
* Fixed Card widget MinSize (#1581)
* Fix missing release tag to enable BuildRelease in Settings.BuildType()
* Dialog shadow does not resize after Refresh (#1370)
* Android Duplicate Number Entry (#1256)
* Support older macOS by default - back to 10.11 (#886)
* Complete certification of macOS App Store releases (#1443)
* Fix compilation errors for early stage Wayland testing
* Fix entry.SetValidationError() not working correctly


## 1.4.1 - 20 November 2020

### Changed

* Table columns can now be different sizes using SetColumnWidth
* Avoid unnecessary validation check on Refresh in widget.Form

### Fixed

* Tree could flicker on mouse hover (#1488)
* Content of table cells could overflow when sized correctly
* file:// based URI on Android would fail to list folder (#1495)
* Images in iOS release were not all correct size (#1498)
* iOS compile failed with Go 1.15 (#1497)
* Possible crash when minimising app containing List on Windows
* File chooser dialog ignores drive Z (#1513)
* Entry copy/paste is crashing on android 7.1 (#1511)
* Fyne package creating invalid windows packages (#1521)
* Menu bar initially doesn't respond to mouse input on macOS (#505) 
* iOS: Missing CFBundleIconName and asset catalog (#1504)
* CenterOnScreen causes crash on MacOS when called from goroutine (#1539)
* desktop.MouseHover Button state is not reliable (#1533)
* Initial validation status in widget.Form is not respected
* Fix nil reference in disabled buttons (#1558)


## 1.4 - 1 November 2020

### Added (highlights)

* List (#156), Table (#157) and Tree collection Widgets
* Card, FileItem, Separator widgets
* ColorPicker dialog
* User selection of primary colour
* Container API package to ease using layouts and container widgets
* Add input validation
* ListableURI for working with directories etc
* Added PaddedLayout

* Window.SetCloseIntercept (#467)
* Canvas.InteractiveArea() to indicate where widgets should avoid
* TextFormatter for ProgressBar
* FileDialog.SetLocation() (#821)
* Added dialog.ShowFolderOpen (#941)
* Support to install on iOS and android with 'fyne install'
* Support asset bundling with go:generate
* Add fyne release command for preparing signed apps
* Add keyboard and focus support to Radio and Select widgets 

### Changed

* Theme update - new blue highlight, move buttons to outline
* Android SDK target updated to 29
* Mobile log entries now start "Fyne" instead of "GoLog"
* Don't expand Select to its largest option (#1247)
* Button.HideShadow replaced by Button.Importance = LowImportance

* Deprecate NewContainer in favour of NewContainerWithoutLayout
* Deprecate HBox and VBox in favour of new container APIs
* Move Container.AddObject to Container.Add matching Container.Remove
* Start move from widget.TabContainer to container.AppTabs
* Replace Radio with RadioGroup
* Deprecate WidgetRenderer.BackgroundColor

### Fixed

* Support focus traversal in dialog (#948), (#948)
* Add missing AbsolutePosition in some mouse events (#1274)
* Don't let scrollbar handle become too small
* Ensure tab children are resized before being shown (#1331)
* Don't hang if OpenURL loads browser (#1332)
* Content not filling dialog (#1360)
* Overlays not adjusting on orientation change in mobile (#1334)
* Fix missing key events for some keypad keys (#1325)
* Issue with non-english folder names in Linux favourites (#1248)
* Fix overlays escaping screen interactive bounds (#1358)
* Key events not blocked by overlays (#814)
* Update scroll container content if it is changed (#1341)
* Respect SelectEntry datta changes on refresh (#1462)
* Incorrect SelectEntry dropdown button position (#1361)
* don't allow both single and double tap events to fire (#1381)
* Fix issue where long or tall images could jump on load (#1266, #1432)
* Weird behaviour when resizing or minimizing a ScrollContainer (#1245)
* Fix panic on NewTextGrid().Text()
* Fix issue where scrollbar could jump after mousewheel scroll
* Add missing raster support in software render
* Respect GOOS/GOARCH in fyne command utilities
* BSD support in build tools
* SVG Cache could return the incorrect resource (#1479)

* Many optimisations and widget performance enhancements
* Various fixes to file creation and saving on mobile devices


## 1.3.3 - 10 August 2020

### Added

* Use icons for file dialog favourites (#1186)
* Add ScrollContainer ScrollToBottom and ScrollToTop

### Changed

* Make file filter case sensitive (#1185)

### Fixed

* Allow popups to create dialogs (#1176)
* Use default cursor for dragging scrollbars (#1172)
* Correctly parse SVG files with missing X/Y for rect
* Fix visibility of Entry placeholder when text is set (#1193)
* Fix encoding issue with Windows notifications (#1191)
* Fix issue where content expanding on Windows could freeze (#1189)
* Fix errors on Windows when reloading Fyne settings (#1165)
* Dialogs not updating theme correctly (#1201)
* Update the extended progressbar on refresh (#1219)
* Segfault if font fails (#1200)
* Slider rendering incorrectly when window maximized (#1223)
* Changing form label not refreshed (#1231)
* Files and folders starting "." show no name (#1235)


## 1.3.2 - 11 July 2020

### Added

* Linux packaged apps now include a Makefile to aid install

### Changed

* Fyne package supports specific architectures for Android
* Reset missing textures on refresh
* Custom confirm callbacks now called on implicitly shown dialogs
* SelectEntry can update drop-down list during OnChanged callback
* TextGrid whitespace color now matches theme changes
* Order of Window Resize(), SetFixedSize() and CenterOnScreen() does no matter before Show()
* Containers now refresh their visuals as well as their Children on Refresh()

### Fixed

* Capped StrokeWidth on canvas.Line (#831)
* Canvas lines, rectangles and circles do not resize and refresh correctly
* Black flickering on resize on MacOS and OS X (possibly not on Catalina) (#1122)
* Crash when resizing window under macOS (#1051, #1140)
* Set SetFixedSize to true, the menus are overlapped (#1105)
* Ctrl+v into text input field crashes app. Presumably clipboard is empty (#1123, #1132)
* Slider default value doesn't stay inside range (#1128)
* The position of window is changed when status change from show to hide, then to show (#1116)
* Creating a windows inside onClose handler causes Fyne to panic (#1106)
* Backspace in entry after SetText("") can crash (#1096)
* Empty main menu causes panic (#1073)
* Installing using `fyne install` on Linux now works on distrubutions that don't use `/usr/local`
* Fix recommendations from staticcheck
* Unable to overwrite file when using dialog.ShowFileSave (#1168)


## 1.3 - 5 June 2020

### Added

* File open and save dialogs (#225)
* Add notifications support (#398)
* Add text wrap support (#332)
* Add Accordion widget (#206)
* Add TextGrid widget (#115)
* Add SplitContainer widget (#205)
* Add new URI type and handlers for cross-platform data access
* Desktop apps can now create splash windows
* Add ScaleMode to images, new ImageScalePixels feature for retro graphics
* Allow widgets to influence mouse cursor style (#726)
* Support changing the text on form submit/cancel buttons
* Support reporting CapsLock key events (#552)
* Add OnClosed callback for Dialog
* Add new image test helpers for validating render output
* Support showing different types of soft keyboard on mobile devices (#971, #975)

### Changed

* Upgraded underlying GLFW library to fix various issues (#183, #61)
* Add submenu support and hover effects (#395)
* Default to non-premultiplied alpha (NRGBA) across toolkit
* Rename FixedGridLayout to GridWrapLayout (deprecate old API) (#836)
* Windows redraw and animations continue on window resize and move
* New...PopUp() methods are being replaced by Show...Popup() or New...Popup().Show()
* Apps started on a goroutine will now panic as this is not supported
* On Linux apps now simulate 120DPI instead of 96DPI
* Improved fyne_settings scale picking user interface
* Reorganised fyne_demo to accommodate growing collection of widgets and containers
* Rendering now happens on a different thread to events for more consistent drawing
* Improved text selection on mobile devices

### Fixed (highlights)

* Panic when trying to paste empty clipboard into entry (#743)
* Scale does not match user configuration in Windows 10 (#635)
* Copy/Paste not working on Entry Field in Windows OS (#981)
* Select widgets with many options overflow UI without scrolling (#675)
* android: typing in entry expands only after full refresh (#972)
* iOS app stops re-drawing mid frame after a while (#950)
* Too many successive GUI updates do not properly update the view (904)
* iOS apps would not build using Apple's new certificates
* Preserve aspect ratio in SVG stroke drawing (#976)
* Fixed many race conditions in widget data handling
* Various crashes and render glitches in extended widgets
* Fix security issues reported by gosec (#742)


## 1.2.4 - 13 April 2020

### Added

 * Added Direction field to ScrollContainer and NewHScrollContainer, NewVScrollContainer constructors (#763)
 * Added Scroller.SetMinSize() to enable better defaults for scrolled content
 * Added "fyne vendor" subcommand to help packaging fyne dependencies in projects
 * Added "fyne version" subcommand to help with bug reporting (#656)
 * Clipboard (cut/copy/paste) is now supported on iOS and Android (#414)
 * Preferences.RemoveValue() now allows deletion of a stored user preference

### Changed

 * Report keys based on name not key code - fixes issue with shortcuts with AZERTY (#790)

### Fixed

 * Mobile builds now support go modules (#660)
 * Building for mobile would try to run desktop build first
 * Mobile apps now draw the full safe area on a screen (#799)
 * Preferences were not stored on mobile apps (#779)
 * Window on Windows is not controllable after exiting FullScreen mode (#727)
 * Soft keyboard not working on some Samsung/LG smart phones (#787)
 * Selecting a tab on extended TabContainer doesn't refresh button (#810)
 * Appending tab to empty TabContainer causes divide by zero on mobile (#820)
 * Application crashes on startup (#816)
 * Form does not always update on theme change (#842)


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

* Correct the color of Entry widget cursor if theme changes
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

