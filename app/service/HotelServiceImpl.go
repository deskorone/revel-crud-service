package service

import (
	"testAuth/app/models"
	"testAuth/app/repo"
	"github.com/revel/revel"
)

type HotelServiceImpl struct {
	Repo *repo.Repository
}


var instanceHotelService HotelService

func (*HotelServiceImpl) AddCommentToHotel(c *revel.Controller, hID int) (*models.Comment, error){
	panic("sad")	
}


func (o *HotelServiceImpl) SaveHotel(h models.Hotel, c *revel.Controller) (*models.HotelResp, error) {
	id, err := instance.AuthService.CheckUser(c)
	if err != nil {
		c.Response.Status = 403
		return nil, err
	}
	hr, err := o.Repo.HotelRepo.SaveHotel(&h, id)
	if err != nil {
		return nil, err
	}
	return hr, nil

}

func (c *HotelServiceImpl) DeleteHotel(uId int, hId int) error {
	err := c.Repo.HotelRepo.DeleteHotel(uId, hId) 
	return err
}

func (o *HotelServiceImpl) GetHotelByUser(c *revel.Controller) ([]models.Hotel, error) {

	id, err := instance.AuthService.CheckUser(c)
	if err != nil {
		return nil, err
	}
	arr, err := o.Repo.HotelRepo.GetHotelsByUser(id)
	if err != nil {
		return nil, err
	}

	return arr, nil

}

func getHotelServiceImpl(r *repo.Repository) HotelService {
	if instanceHotelService == nil {
		instanceHotelService = &HotelServiceImpl{r}
	}
	return instanceHotelService
}
