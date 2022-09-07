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
	var do = 5000
	var wg sync.WaitGroup
	wg.Add(do)
	for i := 0; i < do; i++ {
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
