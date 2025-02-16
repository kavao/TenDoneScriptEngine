package main

import (
	"flag"
	"gameengine/src/engine/ecs"
	"gameengine/src/engine/ecs/components"
	"gameengine/src/engine/ecs/core"
	"gameengine/src/engine/script"
	"gameengine/src/engine/systems"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	debugMode = flag.Bool("debug", false, "デバッグモードで実行")
)

type Game struct {
	world            *core.World
	scriptEngine     *script.ScriptEngine
	renderSystem     *systems.RenderSystem
	inputSystem      *systems.InputSystem
	textSystem       *systems.TextSystem
	physicsSystem    *systems.PhysicsSystem
	screenWidth      int
	screenHeight     int
	cleanupCounter   int
	scriptSelector   *systems.ScriptSelectorSystem
	scriptSelected   chan string
	isScriptSelected bool
}

func NewGame() *Game {
	g := &Game{
		screenWidth:  1280, // HDサイズ
		screenHeight: 720,
	}
	// ... 他の初期化コード ...
	return g
}

func main() {
	flag.Parse()
	world := ecs.NewWorld()
	scriptEngine := script.NewScriptEngine(world, "./scripts")

	// デフォルトの画面設定エンティティを作成
	configEntity := world.CreateEntity()
	configComponent := components.NewScreenConfigComponent()
	configEntity.AddComponent(configComponent)
	configEntity.AddTag("screen_config") // タグを追加

	// レンダリングシステムを作成して追加
	renderSystem := systems.NewRenderSystem()
	inputSystem := systems.NewInputSystem()
	textSystem := systems.NewTextSystem()
	physicsSystem := systems.NewPhysicsSystem()

	// ゲームの初期化
	game := &Game{
		world:            world,
		scriptEngine:     scriptEngine,
		renderSystem:     renderSystem,
		inputSystem:      inputSystem,
		textSystem:       textSystem,
		physicsSystem:    physicsSystem,
		screenWidth:      1280,
		screenHeight:     720,
		scriptSelected:   make(chan string, 1), // バッファ付きチャネル
		isScriptSelected: false,
	}

	if *debugMode {
		game.scriptSelector = systems.NewScriptSelectorSystem(world, "scripts_debug", func(script string) {
			game.scriptSelected <- script
			game.isScriptSelected = true
		})
		world.AddSystem(game.scriptSelector)
	}

	// スクリーン設定システムを追加
	screenConfigSystem := systems.NewScreenConfigSystem(game)
	world.AddSystem(screenConfigSystem)

	// 他のシステムを追加
	world.AddSystem(renderSystem)
	world.AddSystem(inputSystem)
	world.AddSystem(textSystem)
	world.AddSystem(physicsSystem)

	// ウィンドウ設定
	ebiten.SetWindowTitle("Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetScreenClearedEveryFrame(true)

	// 初期ウィンドウサイズを明示的に設定
	ebiten.SetWindowSize(1280, 720)

	// リサイズモードは最後に設定
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	// ゲームループの開始
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	if *debugMode && !g.isScriptSelected {
		// スクリプトが選択されるまで通常の更新をスキップ
		if err := g.world.Update(1.0 / 60.0); err != nil {
			return err
		}

		select {
		case script := <-g.scriptSelected:
			if err := g.scriptEngine.ExecuteFile(script); err != nil {
				return err
			}
			g.isScriptSelected = true
		default:
			// チャネルが空の場合は何もしない
		}
		return nil
	}

	// スクリプトエンジンの更新を最初に行う
	if err := g.scriptEngine.CallUpdate(); err != nil {
		return err
	}

	// ワールドの更新
	if err := g.world.Update(1.0 / 60.0); err != nil {
		return err
	}

	// 非アクティブなエンティティを定期的にクリーンアップ
	g.cleanupCounter++
	if g.cleanupCounter >= 60 { // 1秒ごとにクリーンアップ
		g.world.CleanupInactiveEntities()
		g.cleanupCounter = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255}) // 背景を黒に
	g.renderSystem.SetScreen(screen)
	g.renderSystem.Update(0)
	g.textSystem.SetScreen(screen) // テキストシステムの描画も追加
	g.textSystem.Update(0)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// デバッグ出力を追加
	// fmt.Printf("Layout called: outside=%dx%d, screen=%dx%d\n",
	// 	outsideWidth, outsideHeight, g.screenWidth, g.screenHeight)

	// 画面のアスペクト比を維持
	targetAspectRatio := float64(g.screenWidth) / float64(g.screenHeight)
	currentAspectRatio := float64(outsideWidth) / float64(outsideHeight)

	var width, height int
	if currentAspectRatio > targetAspectRatio {
		// 高さに合わせる
		height = g.screenHeight
		width = g.screenWidth
	} else {
		// 幅に合わせる
		width = g.screenWidth
		height = g.screenHeight
	}

	// fmt.Printf("Returning layout: %dx%d\n", width, height)
	return width, height
}

func (g *Game) SetScreenSize(width, height int) {
	g.screenWidth = width
	g.screenHeight = height
}
