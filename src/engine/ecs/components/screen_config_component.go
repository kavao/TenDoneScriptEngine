package components

import "gameengine/src/engine/ecs/core"

type ScreenConfigComponent struct {
	*core.BaseComponent
	Width  int
	Height int
}

// プリセット解像度の定義
var ScreenResolutions = map[string]struct {
	Width  int
	Height int
}{
	"HD":       {1280, 720},  // HD
	"FULL_HD":  {1920, 1080}, // Full HD
	"SD":       {800, 600},   // Standard
	"MOBILE":   {360, 640},   // モバイル縦向き
	"MOBILE_L": {640, 360},   // モバイル横向き
}

func NewScreenConfigComponent() *ScreenConfigComponent {
	return &ScreenConfigComponent{
		BaseComponent: core.NewBaseComponent(4), // Screen Config ID = 4
		Width:         1280,                     // デフォルトはHD
		Height:        720,
	}
}

func (c *ScreenConfigComponent) SetResolution(preset string) {
	if res, exists := ScreenResolutions[preset]; exists {
		c.Width = res.Width
		c.Height = res.Height
	}
}
