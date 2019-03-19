// Package icns implements an encoder for Apple's `.icns` file format.
// Reference: "https://en.wikipedia.org/wiki/Apple_Icon_Image_format".
//
// icns files allow for high resolution icons to make your apps look sexy.
// The most common ways to generate icns files are 1. use `iconutil` which is
// a Mac native cli utility, or 2. use tools that wrap `ImageMagick` which adds
// a large dependency to your project for such a simple use case.
//
// With this library you can use pure Go to create icns files from any source
// image, given that you can decode it into an `image.Image`, without any
// heavyweight dependencies or subprocessing required. You can also use this
// library to create icns files on windows and linux.
//
// A small CLI app `icnsify` is provided to allow you to create icns files
// using this library from the command line. It supports piping, which is
// something `iconutil` does not do, making it substantially easier to wrap.
//
// Note: All icons within the icns are sized for high dpi retina screens, using
// the appropriate icns OSTypes.
package icns
