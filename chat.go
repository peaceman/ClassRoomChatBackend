package main

import "log"

type chatHub struct {
	clients map[*client]bool
	pastMessages []string
	inboundMessageChan chan string
	registerClientChan chan *client
	unregisterClientChan chan *client
}

func NewChatHub() *chatHub {
	chatHub := new(chatHub)
	chatHub.clients = make(map[*client]bool)
	chatHub.inboundMessageChan = make(chan string)
	chatHub.registerClientChan = make(chan *client)
	chatHub.unregisterClientChan = make(chan *client)

	return chatHub
}

func (chatHub *chatHub) loop() {
	for {
		select {
		case client := <-chatHub.registerClientChan:
			chatHub.registerClient(client)
		case client := <-chatHub.unregisterClientChan:
			chatHub.unregisterClient(client)
		case message := <-chatHub.inboundMessageChan:
			chatHub.broadcastMessage(message)
		}
	}
}

func (chatHub *chatHub) registerClient(client *client) {
	chatHub.clients[client] = true

	for _, message := range chatHub.pastMessages {
		client.send <- []byte(message)
	}

	log.Println("registered client", client)
}

func (chatHub *chatHub) unregisterClient(client *client) {
	if _, ok := chatHub.clients[client]; !ok {
		return
	}

	delete(chatHub.clients, client)
	close(client.send)
	log.Println("unregistered client", client)
}

func (chatHub *chatHub) broadcastMessage(message string) {
	log.Println("broadcastMessage", message)
	for client := range chatHub.clients {
		select {
		case client.send <- []byte(message):
		default:
			chatHub.unregisterClient(client)
		}
	}

	chatHub.pastMessages = append(chatHub.pastMessages, message)
}