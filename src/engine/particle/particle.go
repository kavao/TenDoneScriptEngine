package particle

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// パーティクル1つの情報
type Particle struct {
	Position    Vector2D
	Velocity    Vector2D
	Acceleration Vector2D
	Scale       float64
	ScaleVel    float64
	Rotation    float64
	RotationVel float64
	Color       color.RGBA
	Alpha       float64
	AlphaVel    float64
	Life        float64
	MaxLife     float64
	Active      bool
	Image       *ebiten.Image
}

// 2D座標/ベクトル
type Vector2D struct {
	X, Y float64
}

func NewParticle() *Particle {
	return &Particle{
		Scale:  1.0,
		Alpha:  1.0,
		Active: false,
	}
}

func (p *Particle) Update(dt float64) {
	if !p.Active {
		return
	}

	// 物理演算の更新
	p.Velocity.X += p.Acceleration.X * dt
	p.Velocity.Y += p.Acceleration.Y * dt
	p.Position.X += p.Velocity.X * dt
	p.Position.Y += p.Velocity.Y * dt

	// スケール更新
	p.Scale += p.ScaleVel * dt

	// 回転更新
	p.Rotation += p.RotationVel * dt

	// アルファ値更新
	p.Alpha += p.AlphaVel * dt
	p.Alpha = math.Max(0, math.Min(1, p.Alpha))

	// ライフタイム更新
	p.Life -= dt
	if p.Life <= 0 {
		p.Active = false
	}
}

func (p *Particle) Draw(screen *ebiten.Image) {
	if !p.Active || p.Image == nil {
		return
	}

	op := &ebiten.DrawImageOptions{}

	// スケーリング
	op.GeoM.Scale(p.Scale, p.Scale)

	// 回転
	op.GeoM.Rotate(p.Rotation)

	// 移動
	w, h := p.Image.Bounds().Dx(), p.Image.Bounds().Dy()
	op.GeoM.Translate(-float64(w)/2, -float64(h)/2)
	op.GeoM.Translate(p.Position.X, p.Position.Y)

	// カラー/アルファ値
	op.ColorM.Scale(
		float64(p.Color.R)/255,
		float64(p.Color.G)/255,
		float64(p.Color.B)/255,
		float64(p.Color.A)/255*p.Alpha,
	)

	screen.DrawImage(p.Image, op)
} 