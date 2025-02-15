package ecs

import "sort"

// システムの優先度
type SystemPriority int

const (
	PriorityPhysics SystemPriority = iota
	PriorityUpdate
	PriorityRender
)

// システムインターフェース
type System interface {
	Update(dt float64) error
	GetPriority() SystemPriority
	GetRequiredComponents() []ComponentID
	HasRequiredComponents(entity *Entity) bool
	OnEntityAdded(entity *Entity)
	OnEntityRemoved(entity *Entity)
}

// 基本的なシステム実装
type BaseSystem struct {
	priority            SystemPriority
	requiredComponents  []ComponentID
	entities            []*Entity
	componentSignature map[ComponentID]bool
}

func NewBaseSystem(priority SystemPriority, requiredComponents []ComponentID) *BaseSystem {
	signature := make(map[ComponentID]bool)
	for _, id := range requiredComponents {
		signature[id] = true
	}

	return &BaseSystem{
		priority:            priority,
		requiredComponents:  requiredComponents,
		entities:            make([]*Entity, 0),
		componentSignature: signature,
	}
}

func (s *BaseSystem) GetPriority() SystemPriority {
	return s.priority
}

func (s *BaseSystem) GetRequiredComponents() []ComponentID {
	return s.requiredComponents
}

func (s *BaseSystem) OnEntityAdded(entity *Entity) {
	s.entities = append(s.entities, entity)
	sort.Slice(s.entities, func(i, j int) bool {
		return s.entities[i].GetID() < s.entities[j].GetID()
	})
}

func (s *BaseSystem) OnEntityRemoved(entity *Entity) {
	for i, e := range s.entities {
		if e.GetID() == entity.GetID() {
			s.entities = append(s.entities[:i], s.entities[i+1:]...)
			break
		}
	}
}

func (s *BaseSystem) HasRequiredComponents(entity *Entity) bool {
	for id := range s.componentSignature {
		if !entity.HasComponent(id) {
			return false
		}
	}
	return true
} 