package save

import (
	"sync"

	"go.starlark.net/starlark"
)

// グローバルな状態管理
type GameState struct {
	mutex       sync.RWMutex
	Variables   map[string]interface{} `json:"variables"`    // Starlarkと共有する変数
	Flags       map[string]bool        `json:"flags"`        // フラグ管理
	Progress    map[string]int         `json:"progress"`     // 進行度管理
	SaveVersion int                    `json:"save_version"` // セーブデータのバージョン
}

func NewGameState() *GameState {
	return &GameState{
		Variables:   make(map[string]interface{}),
		Flags:       make(map[string]bool),
		Progress:    make(map[string]int),
		SaveVersion: 1,
	}
}

// Starlark用の変数設定
func (g *GameState) SetVariable(name string, value interface{}) {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.Variables[name] = value
}

// Starlark用の変数取得
func (g *GameState) GetVariable(name string) (interface{}, bool) {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	val, exists := g.Variables[name]
	return val, exists
}

// Starlarkの値をGo側の値に変換
func convertStarlarkValue(v starlark.Value) interface{} {
	switch v := v.(type) {
	case starlark.String:
		return v.GoString()
	case starlark.Int:
		i, _ := v.Int64()
		return i
	case starlark.Float:
		return float64(v)
	case starlark.Bool:
		return bool(v)
	case *starlark.List:
		result := make([]interface{}, 0, v.Len())
		for i := 0; i < v.Len(); i++ {
			result = append(result, convertStarlarkValue(v.Index(i)))
		}
		return result
	case *starlark.Dict:
		result := make(map[string]interface{})
		for _, item := range v.Items() {
			key, _ := starlark.AsString(item[0])
			result[key] = convertStarlarkValue(item[1])
		}
		return result
	default:
		return nil
	}
}

// Go側の値をStarlarkの値に変換
func convertToStarlarkValue(v interface{}) starlark.Value {
	switch v := v.(type) {
	case string:
		return starlark.String(v)
	case int:
		return starlark.MakeInt(v)
	case int64:
		return starlark.MakeInt64(v)
	case float64:
		return starlark.Float(v)
	case bool:
		return starlark.Bool(v)
	case []interface{}:
		list := starlark.NewList(nil)
		for _, item := range v {
			list.Append(convertToStarlarkValue(item))
		}
		return list
	case map[string]interface{}:
		dict := starlark.NewDict(len(v))
		for k, val := range v {
			dict.SetKey(starlark.String(k), convertToStarlarkValue(val))
		}
		return dict
	default:
		return starlark.None
	}
}
