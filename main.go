package main

//
import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Response struct {
	Error     error  `json:"error"`
	ResStatus string `json:"resStatus"`
	Data      string `json:"data"`
}
type Url struct {
	Url      string `json:"url"`
	Data     [100]Method
	NumCalls int `json:"numCalls"`
}
type Method struct {
	HttpReq   string `json:"httpReq"`
	Extension string `json:"extension"`
	NumCalls  int    `json:"numCalls"`
	JsonData  []byte `json:"jsonData"`
}

func (url *Url) getReq(resp chan Response, meth Method) error {
	var e1 error
	res, err := http.Get(url.Url + meth.Extension)
	if err != nil {
		e1 = err
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		e1 = err
	}
	resp <- Response{Error: e1, Data: string(body[:]), ResStatus: res.Status}
	res.Body.Close()
	return nil
}

func (url *Url) postReq(responseStream chan<- Response, meth Method) error {
	var e1 error
	res, err := http.Post(url.Url+meth.Extension, "application/json", bytes.NewBuffer(meth.JsonData))
	res.Header.Set("Content-Type", "application/json")
	if err != nil {
		e1 = err
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		e1 = err
	}
	res.Body.Close()
	responseStream <- Response{Error: e1, Data: string(body[:]), ResStatus: res.Status}
	return nil
}
func checkUrlFormat(url string) string {
	if strings.Contains(url, "http://") {
		return url
	}
	return "http://" + url
}
func (url *Url) putReq(responseStream chan<- Response, meth Method) error {
	res, err := http.NewRequest(http.MethodPut, url.Url+meth.Extension, bytes.NewBuffer(meth.JsonData))
	if err != nil {
		log.Fatal(err)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	res.Body.Close()
	responseStream <- Response{Error: nil, Data: string(body[:])}
	return nil
}
func main() {
	var totalCalls int
	var jsonData []byte
	var reqArray [100]Method
	done := make(chan interface{})
	defer close(done)
	responseStream := make(chan Response)
	defer close(responseStream)

	/*
		Taking and parsing out the Terminal Data
	*/
	terminalCmd := os.Args[1:]
	iterations, err := strconv.Atoi(terminalCmd[0])
	if err != nil {
		log.Fatal(err)
	}
	urlBase := checkUrlFormat(terminalCmd[1])
	terminalCmd = terminalCmd[2:]
	fmt.Println(terminalCmd)
	for iter := 0; iter < iterations; iter++ {
		terminalMethod := strings.ToUpper(terminalCmd[1])
		extension := terminalCmd[0]
		rate, err := strconv.Atoi(terminalCmd[2])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s %s %d\n", extension, terminalMethod, rate)
		if terminalMethod == "POST" || terminalMethod == "PUT" {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter Json: ")
			text, _ := reader.ReadString('\n')
			terminalMethod = "POST"
			jsonData = []byte(text)
		}
		m1 := Method{Extension: extension, HttpReq: terminalMethod, NumCalls: rate, JsonData: jsonData}
		reqArray[iter] = m1
		terminalCmd = terminalCmd[3:]
	}
	testingSet := Url{Url: urlBase, Data: reqArray}
	for x := 0; x < len(testingSet.Data); x++ {
		totalCalls += testingSet.Data[x].NumCalls
		for j := 0; j < testingSet.Data[x].NumCalls; j++ {
			switch testingSet.Data[x].HttpReq {
			case "GET":
				go testingSet.getReq(responseStream, testingSet.Data[x])
			case "POST":
				go testingSet.postReq(responseStream, testingSet.Data[x])
				// case "PUT":
				// 	go testingSet.putReq(responseStream, testingSet.Data[x])
			}
		}
	}
	testingSet.NumCalls = totalCalls
	for call := 0; call < totalCalls; call++ {
		recieved := <-responseStream
		fmt.Printf("%d:Recieved Data: %s ----:%s\n", call, recieved.Data, recieved.ResStatus)
	}
}
