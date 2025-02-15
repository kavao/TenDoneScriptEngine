package ui

import (
	"gameengine/src/engine/font"

	"github.com/hajimehoshi/ebiten/v2"
	ebitentext "github.com/hajimehoshi/ebiten/v2/text"
)

type TextRenderer struct {
	fontManager *font.FontManager
}

func NewTextRenderer(fm *font.FontManager) *TextRenderer {
	return &TextRenderer{
		fontManager: fm,
	}
}

func (tr *TextRenderer) DrawText(screen *ebiten.Image, text string, x, y float64, style font.TextStyle) error {
	face, err := tr.fontManager.GetFontFace(style.FontID, style.Size, style.Style)
	if err != nil {
		return err
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(x, y)

	ebitentext.Draw(screen, text, face, int(x), int(y), style.Color)
	return nil
}
