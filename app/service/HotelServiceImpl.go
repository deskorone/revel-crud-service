package service

import (
	"testAuth/app"
	"testAuth/app/models"
)

type HotelServiceImpl struct {
}

const addHotelQ = "insert into hotel (name, avaible) values($1, $2)"
const deleteHotel = "delete from hotel h  where h.user_id = $1 and h.hotel_id = $2"

var instanceHotelService HotelService

func (c *HotelServiceImpl) SaveHotel(h models.Hotel) (*models.Hotel, error) {
	err := app.DB.QueryRow(addHotelQ, h.Name, h.Avaible).Scan(&h.ID, &h.Name, &h.Avaible)

	if err != nil {
		return nil, err
	}
	return &h, nil

}

func (c *HotelServiceImpl) DeleteHotel(uId int, hId int) error {
	_, err := app.DB.Exec(deleteHotel, uId, hId)
	if err != nil {
		return err
	}
	return nil
}


func getHotelServiceImpl () HotelService{
	once.Do(func() {
		instanceHotelService = new(HotelServiceImpl)
	})
	return instanceHotelService
}
