package audio

import (
	"github.com/hajimehoshi/ebiten/v2/audio"
)

// SEプレイヤー
type SEPlayer struct {
	player *audio.Player
	volume float64
}

func NewSEPlayer(player *audio.Player) *SEPlayer {
	return &SEPlayer{
		player: player,
		volume: 1.0,
	}
}

func (p *SEPlayer) Play(volume float64) error {
	player := p.player
	player.SetVolume(volume)
	player.Play() // 戻り値を無視
	return nil
}

func (p *SEPlayer) SetVolume(volume float64) {
	p.volume = volume
}

func (p *SEPlayer) Stop() error {
	p.player.Pause()
	return nil
}
