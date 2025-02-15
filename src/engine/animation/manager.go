package animation

import (
	"sync"
)

type AnimationManager struct {
	mutex      sync.RWMutex
	animations map[string]Animation
	sequences  map[string]*AnimationSequence
}

func NewAnimationManager() *AnimationManager {
	return &AnimationManager{
		animations: make(map[string]Animation),
		sequences:  make(map[string]*AnimationSequence),
	}
}

func (m *AnimationManager) Update() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// アニメーションの更新
	for id, anim := range m.animations {
		if anim.IsFinished() {
			delete(m.animations, id)
			continue
		}
		if err := anim.Update(); err != nil {
			return err
		}
	}

	// シーケンスの更新
	for id, seq := range m.sequences {
		if seq.IsFinished() {
			delete(m.sequences, id)
			continue
		}
		if err := seq.Update(); err != nil {
			return err
		}
	}

	return nil
}

func (m *AnimationManager) AddAnimation(id string, anim Animation) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.animations[id] = anim
}

func (m *AnimationManager) AddSequence(id string, seq *AnimationSequence) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.sequences[id] = seq
}

func (m *AnimationManager) StopAnimation(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.animations, id)
}

func (m *AnimationManager) StopAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.animations = make(map[string]Animation)
	m.sequences = make(map[string]*AnimationSequence)
}
