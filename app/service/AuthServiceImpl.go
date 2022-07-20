package service

import (
	"fmt"
	"testAuth/app/models"
	"testAuth/app/repo"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImpl struct {
	repos *repo.Repository
}

func (o *AuthServiceImpl) Login(c *revel.Controller, r models.LoginRequest) error {
	pswd := r.Password
	username := r.Username
	fmt.Println("HELLO")
	fmt.Println("HELLO")
	fmt.Println("HELLO", o.repos.UserRepo)
	u, err := o.repos.UserRepo.FindUserByName(username)
	if err != nil {
		return err
	}
	e := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pswd))
	if e != nil {
		return e
	}
	c.Session["user"] = u.Id
	return nil
}

func (c *AuthServiceImpl) GetUserById(Id int) (*models.User, error) {
	u, err := c.repos.UserRepo.FindUserById(Id)
	if err != nil {
		return nil, err
	}
	return u, nil
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
	if err != nil {
		return nil, err
	}
	return o.GetUserById(id)
}

var instanceAuthService AuthService

func getAuthService(repo *repo.Repository) AuthService {
	if instanceAuthService == nil {
		instanceAuthService = &AuthServiceImpl{repos: repo}
	}
	return instanceAuthService
}
