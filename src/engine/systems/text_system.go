package systems

import (
	"gameengine/src/engine/ecs"
	"gameengine/src/engine/ecs/components"
	"gameengine/src/engine/ecs/core"
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
	screen    *ebiten.Image
	font      font.Face
	textCache map[string]*ebiten.Image
}

func NewTextSystem() *TextSystem {
	fontData, err := os.ReadFile("assets/fonts/NotoSansJP-Regular.ttf")
	if err != nil {
		return &TextSystem{
			BaseSystem: ecs.NewBaseSystem(ecs.PriorityRender+1, []core.ComponentID{3}),
			font:       basicfont.Face7x13,
		}
	}

	tt, err := opentype.Parse(fontData)
	if err != nil {
		return &TextSystem{
			BaseSystem: ecs.NewBaseSystem(ecs.PriorityRender+1, []core.ComponentID{3}),
			font:       basicfont.Face7x13,
		}
	}

	const dpi = 72
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return &TextSystem{
			BaseSystem: ecs.NewBaseSystem(ecs.PriorityRender+1, []core.ComponentID{3}),
			font:       basicfont.Face7x13,
		}
	}

	return &TextSystem{
		BaseSystem: ecs.NewBaseSystem(ecs.PriorityRender+1, []core.ComponentID{3}),
		font:       face,
		textCache:  make(map[string]*ebiten.Image),
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

		img, exists := s.textCache[textComp.Text]
		if !exists {
			bounds := text.BoundString(s.font, textComp.Text)
			img = ebiten.NewImage(bounds.Dx(), bounds.Dy())
			text.Draw(img, textComp.Text, s.font, 0, -bounds.Min.Y, color.White)
			s.textCache[textComp.Text] = img
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(textComp.X, textComp.Y)
		s.screen.DrawImage(img, op)
	}
	return nil
}

func (s *TextSystem) SetScreen(screen *ebiten.Image) {
	s.screen = screen
}
