package scene

import (
	"gameengine/src/engine/input"
	"gameengine/src/engine/ui"
)

// タイトルシーン
type TitleScene struct {
	*BaseScene
	menuWindow *ui.MenuWindow
}

func NewTitleScene(uiMgr *ui.UIManager, inputMgr *input.InputManager) *TitleScene {
	return &TitleScene{
		BaseScene: NewBaseScene(uiMgr, inputMgr),
	}
}

func (s *TitleScene) Init() error {
	if err := s.BaseScene.Init(); err != nil {
		return err
	}

	// メニューウィンドウの初期化
	s.menuWindow = ui.NewMenuWindow(nil) // フォントは適切に設定する必要があります
	s.menuWindow.SetPosition(100, 100)
	s.menuWindow.SetSize(200, 300)

	// メニュー項目の追加
	s.menuWindow.AddItem("ゲーム開始", true, s.onStartGame)
	s.menuWindow.AddItem("設定", true, s.onOpenSettings)
	s.menuWindow.AddItem("終了", true, s.onExit)

	s.uiManager.AddComponent("title_menu", s.menuWindow)

	return nil
}

func (s *TitleScene) onStartGame() {
	// ゲーム開始処理
}

func (s *TitleScene) onOpenSettings() {
	// 設定画面を開く処理
}

func (s *TitleScene) onExit() {
	// ゲーム終了処理
} 