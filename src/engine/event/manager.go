package event

import (
	"sync"
	"time"
)

// イベントマネージャー
type EventManager struct {
	mutex       sync.RWMutex
	eventBus    *EventBus
	eventPool   sync.Pool
	deferredMap map[string][]time.Time
}

func NewEventManager() *EventManager {
	return &EventManager{
		eventBus: NewEventBus(),
		eventPool: sync.Pool{
			New: func() interface{} {
				return &BaseEvent{}
			},
		},
		deferredMap: make(map[string][]time.Time),
	}
}

// イベントの発行
func (m *EventManager) Emit(eventType string, data interface{}, priority Priority) error {
	event := m.eventPool.Get().(*BaseEvent)
	event.Type = eventType
	event.Timestamp = time.Now()
	event.Priority = priority
	event.Data = data

	err := m.eventBus.Publish(event)
	m.eventPool.Put(event)
	return err
}

// 遅延イベントの発行
func (m *EventManager) EmitDeferred(eventType string, data interface{}, priority Priority, delay time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	executeTime := time.Now().Add(delay)
	if _, exists := m.deferredMap[eventType]; !exists {
		m.deferredMap[eventType] = make([]time.Time, 0)
	}
	m.deferredMap[eventType] = append(m.deferredMap[eventType], executeTime)

	go func() {
		time.Sleep(delay)
		m.mutex.Lock()
		// 実行時刻のリストから削除
		times := m.deferredMap[eventType]
		for i, t := range times {
			if t.Equal(executeTime) {
				m.deferredMap[eventType] = append(times[:i], times[i+1:]...)
				break
			}
		}
		m.mutex.Unlock()

		m.Emit(eventType, data, priority)
	}()
}

// イベントハンドラーの登録
func (m *EventManager) On(eventType string, handler EventHandler, priority Priority) {
	m.eventBus.Subscribe(eventType, handler, priority)
}

// イベントハンドラーの削除
func (m *EventManager) Off(eventType string, handler EventHandler) {
	m.eventBus.Unsubscribe(eventType, handler)
}

// 更新処理
func (m *EventManager) Update() error {
	return m.eventBus.ProcessQueue()
}

// 遅延イベントの取得
func (m *EventManager) GetDeferredEvents(eventType string) []time.Time {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if times, exists := m.deferredMap[eventType]; exists {
		result := make([]time.Time, len(times))
		copy(result, times)
		return result
	}
	return nil
} 