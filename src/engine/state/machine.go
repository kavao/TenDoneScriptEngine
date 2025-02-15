package state

import (
	"fmt"
	"sync"
)

// ステートマシン
type StateMachine struct {
	mutex          sync.RWMutex
	states         map[string]State
	currentState   State
	previousState  State
	globalState    State
	transitions    map[string]map[string]bool
	onTransition   func(TransitionEvent)
	stateData      interface{}
}

func NewStateMachine() *StateMachine {
	return &StateMachine{
		states:      make(map[string]State),
		transitions: make(map[string]map[string]bool),
	}
}

// ステートの追加
func (m *StateMachine) AddState(state State) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.states[state.GetName()] = state
}

// 遷移の追加
func (m *StateMachine) AddTransition(from, to string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// ステートの存在確認
	if _, exists := m.states[from]; !exists {
		return fmt.Errorf("state not found: %s", from)
	}
	if _, exists := m.states[to]; !exists {
		return fmt.Errorf("state not found: %s", to)
	}

	// 遷移の登録
	if _, exists := m.transitions[from]; !exists {
		m.transitions[from] = make(map[string]bool)
	}
	m.transitions[from][to] = true
	return nil
}

// 初期ステートの設定
func (m *StateMachine) SetInitialState(name string, data interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	state, exists := m.states[name]
	if !exists {
		return fmt.Errorf("state not found: %s", name)
	}

	m.currentState = state
	m.stateData = data
	state.OnEnter(data)
	return nil
}

// グローバルステートの設定
func (m *StateMachine) SetGlobalState(state State) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.globalState = state
}

// ステート遷移
func (m *StateMachine) ChangeState(name string, data interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 遷移可能性のチェック
	if !m.canTransition(m.currentState.GetName(), name) {
		return fmt.Errorf("invalid transition: %s -> %s", m.currentState.GetName(), name)
	}

	newState, exists := m.states[name]
	if !exists {
		return fmt.Errorf("state not found: %s", name)
	}

	// 遷移イベントの作成
	event := TransitionEvent{
		FromState: m.currentState.GetName(),
		ToState:   name,
		Data:      data,
	}

	// 現在のステートを終了
	m.currentState.OnExit()

	// ステートの切り替え
	m.previousState = m.currentState
	m.currentState = newState
	m.stateData = data

	// 新しいステートを開始
	m.currentState.OnEnter(data)

	// 遷移コールバックの呼び出し
	if m.onTransition != nil {
		m.onTransition(event)
	}

	return nil
}

// 更新処理
func (m *StateMachine) Update(dt float64) error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// グローバルステートの更新
	if m.globalState != nil {
		if err := m.globalState.OnUpdate(dt); err != nil {
			return err
		}
	}

	// 現在のステートの更新
	if m.currentState != nil {
		if err := m.currentState.OnUpdate(dt); err != nil {
			return err
		}
	}

	return nil
}

// 遷移可能性のチェック
func (m *StateMachine) canTransition(from, to string) bool {
	transitions, exists := m.transitions[from]
	if !exists {
		return false
	}
	return transitions[to]
}

// 遷移コールバックの設定
func (m *StateMachine) SetTransitionCallback(callback func(TransitionEvent)) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.onTransition = callback
}

// 現在のステート名を取得
func (m *StateMachine) GetCurrentState() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.currentState == nil {
		return ""
	}
	return m.currentState.GetName()
}

// 前のステート名を取得
func (m *StateMachine) GetPreviousState() string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if m.previousState == nil {
		return ""
	}
	return m.previousState.GetName()
} 