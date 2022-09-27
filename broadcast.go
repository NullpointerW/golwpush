package golwpush

import (
	"encoding/json"
	"github.com/NullpointerW/golwpush/pkg"
	"github.com/NullpointerW/golwpush/protocol"
	"time"
)

type lingerBuf struct {
	data  []string
	size  int
	len   int
	alloc int
}

func lingerProcess() {
	lingerMs := time.NewTicker(time.Millisecond * 300)
	defer lingerMs.Stop()
	for {
		select {
		case <-lingerMs.C:
			lingerSend()
		case msg := <-broadcast0:
			select {
			case <-lingerMs.C:
				lingerSend()
			default:
			}
			if mergeMsg(msg) {
				lingerMs.Reset(time.Millisecond * 300)
			}
		}
	}
}
func (lg *lingerBuf) append(msg string, msgSize int) {
	lgLen := lg.len
	if lg.alloc <= lgLen {
		lg.data = append(lg.data, msg)
		lg.alloc = len(lg.data)
	} else {
		lg.data[lgLen] = msg
	}
	lg.len++
	lg.size += msgSize

}

func (lg *lingerBuf) flush() {
	lg.size = 0
	lg.len = 0
	if lg.alloc > 2048*2 {
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
	msgSize, lgLen := len(msg), lgBuf.len
	meSize := lgBuf.size + msgSize
	if s := protocol.MaxLen - ((pkg.MsgModeExtraLen) + 2 + ((lgLen + 1) * 4) + lgLen); lgLen > 0 && meSize >= s {
		if meSize > s {
			lingerSend()
			lgBuf.data[0] = msg
			lgBuf.len++
			lgBuf.size = msgSize
		} else {
			lgBuf.append(msg, msgSize)
			lingerSend()
		}
		return true
	} else if s = protocol.MaxLen - pkg.MsgModeExtraLen; lgLen == 0 && msgSize == s {
		lgBuf.append(msg, msgSize)
		lingerSend()
		return true
	}
	lgBuf.append(msg, msgSize)
	return
}

func lingerSend() {
	lgLen, lgSize := lgBuf.len, lgBuf.size
	if lgLen <= 0 || lgSize <= 0 {
		return
	}
	var b []byte
	if lgLen == 1 {
		b, _ = json.Marshal(lgBuf.data[0])
	} else {
		b, _ = json.Marshal(lgBuf.data[:lgLen])
	}
	msg := string(b)
	p := &pkg.Package{Mode: pkg.MSG,
		Data: msg}
	mergedMsg <- p
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
