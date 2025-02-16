package components

import core "gameengine/src/engine/ecs/core"

type TextComponent struct {
	entity  *core.Entity
	Text    string
	X, Y    float64
	Visible bool
}

func NewTextComponent() *TextComponent {
	return &TextComponent{
		Visible: true,
	}
}

func (c *TextComponent) GetEntity() *core.Entity  { return c.entity }
func (c *TextComponent) SetEntity(e *core.Entity) { c.entity = e }
func (c *TextComponent) GetID() core.ComponentID  { return 3 } // TextComponentのID
func (c *TextComponent) OnAdd()                   {}
func (c *TextComponent) OnRemove()                {}
