package input

import (
	"encoding/json"
	"os"
)

type InputConfig struct {
	KeyBindings map[Action][]InputBinding `json:"key_bindings"`
}

func NewInputConfig() *InputConfig {
	return &InputConfig{
		KeyBindings: make(map[Action][]InputBinding),
	}
}

func (c *InputConfig) Save(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (c *InputConfig) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, c)
}

func (c *InputConfig) ApplyToManager(im *InputManager) {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	// 既存のバインディングをクリア
	im.bindings = make(map[Action][]InputBinding)

	// 設定を適用
	for action, bindings := range c.KeyBindings {
		im.bindings[action] = bindings
	}
} 