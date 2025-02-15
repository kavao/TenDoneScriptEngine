package ui

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type MessageWindow struct {
	BaseComponent
	Text            string
	Font            font.Face
	BackgroundColor color.Color
	TextColor       color.Color
	Padding         float64
	Lines           []string
	MaxLines        int
	LineHeight      float64
}

func NewMessageWindow(font font.Face) *MessageWindow {
	return &MessageWindow{
		BaseComponent: BaseComponent{
			Visible: true,
			ZIndex:  100,
		},
		Font:            font,
		BackgroundColor: color.RGBA{0, 0, 0, 200},
		TextColor:       color.White,
		Padding:         10,
		MaxLines:        4,
		LineHeight:      24,
	}
}

func (w *MessageWindow) SetText(text string) {
	w.Text = text
	w.Lines = w.wrapText(text)
}

func (w *MessageWindow) AppendText(text string) {
	w.Text += text
	w.Lines = w.wrapText(w.Text)
	if len(w.Lines) > w.MaxLines {
		w.Lines = w.Lines[len(w.Lines)-w.MaxLines:]
	}
}

func (w *MessageWindow) Draw(screen *ebiten.Image) {
	if !w.Visible {
		return
	}

	// 背景描画
	bounds := w.GetBounds()
	windowImage := ebiten.NewImage(bounds.Dx(), bounds.Dy())
	windowImage.Fill(w.BackgroundColor)

	// テキスト描画
	for i, line := range w.Lines {
		y := w.Padding + float64(i)*w.LineHeight
		text.Draw(windowImage, line, w.Font, int(w.Padding), int(y+w.LineHeight), w.TextColor)
	}

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(w.X, w.Y)
	screen.DrawImage(windowImage, op)
}

func (w *MessageWindow) wrapText(text string) []string {
	// 簡単な行折り返し処理
	maxWidth := int(w.Width - w.Padding*2)
	var lines []string
	words := strings.Split(text, " ")
	currentLine := ""

	for _, word := range words {
		candidate := currentLine
		if currentLine != "" {
			candidate += " "
		}
		candidate += word

		bounds, _ := font.BoundString(w.Font, candidate)
		width := bounds.Max.X - bounds.Min.X
		if width > fixed.I(maxWidth) && currentLine != "" {
			lines = append(lines, currentLine)
			currentLine = word
		} else {
			currentLine = candidate
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return lines
}
