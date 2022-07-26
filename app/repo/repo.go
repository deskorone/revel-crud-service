package repo

import (
	"database/sql"
	"testAuth/app/models"
)

type Repository struct {
	UserRepo  UserRepo
	HotelRepo HotelRepo
	AuthRepo  AuthRepo
}

type UserRepo interface {
	FindUserById(id int) (*models.User, error)
	Save(u models.User) (*models.User, error)
	FindUserByName(Name string) (*models.User, error)
}

type HotelRepo interface {
	GetPaginationHotels(page, size int) ([]models.Hotel, int, error)
	GetHotelsByUser(id int) ([]models.Hotel, error)
	GetAllHotels() ([]models.Hotel, error)
	Sub(hId int, uId int) (*models.Hotel, error)
	Unsub(hId int, uId int) error
	SaveHotelWithoutUser(h models.Hotel) (*models.Hotel, error)
	SaveHotel(h *models.Hotel, id int) (*models.HotelResp, error)
	SaveHotelByUsername(h *models.Hotel, username string) (*models.Hotel, error)
	DeleteHotel(uId, hid int) error
	AddComment(hId int, uV models.UserView, text string) (*models.Comment, error)
	GetHotelById(hId int) (*models.HotelResp, error)
}

type AuthRepo interface {
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		UserRepo:  NewUserRepo(db),
		HotelRepo: NewHotelRepo(db),
		AuthRepo:  NewAuthRepo(db),
	}
}
