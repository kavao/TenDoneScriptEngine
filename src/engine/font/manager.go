package font

import (
	"embed"
	"fmt"
	"image/color"
	"io/fs"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
)

//go:embed assets\mplus1p.ttf
var fontFS embed.FS

// フォントサイズのプリセット
const (
	SizeSmall  = 14
	SizeMedium = 18
	SizeLarge  = 24
	SizeHuge   = 32
)

// フォントスタイル
type FontStyle int

const (
	StyleRegular FontStyle = iota
	StyleBold
	StyleItalic
)

// フォントマネージャー
type FontManager struct {
	fonts     map[string]*sfnt.Font
	faces     map[fontKey]font.Face
	defaultID string
}

// フォントフェイスのキャッシュキー
type fontKey struct {
	fontID string
	size   float64
	style  FontStyle
}

func NewFontManager() (*FontManager, error) {
	fm := &FontManager{
		fonts: make(map[string]*sfnt.Font),
		faces: make(map[fontKey]font.Face),
	}

	// 組み込みフォントの読み込み
	if err := fm.loadEmbeddedFonts(); err != nil {
		return nil, fmt.Errorf("failed to load embedded fonts: %v", err)
	}

	return fm, nil
}

func (fm *FontManager) loadEmbeddedFonts() error {
	entries, err := fontFS.ReadDir("assets")
	if err != nil {
		return fmt.Errorf("failed to read font directory: %v", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		fontData, err := fs.ReadFile(fontFS, "assets/"+entry.Name())
		if err != nil {
			return fmt.Errorf("failed to read font file %s: %v", entry.Name(), err)
		}

		font, err := opentype.Parse(fontData)
		if err != nil {
			return fmt.Errorf("failed to parse font %s: %v", entry.Name(), err)
		}

		fontID := entry.Name()[:len(entry.Name())-4] // 拡張子を除去
		fm.fonts[fontID] = font

		// 最初に読み込んだフォントをデフォルトとして設定
		if fm.defaultID == "" {
			fm.defaultID = fontID
		}
	}

	return nil
}

func (fm *FontManager) GetFace(size float64, style FontStyle) (font.Face, error) {
	return fm.GetFontFace(fm.defaultID, size, style)
}

func (fm *FontManager) GetFontFace(fontID string, size float64, style FontStyle) (font.Face, error) {
	key := fontKey{fontID, size, style}

	// キャッシュされたフェイスがあれば返す
	if face, exists := fm.faces[key]; exists {
		return face, nil
	}

	// フォントの取得
	font, exists := fm.fonts[fontID]
	if !exists {
		return nil, fmt.Errorf("font not found: %s", fontID)
	}

	// フォントフェイスの生成
	face, err := opentype.NewFace(font, &opentype.FaceOptions{
		Size: size,
		DPI:  72,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %v", err)
	}

	// キャッシュに保存
	fm.faces[key] = face
	return face, nil
}

// テキスト描画用のヘルパー関数
type TextStyle struct {
	FontID      string
	Size        float64
	Style       FontStyle
	Color       color.Color
	LineSpacing float64
}

func NewDefaultTextStyle() TextStyle {
	return TextStyle{
		Size:        SizeMedium,
		Style:       StyleRegular,
		Color:       color.White,
		LineSpacing: 1.2,
	}
}
