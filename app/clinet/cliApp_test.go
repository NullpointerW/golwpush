package main

import (
	"runtime"
	"strconv"
	"testing"
)

func Benchmark(b *testing.B) {
	b.N = 1000
	runtime.GOMAXPROCS(10)
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		if pb.Next() {
			main()
		}
	})
}

func Test(t *testing.T) {
	for i := 0; i < 1000; i++ {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			t.Parallel()

			main()
		})
	}

}
