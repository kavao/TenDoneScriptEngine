package event

import (
	"time"
)

// イベントの優先度
type Priority int

const (
	PriorityLow Priority = iota
	PriorityNormal
	PriorityHigh
	PriorityCritical
)

// イベントインターフェース
type Event interface {
	GetType() string
	GetTimestamp() time.Time
	GetPriority() Priority
	GetData() interface{}
}

// 基本的なイベント実装
type BaseEvent struct {
	Type      string
	Timestamp time.Time
	Priority  Priority
	Data      interface{}
}

func NewEvent(eventType string, data interface{}, priority Priority) *BaseEvent {
	return &BaseEvent{
		Type:      eventType,
		Timestamp: time.Now(),
		Priority:  priority,
		Data:      data,
	}
}

func (e *BaseEvent) GetType() string      { return e.Type }
func (e *BaseEvent) GetTimestamp() time.Time { return e.Timestamp }
func (e *BaseEvent) GetPriority() Priority { return e.Priority }
func (e *BaseEvent) GetData() interface{}  { return e.Data } 