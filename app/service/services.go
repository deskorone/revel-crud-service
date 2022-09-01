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
	SaveUser(user models.User) (*models.User, error)
	SubHotel(c *revel.Controller, hotelId int) (*models.Hotel, error)
	UnsubHotel(c *revel.Controller, hotelId int) error
	GetUserById(c *revel.Controller, userId int) (*models.User, error)
}

type AuthService interface {
	Registration(c *revel.Controller, r models.User) error
	Login(c *revel.Controller, r models.LoginRequest) error
	GetUser(c *revel.Controller) (*models.UserView, error)
	GetUserById(userId int) (*models.User, error)
}

type HotelService interface {
	GetPaginationHotels(page, size int) ([]models.Hotel, int, error)
	GetAllHotels() ([]models.Hotel, error)
	SaveHotel(hotel models.Hotel, c *revel.Controller) (*models.HotelResp, error)
	SaveHotelWithoutUser(h models.Hotel) (*models.Hotel, error)
	DeleteHotel(userId int, hotelId int) error
	GetHotelByUser(c *revel.Controller) ([]models.Hotel, error)
	AddCommentToHotel(c *revel.Controller, hotelId int, text string) (*models.Comment, error)
	ParseHotelsFromUrl(arr []models.Hotel) ([]models.Hotel, error)
}

type WebSocketService interface {
	GetMessage()
	AppendConnection(webSocket revel.ServerWebSocket, closeChan chan int)
	DeleteConnection(webSocket revel.ServerWebSocket)
	GetMap() map[revel.ServerWebSocket]chan int
}

var instance Service
var once sync.Once
var r *repo.Repository

var ch = make(chan models.Hotel)

// GetService Получение сервиса
func GetService() *Service {

	// Иницилизация сервисов
	once.Do(func() {
		r = repo.NewRepo(app.DB)
		instance = Service{
			WebSocketService: getWebSockImpl(ch),
			UserService:      getUserServiceImpl(r),
			HotelService:     getHotelServiceImpl(r, ch),
			AuthService:      getAuthServiceImpl(r),
		}
	})
	//fmt.Println(instance)
	return &instance
}
