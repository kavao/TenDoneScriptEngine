package components

import (
	core "gameengine/src/engine/ecs/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteComponent struct {
	entity *core.Entity
	Image  string
	Sprite *ebiten.Image
	Width  int
	Height int
	Layer  int
	Color  color.Color
}

func (c *SpriteComponent) GetEntity() *core.Entity {
	return c.entity
}

func (c *SpriteComponent) SetEntity(e *core.Entity) {
	c.entity = e
}

func (c *SpriteComponent) GetID() core.ComponentID {
	return 2 // SpriteComponentのID
}

func (c *SpriteComponent) OnAdd()    {}
func (c *SpriteComponent) OnRemove() {}

func NewSpriteComponent() *SpriteComponent {
	// デフォルトで32x32の白い四角形を作成
	img := ebiten.NewImage(32, 32)
	img.Fill(color.White)

	return &SpriteComponent{
		Layer:  0,
		Sprite: img,
		Width:  32,
		Height: 32,
		Color:  color.White,
	}
}

// 色を設定するメソッド
func (c *SpriteComponent) SetColor(col color.Color) {
	c.Color = col
	c.Sprite.Fill(col)
}
