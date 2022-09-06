package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

func main() {
	url := "http://localhost:8000/broadcast"
	var wg sync.WaitGroup
	wg.Add(2000)
	for i := 0; i < 2000; i++ {
		go func(seq int) {
			resp, err := http.Post(url, "text/xml",
				strings.NewReader("push-testing"+strconv.Itoa(seq)))
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(resp)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
}
