package animation

import "math"

type EasingFunc func(float64) float64

var (
	Linear EasingFunc = func(t float64) float64 {
		return t
	}

	EaseInQuad EasingFunc = func(t float64) float64 {
		return t * t
	}

	EaseOutQuad EasingFunc = func(t float64) float64 {
		return -t * (t - 2)
	}

	EaseInOutQuad EasingFunc = func(t float64) float64 {
		t *= 2
		if t < 1 {
			return 0.5 * t * t
		}
		t--
		return -0.5 * (t*(t-2) - 1)
	}

	EaseInCubic EasingFunc = func(t float64) float64 {
		return t * t * t
	}

	EaseOutCubic EasingFunc = func(t float64) float64 {
		t--
		return t*t*t + 1
	}

	EaseInElastic EasingFunc = func(t float64) float64 {
		if t == 0 {
			return 0
		}
		if t == 1 {
			return 1
		}
		p := 0.3
		s := p / 4
		t--
		return -(math.Pow(2, 10*t) * math.Sin((t-s)*(2*math.Pi)/p))
	}

	EaseOutElastic EasingFunc = func(t float64) float64 {
		if t == 0 {
			return 0
		}
		if t == 1 {
			return 1
		}
		p := 0.3
		s := p / 4
		return math.Pow(2, -10*t)*math.Sin((t-s)*(2*math.Pi)/p) + 1
	}
) 