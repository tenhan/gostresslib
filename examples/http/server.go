package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func Ping(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Ping request.")
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(5)+5))
	io.WriteString(w, "Pong")
}

func main() {
	rand.Seed(time.Now().Unix())
	http.HandleFunc("/ping", Ping)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
