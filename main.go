package main

import (
	"RayTracingInOneWeekend/core"
	"RayTracingInOneWeekend/utils"
)

func main() {
	// 材质
	// 网格纹理
	checkerTexture := core.NewCheckerTexture(0.32, core.Color{X: .2, Y: .3, Z: .1}, core.Color{X: .9, Y: .9, Z: .9})
	var groundMaterial = core.LambertianReflectionMaterial{Albedo: core.Color{X: 0.5, Y: 0.5, Z: 0.5}, Tex: checkerTexture}
	groundSphere := core.NewSphere(core.Point{Y: -1000}, 1000).WithMaterial(groundMaterial)
	s2 := core.NewSphere(core.Point{X: 0.0, Y: 0.0, Z: -1.0}, 0.5).WithMaterial(core.LambertianReflectionMaterial{Albedo: core.Color{X: 0.1, Y: 0.2, Z: 0.5}})
	s3 := core.NewSphere(core.Point{X: -1.0, Z: -1.0}, 0.5).WithMaterial(core.DielectricMaterial{RefractionIndex: 1.5})
	s4 := core.NewSphere(core.Point{X: 1.0, Z: -1.0}, 0.5).WithMaterial(core.MetalMaterial{Albedo: core.Color{X: 0.8, Y: 0.6, Z: 0.2}, Fuzz: 1.0})
	// 相机构建
	camera := core.NewCamera(core.Point{}, core.Point{X: 13, Y: 2, Z: 3}, 16.0/9.0, 20, 400, 512, 10, true, 0.6, 10.0)
	camera.Add(groundSphere, s2, s3, s4)
	// 随机构建场景
	for i := -8; i < 8; i++ {
		for j := -8; j < 8; j++ {
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
			center := core.Point{X: float64(i) + 0.9*utils.Random(), Y: R, Z: float64(j) + 0.9*utils.Random()}
			if core.Vec3(center).Sub(core.Vec3{X: 4.0, Y: 2.0}).Length() > 0.9 {
				var material core.MaterialI
				albedo := core.Color(core.Random())
				//
				switch {
				case chooseMat < 0.8:
					// diffuse
					material = core.LambertianReflectionMaterial{Albedo: albedo}
				case chooseMat < 0.90:
					// metal
					fuzz := utils.Random()
					material = core.MetalMaterial{
						Albedo: albedo,
						Fuzz:   fuzz,
					}
				default:
					// glass
					material = core.DielectricMaterial{
						RefractionIndex: 1.5,
					}
				}
				sphere := core.NewSphere(center, R).WithMaterial(material)
				if utils.Random() <= 0.3 {
					sphere.SetUniformLinearMovement(core.Point{Y: utils.RandomBetween(0, 0.3)})
				}
				camera.Add(sphere)
			}
		}
	}
	//camera.Render("o")
	camera.MultithreadedRender("mo", 8, 1000000)
	//camera.EnabledBVH(true)
	// BVH渲染，不一定更快
	//camera.Render("o_bvh")
	//camera.MultithreadedRender("mo_bvh", 6, 1000000)
}
