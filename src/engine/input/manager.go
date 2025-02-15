package input

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// 入力タイプの定義
type InputType int

const (
	KeyboardInput InputType = iota
	MouseInput
	GamepadInput
)

// 入力アクションの定義
type Action string

const (
	ActionUp     Action = "up"
	ActionDown   Action = "down"
	ActionLeft   Action = "left"
	ActionRight  Action = "right"
	ActionOK     Action = "ok"
	ActionCancel Action = "cancel"
	ActionMenu   Action = "menu"
)

// 入力マッピング
type InputBinding struct {
	Type     InputType
	Key      ebiten.Key
	Button   ebiten.GamepadButton
	MouseBtn ebiten.MouseButton
}

// 入力マネージャー
type InputManager struct {
	mutex     sync.RWMutex
	bindings  map[Action][]InputBinding
	states    map[Action]bool
	previous  map[Action]bool
	callbacks map[Action][]func()
}

func NewInputManager() *InputManager {
	im := &InputManager{
		bindings:  make(map[Action][]InputBinding),
		states:    make(map[Action]bool),
		previous:  make(map[Action]bool),
		callbacks: make(map[Action][]func()),
	}

	// デフォルトのキーバインドを設定
	im.SetDefaultBindings()
	return im
}

func (im *InputManager) SetDefaultBindings() {
	// 方向キー
	im.BindAction(ActionUp, InputBinding{KeyboardInput, ebiten.KeyUp, 0, 0})
	im.BindAction(ActionDown, InputBinding{KeyboardInput, ebiten.KeyDown, 0, 0})
	im.BindAction(ActionLeft, InputBinding{KeyboardInput, ebiten.KeyLeft, 0, 0})
	im.BindAction(ActionRight, InputBinding{KeyboardInput, ebiten.KeyRight, 0, 0})

	// アクションキー
	im.BindAction(ActionOK, InputBinding{KeyboardInput, ebiten.KeyZ, 0, 0})
	im.BindAction(ActionCancel, InputBinding{KeyboardInput, ebiten.KeyX, 0, 0})
	im.BindAction(ActionMenu, InputBinding{KeyboardInput, ebiten.KeyEscape, 0, 0})

	// マウス
	im.BindAction(ActionOK, InputBinding{MouseInput, 0, 0, ebiten.MouseButtonLeft})
	im.BindAction(ActionCancel, InputBinding{MouseInput, 0, 0, ebiten.MouseButtonRight})
}

func (im *InputManager) BindAction(action Action, binding InputBinding) {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	if _, exists := im.bindings[action]; !exists {
		im.bindings[action] = make([]InputBinding, 0)
	}
	im.bindings[action] = append(im.bindings[action], binding)
}

func (im *InputManager) Update() error {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	// 前のフレームの状態を保存
	for action := range im.states {
		im.previous[action] = im.states[action]
	}

	// 現在の入力状態を更新
	for action, bindings := range im.bindings {
		im.states[action] = false
		for _, binding := range bindings {
			if im.isBindingActive(binding) {
				im.states[action] = true
				break
			}
		}
	}

	// コールバックの実行
	for action, callbacks := range im.callbacks {
		if im.IsJustPressed(action) {
			for _, callback := range callbacks {
				callback()
			}
		}
	}

	return nil
}

func (im *InputManager) isBindingActive(binding InputBinding) bool {
	switch binding.Type {
	case KeyboardInput:
		return ebiten.IsKeyPressed(binding.Key)
	case MouseInput:
		return ebiten.IsMouseButtonPressed(binding.MouseBtn)
	case GamepadInput:
		// TODO: ゲームパッドの実装
		return false
	default:
		return false
	}
}

func (im *InputManager) IsPressed(action Action) bool {
	im.mutex.RLock()
	defer im.mutex.RUnlock()
	return im.states[action]
}

func (im *InputManager) IsJustPressed(action Action) bool {
	im.mutex.RLock()
	defer im.mutex.RUnlock()
	return im.states[action] && !im.previous[action]
}

func (im *InputManager) IsJustReleased(action Action) bool {
	im.mutex.RLock()
	defer im.mutex.RUnlock()
	return !im.states[action] && im.previous[action]
}

// コールバックの登録
func (im *InputManager) OnAction(action Action, callback func()) {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	if _, exists := im.callbacks[action]; !exists {
		im.callbacks[action] = make([]func(), 0)
	}
	im.callbacks[action] = append(im.callbacks[action], callback)
}
