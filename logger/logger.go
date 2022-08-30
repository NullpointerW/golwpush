package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	colorRed     = iota + 91
	colorGreen   //	绿
	colorYellow  //	黄
	colorBlue    // 蓝
	colorMagenta //	洋红

	fatalPrefix       = "[FATAL] "
	errorPrefix       = "[ERROR] "
	warnPrefix        = "[WARN] "
	infoPrefix        = "[INFO] "
	debugPrefix       = "[DEBUG] "
	Dev         Level = true
	Prod              = !Dev

	ackPrefix     = "[ACK]"
	hbPrefix      = "[HEARTBEAT]"
	pongPrefix    = "[HEARTBEAT|PONG]"
	pingPrefix    = "[HEARTBEAT|PING]"
	cliPrefix     = "[CLI]"
	srvPrefix     = "[SRV]"
	kickPrefix    = "[KICK]"
	addrPrefix    = "[%s]"
	uidPrefix     = "[UID:%d]"
	uidHostPrefix = "[UID|HOST:%d|%s]"

	loginPrefix = "[LOGIN]"
	msgPrefix   = "[MSG]"

	HeartBeat     = uint16(0x0001)
	Ack           = HeartBeat << 1
	Ping          = Ack << 1
	Pong          = Ping << 1
	Cli           = Pong << 1
	Srv           = Cli << 1
	Kick          = Srv << 1
	Host          = Kick << 1 //0x0040
	Uid           = Host << 1
	Login         = Uid << 1
	Msg           = Login << 1
	L_Fatal       = Msg << 1 //0x0100
	L_Err         = L_Fatal << 1
	L_Warn        = L_Err << 1
	L_Info        = L_Warn << 1
	L_Debug       = L_Info << 1
	Addr          = Uid | Host
	PingOutput    = Cli | Ping
	PongOutput    = Srv | Pong
	MsgOutput     = Srv | Msg
	PingErrOutput = L_Err | PingOutput
	MsgErrOutput  = L_Err | MsgOutput
	PongErrOutput = L_Err | PongOutput
	SrvErr        = L_Err | Srv
	CliErr        = L_Err | Cli
	L_Bs          = L_Fatal | L_Info | L_Err | L_Debug | L_Warn //0x1f00
)

type Level bool

var (
	std   = log.New(io.MultiWriter(os.Stderr, createFile()), "", log.Ldate|log.Ltime|log.Lshortfile)
	color = runtime.GOOS != "windows"
	Env   = Dev
	mu    = sync.Mutex{}
	red   = func(s string) string {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorRed, s)
	}
	green = func(s string) string {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorGreen, s)
	}
	yellow = func(s string) string {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorYellow, s)
	}
	blue = func(s string) string {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorBlue, s)
	}
	magenta = func(s string) string {
		return fmt.Sprintf("\x1b[%dm%s\x1b[0m", colorMagenta, s)
	}
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func Errorf(format string, v ...any) {
	mu.Lock()
	defer mu.Unlock()
	p := errorPrefix
	if color {
		p = magenta(errorPrefix)
	}
	std.SetPrefix(p)
	std.Output(2, fmt.Sprintf(format, v...))
}

func Error(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	p := errorPrefix
	if color {
		p = magenta(errorPrefix)
	}
	std.SetPrefix(p)
	std.Output(2, fmt.Sprintln(v...))
}

func Fatalf(format string, v ...any) {
	mu.Lock()
	defer mu.Unlock()
	p := fatalPrefix
	if color {
		p = red(fatalPrefix)
	}
	std.SetPrefix(p)
	std.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Fatal(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	p := fatalPrefix
	if color {
		p = red(fatalPrefix)
	}
	std.SetPrefix(p)
	std.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func Warnf(format string, v ...any) {
	mu.Lock()
	defer mu.Unlock()
	p := warnPrefix
	if color {
		p = yellow(warnPrefix)
	}
	std.SetPrefix(p)
	std.Output(2, fmt.Sprintf(format, v...))

}

func Warn(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	p := warnPrefix
	if color {
		p = yellow(warnPrefix)
	}
	std.SetPrefix(p)
	std.Output(2, fmt.Sprintln(v...))
}

func Infof(format string, v ...any) {
	mu.Lock()
	defer mu.Unlock()
	p := infoPrefix
	if color {
		p = green(infoPrefix)
	}
	std.SetPrefix(p)
	std.Output(2, fmt.Sprintf(format, v...))
}

func Info(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	p := infoPrefix
	if color {
		p = green(infoPrefix)
	}
	std.SetPrefix(p)
	std.Output(2, fmt.Sprintln(v...))
}

func Debugf(format string, v ...any) {
	mu.Lock()
	defer mu.Unlock()
	if Env {
		p := debugPrefix
		if color {
			p = blue(debugPrefix)
		}
		std.SetPrefix(p)
		std.Output(2, fmt.Sprintf(format, v...))
	}
}

func Debug(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	if Env {
		p := debugPrefix
		if color {
			p = blue(debugPrefix)
		}
		std.SetPrefix(p)
		std.Output(2, fmt.Sprintln(v...))
	}
}

func PrintlnWithAddr(cFlag uint16, uid uint64, host string, v ...any) {
	customPrint(cFlag|Addr, false, uid, host, "%v", v...)
}

func PrintlnWithHost(cFlag uint16, uid uint64, host string, v ...any) {
	customPrint(cFlag|Host, false, uid, host, "%v", v...)
}

func PrintlnNonHost(cFlag uint16, uid uint64, v ...any) {
	customPrint(cFlag&^Host, false, uid, "", "%v", v...)
}

func PrintlnNonAddr(cFlag uint16, v ...any) {
	customPrint(cFlag&^Addr, false, 0, "", "%v", v...)
}

func PrintfWithHost(cFlag uint16, uid uint64, host string, format string, v ...any) {
	customPrint(cFlag|Host, true, uid, host, format, v...)
}

func PrintfWithAddr(cFlag uint16, uid uint64, host string, format string, v ...any) {
	customPrint(cFlag|Addr, true, uid, host, format, v...)
}

func PrintfNonHost(cFlag uint16, uid uint64, format string, v ...any) {
	customPrint(cFlag&^Host, true, uid, "", format, v...)
}

func PrintfNonAddr(cFlag uint16, format string, v ...any) {
	customPrint(cFlag&^Addr, true, 0, "", format, v...)
}

func PrintlnNonUid(cFlag uint16, host string, v ...any) {
	customPrint(cFlag&^Uid, false, 0, host, "%v", v...)
}

func PrintfNonUid(cFlag uint16, host string, format string, v ...any) {
	customPrint(cFlag&^Uid, true, 0, host, format, v...)
}
func Println(cFlag uint16, uid uint64, host string, v ...any) {
	customPrint(cFlag, false, uid, host, "%v", v...)
}

func Printf(cFlag uint16, uid uint64, host string, format string, v ...any) {
	customPrint(cFlag, true, uid, host, format, v...)
}

func customPrint(cFlag uint16, _fmt bool, uid uint64, host, format string, v ...any) {
	mu.Lock()
	if !Env && cFlag&L_Debug != 0 { //prod
		mu.Unlock()
		return
	}
	reFlag := std.Flags() & log.Lshortfile
	defer func() {
		std.SetFlags(std.Flags() | reFlag)
		mu.Unlock()
	}()
	std.SetFlags(std.Flags() &^ log.Lshortfile)
	var prefix string
	var fatal = false
	if lFlag := cFlag & L_Bs; lFlag != 0 {
		switch lFlag {
		case L_Fatal:
			prefix += colorF(strings.TrimSuffix(fatalPrefix, " "), red)
			fatal = true
		case L_Err:
			prefix += colorF(strings.TrimSuffix(errorPrefix, " "), magenta)
		case L_Warn:
			prefix += colorF(strings.TrimSuffix(warnPrefix, " "), yellow)
		case L_Info:
			prefix += colorF(strings.TrimSuffix(infoPrefix, " "), green)
		case L_Debug:
			prefix += colorF(strings.TrimSuffix(debugPrefix, " "), blue)
		default:
			prefix += colorF(strings.TrimSuffix(infoPrefix, " "), green)
		}
	}
	if cFlag&Ack != 0 {
		prefix += strings.TrimSpace(colorF(ackPrefix, green))
	}
	if cFlag&(Ping|Pong) != 0 {
		cFlag &^= HeartBeat
	}
	if cFlag&HeartBeat != 0 {
		prefix += strings.TrimSpace(colorF(hbPrefix, green))
	}
	if cFlag&Ping != 0 {
		prefix += strings.TrimSpace(colorF(pingPrefix, green))
	}
	if cFlag&Pong != 0 {
		prefix += strings.TrimSpace(colorF(pongPrefix, green))
	}
	if cFlag&Kick != 0 {
		prefix += strings.TrimSpace(colorF(kickPrefix, magenta))
	}
	if cFlag&Login != 0 {
		prefix += strings.TrimSpace(colorF(loginPrefix, green))
	}
	if cFlag&Msg != 0 {
		prefix += strings.TrimSpace(colorF(msgPrefix, green))
	}
	if cFlag&Cli != 0 {
		cFlag = cFlag &^ Srv
		prefix += strings.TrimSpace(colorF(cliPrefix, blue))
	}
	if cFlag&Srv != 0 {
		prefix += strings.TrimSpace(colorF(srvPrefix, yellow))
	}
	if cFlag&Addr != 0 {
		if cFlag&Uid != 0 && cFlag&Host != 0 {
			prefix += strings.TrimSpace(colorF(fmt.Sprintf(uidHostPrefix, uid, host), yellow))
		} else if cFlag&Uid != 0 {
			prefix += strings.TrimSpace(colorF(fmt.Sprintf(uidPrefix, uid), yellow))
		} else {
			prefix += strings.TrimSpace(colorF(fmt.Sprintf(addrPrefix, host), yellow))
		}
	}
	std.SetPrefix(prefix)
	if _fmt {
		std.Output(2, fmt.Sprintf(format, v...))
		return
	}
	std.Output(2, fmt.Sprintln(v...))
	if fatal {
		os.Exit(1)
	}
}

func ModifyLv(lv Level) {
	mu.Lock()
	defer mu.Unlock()
	Env = lv
}

func colorF(s string, f func(s string) string) string {
	if color {
		return f(s)
	}
	return s
}
