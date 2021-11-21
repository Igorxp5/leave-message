package routes

import (
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/igorxp5/leave-message/queue"
)

type Message struct {
	Text string
}

const queueTimerDuration = 60

var message *Message = &Message{"Go is cool \\o/"}
var clientQueue queue.UniqueQueue = queue.New(queue.Max)
var clients map[string]*websocket.Conn = make(map[string]*websocket.Conn)
var queueTimer = time.NewTicker(queueTimerDuration)

func queueNext() {
	clientQueue.Pop()
	queueTimer.Reset(queueTimerDuration)
	broadcastQueuePosition()
	go func() {
		<-queueTimer.C
		if clientQueue.Size() > 0 {
			queueNext()
		}
	}()
}
