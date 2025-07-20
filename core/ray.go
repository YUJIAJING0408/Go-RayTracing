package core

type Ray struct {
	Origin    Point // 起点
	Direction Vec3  // 方向,单位向量
	TM        float64
}

func NewRay(origin Point, direction Vec3) Ray {
	return Ray{origin, direction.Normalize(), 0.0}
}

func NewRayWithTime(origin Point, direction Vec3, tm float64) Ray {
	return Ray{origin, direction.Normalize(), tm}
}

func (r Ray) Time() float64 {
	return r.TM
}

// At 某时刻射线所在位置
func (r Ray) At(t float64) Point {
	return Point{
		X: r.Origin.X + r.Direction.X*t,
		Y: r.Origin.Y + r.Direction.Y*t,
		Z: r.Origin.Z + r.Direction.Z*t,
	}
}
