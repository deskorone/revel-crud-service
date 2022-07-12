package service

import (
	"fmt"
	"testAuth/app"
	"testAuth/app/models"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
}


const saveUserQ = "INSERT INTO usr (username, password, balance) values($1,$2,$3) returning *"
const subHotelQ = "insert into usr_hotel (user_id, hotel_id) values($1, $2)"
const unsubHotelQ = "delete from usr_hotel t where t.user_id=$1 and t.hotel_id=$2"
const FindHotelById = "select * hotel h where h.hotel_id = $1"

func (c *UserServiceImpl) SaveUser(u models.User) (*models.User, error) {

	pswd, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)


	err := app.DB.QueryRow(saveUserQ, u.Name, string(pswd), 0).Scan(&u.Id, &u.Name, &u.Password, &u.Balance)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &u, nil
}

func (c *UserServiceImpl) SubHotel(uId int, hId int) (*models.Hotel, error) {

	_, err := app.DB.Exec(subHotelQ, uId, hId)
	h := models.Hotel{}

	if err != nil {
		return nil, err
	}

	e := app.DB.QueryRow(FindHotelById, hId).Scan(&h.ID, &h.Name, &h.Avaible)
	if e != nil {
		return nil, e
	}

	return &h, nil
}

func (c *UserServiceImpl) UnsubHotel(uId int, hId int) error {
	_, err := app.DB.Exec(unsubHotelQ, uId, hId)
	if err != nil {
		return err
	}
	return nil
}

var instanceUserService UserService

func getUserServiceImpl() UserService {
	once.Do(func() {
		instanceUserService = new(UserServiceImpl)
	})
	return instanceUserService
}
