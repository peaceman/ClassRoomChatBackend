package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var bindAddress = flag.String("addr", ":1338", "http service address")

func setupHTTPEndpoints(chatHub *chatHub, phoneDataHub *PhoneDataHub) {
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.Handle("/chat", chatHub)
	http.HandleFunc("/phone-data", phoneDataHub.handleIncomingPhoneData)
	http.HandleFunc("/data", phoneDataHub.handleWebClientConnection)
}

func main() {
	flag.Parse()

	chatHub := newChatHub()
	go chatHub.loop()

	phoneDataHub := NewPhoneDataHub()
	go phoneDataHub.loop()

	setupHTTPEndpoints(chatHub, phoneDataHub)

	log.Println("Open HTTP socket at:", *bindAddress)
	err := http.ListenAndServe(*bindAddress, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type WebClient struct {
	webSocket               *websocket.Conn
	SendChan                chan []byte
	unregisterWebClientChan chan *WebClient
}

func NewWebClient(webSocket *websocket.Conn, unregisterWebClientChan chan *WebClient) *WebClient {
	webClient := new(WebClient)
	webClient.webSocket = webSocket
	webClient.SendChan = make(chan []byte, maxMessageSize*2)
	webClient.unregisterWebClientChan = unregisterWebClientChan

	return webClient
}

func (webClient *WebClient) StartServing() {
	go webClient.writeLoop()
	webClient.readLoop()
}

func (webClient *WebClient) close() {
	webClient.unregisterWebClientChan <- webClient
	webClient.webSocket.Close()
}

func (webClient *WebClient) writeLoop() {
	defer webClient.close()

	pingTicker := time.NewTicker(pingInterval)

	for {
		select {
		case <-pingTicker.C:
			if err := webClient.write(websocket.PingMessage, []byte{}); err != nil {
				log.Println("PingTicker Error", err)
				return
			}

		case message, ok := <-webClient.SendChan:
			if !ok {
				webClient.write(websocket.CloseMessage, []byte{})
				return
			}

			if err := webClient.write(websocket.TextMessage, message); err != nil {
				log.Println("WebClient SendChan Error", err)
				return
			}
		}
	}
}

func (webClient *WebClient) readLoop() {
	defer webClient.close()

	webClient.webSocket.SetReadLimit(maxMessageSize)
	webClient.webSocket.SetReadDeadline(time.Now().Add(pongWait))
	webClient.webSocket.SetPongHandler(func(string) error { webClient.webSocket.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, _, err := webClient.webSocket.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (webClient *WebClient) write(messageType int, payload []byte) error {
	webClient.webSocket.SetWriteDeadline(time.Now().Add(writeWait))
	return webClient.webSocket.WriteMessage(messageType, payload)
}

func (webClient *WebClient) String() string {
	return fmt.Sprintf("%v", webClient.webSocket.RemoteAddr())
}

type PhoneDataHub struct {
	webClients              map[*WebClient]bool
	phoneData               []string
	IncomingPhoneDataChan   chan string
	RegisterWebClientChan   chan *WebClient
	UnregisterWebClientChan chan *WebClient
}

func NewPhoneDataHub() *PhoneDataHub {
	phoneDataHub := new(PhoneDataHub)
	phoneDataHub.webClients = make(map[*WebClient]bool)
	phoneDataHub.IncomingPhoneDataChan = make(chan string)
	phoneDataHub.RegisterWebClientChan = make(chan *WebClient)
	phoneDataHub.UnregisterWebClientChan = make(chan *WebClient)

	return phoneDataHub
}

func (phoneDataHub *PhoneDataHub) handleIncomingPhoneData(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		http.Error(responseWriter, "Method not allowed", 405)
		return
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	phoneDataHub.IncomingPhoneDataChan <- string(body)
}

func (phoneDataHub *PhoneDataHub) handleWebClientConnection(responseWriter http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(responseWriter, "Method not allowed", 405)
		return
	}

	webSocket, err := upgrader.Upgrade(responseWriter, request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	webClient := NewWebClient(webSocket, phoneDataHub.UnregisterWebClientChan)
	phoneDataHub.RegisterWebClientChan <- webClient

	webClient.StartServing()
}

func (phoneDataHub *PhoneDataHub) loop() {
	for {
		select {
		case webClient := <-phoneDataHub.RegisterWebClientChan:
			phoneDataHub.registerWebClient(webClient)
		case webClient := <-phoneDataHub.UnregisterWebClientChan:
			phoneDataHub.unregisterWebClient(webClient)
		case phoneData := <-phoneDataHub.IncomingPhoneDataChan:
			phoneDataHub.broadcastPhoneData(phoneData)
		}
	}
}

func (phoneDataHub *PhoneDataHub) registerWebClient(webClient *WebClient) {
	phoneDataHub.webClients[webClient] = true
	phoneDataHub.sendPhoneDataToWebClient(webClient)
	log.Println("Registered WebClient with address", webClient)
}

func (phoneDataHub *PhoneDataHub) sendPhoneDataToWebClient(webClient *WebClient) {
	for _, phoneData := range phoneDataHub.phoneData {
		webClient.SendChan <- []byte(phoneData)
	}
}

func (phoneDataHub *PhoneDataHub) unregisterWebClient(webClient *WebClient) {
	if _, ok := phoneDataHub.webClients[webClient]; !ok {
		return // given WebClient is not registered
	}

	delete(phoneDataHub.webClients, webClient)
	close(webClient.SendChan)
	log.Println("Unregistered WebClient with address", webClient)
}

func (phoneDataHub *PhoneDataHub) broadcastPhoneData(phoneData string) {
	log.Printf("Broadcast PhoneData to %d WebClients | %s", len(phoneDataHub.webClients), phoneData)

	for webClient := range phoneDataHub.webClients {
		select {
		case webClient.SendChan <- []byte(phoneData):
		default:
			phoneDataHub.unregisterWebClient(webClient)
		}
	}

	phoneDataHub.phoneData = append(phoneDataHub.phoneData, phoneData)
}
