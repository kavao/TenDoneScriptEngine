package script

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gameengine/src/engine/ecs/components"
	"gameengine/src/engine/ecs/core"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"go.starlark.net/starlark"
)

// スクリプトエンジン
type ScriptEngine struct {
	mutex        sync.RWMutex
	world        *core.World
	thread       *starlark.Thread
	globals      starlark.StringDict
	scriptDir    string
	stateManager *StateManager
}

func NewScriptEngine(world *core.World, scriptDir string) *ScriptEngine {
	engine := &ScriptEngine{
		world:        world,
		thread:       &starlark.Thread{Name: "game"},
		globals:      make(starlark.StringDict),
		scriptDir:    scriptDir,
		stateManager: NewStateManager(world), // StateManagerを初期化
	}

	// デバッグ用
	fmt.Printf("ScriptEngine created with World: %p\n", world)

	// 基本的なグローバル関数の登録
	engine.registerBuiltins()

	return engine
}

// スクリプトの実行
func (e *ScriptEngine) ExecuteFile(filename string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// scripts.zipが存在する場合はそちらを優先
	if _, err := os.Stat("scripts.zip"); err == nil {
		return e.ExecuteZip("scripts.zip")
	}

	var path string
	if strings.HasPrefix(filename, "scripts_debug/") {
		path = filename
	} else {
		path = filepath.Join(e.scriptDir, filename)
	}

	fmt.Printf("Loading script from: %s\n", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read script file: %v", err)
	}

	fmt.Println("Globals before execution:", e.globals.Keys())

	// スクリプトを実行し、その結果をglobalsに保存
	globals, err := starlark.ExecFile(e.thread, path, data, e.globals)
	if err != nil {
		return err
	}

	// 組み込み関数を保持
	builtins := make(starlark.StringDict)
	for k, v := range e.globals {
		if _, ok := v.(*starlark.Builtin); ok {
			builtins[k] = v
		}
	}

	// グローバル変数を更新
	e.globals = globals

	// 組み込み関数を復元
	for k, v := range builtins {
		e.globals[k] = v
	}

	fmt.Println("Globals after execution:", e.globals.Keys())
	fmt.Printf("World pointer in ScriptEngine: %p\n", e.world) // デバッグ出力を追加
	fmt.Println("Available functions:", e.globals)

	if _, ok := e.globals["update"]; ok {
		fmt.Println("update function is registered")
	}
	if _, ok := e.globals["init"]; ok {
		fmt.Println("init function is registered")
	}

	return nil
}

// スクリプトの実行
func (e *ScriptEngine) ExecuteZip(zipPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	// アセットの読み込み
	for _, f := range r.File {
		// アセットディレクトリの処理
		if strings.HasPrefix(f.Name, "assets/") {
			if err := e.extractAsset(f); err != nil {
				return fmt.Errorf("failed to extract asset %s: %v", f.Name, err)
			}
			continue
		}

		// main.starの処理
		if f.Name == "main.star" {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return err
			}

			globals, err := starlark.ExecFile(e.thread, "main.star", data, e.globals)
			if err != nil {
				return err
			}
			e.globals = globals // グローバル変数を更新
			return nil
		}
	}

	return fmt.Errorf("main.star not found in zip")
}

func (e *ScriptEngine) extractAsset(f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	// 出力先のパスを作成
	outPath := filepath.Join(".", f.Name)
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}

	// ディレクトリの場合はスキップ
	if f.FileInfo().IsDir() {
		return nil
	}

	// ファイルを作成
	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// データをコピー
	_, err = io.Copy(outFile, rc)
	return err
}

// グローバル関数の登録
func (e *ScriptEngine) registerBuiltins() {
	e.globals["create_entity"] = starlark.NewBuiltin("create_entity", e.createEntity)
	e.globals["add_component"] = starlark.NewBuiltin("add_component", e.addComponent)
	e.globals["get_component"] = starlark.NewBuiltin("get_component", e.getComponent)
	e.globals["find_entities_by_tag"] = starlark.NewBuiltin("find_entities_by_tag", e.findEntitiesByTag)
	e.globals["add_tag"] = starlark.NewBuiltin("add_tag", e.addTag)
	e.globals["is_key_pressed"] = starlark.NewBuiltin("is_key_pressed", e.isKeyPressed)
	e.globals["print"] = starlark.NewBuiltin("print", e.print)
	e.globals["set_component"] = starlark.NewBuiltin("set_component", e.setComponent)
	e.globals["get_total_entities"] = starlark.NewBuiltin("get_total_entities", e.getTotalEntities)
	e.globals["set_state"] = starlark.NewBuiltin("set_state", e.setState)
	e.globals["get_state"] = starlark.NewBuiltin("get_state", e.getState)
	e.globals["set_states"] = starlark.NewBuiltin("set_states", e.setStates)
	e.globals["get_states"] = starlark.NewBuiltin("get_states", e.getStates)

	// loadコマンドを追加
	e.thread.Load = func(thread *starlark.Thread, module string) (starlark.StringDict, error) {
		data, err := ioutil.ReadFile(module)
		if err != nil {
			return nil, fmt.Errorf("failed to read module %s: %v", module, err)
		}

		// モジュールを実行し、グローバル変数を取得
		globals, err := starlark.ExecFile(thread, module, data, e.globals)
		if err != nil {
			return nil, err
		}

		return globals, nil
	}

	e.registerFunctions(e.thread, e.globals)
}

func (e *ScriptEngine) registerFunctions(thread *starlark.Thread, predeclared starlark.StringDict) {
	predeclared["set_screen_resolution"] = starlark.NewBuiltin("set_screen_resolution", func(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var preset string
		if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &preset); err != nil {
			return nil, err
		}

		fmt.Printf("Setting screen resolution to: %s\n", preset)
		// スクリーン設定エンティティを探す
		entities := e.world.FindEntitiesByTag("screen_config")
		fmt.Printf("Found %d screen config entities\n", len(entities))
		if len(entities) > 0 {
			if config := entities[0].GetComponent(4).(*components.ScreenConfigComponent); config != nil {
				config.SetResolution(preset)
				fmt.Printf("Updated resolution to: %dx%d\n", config.Width, config.Height)
			}
		}

		return starlark.None, nil
	})
}

// エンティティ作成
func (e *ScriptEngine) createEntity(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	// fmt.Printf("Creating entity in World %p\n", e.world) // デバッグ出力を追加
	entity := e.world.CreateEntity()
	// fmt.Printf("Created entity %d in World %p\n", entity.GetID(), e.world)
	return starlark.MakeInt64(int64(entity.GetID())), nil
}

// コンポーネント追加
func (e *ScriptEngine) addComponent(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		entityID      int64
		componentType string
		properties    *starlark.Dict
	)

	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 3, &entityID, &componentType, &properties); err != nil {
		if len(args) == 2 {
			properties = &starlark.Dict{}
		} else {
			return nil, err
		}
	}

	entity := e.world.GetEntity(core.EntityID(entityID))
	if entity == nil {
		fmt.Printf("Failed to get entity ID: %d (created from createEntity)\n", entityID)
		return nil, fmt.Errorf("entity not found: %d", entityID)
	}

	// コンポーネントの作成と追加
	switch componentType {
	case "transform":
		xVal, _, err := properties.Get(starlark.String("x"))
		if err != nil {
			return nil, err
		}
		yVal, _, err := properties.Get(starlark.String("y"))
		if err != nil {
			return nil, err
		}
		x, _ := starlark.AsFloat(xVal)
		y, _ := starlark.AsFloat(yVal)
		component := &components.TransformComponent{
			X: x,
			Y: y,
		}
		entity.AddComponent(component)

	case "sprite":
		component := components.NewSpriteComponent()
		if widthVal, _, err := properties.Get(starlark.String("width")); err == nil {
			if width, _ := starlark.AsInt32(widthVal); width > 0 {
				component.Width = int(width)
				component.Sprite = ebiten.NewImage(int(width), component.Height)
			}
		}
		if heightVal, _, err := properties.Get(starlark.String("height")); err == nil {
			if height, _ := starlark.AsInt32(heightVal); height > 0 {
				component.Height = int(height)
				component.Sprite = ebiten.NewImage(component.Width, int(height))
			}
		}
		if colorVal, _, err := properties.Get(starlark.String("color")); err == nil {
			if colorStr, _ := starlark.AsString(colorVal); colorStr != "" {
				switch colorStr {
				case "cyan":
					component.SetColor(color.RGBA{0, 255, 255, 255})
				case "white":
					component.SetColor(color.White)
					// 他の色も必要に応じて追加
				}
			}
		}
		entity.AddComponent(component)

	case "text":
		component := components.NewTextComponent()
		if textVal, _, err := properties.Get(starlark.String("text")); err == nil {
			if text, _ := starlark.AsString(textVal); text != "" {
				component.Text = text
				fmt.Printf("Setting text to: %s\n", text) // デバッグ出力を追加
			}
		}
		if xVal, _, err := properties.Get(starlark.String("x")); err == nil {
			if x, _ := starlark.AsFloat(xVal); x != 0 {
				component.X = x
			}
		}
		if yVal, _, err := properties.Get(starlark.String("y")); err == nil {
			if y, _ := starlark.AsFloat(yVal); y != 0 {
				component.Y = y
			}
		}
		entity.AddComponent(component)
		fmt.Printf("Added text component to entity %d\n", entityID) // デバッグ出力を追加

	case "physics":
		component := components.NewPhysicsComponent()
		if vxVal, _, err := properties.Get(starlark.String("velocity_x")); err == nil {
			if vx, _ := starlark.AsFloat(vxVal); vx != 0 {
				component.VelocityX = vx
			}
		}
		if vyVal, _, err := properties.Get(starlark.String("velocity_y")); err == nil {
			if vy, _ := starlark.AsFloat(vyVal); vy != 0 {
				component.VelocityY = vy
			}
		}
		if gravityVal, _, err := properties.Get(starlark.String("gravity")); err == nil {
			if gravity, _ := starlark.AsFloat(gravityVal); gravity != 0 {
				component.Gravity = gravity
			}
		}
		entity.AddComponent(component)
	}

	return starlark.None, nil
}

// コンポーネント取得
func (e *ScriptEngine) getComponent(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var entityID int64
	var componentType string
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &entityID, &componentType); err != nil {
		return nil, err
	}

	entity := e.world.GetEntity(core.EntityID(entityID))
	if entity == nil {
		fmt.Printf("Entity %d not found in World %p\n", entityID, e.world)
		return starlark.None, nil
	}

	var componentID core.ComponentID
	switch componentType {
	case "transform":
		componentID = 1
	case "sprite":
		componentID = 2
	case "text":
		componentID = 3
	case "physics":
		componentID = 5
	default:
		return nil, fmt.Errorf("unknown component type: %s", componentType)
	}

	component := entity.GetComponent(componentID)
	if component == nil {
		// コンポーネントが見つからない場合もNoneを返す
		return starlark.None, nil
	}

	// コンポーネントの型に応じてStarlark辞書を作成
	dict := starlark.NewDict(0)
	switch c := component.(type) {
	case *components.TransformComponent:
		dict.SetKey(starlark.String("x"), starlark.Float(c.X))
		dict.SetKey(starlark.String("y"), starlark.Float(c.Y))
	case *components.PhysicsComponent:
		dict.SetKey(starlark.String("velocity_x"), starlark.Float(c.VelocityX))
		dict.SetKey(starlark.String("velocity_y"), starlark.Float(c.VelocityY))
	}

	return dict, nil
}

// デバッグ出力
func (e *ScriptEngine) print(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	for _, arg := range args {
		fmt.Print(arg.String() + " ")
	}
	fmt.Println()
	return starlark.None, nil
}

// タグによるエンティティの検索
func (e *ScriptEngine) findEntitiesByTag(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var tag string
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &tag); err != nil {
		return nil, err
	}

	entities := e.world.FindEntitiesByTag(tag)
	result := make([]starlark.Value, len(entities))
	for i, entity := range entities {
		result[i] = starlark.MakeInt64(int64(entity.GetID()))
	}
	return starlark.NewList(result), nil
}

// タグの追加
func (e *ScriptEngine) addTag(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var entityID int64
	var tag string
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &entityID, &tag); err != nil {
		return nil, err
	}

	entity := e.world.GetEntity(core.EntityID(entityID))
	if entity == nil {
		return nil, fmt.Errorf("entity not found: %d", entityID)
	}

	entity.AddTag(tag)
	return starlark.None, nil
}

// キー入力の検知
func (e *ScriptEngine) isKeyPressed(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var keyName string
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 1, &keyName); err != nil {
		return nil, err
	}

	var key ebiten.Key
	switch keyName {
	case "Space":
		key = ebiten.KeySpace
	case "ArrowLeft":
		key = ebiten.KeyArrowLeft
	case "ArrowRight":
		key = ebiten.KeyArrowRight
	case "ArrowUp":
		key = ebiten.KeyArrowUp
	case "ArrowDown":
		key = ebiten.KeyArrowDown
	default:
		return nil, fmt.Errorf("unknown key: %s", keyName)
	}

	result := ebiten.IsKeyPressed(key)
	// if result {
	// 	fmt.Printf("Key %s is pressed\n", keyName)
	// }
	return starlark.Bool(result), nil
}

// update関数の呼び出し
func (e *ScriptEngine) CallUpdate() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// グローバルから"update"関数を取得
	updateFn, ok := e.globals["update"]
	if !ok {
		// fmt.Println("No update function found") // コメントアウト
		return nil
	}

	// fmt.Println("Calling update function") // コメントアウト
	// Starlarkの関数として呼び出し
	if fn, ok := updateFn.(starlark.Callable); ok {
		_, err := starlark.Call(e.thread, fn, nil, nil)
		if err != nil {
			return fmt.Errorf("error calling update: %v", err)
		}
		// fmt.Println("Update function called successfully") // コメントアウト
	}

	return nil
}

func (e *ScriptEngine) setComponent(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var entityID int64
	var componentType string
	var properties *starlark.Dict

	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 3, &entityID, &componentType, &properties); err != nil {
		return nil, err
	}

	entity := e.world.GetEntity(core.EntityID(entityID))
	if entity == nil {
		// エンティティが見つからない場合はNoneを返す
		return starlark.None, nil
	}

	switch componentType {
	case "transform":
		component := entity.GetComponent(1)
		if component == nil {
			return starlark.None, nil
		}
		transform := component.(*components.TransformComponent)
		if xVal, _, err := properties.Get(starlark.String("x")); err == nil {
			if x, _ := starlark.AsFloat(xVal); true {
				transform.X = x
			}
		}
		if yVal, _, err := properties.Get(starlark.String("y")); err == nil {
			if y, _ := starlark.AsFloat(yVal); true {
				transform.Y = y
			}
		}
	case "text":
		component := entity.GetComponent(3)
		if component == nil {
			return starlark.None, nil
		}
		textComponent := component.(*components.TextComponent)
		if textVal, _, err := properties.Get(starlark.String("text")); err == nil {
			if textStr, _ := starlark.AsString(textVal); true {
				textComponent.Text = textStr
			}
		}
		if xVal, _, err := properties.Get(starlark.String("x")); err == nil {
			if x, _ := starlark.AsFloat(xVal); true {
				textComponent.X = x
			}
		}
		if yVal, _, err := properties.Get(starlark.String("y")); err == nil {
			if y, _ := starlark.AsFloat(yVal); true {
				textComponent.Y = y
			}
		}
	}

	return starlark.None, nil
}

// 全エンティティ数を取得
func (e *ScriptEngine) getTotalEntities(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	count := e.world.GetTotalEntities()
	return starlark.MakeInt64(int64(count)), nil
}

func (e *ScriptEngine) setState(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var entityID int64
	var key string
	var value starlark.Value
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 3, &entityID, &key, &value); err != nil {
		return nil, err
	}

	// Starlark値をGo値に変換
	var goValue interface{}
	switch v := value.(type) {
	case starlark.Int:
		goValue, _ = v.Int64()
	case starlark.Float:
		goValue = float64(v)
	case starlark.String:
		goValue = string(v)
	case starlark.Bool:
		goValue = bool(v)
		// 必要に応じて他の型も追加
	}

	e.stateManager.SetState(core.EntityID(entityID), key, goValue)
	return starlark.None, nil
}

func (e *ScriptEngine) getState(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var entityID int64
	var key string
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &entityID, &key); err != nil {
		return nil, err
	}

	state := e.stateManager.GetState(core.EntityID(entityID), key)
	if state == nil {
		return starlark.None, nil
	}

	// Starlark値をGo値に変換
	var goValue starlark.Value
	switch v := state.(type) {
	case int64:
		goValue = starlark.MakeInt64(v)
	case float64:
		goValue = starlark.Float(v)
	case string:
		goValue = starlark.String(v)
	case bool:
		goValue = starlark.Bool(v)
	}

	return goValue, nil
}

func (e *ScriptEngine) setStates(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var entityID int64
	var statesDict *starlark.Dict
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &entityID, &statesDict); err != nil {
		return nil, err
	}

	// Starlark辞書をGoのmapに変換
	states := make(map[string]interface{})
	for _, item := range statesDict.Items() {
		key, ok := item[0].(starlark.String)
		if !ok {
			return nil, fmt.Errorf("key must be string, got %s", item[0].Type())
		}

		// 値の型変換
		var value interface{}
		switch v := item[1].(type) {
		case starlark.Int:
			value, _ = v.Int64()
		case starlark.Float:
			value = float64(v)
		case starlark.String:
			value = string(v)
		case starlark.Bool:
			value = bool(v)
		default:
			return nil, fmt.Errorf("unsupported value type: %s", item[1].Type())
		}
		states[string(key)] = value
	}

	e.stateManager.SetStates(core.EntityID(entityID), states)
	return starlark.None, nil
}

func (e *ScriptEngine) getStates(thread *starlark.Thread, b *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var entityID int64
	var keysList *starlark.List
	if err := starlark.UnpackPositionalArgs(b.Name(), args, kwargs, 2, &entityID, &keysList); err != nil {
		return nil, err
	}

	// キーのリストをGo形式に変換
	keys := make([]string, 0, keysList.Len())
	for i := 0; i < keysList.Len(); i++ {
		key, ok := keysList.Index(i).(starlark.String)
		if !ok {
			return nil, fmt.Errorf("key must be string, got %s", keysList.Index(i).Type())
		}
		keys = append(keys, string(key))
	}

	// 状態を取得
	states := e.stateManager.GetStates(core.EntityID(entityID), keys)

	// 結果をStarlark辞書に変換
	result := starlark.NewDict(len(states))
	for key, value := range states {
		var starlarkValue starlark.Value
		switch v := value.(type) {
		case int64:
			starlarkValue = starlark.MakeInt64(v)
		case float64:
			starlarkValue = starlark.Float(v)
		case string:
			starlarkValue = starlark.String(v)
		case bool:
			starlarkValue = starlark.Bool(v)
		default:
			continue
		}
		result.SetKey(starlark.String(key), starlarkValue)
	}

	return result, nil
}
