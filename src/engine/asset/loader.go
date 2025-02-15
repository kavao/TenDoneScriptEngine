package asset

import (
	"bytes"
	"image"
	_ "image/png"
	"io/fs"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type AssetLoader struct {
	fs fs.FS
}

func NewAssetLoader(filesystem fs.FS) *AssetLoader {
	return &AssetLoader{
		fs: filesystem,
	}
}

// 画像のロード
func (l *AssetLoader) LoadImage(path string) (*ebiten.Image, error) {
	data, err := fs.ReadFile(l.fs, path)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return ebiten.NewImageFromImage(img), nil
}

// フォントのロード
func (l *AssetLoader) LoadFont(path string, size float64) (font.Face, error) {
	data, err := fs.ReadFile(l.fs, path)
	if err != nil {
		return nil, err
	}

	tt, err := opentype.Parse(data)
	if err != nil {
		return nil, err
	}

	return opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

// スクリプトのロード
func (l *AssetLoader) LoadScript(path string) (string, error) {
	data, err := fs.ReadFile(l.fs, path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
