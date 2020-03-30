# goinfo

goinfo is a simple tool that collect info about a Go project and the develop environment.

It aims to help with issue reporting answering the following questions:

- what version of Go are you using ?
- what operating system and processor architecture are you using ?
- what version of the project or module are you using ?

## Installation

```
$ go get github.com/lucor/goinfo/cmd/goinfo
```

## Usage

```
Usage:
  goinfo [options...]
Options:
  -work-dir         Path of the working dir. Default to current dir
  -module-path      Go module path to detect info. Default to the module defined in work-dir
  -format           Format output for the report. Supported: text, html, json. Default to text
  -help             Display this help text
```

## Examples

### Text format output

```
$ goinfo
```

```
## Go version info
version="go1.14 linux/amd64"

## Go module info
go_mod="true"
go_path="false"
imported="false"
module="github.com/lucor/goinfo"
path="/code/lucor/goinfo"
version="v0.9.0-g6626b18"

## OS info
BUG_REPORT_URL="https://bugs.launchpad.net/ubuntu/"
HOME_URL="https://www.ubuntu.com/"
ID="ubuntu"
ID_LIKE="debian"
NAME="Ubuntu"
PRETTY_NAME="Ubuntu 19.10"
PRIVACY_POLICY_URL="https://www.ubuntu.com/legal/terms-and-policies/privacy-policy"
SUPPORT_URL="https://help.ubuntu.com/"
UBUNTU_CODENAME="eoan"
VERSION="19.10 (Eoan Ermine)"
VERSION_CODENAME="eoan"
VERSION_ID="19.10"

## Go environment info
AR="ar"
CC="gcc"
CGO_CFLAGS="-g -O2"
CGO_CPPFLAGS=""
CGO_CXXFLAGS="-g -O2"
CGO_ENABLED="1"
CGO_FFLAGS="-g -O2"
CGO_LDFLAGS="-g -O2"
CXX="g++"
GCCGO="gccgo"
GO111MODULE=""
GOARCH="amd64"
GOBIN=""
GOCACHE="/home/lucor/.cache/go-build"
GOENV="/home/lucor/.config/go/env"
GOEXE=""
GOFLAGS=""
...

```

### HTML format output

This format could be useful when reporting issues on services that support the HTML details tag like github or gitlab

```
$ goinfo -format html
```

The report will look like:

<details><summary>Go version info</summary><br><pre>
version=go1.14 linux/amd64
</pre></details>

<details><summary>Go module info</summary><br><pre>
go_mod=true
go_path=false
imported=false
module=github.com/lucor/goinfo
path=/code/lucor/goinfo
version=v0.9.0-g6626b18
</pre></details>

<details><summary>OS info</summary><br><pre>
BUG_REPORT_URL=https://bugs.launchpad.net/ubuntu/
HOME_URL=https://www.ubuntu.com/
ID=ubuntu
ID_LIKE=debian
NAME=Ubuntu
PRETTY_NAME=Ubuntu 19.10
PRIVACY_POLICY_URL=https://www.ubuntu.com/legal/terms-and-policies/privacy-policy
SUPPORT_URL=https://help.ubuntu.com/
UBUNTU_CODENAME=eoan
VERSION=19.10 (Eoan Ermine)
VERSION_CODENAME=eoan
VERSION_ID=19.10
</pre></details>

<details><summary>Go environment info</summary><br><pre>
AR=ar
CC=gcc
CGO_CFLAGS=-g -O2
CGO_CPPFLAGS=
CGO_CXXFLAGS=-g -O2
CGO_ENABLED=1
CGO_FFLAGS=-g -O2
CGO_LDFLAGS=-g -O2
CXX=g&#43;&#43;
GCCGO=gccgo
GO111MODULE=
GOARCH=amd64
GOBIN=
GOCACHE=/home/lucor/.cache/go-build
GOENV=/home/lucor/.config/go/env
GOEXE=
GOFLAGS=
...
</pre></details>

### JSON format output


```
$ goinfo -format json
```

```
[
	{
		"summary": "Go version info",
		"info": {
			"version": "go1.14 linux/amd64"
		}
	},
	{
		"summary": "Go module info",
		"info": {
			"go_mod": true,
			"go_path": false,
			"imported": false,
			"module": "github.com/lucor/goinfo",
			"path": "/code/lucor/goinfo",
			"version": "v0.9.0-g6626b18"
		},
	},
	{
		"summary": "OS info",
		"info": {
			"BUG_REPORT_URL": "https://bugs.launchpad.net/ubuntu/",
			"HOME_URL": "https://www.ubuntu.com/",
			"ID": "ubuntu",
			"ID_LIKE": "debian",
			"NAME": "Ubuntu",
			"PRETTY_NAME": "Ubuntu 19.10",
			"PRIVACY_POLICY_URL": "https://www.ubuntu.com/legal/terms-and-policies/privacy-policy",
			"SUPPORT_URL": "https://help.ubuntu.com/",
			"UBUNTU_CODENAME": "eoan",
			"VERSION": "19.10 (Eoan Ermine)",
			"VERSION_CODENAME": "eoan",
			"VERSION_ID": "19.10"
		}
	},
	{
		"summary": "Go environment info",
		"info": {
			"AR": "ar",
			"CC": "gcc",
			"CGO_CFLAGS": "-g -O2",
			"CGO_CPPFLAGS": "",
			"CGO_CXXFLAGS": "-g -O2",
			"CGO_ENABLED": "1",
			"CGO_FFLAGS": "-g -O2",
			"CGO_LDFLAGS": "-g -O2",
			"CXX": "g++",
			"GCCGO": "gccgo",
			"GO111MODULE": "",
			"GOARCH": "amd64",
			"GOBIN": "",
			"GOCACHE": "/home/lucor/.cache/go-build",
			"GOENV": "/home/lucor/.config/go/env",
			"GOEXE": "",
			"GOFLAGS": "",
			...
		}
	}
]

```

## Contribute

- Fork and clone the repository
- Make and test your changes
- Open a pull request against the `master` branch
