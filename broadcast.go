package golwpush

import (
	"encoding/json"
	"github.com/NullpointerW/golwpush/pkg"
)

type lingerBuf struct {
	data  []string
	size  int
	len   int
	alloc int
}

func (lg *lingerBuf) flush() {
	lg.size = 0
	lg.len = 0
	lg.data = make([]string, 2048)
	lg.alloc = 2048
}

var lgBuf = lingerBuf{
	make([]string, 2048),
	0,
	0,
	2048,
}

type Contents struct {
	Ids []uint64 `json:"ids"`
	Msg string   `json:"msg"`
	Res chan uint64
}

func (c Contents) pkg() *pkg.Package {
	return &pkg.Package{
		Data: c.Msg,
		Mode: pkg.MSG,
	}
}

func broadcaster(broadMsg *pkg.Package) {
	//t := time.Now()
	for _, conn := range conns {
		c := conn
		go func(p pkg.Package) {
			c.write(&p)
		}(*broadMsg)

	}
	//logger.Debugf("encode:%d", time.Now().Sub(t).Milliseconds())
}

func mergeMsg(msg string) (send bool) {
	msgLen := len(msg)
	lgLen := lgBuf.len
	buf := lgBuf.data
	if lgBuf.size+msgLen > 1000 {
		lingerSend()
		lgBuf.data[lgLen] = msg
		lgBuf.len++
		lgBuf.size = msgLen
		return true
	}
	if lgBuf.alloc <= lgLen {
		lgBuf.data = append(buf, msg)
		lgBuf.alloc = len(lgBuf.data)
	} else {
		buf[lgLen] = msg
	}
	lgBuf.len++
	lgBuf.size += msgLen
	return
}

func lingerSend() {
	if lgBuf.len <= 0 || lgBuf.size <= 0 {
		return
	}
	b, _ := json.Marshal(lgBuf.data[:lgBuf.len])
	msg := string(b)
	p := &pkg.Package{Mode: pkg.MSG,
		Data: msg}
	broadcaster(p)
	lgBuf.flush()
}

func multiSend(broadMsg *pkg.Package, ids []uint64, res chan uint64) {
	var success uint64
	for _, id := range ids {
		if _, exist := conns[id]; exist {
			conns[id].write(broadMsg)
			success++
		}
	}
	res <- success //结果返回
}
