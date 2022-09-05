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
	wg.Add(5000)
	for i := 0; i < 5000; i++ {
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

//func main() {
//	c := make([]int, 20, 40)
//	fmt.Println(cap(c))
//	fmt.Printf("%p\n", c)
//	c12 := c[1:20]
//	fmt.Printf("%p\n", c12)
//	c = append(c, 1)
//	fmt.Println(cap(c12))
//	fmt.Printf("%p", c12)
//}
