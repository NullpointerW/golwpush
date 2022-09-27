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
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 400, MaxConnsPerHost: 400,
		},
	}
	var do = 400
	var wg sync.WaitGroup
	wg.Add(do)
	t := time.Now()
	for i := 0; i < do; i++ {
		go func(seq int) {
			for i := 0; i < 400; i++ {
				resp, err := client.Post(url, "text/xml",
					strings.NewReader("push-testing"+strconv.Itoa(seq)+":"+strconv.Itoa(i)))
				if err != nil {
					fmt.Println(err)
				}
				resp.Body.Close()
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
