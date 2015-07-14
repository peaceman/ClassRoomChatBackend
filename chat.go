package main

import "log"

type chatHub struct {
	clients map[*client]bool
	pastMessages []clientMessage
	inboundMessageChan chan clientMessage
	registerClientChan chan *client
	unregisterClientChan chan *client
}

func NewChatHub() *chatHub {
	chatHub := new(chatHub)
	chatHub.clients = make(map[*client]bool)
	chatHub.inboundMessageChan = make(chan clientMessage)
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

	for _, clientMessage := range chatHub.pastMessages {
		client.send <- []byte(clientMessage.message)
	}

	log.Println("Registered client with address", client)
}

func (chatHub *chatHub) unregisterClient(client *client) {
	if _, ok := chatHub.clients[client]; !ok {
		return
	}

	delete(chatHub.clients, client)
	close(client.send)
	log.Println("Unregistered client with address", client)
}

func (chatHub *chatHub) broadcastMessage(clientMessage clientMessage) {
	log.Printf("Broadcast message from %s to %d clients | %s", clientMessage.sender, len(chatHub.clients), clientMessage.message)
	for client := range chatHub.clients {
		select {
		case client.send <- []byte(clientMessage.message):
		default:
			chatHub.unregisterClient(client)
		}
	}

	chatHub.pastMessages = append(chatHub.pastMessages, clientMessage)
}