# gendex

## How to run

From project root:

```console
cd ../../../../cmd/fyne/internal/mobile/
go run ./gendex
cd ../../../../
go install ./cmd/fyne
```

It will generate the `./cmd/fyne/internal/mobile/dex.go` file
that will be used in the `fyne` CLI for your next builds.
