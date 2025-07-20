package utils

import (
	"math/rand"
	"time"
)

func Random() float64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Float64()
}

func RandomInt(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(max-min) + min
}

func RandomBetween(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func Degrees2Radians(degrees float64) float64 {
	return degrees * PI / 180
}
