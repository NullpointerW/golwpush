package golwpush

import (
	"encoding/json"
	"github.com/NullpointerW/golwpush/pkg"
	"github.com/NullpointerW/golwpush/protocol"
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
	if lg.alloc > 2048 {
		lg.data = make([]string, 2048)
		lg.alloc = 2048
	}

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
	msgLen, lgLen, buf := len(msg), lgBuf.len, lgBuf.data
	if lgBuf.size+msgLen > (protocol.MaxLen - pkg.MsgModeExtraLen) {
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
	lgLen, lgSize := lgBuf.len, lgBuf.size
	if lgLen <= 0 || lgSize <= 0 {
		return
	}
	b, _ := json.Marshal(lgBuf.data[:lgLen])
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
