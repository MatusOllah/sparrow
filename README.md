# sparrow

[![Go Reference](https://pkg.go.dev/badge/github.com/MatusOllah/sparrow.svg)](https://pkg.go.dev/github.com/MatusOllah/sparrow)
[![Go Report Card](https://goreportcard.com/badge/github.com/MatusOllah/sparrow)](https://goreportcard.com/report/github.com/MatusOllah/sparrow)

**sparrow** is a Sparrow v2 texture atlas library for Go.

## Features

* Decoding
* Encoding
* Image / Frame extracting

## Basic Usage

```go
package main

import (
    "image/png"
    "github.com/MatusOllah/sparrow"
)

func main() {
    img, err := png.Decode("example.png")
    if err != nil {
        panic(err)
    }

    atlas, err := sparrow.ParseTextureAtlas("example.xml")
    if err != nil {
        panic(err)
    }

    // This gets / extracts the my_frame frame
    frame := atlas.MustGetSubTexture("my_frame")
    frameImg := frame.MustImage(img)

    // This returns all frames as a map
    frames := atlas.EnumerateSubTextures()
}
```

## License

MIT License
