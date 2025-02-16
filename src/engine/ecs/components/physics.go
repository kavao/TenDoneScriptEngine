package components

import (
	core "gameengine/src/engine/ecs/core"
)

type PhysicsComponent struct {
	*core.BaseComponent
	VelocityX float64
	VelocityY float64
	Gravity   float64
	Speed     float64
}

func NewPhysicsComponent() *PhysicsComponent {
	return &PhysicsComponent{
		BaseComponent: core.NewBaseComponent(5),
		VelocityX:     0,
		VelocityY:     0,
		Gravity:       0,
		Speed:         3.0,
	}
}

func (c *PhysicsComponent) GetEntity() *core.Entity  { return c.BaseComponent.GetEntity() }
func (c *PhysicsComponent) SetEntity(e *core.Entity) { c.BaseComponent.SetEntity(e) }
func (c *PhysicsComponent) GetID() core.ComponentID  { return 5 }
func (c *PhysicsComponent) OnAdd()                   {}
func (c *PhysicsComponent) OnRemove()                {}
