package logger

import (
	"github.com/NullpointerW/gopush/utils"
	"log"
	"os"
	"time"
)

func createFile() *os.File {
	fn := "log" + time.Now().Format(utils.FilePrefixTimeLayout) + ".log"
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	return f
}
