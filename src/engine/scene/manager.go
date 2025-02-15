package scene

import (
	"fmt"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"gameengine/src/engine/input"
)

// シーンインターフェース
type Scene interface {
	Init() error
	Update() error
	Draw(screen *ebiten.Image)
	Finalize() error
	OnEnter(params map[string]interface{}) error
	OnExit() error
	IsReady() bool
}

// シーン遷移タイプ
type TransitionType int

const (
	TransitionNone TransitionType = iota
	TransitionPush                // シーンをスタックに追加
	TransitionPop                 // 現在のシーンを削除
	TransitionReplace            // 現在のシーンを置き換え
)

// シーン管理
type SceneManager struct {
	mutex       sync.RWMutex
	scenes      []Scene
	transitions []sceneTransition
	inputMgr    *input.InputManager
}

// シーン遷移情報
type sceneTransition struct {
	transitionType TransitionType
	nextScene      Scene
	params         map[string]interface{}
}

func NewSceneManager(inputMgr *input.InputManager) *SceneManager {
	return &SceneManager{
		scenes:      make([]Scene, 0),
		transitions: make([]sceneTransition, 0),
		inputMgr:    inputMgr,
	}
}

func (m *SceneManager) Push(scene Scene, params map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.transitions = append(m.transitions, sceneTransition{
		transitionType: TransitionPush,
		nextScene:      scene,
		params:         params,
	})
	return nil
}

func (m *SceneManager) Pop() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.scenes) == 0 {
		return fmt.Errorf("no scene to pop")
	}

	m.transitions = append(m.transitions, sceneTransition{
		transitionType: TransitionPop,
	})
	return nil
}

func (m *SceneManager) Replace(scene Scene, params map[string]interface{}) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.transitions = append(m.transitions, sceneTransition{
		transitionType: TransitionReplace,
		nextScene:      scene,
		params:         params,
	})
	return nil
}

func (m *SceneManager) Update() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// シーン遷移の処理
	if err := m.processTransitions(); err != nil {
		return err
	}

	// アクティブなシーンの更新
	if len(m.scenes) > 0 {
		currentScene := m.scenes[len(m.scenes)-1]
		if err := currentScene.Update(); err != nil {
			return err
		}
	}

	return nil
}

func (m *SceneManager) Draw(screen *ebiten.Image) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// 下のシーンから順に描画
	for i := 0; i < len(m.scenes); i++ {
		m.scenes[i].Draw(screen)
	}
}

func (m *SceneManager) processTransitions() error {
	for _, t := range m.transitions {
		switch t.transitionType {
		case TransitionPush:
			if err := t.nextScene.Init(); err != nil {
				return err
			}
			if err := t.nextScene.OnEnter(t.params); err != nil {
				return err
			}
			m.scenes = append(m.scenes, t.nextScene)

		case TransitionPop:
			if len(m.scenes) > 0 {
				current := m.scenes[len(m.scenes)-1]
				if err := current.OnExit(); err != nil {
					return err
				}
				if err := current.Finalize(); err != nil {
					return err
				}
				m.scenes = m.scenes[:len(m.scenes)-1]
			}

		case TransitionReplace:
			if len(m.scenes) > 0 {
				current := m.scenes[len(m.scenes)-1]
				if err := current.OnExit(); err != nil {
					return err
				}
				if err := current.Finalize(); err != nil {
					return err
				}
				m.scenes = m.scenes[:len(m.scenes)-1]
			}
			if err := t.nextScene.Init(); err != nil {
				return err
			}
			if err := t.nextScene.OnEnter(t.params); err != nil {
				return err
			}
			m.scenes = append(m.scenes, t.nextScene)
		}
	}

	// 遷移をクリア
	m.transitions = m.transitions[:0]
	return nil
} 