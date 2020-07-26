package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Product struct {
	ProductName string
	Price       int32
	Desc        string
}
type Message struct {
	Name string
	Msg  string
	Product
}

func main() {
	apple := Product{ProductName: "Granny Smith", Price: 10, Desc: "Super Tasty Apple"}
	sellApple := Message{Name: "Keith Kirtfield", Msg: "Please Buy This Apple", Product: apple}
	var m Message
	jsons, err := json.Marshal(sellApple)
	if err != nil {
		log.Fatal("Error")
	}
	fmt.Println([]byte(jsons))
	err = json.Unmarshal(jsons, &m)
	if err != nil {
		log.Fatal("Error")
	}
	fmt.Println(m)
}
