package service

import (
	"github.com/revel/revel"
	"sync"
	"testAuth/app/models"
)

type WebSocketServiceImpl struct {
	ch chan models.Hotel
	m  map[revel.ServerWebSocket]bool
}

func (w *WebSocketServiceImpl) GetMap() map[revel.ServerWebSocket]bool {
	return w.m
}

func (w *WebSocketServiceImpl) GetChan() <-chan models.Hotel {
	return w.ch
}

func (w *WebSocketServiceImpl) DeleteConnection(ws revel.ServerWebSocket) {
	delete(w.m, ws)
}

func (w *WebSocketServiceImpl) AppendConnection(ws revel.ServerWebSocket) {
	w.m[ws] = true
}

func (w *WebSocketServiceImpl) GetMessage() *models.Hotel {
	for {
		select {
		case h := <-w.ch:
			for i := range w.m {
				if err := i.MessageSendJSON(h); err != nil {
					delete(w.m, i)
				}
			}
		}
	}
}

var instanceWebSockImpl WebSocketService

var o sync.Once

func getWebSockImpl(ch chan models.Hotel) WebSocketService {
	o.Do(func() {
		m := make(map[revel.ServerWebSocket]bool)
		instanceWebSockImpl = &WebSocketServiceImpl{ch: ch, m: m}
	})
	return instanceWebSockImpl
}
