package render

import (
	"sort"
)

type RenderCache struct {
	drawOrder []string
	dirty     bool
}

func NewRenderCache() *RenderCache {
	return &RenderCache{
		drawOrder: make([]string, 0),
		dirty:     true,
	}
}

func (c *RenderCache) Invalidate() {
	c.dirty = true
}

func (c *RenderCache) GetDrawOrder(drawables map[string]*DrawableObject) []string {
	if !c.dirty {
		return c.drawOrder
	}

	// 描画順を再計算
	objects := make([]struct {
		Name   string
		ZIndex int
	}, 0, len(drawables))

	for name, obj := range drawables {
		objects = append(objects, struct {
			Name   string
			ZIndex int
		}{name, obj.ZIndex})
	}

	sort.Slice(objects, func(i, j int) bool {
		return objects[i].ZIndex < objects[j].ZIndex
	})

	c.drawOrder = make([]string, len(objects))
	for i, obj := range objects {
		c.drawOrder[i] = obj.Name
	}

	c.dirty = false
	return c.drawOrder
} 