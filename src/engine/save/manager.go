package save

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.starlark.net/starlark"
)

type SaveManager struct {
	state       *GameState
	savePath    string
	maxSlots    int
	currentSlot int
}

type SaveData struct {
	GameState  *GameState `json:"game_state"`
	SaveTime   time.Time  `json:"save_time"`
	SlotNumber int        `json:"slot_number"`
}

func NewSaveManager(savePath string, maxSlots int) *SaveManager {
	return &SaveManager{
		state:    NewGameState(),
		savePath: savePath,
		maxSlots: maxSlots,
	}
}

func (m *SaveManager) Save(slot int) error {
	if slot < 0 || slot >= m.maxSlots {
		return fmt.Errorf("invalid save slot: %d", slot)
	}

	saveData := &SaveData{
		GameState:  m.state,
		SaveTime:   time.Now(),
		SlotNumber: slot,
	}

	// セーブデータをJSON形式にエンコード
	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %v", err)
	}

	// セーブファイルのパスを生成
	filename := fmt.Sprintf("save_%d.json", slot)
	filepath := filepath.Join(m.savePath, filename)

	// ディレクトリが存在することを確認
	if err := os.MkdirAll(m.savePath, 0755); err != nil {
		return fmt.Errorf("failed to create save directory: %v", err)
	}

	// ファイルに書き込み
	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return fmt.Errorf("failed to write save file: %v", err)
	}

	m.currentSlot = slot
	return nil
}

func (m *SaveManager) Load(slot int) error {
	if slot < 0 || slot >= m.maxSlots {
		return fmt.Errorf("invalid save slot: %d", slot)
	}

	// セーブファイルのパスを生成
	filename := fmt.Sprintf("save_%d.json", slot)
	filepath := filepath.Join(m.savePath, filename)

	// ファイルを読み込み
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read save file: %v", err)
	}

	var saveData SaveData
	if err := json.Unmarshal(data, &saveData); err != nil {
		return fmt.Errorf("failed to unmarshal save data: %v", err)
	}

	// バージョンチェック
	if saveData.GameState.SaveVersion > m.state.SaveVersion {
		return fmt.Errorf("save data version (%d) is newer than current version (%d)",
			saveData.GameState.SaveVersion, m.state.SaveVersion)
	}

	m.state = saveData.GameState
	m.currentSlot = slot
	return nil
}

// Starlark APIのためのメソッド
func (m *SaveManager) GetStateForStarlark() map[string]starlark.Value {
	result := make(map[string]starlark.Value)
	
	m.state.mutex.RLock()
	defer m.state.mutex.RUnlock()

	for k, v := range m.state.Variables {
		result[k] = convertToStarlarkValue(v)
	}

	return result
}

func (m *SaveManager) SetStateFromStarlark(name string, value starlark.Value) {
	m.state.SetVariable(name, convertStarlarkValue(value))
} 