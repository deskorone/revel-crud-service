package service

import (
	"errors"
	"testAuth/app/models"
	"testAuth/app/repo"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	repos *repo.Repository
}

// Registration implements AuthService
func (*AuthServiceImpl) Registration(c *revel.Controller, r models.User) error {
	u, err := instance.UserService.SaveUser(r)
	if err != nil {
		c.Response.Status = 400
		return err
	}
	uv := models.UserView{}
	uv.Id = u.Id
	uv.Name = u.Name
	c.Session["user"] = uv
	c.Session.SetNoExpiration()
	return nil
}

func (o *AuthServiceImpl) Login(c *revel.Controller, r models.LoginRequest) error {
	pswd := r.Password
	username := r.Username
	u, err := o.repos.UserRepo.FindUserByName(username)
	if err != nil {
		return err
	}
	e := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pswd))
	if e != nil {
		return e
	}
	uv := models.UserView{}
	uv.Id = u.Id
	uv.Name = u.Name
	c.Session["user"] = uv
	return nil
}

func (c *AuthServiceImpl) GetUserById(Id int) (*models.User, error) {
	u, err := c.repos.UserRepo.FindUserById(Id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (o *AuthServiceImpl) GetUser(c *revel.Controller) (*models.UserView, error) {
	uv, err := c.Session.Get("user")
	if err != nil {
		return nil, err
	}

	nv, ok := uv.(map[string]interface{})
	if !ok {
		return nil, err
	}

	u := &models.UserView{}
	u.Id = int(nv["id"].(float64))
	if !ok {
		return nil, errors.New("error cast")
	}

	u.Name, ok = nv["username"].(string)
	if !ok {
		return nil, errors.New("error cast")
	}
	return u, nil
}

var instanceAuthService AuthService

func getAuthService(repo *repo.Repository) AuthService {
	if instanceAuthService == nil {
		instanceAuthService = &AuthServiceImpl{repos: repo}
	}
	return instanceAuthService
}

type castType struct {
	id       float64
	username string
}
