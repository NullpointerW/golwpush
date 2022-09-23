package persist

import (
	"strconv"
	"sync"
	"testing"
)

func BenchmarkNamedkeyCache(b *testing.B) {
	cache := NamedkeyCache{
		sync.Map{},
		func(uid string) string {
			return RedisKeyPrefix + uid
		},
	}
	for i := 0; i < 450000; i++ {
		cache.Key(strconv.Itoa(i))
	}
	wg := sync.WaitGroup{}
	wg.Add(900000)
	b.ResetTimer()
	for i := 0; i < 900000; i++ {
		go func(i int) {
			cache.Key(strconv.Itoa(i))
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkNamedkeyCacheV2(b *testing.B) {
	cache := NamedkeyCacheV2{
		sync.RWMutex{},
		make(map[string]string, 500),
		func(uid string) string {
			return RedisKeyPrefix + uid
		},
	}
	for i := 0; i < 450000; i++ {
		cache.Key(strconv.Itoa(i))
	}
	wg := sync.WaitGroup{}
	wg.Add(900000)
	b.ResetTimer()
	for i := 0; i < 900000; i++ {
		go func(i int) {
			cache.Key(strconv.Itoa(i))
			wg.Done()
		}(i)
	}
	wg.Wait()
}

//func BenchmarkNonNamedkeyCache(b *testing.B) {
//
//	b.ResetTimer()
//	for i := 0; i < 9000000; i++ {
//		_ = RedisKeyPrefix + strconv.Itoa(i)
//	}
//}
