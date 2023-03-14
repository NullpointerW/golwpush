package golwpush

import (
	"encoding/json"
	"github.com/NullpointerW/golwpush/pkg"
	"github.com/NullpointerW/golwpush/protocol"
	"github.com/NullpointerW/golwpush/utils"
	"time"
)

type lingerBuf struct {
	data  []json.RawMessage
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
func (lg *lingerBuf) append(msg json.RawMessage, msgSize int) {
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
		lg.data = make([]json.RawMessage, 2048)
		lg.alloc = 2048
	}

}

var lgBuf = lingerBuf{
	make([]json.RawMessage, 2048),
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
		Data: utils.Scb(c.Msg),
		Mode: pkg.MSG,
	}
}

func broadcaster(broadMsg *pkg.Package) {
	for _, conn := range conns {
		c := conn
		go func(p pkg.Package) {
			c.write(&p)
		}(*broadMsg)
	}
}

func mergeMsg(msg string) (send bool) {
	msgSize, lgLen := len(msg), lgBuf.len
	raw := utils.Scb(msg)
	meSize := lgBuf.size + msgSize
	//todo 需要更改下数据包的计算规则？
	if s := protocol.MaxLen - ((pkg.MsgModeExtraLen) + 2 + ((lgLen + 1) * 4) + lgLen); lgLen > 0 && meSize >= s {
		if meSize > s {
			lingerSend()
			lgBuf.data[0] = raw
			lgBuf.len++
			lgBuf.size = msgSize
		} else {
			lgBuf.append(raw, msgSize)
			lingerSend()
		}
		return true
	} else if s = protocol.MaxLen - pkg.MsgModeExtraLen; lgLen == 0 && msgSize == s {
		lgBuf.append(raw, msgSize)
		lingerSend()
		return true
	}
	lgBuf.append(raw, msgSize)
	return
}

func lingerSend() {
	lgLen, lgSize := lgBuf.len, lgBuf.size
	if lgLen <= 0 || lgSize <= 0 {
		return
	}
	var (
		b []byte
	)
	if lgLen == 1 {
		b = lgBuf.data[0]
	} else {
		b, _ = json.Marshal(lgBuf.data[:lgLen])
	}
	p := &pkg.Package{Mode: pkg.MSG,
		Data: b}
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
