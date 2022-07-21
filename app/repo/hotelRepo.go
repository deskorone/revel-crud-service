package repo

import (
	"database/sql"
	"testAuth/app/models"
)

type HotelRepoImpl struct {
	DB *sql.DB
}

// AddComment implements HotelRepo
func (c HotelRepoImpl) AddComment(hId int, uv models.UserView, text string) (*models.CommentResp, error) {
	cr := models.CommentResp{}
	err := c.DB.QueryRow(addCommentQ, text, hId, uv.Id).Scan(&cr.ID, &cr.Text)
	if err != nil {
		return nil, err
	}
	cr.UserView = uv
	return &cr, nil
}

// GetHotelById implements HotelRepo
func (HotelRepoImpl) GetHotelById(hId int) (*models.HotelResp, error) {
	panic("unimplemented")
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
	_, err := c.DB.Exec(subHotelQ, uId, hId)
	if err != nil {
		return nil, err
	}
	h := &models.Hotel{}
	err = c.DB.QueryRow(selectHotelQ, hId).Scan(h.ID, h.Name, h.Avaible)
	if err != nil {
		return nil, err
	}
	return h, err
}

// Unsub implements HotelRepo
func (c HotelRepoImpl) Unsub(hId int, uId int) error {
	_, err := c.DB.Exec(unsubHotelQ, uId, hId)
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
	getAllQ          = "select h.hotel_id, h.name, h.avaible, u.user_id, u.username from hotel h inner join usr as u on u.user_id = h.user_id"
	getHotelByUserQ  = "select h.hotel_id, h.name, h.avaible from hotel h where h.user_id = $1"
	subHotelQ        = "insert into usr_hotel (user_id, hotel_id) values($1,$2) returning *"
	unsubHotelQ      = "delete from usr_hotel t where t.user_id = $1 and t.hotel_id = $2"
	selectHotelQ     = "select h.hotel_id, h.name, h.avaible from hotel h where h.hotel_id = $1"
	saveHotelQ       = "insert into hotel (name, avaible, user_id) values ($1,$2,$3) returning hotel_id, name, avaible"
	getUserViewbyIdQ = "select user_id, username from usr where user_id = $1"
	deleteHoteQ      = "delete from hotel h where h.hotel_id = $1 and h.user_id = $2"
	addCommentQ      = "insert into comment (text, hotel_id, user_id) values ($1,$2,$3) returning comment_id, text"
)
