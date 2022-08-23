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
	Prod        Level = !Dev

	ackPrefix  = "[ACK]"
	pongPrefix = "[PONG]"
	pingPrefix = "[PING]"
	cliPrefix  = "[CLI]"
	srvPrefix  = "[SRV]"
	kickPrefix = "[KICK]"
	addrPrefix = "[%s]"

	ACK  = uint8(0x1)
	PING = ACK << 1
	PONG = PING << 1
	CLI  = PONG << 1
	SRV  = CLI << 1
	KICK = SRV << 1
	ADDR = KICK << 1
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

func PrintlnWithAddr(cFlag uint8, addr net.Addr, v ...any) {
	customPrint(cFlag|ADDR, false, addr.String(), "%v", v...)
}
func Println(cFlag uint8, v ...any) {
	customPrint(cFlag&^ADDR, false, "", "%v", v...)
}
func PrintfWithAddr(cFlag uint8, addr net.Addr, format string, v ...any) {
	customPrint(cFlag|ADDR, true, addr.String(), format, v...)
}
func Printf(cFlag uint8, format string, v ...any) {
	customPrint(cFlag&^ADDR, true, "", format, v...)
}
func customPrint(cFlag uint8, _fmt bool, addr, format string, v ...any) {
	mu.Lock()
	defer mu.Unlock()
	var prefix string
	if cFlag&ACK != 0 {
		prefix += strings.TrimSpace(green(ackPrefix))
	}
	if cFlag&PING != 0 {
		prefix += strings.TrimSpace(green(pingPrefix))
	}
	if cFlag&PONG != 0 {
		prefix += strings.TrimSpace(green(pongPrefix))
	}
	if cFlag&CLI != 0 {
		prefix += strings.TrimSpace(blue(cliPrefix))
	}
	if cFlag&SRV != 0 {
		prefix += strings.TrimSpace(yellow(srvPrefix))
	}
	if cFlag&KICK != 0 {
		prefix += strings.TrimSpace(magenta(kickPrefix))
	}
	if cFlag&ADDR != 0 {
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
