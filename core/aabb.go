package core

type AABB struct {
	X, Y, Z Interval
}

func NewAABB(x, y, z Interval) *AABB {
	return &AABB{x, y, z}
}

func NewAABBFromAABB(aabb1 *AABB, aabb2 *AABB) *AABB {
	return &AABB{
		NewIntervalFromInterval(aabb1.AxisInterval(0), aabb2.AxisInterval(0)),
		NewIntervalFromInterval(aabb1.AxisInterval(1), aabb2.AxisInterval(1)),
		NewIntervalFromInterval(aabb1.AxisInterval(2), aabb2.AxisInterval(2)),
	}
}

func NewAABBFromPoints(a, b Point) *AABB {
	var x, y, z Interval
	if a.X <= b.X {
		x = Interval{a.X, b.X}
	} else {
		x = Interval{b.X, a.X}
	}
	if a.Y <= b.Y {
		y = Interval{a.Y, b.Y}
	} else {
		y = Interval{b.Y, a.Y}
	}
	if a.Z <= b.Z {
		z = Interval{a.Z, b.Z}
	} else {
		z = Interval{b.Z, a.Z}
	}
	return &AABB{x, y, z}
}

func (aabb AABB) AxisInterval(n int) *Interval {
	switch n {
	case 0:
		return &aabb.X
	case 1:
		return &aabb.Y
	case 2:
		return &aabb.Z
	default:
		panic("bad axis")
	}
}

func (aabb AABB) Hit(r *Ray) bool {
	var rayInterval Interval
	var rayOrigin = r.Origin
	var rayDirection = r.Direction
	for i := 0; i < 3; i++ {
		// 获取轴
		axisInterval := aabb.AxisInterval(i)
		adinv, t0, t1 := 0.0, 0.0, 0.0
		switch i {
		case 0:
			adinv = 1.0 / rayDirection.X
			t0 = (axisInterval.Min - rayOrigin.X) * adinv
			t1 = (axisInterval.Max - rayOrigin.X) * adinv

		case 1:
			adinv = 1.0 / rayDirection.Y
			t0 = (axisInterval.Min - rayOrigin.Y) * adinv
			t1 = (axisInterval.Max - rayOrigin.Y) * adinv
		default:
			adinv = 1.0 / rayDirection.Z
			t0 = (axisInterval.Min - rayOrigin.Z) * adinv
			t1 = (axisInterval.Max - rayOrigin.Z) * adinv
		}
		if t0 < t1 {
			if t0 > rayInterval.Min {
				rayInterval.Min = t0
			}
			if t1 < rayInterval.Max {
				rayInterval.Max = t1
			}
		} else {
			if t0 > rayInterval.Max {
				rayInterval.Max = t0
			}
			if t1 > rayInterval.Min {
				rayInterval.Min = t1
			}
		}
		if rayInterval.Max < rayInterval.Min {
			return false
		}
	}
	return true
}

func (aabb AABB) LongestAxis() (index int) {
	if aabb.X.Size() > aabb.Y.Size() {
		if aabb.X.Size() > aabb.Z.Size() {
			index = 0
		} else {
			index = 2
		}
	} else {
		if aabb.Y.Size() > aabb.Z.Size() {
			index = 1
		} else {
			index = 2
		}
	}
	return index
}
