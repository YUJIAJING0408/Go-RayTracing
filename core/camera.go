package core

import (
	"RayTracingInOneWeekend/utils"
	"fmt"
	"math"
	"sync"
)

/*
一个在 3D 空间中的点，所有场景光线都将从此点发出（这通常也被称为视点）。从相机中心到视口中心的向量将垂直于视口。
我们最初将视口与相机中心点之间的距离设置为 1 个单位。这个距离通常被称为焦距。
为了简化，我们将相机中心设置在 (0,0,0)。同时，y 轴向上，x 轴向右，负 z 轴指向观察方向。（这通常被称为右手坐标系。）
现在到了不可避免地棘手的部分。虽然我们的 3D 空间遵循上述约定，但这与我们图像坐标的约定相冲突。
我们希望将零像素放在左上角，并逐行向下工作到最后一个像素在右下角。
这意味着我们的图像坐标 Y 轴是倒置的：Y 值随着图像向下增加。
在扫描图像时，我们将从左上角的像素（像素 0,0）开始，逐行从左到右扫描，然后逐行从上到下扫描。
为了帮助导航像素网格，我们将使用一个从左边缘到右边缘的向量Vu，以及一个从上边缘到下边缘的向量Vv。
我们的像素网格将从视口边缘内缩像素间距的一半。
这样，我们的视口区域将被均匀地划分为宽度 × 高度 个相同的区域。
*/

type Camera struct {
	ImageWidth                 int     // 渲染窗口宽
	ImageHeight                int     // 渲染窗口高
	SamplesPerPixel            int     // 每像素采样数量
	MaxDepth                   int     // 光线最大递归深度
	AspectRatio                float64 // 宽高比
	ViewportHeight             float64 // 视口高度
	ViewportWidth              float64 // 视口宽度
	FocalLength                float64 // 焦距
	PixelSamplesScale          float64 // 每采样权重
	VFov                       float64 // 视野
	DefocusAngle               float64 // 每像素通过的光线变化角度（景深）
	FocusDist                  float64 // 相机观察点与完美焦点平面之间的距离
	CameraCenter               Point   // 相机位置
	LookAt                     Point   // 光线出发点
	LookFrom                   Point   // 焦点
	ViewportU                  Vec3    // 视口水平长度向量
	ViewportV                  Vec3    // 视口垂直长度向量
	PixelDeltaU                Vec3    // 像素水平间隔
	PixelDeltaV                Vec3    // 像素垂直间隔
	ViewportUpperLeft          Vec3    // 视口左上向量
	Pixel00Local               Vec3    // 视口原点
	world                      Scenes  // 场景
	IsAntialiased              bool    // 抗锯齿
	u, v, w, vup               Vec3    // Camera frame basis vectors and Camera-relative "up" direction
	defocusDistU, defocusDistV Vec3
}

func (c *Camera) Add(hittableItem ...HittableItemI) {
	c.world.Add(hittableItem...)
}

func (c *Camera) SetLookAt(p Point) {
	c.LookAt = p
}

func (c *Camera) SetLookFrom(p Point) {
	c.CameraCenter = p
	c.LookFrom = p
}

func NewCamera(lookAt, lookFrom Point, aspectRatio, vFov float64, imageWidth, samplesPerPixel, maxDepth int, isAntialiased bool, defocusAngle, focusDist float64) *Camera {
	center := lookFrom
	vup := Vec3{0, 1, 0}
	// 焦距
	//focalLength := Vec3(lookFrom).Sub(Vec3(lookAt)).Length()
	// 计算u v w
	w := Vec3(lookFrom).Sub(Vec3(lookAt)).Normalize()
	u := vup.Cross(w).Normalize()
	v := w.Cross(u)

	theta := utils.Degrees2Radians(vFov)
	h := math.Tan(theta / 2)
	imageHeight := int(float64(imageWidth) / aspectRatio)
	viewportHeight := 2 * h * focusDist
	viewportWidth := viewportHeight * aspectRatio
	viewportU := u.MultiplicationNum(viewportWidth)
	viewportV := v.MultiplicationNum(-viewportHeight)
	pixelDeltaU := viewportU.Div(float64(imageWidth))
	pixelDeltaV := viewportV.Div(float64(imageHeight))
	//halfU := viewportU.Div(2.0)
	//halfV := viewportV.Div(2.0)
	//c := Vec3(origin).Sub(Vec3{0, 0, focalLength})
	//viewportUpperLeft := c.Sub(halfU).Sub(halfV)
	halfU := viewportU.Div(2.0)
	halfV := viewportV.Div(2.0)
	c := Vec3(lookFrom).Sub(w.MultiplicationNum(focusDist))
	viewportUpperLeft := c.Sub(halfU).Sub(halfV)
	pixel00Local := viewportUpperLeft.Add(pixelDeltaV.MultiplicationNum(0.5).Add(pixelDeltaU.MultiplicationNum(0.5)))
	defocusRadius := focusDist * math.Tan(utils.Degrees2Radians(defocusAngle/2))
	return &Camera{
		MaxDepth:          maxDepth,
		SamplesPerPixel:   samplesPerPixel,
		ImageWidth:        imageWidth,
		ImageHeight:       imageHeight,
		AspectRatio:       aspectRatio,
		ViewportHeight:    viewportHeight,
		ViewportWidth:     viewportWidth,
		CameraCenter:      center,
		LookFrom:          lookFrom,
		ViewportU:         viewportU,
		ViewportV:         viewportV,
		PixelDeltaU:       pixelDeltaU,
		PixelDeltaV:       pixelDeltaV,
		ViewportUpperLeft: viewportUpperLeft,
		Pixel00Local:      pixel00Local,
		LookAt:            Point{0, 0, -1},
		// 抗锯齿开启
		IsAntialiased:     isAntialiased,
		PixelSamplesScale: 1.0 / float64(samplesPerPixel),
		DefocusAngle:      defocusAngle,
		FocusDist:         focusDist,
		defocusDistU:      u.MultiplicationNum(defocusRadius),
		defocusDistV:      v.MultiplicationNum(defocusRadius),
		world: Scenes{
			HittableList: nil,
			HittableAABB: &AABB{
				X: Interval{},
				Y: Interval{},
				Z: Interval{},
			},
		},
	}
}

func (c *Camera) EnabledBVH(enabled bool) {
	if enabled {
		c.world.EnabledBVH = true
		c.world.BVH = NewBVHNode(c.world.HittableList)
	} else {
		c.world.EnabledBVH = false
	}
}

func (c *Camera) Render(name string) {
	var ppm = utils.PPMImage{
		Width:  c.ImageWidth,
		Height: c.ImageHeight,
		Max:    255,
	}
	var pixels = make([]utils.Pixel, 0, c.ImageHeight*c.ImageWidth)
	bar := utils.NewProgressBar(int64(c.ImageWidth*c.ImageHeight), 50)
	for j := 0; j < c.ImageHeight; j++ { // 行
		for i := 0; i < c.ImageWidth; i++ { // 列
			var color = Color{}
			if c.IsAntialiased {
				for _ = range c.SamplesPerPixel {
					var r = c.GetRay(i, j)
					color = Color(Vec3(color).Add(Vec3(c.RayColor(r, c.MaxDepth))))
				}
				color = Color(Vec3(color).MultiplicationNum(c.PixelSamplesScale))
			} else {
				pixelCenter := c.Pixel00Local.Add(c.PixelDeltaU.MultiplicationNum(float64(i))).Add(c.PixelDeltaV.MultiplicationNum(float64(j)))
				rayDirection := pixelCenter.Sub(Vec3(c.CameraCenter))
				ray := NewRay(c.CameraCenter, rayDirection)
				color = c.RayColor(&ray, c.MaxDepth) // 单线程
			}
			bar.Add(1)
			pixels = append(pixels, color.Color2Pixel())
		}
	}
	ppm.Full(pixels)
	ppm.FastWriteAndSave(name)

}

// MultithreadedRender 多线程渲染，不一定速度会更快，开销全花在通信上面
func (c *Camera) MultithreadedRender(name string, maxWorkers, buffer int) {
	wgTask := sync.WaitGroup{}
	wgWorker := sync.WaitGroup{}
	var ppm = utils.PPMImage{
		Width:  c.ImageWidth,
		Height: c.ImageHeight,
		Max:    255,
	}
	var pixels = make([]utils.Pixel, c.ImageHeight*c.ImageWidth) // 这里与上面的Render不同，不是进行Cap分配，是直接进行初始化赋值这么多个
	colorTaskChan := make(chan ColorTask, c.ImageHeight*c.ImageWidth)
	colorResChan := make(chan ColorRes, buffer)

	pb := utils.NewProgressBar(int64(c.ImageWidth*c.ImageHeight), 50)
	// worker
	for w := 0; w < maxWorkers; w++ {
		wgWorker.Add(1)
		go func(resChan chan ColorRes, taskChan chan ColorTask, pb *utils.ProgressBar, w *sync.WaitGroup) {
			defer wgWorker.Done()
			for task := range taskChan {
				// 处理完直接发送，直到taskChan没东西或关闭
				var color = Color{}
				if c.IsAntialiased {
					for _ = range c.SamplesPerPixel {
						var r = c.GetRay(task.WidthIndex, task.HeightIndex)
						color = Color(Vec3(color).Add(Vec3(c.RayColor(r, c.MaxDepth))))
					}
					color = Color(Vec3(color).MultiplicationNum(c.PixelSamplesScale))
				} else {
					color = c.RayColor(task.R, c.MaxDepth)
				}
				pb.Add(1)
				//println(task.WidthIndex, " ", task.HeightIndex)
				resChan <- ColorRes{
					C:           color,
					WidthIndex:  task.WidthIndex,
					HeightIndex: task.HeightIndex,
					Index:       task.WidthIndex + task.HeightIndex*c.ImageWidth,
				}
			}
		}(colorResChan, colorTaskChan, pb, &wgWorker)
	}

	// 协程处理结果
	wgTask.Add(1)
	go func() {
		defer wgTask.Done()
		for colorRes := range colorResChan {
			//println(colorRes.WidthIndex, " ", colorRes.HeightIndex, " ", colorRes.Index)
			pixels[colorRes.HeightIndex*c.ImageWidth+colorRes.WidthIndex] = colorRes.C.Color2Pixel()
		}
	}()

	// 分发task
	for j := 0; j < c.ImageHeight; j++ { // 行
		for i := 0; i < c.ImageWidth; i++ { // 列
			//wgTask.Add(1)
			pixelCenter := c.Pixel00Local.Add(c.PixelDeltaU.MultiplicationNum(float64(i))).Add(c.PixelDeltaV.MultiplicationNum(float64(j)))
			rayDirection := pixelCenter.Sub(Vec3(c.CameraCenter))
			ray := NewRay(c.CameraCenter, rayDirection)
			//println(i, " ", j,)
			colorTaskChan <- ColorTask{
				R:           &ray,
				WidthIndex:  i,
				HeightIndex: j,
			}
		}
	}
	// 关闭任务通道
	close(colorTaskChan)

	// 等待工作协程处理完成
	wgWorker.Wait()

	// 关闭结果通道
	close(colorResChan)

	// 等待结果处理完
	wgTask.Wait()
	fmt.Println("All Pixels have been rendered")
	ppm.Full(pixels)
	ppm.FastWriteAndSave(name)
}

func (c *Camera) RayColor(r *Ray, maxDepth int) Color {
	if maxDepth <= 0 {
		return Color{0, 0, 0}
	}
	if c.world.EnabledBVH {
		hit, hitRecord := c.world.BVH.Hittable(*r, NewInterval(1e-5, utils.Infinity))
		if hit {
			hit, attenuation, scattered := hitRecord.Material.Scatter(r, hitRecord)
			if hit {
				return Color(Vec3(attenuation).MultiplicationVec3(Vec3(c.RayColor(scattered, maxDepth-1))))
			}
			return Color{}
		}
	} else {
		hit, hitRecord := c.world.HitAnything(*r, NewInterval(1e-5, utils.Infinity))
		if hit {
			hit, attenuation, scattered := hitRecord.Material.Scatter(r, hitRecord)
			if hit {
				return Color(Vec3(attenuation).MultiplicationVec3(Vec3(c.RayColor(scattered, maxDepth-1))))
			}
			return Color{}
		}
	}
	var directionNormalized = r.Direction.Normalize()
	a := 0.5 * (directionNormalized.Y + 1.0)
	return Color(Vec3(Color{1.0, 1.0, 1.0}).MultiplicationNum(1 - a).Add(Vec3{0.5, 0.7, 1.0}.MultiplicationNum(a)))

}

func (c *Camera) GetRay(i, j int) *Ray {
	// 从散焦盘构造一条指向像素位置i, j周围随机采样点的相机射线。
	var offset = SampleSquare()
	// 在每个像素邻域内进行随机采样
	pixelSample := c.Pixel00Local.Add(c.PixelDeltaU.MultiplicationNum(float64(i) + offset.X)).Add(c.PixelDeltaV.MultiplicationNum(float64(j) + offset.Y))
	rayOrigin := c.CameraCenter
	if c.DefocusAngle > 0 {
		rayOrigin = c.DefocusDiskSample()
	}
	rayDirection := pixelSample.Sub(Vec3(rayOrigin))
	rayTime := utils.Random()
	ray := NewRayWithTime(rayOrigin, rayDirection, rayTime)
	return &ray
}

// SampleSquare x,y in [-0.5,0.5]
func SampleSquare() Vec3 {
	return Vec3{
		X: utils.Random() - 0.5,
		Y: utils.Random() - 0.5,
	}
}

func (c *Camera) DefocusDiskSample() Point {
	p := RandomInUnitDisk()
	return Point(Vec3(c.CameraCenter).Add(c.defocusDistU.MultiplicationNum(p.X)).Add(c.defocusDistV.MultiplicationNum(p.Y)))
}
