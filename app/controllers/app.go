package controllers

import (
	"fmt"
	"strconv"
	"testAuth/app/models"
	"testAuth/app/service"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Hotels() revel.Result {
	return c.Render()
}

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

func (c App) Logout() revel.Result {
	for k := range c.Session {
		delete(c.Session, k)
	}
	return c.RenderJSON(map[string]interface{}{"Good": "logout"})
}

func (c App) Login() revel.Result {
	r := models.LoginRequest{}
	c.Params.BindJSON(&r)
	err := service.GetService().AuthService.Login(c.Controller, r)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Error": "No valide data"})
	}
	return c.RenderJSON(map[string]interface{}{"Auth": "Ok"})
}

func (c App) AddHotel() revel.Result {
	req := models.Hotel{}
	err := c.Params.BindJSON(&req)
	if err != nil {
		return BuildCredError(c.Controller, "No valid data")
	}

	h, err := service.GetService().HotelService.SaveHotel(req, c.Controller)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Error": err.Error()})
	}

	return c.RenderJSON(h)
}

func (c App) GetAllHotels() revel.Result {
	arr, err := service.GetService().HotelService.GetAllHotels()
	if err != nil {
		return BuildCredError(c.Controller, err.Error())
	}
	return c.RenderJSON(arr)
}

func (c App) UnsubToHotel() revel.Result {
	id := c.Params.Query.Get("id")
	n, err := strconv.Atoi(id)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Error": err.Error()})
	}
	err = service.GetService().UserService.UnsubHotel(c.Controller, n)
	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Error": err.Error()})
	}
	return c.RenderJSON("OK")
}

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

func (c App) ParseHtml() revel.Result {

	var arr []models.Hotel

	err := c.Params.BindJSON(&arr)
	if err != nil {
		return BuildCredError(c.Controller, "No valide date")
	}
	_, err = service.GetService().HotelService.ParseHotelsFromUrl(arr)

	if err != nil {
		return BuildCredError(c.Controller, err.Error())
	}

	return c.RenderJSON("ok")
}

func BuildCredError(c *revel.Controller, msg string) revel.Result {
	c.Response.Status = 400
	return c.RenderJSON(map[string]interface{}{"Error": msg})
}

func (c App) HotelsPagination(page, size int) revel.Result {

	arr, _, err := service.GetService().HotelService.GetPaginationHotels(page, size)

	if err != nil {
		return BuildCredError(c.Controller, err.Error())
	}

	return c.RenderJSON(arr)
}

func (c App) HotelsWs(w revel.ServerWebSocket) revel.Result {

	ch := make(chan models.Hotel)
	defer close(ch)

	service.GetService().WebSocketService.AddChanel(ch)
	defer service.GetService().WebSocketService.DeleteChan(ch)

	go service.GetService().WebSocketService.GetMessage(w)

	for {
		select {
		case h := <-ch:
			err := w.MessageSendJSON(h)
			if err != nil {
				return nil
			}
		}
	}

	return nil
}
