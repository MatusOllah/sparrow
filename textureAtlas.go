package sparrow

import (
	"cmp"
	"encoding/xml"
	"errors"
	"io/fs"
	"log/slog"
	"slices"
)

var ErrSubTextureDoesNotExist error = errors.New("subtexture does not exist")

type TextureAtlas struct {
	ImagePath   string        `xml:"imagePath,attr"`
	SubTextures []*SubTexture `xml:"SubTexture"`
}

// NewTextureAtlas returns a new TextureAtlas.
func NewTextureAtlas() *TextureAtlas {
	return &TextureAtlas{}
}

// ParseTextureAtlas parses a Sparrow v2 texture atlas in XML format.
func ParseTextureAtlas(xmlData []byte) (*TextureAtlas, error) {
	var atlas TextureAtlas
	if err := xml.Unmarshal(xmlData, &atlas); err != nil {
		return nil, err
	}

	return &atlas, nil
}

// ParseTextureAtlasFromFS parses a Sparrow v2 texture atlas in XML format from a filesystem.
func ParseTextureAtlasFromFS(fsys fs.FS, path string) (*TextureAtlas, error) {
	xmlData, err := fs.ReadFile(fsys, path)
	if err != nil {
		return nil, err
	}

	atlas, err := ParseTextureAtlas(xmlData)
	if err != nil {
		return nil, err
	}

	return atlas, nil
}

// GetSubTexture finds and returns a SubTexture with the specified name.
func (ta *TextureAtlas) GetSubTexture(name string) (*SubTexture, error) {
	i, ok := slices.BinarySearchFunc(ta.SubTextures, &SubTexture{Name: name}, func(st1, st2 *SubTexture) int {
		return cmp.Compare(st1.Name, st2.Name)
	})
	if !ok {
		return nil, ErrSubTextureDoesNotExist
	}

	return ta.SubTextures[i], nil
}

// MustGetSubTexture simplay calls GetSubTexture and returns nil if an error occurs.
func (ta *TextureAtlas) MustGetSubTexture(name string) *SubTexture {
	st, err := ta.GetSubTexture(name)
	if err != nil {
		slog.Error(err.Error())
		return nil
	}

	return st
}

// EnumerateSubTextures returns a map that contains all subtextures.
func (ta *TextureAtlas) EnumerateSubTextures() map[string]*SubTexture {
	m := map[string]*SubTexture{}
	for _, st := range ta.SubTextures {
		m[st.Name] = st
	}

	return m
}

// Encode encodes a TextureAtlas as XML.
func (ta *TextureAtlas) Encode() ([]byte, error) {
	return xml.Marshal(ta)
}
