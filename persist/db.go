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
	KeyCache       NamedkeyCache
)

type NamedkeyCache interface {
	Key(string) string
}

type namedkeyCachev1 struct {
	pool    sync.Map
	putfunc func(basekey string) string
}
type namedkeyCachev2 struct {
	mu      sync.RWMutex
	pool    map[string]string
	putfunc func(string) string
}

func (c *namedkeyCachev2) Key(k string) string {
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

func (c *namedkeyCachev1) Key(k string) string {
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
		//DB:       0,
	})
	_, err := Redis.Ping().Result()
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Infof("redis[host:%s]connected,conn_pool_num:%d", "localhost:6379", runtime.NumCPU()*10)
	KeyCache = &namedkeyCachev2{
		sync.RWMutex{},
		make(map[string]string),
		func(uid string) string {
			return RedisKeyPrefix + uid
		},
	}
}
