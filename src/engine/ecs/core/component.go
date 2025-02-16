package core

import (
	"fmt"
	"sort"
	"sync"
)

// 基本的な型定義
type EntityID uint64
type ComponentID uint32
type SystemPriority int

// コンポーネントインターフェース
type Component interface {
	GetID() ComponentID
	GetEntity() *Entity
	SetEntity(entity *Entity)
	OnAdd()
	OnRemove()
}

// Systemインターフェース
type System interface {
	Update(dt float64) error
	GetPriority() SystemPriority
	GetRequiredComponents() []ComponentID
	HasRequiredComponents(entity *Entity) bool
	OnEntityAdded(entity *Entity)
	OnEntityRemoved(entity *Entity)
}

// World structを定義
type World struct {
	Mutex        sync.RWMutex
	Entities     map[EntityID]*Entity
	Systems      []System
	ToAdd        []*Entity
	ToRemove     []EntityID
	NextEntityID EntityID
}

// Entity型の定義
type Entity struct {
	ID         EntityID
	World      *World
	Components map[ComponentID]Component
	mutex      sync.RWMutex
	Active     bool
	Tags       map[string]bool
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

var DebugMode = false // パッケージレベルで定義

// World structのメソッド
func (w *World) Update(dt float64) error {
	w.Mutex.Lock()
	toAdd := w.ToAdd
	toRemove := w.ToRemove
	w.ToAdd = w.ToAdd[:0]
	w.ToRemove = w.ToRemove[:0]
	w.Mutex.Unlock()

	// 削除処理を追加
	for _, id := range toRemove {
		w.Mutex.Lock()
		if entity, exists := w.Entities[id]; exists {
			delete(w.Entities, id)
			systems := w.Systems
			w.Mutex.Unlock()

			for _, system := range systems {
				system.OnEntityRemoved(entity)
			}
		} else {
			w.Mutex.Unlock()
		}
	}

	// システムの更新（ロック外で実行）
	for _, entity := range toAdd {
		w.Mutex.Lock()
		w.Entities[entity.GetID()] = entity
		systems := w.Systems // コピーを作成
		w.Mutex.Unlock()

		for _, system := range systems {
			if system.HasRequiredComponents(entity) {
				system.OnEntityAdded(entity)
			}
		}
	}

	// システムの更新
	w.Mutex.RLock()
	systems := w.Systems // コピーを作成
	w.Mutex.RUnlock()

	for _, system := range systems {
		if err := system.Update(dt); err != nil {
			return err
		}
	}

	return nil
}

func (w *World) CreateEntity() *Entity {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()

	entity := &Entity{
		ID:         w.NextEntityID,
		World:      w,
		Components: make(map[ComponentID]Component),
		Active:     true,
		Tags:       make(map[string]bool),
	}
	w.Entities[entity.GetID()] = entity
	w.NextEntityID++
	return entity
}

func (w *World) GetEntity(id EntityID) *Entity {
	w.Mutex.RLock()
	defer w.Mutex.RUnlock()
	return w.Entities[id]
}

// Entityのメソッド
func (e *Entity) AddComponent(component Component) {
	e.mutex.Lock()
	id := component.GetID()
	e.Components[id] = component
	component.SetEntity(e)
	e.mutex.Unlock()

	if DebugMode {
		fmt.Printf("Added component %d to entity %d\n", id, e.ID)
	}

	e.World.Mutex.Lock()
	for _, system := range e.World.Systems {
		if system.HasRequiredComponents(e) {
			if DebugMode {
				fmt.Printf("Entity %d now matches system requirements\n", e.ID)
			}
			system.OnEntityAdded(e)
		}
	}
	e.World.Mutex.Unlock()
}

func (e *Entity) GetID() EntityID {
	return e.ID
}

// Entityのその他のメソッド
func (e *Entity) HasComponent(id ComponentID) bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	_, exists := e.Components[id]
	return exists
}

func (e *Entity) GetComponent(id ComponentID) Component {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.Components[id]
}

// Entityの追加メソッド
func (e *Entity) Deactivate() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.Active = false
}

func (e *Entity) Activate() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.Active = true
}

func (e *Entity) IsActive() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.Active
}

func (e *Entity) AddTag(tag string) {
	if e.Tags == nil {
		e.Tags = make(map[string]bool)
	}
	e.Tags[tag] = true
}

func (e *Entity) RemoveTag(tag string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.Tags, tag)
}

func (e *Entity) HasTag(tag string) bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	value, exists := e.Tags[tag]
	return exists && value
}

func (w *World) FindEntitiesByTag(tag string) []*Entity {
	w.Mutex.RLock()
	defer w.Mutex.RUnlock()

	var entities []*Entity
	for _, entity := range w.Entities {
		if entity.IsActive() && entity.HasTag(tag) {
			entities = append(entities, entity)
		}
	}
	return entities
}

func (w *World) GetTotalEntities() int {
	w.Mutex.RLock()
	defer w.Mutex.RUnlock()
	return len(w.Entities)
}

func NewWorld() *World {
	return &World{
		Entities: make(map[EntityID]*Entity),
		Systems:  make([]System, 0),
		ToAdd:    make([]*Entity, 0),
		ToRemove: make([]EntityID, 0),
	}
}

func (w *World) AddSystem(system System) {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()

	w.Systems = append(w.Systems, system)

	// 優先度でソート
	sort.Slice(w.Systems, func(i, j int) bool {
		return w.Systems[i].GetPriority() < w.Systems[j].GetPriority()
	})

	fmt.Printf("Adding system, checking %d entities\n", len(w.Entities))
	// 既存のエンティティをシステムに追加
	for id, entity := range w.Entities {
		fmt.Printf("Checking entity %d for system\n", id)
		if entity.IsActive() && system.HasRequiredComponents(entity) {
			fmt.Printf("Adding entity %d to system\n", id)
			system.OnEntityAdded(entity)
		}
	}
}

func (w *World) CleanupInactiveEntities() {
	w.Mutex.Lock()
	defer w.Mutex.Unlock()

	for id, entity := range w.Entities {
		if !entity.IsActive() {
			delete(w.Entities, id)
		}
	}
}

func (w *World) DestroyEntity(id EntityID) {
	// 単一のロックで処理
	w.Mutex.Lock()
	defer w.Mutex.Unlock()

	if entity, exists := w.Entities[id]; exists {
		entity.Active = false
	}
}
