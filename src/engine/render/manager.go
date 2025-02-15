package render

import (
	"sync"
	"github.com/hajimehoshi/ebiten/v2"
)

type DrawableObject struct {
	Image     *ebiten.Image
	X, Y      float64
	ZIndex    int
	Visible   bool
	Options   *ebiten.DrawImageOptions
}

type RenderManager struct {
	mutex     sync.RWMutex
	drawables map[string]*DrawableObject
	cache     *RenderCache
}

func NewRenderManager() *RenderManager {
	return &RenderManager{
		drawables: make(map[string]*DrawableObject),
		cache:     NewRenderCache(),
	}
}

func (r *RenderManager) RegisterDrawable(name string, obj *DrawableObject) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.drawables[name] = obj
	r.cache.Invalidate()
}

func (r *RenderManager) Draw(screen *ebiten.Image) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// キャッシュされた描画順を使用
	drawOrder := r.cache.GetDrawOrder(r.drawables)
	
	for _, name := range drawOrder {
		obj := r.drawables[name]
		if obj.Visible {
			screen.DrawImage(obj.Image, obj.Options)
		}
	}
} 