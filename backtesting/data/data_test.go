package data

import (
	"testing"

	"github.com/ivelsantos/cryptor/models"
)

func BenchmarkGetData1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getDataBench(1, b)
	}
}

func BenchmarkGetData10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getDataBench(10, b)
	}
}

func BenchmarkGetData100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		getDataBench(100, b)
	}
}

// func BenchmarkGetData10000(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		getDataBench(10000, b)
// 	}
// }

// func BenchmarkGetData100000(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		getDataBench(100000, b)
// 	}
// }

// func BenchmarkGetData1000000(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		getDataBench(1000000, b)
// 	}
// }

func getDataBench(window_size int, b *testing.B) {
	err := models.InitDB("../../algor.db")
	if err != nil {
		b.Fatal(err)
	}

	algos, err := models.GetAllAlgos()
	if err != nil {
		b.Errorf("Failed to get algos: %v", err)
		return
	}
	kls, err := GetData(algos[0], window_size)
	if err != nil {
		b.Fatal(err)
	}
	b.Logf("klines length: %v", len(kls))
	// b.Logf("%v", kls[499])
}
