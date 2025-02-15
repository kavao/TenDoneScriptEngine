package animation

type AnimationSequence struct {
	animations []Animation
	current    int
	onComplete func()
}

func NewAnimationSequence(animations []Animation) *AnimationSequence {
	return &AnimationSequence{
		animations: animations,
		current:    0,
	}
}

func (s *AnimationSequence) Update() error {
	if s.IsFinished() {
		if s.onComplete != nil {
			s.onComplete()
		}
		return nil
	}

	currentAnim := s.animations[s.current]
	if err := currentAnim.Update(); err != nil {
		return err
	}

	if currentAnim.IsFinished() {
		s.current++
	}

	return nil
}

func (s *AnimationSequence) IsFinished() bool {
	return s.current >= len(s.animations)
}

func (s *AnimationSequence) Reset() {
	s.current = 0
	for _, anim := range s.animations {
		anim.Reset()
	}
}

func (s *AnimationSequence) OnComplete(callback func()) {
	s.onComplete = callback
}
