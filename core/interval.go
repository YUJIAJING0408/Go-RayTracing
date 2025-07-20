package core

import "RayTracingInOneWeekend/utils"

type Interval struct {
	Min, Max float64
}

func NewInterval(min, max float64) Interval {
	return Interval{Min: min, Max: max}
}

// NewNormalInterval [0,+∞)
func NewNormalInterval() Interval {
	return Interval{
		Min: 0,
		Max: utils.Infinity,
	}
}

// NewUniverseInterval (-∞,+∞)
func NewUniverseInterval() Interval {
	return Interval{
		Min: -utils.Infinity,
		Max: utils.Infinity,
	}
}

func NewIntervalFromInterval(i1, i2 *Interval) Interval {
	return Interval{
		Min: min(i1.Min, i2.Min),
		Max: max(i1.Max, i2.Max),
	}
}

// NewEmptyInterval (+∞,-∞)
func NewEmptyInterval() Interval {
	return Interval{
		Min: utils.Infinity,
		Max: -utils.Infinity,
	}
}

func (i *Interval) Copy() Interval {
	return Interval{i.Min, i.Max}
}

func (i *Interval) Size() float64 {
	return i.Max - i.Min
}

// Contains 是否在双开区间内
func (i *Interval) Contains(x float64) bool {
	return x >= i.Min && x <= i.Max
}

// Surrounds 是否在双闭区间内
func (i *Interval) Surrounds(x float64) bool {
	return x > i.Min && x < i.Max
}

// Clamp 限制X在范围内
func (i *Interval) Clamp(x float64) float64 {
	if x < i.Min {
		return i.Min
	}
	if x > i.Max {
		return i.Max
	}
	return x
}

func (i *Interval) Expand(delta float64) *Interval {
	padding := delta / 2
	return &Interval{
		Min: i.Min - padding,
		Max: i.Max + padding,
	}
}
