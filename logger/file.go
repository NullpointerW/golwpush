package logger

import (
	"github.com/NullpointerW/golwpush/utils"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type MultiWriter struct {
	io.Writer
	fd *os.File
}

func createFile() io.Writer {
	path := "log"
	if fd, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("crate dir:%s", path)
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	} else if !fd.IsDir() {
		log.Fatalf("path:%s is not a dir", path)
	}

	fn := path + "\\log" + time.Now().Format(utils.FilePrefixTimeLayout) + ".log"
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	return MultiWriter{io.MultiWriter(os.Stderr, f), f}
}

func cleaner(d time.Duration) {
	path := "log"
clean:
	Info("start clean log file")
	if fd, err := os.Stat(path); err == nil && fd.IsDir() {
		dir, err := os.Open(path)
		if err != nil {
			return
		}
		fs, err := dir.ReadDir(0)
		for _, f := range fs {
			if !f.IsDir() {
				ln := f.Name()
				n := strings.ReplaceAll(ln, "log", "")
				t := strings.Trim(n, ".")
				rt, _ := time.ParseInLocation(utils.FilePrefixTimeLayout, t, utils.TimeLoc)
				if time.Now().Sub(rt) >= d {
					if err = os.Remove(path + "\\" + ln); err != nil {
						Errorf("del log err :%s", err.Error())
					} else {
						Infof("deleted log %s", ln)
					}
				}
			}
		}
	}
	for {
		time.Sleep(d)
		changeFile()
		goto clean
	}
}

func changeFile() {
	mu.Lock()
	defer mu.Unlock()
	cls, ok := std.Writer().(MultiWriter)
	if ok {
		err := cls.fd.Close()
		if err != nil {
			Errorf("close log err:%s", err.Error())
		}
	}
	std.SetOutput(io.MultiWriter(os.Stderr, createFile()))
}
