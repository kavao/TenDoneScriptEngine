package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
)

type MenuItem struct {
	Text     string
	Enabled  bool
	OnSelect func()
}

type MenuWindow struct {
	BaseComponent
	Font            font.Face
	BackgroundColor color.Color
	TextColor       color.Color
	DisabledColor   color.Color
	SelectedColor   color.Color
	Padding         float64
	Items           []*MenuItem
	SelectedIndex   int
	LineHeight      float64
}

func NewMenuWindow(font font.Face) *MenuWindow {
	return &MenuWindow{
		BaseComponent: BaseComponent{
			Visible: true,
			ZIndex:  110,
		},
		Font:            font,
		BackgroundColor: color.RGBA{0, 0, 0, 200},
		TextColor:       color.White,
		DisabledColor:   color.RGBA{128, 128, 128, 255},
		SelectedColor:   color.RGBA{255, 255, 0, 255},
		Padding:         10,
		LineHeight:      24,
	}
}

func (w *MenuWindow) AddItem(text string, enabled bool, onSelect func()) {
	w.Items = append(w.Items, &MenuItem{
		Text:     text,
		Enabled:  enabled,
		OnSelect: onSelect,
	})
}

func (w *MenuWindow) ClearItems() {
	w.Items = nil
	w.SelectedIndex = 0
}

func (w *MenuWindow) Update() error {
	// キー入力処理
	// TODO: 入力システムとの連携
	return nil
}

func (w *MenuWindow) Draw(screen *ebiten.Image) {
	if !w.Visible {
		return
	}

	// 背景描画
	bounds := w.GetBounds()
	windowImage := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	windowImage.Fill(w.BackgroundColor)

	// メニュー項目描画
	for i, item := range w.Items {
		y := w.Padding + float64(i)*w.LineHeight
		textColor := w.TextColor
		if !item.Enabled {
			textColor = w.DisabledColor
		}
		if i == w.SelectedIndex {
			textColor = w.SelectedColor
		}

		text.Draw(
			windowImage,
			item.Text,
			w.Font,
			int(w.Padding),
			int(y+w.LineHeight),
			textColor,
		)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.X, w.Y)
	screen.DrawImage(windowImage, op)
}

func (w *MenuWindow) Select() {
	if w.SelectedIndex >= 0 && w.SelectedIndex < len(w.Items) {
		item := w.Items[w.SelectedIndex]
		if item.Enabled && item.OnSelect != nil {
			item.OnSelect()
		}
	}
} 