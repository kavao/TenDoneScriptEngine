package input

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type MouseState struct {
	X, Y          int
	PrevX, PrevY  int
	ScrollX       float64
	ScrollY       float64
	DragStartX    int
	DragStartY    int
	IsDragging    bool
	DragThreshold int
}

func NewMouseState() *MouseState {
	return &MouseState{
		DragThreshold: 5, // ドラッグ開始の閾値（ピクセル）
	}
}

func (m *MouseState) Update() {
	m.PrevX = m.X
	m.PrevY = m.Y
	m.X, m.Y = ebiten.CursorPosition()
	m.ScrollX, m.ScrollY = ebiten.Wheel()

	// ドラッグ処理
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if !m.IsDragging {
			dx := m.X - m.DragStartX
			dy := m.Y - m.DragStartY
			if dx*dx+dy*dy >= m.DragThreshold*m.DragThreshold {
				m.IsDragging = true
			}
		}
	} else {
		m.IsDragging = false
		m.DragStartX = m.X
		m.DragStartY = m.Y
	}
}

func (m *MouseState) GetDragDelta() (dx, dy int) {
	if m.IsDragging {
		return m.X - m.PrevX, m.Y - m.PrevY
	}
	return 0, 0
} 