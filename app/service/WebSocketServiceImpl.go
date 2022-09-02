package service

import (
	"context"
	"github.com/revel/revel"
	"sync"
	"testAuth/app/models"
)

type WebSocketServiceImpl struct {
	ch             chan models.Hotel
	connectionsMap map[revel.ServerWebSocket]context.CancelFunc
}

func (service *WebSocketServiceImpl) GetMap() map[revel.ServerWebSocket]context.CancelFunc {
	return service.connectionsMap
}

// DeleteConnection удаление соединения из мапы срединений
func (service *WebSocketServiceImpl) DeleteConnection(ws revel.ServerWebSocket) {
	delete(service.connectionsMap, ws)
}

// AppendConnection Добавление соединения в мапу соединений
func (service *WebSocketServiceImpl) AppendConnection(ws revel.ServerWebSocket) {
	ctx := context.Background()
	c, cancel := context.WithCancel(ctx)
	service.connectionsMap[ws] = cancel
	select {
	case <-c.Done():
		// Поле того как контекст завершится, функция завершается и соединенин для этого сокета закрывается
	}

}

// GetMessage Ожидание сообщения от HotelSevrice
func (service *WebSocketServiceImpl) GetMessage() {
	for {
		// Ждем пока из HotelService придет сообщение
		select {
		case hotel := <-service.ch:
			for webSocket := range service.connectionsMap {
				if err := webSocket.MessageSendJSON(hotel); err != nil {
					// Прерываю контекст
					service.connectionsMap[webSocket]() // cancelFunc
					delete(service.connectionsMap, webSocket)
				}
			}
		}
	}
}

var instanceWebSockImpl WebSocketService

var o sync.Once

func getWebSockImpl(ch chan models.Hotel) WebSocketService {
	o.Do(func() {
		m := make(map[revel.ServerWebSocket]context.CancelFunc)
		instanceWebSockImpl = &WebSocketServiceImpl{ch: ch, connectionsMap: m}
		// Запуск ожидания сообщения в фоне
		go instanceWebSockImpl.GetMessage()
	})
	return instanceWebSockImpl
}
