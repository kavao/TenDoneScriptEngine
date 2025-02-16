package systems

import (
	"gameengine/src/engine/ecs"
	"gameengine/src/engine/ecs/components"
	"gameengine/src/engine/ecs/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type RenderSystem struct {
	*ecs.BaseSystem
	screen *ebiten.Image
}

func NewRenderSystem() *RenderSystem {
	return &RenderSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.PriorityRender, []core.ComponentID{1, 2}), // Transform と Sprite のID
	}
}

func (s *RenderSystem) Update(dt float64) error {
	//	fmt.Printf("RenderSystem Update: screen=%v\n", s.screen != nil)
	if s.screen == nil {
		return nil
	}

	// fmt.Printf("RenderSystem checking required components: %v\n", s.BaseSystem.GetRequiredComponents())
	// for _, entity := range s.BaseSystem.Entities() {
	// 	fmt.Printf("Entity %d components: %v\n", entity.ID, entity.Components)
	// 	if s.BaseSystem.HasRequiredComponents(entity) {
	// 		fmt.Printf("Entity %d has required components\n", entity.ID)
	// 	}
	// }

	// fmt.Printf("RenderSystem entities: %d\n", len(s.BaseSystem.Entities()))
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
