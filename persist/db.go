package persist

import (
	"github.com/NullpointerW/golwpush/logger"
	"github.com/go-redis/redis"
	"runtime"
	"sync"
)

var (
	Redis          *redis.Client
	RedisKeyPrefix = "golwpush."
)

type NamedkeyCache struct {
	pool    sync.Map
	putfunc func(basekey string) string
}
type NamedkeyCacheV2 struct {
	mu      sync.RWMutex
	pool    map[string]string
	putfunc func(basekey string) string
}

func (c *NamedkeyCacheV2) Key(k string) string {
	c.mu.RLock()
	if v, e := c.pool[k]; e {
		c.mu.RUnlock()
		return v
	}
	c.mu.RUnlock()
	r := c.putfunc(k)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.pool[k] = r
	return r
}

func (c *NamedkeyCache) Key(k string) string {
	if v, e := c.pool.Load(k); e {
		k, _ = v.(string)
		return k
	}
	r := c.putfunc(k)
	c.pool.Store(k, r)
	return r
}

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		//DB:       32,
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Infof("redis connected,conn pool num:%d", runtime.NumCPU()*10)
}
