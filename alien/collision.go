package alien

// CheckCollision 检查子弹和外星人之间是否有碰撞，任意一个角碰到了就算碰撞
func CheckCollision(entityA, entityB Entity) bool {
	//entityA := entityAA.(Entity)
	//entityB := entityBB.(Entity)
	top, left := entityB.Y(), entityB.X()                                                        //344 220
	bottom, right := entityB.Y()+float64(entityB.Height()), entityB.X()+float64(entityB.Width()) //
	// 左上角
	x, y := entityA.X(), entityA.Y()
	if y >= top && y <= bottom && x >= left && x <= right {
		return true
	}

	// 右上角
	x, y = entityA.X()+float64(entityA.Width()), entityA.Y()
	if y >= top && y <= bottom && x >= left && x <= right {
		return true
	}

	// 左下角
	x, y = entityA.X(), entityA.Y()+float64(entityA.Height())
	if y >= top && y <= bottom && x >= left && x <= right {
		return true
	}

	// 右下角
	x, y = entityA.X()+float64(entityA.Width()), entityA.Y()+float64(entityA.Height())
	if y >= top && y <= bottom && x >= left && x <= right {
		return true
	}

	return false
}

type GameObject struct {
	width  int
	height int
	x      float64
	y      float64
}

func (gameObj *GameObject) Width() int {
	return gameObj.width
}

func (gameObj *GameObject) Height() int {
	return gameObj.height
}

func (gameObj *GameObject) X() float64 {
	return gameObj.x
}

func (gameObj *GameObject) Y() float64 {
	return gameObj.y
}

type Entity interface {
	Width() int
	Height() int
	X() float64
	Y() float64
}
