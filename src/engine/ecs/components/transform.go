package components

import core "gameengine/src/engine/ecs/core"

type TransformComponent struct {
	entity   *core.Entity
	X, Y     float64
	ScaleX   float64
	ScaleY   float64
	Rotation float64
}

func (c *TransformComponent) GetEntity() *core.Entity {
	return c.entity
}

func (c *TransformComponent) SetEntity(e *core.Entity) {
	c.entity = e
}

func (c *TransformComponent) GetID() core.ComponentID {
	return 1 // TransformComponentのID
}

func (c *TransformComponent) OnAdd()    {}
func (c *TransformComponent) OnRemove() {}

func NewTransformComponent() *TransformComponent {
	return &TransformComponent{
		ScaleX: 1.0,
		ScaleY: 1.0,
	}
}
