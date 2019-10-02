This document outlines the current additional information for mobile development.
We should wrap this into our current tooling before release!

# Pre-requisites

Install gomobile (the go compiler for mobile apps)

`go get golang.org/x/mobile/cmd/gomobile`

## Android

Get and install the Andoid Native Development Kit (**NDK**). This is platform dependent and may be packaged 
with various tools that you might already have installed (like Android Studio)
Once installed set up the appropriate path, the directory you reference should contain `ndk-build`,
`platforms` and `toolchains`.

`export ANDROID_NDK_HOME=/<path to>/android-ndk`

## iOS

Install **XCode** from the app store and sign up for an [Apple Developer](https://developer.apple.com) account.
Before compiling you will need to add a device ID to your account, create a development certificate and
then a providioning profile.

# Compilation

Apps can then be compiled simply like the following:

`gomobile build -target=android fyne.io/fyne/cmd/hello`

For iOS apps you must also specify the bundleID of the app which must match an installed provisioning profile.

`gomobile build -target=ios -bundleid=com.teamid.appid fyne.io/fyne/cmd/hello`

