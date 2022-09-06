package utils

import (
	"bytes"
	"encoding/json"
	"testing"
)

type Package struct {
	Uid  string `json:"uid,omitempty"`
	Id   string `json:"id"`
	Mode int    `json:"mode"`
	Data string `json:"data,omitempty"` // ACK
}

var s = "adsfasdfadsfadsfasdfadfadfasdfasdfadsfasdfasdfasdfsadfas"

//func BenchmarkB2sNew(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		_ = Scb(s)
//	}
//}
//func BenchmarkB2sNormal(b *testing.B) {
//	var _ []byte
//	for i := 0; i < b.N; i++ {
//		_ = []byte(s)
//	}
//}

var (
	bt       = []byte("adsfasdfadsfadsfasdfadfadfasdfasdfadsfasdfasdfasdfsadfas")
	cha byte = 'a'
)

const concha byte = 'a'

//func BenchmarkS2BNew(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		_ = Bcs(bt)
//	}
//}
//func BenchmarkS2BNormal(b *testing.B) {
//	var _ []byte
//	for i := 0; i < b.N; i++ {
//		_ = string(bt)
//	}
//}

//func BenchmarkBcsCharConst(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		_ = string(concha)
//	}
//}

//func BenchmarkBcsCharNewConst(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		s := BcsChar(concha)
//		fmt.Printf(s)
//	}
//}

//func BenchmarkBcsCharNew(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		_ = BcsChar(cha)
//		//fmt.Println(s + "vc")
//	}
//}
//
//func BenchmarkBcsChar(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		_ = string(cha)
//		//fmt.Println(s)
//	}
//}

//func BenchmarkJson(b *testing.B) {
//	s := "{\"id\":\"151:1662429247881144800:1\",\"mode\":6,\"data\":\"2022-09-06 09:54:08\"}"
//	for i := 0; i < b.N; i++ {
//		d := json.NewDecoder(strings.NewReader(s))
//		var p Package
//		d.Decode(&p)
//		//fmt.Println(p)
//	}
//}
func BenchmarkJsonEn(b *testing.B) {
	p := Package{
		Id:   "151:1662429247881144800:1",
		Mode: 2,
		Data: "2022-09-06 09:54:08",
	}
	var bs []byte
	buf := bytes.NewBuffer(bs)
	d := json.NewEncoder(buf)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d.Encode(&p)
		buf.WriteByte('|')
		//fmt.Println(string(buf.Bytes()))
	}
}

//func BenchmarkJsonOld(b *testing.B) {
//	s := "{\"id\":\"151:1662429247881144800:1\",\"mode\":6,\"data\":\"2022-09-06 09:54:08\"}"
//	for i := 0; i < b.N; i++ {
//		var p Package
//		json.Unmarshal([]byte(s), &p)
//		//fmt.Println(p)
//	}
//}

func BenchmarkJsonEnOld(b *testing.B) {
	p := Package{
		Id:   "151:1662429247881144800:1",
		Mode: 2,
		Data: "2022-09-06 09:54:08",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//var b []byte
		//d := json.NewEncoder(bytes.NewBuffer(b))

		b, _ := json.Marshal(p)
		b = append(b, '|')
		//fmt.Println(string(b))
		//fmt.Println(p)
	}
}
