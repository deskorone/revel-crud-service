package service

import (
	"testAuth/app/models"
	"testAuth/app/repo"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repos *repo.Repository
}

// GetUserById implements UserService
func (o *UserServiceImpl) GetUserById(c *revel.Controller, id int) (*models.User, error) {
	u, err := o.repos.UserRepo.FindUserById(id)

	if err != nil {
		c.Response.Status = 404
		return nil, err
	}
	return u, nil
}

func (c *UserServiceImpl) SaveUser(u models.User) (*models.User, error) {

	pswd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	u.Password = string(pswd)

	user, err := c.repos.UserRepo.Save(u)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *UserServiceImpl) SubHotel(uId int, hId int) (*models.Hotel, error) {
	h, err := c.repos.HotelRepo.Sub(hId, uId)
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (c *UserServiceImpl) UnsubHotel(uId int, hId int) error {
	err := c.repos.HotelRepo.Unsub(hId, uId)
	if err != nil {
		return err
	}
	return nil
}

var instanceUserService UserService

func getUserServiceImpl(r *repo.Repository) UserService {
	if instanceUserService == nil {
		instanceUserService = &UserServiceImpl{r}
	}
	return instanceUserService
}
