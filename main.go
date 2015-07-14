package main

import (
	"flag"
	"time"
	"fmt"
	"net/http"
	"log"
)

var bindAddress = flag.String("addr", ":1338", "http service address")



func setupHttpEndpoints(chatHub *chatHub) {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/chat", chatHub)
}

func main() {
	flag.Parse()

	chatHub := NewChatHub()
	go chatHub.loop()

	setupHttpEndpoints(chatHub)

	log.Println("Open HTTP socket at:", *bindAddress)
	err := http.ListenAndServe(*bindAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	ticker := time.NewTicker(time.Second)
	go func() { 
		defer fmt.Println("deferred message from the go routine")

		for {
			result := <-ticker.C
			fmt.Printf("ticker chan: %T(%v)\n", result, result)
		}
	}()

	timer := time.NewTimer(5 * time.Second)
	<-timer.C
}