package logger

import (
	"fmt"
	"log"
	"os"
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
	Env = Dev
	mu  = sync.Mutex{}

	red = func(s string) string {
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
	log.SetPrefix(magenta(errorPrefix))
	log.Output(2, fmt.Sprintf(format, v...))
}

func Error(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	log.SetPrefix(magenta(errorPrefix))
	log.Output(2, fmt.Sprintln(v...))
}

func Fatalf(format string, v ...any) {
	mu.Lock()
	defer mu.Unlock()
	log.SetPrefix(red(fatalPrefix))
	log.Output(2, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func Fatal(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	log.SetPrefix(red(fatalPrefix))
	log.Output(2, fmt.Sprintln(v...))
	os.Exit(1)
}

func Warnf(format string, v ...any) {
	mu.Lock()
	defer mu.Unlock()
	log.SetPrefix(yellow(warnPrefix))
	log.Output(2, fmt.Sprintf(format, v...))

}

func Warn(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	log.SetPrefix(yellow(warnPrefix))
	log.Output(2, fmt.Sprintln(v...))
}

func Infof(format string, v ...any) {
	mu.Lock()
	defer mu.Unlock()
	log.SetPrefix(green(infoPrefix))
	log.Output(2, fmt.Sprintf(format, v...))
}

func Info(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	log.SetPrefix(green(infoPrefix))
	log.Output(2, fmt.Sprintln(v...))
}

func Debugf(format string, v ...any) {
	mu.Lock()
	defer mu.Unlock()
	if Env {
		log.SetPrefix(blue(debugPrefix))
		log.Output(2, fmt.Sprintf(format, v...))
	}
}

func Debug(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	if Env {
		log.SetPrefix(blue(debugPrefix))
		log.Output(2, fmt.Sprintln(v...))
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
	reFlag := log.Flags() & log.Lshortfile
	defer func() {
		log.SetFlags(log.Flags() | reFlag)
		mu.Unlock()
	}()
	log.SetFlags(log.Flags() &^ log.Lshortfile)
	var prefix string
	var fatal = false
	if lFlag := cFlag & L_Bs; lFlag != 0 {
		switch lFlag {
		case L_Fatal:
			prefix += red(strings.TrimSuffix(fatalPrefix, " "))
			fatal = true
		case L_Err:
			prefix += magenta(strings.TrimSuffix(errorPrefix, " "))
		case L_Warn:
			prefix += yellow(strings.TrimSuffix(warnPrefix, " "))
		case L_Info:
			prefix += green(strings.TrimSuffix(infoPrefix, " "))
		case L_Debug:
			prefix += blue(strings.TrimSuffix(debugPrefix, " "))
		default:
			prefix += green(strings.TrimSuffix(infoPrefix, " "))
		}
	}
	if cFlag&Ack != 0 {
		prefix += strings.TrimSpace(green(ackPrefix))
	}
	if cFlag&(Ping|Pong) != 0 {
		cFlag &^= HeartBeat
	}
	if cFlag&HeartBeat != 0 {
		prefix += strings.TrimSpace(green(hbPrefix))
	}
	if cFlag&Ping != 0 {
		prefix += strings.TrimSpace(green(pingPrefix))
	}
	if cFlag&Pong != 0 {
		prefix += strings.TrimSpace(green(pongPrefix))
	}
	if cFlag&Kick != 0 {
		prefix += strings.TrimSpace(magenta(kickPrefix))
	}
	if cFlag&Login != 0 {
		prefix += strings.TrimSpace(green(loginPrefix))
	}
	if cFlag&Msg != 0 {
		prefix += strings.TrimSpace(green(msgPrefix))
	}
	if cFlag&Cli != 0 {
		cFlag = cFlag &^ Srv
		prefix += strings.TrimSpace(blue(cliPrefix))
	}
	if cFlag&Srv != 0 {
		prefix += strings.TrimSpace(yellow(srvPrefix))
	}
	if cFlag&Addr != 0 {
		if cFlag&Uid != 0 && cFlag&Host != 0 {
			prefix += strings.TrimSpace(yellow(fmt.Sprintf(uidHostPrefix, uid, host)))
		} else if cFlag&Uid != 0 {
			prefix += strings.TrimSpace(yellow(fmt.Sprintf(uidPrefix, uid)))
		} else {
			prefix += strings.TrimSpace(yellow(fmt.Sprintf(addrPrefix, host)))
		}
	}
	log.SetPrefix(prefix)
	if _fmt {
		log.Output(2, fmt.Sprintf(format, v...))
		return
	}
	log.Output(2, fmt.Sprintln(v...))

	if fatal {
		os.Exit(1)
	}

}

func ModifyLv(lv Level) {
	mu.Lock()
	defer mu.Unlock()
	Env = lv
}
