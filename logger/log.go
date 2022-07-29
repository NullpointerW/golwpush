package logger

import (
	"fmt"
	"log"
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
	log.SetPrefix(magenta(errorPrefix))
	log.Printf(format, v...)
}

func Error(v ...any) {
	log.SetPrefix(magenta(errorPrefix))
	log.Println(v...)
}

func Fatalf(format string, v ...any) {
	log.SetPrefix(red(fatalPrefix))
	log.Printf(format, v...)
}

func Fatal(v ...any) {
	log.SetPrefix(red(fatalPrefix))
	log.Println(v...)
}

func Warnf(format string, v ...any) {
	log.SetPrefix(yellow(warnPrefix))
	log.Printf(format, v...)
}

func Warn(v ...any) {
	log.SetPrefix(yellow(warnPrefix))
	log.Println(v...)
}

func Infof(format string, v ...any) {
	log.SetPrefix(green(infoPrefix))
	log.Printf(format, v...)
}

func Info(v ...any) {
	log.SetPrefix(green(infoPrefix))
	log.Println(v...)
}

func Debugf(format string, v ...any) {
	log.SetPrefix(blue(debugPrefix))
	log.Printf(format, v...)
}

func Debug(v ...any) {
	log.SetPrefix(blue(debugPrefix))
	log.Println(v...)
}
