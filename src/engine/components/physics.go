package components

import (
	"gameengine/src/engine/ecs"
)

type PhysicsComponent struct {
	entity    *ecs.Entity
	VelocityX float64
	VelocityY float64
	Gravity   float64
	Speed     float64
}

func NewPhysicsComponent() *PhysicsComponent {
	return &PhysicsComponent{
		VelocityX: 0,
		VelocityY: 0,
		Gravity:   0,
		Speed:     3.0,
	}
}

func (c *PhysicsComponent) GetEntity() *ecs.Entity  { return c.entity }
func (c *PhysicsComponent) SetEntity(e *ecs.Entity) { c.entity = e }
func (c *PhysicsComponent) GetID() ecs.ComponentID  { return 5 } // PhysicsComponentのID
func (c *PhysicsComponent) OnAdd()                  {}
func (c *PhysicsComponent) OnRemove()               {}
