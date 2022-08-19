package GoPush

import "GoPush/pkg"

type Contents struct {
	Ids []uint64 `json:"ids"`
	Msg string   `json:"msg"`
	res chan uint
}

func (c Contents) pkg() *pkg.Package {
	return &pkg.Package{
		Data: c.Msg,
		Mode: pkg.MSG,
	}
}

func broadcaster(broadMsg *pkg.Package) {
	for _, conn := range conns {
		conn.write(broadMsg)
	}
}

func multiSend(broadMsg *pkg.Package, ids []uint64, res chan uint) {
	var success uint
	for _, id := range ids {
		if _, exist := conns[id]; exist {
			conns[id].write(broadMsg)
			success++
		}
	}
	res <- success //结果返回
}
