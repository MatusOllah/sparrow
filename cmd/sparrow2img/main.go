package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/MatusOllah/sparrow"
	"github.com/jessevdk/go-flags"
	"github.com/pterm/pterm"
)

var img image.Image

func main() {
	// parse flags
	if _, err := flags.NewParser(&opts, flags.HelpFlag|flags.IgnoreUnknown|flags.PassDoubleDash).Parse(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// create output dir
	if err := os.MkdirAll(opts.Output, 0750); err != nil {
		panic(err)
	}

	// parse atlas
	parseSpin, err := pterm.DefaultSpinner.Start("Parsing atlas...")
	if err != nil {
		panic(err)
	}

	xmlData, err := os.ReadFile(opts.PosArgs.XMLPath)
	if err != nil {
		parseSpin.Fail(err)
		panic(err)
	}

	atlas, err := sparrow.ParseTextureAtlas(xmlData)
	if err != nil {
		parseSpin.Fail(err)
		panic(err)
	}

	// open image
	imgFile, err := os.Open(opts.PosArgs.PNGPath)
	if err != nil {
		parseSpin.Fail(err)
		panic(err)
	}
	defer imgFile.Close()

	_img, err := png.Decode(imgFile)
	if err != nil {
		parseSpin.Fail(err)
		panic(err)
	}
	img = _img

	parseSpin.Success()

	// extract
	if opts.SubTexture != "" {
		spin, err := pterm.DefaultSpinner.Start("Extracting " + opts.SubTexture)
		if err != nil {
			panic(err)
		}

		st, err := atlas.GetSubTexture(opts.SubTexture)
		if err != nil {
			spin.Fail(err)
			panic(err)
		}

		if err := extractSubImage(st); err != nil {
			spin.Fail(err)
			panic(err)
		}

		spin.Success()
		os.Exit(0)
	}

	pb, err := pterm.DefaultProgressbar.WithTotal(len(atlas.SubTextures)).WithTitle("Extracting ").Start()
	if err != nil {
		panic(err)
	}

	for _, st := range atlas.SubTextures {
		if opts.Verbose {
			pterm.Info.Printfln("Extracting %s", st.Name)
		}
		pb.UpdateTitle("Extracting " + st.Name)

		if err := extractSubImage(st); err != nil {
			pterm.Error.Printfln("Extracting %s failed: %v", st.Name, err)
			pb.Stop()
			panic(err)
		}

		if opts.Verbose {
			pterm.Success.Printfln("Extracting %s", st.Name)
		}
		pb.Increment()
	}
}

func extractSubImage(st *sparrow.SubTexture) error {
	f, err := os.Create(filepath.Join(opts.Output, st.Name+".png"))
	if err != nil {
		return err
	}
	defer f.Close()

	i, err := st.Image(img)
	if err != nil {
		return err
	}

	err = png.Encode(f, imageToNRGBA(i))
	if err != nil {
		return err
	}

	return nil
}

func imageToNRGBA(src image.Image) *image.NRGBA {
	if dst, ok := src.(*image.NRGBA); ok {
		return dst
	}

	b := src.Bounds()
	dst := image.NewNRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}
