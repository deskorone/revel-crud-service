package service

import (
	"sync"
	"testAuth/app"
	"testAuth/app/models"
	"testAuth/app/repo"

	"github.com/revel/revel"
)

type Service struct {
	UserService
	HotelService
	AuthService
}

type UserService interface {
	SaveUser(u models.User) (*models.User, error)
	SubHotel(c *revel.Controller,hID int) (*models.Hotel, error)
	UnsubHotel(c *revel.Controller, hId int) error
	GetUserById(c *revel.Controller, id int) (*models.User, error)
}

type AuthService interface {
	Registration(c *revel.Controller, r models.User) error
	Login(c *revel.Controller, r models.LoginRequest) error
	GetUser(c *revel.Controller) (*models.UserView, error)
	GetUserById(Id int) (*models.User, error)
}

type HotelService interface {
	SaveHotel(h models.Hotel, c *revel.Controller) (*models.HotelResp, error)
	DeleteHotel(uId int, hId int) error
	GetHotelByUser(c *revel.Controller) ([]models.Hotel, error)
	AddCommentToHotel(c *revel.Controller, hID int) (*models.Comment, error)
}

var instance Service
var once sync.Once
var r *repo.Repository

func GetService() *Service {

	once.Do(func() {
		r = repo.NewRepo(app.DB)
		instance = Service{
			UserService:  getUserServiceImpl(r),
			HotelService: getHotelServiceImpl(r),
			AuthService:  getAuthService(r),
		}
	})
	return &instance
}
