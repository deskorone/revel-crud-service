package service

import (
	"github.com/revel/revel"
	"sync"
	"testAuth/app/models"
)

type WebSocketServiceImpl struct {
	arr []chan models.Hotel
	ch  chan models.Hotel
	m   map[revel.ServerWebSocket]bool
}

func (w *WebSocketServiceImpl) DeleteConnection(ws revel.ServerWebSocket) {
	delete(w.m, ws)
}

func (w *WebSocketServiceImpl) AppendConnection(ws revel.ServerWebSocket) {
	w.m[ws] = true
}

func (w *WebSocketServiceImpl) GetChan() <-chan models.Hotel {
	return w.ch
}

func (w *WebSocketServiceImpl) GetMessage(ws revel.ServerWebSocket) *models.Hotel {
	for {
		select {
		case h := <-w.ch:
			for i := range w.m {
				i.MessageSendJSON(h)
			}
		}
	}
	return nil
}

func (w *WebSocketServiceImpl) DeleteChan(ch chan models.Hotel) {
	for n, i := range w.arr {
		if i == ch {
			w.arr[n] = w.arr[len(w.arr)-1]
			w.arr = w.arr[:len(w.arr)-1]
			break
		}
	}
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
		m := make(map[revel.ServerWebSocket]bool)
		arr := make([]chan models.Hotel, 0)
		instanceWebSockImpl = &WebSocketServiceImpl{arr: arr, ch: ch, m: m}
	})
	return instanceWebSockImpl
}
