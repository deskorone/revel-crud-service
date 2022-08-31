package service

import (
	"sync"
	"testAuth/app"
	"testAuth/app/models"
	"testAuth/app/repo"

	"github.com/revel/revel"
)

type Service struct {
	UserService      UserService
	HotelService     HotelService
	AuthService      AuthService
	WebSocketService WebSocketService
}

type UserService interface {
	SaveUser(u models.User) (*models.User, error)
	SubHotel(c *revel.Controller, hID int) (*models.Hotel, error)
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
	GetPaginationHotels(page, size int) ([]models.Hotel, int, error)
	GetAllHotels() ([]models.Hotel, error)
	SaveHotel(h models.Hotel, c *revel.Controller) (*models.HotelResp, error)
	SaveHotelWithoutUser(h models.Hotel) (*models.Hotel, error)
	DeleteHotel(uId int, hId int) error
	GetHotelByUser(c *revel.Controller) ([]models.Hotel, error)
	AddCommentToHotel(c *revel.Controller, hID int, text string) (*models.Comment, error)
	ParseHotelsFromUrl(arr []models.Hotel) ([]models.Hotel, error)
}

type WebSocketService interface {
	GetChanels() []chan models.Hotel
	AddChanel(ch chan models.Hotel)
	DeleteChan(ch chan models.Hotel)
	GetMessage(ws revel.ServerWebSocket) *models.Hotel
	AppendConnection(ws revel.ServerWebSocket)
	GetChan() <-chan models.Hotel
	DeleteConnection(ws revel.ServerWebSocket)
}

var instance Service
var once sync.Once
var r *repo.Repository

var ch = make(chan models.Hotel)

func GetService() *Service {

	once.Do(func() {
		r = repo.NewRepo(app.DB)
		instance = Service{
			UserService:      getUserServiceImpl(r),
			HotelService:     getHotelServiceImpl(r, ch),
			AuthService:      getAuthServiceImpl(r),
			WebSocketService: getWebSockImpl(ch),
		}
	})
	//fmt.Println(instance)
	return &instance
}
