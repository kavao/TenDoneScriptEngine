package systems

import (
	"gameengine/src/engine/ecs"
	"gameengine/src/engine/ecs/core"
)

type InputSystem struct {
	*ecs.BaseSystem
}

func NewInputSystem() *InputSystem {
	return &InputSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.PriorityUpdate, []core.ComponentID{}), // 必要なコンポーネントなし
	}
}

func (s *InputSystem) Update(dt float64) error {
	return nil
}
