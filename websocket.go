package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"time"
	"log"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingInterval = (pongWait * 9) / 10
	maxMessageSize = 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: maxMessageSize * 4,
	WriteBufferSize: maxMessageSize * 4,
}

type client struct {
	webSocket *websocket.Conn
	send chan []byte
	chatHub *chatHub
}

func newClient(webSocket *websocket.Conn, chatHub *chatHub) *client {
	client := new(client)
	client.webSocket = webSocket
	client.send = make(chan []byte, maxMessageSize * 2)
	client.chatHub = chatHub
	return client
}

func (c *client) StartServing() {
	c.chatHub.registerClientChan <- c
	go c.writeLoop()
	c.readLoop()
}

func (c *client) writeLoop() {
	pingTicker := time.NewTicker(pingInterval)
	
	defer func() {
		pingTicker.Stop()
		c.webSocket.Close()
	}()

	for {
		select {
		case <-pingTicker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				log.Println("writeLoop", err)
				return
			}
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func (c *client) write(messageType int, payload []byte) error {
	c.webSocket.SetWriteDeadline(time.Now().Add(writeWait))
	return c.webSocket.WriteMessage(messageType, payload)
}

func (c *client) readLoop() {
	defer func() {
		c.chatHub.unregisterClientChan <- c
		c.webSocket.Close()
	}()

	c.webSocket.SetReadLimit(maxMessageSize)
	c.webSocket.SetReadDeadline(time.Now().Add(pongWait))
	c.webSocket.SetPongHandler(func(string) error { c.webSocket.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		messageType, message, err := c.webSocket.ReadMessage()
		if err != nil {
			break
		}

		log.Println("Received message", messageType, string(message))
		c.chatHub.inboundMessageChan <- string(message)
	}	
}

type WebSocketHandler struct {
	chatHub *chatHub
}

func (wsh *WebSocketHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method != "GET" {
		http.Error(response, "Method not allowed", 405)
		return
	}

	webSocket, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		log.Println(err)
	}

	client := newClient(webSocket, wsh.chatHub)
	client.StartServing()
}

