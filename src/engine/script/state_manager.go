package script

import (
	"gameengine/src/engine/ecs"
	"sync"
)

// StateManager は各エンティティの状態を管理します
type StateManager struct {
	mutex  sync.RWMutex
	states map[ecs.EntityID]map[string]interface{}
	world  *ecs.World
}

func NewStateManager(world *ecs.World) *StateManager {
	return &StateManager{
		states: make(map[ecs.EntityID]map[string]interface{}),
		world:  world,
	}
}

// エンティティの状態を設定
func (sm *StateManager) SetState(entityID ecs.EntityID, key string, value interface{}) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, exists := sm.states[entityID]; !exists {
		sm.states[entityID] = make(map[string]interface{})
	}
	sm.states[entityID][key] = value
}

// エンティティの状態を取得
func (sm *StateManager) GetState(entityID ecs.EntityID, key string) interface{} {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	if state, exists := sm.states[entityID]; exists {
		return state[key]
	}
	return nil
}

// SetStates は複数の状態を一度に設定します
func (sm *StateManager) SetStates(entityID ecs.EntityID, states map[string]interface{}) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if _, exists := sm.states[entityID]; !exists {
		sm.states[entityID] = make(map[string]interface{})
	}

	for key, value := range states {
		sm.states[entityID][key] = value
	}
}

// GetStates は指定したエンティティの全ての状態を取得します
func (sm *StateManager) GetStates(entityID ecs.EntityID, keys []string) map[string]interface{} {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	result := make(map[string]interface{})
	if state, exists := sm.states[entityID]; exists {
		for _, key := range keys {
			if value, ok := state[key]; ok {
				result[key] = value
			}
		}
	}
	return result
}
