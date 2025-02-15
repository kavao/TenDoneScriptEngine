package systems

import (
	"gameengine/src/engine/ecs"
)

type InputSystem struct {
	*ecs.BaseSystem
}

func NewInputSystem() *InputSystem {
	return &InputSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.PriorityUpdate, []ecs.ComponentID{}), // 必要なコンポーネントなし
	}
}

func (s *InputSystem) Update(dt float64) error {
	return nil
}
