package systems

import (
	"gameengine/src/engine/ecs"
	"gameengine/src/engine/ecs/components"
	"gameengine/src/engine/ecs/core"
	"os"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// キー状態を追跡する構造体を追加
type keyStates struct {
	upPressed    bool
	downPressed  bool
	enterPressed bool
}

type ScriptSelectorSystem struct {
	*ecs.BaseSystem
	world          *core.World
	scripts        []string
	selectedIndex  int
	isActive       bool
	onScriptSelect func(string)
	textEntities   []core.EntityID // 現在表示中のテキストエンティティを追跡
	keyStates      keyStates       // キー状態を追加
}

func NewScriptSelectorSystem(world *core.World, scriptDir string, onSelect func(string)) *ScriptSelectorSystem {
	// スクリプトファイルの一覧を取得
	files, _ := os.ReadDir(scriptDir)
	var scripts []string
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".star") && !strings.HasPrefix(f.Name(), "_") {
			scripts = append(scripts, filepath.ToSlash(filepath.Join(scriptDir, f.Name())))
		}
	}

	return &ScriptSelectorSystem{
		BaseSystem:     ecs.NewBaseSystem(ecs.PriorityUpdate, []core.ComponentID{3}), // Text component
		world:          world,
		scripts:        scripts,
		selectedIndex:  0,
		isActive:       true,
		onScriptSelect: onSelect,
	}
}

func (s *ScriptSelectorSystem) Update(dt float64) error {
	// システムが非アクティブな場合は即座にリターン
	if !s.isActive {
		return nil
	}

	// 前回のテキストエンティティを削除
	for _, entityID := range s.textEntities {
		s.world.Mutex.Lock()
		s.world.ToRemove = append(s.world.ToRemove, entityID)
		s.world.Mutex.Unlock()
	}
	s.textEntities = make([]core.EntityID, 0)

	// 選択肢の表示
	for i, script := range s.scripts {
		entity := s.world.CreateEntity()
		textComp := components.NewTextComponent()

		// プレフィックスの幅を統一
		prefix := "[_] "
		if i == s.selectedIndex {
			prefix = "[o] "
		}

		// パスを短く表示するように変更
		displayPath := filepath.Base(script)
		textComp.Text = prefix + displayPath

		// 表示位置を調整
		textComp.X = 50
		textComp.Y = 100 + float64(i*30)
		entity.AddComponent(textComp)
		s.textEntities = append(s.textEntities, entity.GetID())
	}

	// キー入力での選択
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		if !s.keyStates.upPressed {
			s.selectedIndex--
			if s.selectedIndex < 0 {
				s.selectedIndex = len(s.scripts) - 1
			}
			s.keyStates.upPressed = true
		}
	} else {
		s.keyStates.upPressed = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		if !s.keyStates.downPressed {
			s.selectedIndex++
			if s.selectedIndex >= len(s.scripts) {
				s.selectedIndex = 0
			}
			s.keyStates.downPressed = true
		}
	} else {
		s.keyStates.downPressed = false
	}

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if !s.keyStates.enterPressed {
			if s.selectedIndex >= 0 && s.selectedIndex < len(s.scripts) {
				s.onScriptSelect(s.scripts[s.selectedIndex])
				s.isActive = false
				// 全てのテキストエンティティを削除
				for _, entityID := range s.textEntities {
					s.world.Mutex.Lock()
					s.world.ToRemove = append(s.world.ToRemove, entityID)
					s.world.Mutex.Unlock()
				}
				s.textEntities = make([]core.EntityID, 0)
			}
			s.keyStates.enterPressed = true
		}
	} else {
		s.keyStates.enterPressed = false
	}

	return nil
}
