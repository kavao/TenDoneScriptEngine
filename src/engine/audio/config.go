package audio

type AudioConfig struct {
	BGMVolume float64 `json:"bgm_volume"`
	SEVolume  float64 `json:"se_volume"`
}

func NewAudioConfig() *AudioConfig {
	return &AudioConfig{
		BGMVolume: 1.0,
		SEVolume:  1.0,
	}
}

func (c *AudioConfig) SetBGMVolume(volume float64) {
	if volume < 0.0 {
		volume = 0.0
	}
	if volume > 1.0 {
		volume = 1.0
	}
	c.BGMVolume = volume
}

func (c *AudioConfig) SetSEVolume(volume float64) {
	if volume < 0.0 {
		volume = 0.0
	}
	if volume > 1.0 {
		volume = 1.0
	}
	c.SEVolume = volume
} 