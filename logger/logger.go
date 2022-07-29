package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	colorRed     = iota + 91
	colorGreen   //	绿
	colorYellow  //	黄
	colorBlue    // 蓝
	colorMagenta //	洋红

	fatalPrefix = "[FATAL] "
	errorPrefix = "[ERROR] "
	warnPrefix  = "[WARN] "
	infoPrefix  = "[INFO] "
	debugPrefix = "[DEBUG] "
)

var (
	mu = sync.Mutex{}

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
	log.SetPrefix(blue(debugPrefix))
	log.Output(2, fmt.Sprintf(format, v...))
}

func Debug(v ...any) {
	mu.Lock()
	defer mu.Unlock()
	log.SetPrefix(blue(debugPrefix))
	log.Output(2, fmt.Sprintln(v...))
}
