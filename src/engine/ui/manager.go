package ui

import (
	"sort"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

type UIManager struct {
	mutex      sync.RWMutex
	components map[string]Component
	zOrder     []string
}

func NewUIManager() *UIManager {
	return &UIManager{
		components: make(map[string]Component),
	}
}

func (m *UIManager) AddComponent(name string, component Component) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.components[name] = component
	m.updateZOrder()
}

func (m *UIManager) GetComponent(name string) Component {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return m.components[name]
}

func (m *UIManager) RemoveComponent(name string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.components, name)
	m.updateZOrder()
}

func (m *UIManager) Draw(screen *ebiten.Image) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, name := range m.zOrder {
		if component := m.components[name]; component != nil {
			component.Draw(screen)
		}
	}
}

func (m *UIManager) Update() error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	for _, component := range m.components {
		if err := component.Update(); err != nil {
			return err
		}
	}
	return nil
}

func (m *UIManager) updateZOrder() {
	type componentWithZ struct {
		name    string
		zIndex  int
		visible bool
	}

	components := make([]componentWithZ, 0, len(m.components))
	for name, comp := range m.components {
		if base, ok := comp.(*BaseComponent); ok {
			components = append(components, componentWithZ{
				name:    name,
				zIndex:  base.ZIndex,
				visible: comp.IsVisible(),
			})
		}
	}

	sort.Slice(components, func(i, j int) bool {
		return components[i].zIndex < components[j].zIndex
	})

	m.zOrder = make([]string, len(components))
	for i, comp := range components {
		m.zOrder[i] = comp.name
	}
} 