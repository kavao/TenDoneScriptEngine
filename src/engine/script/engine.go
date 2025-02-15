package script

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sync"

	"go.starlark.net/starlark"
	"gameengine/src/engine/ecs"
)

// スクリプトエンジン
type ScriptEngine struct {
	mutex       sync.RWMutex
	world       *ecs.World
	thread      *starlark.Thread
	globals     starlark.StringDict
	scriptDir   string
}

func NewScriptEngine(world *ecs.World, scriptDir string) *ScriptEngine {
	engine := &ScriptEngine{
		world:     world,
		thread:    &starlark.Thread{Name: "game"},
		globals:   make(starlark.StringDict),
		scriptDir: scriptDir,
	}

	// 基本的なグローバル関数の登録
	engine.registerBuiltins()

	return engine
}

// スクリプトの実行
func (e *ScriptEngine) ExecuteFile(filename string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	path := filepath.Join(e.scriptDir, filename)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read script file: %v", err)
	}

	_, err = starlark.ExecFile(e.thread, path, data, e.globals)
	return err
}

// グローバル関数の登録
func (e *ScriptEngine) registerBuiltins() {
	e.globals["create_entity"] = starlark.NewBuiltin("create_entity", e.createEntity)
	e.globals["add_component"] = starlark.NewBuiltin("add_component", e.addComponent)
	e.globals["get_component"] = starlark.NewBuiltin("get_component", e.getComponent)
	e.globals["print"] = starlark.NewBuiltin("print", e.print)
}

// エンティティ作成
func (e *ScriptEngine) createEntity(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	entity := e.world.CreateEntity()
	return starlark.MakeInt64(int64(entity.GetID())), nil
}

// コンポーネント追加
func (e *ScriptEngine) addComponent(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var entityID int64
	var componentType string
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &entityID, &componentType); err != nil {
		return nil, err
	}

	entity := e.world.GetEntity(ecs.EntityID(entityID))
	if entity == nil {
		return nil, fmt.Errorf("entity not found: %d", entityID)
	}

	// TODO: コンポーネントの作成と追加
	return starlark.None, nil
}

// コンポーネント取得
func (e *ScriptEngine) getComponent(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var entityID int64
	var componentType string
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &entityID, &componentType); err != nil {
		return nil, err
	}

	entity := e.world.GetEntity(ecs.EntityID(entityID))
	if entity == nil {
		return nil, fmt.Errorf("entity not found: %d", entityID)
	}

	// TODO: コンポーネントの取得と変換
	return starlark.None, nil
}

// デバッグ出力
func (e *ScriptEngine) print(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	for _, arg := range args {
		fmt.Print(arg.String() + " ")
	}
	fmt.Println()
	return starlark.None, nil
} 