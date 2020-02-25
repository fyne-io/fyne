This document outlines the current additional information for mobile development.
We should wrap this into our current tooling before release!

# Pre-requisites

Please ensure that your fyne install is the latest development version (both library and "fyne" command line tool).
```
go get -u fyne.io/fyne fyne.io/fyne/cmd/fyne
```

## Android

Get and install both the Android Software Development Kit (**SDK**) the Andoid Native Development Kit (**NDK**). This is platform dependent and may be packaged with various tools that you might already have installed (like Android Studio). The recommended approch is to simply install [Android Studio](https://developer.android.com/studio/index.html) and then install the the **NDK** by going to **Tools > SDK Manager** and from **SDK Tools** install the `NDK (Side by side)` package.

Once installed, the appropriate path to the **NDK** need to be set up. The directory you reference should contain `ndk-build`, `platforms` and `toolchains`. The default on Linux is `$HOME/Android/Sdk/ndk/version_of_your_installed_ndk`.

```bash
export ANDROID_NDK_HOME=/<path to>/android-ndk
```

## iOS

Install **XCode** from the app store and sign up for an [Apple Developer](https://developer.apple.com) account.
Before compiling you will need to add a device ID to your account, create a development certificate and
then a provisioning profile.

# Compilation

Apps can then be compiled simply like the following:

`fyne package -os android fyne.io/fyne/cmd/hello`

For iOS apps you must also specify the bundleID of the app which must match an installed provisioning profile.

`fyne package -os ios -appid=com.teamid.appid fyne.io/fyne/cmd/hello`

# Distribution

The compiled applications above can be tested locally and even shared with
developers that are part of your team. To deliver the applications to an app
store or marketplace takes more work which has not yet been automated.

## Android

To release your app needs to have debugging turned off. If you are using a custom manifest file then
you should add the following to AndroidManifest.xml:

`<application android:label="..." android:debuggable="false">`

This file should be placed in the `main` package, where you are running your builds from.
If you have not set up a manifest file then fyne can make the changes for you by 
adding a `-release` parameter to your package build:

`$ fyne package -os "android" -release ...`

You may also need to correct the alignment of the APK file before uploading it to the 
Play Store, the following command should work:

`${ANDROID_HOME}/build-tools/zipalign 4 myapp.apk aligned.apk`

The aligned file can then be uploaded for distribution.

## iOS

The above steps will create a .app folder that can be installed to your
iPhone directly. These steps work on that content to prepare it for uploading
to the App Store...
Be sure to find the appropriate values for each of the <placeholder> items!

The important part of this process is to re-codesign using a distribution certificate
which must match the app store provisioning profile in step 3.

* `mkdir Payload`
* `cp -r myapp.app Payload/`
* `cp ~/Library/MobileDevice/Provisioning\ Profiles/<app-store-distribution>.mobileprovision Payload/myapp.app/embedded.mobileprovision`
* `open Payload/myapp.app/Info.plist` (and check the build etc are correct - then save)
* create a new text file called entitlements.plist, enter the following content:
```
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>application-identifier</key>
    <string><teamId>.<appId></string>
</dict>
</plist>
```
* `codesign -f -vv -s "<distribution certificate name>" --entitlements entitlements.plist Payload/myapp.app/`
* If you have not already done so create an app specific password for your
Apple account ([appleid.com](https://appleid.apple.com))
* `zip -r myapp.ipa Payload/`
* `xcrun altool --upload-app --file myapp.ipa --username <appleId> --password <app-specific-password>`

The upload make take some time but once it is complete you should see the build
appear alongside the appropriate entry in your app store connect metadata.

This process will be improved soon!

