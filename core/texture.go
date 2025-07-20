package core

import "math"

type TextureI interface {
	Value(u, v float64, p Point) Color
}

type SolidColorTexture struct {
	albedo Color
}

func NewSolidColorTexture(albedo Color) *SolidColorTexture {
	return &SolidColorTexture{albedo}
}

func (s SolidColorTexture) Value(u, v float64, p Point) Color {
	return s.albedo
}

type CheckerTexture struct {
	invScale  float64
	even, odd TextureI
}

func NewCheckerTexture(invScale float64, even, odd Color) *CheckerTexture {
	return &CheckerTexture{
		invScale: invScale,
		even:     NewSolidColorTexture(even),
		odd:      NewSolidColorTexture(odd),
	}
}

func (c CheckerTexture) Value(u, v float64, p Point) Color {
	xInteger := int(math.Floor(c.invScale * p.X))
	yInteger := int(math.Floor(c.invScale * p.Y))
	zInteger := int(math.Floor(c.invScale * p.Z))
	isEven := (xInteger+yInteger+zInteger)%2 == 0
	if isEven {
		return c.even.Value(u, v, p)
	} else {
		return c.odd.Value(u, v, p)
	}
}
