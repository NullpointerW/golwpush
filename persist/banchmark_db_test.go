package persist

import (
	"encoding/json"
	"fmt"
	"github.com/NullpointerW/golwpush/pkg"
	"github.com/NullpointerW/golwpush/utils"
	"math"
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

func TestZPOPMin(t *testing.T) {
	k := KeyCache.Key(strconv.FormatUint(307, 10))
	c := Redis.ZCard(k).Val()
	pageSize := 5000
	pageNum := int(math.Ceil(float64(c) / float64(pageSize)))
	for i := 0; i < pageNum; i++ {
		//r, err := persist.Redis.ZRange(k, int64(i*pageSize), int64(i*pageSize+pageSize)-1).Result()
		r, err := Redis.ZPopMin(k, int64(pageSize)).Result()
		if err != nil {
			fmt.Println(err)
			//todo
			return
		}
		for _, z := range r {
			jsonRaw, ok := z.Member.(string)
			if !ok {
				fmt.Printf("%v", z)
				return
			}
			var p pkg.SendMarshal
			err = json.Unmarshal(utils.Scb(jsonRaw), &p)
			if err != nil {
				fmt.Println(err)
				//todo handle
				return
			}
			fmt.Println(p)
		}
	}
}

func TestMsgRetransmission(t *testing.T) {
	k := KeyCache.Key(strconv.FormatUint(134, 10))
	c := Redis.ZCard(k).Val()
	pageSize := 5000
	pageNum := int(math.Ceil(float64(c) / float64(pageSize)))
	for i := 0; i < pageNum; i++ {
		var del []interface{}
		r, err := Redis.ZRange(k, int64(i*pageSize), int64(i*pageSize+pageSize)-1).Result()
		//r, err := persist.Redis.ZPopMin(k, int64(pageSize)).Result()
		if err != nil {
			//todo
			return
		}
		for _, jsonRaw := range r {
			//jsonRaw, ok := z.Member.(string)
			//if !ok {
			//	return
			//}
			var p pkg.SendMarshal
			err = json.Unmarshal(utils.Scb(jsonRaw), &p)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(p)
			del = append(del, jsonRaw)
		}
		status := Redis.ZRem(k, del...)
		ok, err := status.Result()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(ok)
	}
}
