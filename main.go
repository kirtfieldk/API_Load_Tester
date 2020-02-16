package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Response struct {
	error error
	data  string
}
type URL struct {
	url string
}

func (url *URL) getReq(rate int, responseStream chan<- Response, done <-chan interface{}) error {
	var returnValue Response
	go func() {
		for x := 0; x < rate; x++ {
			res, err := http.Get(url.url)
			if err != nil {
				returnValue.error = err
			}
			// resBody, err := ioutil.ReadAll(res.Body)
			if err != nil {
				log.Fatal(err)
			}
			res.Body.Close()
			returnValue.data = "fin"
			select {
			case <-done:
				fmt.Println("QUITE the GET")
				return
			case responseStream <- returnValue:
			}

		}
	}()

	return nil
}
func (url *URL) postReq(reqBody []byte, rate int, responseStream chan<- Response, done <-chan interface{}) error {
	var returnValue Response
	go func() {
		for x := 0; x < rate; x++ {
			res, err := http.Post(url.url, "application/json", bytes.NewBuffer(reqBody))
			if err != nil {
				returnValue.error = err
			}
			res.Body.Close()
			returnValue.data = "Finish"
			select {
			case <-done:
				fmt.Println("Quite the POST")
				return
			case responseStream <- returnValue:
			}
		}
	}()
	return nil
}

func main() {
	fmt.Println("Helscscfewfefweflo")
	defer fmt.Println("finished")
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
	responseStream := make(chan Response, limit)
	paramsArr := args[1:]
	go func() {
		done := make(chan interface{})
		defer close(done)
		for num := 0; num < limit; num++ {

			userUrl := URL{url: args[1] + paramsArr[1]}
			method := paramsArr[2]
			reqLimit, err := strconv.Atoi(paramsArr[3])
			if err != nil {
				log.Fatal(err)
			}
			switch method {
			case <-done:
				return
			case "GET":
				userUrl.getReq(reqLimit, responseStream, done)
			case "POST":
				userUrl.postReq(reqBody, reqLimit, responseStream, done)
			}
			paramsArr = paramsArr[3:]
		}
	}()
	for y := range responseStream {
		fmt.Println(y.data)
	}
}
