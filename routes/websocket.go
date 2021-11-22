package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

type SocketMessage struct {
	event   string
	message interface{}
}

func (message *SocketMessage) Bytes() ([]byte, error) {
	data := make(map[string]interface{})
	data["event"] = message.event
	data["message"] = message.message
	return json.Marshal(data)
}

func WebSocket(c *websocket.Conn) {
	var err error
	clientId := uuid.New()

	clients[clientId.String()] = c
	defer delete(clients, clientId.String())

	newClientChannel <- clientId.String()

	for {
		if _, _, err = c.ReadMessage(); err != nil {
			//Client has disconnected!
			break
		}
	}

	disconnectClientChannel <- clientId.String()
	log.Printf("socket (%s) connection closed!", clientId)
}

func sendClientId(c *websocket.Conn, clientId string) error {
	message := &SocketMessage{"clientId", clientId}
	return sendMessage(c, clientId, message)
}

func sendError(c *websocket.Conn, clientId string, errorString string) error {
	message := &SocketMessage{"error", errorString}
	return sendMessage(c, clientId, message)
}

func sendMessage(c *websocket.Conn, clientId string, message *SocketMessage) error {
	data, err := message.Bytes()
	if err != nil {
		return errors.New(fmt.Sprintf("socket (%s): %#v", clientId, err))
	}
	if err = c.WriteMessage(1, data); err != nil {
		return errors.New(fmt.Sprintf("socket (%s) write: %v", clientId, err))
	}
	return nil
}

func sendQueuePosition(c *websocket.Conn, clientId string, position uint64) error {
	message := &SocketMessage{"queuePosition", position}
	return sendMessage(c, clientId, message)
}

func broadcastQueuePosition() {
	var err error
	var index uint64

	for clientId, c := range clients {
		index, err = clientQueue.Index(clientId)
		if err != nil {
			log.Printf("[broadcastQueuePosition] socket (%s): %v", clientId, err)
			continue
		}
		err = sendQueuePosition(c, clientId, index+1)
		if err != nil {
			log.Println(err)
		}
	}
}
