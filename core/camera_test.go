package core

import (
	"math"
	"testing"
)

func TestCamera(t *testing.T) {
	camera := NewCamera(Point{0, 0, 0}, Point{0, 0, -1}, 16.0/9.0, 90, 160, 100, 20, true, 0.4, 10)
	R := math.Cos(math.Pi / 4)
	s1 := Sphere{
		Center:   Point{X: -R, Z: -1.0},
		Radius:   R,
		Material: LambertianReflectionMaterial{Albedo: Color{X: 1.0}},
	}
	s2 := Sphere{
		Center:   Point{X: R, Z: -1.0},
		Radius:   R,
		Material: LambertianReflectionMaterial{Albedo: Color{Z: 1.0}},
	}
	camera.Add(&s1, &s2)
	camera.Render()
	//camera.MultithreadedRender(12, 10000000)
}
