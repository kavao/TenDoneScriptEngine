package ecs

import (
	"gameengine/src/engine/ecs/components"
	"gameengine/src/engine/ecs/core"
)

func NewWorld() *core.World {
	world := core.NewWorld()

	// デフォルトの画面設定エンティティを作成
	configEntity := world.CreateEntity()
	configComponent := components.NewScreenConfigComponent()
	configEntity.AddComponent(configComponent)

	return world
}
