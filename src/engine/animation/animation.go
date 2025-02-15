package animation

import (
	"time"
)

// アニメーションインターフェース
type Animation interface {
	Update() error
	IsFinished() bool
	Reset()
	GetValue() float64
}

// 基本的なアニメーション構造
type BaseAnimation struct {
	startTime    time.Time
	duration     time.Duration
	currentValue float64
	startValue   float64
	endValue     float64
	isLoop       bool
}

func NewBaseAnimation(duration time.Duration, start, end float64, loop bool) *BaseAnimation {
	return &BaseAnimation{
		startTime:    time.Now(),
		duration:     duration,
		currentValue: start,
		startValue:   start,
		endValue:     end,
		isLoop:       loop,
	}
}

func (a *BaseAnimation) Update() error {
	elapsed := time.Since(a.startTime)
	progress := float64(elapsed) / float64(a.duration)

	if progress >= 1.0 {
		if a.isLoop {
			a.Reset()
			return nil
		}
		a.currentValue = a.endValue
		return nil
	}

	a.currentValue = a.startValue + (a.endValue-a.startValue)*progress
	return nil
}

func (a *BaseAnimation) IsFinished() bool {
	return !a.isLoop && time.Since(a.startTime) >= a.duration
}

func (a *BaseAnimation) Reset() {
	a.startTime = time.Now()
	a.currentValue = a.startValue
}

func (a *BaseAnimation) GetValue() float64 {
	return a.currentValue
} 