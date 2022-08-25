package utils

type ChanMap[K comparable, V any] struct {
	map0 map[K]V
	Cap  int
	RmCh chan K
}

func (cMap ChanMap[K, V]) put(key K, val V) (ok bool) {
	ok = true
	if len(cMap.map0) < cMap.Cap {
		cMap.map0[key] = val
		return
	}
	return false
}

func (cMap ChanMap[K, V]) rm(key K) {
	delete(cMap.map0, key)
}

func (cMap ChanMap[K, V]) Rm(key K) {
	cMap.RmCh <- key
}

func (cMap ChanMap[K, V]) Len() int {
	return len(cMap.map0)
}
func NewChMap[K comparable, V any](cap int) (ChanMap[K, V], map[K]V) {
	innerMap := make(map[K]V, cap)
	return ChanMap[K, V]{
		Cap:  cap,
		RmCh: make(chan K, cap),
		map0: innerMap,
	}, innerMap
}
