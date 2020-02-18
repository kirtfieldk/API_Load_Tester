package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Response struct {
	error error
	data  string
}
type URL struct {
	url string
}

func (url *URL) getReq(rate int, responseStream chan<- Response) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for x := 0; x <= rate; x++ {
			if x == rate {

				return
			}
			res, err := http.Get(url.url)
			if err != nil {
				continue
			}
			// resBody, err := ioutil.ReadAll(res.Body)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			res.Body.Close()
			responseStream <- Response{error: nil, data: "GOTTED"}
		}

	}()
	wg.Wait()
	fmt.Println("ENDED GET REQ")
	return nil
}
func (url *URL) postReq(reqBody []byte, rate int, responseStream chan<- Response) error {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for x := 0; x <= rate; x++ {
			if x == rate {
				return
			}
			res, err := http.Post(url.url, "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				continue
			}
			res.Body.Close()
			responseStream <- Response{error: nil, data: "POSTED"}
		}

	}()
	wg.Wait()
	fmt.Println("ENDED POST REQ")
	return nil
}

func main() {
	defer fmt.Println("Cloesing program")
	responseStream := make(chan Response)
	defer close(responseStream)
	reqBody, err := json.Marshal(map[string]string{
		"name": "dd",
	})
	if err != nil {
		log.Fatal(err)
	} //
	args := os.Args[1:]
	limit, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal(err)
	}
	paramsArr := args[1:]
	go func() {
		done := make(chan interface{})
		defer close(done)
		for num := 0; num <= limit; num++ {
			if num == limit {
				fmt.Println("Exiting the anon goroutine")
				return
			}
			userUrl := URL{url: args[1] + paramsArr[1]}
			method := paramsArr[2]
			reqLimit, err := strconv.Atoi(paramsArr[3])
			if err != nil {
				log.Fatal(err)
			}
			switch method {
			case "GET":
				userUrl.getReq(reqLimit, responseStream)
			case "POST":
				userUrl.postReq(reqBody, reqLimit, responseStream)
			}
			paramsArr = paramsArr[3:]
		}
	}()

	for y := range responseStream {
		fmt.Println(y.data)
	}

}
