# Changelog

This file lists the main changes with each version of the Fyne toolkit.
More detailed release notes can be found on the [releases page](https://github.com/fyne-io/fyne/releases). 

## 2.4.4 - 13 February 2024

### Fixed

* Spaces could be appended to linux Exec command during packaging
* Secondary mobile windows would not size correctly when padded
* Setting Icon.Resource to nil will not clear rendering
* Dismiss iOS keyboard if "Done" is tapped
* Large speed improvement in Entry and GridWrap widgets
* tests fail with macOS Assertion failure in NSMenu (#4572)
* Fix image test failures on Apple Silicon
* High CPU use when showing CustomDialogs (#4574)
* Entry does not show the last (few) changes when updating a binding.String in a fast succession (#4082)
* Calling Entry.SetText and then Entry.Bind immediately will ignore the bound value (#4235)
* Changing theme while application is running doesn't change some parameters on some widgets (#4344)
* Check widget: hovering/tapping to the right of the label area should not activate widget (#4527)
* Calling entry.SetPlaceHolder inside of OnChanged callback freezes app (#4516)
* Hyperlink enhancement: underline and tappable area shouldn't be wider than the text label (#3528)
* Fix possible compile error from go-text/typesetting


## 2.4.3 - 23 December 2023

### Fixed

* Fix OpenGL init for arm64 desktop devices 
* System tray icon on Mac is showing the app ID (#4416)
* Failure with fyne release -os android/arm (#4174)
* Android GoBack with forcefully close the app even if the keyboard is up (#4257)
* *BSD systems using the wrong (and slow) window resize
* Optimisations to reduce memory allocations in List, GridWrap, driver and mime type handling
* Reduce calls to C and repeated size checks in painter and driver code


## 2.4.2 - 22 November 2023

### Fixed

* Markdown only shows one horizontal rule (#4216)
* Spacer in HBox with hidden item will cause an additional trailing padding (#4259)
* Application crash when fast clicking the folders inside the file dialog (#4260)
* failed to initialise OpenGL (#437)
* App panic when clicking on a notification panel if there's a systray icon (#4385)
* Systray cannot be shown on Ubuntu (#3678, #4381)
* failed to initialise OpenGL on Windows dual-chip graphics cards (#437)
* Reduce memory allocations for each frame painted
* RichText may not refresh if segments manually replaced
* Correct URI.Extension() documentation
* Update for security fixes to x/sys and x/net
* Inconsistent rendering of Button widget (#4243)
* PasswordEntry initial text is not obscured (#4312)
* Pasting text in Entry does not update cursor position display (#4181)


## 2.4.1 - 9 October 2023

### Fixed

* Left key on tree now collapses open branch
* Avoid memory leak in Android driver code
* Entry Field on Android in Landscape Mode Shows "0" (#4036)
* DocTabs Indicator remains visible after last tab is removed (#4220)
* Some SVG resources don't update appearance correctly with the theme (#3900)
* Fix mobile simulation builds on OpenBSD
* Fix alignment of menu button on mobile
* Fix Compilation with Android NDK r26
* Clicking table headers causes high CPU consumption (#4264)
* Frequent clicking on table may cause the program to not respond (#4210)
* Application stops responding when scrolling a table (#4263)
* Possible crash parsing malformed JSON color (#4270)
* NewFolderOpen: incomplete filenames (#2165)
* Resolve issue where storage.List could crash with short URI (#4271)
* TextTruncateEllipsis abnormally truncates strings with multi-byte UTF-8 characters (#4283)
* Last character doesn't appear in Select when there is a special character (#4293)
* Resolve random crash in DocTab (#3909)
* Selecting items from a list caused the keyboard to popup on Android (#4236)


## 2.4.0 - 1 September 2023

### Added

* Rounded corners in rectangle (#1090)
* Support for emoji in text
* Layout debugging (with `-tags debug` build flag) (#3314)
* GridWrap collection widget
* Add table headers (#1658, #3594)
* Add mobile back button handling (#2910)
* Add option to disable UI animations (#1813)
* Text truncation ellipsis (#1659)
* Add support for binding tree data, include new `NewTreeWithData`
* Add support for OpenType fonts (#3245)
* Add `Window.SetOnDropped` to handle window-wide item drop on desktop
* Add lists to the types supported by preferences API
* Keyboard focus handling for all collection widgets
* Add APIs for refreshing individual items in collections (#3826)
* Tapping slider moves it to that position (#3650)
* Add `OnChangeEnded` callback to `Slider` (#3652)
* Added keyboard controls to `Slider`
* Add `NewWarningThemedResource` and `NewSuccessThemedResource` along with `NewColoredResource` (#4040)
* Custom hyperlink callback for rich text hyperlinks (#3335)
* Added `dialog.NewCustomWithoutButtons`, with a `SetButtons` method (#2127, #2782)
* Added `SetConfirmImportance` to `dialog.ConfirmDialog`.
* Added `FormDialog.Submit()` to close and submit the dialog if validation passes
* Rich Text image alignment (#3810)
* Bring back `theme.HyperlinkColor` (#3867)
* Added `Importance` field on `Label` to color the text
* Navigating in entry quickly with ctrl key (#2462)
* Support `.desktop` file metadata in `FyneApp.toml` for Linux and BSD
* Support mobile simulator on FreeBSD
* Add data binding boolean operators `Not`, `And` and `Or`
* Added `Entry.Append`, `Select.SetOptions`, `Check.SetText`, `FormDialog.Submit`
* Add `ShowPopUpAtRelativePosition` and `PopUp.ShowAtRelativePosition`
* Add desktop support to get key modifiers with `CurrentKeyModifiers`
* Add geometry helpers `NewSquareSize` and `NewSquareOffsetPos`
* Add `--pprof` option to fyne build commands to enable profiling
* Support compiling from Android (termux)

### Changed

* Go 1.17 or later is now required.
* Theme updated for rounded corners on buttons and input widgets
* `widget.ButtonImportance` is now `widget.Importance`
* The `Max` container and layout have been renamed `Stack` for clarity
* Refreshing an image will now happen in app-thread not render process, apps may wish to add async image load
* Icons for macOS bundles are now padded and rounded, disable with "-use-raw-icon" (#3752)
* Update Android target SDK to 33 for Play Store releases
* Focus handling for List/Tree/Table are now at the parent widget not child elements
* Accordion widget now fills available space - put it inside a `VBox` container for old behavior (#4126)
* Deprecated theme.FyneLogo() for later removal (#3296)
* Improve look of menu shortcuts (#2722)
* iOS and macOS packages now default to using "XCWildcard" provisioning profile
* Improving performance of lookup for theme data
* Improved application startup time

### Fixed

* Rendering performance enhancements
* `dialog.NewProgressInfinite` is deprecated, but dialog.NewCustom isn't equivalent
* Mouse cursor desync with Split handle when dragging (#3791)
* Minor graphic glitch with checkbox (#3792)
* binding.String===>Quick refresh *b.val will appear with new data reset by a call to OnChange (#3774)
* Fyne window becomes unresponsive when in background for a while (#2791)
* Hangs on repeated calls to `Select.SetSelected` in table. (#3684)
* `Select` has wrong height, padding and border (#4142)
* `widget.ImageSegment` can't be aligned. (#3505)
* Memory leak in font metrics cache (#4108)
* Don't panic when loading preferences with wrong type (#4039)
* Button with icon has wrong padding on right (#4124)
* Preferences don't all save when written in `CloseIntercept` (#3170)
* Text size does not update in Refresh for TextGrid
* DocTab selection underline not updated when deleting an Item (#3905)
* Single line Entry throws away selected text on submission (#4026)
* Significantly improve performance of large `TextGrid` and `Tree` widgets
* `List.ScrollToBottom` not scrolling to show the totality of the last Item (#3829)
* Setting `Position1` of canvas.Circle higher than `Position2` causes panic. (#3949)
* Enhance scroll wheel/touchpad scroll speed on desktop (#3492)
* Possible build issue on Windows with app metadata
* `Form` hint text has confusing padding to next widget (#4137)
* `Entry` Placeholder Style Only Applied On Click (#4035)
* Backspace and Delete key Do not Fire OnChanged Event (#4117)
* Fix `ProgressBar` text having the wrong color sometimes
* Window doesn't render when called for the first time from system tray and the last window was closed (#4163)
* Possible race condition in preference change listeners
* Various vulnerabilities resolved through updating dependencies 
* Wrong background for color dialog (#4199)


## 2.3.5 - 6 June 2023

### Fixed

* Panic with unsupported font (#3646)
* Temporary manifest file not closed after building on Windows
* Panic when using autogenerated quit menu and having unshown windows (#3870)
* Using `canvas.ImageScaleFastest` not working on arm64 (#3891)
* Disabled password Entry should also disable the ActionItem (#3908)
* Disabled RadioGroup does not display status (#3882)
* Negative TableCellID Row (#2857)
* Make sure we have sufficient space for the bar as well if content is tiny (#3898)
* Leak in image painter when replacing image.Image source regularly
* Links in Markdown/Rich Text lists breaks formatting (#2911)
* Crash when reducing window to taskbar with popup opened (#3877)
* RichText vertical scroll will truncate long content with horizontal lines (#3929)
* Custom metadata would not apply with `fyne release` command
* Horizontal CheckGroup overlap when having long text (#3005)
* Fix focused colour of coloured buttons (#3462)
* Menu separator not visible with light theme (#3814)


## 2.3.4 - 3 May 2023

### Fixed

* Memory leak when switching theme (#3640)
* Systray MenuItem separators not rendered in macOS root menu (#3759)
* Systray leaks window handles on Windows (#3760)
* RadioGroup miscalculates label widths in horizontal mode (#3386)
* Start of selection in entry is shifted when moving too fast (#3804)
* Performance issue in widget.List (#3816)
* Moving canvas items (e.g. Images) does not cause canvas repaint (#2205)
* Minor graphic glitch with checkbox (#3792)
* VBox and HBox using heap memory that was not required
* Menu hover is slow on long menus

## 2.3.3 - 24 March 2023

### Fixed

* Linux, Windows and BSD builds could fail if gles was missing


## 2.3.2 - 20 March 2023

### Fixed

* Fyne does not run perfectly on ARM-based MacOS platforms (#3639) *
* Panic on closing window in form submit on Мac M2 (#3397) *
* Wobbling slider effect for very small steps (#3648)
* Fix memory leak in test canvas refresh
* Optimise text texture memory by switching to single channel
* Packaging an android fyne app that uses tags can fail (#3641)
* NewAdaptiveGrid(0) blanks app window on start until first resize on Windows (#3669)
* Unnecessary refresh when sliding Split container
* Linux window resize refreshes all content
* Themed and unthemed svg resources can cache collide
* When packaging an ampersand in "Name" causes an error (#3195)
* Svg in ThemedResource without viewBox does not match theme (#3714)
* Missing menu icons in Windows system tray
* Systray Menu Separators don't respect the submenu placement (#3642)
* List row focus indicator disappears on scrolling (#3699)
* List row focus not reset when row widget is reused to display a new item (#3700)
* Avoid panic if accidental 5th nil is passed to Border container
* Mobile simulator not compiling on Apple M1/2
* Cropped letters in certain cases with the new v2.3.0 theme (#3500)

Many thanks indeed to [Dymium](https://dymium.io) for sponsoring an Apple
M2 device which allowed us to complete the marked (*) issues.


## 2.3.1 - 13 February 2023

### Changed

* Pad app version to ensure Windows packages correctly (#3638)

### Fixed

* Custom shortcuts with fyne.KeyTab is not working (#3087)
* Running a systray app with root privileges resulted in panic (#3120)
* Markdown image with no title is not parsed (#3577)
* Systray app on macOS panic when started while machine sleeps (#3609)
* Runtime error with VNC on RaspbianOS (#2972)
* Hovered background in List widget isn't reset when scrolling reuses an existing list item (#3584)
* cmd/fyne package can't find FyneApp.toml when -src option has given (#3459)
* TextWrapWord will cause crash in RichText unverified (#3498)
* crash in widget.(*RichText).lineSizeToColumn (#3292)
* Crash in widget.(*Entry).SelectedText (#3290)
* Crash in widget.(*RichText).updateRowBounds.func1 (#3291)
* window is max size at all times (#3507)
* systray.Quit() is not called consistently when the app is closing (#3597)
* Software rendering would ignore scale for text
* crash when minimize a window which contains a stroked rectangle (#3552)
* Menu item would not appear disabled initially
* Wrong icon colour for danger and warning buttons
* Embedding Fyne apps in iFrame alignment issue
* Generated metadata can be in wrong directory
* Android RootURI may not exist when used for storage (#3207)


## 2.3.0 - 24 December 2022

### Added

* Shiny new theme that was designed for us
* Improved text handling to support non-latin alphabets
* Add cloud storage and preference support
* Add menu icon and submenu support to system tray menus
* More button importance levels `ErrorImportance`, `WarningImportance`
* Support disabling of `AppTabs` and `DocTabs` items
* Add image support to rich text (#2366)
* Add CheckGroup.Remove (#3124)

### Changed

* The buttons on the default theme are no longer transparent, but we added more button importance types
* Expose a storage.ErrNotExists for non existing documents (#3083)
* Update `go-gl/glfw` to build against latest Glfw 3.3.8
* List items in `widget.List` now implement the Focusable interface

### Fixed

* Displaying unicode or different language like Bengali doesn't work (#598)
* Cannot disable container.TabItem (#1904)
* Update Linux/XDG application theme to follow the FreeDesktop Dark Style Preference (#2657)
* Running `fyne package -os android` needs NDK 16/19c (#3066)
* Caret position lost when resizing a MultilineEntry (#3024)
* Fix possible crash in table resize (#3369)
* Memory usage surge when selecting/appending MultilineEntry text (#3426)
* Fyne bundle does not support appending when parameter is a directory
* Crash parsing invalid file URI (#3275)
* Systray apps on macOS can only be terminated via the systray menu quit button (#3395)
* Wayland Scaling support: sizes and distances are scaled wrong (#2850)
* Google play console minimum API level 31 (#3375)
* Data bound entry text replacing selection is ignored (#3340)
* Split Container does not respect item's Visible status (#3232)
* Android - Entry - OnSubmitted is not working (#3267)
* Can't set custom CGO_CFLAGS and CGO_LDFLAGS with "fyne package" on darwin (#3276)
* Text line not displayed in RichText (#3117)
* Segfault when adding items directly in form struct (#3153)
* Preferences RemoveValue does not save (#3229)
* Create new folder directly from FolderDialog (#3174)
* Slider drag handle is clipped off at minimum size (#2966)
* Entry text "flickering" while typing (#3461)
* Rendering of not changed canvas objects after an event (#3211)
* Form dialog not displaying hint text and validation errors (#2781)


## 2.2.4 - 9 November 2022

### Fixes

* Iphone incorrect click coordinates in zoomed screen view (#3122)
* CachedFontFace seems to be causing crash (#3134)
* Fix possible compile error if "fyne build" is used without icon metadata
* Detect and use recent Android NDK toolchain
* Handle fyne package -release and fyne release properly for Android and iOS
* Fix issue with mobile simulation when systray used
* Fix incorrect size and position for radio focus indicator (#3137)


## 2.2.3 - 8 July 2022

### Fixed

* Regression: Preferences are not parsed at program start (#3125)
* Wrappable RichText in a Split container causes crash (#3003, #2961)
* meta.Version is always 1.0.0 on android & ios (#3109)


## 2.2.2 - 30 June 2022

### Fixed

* Windows missing version metadata when packaged (#3046)
* Fyne package would not build apps using old Fyne versions
* System tray icon may not be removed on app exit in Windows
* Emphasis in Markdown gives erroneous output in RichText (#2974)
* When last visible window is closed, hidden window is set visible (#3059)
* Do not close app when last window is closed but systrayMenu exists (#3092)
* Image with ImageFillOriginal not showing (#3102)


## 2.2.1 - 12 June 2022

### Fixed

* Fix various race conditions and compatibility issues with System tray menus
* Resolve issue where macOS systray menu may not appear
* Updated yaml dependency to fix CVE-2022-28948
* Tab buttons stop working after removing a tab (#3050)
* os.SetEnv("FYNE_FONT") doesn't work in v2.2.0 (#3056)


## 2.2.0 - 7 June 2022

### Added

* Add SetIcon method on ToolbarAction (#2475)
* Access compiled app metadata using new `App.Metadata()` method
* Add support for System tray icon and menu (#283)
* Support for Android Application Bundle (.aab) (#2663)
* Initial support for OpenBSD and NetBSD
* Add keyboard shortcuts to menu (#682)
* Add technical preview of web driver and `fyne serve` command
* Added `iossimulator` build target (#1917)
* Allow dynamic themes via JSON templates (#211)
* Custom hyperlink callback (#2979)
* Add support for `.ico` file when compiling for windows (#2412)
* Add binding.NewStringWithFormat (#2890)
* Add Entry.SetMinRowsVisible
* Add Menu.Refresh() and MainMenu.Refresh() (#2853)
* Packages for Linux and BSD now support installing into the home directory
* Add `.RemoveAll()` to containers
* Add an AllString validator for chaining together string validators

### Changed

* Toolbar item constructors now return concrete types instead of ToolbarItem
* Low importance buttons no longer draw button color as a background
* ProgressBar widget height is now consistent with other widgets
* Include check in DocTabs menu to show current tab
* Don't call OnScrolled if offset did not change (#2646)
* Prefer ANDROID_NDK_HOME over the ANDROID_HOME ndk-bundle location (#2920)
* Support serialisation / deserialisation of the widget tree (#5)
* Better error reporting / handling when OpenGL is not available (#2689)
* Memory is now better reclaimed on Android when the OS requests it
* Notifications on Linux and BSD now show the application icon
* Change listeners for preferences no longer run when setting the same value
* The file dialog now shows extensions in the list view for better readability
* Many optimisations and widget performance enhancements
* Updated various dependencies to their latest versions

### Fixed

* SendNotification does not show app name on Windows (#1940)
* Copy-paste via keyboard don't work translated keyboard mappings on Windows (#1220)
* OnScrolled triggered when offset hasn't changed (#1868)
* Carriage Return (\r) is rendered as space (#2456)
* storage.List() returns list with nil elements for empty directories (#2858)
* Entry widget, position of cursor when clicking empty space (#2877)
* SelectEntry cause UI hang (#2925)
* Font cutoff with bold italics (#3001)
* Fyne error: Preferences load error (#2936, 3015)
* Scrolled List bad redraw when window is maximized (#3013)
* Linux and BSD packages not being installable if the name contained spaces


## 2.1.4 - 17 March 2022

### Fixed

* SetTheme() is not fully effective for widget.Form (#2810)
* FolderOpenDialog SetDismissText is ineffective (#2830)
* window.Resize() does not work if SetFixedSize(true) is set after (#2819)
* Container.Remove() race causes crash (#2826, #2775, #2481)
* FixedSize Window improperly sized if contains image with ImageFillOriginal (#2800)


## 2.1.3 - 24 February 2022

### Fixed

* The text on button can't be show correctly when use imported font (#2512)
* Fix issues with DocTabs scrolling (#2709)
* Fix possible crash for tapping extended Radio or Check item
* Resolve lookup of relative icons in FyneApp.toml
* Window not shown when SetFixedSize is used without Resize (#2784)
* Text and links in markdown can be rendered on top of each other (#2695)
* Incorrect cursor movement in a multiline entry with wrapping (#2698)


## 2.1.2 - 6 December 2021

### Fixed

* Scrolling list bound to data programmatically causes nil pointer dereference (#2549)
* Rich text from markdown can get newlines wrong (#2589)
* Fix crash on 32bit operating systems (#2603)
* Compile failure on MacOS 10.12 Sierra (#2478)
* Don't focus widgets on mobile where keyboard should not display (#2598)
* storage.List doesn't return complete URI on Android for "content:" scheme (#2619)
* Last word of the line and first word of the next line are joined in markdown parse (#2647)
* Support for building `cmd/fyne` on Windows arm64
* Fixed FreeBSD requiring installed glfw library dependency (#1928)
* Apple M1: error when using mouse drag to resize window (#2188)
* Struct binding panics in reload with slice field (#2607)
* File Dialog favourites can break for certain locations (#2595)
* Define user friendly names for Android Apps (#2653)
* Entry validator not updating if content is changed via data binding after SetContent (#2639)
* CenterOnScreen not working for FixedSize Window (#2550)
* Panic in boundStringListItem.Get() (#2643)
* Can't set an app/window icon to be an svg. (#1196)
* SetFullScreen(false) can give error (#2588)


## 2.1.1 - 22 October 2021

### Fixed

* Fix issue where table could select cells beyond data bound
* Some fast taps could be ignored (#2484)
* iOS app stops re-drawing mid-frame after a while (#950)
* Mobile simulation mode did not work on Apple M1 computers
* TextGrid background color can show gaps in render (#2493)
* Fix alignment of files in list view of file dialog
* Crash setting visible window on macOS to fixed size (#2488)
* fyne bundle ignores -name flag in windows (#2395)
* Lines with nil colour would crash renderer
* Android -nm tool not found with NDK 23 (#2498)
* Runtime panic because out of touchID (#2407)
* Long text in Select boxes overflows out of the box (#2522)
* Calling SetText on Label may not refresh correctly
* Menu can be triggered by # key but not always Alt
* Cursor position updates twice with delay (#2525)
* widgets freeze after being in background and then a crash upon pop-up menu (#2536)
* too many Refresh() calls may now cause visual artifacts in the List widget (#2548)
* Entry.SetText may panic if called on a multiline entry with selected text (#2482)
* TextGrid not always drawing correctly when resized (#2501)


## 2.1.0 - 17 September 2021

### Added

* DocTabs container for handling multiple open files
* Lifecycle API for handling foreground, background and other event
* Add RichText widget and Markdown parser
* Add TabWidth to TextStyle to specify tab size in spaces
* Add CheckGroup widget for multi-select
* Add FyneApp.toml metadata file to ease build commands
* Include http and https in standard repositories
* Add selection color to themes
* Include baseline information in driver font measurement
* Document storage API (App.Storage().Create() and others)
* Add "App Files" to file dialog for apps that use document storage
* Tab overflow on AppTabs
* Add URI and Unbound type to data bindings
* Add keyboard support for menus, pop-ups and buttons
* Add SimpleRenderer to help make simple widgets (#709)
* Add scroll functions for List, Table, Tree (#1892)
* Add selection and disabling to MenuItem
* Add Alignment to widget.Select (#2329)
* Expose ScanCode for keyboard events originating from hardware (#1523)
* Support macOS GPU switching (#2423)

### Changed

* Focusable widgets are no longer focused on tap, add canvas.Focus(obj) in Tapped handler if required
* Move to background based selection for List, Table and Tree
* Update fyne command line tool to use --posix style parameters
* Switch from gz to xz compression for unix packages
* Performance improvements with line, text and raster rendering
* Items not yet visible can no longer be focused
* Lines can now be drawn down to 1px (instead of 1dp) (#2298)
* Support multiple lines of text on button (#2378)
* Improved text layout speed by caching string size calculations
* Updated to require Go 1.14 so we can use some new features
* Window Resize request is now asynchronous
* Up/Down keys take cursor home/end when on first/last lines respectively

### Fixed

* Correctly align text tabs (#1791)
* Mobile apps theme does not match system (#472)
* Toolbar with widget.Label makes the ToolbarAction buttons higher (#2257)
* Memory leaks in renderers and canvases cache maps (#735)
* FileDialog SetFilter does not work on Android devices (#2353)
* Hover fix for List and Tree with Draggable objects
* Line resize can flip slope (#2208)
* Deadlocks when using widgets with data (#2348)
* Changing input type with keyboard visible would not update soft keyboards
* MainMenu() Close item does NOT call function defined in SetCloseIntercept (#2355)
* Entry cursor position with mouse is offset vertically by theme.SizeNameInputBorder (#2387)
* Backspace key is not working on Android AOSP (#1941)
* macOS: 'NSUserNotification' has been deprecated (#1833)
* macOS: Native menu would add new items if refreshed
* iOS builds fail since Go 1.16
* Re-add support for 32 bit iOS devices, if built with Go 1.14
* Android builds fail on Apple M1 (#2439)
* SetFullScreen(true) before ShowAndRun fails (#2446)
* Interacting with another app when window.SetFullScreen(true) will cause the application to hide itself. (#2448)
* Sequential writes to preferences does not save to file (#2449)
* Correct Android keyboard handling (#2447)
* MIUI-Android: The widget’s Hyperlink cannot open the URL (#1514)
* Improved performance of data binding conversions and text MinSize


## 2.0.4 - 6 August 2021

### Changed

* Disable Form labels when the element it applys to is disabled (#1530)
* Entry popup menu now fires shortcuts so extended widgets can intercept
* Update Android builds to SDK 30

### Fixed

* sendnotification show appID for name on windows (#1940)
* Fix accidental removal of windows builds during cross-compile
* Removing an item from a container did not update layout
* Update title bar on Windows 10 to match OS theme (#2184)
* Tapped triggered after Drag (#2235)
* Improved documentation and example code for file dialog (#2156)
* Preferences file gets unexpectedly cleared (#2241)
* Extra row dividers rendered on using SetColumnWidth to update a table (#2266)
* Fix resizing fullscreen issue
* Fullscreen changes my display resolution when showing a dialog (#1832)
* Entry validation does not work for empty field (#2179)
* Tab support for focus handling missing on mobile
* ScrollToBottom not always scrolling all the way when items added to container.Scroller
* Fixed scrollbar disappearing after changing content (#2303)
* Calling SetContent a second time with the same content will not show
* Drawing text can panic when Color is nil (#2347)
* Optimisations when drawing transparent rectangle or whitespace strings


## 2.0.3 - 30 April 2021

### Fixed

* Optimisations for TextGrid rendering
* Data binding with widget.List sometimes crash while scrolling (#2125)
* Fix compilation on FreeBSD 13
* DataLists should notify only once when change.
* Keyboard will appear on Android in disabled Entry Widget (#2139)
* Save dialog with filename for Android
* form widget can't draw hinttext of appended item. (#2028)
* Don't create empty shortcuts (#2148)
* Install directory for windows install command contains ".exe"
* Fix compilation for Linux Wayland apps
* Fix tab button layout on mobile (#2117)
* Options popup does not move if a SelectEntry widget moves with popup open
* Speed improvements to Select and SelectEntry drop down
* theme/fonts has an apache LICENSE file but it should have SIL OFL (#2193)
* Fix build requirements for target macOS platforms (#2154)
* ScrollEvent.Position and ScrollEvent.AbsolutePosition is 0,0 (#2199)


## 2.0.2 - 1 April 2021

### Changed

* Text can now be copied from a disable Entry using keyboard shortcuts

### Fixed

* Slider offset position could be incorrect for mobile apps
* Correct error in example code
* When graphics init fails then don't try to continue running (#1593)
* Don't show global settings on mobile in fyne_demo as it's not supported (#2062)
* Empty selection would render small rectangle in Entry
* Do not show validation state for disabled Entry
* dialog.ShowFileSave did not support mobile (#2076)
* Fix issue that storage could not write to files on iOS and Android
* mobile app could crash in some focus calls
* Duplicate symbol error when compiling for Android with NDK 23 (#2064)
* Add internet permission by default for Android apps (#1715)
* Child and Parent support in storage were missing for mobile appps
* Various crashes with Entry and multiline selections (including #1989)
* Slider calls OnChanged for each value between steps (#1748)
* fyne command doesn't remove temporary binary from src (#1910)
* Advanced Color picker on mobile keeps updating values forever after sliding (#2075)
* exec.Command and widget.Button combination not working (#1857)
* After clicking a link on macOS, click everywhere in the app will be linked (#2112)
* Text selection - Shift+Tab bug (#1787)


## 2.0.1 - 4 March 2021

### Changed

* An Entry with `Wrapping=fyne.TextWrapOff` no longer blocks scroll events from a parent

### Fixed

* Dialog.Resize() has no effect if called before Dialog.Show() (#1863)
* SelectTab does not always correctly set the blue underline to the selected tab (#1872)
* Entry Validation Broken when using Data binding (#1890)
* Fix background colour not applying until theme change
* android runtime error with fyne.dialog (#1896)
* Fix scale calculations for Wayland phones (PinePhone)
* Correct initial state of entry validation
* fix entry widget mouse drag selection when scrolled
* List widget panic when refreshing after changing content length (#1864)
* Fix image caching that was too aggressive on resize
* Pointer and cursor misalignment in widget.Entry (#1937)
* SIGSEGV Sometimes When Closing a Program by Clicking a Button (#1604)
* Advanced Color Picker shows Black for custom primary color as RGBA (#1970)
* Canvas.Focus() before window visible causes application to crash (#1893)
* Menu over Content (#1973)
* Error compiling fyne on Apple M1 arm64 (#1739)
* Cells are not getting draw in correct location after column resize. (#1951)
* Possible panic when selecting text in a widget.Entry (#1983)
* Form validation doesn't enable submit button (#1965)
* Creating a window shows it before calling .Show() and .Hide() does not work (#1835)
* Dialogs are not refreshed correctly on .Show() (#1866)
* Failed creating setting storage : no such directory (#2023)
* Erroneous custom filter types not supported error on mobile (#2012)
* High importance button show no hovered state (#1785)
* List widget does not render all visible content after content data gets shorter (#1948)
* Calling Select on List before draw can crash (#1960)
* Dialog not resizing in newly created window (#1692)
* Dialog not returning to requested size (#1382)
* Entry without scrollable content prevents scrolling of outside scroller (#1939)
* fyne_demo crash after selecting custom Theme and table (#2018)
* Table widget crash when scrolling rapidly (#1887)
* Cursor animation sometimes distorts the text (#1778)
* Extended password entry panics when password revealer is clicked (#2036)
* Data binding limited to 1024 simultaneous operations (#1838)
* Custom theme does not refresh when variant changes (#2006)


## 2.0 - 22 January 2021

### Changes that are not backward compatible

These changes may break some apps, please read the 
[upgrading doc](https://developer.fyne.io/api/v2.0/upgrading) for more info
The import path is now `fyne.io/fyne/v2` when you are ready to make the update.

* Coordinate system to float32
  * Size and Position units were changed from int to float32
  * `Text.TextSize` moved to float32 and `fyne.MeasureText` now takes a float32 size parameter
  * Removed `Size.Union` (use `Size.Max` instead)
  * Added fyne.Delta for difference-based X, Y float32 representation
  * DraggedEvent.DraggedX and DraggedY (int, int) to DraggedEvent.Dragged (Delta)
  * ScrollEvent.DeltaX and DeltaY (int, int) moved to ScrollEvent.Scrolled (Delta)

* Theme API update
  * `fyne.Theme` moved to `fyne.LegacyTheme` and can be load to a new theme using `theme.FromLegacy`
  * A new, more flexible, Theme interface has been created that we encourage developers to use

* The second parameter of `theme.NewThemedResource` was removed, it was previously ignored
* The desktop.Cursor definition was renamed desktop.StandardCursor to make way for custom cursors
* Button `Style` and `HideShadow` were removed, use `Importance`

* iOS apps preferences will be lost in this upgrade as we move to more advanced storage
* Dialogs no longer show when created, unless using the ShowXxx convenience methods
* Entry widget now contains scrolling so should no longer be wrapped in a scroll container

* Removed deprecated types including:
  - `dialog.FileIcon` (now `widget.FileIcon`)
  - `widget.Radio` (now `widget.RadioGroup`)
  - `widget.AccordionContainer` (now `widget.Accordion`)
  - `layout.NewFixedGridLayout()` (now `layout.NewGridWrapLayout()`)
  - `widget.ScrollContainer` (now `container.Scroll`)
  - `widget.SplitContainer` (now `container.Spilt`)
  - `widget.Group` (replaced by `widget.Card`)
  - `widget.Box` (now `container.NewH/VBox`, with `Children` field moved to `Objects`)
  - `widget.TabContainer` and `widget.AppTabs` (now `container.AppTabs`)
* Many deprecated fields have been removed, replacements listed in API docs 1.4
  - for specific information you can browse https://developer.fyne.io/api/v1.4/

### Added

* Data binding API to connect data sources to widgets and sync data
  - Add preferences data binding and `Preferences.AddChangeListener`
  - Add bind support to `Check`, `Entry`, `Label`, `List`, `ProgressBar` and `Slider` widgets
* Animation API for handling smooth element transitions
  - Add animations to buttons, tabs and entry cursor
* Storage repository API for connecting custom file sources
  - Add storage functions `Copy`, `Delete` and `Move` for `URI`
  - Add `CanRead`, `CanWrite` and `CanList` to storage APIs
* New Theme API for easier customisation of apps
  - Add ability for custom themes to support light/dark preference
  - Support for custom icons in theme definition
  - New `theme.FromLegacy` helper to use old theme API definitions
* Add fyne.Vector for managing x/y float32 coordinates
* Add MouseButtonTertiary for middle mouse button events on desktop
* Add `canvas.ImageScaleFastest` for faster, less precise, scaling
* Add new `dialog.Form` that will phase out `dialog.Entry`
* Add keyboard control for main menu
* Add `Scroll.OnScrolled` event for seeing changes in scroll container
* Add `TextStyle` and `OnSubmitted` to `Entry` widget
* Add support for `HintText` and showing validation errors in `Form` widget
* Added basic support for tab character in `Entry`, `Label` and `TextGrid`

### Changed

* Coordinate system is now float32 - see breaking changes above
* ScrollEvent and DragEvent moved to Delta from (int, int)
* Change bundled resources to use more efficient string storage
* Left and Right mouse buttons on Desktop are being moved to `MouseButtonPrimary` and `MouseButtonSecondary`
* Many optimisations and widget performance enhancements

* Moving to new `container.New()` and `container.NewWithoutLayout()` constructors (replacing `fyne.NewContainer` and `fyne.NewContainerWithoutLayout`)
* Moving storage APIs `OpenFileFromURI`, `SaveFileToURI` and `ListerForURI` to `Reader`, `Writer` and `List` functions

### Fixed

* Validating a widget in widget.Form before renderer was created could cause a panic
* Added file and folder support for mobile simulation support (#1470)
* Appending options to a disabled widget.RadioGroup shows them as enabled (#1697)
* Toggling toolbar icons does not refresh (#1809)
* Black screen when slide up application on iPhone (#1610)
* Properly align Label in FormItem (#1531)
* Mobile dropdowns are too low (#1771)
* Cursor does not go down to next line with wrapping (#1737)
* Entry: while adding text beyond visible reagion there is no auto-scroll (#912)


## 1.4.3 - 4 January 2021

### Fixed

* Fix crash when showing file open dialog on iPadOS
* Fix possible missing icon on initial show of disabled button
* Capturing a canvas on macOS retina display would not capture full resolution
* Fix the release build flag for mobile
* Fix possible race conditions for canvas capture
* Improvements to `fyne get` command downloader
* Fix tree, so it refreshes visible nodes on Refresh()
* TabContainer Panic when removing selected tab (#1668)
* Incorrect clipping behaviour with nested scroll containers (#1682)
* MacOS Notifications are not shown on subsequent app runs (#1699)
* Fix the behavior when dragging the divider of split container (#1618)


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

