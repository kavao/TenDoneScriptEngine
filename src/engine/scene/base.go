package scene

import (
	"github.com/hajimehoshi/ebiten/v2"
	"gameengine/src/engine/input"
	"gameengine/src/engine/ui"
)

// 基本的なシーン実装
type BaseScene struct {
	initialized bool
	uiManager   *ui.UIManager
	inputMgr    *input.InputManager
}

func NewBaseScene(uiMgr *ui.UIManager, inputMgr *input.InputManager) *BaseScene {
	return &BaseScene{
		uiManager: uiMgr,
		inputMgr:  inputMgr,
	}
}

func (s *BaseScene) Init() error {
	s.initialized = true
	return nil
}

func (s *BaseScene) Update() error {
	return s.uiManager.Update()
}

func (s *BaseScene) Draw(screen *ebiten.Image) {
	s.uiManager.Draw(screen)
}

func (s *BaseScene) Finalize() error {
	s.initialized = false
	return nil
}

func (s *BaseScene) OnEnter(params map[string]interface{}) error {
	return nil
}

func (s *BaseScene) OnExit() error {
	return nil
}

func (s *BaseScene) IsReady() bool {
	return s.initialized
} 