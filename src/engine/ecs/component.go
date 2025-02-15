package ecs

import "sync"

// コンポーネントID
type ComponentID uint32

// コンポーネントインターフェース
type Component interface {
	GetID() ComponentID
	GetEntity() *Entity
	SetEntity(entity *Entity)
	OnAdd()
	OnRemove()
}

// 基本的なコンポーネント実装
type BaseComponent struct {
	mutex  sync.RWMutex
	id     ComponentID
	entity *Entity
}

func NewBaseComponent(id ComponentID) *BaseComponent {
	return &BaseComponent{
		id: id,
	}
}

func (c *BaseComponent) GetID() ComponentID {
	return c.id
}

func (c *BaseComponent) GetEntity() *Entity {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.entity
}

func (c *BaseComponent) SetEntity(entity *Entity) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.entity = entity
}

func (c *BaseComponent) OnAdd()    {}
func (c *BaseComponent) OnRemove() {} 