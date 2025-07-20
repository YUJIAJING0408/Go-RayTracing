package core

import (
	"RayTracingInOneWeekend/utils"
	"math"
)

type ColorI interface {
	Color2Pixel() utils.Pixel
}

type ColorTask struct {
	R           *Ray
	WidthIndex  int
	HeightIndex int
}

type ColorRes struct {
	C           Color
	WidthIndex  int
	HeightIndex int
	Index       int
}

// [0.0,1.0]
type Color Vec3

func (c *Color) Color2Pixel() utils.Pixel {
	c.Linear2Gamma(2.0)
	var interval = Interval{0, 1.0}
	return utils.Pixel{
		R: int(interval.Clamp(c.X) * 255),
		G: int(interval.Clamp(c.Y) * 255),
		B: int(interval.Clamp(c.Z) * 255),
	}
}

func (c *Color) Linear2Gamma(gamma float64) {
	if gamma == 0 {
		gamma = 2.0
	}
	c.X = math.Pow(c.X, 1.0/gamma)
	c.Y = math.Pow(c.Y, 1.0/gamma)
	c.Z = math.Pow(c.Z, 1.0/gamma)
}
