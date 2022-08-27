package repo

import (
	"database/sql"
	"fmt"
	"testAuth/app/models"
)

type HotelRepoImpl struct {
	DB *sql.DB
}

// SaveHotelWithoutUser implements HotelRepo
func (c HotelRepoImpl) SaveHotelWithoutUser(h models.Hotel) (*models.Hotel, error) {

	q := "insert into hotel (name, avaible, rating, price) values ($1, $2, $3, $4) returning hotel_id, name, avaible,  rating, price"

	if err := c.DB.QueryRow(q, h.Name, h.Avaible, h.Rating, h.Price).Scan(&h.ID, &h.Name, &h.Avaible, &h.Rating, &h.Price); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return &h, nil
}

// SaveHotelByUsername implements HotelRepo
func (c HotelRepoImpl) SaveHotelByUsername(h *models.Hotel, username string) (*models.Hotel, error) {

	u := models.User{}
	if err := c.DB.QueryRow(findUserByName, username).Scan(&u.Id, &u.Name, &u.Password, &u.Password); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if err := c.DB.QueryRow(saveHotelQ, h.Name, h.Avaible, u.Id).Scan(&h.ID, &h.Name, &h.Avaible); err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return h, nil
}

// AddComment implements HotelRepo
func (c HotelRepoImpl) AddComment(hId int, uv models.UserView, text string) (*models.Comment, error) {
	cr := models.Comment{}

	err := c.DB.QueryRow(addCommentQ, text, hId, uv.Id).Scan(&cr.ID, &cr.Text, &cr.HotelID, &cr.UserID)
	return &cr, err
}

// GetHotelById implements HotelRepo
func (c HotelRepoImpl) GetHotelById(hId int) (*models.HotelResp, error) {
	h := models.HotelResp{}

	err := c.DB.QueryRow(getHotelByID, hId).Scan(&h.ID, &h.Name, &h.UserView.Id, &h.UserView.Name)

	return &h, err
}

// DeleteHotel implements HotelRepo
func (c HotelRepoImpl) DeleteHotel(uId int, hid int) error {
	_, err := c.DB.Exec(deleteHoteQ, hid, uId)
	if err != nil {
		return err
	}
	return nil
}

// SaveHotel implements HotelRepo
func (c HotelRepoImpl) SaveHotel(h *models.Hotel, id int) (*models.HotelResp, error) {
	hr := models.HotelResp{}
	err := c.DB.QueryRow(saveHotelQ, h.ID, h.Name, h.Avaible, id).Scan(&hr.ID, &hr.Name, &hr.Avaible)
	if err != nil {
		return nil, err
	}
	u := models.UserView{}
	err = c.DB.QueryRow(getUserViewbyIdQ, id).Scan(&u.Id, &u.Name)
	if err != nil {
		return nil, err
	}
	hr.UserView = u
	return &hr, nil
}

// Sub implements HotelRepo
func (c HotelRepoImpl) Sub(hId int, uId int) (*models.Hotel, error) {

	h := models.Hotel{}
	err := c.DB.QueryRow(selectHotelQ, hId).Scan(&h.ID, &h.Name, &h.Avaible)
	if err != nil {
		return nil, err
	}

	if h.Avaible > 0 {
		_, err := c.DB.Exec(subHotelQ, uId, hId)
		if err != nil {
			return nil, err
		}

		_, err = c.DB.Exec(insertAvaible, h.Avaible-1, hId)

		if err != nil {
			return nil, err
		}
	}
	h.Avaible = h.Avaible - 1
	return &h, nil
}

// Unsub implements HotelRepo
func (c HotelRepoImpl) Unsub(hId int, uId int) error {
	h := models.Hotel{}
	err := c.DB.QueryRow(selectHotelQ, hId).Scan(&h.ID, &h.Name, &h.Avaible)
	if err != nil {
		return err
	}
	_, err = c.DB.Exec(unsubHotelQ, uId, hId)
	if err != nil {
		return err
	}
	_, err = c.DB.Exec(insertAvaible, h.Avaible+1, hId)
	return err
}

// GetAllHotels implements HotelRepo
func (c HotelRepoImpl) GetAllHotels() ([]models.HotelResp, error) {
	r, err := c.DB.Query(getAllQ)
	if err != nil {
		return nil, err
	}
	arr := make([]models.HotelResp, 0)

	for r.Next() {
		h := models.HotelResp{}
		if err := r.Scan(&h.ID, &h.Name, &h.Avaible, &h.UserView.Id, &h.UserView.Name); err != nil {
			return nil, err
		}
		arr = append(arr, h)
	}
	return arr, nil
}

// GetHotelByUser implements HotelRepo
func (c HotelRepoImpl) GetHotelsByUser(id int) ([]models.Hotel, error) {
	r, err := c.DB.Query(getHotelByUserQ, id)
	if err != nil {
		return nil, err
	}

	arr := make([]models.Hotel, 0)
	for r.Next() {
		h := models.Hotel{}
		if err := r.Scan(&h.ID, &h.Name, &h.Avaible); err != nil {
			return nil, err
		}
		arr = append(arr, h)
	}
	return arr, nil
}

var hotelRepoInstance HotelRepo

func NewHotelRepo(DB *sql.DB) HotelRepo {
	if hotelRepoInstance == nil {
		hotelRepoInstance = HotelRepoImpl{DB: DB}
	}
	return hotelRepoInstance
}

const (
	saveHotelWithoutUser = "insert into hotel (name, avaible) values ($1, $2) returning hotel_id, name, avaible"
	getAllQ              = "select h.hotel_id, h.name, h.avaible, u.user_id, u.username from hotel h inner join usr as u on u.user_id = h.user_id"
	getHotelByUserQ      = "select h.hotel_id, h.name, h.avaible from hotel h where h.user_id = $1"
	subHotelQ            = "insert into usr_hotel (user_id, hotel_id) values($1,$2) returning *"
	unsubHotelQ          = "delete from usr_hotel t where t.user_id = $1 and t.hotel_id = $2"
	selectHotelQ         = "select h.hotel_id, h.name, h.avaible from hotel h where h.hotel_id = $1"
	saveHotelQ           = "insert into hotel (name, avaible, user_id) values ($1,$2,$3) returning hotel_id, name, avaible"
	getUserViewbyIdQ     = "select user_id, username from usr where user_id = $1"
	deleteHoteQ          = "delete from hotel h where h.hotel_id = $1 and h.user_id = $2"
	addCommentQ          = "insert into comment (text, hotel_id, user_id) values ($1,$2,$3) returning *"
	insertAvaible        = "update hotel set avaible = $1 where hotel_id = $2"
	getHotelByID         = `select

						h.hotel_id,
	  					h.name,
	   					h.avaible,
	    				u.user_id,
		 				u.username 

						from hotel 
		 				h inner join usr 
						as u on u.user_id = h.user_id 
						
						where h.hotel_id = $1`
	findUserByName = "select * from usr where username = $1"
)
