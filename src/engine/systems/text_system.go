package systems

import (
	"gameengine/src/engine/components"
	"gameengine/src/engine/ecs"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
)

type TextSystem struct {
	*ecs.BaseSystem
	screen *ebiten.Image
	font   font.Face
}

func NewTextSystem() *TextSystem {
	// 日本語フォントの読み込み
	fontData, err := os.ReadFile("assets/fonts/NotoSansJP-Regular.ttf")
	if err != nil {
		// エラー時は基本フォントにフォールバック
		return &TextSystem{
			BaseSystem: ecs.NewBaseSystem(ecs.PriorityRender+1, []ecs.ComponentID{3}),
			font:       basicfont.Face7x13,
		}
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		return &TextSystem{
			BaseSystem: ecs.NewBaseSystem(ecs.PriorityRender+1, []ecs.ComponentID{3}),
			font:       basicfont.Face7x13,
		}
	}

	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return &TextSystem{
			BaseSystem: ecs.NewBaseSystem(ecs.PriorityRender+1, []ecs.ComponentID{3}),
			font:       basicfont.Face7x13,
		}
	}

	return &TextSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.PriorityRender+1, []ecs.ComponentID{3}),
		font:       face,
	}
}

func (s *TextSystem) Update(dt float64) error {
	if s.screen == nil {
		return nil
	}

	for _, entity := range s.BaseSystem.Entities() {
		textComp := entity.GetComponent(3).(*components.TextComponent)
		if !textComp.Visible {
			continue
		}

		text.Draw(s.screen, textComp.Text, s.font,
			int(textComp.X), int(textComp.Y),
			color.White)
	}
	return nil
}

func (s *TextSystem) SetScreen(screen *ebiten.Image) {
	s.screen = screen
}
