# icns

Easily convert `.jpg` and `.png` to `.icns` with the command line tool `icnsify`, or use the library to convert from any `image.Image` to `.icns`.

`go get github.com/jackmordaunt/icns`

`icns` files allow for high resolution icons to make your apps look sexy. The most common ways to generate icns files are:

1. `iconutil`, which is a Mac native cli utility.
2. `ImageMagick` which adds a large dependency to your project for such a simple use case.

With this library you can use pure Go to create `icns` files from any source image, given that you can decode it into an `image.Image`, without any heavyweight dependencies or subprocessing required. You can also use it to create icns files on windows and linux (thanks Go).

A small CLI app `icnsify` is provided allowing you to create icns files using this library from the command line. It supports piping, which is something `iconutil` does not do, making it substantially easier to wrap or chuck into a shell pipeline.

Note: All icons within the `icns` are sized for high dpi retina screens, using the appropriate `icns` OSTypes.

## GUI

`preview` is a gui for displaying `icns` files cross-platform.

### Go Tool

```
go install github.com/jackmordaunt/icns/cmd/preview@latest
```

### Clone

```
git clone https://github.com/jackmordaunt/icns
cd icns/cmd/preview && go install .
```

Note: Gio cannot be cross-compiled right now, so there are no `preview` builds in releases.
Note: `preview` has it's own `go.mod` and therefore is versioned independently (unversioned).

![preview](docs/preview.png)

## Command Line

### Go Tool

```
go install github.com/jackmordaunt/icns/v2/cmd/icnsify@latest
```

### Clone

```
git clone https://github.com/jackmordaunt/icns
cd icns && go install ./cmd/icnsify
```

Pipe it

`cat icon.png | icnsify > icon.icns`

`cat icon.icns | icnsify > icon.png`

Standard

`icnsify -i icon.png -o icon.icns`

`icnsify -i icon.icns -o icon.png`

## Library

`go get github.com/jackmordaunt/icns/v2`

```go
func main() {
        pngf, err := os.Open("path/to/icon.png")
        if err != nil {
                log.Fatalf("opening source image: %v", err)
        }
        defer pngf.Close()
        srcImg, _, err := image.Decode(pngf)
        if err != nil {
                log.Fatalf("decoding source image: %v", err)
        }
        dest, err := os.Create("path/to/icon.icns")
        if err != nil {
                log.Fatalf("opening destination file: %v", err)
        }
        defer dest.Close()
        if err := icns.Encode(dest, srcImg); err != nil {
                log.Fatalf("encoding icns: %v", err)
        }
}
```

## Roadmap

- [x] Encoder: `image.Image -> .icns`
- [x] Command Line Interface
  - [x] Encoding
  - [x] Pipe support
  - [x] Decoding
- [x] Implement Decoder: `.icns -> image.Image`
- [ ] Symmetric test: `decode(encode(img)) == img`

## Coffee

If this software is useful to you, consider buying me a coffee!

[https://liberapay.com/JackMordaunt](https://liberapay.com/JackMordaunt)
