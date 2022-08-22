package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"testing"
	"unsafe"
)

func Benchmark(b *testing.B) {
	b.N = 10000
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
	_, _, _ = bufio.NewReader(os.Stdin).ReadLine()
}

func TestStr(t *testing.T) {
	a := "void"
	fmt.Printf("%p\n", (unsafe.Pointer)(&a))
}
