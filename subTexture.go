package sparrow

import (
	"errors"
	"image"
	"image/draw"
	"log/slog"
)

var ErrCroppingUnsupported error = errors.New("image does not support cropping")

type subImager interface {
	SubImage(r image.Rectangle) image.Image
}

type SubTexture struct {
	Name        string `xml:"name,attr"`
	X           int    `xml:"x,attr"`
	Y           int    `xml:"y,attr"`
	Width       int    `xml:"width,attr"`
	Height      int    `xml:"height,attr"`
	FrameX      int    `xml:"frameX,attr"`
	FrameY      int    `xml:"frameY,attr"`
	FrameWidth  int    `xml:"frameWidth,attr"`
	FrameHeight int    `xml:"frameHeight,attr"`
}

// NewSubTexture returns a new SubTexture.
func NewSubTexture(name string, x, y, width, height int) *SubTexture {
	return &SubTexture{
		Name:   name,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// NewSubTextureFromRect returns a new SubTexture from r.
func NewSubTextureFromRect(name string, r image.Rectangle) *SubTexture {
	return NewSubTexture(name, r.Min.X, r.Min.Y, r.Dx(), r.Dy())
}

// Rect returns the st's uncropped rect.
func (st *SubTexture) Rect() image.Rectangle {
	return image.Rect(st.X, st.Y, st.X+st.Width, st.Y+st.Height)
}

// Image returns an true, cropped image representing the portion of the image i visible through the st's rect.
func (st *SubTexture) Image(i image.Image) (image.Image, error) {
	simg, ok := i.(subImager)
	if !ok {
		return nil, ErrCroppingUnsupported
	}

	return getTrueImage(simg.SubImage(st.Rect()), st.FrameX, st.FrameY, st.FrameWidth, st.FrameHeight), nil
}

// MustImage simply calls Image and returns nil if an error occurs.
func (st *SubTexture) MustImage(i image.Image) image.Image {
	img, err := st.Image(i)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	return img
}

func getTrueImage(i image.Image, fx, fy, fw, fh int) image.Image {
	final := imageToRGBA(i)

	if fx == 0 && fy == 0 && fw == 0 && fh == 0 {
		return final
	}

	if fx < 0 {
		final = padImg(final, 0, 0, 0, -fx)
	} else {
		final = final.SubImage(image.Rect(fx, 0, final.Bounds().Dx(), final.Bounds().Dy())).(*image.RGBA)
	}

	if fy < 0 {
		final = padImg(final, -fy, 0, 0, 0)
	} else {
		final = final.SubImage(image.Rect(0, fy, final.Bounds().Dx(), final.Bounds().Dy())).(*image.RGBA)
	}

	if fx+fw > i.Bounds().Dx() {
		final = padImg(final, 0, fx+fw-i.Bounds().Dx(), 0, 0)
	} else {
		final = final.SubImage(image.Rect(0, 0, fw, final.Bounds().Dy())).(*image.RGBA)
	}

	if fy+fh > i.Bounds().Dy() {
		final = padImg(final, 0, 0, fy+fh-i.Bounds().Dy(), 0)
	} else {
		final = final.SubImage(image.Rect(0, 0, final.Bounds().Dx(), fh)).(*image.RGBA)
	}

	return final
}

func padImg(i *image.RGBA, top, right, bottom, left int) *image.RGBA {
	width := i.Bounds().Dx()
	height := i.Bounds().Dy()

	newWidth := width + right + left
	newHeight := height + top + bottom

	result := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.Draw(result, result.Bounds(), i, i.Rect.Min.Sub(image.Pt(left, top)), draw.Over) // konecne to ide :D

	return result
}

func imageToRGBA(src image.Image) *image.RGBA {
	if dst, ok := src.(*image.RGBA); ok {
		return dst
	}

	b := src.Bounds()
	dst := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(dst, dst.Bounds(), src, b.Min, draw.Src)
	return dst
}
