package service

import (
	"fmt"
	"net/http"
	// "sync"
	"github.com/PuerkitoBio/goquery"
	"github.com/revel/revel"
	"testAuth/app/models"
	"testAuth/app/repo"
	"time"
)

type HotelServiceImpl struct {
	Repo *repo.Repository
}

// Константы которые равны порядковым номерам нужных полей при парсинге селектора отеля
const (
	nameCounter    = 0
	avaibleCounter = 3
)

// SaveHotelWithoutUser implements HotelService
func (c *HotelServiceImpl) SaveHotelWithoutUser(h models.Hotel) (*models.Hotel, error) {
	return c.Repo.HotelRepo.SaveHotelWithoutUser(h)

}

// ParseHotelsFromUrl implements HotelService
func (c *HotelServiceImpl) ParseHotelsFromUrl(url string) ([]models.Hotel, error) {
	res, err := http.Get(url)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	
	arr := make([]models.Hotel, 0)
	doc.Find(".hotelInfo_desktop").Each(func(i int, s *goquery.Selection) {
		h := models.Hotel{}
		s.Children().Each(func(i int, s *goquery.Selection) {
			switch i {
			case nameCounter:
				h.Name = s.Text()
			case avaibleCounter:
				h.Avaible = i
			}
		})
		hres, err := c.Repo.HotelRepo.SaveHotelWithoutUser(h)
		if err != nil {
			return
		}
		arr = append(arr, *hres)
	})
	return arr, nil
}

var instanceHotelService HotelService

func (q *HotelServiceImpl) AddCommentToHotel(c *revel.Controller, hID int, text string) (*models.Comment, error) {
	u, err := instance.AuthService.GetUser(c)
	if err != nil {
		return nil, err
	}
	r, err := q.Repo.HotelRepo.AddComment(hID, *u, text)
	return r, err
}

func (o *HotelServiceImpl) SaveHotel(h models.Hotel, c *revel.Controller) (*models.HotelResp, error) {
	u, err := instance.AuthService.GetUser(c)
	if err != nil {
		c.Response.Status = 403
		return nil, err
	}
	hr, err := o.Repo.HotelRepo.SaveHotel(&h, u.Id)
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

	u, err := instance.AuthService.GetUser(c)
	if err != nil {
		return nil, err
	}
	arr, err := o.Repo.HotelRepo.GetHotelsByUser(u.Id)
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
