package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	url := "http://localhost:8000/broadcast"
	var do = 20
	var wg sync.WaitGroup
	wg.Add(do)
	t := time.Now()
	for i := 0; i < do; i++ {
		go func(seq int) {
			for i := 0; i < 4; i++ {
				_, err := http.Post(url, "text/xml",
					strings.NewReader("push-testing"+strconv.Itoa(seq)+":"+strconv.Itoa(i)))
				if err != nil {
					fmt.Println(err)
				}
			}
			//} else {
			//	fmt.Println(resp)
			//}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("spent %d ms", time.Now().Sub(t).Milliseconds())
}
