package ecs

import (
	"gameengine/src/engine/ecs/core"
	"sort"
)

// システムの優先度
type SystemPriority = core.SystemPriority

const (
	PriorityPhysics SystemPriority = iota
	PriorityUpdate
	PriorityRender
)

// システムインターフェース
type System interface {
	Update(dt float64) error
	GetPriority() core.SystemPriority
	GetRequiredComponents() []core.ComponentID
	HasRequiredComponents(entity *core.Entity) bool
	OnEntityAdded(entity *core.Entity)
	OnEntityRemoved(entity *core.Entity)
}

// 基本的なシステム実装
type BaseSystem struct {
	priority           core.SystemPriority
	requiredComponents []core.ComponentID
	entities           []*core.Entity
	componentSignature map[core.ComponentID]bool
}

func NewBaseSystem(priority core.SystemPriority, requiredComponents []core.ComponentID) *BaseSystem {
	signature := make(map[core.ComponentID]bool)
	for _, id := range requiredComponents {
		signature[id] = true
	}

	return &BaseSystem{
		priority:           priority,
		requiredComponents: requiredComponents,
		entities:           make([]*core.Entity, 0),
		componentSignature: signature,
	}
}

func (s *BaseSystem) GetPriority() SystemPriority {
	return s.priority
}

func (s *BaseSystem) GetRequiredComponents() []core.ComponentID {
	return s.requiredComponents
}

func (s *BaseSystem) OnEntityAdded(entity *core.Entity) {
	s.entities = append(s.entities, entity)
	sort.Slice(s.entities, func(i, j int) bool {
		return s.entities[i].GetID() < s.entities[j].GetID()
	})
}

func (s *BaseSystem) OnEntityRemoved(entity *core.Entity) {
	for i, e := range s.entities {
		if e.GetID() == entity.GetID() {
			s.entities = append(s.entities[:i], s.entities[i+1:]...)
			break
		}
	}
}

func (s *BaseSystem) HasRequiredComponents(entity *core.Entity) bool {
	// デバッグ出力は開発時のみ有効に
	// fmt.Printf("Checking entity %d for required components\n", entity.ID)
	for id := range s.componentSignature {
		// fmt.Printf("  Checking component %d: %v\n", id, entity.HasComponent(id))
		if !entity.HasComponent(id) {
			return false
		}
	}
	// fmt.Printf("  Entity has all required components\n")
	return true
}

func (s *BaseSystem) Entities() []*core.Entity {
	return s.entities
}
