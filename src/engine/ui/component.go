package ui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// 基本的なUIコンポーネント
type Component interface {
	Draw(*ebiten.Image)
	Update() error
	SetPosition(x, y float64)
	SetSize(width, height float64)
	GetBounds() image.Rectangle
	IsVisible() bool
	SetVisible(bool)
}

// 基本的なコンポーネントの実装
type BaseComponent struct {
	X, Y          float64
	Width, Height float64
	Visible       bool
	ZIndex        int
}

func (b *BaseComponent) SetPosition(x, y float64) {
	b.X = x
	b.Y = y
}

func (b *BaseComponent) SetSize(width, height float64) {
	b.Width = width
	b.Height = height
}

func (b *BaseComponent) GetBounds() image.Rectangle {
	return image.Rect(
		int(b.X), int(b.Y),
		int(b.X+b.Width), int(b.Y+b.Height),
	)
}

func (b *BaseComponent) IsVisible() bool {
	return b.Visible
}

func (b *BaseComponent) SetVisible(visible bool) {
	b.Visible = visible
}

func (b *BaseComponent) Draw(screen *ebiten.Image) {}
func (b *BaseComponent) Update() error             { return nil }
