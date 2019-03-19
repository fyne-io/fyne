rsrc - Tool for embedding binary resources in Go programs.

INSTALL: go get github.com/akavel/rsrc

USAGE:

rsrc [-manifest FILE.exe.manifest] [-ico FILE.ico[,FILE2.ico...]] -o FILE.syso
  Generates a .syso file with specified resources embedded in .rsrc section.
  The .syso file can be linked by Go linker when building Win32 executables.
  Icon embedded this way will show up on application's .exe instead of empty icon.
  Manifest file embedded this way will be recognized and detected by Windows.

The generated *.syso files should get automatically recognized by 'go build'
command and linked into an executable/library, as long as there are any *.go
files in the same directory.

OPTIONS:
  -arch="386": architecture of output file - one of: 386, [EXPERIMENTAL: amd64]
  -data="": path to raw data file to embed [WARNING: useless for Go 1.4+]
  -ico="": comma-separated list of paths to .ico files to embed
  -manifest="": path to a Windows manifest file to embed
  -o="rsrc.syso": name of output COFF (.res or .syso) file

Based on ideas presented by Minux.

In case anything does not work, it'd be nice if you could report (either via Github
issues, or via email to czapkofan@gmail.com), and please attach the input file(s)
which resulted in a problem, plus error message & symptoms, and/or any other details.

TODO MAYBE/LATER:
- fix or remove FIXMEs

LICENSE: MIT
  Copyright 2013-2017 The rsrc Authors.

http://github.com/akavel/rsrc
