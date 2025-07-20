package core

import (
	"RayTracingInOneWeekend/utils"
	"testing"
)

func TestRender(t *testing.T) {
	// 材质
	groundMaterial := LambertianReflectionMaterial{Albedo: Color{0.5, 0.5, 0.5}}
	groundSphere := NewSphere(Point{0, -1000, 0}, 1000).WithMaterial(groundMaterial)
	s2 := NewSphere(Point{X: 0.0, Y: 0.0, Z: -1.0}, 0.5).WithMaterial(LambertianReflectionMaterial{Albedo: Color{X: 0.1, Y: 0.2, Z: 0.5}})
	s3 := NewSphere(Point{X: -1.0, Z: -1.0}, 0.5).WithMaterial(DielectricMaterial{RefractionIndex: 1.5})
	s4 := NewSphere(Point{X: 1.0, Z: -1.0}, 0.5).WithMaterial(MetalMaterial{Albedo: Color{X: 0.8, Y: 0.6, Z: 0.2}, Fuzz: 1.0})
	// 相机构建
	camera := NewCamera(Point{0, 0, 0}, Point{13, 2, 3}, 16.0/9.0, 20, 160, 32, 10, true, 0.6, 10.0)
	camera.Add(groundSphere, s2, s3, s4)
	// 随机构建场景
	for i := -10; i < 10; i++ {
		for j := -10; j < 10; j++ {
			chooseSize := utils.Random()
			R := 0.0
			switch {
			case chooseSize < 0.8:
				R = utils.RandomBetween(0.1, 0.25)
			case chooseSize < 0.95:
				R = utils.RandomBetween(0.25, 0.4)
			default:
				R = utils.RandomBetween(0.4, 0.6)

			}
			chooseMat := utils.Random()
			center := Point{float64(i) + 0.9*utils.Random(), R, float64(j) + 0.9*utils.Random()}
			if Vec3(center).Sub(Vec3{4.0, 2.0, 0.0}).Length() > 0.9 {
				var material MaterialI
				albedo := Color(Random())
				//
				switch {
				case chooseMat < 0.8:
					// diffuse
					material = LambertianReflectionMaterial{Albedo: albedo}
				case chooseMat < 0.90:
					// metal
					fuzz := utils.Random()
					material = MetalMaterial{
						Albedo: albedo,
						Fuzz:   fuzz,
					}
				default:
					// glass
					material = DielectricMaterial{
						RefractionIndex: 1.5,
					}
				}
				sphere := NewSphere(center, R).WithMaterial(material)
				if utils.Random() <= 0.3 {
					sphere.SetUniformLinearMovement(Point{0, utils.RandomBetween(0, 0.3), 0})
				}
				camera.Add(sphere)
			}
		}
	}
	//camera.EnabledBVH(true)
	// 渲染
	camera.MultithreadedRender("mo", 10, 1000000)
	//camera.Render("o")
}
