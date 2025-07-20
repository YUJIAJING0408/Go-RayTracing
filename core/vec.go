// 核心对象
package core

import (
	"RayTracingInOneWeekend/utils"
	"math"
)

type Vec3 struct {
	X, Y, Z float64
}

func NewVec3(x, y, z float64) Vec3 {
	return Vec3{x, y, z}
}

// Add 加
func (v Vec3) Add(other Vec3) Vec3 {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	return v
}

// Sub 减
func (v Vec3) Sub(other Vec3) Vec3 {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	return v
}

// MultiplicationNum 乘数字
func (v Vec3) MultiplicationNum(num float64) Vec3 {
	v.X *= num
	v.Y *= num
	v.Z *= num
	return v
}

// MultiplicationVec3 与另一个向量逐元素相乘
func (v Vec3) MultiplicationVec3(other Vec3) Vec3 {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	return v
}

// Div 除法
func (v Vec3) Div(num float64) Vec3 {
	return v.MultiplicationNum(1 / num)
}

// Dot 点乘
func (v Vec3) Dot(other Vec3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross 叉乘
func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// LengthSquared 向量平方距离
func (v Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// Length 向量长度
func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// Normalize 获得单位向量
func (v Vec3) Normalize() Vec3 {
	return v.Div(v.Length())
}

func (v Vec3) IsNormalized() bool {
	if v.LengthSquared() == 1 {
		return true
	} else {
		return false
	}
}

func Random() Vec3 {
	return Vec3{utils.Random(), utils.Random(), utils.Random()}
}

func RandomBetween(min, max float64) Vec3 {
	return Vec3{utils.RandomBetween(min, max), utils.RandomBetween(min, max), utils.RandomBetween(min, max)}
}

func RandomNormalizedVec3() Vec3 {
	for {
		var p = RandomBetween(-1.0, 1.0)
		var l = p.Length()
		if l >= 1e-160 && l <= 1.0 {
			return p.Div(l)
		}
	}
}

// RandomOnHemisphere 通过计算表面法向量和随机向量的点积来判断它是否位于正确的半球。如果点积为正，则向量位于正确的半球。如果点积为负，则需要反转向量。
func RandomOnHemisphere(normal Vec3) Vec3 {
	onNormalizedVec3 := RandomNormalizedVec3()
	if onNormalizedVec3.Dot(normal) < 0 {
		return onNormalizedVec3.MultiplicationNum(-1.0)
	} else {
		return onNormalizedVec3
	}
}

// NearZero 如果向量在所有维度上都非常接近零，则返回 true。
func (v Vec3) NearZero() bool {
	return v.X <= 1e-8 && v.Y <= 1e-8 && v.Z <= 1e-8
}

// Reflect v是入射光线，n是法线，b是单位化后的n， v+2b就得到反射光线, 避免n不是单位向量
func (v Vec3) Reflect(n Vec3) Vec3 {
	return v.Sub(n.MultiplicationNum(v.Dot(n) * 2.0))
}

// Refract 折射由斯涅尔定律描述
func (v Vec3) Refract(n Vec3, etaiOverEtat float64) Vec3 {
	cosTheta := math.Min(v.MultiplicationNum(-1.0).Dot(n), 1.0)
	rOutPerp := v.Add(n.MultiplicationNum(cosTheta)).MultiplicationNum(etaiOverEtat) // 垂直
	rOutParallel := n.MultiplicationNum(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared())))
	return rOutPerp.Add(rOutParallel)
}

// 在单位盘内生成随机点,长度小于1的结果保留
func RandomInUnitDisk() Vec3 {
	for {
		p := Vec3{utils.RandomBetween(-1, 1), utils.RandomBetween(-1, 1), 0}
		if p.LengthSquared() < 1 {
			return p
		}
	}
}

// 点
type Point Vec3
