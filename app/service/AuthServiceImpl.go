package service

import (
	"testAuth/app/models"
	"testAuth/app/repo"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	repos *repo.Repository
}

// Registration implements AuthService регистрация пользователя
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

// Login Авторизация пользователя
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

// GetUserById получить пользователя по id
func (c *AuthServiceImpl) GetUserById(Id int) (*models.User, error) {
	u, err := c.repos.UserRepo.FindUserById(Id)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// GetUser Получить пользователя из сессии
func (o *AuthServiceImpl) GetUser(c *revel.Controller) (*models.UserView, error) {

	user := models.UserView{}

	_, err := c.Session.GetInto("user", &user, false)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

var instanceAuthService AuthService

func getAuthServiceImpl(repo *repo.Repository) AuthService {
	if instanceAuthService == nil {
		instanceAuthService = &AuthServiceImpl{repos: repo}
	}
	return instanceAuthService
}
