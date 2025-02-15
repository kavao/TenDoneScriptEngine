package systems

import (
	"gameengine/src/engine/components"
	"gameengine/src/engine/ecs"
)

type PhysicsSystem struct {
	*ecs.BaseSystem
}

func NewPhysicsSystem() *PhysicsSystem {
	return &PhysicsSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.PriorityPhysics, []ecs.ComponentID{1, 5}), // Transform と Physics
	}
}

func (s *PhysicsSystem) Update(dt float64) error {
	for _, entity := range s.BaseSystem.Entities() {
		physics := entity.GetComponent(5).(*components.PhysicsComponent)
		transform := entity.GetComponent(1).(*components.TransformComponent)

		// 速度に重力を加える
		physics.VelocityY += physics.Gravity * dt

		// 位置の更新
		transform.X += physics.VelocityX * dt
		transform.Y += physics.VelocityY * dt

		// 画面外に出たエンティティの削除
		if transform.Y < -50 || transform.Y > 650 {
			// DestroyEntityの代わりにDeactivateを使用
			entity.Deactivate()
			// タグを削除して、カウントから除外
			entity.RemoveTag("bullet")
		}
	}
	return nil
} 