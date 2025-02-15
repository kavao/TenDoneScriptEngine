package particle

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

// エミッター設定
type EmitterConfig struct {
	Position      Vector2D
	SpawnRate     float64  // 1秒あたりの生成数
	ParticleLife  Range    // パーティクルの寿命範囲
	InitialScale  Range    // 初期スケール範囲
	ScaleVelocity Range    // スケール変化速度範囲
	Speed         Range    // 速度範囲
	Angle         Range    // 角度範囲（ラジアン）
	Gravity       Vector2D // 重力加速度
	Color         color.RGBA
	AlphaVelocity Range // アルファ値変化速度範囲
	Image         *ebiten.Image
}

// 値の範囲
type Range struct {
	Min, Max float64
}

// パーティクルエミッター
type Emitter struct {
	config     EmitterConfig
	particles  []*Particle
	active     bool
	spawnTimer float64
}

func NewEmitter(config EmitterConfig, maxParticles int) *Emitter {
	particles := make([]*Particle, maxParticles)
	for i := range particles {
		particles[i] = NewParticle()
	}

	return &Emitter{
		config:    config,
		particles: particles,
		active:    true,
	}
}

func (e *Emitter) Update(dt float64) {
	if !e.active {
		return
	}

	// パーティクル生成
	e.spawnTimer += dt
	spawnInterval := 1.0 / e.config.SpawnRate
	for e.spawnTimer >= spawnInterval {
		e.spawnTimer -= spawnInterval
		e.spawnParticle()
	}

	// パーティクル更新
	for _, p := range e.particles {
		p.Update(dt)
	}
}

func (e *Emitter) Draw(screen *ebiten.Image) {
	for _, p := range e.particles {
		p.Draw(screen)
	}
}

func (e *Emitter) spawnParticle() {
	// 非アクティブなパーティクルを探す
	var particle *Particle
	for _, p := range e.particles {
		if !p.Active {
			particle = p
			break
		}
	}
	if particle == nil {
		return
	}

	// パーティクルの初期化
	angle := e.config.Angle.Min + rand.Float64()*(e.config.Angle.Max-e.config.Angle.Min)
	speed := e.config.Speed.Min + rand.Float64()*(e.config.Speed.Max-e.config.Speed.Min)

	particle.Position = e.config.Position
	particle.Velocity = Vector2D{
		X: speed * math.Cos(angle),
		Y: speed * math.Sin(angle),
	}
	particle.Acceleration = e.config.Gravity
	particle.Scale = e.config.InitialScale.Min + rand.Float64()*(e.config.InitialScale.Max-e.config.InitialScale.Min)
	particle.ScaleVel = e.config.ScaleVelocity.Min + rand.Float64()*(e.config.ScaleVelocity.Max-e.config.ScaleVelocity.Min)
	particle.Color = e.config.Color
	particle.AlphaVel = e.config.AlphaVelocity.Min + rand.Float64()*(e.config.AlphaVelocity.Max-e.config.AlphaVelocity.Min)
	particle.Life = e.config.ParticleLife.Min + rand.Float64()*(e.config.ParticleLife.Max-e.config.ParticleLife.Min)
	particle.MaxLife = particle.Life
	particle.Image = e.config.Image
	particle.Active = true
}

func (e *Emitter) SetPosition(x, y float64) {
	e.config.Position.X = x
	e.config.Position.Y = y
}

func (e *Emitter) SetActive(active bool) {
	e.active = active
} 