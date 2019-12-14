package engine

import (
	"sync"
)

// Command represents actions that can be performed in a single event loop iteration.
type Command interface {
	Execute(handler Handler)
}
// Handler allows to send commands to an event loop it's associated with.
type Handler interface {
	Post(cmd Command)
}

type messageQueue struct {
	sync.Mutex
	data []Command
	waiting bool
	//receiveSignal chan struct{}
}

var receiveSignal = make(chan struct{})

func (mq *messageQueue) push(cmd Command) {
	mq.Lock()
	defer mq.Unlock()
	mq.data = append(mq.data, cmd)
	//mq.receiveSignal = make(chan struct{})
	if mq.waiting {
		mq.waiting = false
		receiveSignal <- struct{}{}
	}
}

func (mq *messageQueue) pull() Command {
	mq.Lock()
	defer mq.Unlock()
	if len(mq.data) == 0 {
		mq.waiting = true;
		mq.Unlock()
		<- receiveSignal
		mq.Lock()
	}
	res := mq.data[0]
	mq.data[0] = nil
	mq.data = mq.data[1:]
	return res
}

func (mq *messageQueue) size() int {
	return len(mq.data)
}

type EventLoop struct {
	queue *messageQueue
	terminateReceived bool
	stopSignal chan struct{}
}
