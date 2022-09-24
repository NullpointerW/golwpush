package persist

import (
	"strconv"
	"strings"
	"sync"
	"testing"
)

func BenchmarkNamedkeyCache(b *testing.B) {
	cache := namedkeyCachev1{
		sync.Map{},
		func(uid string) string {
			return RedisKeyPrefix + uid
		},
	}
	for i := 0; i < 20000000; i++ {
		cache.Key(strconv.Itoa(i))
	}
	wg := sync.WaitGroup{}
	wg.Add(20000000)
	b.ResetTimer()
	for i := 0; i < 20000000; i++ {
		go func(i int) {
			str := cache.Key(strconv.Itoa(i))
			strings.Split(str, ".")
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkNamedkeyCacheV2(b *testing.B) {
	cache := namedkeyCachev2{
		sync.RWMutex{},
		make(map[string]string, 20000000),
		func(uid string) string {
			return RedisKeyPrefix + uid
		},
	}
	for i := 0; i < 20000000; i++ {
		cache.Key(strconv.Itoa(i))
	}
	wg := sync.WaitGroup{}
	wg.Add(20000000)
	b.ResetTimer()
	for i := 0; i < 20000000; i++ {
		go func(i int) {
			str := cache.Key(strconv.Itoa(i))
			strings.Split(str, ".")
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func BenchmarkNonNamedkeyCache(b *testing.B) {
	wg := sync.WaitGroup{}
	wg.Add(20000000)
	b.ResetTimer()
	for i := 0; i < 20000000; i++ {
		go func(i int) {
			str := RedisKeyPrefix + strconv.Itoa(i)
			strings.Split(str, ".")
			wg.Done()
		}(i)
	}
	wg.Wait()
}
