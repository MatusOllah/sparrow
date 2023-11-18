# sparrow

[![Go Reference](https://pkg.go.dev/badge/github.com/MatusOllah/sparrow.svg)](https://pkg.go.dev/github.com/MatusOllah/sparrow)
[![Go Report Card](https://goreportcard.com/badge/github.com/MatusOllah/sparrow)](https://goreportcard.com/report/github.com/MatusOllah/sparrow)

**sparrow** is a Sparrow v2 texture atlas (PNG and XML pair, the same format that FNF uses) library for Go.

Handy if you're making a game that uses Sparrow v2 or rewriting FNF in Go (like me).

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
    // Here you can use any png & atlas, I'm using the awesome BOYFRIEND.xml
    img, err := png.Decode("BOYFRIEND.png")
    if err != nil {
        panic(err)
    }

    atlas, err := sparrow.ParseTextureAtlas("BOYFRIEND.xml")
    if err != nil {
        panic(err)
    }

    // This gets / extracts the BF HEY!!0025 frame
    bfHey := atlas.MustGetSubTexture("BF HEY!!0025")
    bfHeyImg := bfHey.MustImage(img)

    // This returns all frames as a map
    frames := atlas.EnumerateSubTextures()
}
```
