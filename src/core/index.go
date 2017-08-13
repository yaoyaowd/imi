package core

import (
	data "../data"
	"sort"
)

type PostingList []int

type Index struct {
	D *data.Dataset
	D1 []*data.Point
	D2 []*data.Point
	Matrix [][]PostingList
}

func Locate(p *data.Point, c []*data.Point, a, b int) int {
	id := 0
	dis := p.L2SquareSubarray(c[0], a, b)
	for i := 1; i < len(c); i++ {
		dis2 := p.L2SquareSubarray(c[i], a, b)
		if dis2 < dis {
			id = i
			dis = dis2
		}
	}
	return id
}

func NewIndex(dataset *data.Dataset, d1Limit, d2Limit int) *Index {
	d1 := NewKmeans(dataset, 0, dataset.VecLen / 2, d1Limit, 30)
	d2 := NewKmeans(dataset, dataset.VecLen / 2 + 1, dataset.VecLen, d2Limit, 20)
	matrix := make([][]PostingList, d1Limit)
	for i := 0; i < d1Limit; i++ {
		matrix[i] = make([]PostingList, d2Limit)
		for j := 0; j < d2Limit; j++ {
			matrix[i][j] = PostingList{}
		}
	}

	for i, p := range dataset.Points {
		id1 := Locate(p, d1, 0, dataset.VecLen / 2)
		id2 := Locate(p, d2, dataset.VecLen / 2 + 1, dataset.VecLen)
		matrix[id1][id2] = append(matrix[id1][id2], i)
	}
	return &Index{dataset, d1, d2, matrix}
}

type CheckPoint struct {
	X int
	Y int
	Score float32
}

type CheckPoints []CheckPoint

func (c CheckPoints) Len() int {
	return len(c)
}
func (c CheckPoints) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c CheckPoints) Less(i, j int) bool {
	return c[i].Score < c[j].Score
}

func (index *Index) Query(p *data.Point, k int, out chan int) {
	seq := []CheckPoint{}
	for i := 0; i < len(index.D1); i++ {
		for j := 0; j < len(index.D2); j++ {
			s1 := p.L2SquareSubarray(index.D1[i], 0, index.D.VecLen / 2)
			s2 := p.L2SquareSubarray(index.D2[j], index.D.VecLen / 2 + 1, index.D.VecLen)
			seq = append(seq, CheckPoint{i, j, s1 + s2})
		}
	}
	sort.Sort(CheckPoints(seq))

	cpIdx := 0
	x := seq[cpIdx].X
	y := seq[cpIdx].Y
	idx := 0
	tot := 0
	for ; ; {
		if idx < len(index.Matrix[x][y]) {
			out <- index.Matrix[x][y][idx]
			tot += 1
			idx += 1
		} else {
			if tot > k {
				break
			}

			cpIdx += 1
			if cpIdx < len(seq) {
				x = seq[cpIdx].X
				y = seq[cpIdx].Y
			} else {
				break
			}
			idx = 0
		}
	}
	out <- -1
}
