package utils

// ChanMap 需要一个goroutine作为监听者来监听管道信息,
//其他goroutine可通过此channel来通知监听者安全删除元素

//ChanMap.Cap指定了map的容量,此值固定,如果添加的元素达到cap则拒绝添加并返回false
//该结构只适用于极其有限的特定场景,只有监听者才能安全的操作该map,
//其他goroutine只能通过channel来通知监听者删除元素

//监听者调用 NewChMap(cap int)来获取该结构和内部map

type ChanMap[K comparable, V any] struct {
	map0 map[K]V
	cap  int
	Del  chan K
}

//type  Store[K comparable, V any] struct {
//	Map0 map[K]V
//	RmCh chan  K
//	Cap int
//}
//type ChanMap[K comparable, V any] interface{
//	put(key K, val V) (ok bool)
//	Rm(key K)
//}

func (cMap ChanMap[any, any]) Cap() int {
	return cMap.cap
}

func (cMap ChanMap[K, V]) Put(key K, val V) (ok bool) {
	ok = true
	if len(cMap.map0) < cMap.cap {
		cMap.map0[key] = val
		return
	}
	return false
}

//func (cMap ChanMap[K, V]) Get(key K) (val V) {
//	val = cMap.map0[key]
//	return
//}
//
//func (cMap ChanMap[K, V]) MonitorRm(key K) {
//	delete(cMap.map0, key)
//}
//
//func (cMap ChanMap[K, V]) Rm(key K) {
//	cMap.RmCh <- key
//}
//
//func (cMap ChanMap[K, V]) Len() int {
//	return len(cMap.map0)
//}

func NewChMap[K comparable, V any](cap int) (ChanMap[K, V], map[K]V) {
	innerMap := make(map[K]V, cap)
	return ChanMap[K, V]{
		cap:  cap,
		Del:  make(chan K, 2*cap), //避免阻塞
		map0: innerMap,
	}, innerMap
}
