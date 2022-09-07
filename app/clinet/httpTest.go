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
	var do = 2000
	var wg sync.WaitGroup
	wg.Add(do)
	t := time.Now()
	for i := 0; i < do; i++ {
		go func(seq int) {
			_, err := http.Post(url, "text/xml",
				strings.NewReader("push-testing"+strconv.Itoa(seq)))
			if err != nil {
				fmt.Println(err)
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
