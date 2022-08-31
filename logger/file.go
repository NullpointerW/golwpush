package logger

import (
	"github.com/NullpointerW/golwpush/utils"
	"log"
	"os"
	"time"
)

func createFile() *os.File {
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
	return f
}
