package collision

import (
	"sync"
)

// 衝突オブジェクト
type CollisionObject struct {
	ID       string
	Shape    Shape
	Layer    uint32
	Mask     uint32
	UserData interface{}
}

// 衝突イベント
type CollisionEvent struct {
	ObjectA   *CollisionObject
	ObjectB   *CollisionObject
	Collision bool
}

// 衝突判定マネージャー
type CollisionManager struct {
	mutex    sync.RWMutex
	objects  map[string]*CollisionObject
	events   []CollisionEvent
	handlers map[string]func(CollisionEvent)
}

func NewCollisionManager() *CollisionManager {
	return &CollisionManager{
		objects:  make(map[string]*CollisionObject),
		events:   make([]CollisionEvent, 0),
		handlers: make(map[string]func(CollisionEvent)),
	}
}

// オブジェクトの追加
func (m *CollisionManager) AddObject(obj *CollisionObject) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.objects[obj.ID] = obj
}

// オブジェクトの削除
func (m *CollisionManager) RemoveObject(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.objects, id)
}

// 衝突判定の更新
func (m *CollisionManager) Update() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.events = m.events[:0] // イベントリストをクリア

	// 全オブジェクト間の衝突判定
	for idA, objA := range m.objects {
		for idB, objB := range m.objects {
			if idA >= idB {
				continue // 重複判定を避ける
			}

			// レイヤーマスクのチェック
			if objA.Layer&objB.Mask == 0 && objB.Layer&objA.Mask == 0 {
				continue
			}

			collision := CheckCollision(objA.Shape, objB.Shape)
			event := CollisionEvent{
				ObjectA:   objA,
				ObjectB:   objB,
				Collision: collision,
			}
			m.events = append(m.events, event)

			// ハンドラーの呼び出し
			if handler, exists := m.handlers[idA]; exists && collision {
				handler(event)
			}
			if handler, exists := m.handlers[idB]; exists && collision {
				handler(event)
			}
		}
	}
}

// 衝突ハンドラーの登録
func (m *CollisionManager) SetCollisionHandler(id string, handler func(CollisionEvent)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.handlers[id] = handler
} 