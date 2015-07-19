package main

import (
	"flag"
	"log"
	"net/http"
)

var bindAddress = flag.String("addr", ":1338", "http service address")

func setupHTTPEndpoints(chatHub *chatHub) {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/chat", chatHub)
}

func main() {
	flag.Parse()

	chatHub := newChatHub()
	go chatHub.loop()

	setupHTTPEndpoints(chatHub)

	log.Println("Open HTTP socket at:", *bindAddress)
	err := http.ListenAndServe(*bindAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
