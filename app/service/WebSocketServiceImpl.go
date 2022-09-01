package service

import (
	"github.com/revel/revel"
	"sync"
	"testAuth/app/models"
)

type WebSocketServiceImpl struct {
	ch             chan models.Hotel
	connectionsMap map[revel.ServerWebSocket]chan int
}

func (service *WebSocketServiceImpl) GetMap() map[revel.ServerWebSocket]chan int {
	return service.connectionsMap
}

// DeleteConnection удаление соединения из мапы срединений
func (service *WebSocketServiceImpl) DeleteConnection(ws revel.ServerWebSocket) {
	delete(service.connectionsMap, ws)
}

// AppendConnection Добавление соединения в мапу соединений
func (service *WebSocketServiceImpl) AppendConnection(ws revel.ServerWebSocket, closeChan chan int) {
	service.connectionsMap[ws] = closeChan
}

// GetMessage Ожидание сообщения от HotelSevrice
func (service *WebSocketServiceImpl) GetMessage() {
	for {
		// Ждем пока из HotelService придет сообщение
		select {
		case hotel := <-service.ch:
			for webSocket := range service.connectionsMap {
				if err := webSocket.MessageSendJSON(hotel); err != nil {
					// Шлю для данного соединения сообщения что соеднинение прервано
					service.connectionsMap[webSocket] <- 0
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
		m := make(map[revel.ServerWebSocket]chan int)
		instanceWebSockImpl = &WebSocketServiceImpl{ch: ch, connectionsMap: m}
		// Запуск ожидания сообщения в фоне
		go instanceWebSockImpl.GetMessage()
	})
	return instanceWebSockImpl
}
