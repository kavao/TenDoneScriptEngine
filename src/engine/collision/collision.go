package collision

// 衝突判定関数
func CheckCollision(a, b Shape) bool {
	// まずバウンディングボックスで大まかな判定
	if !checkBoundsOverlap(a.GetBounds(), b.GetBounds()) {
		return false
	}

	// 形状ごとの詳細な判定
	switch a.Type() {
	case ShapeTypeBox:
		switch b.Type() {
		case ShapeTypeBox:
			return checkBoxBox(a.(*BoxShape), b.(*BoxShape))
		case ShapeTypeCircle:
			return checkBoxCircle(a.(*BoxShape), b.(*CircleShape))
		}
	case ShapeTypeCircle:
		switch b.Type() {
		case ShapeTypeBox:
			return checkBoxCircle(b.(*BoxShape), a.(*CircleShape))
		case ShapeTypeCircle:
			return checkCircleCircle(a.(*CircleShape), b.(*CircleShape))
		}
	}

	return false
}

// バウンディングボックスの重なり判定
func checkBoundsOverlap(a, b Bounds) bool {
	return !(a.Max.X < b.Min.X || a.Min.X > b.Max.X ||
		a.Max.Y < b.Min.Y || a.Min.Y > b.Max.Y)
}

// 矩形同士の衝突判定
func checkBoxBox(a, b *BoxShape) bool {
	boundsA := a.GetBounds()
	boundsB := b.GetBounds()
	return checkBoundsOverlap(boundsA, boundsB)
}

// 円同士の衝突判定
func checkCircleCircle(a, b *CircleShape) bool {
	dx := a.Center.X - b.Center.X
	dy := a.Center.Y - b.Center.Y
	distSq := dx*dx + dy*dy
	return distSq <= (a.Radius+b.Radius)*(a.Radius+b.Radius)
}

// 矩形と円の衝突判定
func checkBoxCircle(box *BoxShape, circle *CircleShape) bool {
	// 円の中心と矩形の最近接点を求める
	closestX := clamp(circle.Center.X, box.Position.X, box.Position.X+box.Width)
	closestY := clamp(circle.Center.Y, box.Position.Y, box.Position.Y+box.Height)

	// 最近接点と円の中心との距離を計算
	dx := circle.Center.X - closestX
	dy := circle.Center.Y - closestY
	distSq := dx*dx + dy*dy

	return distSq <= circle.Radius*circle.Radius
}

// 値を範囲内に制限するヘルパー関数
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
