package perf

import "testing"

func BenchmarkScoreMap(b *testing.B) {
	w := []float64{}
	f := map[int]float64{}
	for i := 0; i < 300; i++ {
		f[i] = 1.0
		w = append(w, (float64(i) + 1.0) / 300.0)
	}
        for n := 0; n < b.N; n++ {
		ScoreMap(&w, &f)
        }
}


func BenchmarkScoreList(b *testing.B) {
	w := []float64{}
	f := []float64{}
	for i := 0; i < 300; i++ {
		f = append(f, 1.0)
		w = append(w, (float64(i) + 1.0) / 300.0)
	}
        for n := 0; n < b.N; n++ {
                ScoreList(&w, &f)
        }
}

func BenchmarkWriteData(b *testing.B) {
	WriteData2("/Users/dwang/test_dict")
}

func BenchmarkLoadData(b *testing.B) {
	LoadData2("/Users/dwang/test_dict")
}