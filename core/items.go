package core

import "math"

type RemovableI interface {
	NowAt(time float64) Point // 获取任意位置的中心
}

type RemovableSetting struct {
	MovingFunc  func(t float64) Point
	IsRemovable bool
}

type HitRecord struct {
	Time      float64   // 负根
	HitPoint  Point     // 交叉点
	Material  MaterialI // 击中点材质
	Normal    Vec3      // 交叉点法线
	FrontFace bool      // 法线方向
	U, V      float64   // UV坐标
}

type HittableItemI interface {
	Hittable(ray Ray, rayT Interval) (hit bool, hitRecord HitRecord)
	SetBoundingBox(aabb *AABB)
	GetBoundingBox() *AABB
}

type Scenes struct {
	HittableList []HittableItemI
	HittableAABB *AABB
	EnabledBVH   bool
	BVH          *BVHNode
}

func (s *Scenes) Len() int {
	return len(s.HittableList)
}

func (s *Scenes) Add(items ...HittableItemI) {
	for _, item := range items {
		s.HittableList = append(s.HittableList, item)
		s.HittableAABB = NewAABBFromAABB(s.HittableAABB, item.GetBoundingBox())
	}
}

func (s *Scenes) GetBoundingBox() *AABB {
	return s.HittableAABB
}

func (s *Scenes) HitAnything(ray Ray, rayT Interval) (hitAnything bool, record HitRecord) {
	closestSoFar := rayT.Max
	for _, hittableItem := range s.HittableList {
		// 只保留最近的那个
		if hit, hitRecord := hittableItem.Hittable(ray, rayT); hit {
			// 打击时间小于才有效，大于直接抛弃
			if hitRecord.Time < closestSoFar {
				hitAnything = true
				closestSoFar = hitRecord.Time
				record = hitRecord
			}
		} else {
			rayT.Max = closestSoFar
		}
	}
	return hitAnything, record
}

/*
圆球定义：
半径为 r,且以原点为中心的球体的方程是一个重要的数学方程：x^2 + y^2 + z^2 = r^2
你也可以这样理解：
即如果给定的点 (x,y,z)在球面上，那么 x^2 + y^2 + z^2 = r^2，
如果给定的点 (x,y,z)在球体内，那么 x^2 + y^2 + z^2 < r^2，
如果给定的点 (x,y,z)在球体外，那么 x^2 + y^2 + z^2 > r^2。
如果我们想让球心位于任意点 (Cx,Cy,Cz)，那么方程就变得不那么简洁了：
(Cx−x)^2 + (Cy−y)^2 + (Cz−z)^2 = r^2
在图形学中，你几乎总是希望你的公式用向量表示，这样所有的 x/y/z 这些内容都可以简单地用 vec3 类来表示。
你可能注意到，从点 P=(x,y,z)到中心C=(Cx,Cy,Cz)的向量是(C−P)。
如果我们使用点积的定义：(C−P) ⋅ (C−P) = (Cx−x)^2 + (Cy−y)^2 + (Cz−z)^2
然后我们可以将球的方程写成向量形式：(C−P) ⋅ (C−P) = r^2
我们可以将其理解为“任何满足此方程的点P都位于球面上。
我们想知道我们的射线P(t) = Q + td是否会与球体相交。
如果它与球体相交，那么就存在某个t，使得P(t)满足球体方程。
因此，我们正在寻找任何t，使得此条件成立：(C − P(t)) ⋅ (C − P(t)) = r^2
可以通过将 P(t)替换为其展开形式来找到：(C − (Q + td)) ⋅ (C − (Q + td)) = r^2
我们在左边有三个向量，右边有三个向量进行点乘。如果我们计算完整的点乘，我们会得到九个向量。
你当然可以一步步写出来，但我们不需要那么费劲。如果你记得，我们要解的是t，所以我们会根据是否存在t来分离各项：
(−td + (C − Q)) ⋅ (−td + (C − Q)) = r^2
现在我们遵循向量代数的规则来分配点乘：
t^2d ⋅ d − 2td ⋅ (C−Q) + (C − Q) ⋅ (C − Q) = r^2
将半径的平方移到左边：
t^2d ⋅ d − 2td ⋅ (C−Q) + (C − Q) ⋅ (C − Q) - r^2 = 0
很难看出这个方程具体是什么，但方程中的向量和r都是常数且已知。
此外，我们拥有的向量通过点积都简化为标量。
唯一未知的是t，并且我们有一个t^2，这意味着这个方程是二次方程。
你可以通过使用二次公式来解二次方程 ax^2+bx+c=0：
( −b ± √(b^2 − 4ac)) / 2a
所以解出光线-球体相交方程中的t，就得到了a、b和c的这些值：
a = d ⋅ d
b = −2d ⋅ (C−Q)
c = (C − Q) ⋅ (C − Q) − r^2
继续化简令b = -2h，h = d ⋅ (C−Q)
最后求根公式为：h ± √(b^2 − 4ac) / a
*/

type Sphere struct {
	// 物理属性
	Center Point   // 圆心
	Radius float64 // 半径
	// 渲染属性
	Material         MaterialI        // 材质定义
	RemovableSetting RemovableSetting // 移动设置
	AABB             *AABB
}

func NewSphere(center Point, radius float64) *Sphere {
	return &Sphere{
		Center: center,
		Radius: radius,
		AABB:   NewAABBFromPoints(Point(Vec3(center).Sub(Vec3{radius, radius, radius})), Point(Vec3(center).Add(Vec3{radius, radius, radius}))),
	}
}

func (sphere *Sphere) WithMaterial(mat MaterialI) *Sphere {
	sphere.Material = mat
	return sphere
}

func (sphere *Sphere) SetBoundingBox(aabb *AABB) {
	if sphere.RemovableSetting.IsRemovable {
		origin := sphere.Center
		radius := sphere.Radius
		end := sphere.NowAt(1.0)
		aabbOrigin := NewAABBFromPoints(Point(Vec3(origin).Sub(Vec3{radius, radius, radius})), Point(Vec3(origin).Add(Vec3{radius, radius, radius})))
		aabbEnd := NewAABBFromPoints(Point(Vec3(end).Sub(Vec3{radius, radius, radius})), Point(Vec3(end).Add(Vec3{radius, radius, radius})))
		sphere.AABB = NewAABBFromAABB(aabbOrigin, aabbEnd)
	} else {
		sphere.AABB = aabb
	}
}

func (sphere *Sphere) GetBoundingBox() *AABB {
	return sphere.AABB
}

func (sphere *Sphere) NowAt(time float64) Point {
	// 判断物体是否可以移动
	if sphere.RemovableSetting.IsRemovable {
		return sphere.RemovableSetting.MovingFunc(time)
	} else {
		// 不设置默认为静止
		return sphere.Center
	}
}

// SetUniformLinearMovement 设置球体匀速直线运动
func (sphere *Sphere) SetUniformLinearMovement(end Point) {
	sphere.RemovableSetting.IsRemovable = true
	sphere.RemovableSetting.MovingFunc = func(t float64) Point {
		moveRay := Ray{
			Origin:    sphere.Center,
			Direction: Vec3(end).Sub(Vec3(sphere.Center)).Normalize(),
			TM:        t,
		}
		return moveRay.At(t)
	}
}

func (sphere *Sphere) Hittable(ray Ray, rayT Interval) (hit bool, hitRecord HitRecord) {
	currentCenter := sphere.NowAt(ray.TM)

	var oc = Vec3(currentCenter).Sub(Vec3(ray.Origin))
	a := ray.Direction.LengthSquared()
	h := ray.Direction.Dot(oc)
	c := oc.LengthSquared() - sphere.Radius*sphere.Radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		// 无实根（无交点）
		return false, hitRecord
	}
	// 有交点，只需要找出正负根都存在范围内的情况为True
	// 整个Z轴是反的，0最大
	sqrtDiscriminant := math.Sqrt(discriminant)
	root := (h - sqrtDiscriminant) / a // 负根（离相机更近）
	if !rayT.Surrounds(root) {
		root = (h + sqrtDiscriminant) / a
		if !rayT.Surrounds(root) {
			return false, hitRecord // 负根在范围内，正根在范围外
		}
	}
	hitRecord.Time = root
	hitRecord.HitPoint = ray.At(hitRecord.Time)
	outwardNormal := Vec3(hitRecord.HitPoint).Sub(Vec3(currentCenter)).Div(sphere.Radius) //外部法向
	if ray.Direction.Dot(outwardNormal) < 0 {
		hitRecord.FrontFace = true
		hitRecord.Normal = outwardNormal
	} else {
		hitRecord.FrontFace = false
		hitRecord.Normal = outwardNormal.MultiplicationNum(-1.0)
	}
	// 记录击中位置的材质
	hitRecord.Material = sphere.Material
	return true, hitRecord
}
