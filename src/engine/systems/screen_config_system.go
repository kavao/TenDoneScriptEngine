package systems

import (
	"fmt"
	"gameengine/src/engine/ecs"
	"gameengine/src/engine/ecs/components"
	"gameengine/src/engine/ecs/core"

	"github.com/hajimehoshi/ebiten/v2"
)

type ScreenConfigSystem struct {
	*ecs.BaseSystem
	game interface {
		SetScreenSize(width, height int)
	}
	currentWidth  int
	currentHeight int
}

func NewScreenConfigSystem(game interface{ SetScreenSize(width, height int) }) *ScreenConfigSystem {
	return &ScreenConfigSystem{
		BaseSystem:    ecs.NewBaseSystem(ecs.PriorityUpdate, []core.ComponentID{4}),
		game:          game,
		currentWidth:  1280,
		currentHeight: 720,
	}
}

func (s *ScreenConfigSystem) Update(dt float64) error {
	for _, entity := range s.BaseSystem.Entities() {
		config := entity.GetComponent(4).(*components.ScreenConfigComponent)
		if s.currentWidth != config.Width || s.currentHeight != config.Height {
			fmt.Printf("Screen config changed: %dx%d\n", config.Width, config.Height)

			// 一時的にリサイズモードを無効化
			ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

			// ウィンドウサイズを設定
			ebiten.SetWindowSize(config.Width, config.Height)

			// ゲーム内部の解像度を設定
			s.game.SetScreenSize(config.Width, config.Height)

			// リサイズモードを再度有効化
			ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

			s.currentWidth = config.Width
			s.currentHeight = config.Height

			fmt.Printf("Window size set to: %dx%d\n", config.Width, config.Height)
		}
	}
	return nil
}
