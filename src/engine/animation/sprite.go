package animation

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type SpriteAnimation struct {
	frames       []*ebiten.Image
	frameTime    time.Duration
	currentFrame int
	lastUpdate   time.Time
	isLoop       bool
	isFinished   bool
}

func NewSpriteAnimation(frames []*ebiten.Image, frameTime time.Duration, loop bool) *SpriteAnimation {
	return &SpriteAnimation{
		frames:     frames,
		frameTime:  frameTime,
		lastUpdate: time.Now(),
		isLoop:     loop,
	}
}

func (s *SpriteAnimation) Update() error {
	if s.isFinished {
		return nil
	}

	if time.Since(s.lastUpdate) >= s.frameTime {
		s.currentFrame++
		if s.currentFrame >= len(s.frames) {
			if s.isLoop {
				s.currentFrame = 0
			} else {
				s.currentFrame = len(s.frames) - 1
				s.isFinished = true
			}
		}
		s.lastUpdate = time.Now()
	}
	return nil
}

func (s *SpriteAnimation) GetCurrentFrame() *ebiten.Image {
	return s.frames[s.currentFrame]
}

func (s *SpriteAnimation) IsFinished() bool {
	return s.isFinished
}

func (s *SpriteAnimation) Reset() {
	s.currentFrame = 0
	s.lastUpdate = time.Now()
	s.isFinished = false
} 