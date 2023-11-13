package sparrow

import (
	"encoding/xml"
)

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

// Encode encodes a TextureAtlas as XML.
func (ta *TextureAtlas) Encode() ([]byte, error) {
	return xml.Marshal(ta)
}
