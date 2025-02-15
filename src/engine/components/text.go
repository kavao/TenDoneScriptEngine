package components

import (
	"gameengine/src/engine/ecs"
)

type TextComponent struct {
	entity  *ecs.Entity
	Text    string
	X, Y    float64
	Visible bool
}

func NewTextComponent() *TextComponent {
	return &TextComponent{
		Visible: true,
	}
}

func (c *TextComponent) GetEntity() *ecs.Entity { return c.entity }
func (c *TextComponent) SetEntity(e *ecs.Entity) { c.entity = e }
func (c *TextComponent) GetID() ecs.ComponentID { return 3 } // TextComponentのID
func (c *TextComponent) OnAdd() {}
func (c *TextComponent) OnRemove() {} 