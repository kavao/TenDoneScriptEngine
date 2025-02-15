package state

// ステートインターフェース
type State interface {
	OnEnter(data interface{})
	OnUpdate(dt float64) error
	OnExit()
	GetName() string
}

// 基本的なステート実装
type BaseState struct {
	name string
}

func NewBaseState(name string) *BaseState {
	return &BaseState{
		name: name,
	}
}

func (s *BaseState) OnEnter(data interface{}) {}
func (s *BaseState) OnUpdate(dt float64) error { return nil }
func (s *BaseState) OnExit()                   {}
func (s *BaseState) GetName() string          { return s.name }

// ステート遷移イベント
type TransitionEvent struct {
	FromState string
	ToState   string
	Data      interface{}
} 