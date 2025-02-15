package main

import (
	"flag"
	"gameengine/src/engine/ecs"
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
	world          *ecs.World
	scriptEngine   *script.ScriptEngine
	renderSystem   *systems.RenderSystem
	textSystem     *systems.TextSystem
	cleanupCounter int
}

func main() {
	flag.Parse()

	// ECSワールドの初期化
	world := ecs.NewWorld()

	// レンダリングシステムを作成して追加
	renderSystem := systems.NewRenderSystem()
	inputSystem := systems.NewInputSystem()
	textSystem := systems.NewTextSystem()
	physicsSystem := systems.NewPhysicsSystem()
	world.AddSystem(renderSystem)
	world.AddSystem(inputSystem)
	world.AddSystem(textSystem)
	world.AddSystem(physicsSystem)

	// スクリプトの選択
	scriptPath := "main.star" // デフォルト
	if *debugMode {
		// デバッグモードの場合、スクリプトを選択
		var err error
		scriptPath, err = script.ShowScriptSelector("scripts_debug")
		if err != nil {
			log.Fatal(err)
		}
	}

	// スクリプトエンジンの初期化
	scriptEngine := script.NewScriptEngine(world, "./scripts")

	// ゲームの初期化
	game := &Game{
		world:        world,
		scriptEngine: scriptEngine,
		renderSystem: renderSystem,
		textSystem:   textSystem,
	}

	// スクリプトの実行（scripts.zipがある場合はそちらを優先）
	if err := scriptEngine.ExecuteFile(scriptPath); err != nil {
		log.Fatal(err)
	}

	// ウィンドウ設定
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Game")
	ebiten.SetWindowResizable(true)
	ebiten.SetScreenClearedEveryFrame(true)
	ebiten.SetWindowFloating(true) // ウィンドウを最前面に表示

	// ゲームループの開始
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	// ワールドの更新（これにより全てのシステムが更新される）
	if err := g.world.Update(1.0 / 60.0); err != nil {
		return err
	}

	// スクリプトの更新
	if err := g.scriptEngine.CallUpdate(); err != nil {
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
	return 800, 600
}
