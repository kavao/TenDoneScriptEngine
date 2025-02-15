package audio

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2/audio"
)

// BGMプレイヤー
type BGMPlayer struct {
	audioPlayer   *audio.Player
	volume        float64
	targetVolume  float64
	fadeStartTime time.Time
	fadeDuration  time.Duration
	isFading      bool
	isLooping     bool
}

func NewBGMPlayer(stream *audio.Player) *BGMPlayer {
	return &BGMPlayer{
		audioPlayer:  stream,
		volume:       1.0,
		targetVolume: 1.0,
	}
}

func (p *BGMPlayer) Play(volume float64, loop bool) error {
	p.volume = volume
	p.targetVolume = volume
	p.isLooping = loop

	if err := p.audioPlayer.Rewind(); err != nil {
		return err
	}

	p.audioPlayer.SetVolume(volume)
	p.audioPlayer.Play() // 戻り値を無視
	return nil
}

func (p *BGMPlayer) Stop() error {
	p.audioPlayer.Pause()
	return nil
}

func (p *BGMPlayer) Pause() error {
	p.audioPlayer.Pause()
	return nil
}

func (p *BGMPlayer) Resume() error {
	p.audioPlayer.Play()
	return nil
}

func (p *BGMPlayer) SetVolume(volume float64) {
	p.volume = volume
	p.targetVolume = volume
	p.audioPlayer.SetVolume(volume)
}

func (p *BGMPlayer) FadeTo(targetVolume float64, duration time.Duration) {
	p.targetVolume = targetVolume
	p.fadeStartTime = time.Now()
	p.fadeDuration = duration
	p.isFading = true
}

func (p *BGMPlayer) Update() error {
	if p.isFading {
		elapsed := time.Since(p.fadeStartTime)
		if elapsed >= p.fadeDuration {
			p.volume = p.targetVolume
			p.isFading = false
		} else {
			progress := float64(elapsed) / float64(p.fadeDuration)
			p.volume = p.volume + (p.targetVolume-p.volume)*progress
		}
		p.audioPlayer.SetVolume(p.volume)
	}

	if p.isLooping && !p.audioPlayer.IsPlaying() {
		if err := p.audioPlayer.Rewind(); err != nil {
			return err
		}
		p.audioPlayer.Play() // 戻り値を無視
		return nil
	}

	return nil
}
