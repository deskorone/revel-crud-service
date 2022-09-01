package controllers

import (
	"fmt"
	"strconv"
	"testAuth/app/models"
	"testAuth/app/service"

	"github.com/revel/revel"
)

// App Структура контроллера
type App struct {
	*revel.Controller
}

// Hotels Возвращает вид отелей
func (c App) Hotels() revel.Result {
	return c.Render()
}

// Register Регистрация пользователя в теле запроса отправляется json типа models.User
func (c App) Register() revel.Result {
	u := models.User{}
	err := c.Params.BindJSON(&u)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = service.GetService().AuthService.Registration(c.Controller, u)
	if err != nil {
		c.Response.Status = 400
		return c.RenderJSON(map[string]interface{}{"Error": err})
	}
	return c.RenderJSON(map[string]interface{}{"You": "auth"})
}

// SubToHotel Подписка на отель доступен только залогининым пользователям (в параметрах принимает id отеля)
func (c App) SubToHotel() revel.Result {
	id := c.Params.Query.Get("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Error": err.Error()})
	}
	a, err := service.GetService().UserService.SubHotel(c.Controller, n)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Error": err.Error()})
	}
	return c.RenderJSON(a)
}

// Logout Выход из аккаунта
func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.RenderJSON(map[string]interface{}{"Good": "logout"})
}

// Login Вход в аккаунт принимает в теле models.LoginRequest
func (c App) Login() revel.Result {
	request := models.LoginRequest{}
	err := c.Params.BindJSON(&request)
	if err != nil {
		return BuildCredError(c.Controller, err.Error())
	}
	err = service.GetService().AuthService.Login(c.Controller, request)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Error": "No valide data"})
	}
	return c.RenderJSON(map[string]interface{}{"Auth": "Ok"})
}

// AddHotel Добавить отель (доступно только залогининым)
func (c App) AddHotel() revel.Result {
	req := models.Hotel{}
	err := c.Params.BindJSON(&req)
	if err != nil {
		return BuildCredError(c.Controller, "No valid data")
	}

	hotel, err := service.GetService().HotelService.SaveHotel(req, c.Controller)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Error": err.Error()})
	}

	return c.RenderJSON(hotel)
}

// GetAllHotels Получить все отели из базы
func (c App) GetAllHotels() revel.Result {
	hotels, err := service.GetService().HotelService.GetAllHotels()
	if err != nil {
		return BuildCredError(c.Controller, err.Error())
	}
	return c.RenderJSON(hotels)
}

// UnsubToHotel Отписка от отеля
func (c App) UnsubToHotel() revel.Result {
	id := c.Params.Query.Get("id")
	numId, err := strconv.Atoi(id)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Error": err.Error()})
	}
	err = service.GetService().UserService.UnsubHotel(c.Controller, numId)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Error": err.Error()})
	}
	return c.RenderJSON("OK")
}

// SaveHotelWithoutUser Сохранить отель (доступен не авторизированным пользователям) принимает в теле запоса json типа models.Hotel
func (c App) SaveHotelWithoutUser() revel.Result {

	req := models.Hotel{}
	err := c.Params.BindJSON(&req)
	if err != nil {
		return BuildCredError(c.Controller, "No valid data")
	}
	r, err := service.GetService().HotelService.SaveHotelWithoutUser(req)
	if err != nil {
		return BuildCredError(c.Controller, err.Error())
	}
	return c.RenderJSON(r)
}

// BuildCredError Функция которая делает ошибку с статусом 400
func BuildCredError(c *revel.Controller, msg string) revel.Result {
	c.Response.Status = 400
	return c.RenderJSON(map[string]interface{}{"Error": msg})
}

// HotelsPagination Функция которая выдает пагинацию страницы
func (c App) HotelsPagination(page, size int) revel.Result {

	arr, _, err := service.GetService().HotelService.GetPaginationHotels(page, size)

	if err != nil {
		return BuildCredError(c.Controller, err.Error())
	}

	return c.RenderJSON(arr)
}

// HotelsWs Функция которая отвечает за веб сокет соединение
func (c App) HotelsWs(webSocket revel.ServerWebSocket) revel.Result {

	// Канал оповещающий что соединение закрыто
	closeCh := make(chan int)
	service.GetService().WebSocketService.AppendConnection(webSocket, closeCh)
	defer close(closeCh)
	for {
		//Ждем пока придет оповещение о закрытии соединения и выходим из цикла
		select {
		case <-closeCh:
			return nil
		}
	}
}
