package service

import (
	"sync"
	"testAuth/app/models"
)

type WebSocketServiceImpl struct {
	arr []chan models.Hotel
	ch  chan models.Hotel
}

func (w *WebSocketServiceImpl) GetChan() <-chan models.Hotel {
	return w.ch
}

func (w *WebSocketServiceImpl) GetMessage() *models.Hotel {
	for {
		select {
		case h := <-w.ch:
			return &h
		}
	}
}

func (w *WebSocketServiceImpl) DeleteChan(ch chan models.Hotel) {
	panic("NOt impl")
}

func (w *WebSocketServiceImpl) GetChanels() []chan models.Hotel {
	return w.arr
}

func (w *WebSocketServiceImpl) AddChanel(ch chan models.Hotel) {
	w.arr = append(w.arr, ch)
}

var instanceWebSockImpl WebSocketService

var o sync.Once

func getWebSockImpl(ch chan models.Hotel) WebSocketService {
	o.Do(func() {
		arr := make([]chan models.Hotel, 0)
		instanceWebSockImpl = &WebSocketServiceImpl{arr: arr, ch: ch}
	})
	return instanceWebSockImpl
}
