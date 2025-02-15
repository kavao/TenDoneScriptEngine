package ecs

import (
	"sync"
	"sync/atomic"
)

// エンティティID
type EntityID uint64

var nextEntityID uint64

// エンティティ
type Entity struct {
	id         EntityID
	components map[ComponentID]Component
	mutex      sync.RWMutex
	world      *World
	active     bool
	tags       map[string]bool
}

func NewEntity(world *World) *Entity {
	return &Entity{
		id:         EntityID(atomic.AddUint64(&nextEntityID, 1)),
		components: make(map[ComponentID]Component),
		world:      world,
		active:     true,
		tags:       make(map[string]bool),
	}
}

// コンポーネントの追加
func (e *Entity) AddComponent(component Component) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	id := component.GetID()
	e.components[id] = component
	component.SetEntity(e)

	// システムの更新
	e.world.entityComponentAdded(e, component)
}

// コンポーネントの削除
func (e *Entity) RemoveComponent(id ComponentID) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if component, exists := e.components[id]; exists {
		delete(e.components, id)
		e.world.entityComponentRemoved(e, component)
	}
}

// コンポーネントの取得
func (e *Entity) GetComponent(id ComponentID) Component {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.components[id]
}

// コンポーネントの存在確認
func (e *Entity) HasComponent(id ComponentID) bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	_, exists := e.components[id]
	return exists
}

// タグの追加
func (e *Entity) AddTag(tag string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.tags[tag] = true
}

// タグの削除
func (e *Entity) RemoveTag(tag string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.tags, tag)
}

// タグの確認
func (e *Entity) HasTag(tag string) bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.tags[tag]
}

// エンティティの無効化
func (e *Entity) Deactivate() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.active = false
}

// エンティティの有効化
func (e *Entity) Activate() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.active = true
}

// エンティティの状態確認
func (e *Entity) IsActive() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.active
}

// エンティティIDの取得
func (e *Entity) GetID() EntityID {
	return e.id
} 