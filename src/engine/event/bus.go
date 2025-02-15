package event

import (
	"fmt"
	"sync"
)

// イベントハンドラー
type EventHandler func(Event) error

// イベントバス
type EventBus struct {
	mutex     sync.RWMutex
	handlers  map[string][]handlerEntry
	queue     []Event
	queueLock sync.Mutex
}

type handlerEntry struct {
	handler  EventHandler
	priority Priority
}

func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[string][]handlerEntry),
		queue:    make([]Event, 0),
	}
}

// イベントハンドラーの登録
func (b *EventBus) Subscribe(eventType string, handler EventHandler, priority Priority) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	entry := handlerEntry{
		handler:  handler,
		priority: priority,
	}

	if _, exists := b.handlers[eventType]; !exists {
		b.handlers[eventType] = make([]handlerEntry, 0)
	}

	// 優先度に基づいてハンドラーを挿入
	handlers := b.handlers[eventType]
	insertIdx := len(handlers)
	for i, h := range handlers {
		if priority > h.priority {
			insertIdx = i
			break
		}
	}

	// スライスへの挿入
	handlers = append(handlers, handlerEntry{})
	copy(handlers[insertIdx+1:], handlers[insertIdx:])
	handlers[insertIdx] = entry
	b.handlers[eventType] = handlers
}

// イベントハンドラーの削除
func (b *EventBus) Unsubscribe(eventType string, handler EventHandler) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if handlers, exists := b.handlers[eventType]; exists {
		for i, h := range handlers {
			if fmt.Sprintf("%p", h.handler) == fmt.Sprintf("%p", handler) {
				b.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// イベントの発行（同期）
func (b *EventBus) Publish(event Event) error {
	b.mutex.RLock()
	handlers, exists := b.handlers[event.GetType()]
	b.mutex.RUnlock()

	if !exists {
		return nil
	}

	for _, h := range handlers {
		if err := h.handler(event); err != nil {
			return err
		}
	}

	return nil
}

// イベントの発行（非同期）
func (b *EventBus) PublishAsync(event Event) {
	b.queueLock.Lock()
	b.queue = append(b.queue, event)
	b.queueLock.Unlock()
}

// キューの処理
func (b *EventBus) ProcessQueue() error {
	b.queueLock.Lock()
	queue := b.queue
	b.queue = make([]Event, 0)
	b.queueLock.Unlock()

	for _, event := range queue {
		if err := b.Publish(event); err != nil {
			return err
		}
	}

	return nil
} 