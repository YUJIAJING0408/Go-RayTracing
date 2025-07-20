package core

import (
	"sort"
)

type BVHNode struct {
	Left  HittableItemI
	Right HittableItemI
	AABB  *AABB
}

func (B BVHNode) Hittable(ray Ray, rayT Interval) (hit bool, hitRecord HitRecord) {
	if !B.AABB.Hit(&ray) {
		// 没打到，直接返回
		return false, HitRecord{}
	}
	// 打到AABB
	hitLeft, hitRecordLeft := B.Left.Hittable(ray, rayT)
	maxTime := rayT.Max
	if hitLeft {
		maxTime = hitRecordLeft.Time
	}
	hitRight, hitRecordRight := B.Right.Hittable(ray, Interval{
		Min: rayT.Min,
		Max: maxTime,
	})

	// 都未击中
	if !hitRight && !hitLeft {
		return false, HitRecord{}
	}
	// 以下情况为保证击中
	// 如果都击中，比价时间，取较短者
	if hitRight && hitLeft {
		if hitRecordLeft.Time > hitRecordRight.Time {
			hitRecord = hitRecordRight
		} else {
			hitRecord = hitRecordLeft
		}
	}
	// 左侧击中，右侧未击中
	if hitLeft && !hitRight {
		hitRecord = hitRecordLeft
	}
	// 左侧未击中，右侧击中
	if hitRight && !hitLeft {
		hitRecord = hitRecordRight
	}
	//if hitRecord.Material == nil {
	//	println("BVHNode.Hittable")
	//}
	return true, hitRecord
}

func (B BVHNode) SetBoundingBox(aabb *AABB) {
	B.AABB = aabb
}

func (B BVHNode) GetBoundingBox() *AABB {
	return B.AABB
}

func NewBVHNode(hittableList []HittableItemI) (bvhNode *BVHNode) {
	bvhNode = new(BVHNode)
	aabb := &AABB{
		X: Interval{},
		Y: Interval{},
		Z: Interval{},
	}
	// 最长的轴分割
	for _, itemI := range hittableList {
		aabb = NewAABBFromAABB(aabb, itemI.GetBoundingBox())
	}

	axisIndex := aabb.LongestAxis()
	// 随机轴分割
	// axisIndex := utils.RandomInt(0, 2)
	if len(hittableList) == 1 {
		bvhNode.Left = hittableList[0]
		bvhNode.Right = hittableList[0]
	} else if len(hittableList) == 2 {
		bvhNode.Left = hittableList[0]
		bvhNode.Right = hittableList[1]
	} else {
		// 按照随机方向（x，y，z）排序
		SortHittableItems(hittableList, axisIndex)
		midIndex := len(hittableList) / 2
		bvhNode.Left = NewBVHNode(hittableList[:midIndex])
		bvhNode.Right = NewBVHNode(hittableList[midIndex:])
	}
	bvhNode.AABB = aabb
	return bvhNode
}

func SortHittableItems(hittableList []HittableItemI, sortAxisIndex int) {
	sort.Slice(hittableList, func(i, j int) bool {
		return BoxCompare(hittableList[i], hittableList[j], sortAxisIndex)
	})
}

func BoxCompare(a, b HittableItemI, axisIndex int) bool {
	aInterval := a.GetBoundingBox().AxisInterval(axisIndex)
	bInterval := b.GetBoundingBox().AxisInterval(axisIndex)
	return aInterval.Min < bInterval.Min
}

func BoxCompareX(a, b HittableItemI) bool {
	return BoxCompare(a, b, 0)
}

func BoxCompareY(a, b HittableItemI) bool {
	return BoxCompare(a, b, 1)
}

func BoxCompareZ(a, b HittableItemI) bool {
	return BoxCompare(a, b, 2)
}
