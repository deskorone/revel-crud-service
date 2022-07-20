package repo

import (
	"database/sql"
	"testAuth/app/models"
)

//TODO: make const

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
	GetHotelsByUser(id int) ([]models.Hotel, error)
	GetAllHotels() ([]models.HotelResp, error)
	Sub(hId int, uId int) (*models.Hotel, error)
	Unsub(hId int, uId int) error
	SaveHotel(h *models.Hotel, id int) (*models.HotelResp, error)
	DeleteHotel(uId, hid int) error
}

type AuthRepo interface {
	Registration(u models.User) (*models.User, error)
}

func NewRepo(db *sql.DB) *Repository {
	return &Repository{
		UserRepo:  NewUserRepo(db),
		HotelRepo: NewHotelRepo(db),
		AuthRepo:  NewAuthRepo(db),
	}
}
