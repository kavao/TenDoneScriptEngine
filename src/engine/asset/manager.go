package asset

import (
	"fmt"
	"gameengine/src/engine/audio"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

// アセットの種類
type AssetType int

const (
	AssetTypeImage AssetType = iota
	AssetTypeAudio
	AssetTypeFont
	AssetTypeScript
)

// アセット情報
type AssetInfo struct {
	Type     AssetType
	Path     string
	LoadFunc func() (interface{}, error)
}

// アセットマネージャー
type AssetManager struct {
	mutex       sync.RWMutex
	assets      map[string]interface{}
	assetInfo   map[string]AssetInfo
	audioMgr    *audio.AudioManager
	loadingChan chan string
}

func NewAssetManager(audioMgr *audio.AudioManager) *AssetManager {
	return &AssetManager{
		assets:      make(map[string]interface{}),
		assetInfo:   make(map[string]AssetInfo),
		audioMgr:    audioMgr,
		loadingChan: make(chan string, 100),
	}
}

// アセットの登録
func (m *AssetManager) RegisterAsset(id string, info AssetInfo) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.assetInfo[id] = info
}

// アセットの非同期ロード
func (m *AssetManager) LoadAsync(ids ...string) {
	for _, id := range ids {
		m.loadingChan <- id
	}
}

// アセットのロード状態を更新
func (m *AssetManager) Update() error {
	select {
	case id := <-m.loadingChan:
		if err := m.loadAsset(id); err != nil {
			return fmt.Errorf("failed to load asset %s: %v", id, err)
		}
	default:
		// ロード待ちなし
	}
	return nil
}

// アセットのロード
func (m *AssetManager) loadAsset(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	info, exists := m.assetInfo[id]
	if !exists {
		return fmt.Errorf("asset not registered: %s", id)
	}

	asset, err := info.LoadFunc()
	if err != nil {
		return err
	}

	m.assets[id] = asset
	return nil
}

// 画像アセットの取得
func (m *AssetManager) GetImage(id string) (*ebiten.Image, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	asset, exists := m.assets[id]
	if !exists {
		return nil, fmt.Errorf("image asset not loaded: %s", id)
	}

	img, ok := asset.(*ebiten.Image)
	if !ok {
		return nil, fmt.Errorf("asset is not an image: %s", id)
	}

	return img, nil
}

// フォントアセットの取得
func (m *AssetManager) GetFont(id string) (font.Face, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	asset, exists := m.assets[id]
	if !exists {
		return nil, fmt.Errorf("font asset not loaded: %s", id)
	}

	f, ok := asset.(font.Face)
	if !ok {
		return nil, fmt.Errorf("asset is not a font: %s", id)
	}

	return f, nil
}

// スクリプトアセットの取得
func (m *AssetManager) GetScript(id string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	asset, exists := m.assets[id]
	if !exists {
		return "", fmt.Errorf("script asset not loaded: %s", id)
	}

	script, ok := asset.(string)
	if !ok {
		return "", fmt.Errorf("asset is not a script: %s", id)
	}

	return script, nil
}

// アセットの解放
func (m *AssetManager) UnloadAsset(id string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	delete(m.assets, id)
}

// 全アセットの解放
func (m *AssetManager) UnloadAll() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.assets = make(map[string]interface{})
}
