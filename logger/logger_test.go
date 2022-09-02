package logger

import (
	"fmt"
	"log"
	"testing"
)

var host = "192.168.199:3041"
var uid uint64 = 114514

func TestDebug(t *testing.T) {

	tests := []string{
		"Debug", "Debugf", "Info", "Infof", "Warn", "Warnf", "Error", "Errorf", "Fatal", "Fatalf",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			Debug(tt)
		})
	}
}

func TestDebugf(t *testing.T) {
	type args struct {
		format string
		v      []any
	}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{"test1%d%s", []any{12, "test"}}},
		{"test2", args{"test2%d%s", []any{12, "test"}}},
		{"test3", args{"test3%d%s", []any{12, "test"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Debugf(tt.args.format, tt.args.v...)
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		v []any
	}

	tests := []struct {
		name string
		args args
	}{{"w", args{
		make([]any, 10, 11),
	}}, {}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Error(tt.args.v...)
		})
	}
}

func TestCustomPrint(t *testing.T) {
	PlnWAddr(Kick|Ping|Cli|Msg|Srv, uid, host, "testing for 3")
	Println(Kick|Ping|Cli|Host|Msg, 0, "", "testing for 3")
	PfWAddr(Kick|Ping|Cli|Msg|Srv, 0, host, "testing for %d", uint64(3))
	Printf(Kick|Ping|Uid|Msg|Srv, uid, host, "testing for %d", 3)
	PfNAddr(Kick|Ping|Addr|Msg|Srv, "testing for %d", 3)
	PfNUid(Kick|Ping|Addr|Msg|Srv, host, "testing for %d", 3)
	log.Print("log")
	//PlnWAddr(L_Info|Kick|Ping|Cli|Host|Msg|Srv, 0, host, "testing")
	//PlnWAddr(PingErrOutput|Cli|Uid, uid, host, "testing")

	fmt.Printf("%b\n", L_Fatal)
}

func TestConstVal(t *testing.T) {

	//fmt.Printf("%b\n", Addr)
	//fmt.Printf("%x\n", Addr)
	fmt.Printf("%016b\n", HeartBeat|Srv)
	fmt.Println(HeartBeat|Srv&HeartBeat != 0)
	fmt.Printf("%x\n", L_Fatal)
	fmt.Printf("%x\n", L_Fatal|L_Info|L_Err|L_Debug|L_Warn)

}

func TestCustomPrintf(t *testing.T) {

	ModifyLv(Prod)
	//PlnWAddr(L_Debug|Srv, _addr, "testing for 3")
	PlnWAddr(HeartBeat|Srv, uid, host, "testing for 3")
	//PlnWAddr(HeartBeat|Pong|Srv, _addr, "testing for 3")

}

func TestErrorf(t *testing.T) {
	type args struct {
		format string
		v      []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Errorf(tt.args.format, tt.args.v...)
		})
	}
}

func TestFatal(t *testing.T) {
	type args struct {
		v []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Fatal(tt.args.v...)
		})
	}
}

func TestFatalf(t *testing.T) {
	type args struct {
		format string
		v      []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Fatalf(tt.args.format, tt.args.v...)
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		v []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args.v...)
		})
	}
}

func TestInfof(t *testing.T) {
	type args struct {
		format string
		v      []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Infof(tt.args.format, tt.args.v...)
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		v []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warn(tt.args.v...)
		})
	}
}

func TestWarnf(t *testing.T) {
	type args struct {
		format string
		v      []any
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Warnf(tt.args.format, tt.args.v...)
		})
	}
}
