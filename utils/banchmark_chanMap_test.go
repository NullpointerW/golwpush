package utils

import (
	"sync"
	"testing"
)

//func TestV(t *testing.T) {
//	//sl := make([]int, 3, 4)
//	//slAp(sl)
//	//fmt.Println(len(sl))
//	mp := make(map[int]int)
//	fmt.Println(len(mp))
//	mapAp(mp)
//	fmt.Println(len(mp))
//}

func BenchmarkSyncMap(b *testing.B) {
	var smap = sync.Map{}
	for i := 0; i < 1000000; i++ {
		smap.Store(i, i)
	}
	wg := sync.WaitGroup{}
	wg.Add(1000000)
	for i := 0; i < 1000000; i++ {
		key := i
		go func() {
			smap.Delete(key)
			wg.Done()
		}()
	}
	wg.Wait()
	var count int64
	smap.Range(func(k, v any) bool {
		count++
		return true
	})
	//fmt.Printf("syncMap len %d\n", count)
}

func BenchmarkChMap(b *testing.B) {
	var cmap = ChanMap[int, int]{RmCh: make(chan int, 1000000)}

	for i := 0; i < 1000000; i++ {
		cmap.put(i, i)
	}
	wg := sync.WaitGroup{}
	wg.Add(1000000)
	for i := 0; i < 1000000; i++ {
		key := i
		go func() {
			cmap.Rm(key)
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(cmap.RmCh)
	}()
	for k := range cmap.RmCh {
		cmap.rm(k)
	}
	//fmt.Printf("chMapCount :%d\n", len(cmap.map0))
}
