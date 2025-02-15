package save

import (
	"fmt"
	"go.starlark.net/starlark"
)

// Starlark用のグローバル変数アクセサ
type StateAccess struct {
	manager *SaveManager
}

func NewStateAccess(manager *SaveManager) *StateAccess {
	return &StateAccess{manager: manager}
}

// Starlarkインターフェースの実装
func (s *StateAccess) String() string {
	return "GameState"
}

func (s *StateAccess) Type() string {
	return "GameState"
}

func (s *StateAccess) Freeze() {}

func (s *StateAccess) Truth() starlark.Bool {
	return starlark.True
}

func (s *StateAccess) Hash() (uint32, error) {
	return 0, fmt.Errorf("unhashable type: GameState")
}

// Starlarkからの変数アクセス
func (s *StateAccess) Attr(name string) (starlark.Value, error) {
	if val, exists := s.manager.state.GetVariable(name); exists {
		return convertToStarlarkValue(val), nil
	}
	return starlark.None, nil
}

// Starlarkからの変数設定
func (s *StateAccess) SetAttr(name string, value starlark.Value) error {
	s.manager.SetStateFromStarlark(name, value)
	return nil
}

// 属性一覧の取得
func (s *StateAccess) AttrNames() []string {
	s.manager.state.mutex.RLock()
	defer s.manager.state.mutex.RUnlock()

	names := make([]string, 0, len(s.manager.state.Variables))
	for name := range s.manager.state.Variables {
		names = append(names, name)
	}
	return names
} 