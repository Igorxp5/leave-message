package routes

import (
	"log"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/igorxp5/leave-message/queue"
)

type Message struct {
	Text string
}

const queueTimerDuration = 60 * time.Second

var message *Message = &Message{"Go is cool \\o/"}
var clientQueue queue.UniqueQueue = queue.New(queue.Max)
var clients map[string]*websocket.Conn = make(map[string]*websocket.Conn)
var queueTimer = time.NewTicker(queueTimerDuration)

//Channels
var newClientChannel = make(chan string)
var disconnectClientChannel = make(chan string)
var postMessageChannel = make(chan bool)

func StartQueueManager() {
	var index uint64
	var err error

	for {
		select {
		case clientId := <-newClientChannel:
			index, err = clientQueue.Add(clientId)
			if err != nil {
				log.Println("[clientQueue.Add]", err)
				err = sendError(clients[clientId], clientId, "Could not add you to the queue :(")
				if err != nil {
					log.Println(err)
				}
				continue
			}

			err = sendClientId(clients[clientId], clientId)
			if err != nil {
				log.Println("[sendClientId]", err)
				clientQueue.Remove(clientId)
				continue
			}

			err = sendQueuePosition(clients[clientId], clientId, index+1)
			if err != nil {
				log.Println("[sendQueuePosition]", err)
				if index != 0 {
					clientQueue.Remove(clientId)
				}
			}
		case clientId := <-disconnectClientChannel:
			if index, err = clientQueue.Index(clientId); err != nil && index != 0 {
				clientQueue.Remove(clientId)
			}
		case <-postMessageChannel:
			clientQueue.Pop()
			queueTimer.Reset(queueTimerDuration)
			broadcastQueuePosition()
		case <-queueTimer.C:
			log.Println("[Timer] Cleaning queue...")
			if !clientQueue.Empty() {
				currentClientId, err := clientQueue.First()
				if err != nil {
					log.Println("[Timer]", err)
				} else {
					log.Printf("[Timer] \"%s\" session expired", currentClientId)
					clientQueue.Pop()
					queueTimer.Reset(queueTimerDuration)
					broadcastQueuePosition()
				}
			}
		}
	}
}
