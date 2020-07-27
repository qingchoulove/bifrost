package util

import "sync"

type CloseChannel struct {
	ch     chan int
	closed bool
	once   sync.Once
}

func NewCloseChannel() *CloseChannel {
	return &CloseChannel{
		ch: make(chan int),
	}
}

func (ch *CloseChannel) Wait() chan int {
	return ch.ch
}

func (ch *CloseChannel) Close() {
	ch.once.Do(func() {
		close(ch.ch)
		ch.closed = true
	})
}
