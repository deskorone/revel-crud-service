package service

import (
	"sync"
	"testAuth/app/models"

	"github.com/revel/revel"
)

type Service struct {
	UserService
	HotelService
	AuthService
}

type UserService interface {
	SaveUser(u models.User) (*models.User, error)
	SubHotel(uId int, hId int) (*models.Hotel, error)
	UnsubHotel(uId int, hId int) error
}

type AuthService interface {
	Login(c *revel.Controller, r models.LoginRequest) error
	GetUser(c *revel.Controller) (*models.User, error)
	CheckUser(c *revel.Controller) (int, error)
	FindUserById(Id int) (*models.User, error)
}

var instance *Service

var once sync.Once

func GetService() *Service {
	if instance == nil {
		instance = &Service{
			UserService:  getUserServiceImpl(),
			HotelService: &HotelServiceImpl{},
			AuthService: getAuthService(),
		}
	}
	return instance
}

type HotelService interface {
	SaveHotel(h models.Hotel) (*models.Hotel, error)
	DeleteHotel(uId int, hId int) error
}
