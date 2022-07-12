package controllers

import (
	"fmt"
	"net/http"
	"testAuth/app/models"
	"testAuth/app/service"

	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	return c.Render()
}

func (c App) Register() revel.Result {

	u := models.User{}

	err := c.Params.BindJSON(&u)

	fmt.Print(u)

	if err != nil {
		fmt.Println(err.Error())
	}

	r, err := service.GetService().UserService.SaveUser(u)

	if err != nil {
		return c.RenderJSON(map[string]interface{}{"Error": "registration"})
	}

	// c.Session["user"] = r.Id
	c.Session["user"] = r.Id
	c.Session.SetNoExpiration()

	return c.RenderJSON(map[string]interface{}{"You": "auth"})
}

func (c App) SubToHotel() revel.Result {
	u, err := service.GetService().AuthService.GetUser(c.Controller)
	if err != nil {
		c.Response.Status = http.StatusForbidden
		return c.RenderJSON(map[string]interface{}{"Not" : "Authorize"})
	}
	fmt.Println(u)
	fmt.Println(u)
	fmt.Println(u)
	fmt.Println(u)
	return c.RenderJSON(map[string]interface{}{"hello": "world"})
}
