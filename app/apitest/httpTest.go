package main

import (
	"encoding/json"
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
			MaxIdleConns: 0, MaxConnsPerHost: 1000,
		},
	}
	var do = 200
	var wg sync.WaitGroup
	wg.Add(do)
	t := time.Now()
	for i := 0; i < do; i++ {
		go func(seq int) {
			for i := 0; i < 5; i++ {
				raw, _ := json.Marshal("push?testing-vcs-for loooooooooooooooooooooong msg" + strconv.Itoa(seq) + ":" + strconv.Itoa(i))
				resp, err := client.Post(url, "text/xml",
					strings.NewReader(string(raw)))
				if err != nil {
					fmt.Println(err)
				}
				resp.Body.Close()
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Printf("spent %d ms", time.Now().Sub(t).Milliseconds())
}

//func main() {
//	url := "http://localhost:8000/broadcast"
//	var do = 10
//	var wg sync.WaitGroup
//	wg.Add(do)
//	t := time.Now()
//	for i := 0; i < do; i++ {
//		go func(seq int) {
//			for i := 0; i < 1; i++ {
//				//raw, _ := json.Marshal("push-testing" + strconv.Itoa(seq) + ":" + strconv.Itoa(i))
//				_, err := http.Post(url, "text/xml",
//					strings.NewReader(`"test"`))
//				if err != nil {
//					fmt.Println(err)
//				}
//				//resp.Body.Close()
//			}
//			wg.Done()
//		}(i)
//	}
//	wg.Wait()
//	fmt.Printf("spent %d ms", time.Now().Sub(t).Milliseconds())
//}
//func main() {
//	url := "http://localhost:8000/broadcast"
//	var do = 10
//	t := time.Now()
//	for i := 0; i < do; i++ {
//		func(seq int) {
//			for i := 0; i < 1; i++ {
//				//raw, _ := json.Marshal("push-testing" + strconv.Itoa(seq) + ":" + strconv.Itoa(i))
//				_, err := http.Post(url, "text/xml",
//					strings.NewReader(`"test"`))
//				if err != nil {
//					fmt.Println(err)
//				}
//				//resp.Body.Close()
//			}
//		}(i)
//	}
//	fmt.Printf("spent %d ms", time.Now().Sub(t).Milliseconds())
//}
