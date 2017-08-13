package data

import (
	"math"
)

func (u *Point) Dot(v *Point) float32 {
	s := float32(0.0)
	for i := 0; i < len(u.Features); i++ {
		s += u.Features[i] * v.Features[i]
	}
	return s
}

func (u *Point) L2SquareSubarray(v *Point, a, b int) float32 {
	s := 0.0
	for i := a; i < b; i++ {
		d := float64(u.Features[i] - v.Features[i - a])
		s += d * d
	}
	return float32(s)
}

func (u *Point) L2(v *Point) float32 {
	s := 0.0
	for i := 0; i < len(u.Features); i++ {
		d := float64(u.Features[i] - v.Features[i])
		s += d * d
	}
	return float32(math.Sqrt(s))
}
