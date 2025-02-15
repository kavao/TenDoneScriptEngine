package components

import (
	"gameengine/src/engine/ecs"
)

type TransformComponent struct {
	entity  *ecs.Entity
	X, Y    float64
	ScaleX  float64
	ScaleY  float64
	Rotation float64
}

func (c *TransformComponent) GetEntity() *ecs.Entity {
	return c.entity
}

func (c *TransformComponent) SetEntity(e *ecs.Entity) {
	c.entity = e
}

func (c *TransformComponent) GetID() ecs.ComponentID {
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