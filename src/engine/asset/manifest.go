package asset

import "encoding/json"

// アセットマニフェスト
type AssetManifest struct {
	Images   map[string]ImageAssetInfo   `json:"images"`
	Audio    map[string]AudioAssetInfo   `json:"audio"`
	Fonts    map[string]FontAssetInfo    `json:"fonts"`
	Scripts  map[string]ScriptAssetInfo  `json:"scripts"`
}

type ImageAssetInfo struct {
	Path string `json:"path"`
}

type AudioAssetInfo struct {
	Path string `json:"path"`
	Type string `json:"type"` // "bgm" or "se"
}

type FontAssetInfo struct {
	Path string  `json:"path"`
	Size float64 `json:"size"`
}

type ScriptAssetInfo struct {
	Path string `json:"path"`
}

// マニフェストのロード
func LoadManifest(data []byte) (*AssetManifest, error) {
	var manifest AssetManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, err
	}
	return &manifest, nil
} 