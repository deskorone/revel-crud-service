package service

import (
	"testAuth/app"
	"testAuth/app/models"

	"github.com/revel/revel"
)

type AuthServiceImpl struct {
}

const findByIdQ = "select * from usr u where u.user_id = $1"

func (o *AuthServiceImpl) Login(c *revel.Controller, r models.LoginRequest) error {
	return nil
}

func (c *AuthServiceImpl) FindUserById(Id int) (*models.User, error) {
	u := models.User{}

	if err := app.DB.QueryRow(findByIdQ, Id).Scan(&u.Id, &u.Name, &u.Password, &u.Balance); err != nil {
		return nil, err
	}
	return &u, nil
}

func (o *AuthServiceImpl) CheckUser(c *revel.Controller) (int, error) {
	id, err := c.Session.Get("user")
	if err != nil {
		return 0, err
	}
	var iId = int(id.(float64))
	return iId, nil
}

func (o *AuthServiceImpl) GetUser(c *revel.Controller) (*models.User, error) {
	id, err := o.CheckUser(c)
	if err != nil{
		return nil, err
	}	
	return o.FindUserById(id)
}

var instanceAuthService AuthService

func getAuthService() AuthService {
	if instanceAuthService == nil {
		instanceAuthService = new(AuthServiceImpl)
	}
	return instanceAuthService
}
