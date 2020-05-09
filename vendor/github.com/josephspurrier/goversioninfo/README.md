GoVersionInfo
==========
[![Build Status](https://travis-ci.org/josephspurrier/goversioninfo.svg)](https://travis-ci.org/josephspurrier/goversioninfo) [![Coverage Status](https://coveralls.io/repos/josephspurrier/goversioninfo/badge.svg)](https://coveralls.io/r/josephspurrier/goversioninfo) [![GoDoc](https://godoc.org/github.com/josephspurrier/goversioninfo?status.svg)](https://godoc.org/github.com/josephspurrier/goversioninfo)

Microsoft Windows File Properties/Version Info and Icon Resource Generator for the Go Language

Package creates a syso file which contains Microsoft Windows Version Information and an optional icon. When you run "go build", Go will embed the version information and an optional icon and an optional manifest in the executable. Go will automatically use the syso file if it's in the same directory as the main() function.

Example of the file properties you can set using this package:

![Image of File Properties](https://cloud.githubusercontent.com/assets/2394539/12073634/0b32cb04-b0f6-11e5-9d8e-f9923ca554cf.jpg)

## Usage

To install, run the following command:
~~~
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
~~~

Copy testdata/resource/versioninfo.json into your working directory and then modify the file with your own settings.

Add a similar text to the top of your Go source code (-icon and -manifest are optional, but can also be specified in the versioninfo.json file):
~~~ go
//go:generate goversioninfo -icon=testdata/resource/icon.ico -manifest=testdata/resource/goversioninfo.exe.manifest
~~~

Run the Go commands in this order so goversioninfo will create a file called resource.syso in the same directory as the Go source code.
~~~
go generate
go build
~~~

## Command-Line Flags

Complete list of the flags for goversioninfo:

~~~
  -charset=0: charset ID
  -comment="": StringFileInfo.Comments
  -company="": StringFileInfo.CompanyName
  -copyright="": StringFileInfo.LegalCopyright
  -description="": StringFileInfo.FileDescription
  -example=false: just dump out an example versioninfo.json to stdout
  -file-version="": StringFileInfo.FileVersion
  -icon="": icon file name
  -internal-name="": StringFileInfo.InternalName
  -manifest="": manifest file name
  -o="resource.syso": output file name
  -platform-specific=false: output i386 and amd64 named resource.syso, ignores -o
  -original-name="": StringFileInfo.OriginalFilename
  -private-build="": StringFileInfo.PrivateBuild
  -product-name="": StringFileInfo.ProductName
  -product-version="": StringFileInfo.ProductVersion
  -special-build="": StringFileInfo.SpecialBuild
  -trademark="": StringFileInfo.LegalTrademarks
  -translation=0: translation ID
  -64:false: generate 64-bit binaries on true
  -ver-major=-1: FileVersion.Major
  -ver-minor=-1: FileVersion.Minor
  -ver-patch=-1: FileVersion.Patch
  -ver-build=-1: FileVersion.Build
  -product-ver-major=-1: ProductVersion.Major
  -product-ver-minor=-1: ProductVersion.Minor
  -product-ver-patch=-1: ProductVersion.Patch
  -product-ver-build=-1: ProductVersion.Build
~~~

You can look over the Microsoft Resource Information: [VERSIONINFO resource](https://msdn.microsoft.com/en-us/library/windows/desktop/aa381058(v=vs.85).aspx)

You can look through the Microsoft Version Information structures: [Version Information Structures](https://msdn.microsoft.com/en-us/library/windows/desktop/ff468916(v=vs.85).aspx)

## PowerShell Differences

In PowerShell, the version components are named differently than the fields in
the versioninfo.json file:

```
PowerShell:          versioninfo.json:
-----------          -----------------
FileMajorPart      = FileVersion.Major
FileMinorPart      = FileVersion.Minor
FileBuildPart      = FileVersion.Patch
FilePrivatePart    = FileVersion.Build
ProductMajorPart   = ProductVersion.Major
ProductMinorPart   = ProductVersion.Minor
ProductBuildPart   = ProductVersion.Patch
ProductPrivatePart = ProductVersion.Build

```

If you find any other differences, let me know.

## Alternatives to this Tool

You can also use [windres](https://sourceware.org/binutils/docs/binutils/windres.html) to create the syso file. The windres executable is available in either [MinGW](http://www.mingw.org/) or [tdm-gcc](http://tdm-gcc.tdragon.net/).

Below is a sample batch file you can use to create a .syso file from a .rc file. There are sample .rc files in the testdata/rc folder.

~~~
@ECHO OFF

SET PATH=C:\TDM-GCC-64\bin;%PATH%
REM SET PATH=C:\mingw64\bin;%PATH%

windres -i testdata/rc/versioninfo.rc -O coff -o versioninfo.syso

PAUSE
~~~

The information on how to create a .rc file is available [here](https://msdn.microsoft.com/en-us/library/windows/desktop/aa381043(v=vs.85).aspx). You can use the testdata/rc/versioninfo.rc file to create a .syso file that contains version info, icon, and manifest.

## Issues

The majority of the code for the creation of the syso file is from this package: [https://github.com/akavel/rsrc](https://github.com/akavel/rsrc)

There is an [issue](https://github.com/akavel/rsrc/issues/12) with adding the icon resource that prevents your application from being compressed or modified with a resource editor. Please use with caution.

## Major Contributions

Thanks to [Tamás Gulácsi](https://github.com/tgulacsi) for his superb code additions, refactoring, and optimization to make this a solid package.

Thanks to [Mateusz Czaplinski](https://github.com/akavel/rsrc) for his embedded binary resource package with icon and manifest functionality.
