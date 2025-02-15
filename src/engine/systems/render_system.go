package systems

import (
	"gameengine/src/engine/components"
	"gameengine/src/engine/ecs"

	"github.com/hajimehoshi/ebiten/v2"
)

type RenderSystem struct {
	*ecs.BaseSystem
	screen *ebiten.Image
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.PriorityRender, []ecs.ComponentID{1, 2}), // Transform と Sprite のID
	}
}

func (s *RenderSystem) Update(dt float64) error {
	if s.screen == nil {
		return nil
	}

	for _, entity := range s.BaseSystem.Entities() {
		var transform *components.TransformComponent
		var sprite *components.SpriteComponent

		transform = entity.GetComponent(1).(*components.TransformComponent)
		sprite = entity.GetComponent(2).(*components.SpriteComponent)

		if sprite.Sprite == nil {
			continue
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(transform.X, transform.Y)
		s.screen.DrawImage(sprite.Sprite, op)
	}
	return nil
}

func (s *RenderSystem) SetScreen(screen *ebiten.Image) {
	s.screen = screen
}
