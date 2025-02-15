package ecs

import (
	"sort"
	"sync"
)

// ワールド
type World struct {
	mutex        sync.RWMutex
	entities     map[EntityID]*Entity // マップを使用してエンティティを保存
	systems      []System
	toAdd        []*Entity
	toRemove     []EntityID
	nextEntityID EntityID
}

func NewWorld() *World {
	return &World{
		entities: make(map[EntityID]*Entity),
		systems:  make([]System, 0),
		toAdd:    make([]*Entity, 0),
		toRemove: make([]EntityID, 0),
	}
}

// エンティティの作成
func (w *World) CreateEntity() *Entity {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	entity := &Entity{
		id:         w.nextEntityID,
		world:      w,
		components: make(map[ComponentID]Component),
		tags:       make(map[string]bool),
		active:     true,
	}
	w.entities[entity.GetID()] = entity // 即座にマップに追加
	w.nextEntityID++
	return entity
}

// エンティティの削除
func (w *World) DestroyEntity(id EntityID) {
	w.toRemove = append(w.toRemove, id)
}

// システムの追加
func (w *World) AddSystem(system System) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.systems = append(w.systems, system)

	// 優先度でソート
	sort.Slice(w.systems, func(i, j int) bool {
		return w.systems[i].GetPriority() < w.systems[j].GetPriority()
	})

	// 既存のエンティティをシステムに追加
	for _, entity := range w.entities {
		if entity.IsActive() && system.HasRequiredComponents(entity) {
			system.OnEntityAdded(entity)
		}
	}
}

// 更新処理
func (w *World) Update(dt float64) error {
	w.mutex.Lock()
	// 保留中のエンティティの追加
	for _, entity := range w.toAdd {
		w.entities[entity.GetID()] = entity
		for _, system := range w.systems {
			if system.HasRequiredComponents(entity) {
				system.OnEntityAdded(entity)
			}
		}
	}
	w.toAdd = w.toAdd[:0]

	// 保留中のエンティティの削除
	for _, id := range w.toRemove {
		if entity, exists := w.entities[id]; exists {
			for _, system := range w.systems {
				system.OnEntityRemoved(entity)
			}
			delete(w.entities, id)
		}
	}
	w.toRemove = w.toRemove[:0]
	w.mutex.Unlock()

	// システムの更新
	for _, system := range w.systems {
		if err := system.Update(dt); err != nil {
			return err
		}
	}

	return nil
}

// コンポーネントが追加された時の処理
func (w *World) entityComponentAdded(entity *Entity, component Component) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	for _, system := range w.systems {
		// エンティティがすでにロックされているので、HasComponentの代わりに
		// 直接コンポーネントの存在チェックを行う
		hasAllComponents := true
		for _, requiredID := range system.GetRequiredComponents() {
			if _, exists := entity.components[requiredID]; !exists {
				hasAllComponents = false
				break
			}
		}

		if hasAllComponents {
			system.OnEntityAdded(entity)
		}
	}
}

// コンポーネントが削除された時の処理
func (w *World) entityComponentRemoved(entity *Entity, component Component) {
	for _, system := range w.systems {
		if !system.HasRequiredComponents(entity) {
			system.OnEntityRemoved(entity)
		}
	}
}

// エンティティの検索
func (w *World) FindEntitiesByTag(tag string) []*Entity {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	var entities []*Entity
	for _, entity := range w.entities {
		// fmt.Printf("Checking entity %d for tag '%s': active=%v, hasTag=%v\n",
		// 	entity.GetID(), tag, entity.IsActive(), entity.HasTag(tag))
		if entity.IsActive() && entity.HasTag(tag) {
			entities = append(entities, entity)
		}
	}
	return entities
}

// エンティティの取得
func (w *World) GetEntity(id EntityID) *Entity {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.entities[id]
}

// 非アクティブなエンティティを削除
func (w *World) CleanupInactiveEntities() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// 非アクティブなエンティティを削除
	for id, entity := range w.entities {
		if !entity.IsActive() {
			delete(w.entities, id)
		}
	}
}

// アクティブなエンティティの総数を取得
func (w *World) GetTotalEntities() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	count := 0
	for _, entity := range w.entities {
		if entity.IsActive() {
			count++
		}
	}
	return count
}
