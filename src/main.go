package main

import (
	"gameengine/src/engine/ecs"
	"gameengine/src/engine/script"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	world        *ecs.World
	scriptEngine *script.ScriptEngine
}

func main() {
	// ECSワールドの初期化
	world := ecs.NewWorld()

	// スクリプトエンジンの初期化
	scriptEngine := script.NewScriptEngine(world, "./scripts")

	// ゲームの初期化
	game := &Game{
		world:        world,
		scriptEngine: scriptEngine,
	}

	// メインスクリプトの実行
	if err := scriptEngine.ExecuteFile("main.star"); err != nil {
		log.Fatal(err)
	}

	// ウィンドウ設定
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Game")

	// ゲームループの開始
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	return g.world.Update(1.0 / 60.0)
}

func (g *Game) Draw(screen *ebiten.Image) {
	// 仮の描画処理
	screen.Fill(color.RGBA{0, 0, 0, 255})
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 800, 600
}
