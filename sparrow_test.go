package sparrow_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"image"
	"image/draw"
	"image/png"
	"os"
	"testing"

	. "github.com/MatusOllah/sparrow"
)

func TestTextureAtlas(t *testing.T) {
	boyfriendXML, err := os.ReadFile("testdata/BOYFRIEND.xml")
	if err != nil {
		t.Error(err)
	}

	bf, err := ParseTextureAtlas(boyfriendXML)
	if err != nil {
		t.Error(err)
	}

	bfHey, err := bf.GetSubTexture("BF HEY!!0025")
	if err != nil {
		t.Error(err)
	}

	expectedBFHey := &SubTexture{
		Name:        "BF HEY!!0025",
		X:           6216,
		Y:           509,
		Width:       414,
		Height:      412,
		FrameX:      -1,
		FrameY:      -6,
		FrameWidth:  415,
		FrameHeight: 418,
	}

	if *bfHey != *expectedBFHey {
		t.Errorf("expected sub-texture %v, got %v\n", expectedBFHey, bfHey)
	}
}

func TestEnumerateSubTextures(t *testing.T) {
	boyfriendXML, err := os.ReadFile("testdata/BOYFRIEND.xml")
	if err != nil {
		t.Error(err)
	}

	bf, err := ParseTextureAtlas(boyfriendXML)
	if err != nil {
		t.Error(err)
	}

	expectedBFHey := &SubTexture{
		Name:        "BF HEY!!0025",
		X:           6216,
		Y:           509,
		Width:       414,
		Height:      412,
		FrameX:      -1,
		FrameY:      -6,
		FrameWidth:  415,
		FrameHeight: 418,
	}

	bfHey := bf.EnumerateSubTextures()["BF HEY!!0025"]

	if *bfHey != *expectedBFHey {
		t.Errorf("expected sub-texture %v, got %v\n", expectedBFHey, bfHey)
	}
}

func TestImage(t *testing.T) {
	boyfriendXML, err := os.ReadFile("testdata/BOYFRIEND.xml")
	if err != nil {
		t.Error(err)
	}

	boyfriendPNG, err := os.Open("testdata/BOYFRIEND.png")
	if err != nil {
		t.Error(err)
	}
	defer boyfriendPNG.Close()

	bfImg, err := png.Decode(boyfriendPNG)
	if err != nil {
		t.Error(err)
	}

	bfHeyFile, err := os.ReadFile("testdata/BF HEY!!0025.png")
	if err != nil {
		t.Error(err)
	}

	bf, err := ParseTextureAtlas(boyfriendXML)
	if err != nil {
		t.Error(err)
	}

	bfHey, err := bf.MustGetSubTexture("BF HEY!!0025").Image(bfImg)
	if err != nil {
		t.Error(err)
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, imageToNRGBA(bfHey))
	if err != nil {
		t.Error(err)
	}

	bfHeySum := sha256sum(buf.Bytes())
	expectedBFHeySum := sha256sum(bfHeyFile)

	if bfHeySum != expectedBFHeySum {
		t.Errorf("expected image sha256 sum %v, got %v\n", expectedBFHeySum, bfHeySum)
	}
}

func sha256sum(b []byte) string {
	hasher := sha256.New()
	hasher.Write(b)
	return hex.EncodeToString(hasher.Sum(nil))
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
