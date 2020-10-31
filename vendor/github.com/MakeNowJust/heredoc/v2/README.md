# heredoc

[![Version](https://img.shields.io/github/v/release/MakeNowJust/heredoc)](https://github.com/MakeNowJust/heredoc/releases)
[![Build Status](https://circleci.com/gh/MakeNowJust/heredoc.svg?style=svg)](https://circleci.com/gh/MakeNowJust/heredoc)
[![GoDoc](https://godoc.org/github.com/MakeNowJusti/heredoc?status.svg)](https://godoc.org/github.com/MakeNowJust/heredoc)

## About

Package heredoc provides the here-document with keeping indent.

## Import

```go
import "github.com/MakeNowJust/heredoc/v2"
```

## Example

```go
package main

import (
	"fmt"

	"github.com/MakeNowJust/heredoc/v2"
)

func main() {
	fmt.Println(heredoc.Doc(`
		Lorem ipsum dolor sit amet, consectetur adipisicing elit,
		sed do eiusmod tempor incididunt ut labore et dolore magna
		aliqua. Ut enim ad minim veniam, ...
	`))
	// Output:
	// Lorem ipsum dolor sit amet, consectetur adipisicing elit,
	// sed do eiusmod tempor incididunt ut labore et dolore magna
	// aliqua. Ut enim ad minim veniam, ...
	//
}
```

## API Document

 - [heredoc - GoDoc](https://godoc.org/github.com/MakeNowJust/heredoc)

## License

This software is released under the MIT License, see LICENSE.
