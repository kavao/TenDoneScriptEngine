package particle

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// パーティクルマネージャー
type ParticleManager struct {
	mutex    sync.RWMutex
	emitters map[string]*Emitter
}

func NewParticleManager() *ParticleManager {
	return &ParticleManager{
		emitters: make(map[string]*Emitter),
	}
}

func (m *ParticleManager) AddEmitter(id string, emitter *Emitter) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.emitters[id] = emitter
}

func (m *ParticleManager) RemoveEmitter(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.emitters, id)
}

func (m *ParticleManager) GetEmitter(id string) *Emitter {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.emitters[id]
}

func (m *ParticleManager) Update(dt float64) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, emitter := range m.emitters {
		emitter.Update(dt)
	}
}

func (m *ParticleManager) Draw(screen *ebiten.Image) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, emitter := range m.emitters {
		emitter.Draw(screen)
	}
} 