package audio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/jfreymuth/oggvorbis"
)

const (
	sampleRate     = 44100
	bytesPerSample = 4 // float32 stereo
)

// サウンドタイプ
type SoundType int

const (
	BGM SoundType = iota
	SE
)

// サウンドマネージャー
type AudioManager struct {
	mutex        sync.RWMutex
	audioContext *audio.Context
	bgmPlayers   map[string]*BGMPlayer
	sePlayers    map[string]*SEPlayer
	bgmVolume    float64
	seVolume     float64
}

func NewAudioManager() (*AudioManager, error) {
	audioContext := audio.NewContext(sampleRate)

	return &AudioManager{
		audioContext: audioContext,
		bgmPlayers:   make(map[string]*BGMPlayer),
		sePlayers:    make(map[string]*SEPlayer),
		bgmVolume:    1.0,
		seVolume:     1.0,
	}, nil
}

// サウンドのロード
func (m *AudioManager) LoadSound(id string, data []byte, soundType SoundType) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var err error
	switch soundType {
	case BGM:
		err = m.LoadBGM(id, data)
	case SE:
		err = m.LoadSE(id, data)
	default:
		return fmt.Errorf("unknown sound type: %v", soundType)
	}

	return err
}

// BGMの再生
func (m *AudioManager) PlayBGM(id string, volume float64, loop bool) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	player, exists := m.bgmPlayers[id]
	if !exists {
		return fmt.Errorf("BGM not found: %s", id)
	}

	return player.Play(volume*m.bgmVolume, loop)
}

// BGMの停止
func (m *AudioManager) StopBGM(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	player, exists := m.bgmPlayers[id]
	if !exists {
		return fmt.Errorf("BGM not found: %s", id)
	}

	return player.Stop()
}

// SEの再生
func (m *AudioManager) PlaySE(id string, volume float64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	player, exists := m.sePlayers[id]
	if !exists {
		return fmt.Errorf("SE not found: %s", id)
	}

	return player.Play(volume * m.seVolume)
}

// 音量の設定
func (m *AudioManager) SetVolume(soundType SoundType, volume float64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	switch soundType {
	case BGM:
		m.bgmVolume = volume
		for _, player := range m.bgmPlayers {
			player.SetVolume(volume)
		}
	case SE:
		m.seVolume = volume
	}
}

// 更新処理
func (m *AudioManager) Update() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	for _, player := range m.bgmPlayers {
		if err := player.Update(); err != nil {
			return err
		}
	}

	return nil
}

func (a *AudioManager) LoadBGM(name string, data []byte) error {
	var stream io.Reader
	var err error

	if isOggFormat(data) {
		stream = bytes.NewReader(data)
		oggReader, err := oggvorbis.NewReader(stream)
		if err != nil {
			return fmt.Errorf("failed to decode audio: %v", err)
		}
		stream = &oggAdapter{Reader: oggReader}
	} else {
		stream = bytes.NewReader(data)
		stream, err = mp3.DecodeWithSampleRate(sampleRate, stream)
		if err != nil {
			return fmt.Errorf("failed to decode audio: %v", err)
		}
	}

	// 既存のBGMを停止
	if a.bgmPlayers[name] != nil {
		a.bgmPlayers[name].Stop()
	}

	// 新しいBGMプレイヤーを作成
	player, err := a.audioContext.NewPlayer(stream)
	if err != nil {
		return fmt.Errorf("failed to create player: %v", err)
	}

	bgmPlayer := NewBGMPlayer(player)
	a.bgmPlayers[name] = bgmPlayer

	return nil
}

func (a *AudioManager) LoadSE(name string, data []byte) error {
	var stream io.Reader
	var err error

	if isOggFormat(data) {
		stream = bytes.NewReader(data)
		oggReader, err := oggvorbis.NewReader(stream)
		if err != nil {
			return fmt.Errorf("failed to decode audio: %v", err)
		}
		stream = &oggAdapter{Reader: oggReader}
	} else {
		stream = bytes.NewReader(data)
		stream, err = mp3.DecodeWithSampleRate(sampleRate, stream)
		if err != nil {
			return fmt.Errorf("failed to decode audio: %v", err)
		}
	}

	player, err := a.audioContext.NewPlayer(stream)
	if err != nil {
		return fmt.Errorf("failed to create player: %v", err)
	}

	sePlayer := NewSEPlayer(player)
	a.sePlayers[name] = sePlayer

	return nil
}

func (a *AudioManager) Close() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for _, player := range a.bgmPlayers {
		player.Stop()
	}
	for _, player := range a.sePlayers {
		player.Stop()
	}
}

// ユーティリティ関数
func isOggFormat(data []byte) bool {
	// OGGファイルのマジックナンバーをチェック
	if len(data) < 4 {
		return false
	}
	return string(data[:4]) == "OggS"
}

type oggAdapter struct {
	*oggvorbis.Reader
}

func (o *oggAdapter) Read(p []byte) (int, error) {
	samples := make([]float32, len(p)/4)
	n, err := o.Reader.Read(samples)
	if err != nil {
		return 0, err
	}

	for i, s := range samples[:n] {
		binary.LittleEndian.PutUint32(p[i*4:], math.Float32bits(s))
	}
	return n * 4, nil
}
