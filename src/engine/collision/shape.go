package collision

// 形状の種類
type ShapeType int

const (
	ShapeTypeBox ShapeType = iota
	ShapeTypeCircle
	ShapeTypePolygon
)

// 基本的な形状インターフェース
type Shape interface {
	Type() ShapeType
	GetBounds() Bounds
	GetCenter() Vector2D
	SetPosition(x, y float64)
	Intersects(other Shape) bool
}

// 2D座標
type Vector2D struct {
	X, Y float64
}

// バウンディングボックス
type Bounds struct {
	Min, Max Vector2D
}

// 矩形の衝突判定形状
type BoxShape struct {
	Position Vector2D
	Width    float64
	Height   float64
}

func NewBoxShape(x, y, width, height float64) *BoxShape {
	return &BoxShape{
		Position: Vector2D{x, y},
		Width:    width,
		Height:   height,
	}
}

func (b *BoxShape) Type() ShapeType {
	return ShapeTypeBox
}

func (b *BoxShape) GetBounds() Bounds {
	return Bounds{
		Min: Vector2D{b.Position.X, b.Position.Y},
		Max: Vector2D{b.Position.X + b.Width, b.Position.Y + b.Height},
	}
}

func (b *BoxShape) GetCenter() Vector2D {
	return Vector2D{
		X: b.Position.X + b.Width/2,
		Y: b.Position.Y + b.Height/2,
	}
}

func (b *BoxShape) SetPosition(x, y float64) {
	b.Position.X = x
	b.Position.Y = y
}

func (b *BoxShape) Intersects(other Shape) bool {
	return CheckCollision(b, other)
}

// 円形の衝突判定形状
type CircleShape struct {
	Center Vector2D
	Radius float64
}

func NewCircleShape(x, y, radius float64) *CircleShape {
	return &CircleShape{
		Center: Vector2D{x, y},
		Radius: radius,
	}
}

func (c *CircleShape) Type() ShapeType {
	return ShapeTypeCircle
}

func (c *CircleShape) GetBounds() Bounds {
	return Bounds{
		Min: Vector2D{c.Center.X - c.Radius, c.Center.Y - c.Radius},
		Max: Vector2D{c.Center.X + c.Radius, c.Center.Y + c.Radius},
	}
}

func (c *CircleShape) GetCenter() Vector2D {
	return c.Center
}

func (c *CircleShape) SetPosition(x, y float64) {
	c.Center.X = x
	c.Center.Y = y
}

func (c *CircleShape) Intersects(other Shape) bool {
	return CheckCollision(c, other)
}
