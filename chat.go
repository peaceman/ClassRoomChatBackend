package main

import "log"

type chatHub struct {
	clients              map[*client]bool
	pastMessages         []message
	inboundMessageChan   chan message
	registerClientChan   chan *client
	unregisterClientChan chan *client
}

func newChatHub() *chatHub {
	chatHub := new(chatHub)
	chatHub.clients = make(map[*client]bool)
	chatHub.inboundMessageChan = make(chan message)
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

	chatHub.sendPastMessagesToClient(client)

	log.Println("Registered client with address", client)
}

func (chatHub *chatHub) sendPastMessagesToClient(receiver *client) {
	for _, msg := range chatHub.pastMessages {
		receiver.sendChan <- []byte(msg.content)
	}
}

func (chatHub *chatHub) unregisterClient(client *client) {
	if _, ok := chatHub.clients[client]; !ok {
		return
	}

	delete(chatHub.clients, client)
	close(client.sendChan)
	log.Println("Unregistered client with address", client)
}

func (chatHub *chatHub) broadcastMessage(msg message) {
	log.Printf("Broadcast message from %s to %d clients | %s", msg.sender, len(chatHub.clients), msg.content)
	for client := range chatHub.clients {
		select {
		case client.sendChan <- []byte(msg.content):
		default:
			chatHub.unregisterClient(client)
		}
	}

	chatHub.pastMessages = append(chatHub.pastMessages, msg)
}
