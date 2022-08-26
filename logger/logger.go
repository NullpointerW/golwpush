package logger

import (
	"fmt"
	"log"
	"net"
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

	ackPrefix   = "[ACK]"
	pongPrefix  = "[PONG]"
	pingPrefix  = "[PING]"
	cliPrefix   = "[CLI]"
	srvPrefix   = "[SRV]"
	kickPrefix  = "[KICK]"
	addrPrefix  = "[%s]"
	loginPrefix = "[LOGIN]"
	msgPrefix   = "[MSG]"

	Ack           = uint16(0x01)
	Ping          = Ack << 1
	Pong          = Ping << 1
	Cli           = Pong << 1
	Srv           = Cli << 1
	Kick          = Srv << 1
	Addr          = Kick << 1 //0x0040
	Login         = Addr << 1
	Msg           = Login << 1
	L_Fatal       = Msg << 1 //0x0100
	L_Err         = L_Fatal << 1
	L_Warn        = L_Err << 1
	L_Info        = L_Warn << 1
	L_Debug       = L_Info << 1
	PingOutput    = Srv | Pong
	MsgOutput     = Srv | Msg
	PingOutputErr = L_Err | PingOutput
	MsgOutputErr  = L_Err | MsgOutput
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

func PrintlnWithAddr(cFlag uint16, addr net.Addr, v ...any) {
	customPrint(cFlag|Addr, false, addr.String(), "%v", v...)
}
func Println(cFlag uint16, v ...any) {
	customPrint(cFlag&^Addr, false, "", "%v", v...)
}
func PrintfWithAddr(cFlag uint16, addr net.Addr, format string, v ...any) {
	customPrint(cFlag|Addr, true, addr.String(), format, v...)
}
func Printf(cFlag uint16, format string, v ...any) {
	customPrint(cFlag&^Addr, true, "", format, v...)
}
func customPrint(cFlag uint16, _fmt bool, addr, format string, v ...any) {
	mu.Lock()
	reFlag := log.Flags() & log.Lshortfile
	defer func() {
		log.SetFlags(log.Flags() | reFlag)
		mu.Unlock()
	}()
	log.SetFlags(log.Flags() &^ log.Lshortfile)
	var prefix string
	if lFlag := cFlag & L_Bs; lFlag != 0 {
		switch lFlag {
		case L_Fatal:
			prefix += red(strings.TrimSuffix(fatalPrefix, " "))
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
		prefix += strings.TrimSpace(yellow(fmt.Sprintf(addrPrefix, addr)))
	}
	log.SetPrefix(prefix)
	if _fmt {
		log.Output(2, fmt.Sprintf(format, v...))
		return
	}
	log.Output(2, fmt.Sprintln(v...))

}

func ModifyLv(lv Level) {
	mu.Lock()
	defer mu.Unlock()
	Env = lv
}
