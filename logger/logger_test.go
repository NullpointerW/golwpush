package logger

import (
	"fmt"
	"testing"
)

type addr struct {
}

func (s addr) String() string {
	return "192.168.1.30:2548"
}

func (s addr) Network() string {
	return "192.168.1.30:2548"
}

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
	var _addr addr
	PrintlnWithAddr(Kick|Ping|Cli|Msg, _addr, "testing for 3")
	Println(Kick|Ping|Cli|Addr|Msg, "testing for 3")
	PrintfWithAddr(Kick|Ping|Cli|Msg, _addr, "testing for %d", uint64(3))
	Printf(Kick|Ping|Cli|Addr|Msg, "testing for %d", uint64(3))
	PrintlnWithAddr(Kick|Ping|Cli|Addr|Msg, _addr, "testing")
	fmt.Printf("%b\n", Addr)
	fmt.Printf("%x\n", Addr)
	fmt.Printf("%d\n", Addr)
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
