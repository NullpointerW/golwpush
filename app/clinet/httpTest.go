package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

func main() {
	url := "http://localhost:8000/broadcast"
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 0, MaxConnsPerHost: 100,
		},
	}
	var do = 100
	var wg sync.WaitGroup
	wg.Add(do)
	t := time.Now()
	for i := 0; i < do; i++ {
		go func(seq int) {
			for i := 0; i < 1; i++ {
				//raw, _ := json.Marshal("push-testing" + strconv.Itoa(seq) + ":" + strconv.Itoa(i))
				_, err := client.Post(url, "text/xml",
					strings.NewReader(`"test"`))
				if err != nil {
					fmt.Println(err)
				}
				//resp.Body.Close()
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("spent %d ms", time.Now().Sub(t).Milliseconds())
}
