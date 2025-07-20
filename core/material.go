package core

import (
	"RayTracingInOneWeekend/utils"
	"math"
)

// MaterialI 材质接口
type MaterialI interface {
	Scatter(r *Ray, h HitRecord) (hit bool, attenuation Color, scattered *Ray) // 散射    还需要实现接口的材质提供衰减后的颜色attenuation，和散射出的新射线
}

/*
现在我们有了物体和多束每像素光线，我们可以制作一些看起来更真实的材质。
我们将从漫反射材质（也称为哑光材质）开始。
一个问题是我们是否要混合搭配几何体和材质（这样我们就可以将材质分配给多个球体，反之亦然），
或者几何体和材质是紧密绑定的（这对于几何体和材质相关的程序化对象可能很有用）。我们将选择分离——这是大多数渲染器采用的方式——但请注意，还有其他替代方法。
*/

/*
漫反射物体不会自身发光，它们只是吸收周围环境的光，但会通过自身的固有颜色进行调节。
光线在漫反射表面上的反射方向是随机的，所以，如果向两个漫反射表面之间的缝隙发射三束光线，
它们将各自具有不同的随机行为（Vec3.Random）

如果一条射线从材料上反射并保持 100%的颜色，那么我们就说这个材料是白色的。
如果一条射线从材料上反射并保持 0%的颜色，那么我们就说这个材料是黑色的。
作为我们新漫反射材料的一个初步演示，我们将 ray_color 函数设置为返回反弹时 50%的颜色。
我们预计会得到一个漂亮的灰色。
*/
type DiffuseMaterial struct{}

// SimpleDiffuseMaterial 简易漫反射材质
type SimpleDiffuseMaterial struct{}

/*
在半球上均匀地散射反射光线可以产生一个柔和的漫反射模型，但我们当然可以做得更好。
对真实漫反射物体更准确的表示是朗伯分布。这种分布以与反射光线和表面法线之间的角度 cos(ϕ)成正比的方式散射反射光线，
其中 ϕ是反射光线和表面法线之间的角度。这意味着反射光线最有可能散射到接近表面法线的方向，而在远离法线的方向散射的可能性较小。
这种非均匀的朗伯分布比我们之前的均匀散射更好地模拟了现实世界中的材料反射。

一个球将沿表面的法向量方向（n）偏移，另一个球将沿相反方向（−n）偏移。
这给我们留下了两个单位大小的球，它们只会在交点处接触表面。由此，一个球的中心位于 (P+n)，
另一个球的中心位于 (P−n)。中心位于 (P−n)的球被认为是表面内部，而中心为 (P+n)的球被认为是表面外部。
*/

type Material struct {
}

/*
朗伯（漫反射）反射率可以始终根据其反射率R散射和衰减光线，
或者它有时可以以概率1−R散射而不衰减（其中未散射的光线只是被材料吸收）。
它也可以是这两种策略的混合。
我们将选择始终散射，因此实现朗伯材料就变成一项简单的任务：
*/
type LambertianReflectionMaterial struct {
	Albedo Color    // 反照率
	Tex    TextureI // 材质
}

func (l LambertianReflectionMaterial) Scatter(r *Ray, hitRecord HitRecord) (hit bool, attenuation Color, scattered *Ray) {
	scatterDirection := hitRecord.Normal.Add(RandomNormalizedVec3())
	if scatterDirection.NearZero() {
		scatterDirection = hitRecord.Normal
	}
	scattered = &Ray{hitRecord.HitPoint, scatterDirection, r.Time()}

	//
	if l.Tex != nil {
		attenuation = l.Tex.Value(hitRecord.U, hitRecord.V, hitRecord.HitPoint)
	} else {
		attenuation = l.Albedo
	}
	return true, attenuation, scattered
}

/*
Fuzz:
通过使用一个小球来随机化反射方向，为光线选择一个新的端点。
使用以原始端点为中心的球面上的一个随机点，该点按毛糙因子缩放。
模糊球体越大，反射就越模糊。这表明需要添加一个模糊度参数，该参数就是球的半径（因此零表示无扰动）。
问题是对于大球体或掠射光线，可能会在表面以下散射。我们可以让表面吸收这些。
为了使模糊球体有意义，它需要与反射向量保持一致的缩放比例，而反射向量的长度可以是任意的。
为了解决这个问题，我们需要对反射光线进行归一化。
*/
type MetalMaterial struct {
	Albedo Color
	Fuzz   float64 // 模糊度
}

func (m MetalMaterial) Scatter(r *Ray, h HitRecord) (hit bool, attenuation Color, scattered *Ray) {
	directReflected := r.Direction.Reflect(h.Normal) // 反射光线方向
	reflected := directReflected.Normalize().Add(RandomNormalizedVec3().MultiplicationNum(m.Fuzz))
	scattered = &Ray{h.HitPoint, reflected, r.Time()} // 反射光线
	attenuation = m.Albedo
	hit = scattered.Direction.Dot(h.Normal) > 0 // 判断是否反射光线与入射点法线同向，否的话无法进行下次
	return hit, attenuation, scattered
}

/*
折射
*/
type DielectricMaterial struct {
	RefractionIndex float64 //  真空或空气中的折射率，或材料的折射率与封闭介质的折射率之比
}

func (d DielectricMaterial) Scatter(r *Ray, h HitRecord) (hit bool, attenuation Color, scattered *Ray) {
	attenuation = Color{1.0, 1.0, 1.0}
	ri := d.RefractionIndex
	if h.FrontFace {
		ri = 1.0 / d.RefractionIndex
	}
	rayInDirectionNormalized := r.Direction.Normalize() // 入射线单位向量
	cosTheta := math.Min(rayInDirectionNormalized.MultiplicationNum(-1.0).Dot(h.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)
	// 判断是反射还是折射
	cannotRefract := ri*sinTheta > 1.0
	var refractedDirection Vec3
	// 是否在里面
	if cannotRefract || Reflectance(cosTheta, ri) > utils.Random() {
		// 折射
		refractedDirection = rayInDirectionNormalized.Reflect(h.Normal)
	} else {
		// 折射
		refractedDirection = rayInDirectionNormalized.Refract(h.Normal, ri) // 折射线方向
	}
	scattered = &Ray{h.HitPoint, refractedDirection, r.Time()} //折射或反射线
	return true, attenuation, scattered
}

func Reflectance(cosine, refractionIndex float64) float64 {
	r0 := (1 - refractionIndex) / (1 + refractionIndex)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow(1-cosine, 5)
}
