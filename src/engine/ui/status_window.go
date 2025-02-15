package ui

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type StatusWindow struct {
	BaseComponent
	Font            font.Face
	BackgroundColor color.Color
	TextColor       color.Color
	Padding         float64
	Stats           map[string]interface{}
	LineHeight      float64
}

func NewStatusWindow(font font.Face) *StatusWindow {
	return &StatusWindow{
		BaseComponent: BaseComponent{
			Visible: true,
			ZIndex:  90,
		},
		Font:            font,
		BackgroundColor: color.RGBA{0, 0, 0, 200},
		TextColor:       color.White,
		Padding:         10,
		Stats:           make(map[string]interface{}),
		LineHeight:      24,
	}
}

func (w *StatusWindow) SetStat(key string, value interface{}) {
	w.Stats[key] = value
}

func (w *StatusWindow) Draw(screen *ebiten.Image) {
	if !w.Visible {
		return
	}

	// 背景描画
	bounds := w.GetBounds()
	windowImage := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	windowImage.Fill(w.BackgroundColor)

	// ステータス描画
	i := 0
	for key, value := range w.Stats {
		y := w.Padding + float64(i)*w.LineHeight
		text.Draw(
			windowImage,
			fmt.Sprintf("%s: %v", key, value),
			w.Font,
			int(w.Padding),
			int(y+w.LineHeight),
			w.TextColor,
		)
		i++
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.X, w.Y)
	screen.DrawImage(windowImage, op)
} 