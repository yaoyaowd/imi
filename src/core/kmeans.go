package core

import (
	data "../data"
	"math/rand"
	"log"
)

func NewKmeans(dataset *data.Dataset, a, b, limit, iterLimit int) []*data.Point {
	initial := rand.Perm(dataset.Rows)[:limit]
	points := []*data.Point{}
	for _, v := range initial {
		p := data.Point{}
		p.Features = []float32{}
		for i := a; i < b; i++ {
			p.Features = append(p.Features, dataset.Points[v].Features[i])
		}
		points = append(points, &p)
	}

	indexes := make([]int, dataset.Rows)
	for iter := 0; iter < iterLimit; iter++ {
		numChanged := 0
		sumDot := float32(0.0)
		for i, p := range dataset.Points {
			id := 0
			dot := p.L2SquareSubarray(points[0], a, b)
			for j := 1; j < limit; j++ {
				dot2 := p.L2SquareSubarray(points[j], a, b)
				if dot2 < dot {
					dot = dot2
					id = j
				}
			}
			if id != indexes[i] {
				indexes[i] = id
				numChanged += 1
			}
			sumDot += dot
		}
		log.Printf("Iteration %d: numChanged %d, total distance: %f", iter, numChanged, sumDot)

		for i, _ := range points {
			for j, _ := range points[i].Features {
				points[i].Features[j] = 0
			}
			points[i].Weight = 0
		}
		for i, idx := range indexes {
			points[idx].Weight += 1.0
			for j, _ := range points[idx].Features {
				points[idx].Features[j] += dataset.Points[i].Features[j + a]
			}
		}
		for i, _ := range points {
			for j, _ := range points[i].Features {
				points[i].Features[j] /= points[i].Weight
			}
		}

		if numChanged == 0 {
			break
		}
	}

	return points
}
