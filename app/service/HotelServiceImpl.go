package service

import (
	"errors"
	"testAuth/app/models"
	"testAuth/app/repo"

	"github.com/revel/revel"
)

type HotelServiceImpl struct {
	Repo *repo.Repository
	ch   chan models.Hotel
}

func (c HotelServiceImpl) GetPaginationHotels(page, size int) ([]models.Hotel, int, error) {
	if page < 1 || size < 1 {
		return nil, 0, errors.New("inccorrect value")
	}
	return c.Repo.HotelRepo.GetPaginationHotels(page, size)
}

// GetAllHotels implements HotelService
func (c HotelServiceImpl) GetAllHotels() ([]models.Hotel, error) {

	arr, err := c.Repo.HotelRepo.GetAllHotels()
	return arr, err
}

// SaveHotelWithoutUser implements HotelService
func (c *HotelServiceImpl) SaveHotelWithoutUser(h models.Hotel) (*models.Hotel, error) {
	hotel := &models.Hotel{}
	hotel, err := c.Repo.HotelRepo.SaveHotelWithoutUser(h)
	if err != nil {
		return nil, err
	}
	select {
	case c.ch <- h:
	}
	return hotel, nil
}

// ParseHotelsFromUrl implements HotelService
func (c *HotelServiceImpl) ParseHotelsFromUrl(arr []models.Hotel) ([]models.Hotel, error) {

	return arr, nil
}

var instanceHotelService HotelService

// AddCommentToHotel Добавить комментарий к отелю, принимает id отеля и текст коментария (доступно только авторизированному пользователю)
func (q *HotelServiceImpl) AddCommentToHotel(c *revel.Controller, hotelID int, text string) (*models.Comment, error) {
	u, err := instance.AuthService.GetUser(c)
	if err != nil {
		return nil, err
	}
	r, err := q.Repo.HotelRepo.AddComment(hotelID, *u, text)
	return r, err
}

// SaveHotel Созранение отеля принимает modes.Hotel
func (o *HotelServiceImpl) SaveHotel(hotel models.Hotel, c *revel.Controller) (*models.HotelResp, error) {
	u, err := instance.AuthService.GetUser(c)
	if err != nil {
		c.Response.Status = 403
		return nil, err
	}
	hr, err := o.Repo.HotelRepo.SaveHotel(&hotel, u.Id)
	return hr, err

}

// DeleteHotel удаляет отель примает id владельца и отеля
func (c *HotelServiceImpl) DeleteHotel(userId int, hotelId int) error {
	err := c.Repo.HotelRepo.DeleteHotel(userId, hotelId)
	return err
}

// GetHotelByUser Получение своих отелей
func (o *HotelServiceImpl) GetHotelByUser(c *revel.Controller) ([]models.Hotel, error) {

	user, err := instance.AuthService.GetUser(c)
	if err != nil {
		return nil, err
	}
	hotels, err := o.Repo.HotelRepo.GetHotelsByUser(user.Id)

	return hotels, err

}

func getHotelServiceImpl(r *repo.Repository, ch chan models.Hotel) HotelService {
	if instanceHotelService == nil {
		instanceHotelService = &HotelServiceImpl{r, ch}
	}
	return instanceHotelService
}
