package services

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Ilya-Tuk/Weather/internal/models"
	"github.com/Ilya-Tuk/Weather/internal/outerapi"
	"github.com/gin-gonic/gin"
)

type UsersRepository interface {
	AddUser(models.User) bool
	UserExist(string) bool
	SetUsersFavourite(string, []string) error
	GetUsersFavourite(string) ([]string, error)
	FindUser(string) (models.FullUser, bool)
	WalkByUsers(func(*models.FullUser))
}

type Service struct {
	repo UsersRepository
}

func New(repo UsersRepository) Service {
	return Service{
		repo: repo,
	}
}

func (serv *Service) CreateNewUser(name string, password string) bool {
	if serv.repo.UserExist(name) {
		return false
	}

	return serv.repo.AddUser(models.User{Name: name, Password: password})
}

func (serv *Service) UserExists(name string) bool {
	return serv.repo.UserExist(name)
}

func (serv *Service) FindUser(name string) (models.FullUser, bool) {
	return serv.repo.FindUser(name)
}

func (serv *Service) GetUsersFavourites(name string) ([]string, bool) {
	favs, err := serv.repo.GetUsersFavourite(name)
	if err != nil {
		return favs, false
	}
	return favs, true
}

func (serv *Service) AddUsersFavourite(name string, fav string) bool {
	if !serv.repo.UserExist(name) {
		return false
	}

	favs, _ := serv.repo.GetUsersFavourite(name)
	favs = append(favs, fav)

	err := serv.repo.SetUsersFavourite(name, favs)

	return err != nil
}

func (serv *Service) DeleteUsersFavourite(name string, fav string) bool {
	if !serv.repo.UserExist(name) {
		return false
	}

	favs, _ := serv.repo.GetUsersFavourite(name)

	aimFav := -1

	for i, el := range favs {
		if el == fav {
			aimFav = i
			break
		}
	}

	if aimFav == -1 {
		return false
	}

	favs = append(favs[:aimFav], favs[aimFav+1:]...)

	err := serv.repo.SetUsersFavourite(name, favs)

	return err != nil
}

func (serv *Service) SetAlerts(context.Context) {
	serv.repo.WalkByUsers(func(user *models.FullUser) {
		user.Alerts = []models.Alert{}
	})

	serv.repo.WalkByUsers(func(user *models.FullUser) {
		if len(user.Favourites) == 0 {
			return
		}
		for _, el := range user.Favourites {
			Current, _ := outerapi.GetWeather(&gin.Context{}, el)

			if Current.IsError() {
				continue
			}

			NormalCurrent := make(models.Intermap)

			json.Unmarshal(Current.Body(), &NormalCurrent)

			if NormalCurrent["Current"].(models.Intermap)["precip_mm"].(float64) > 0.5 {
				user.Alerts = append(user.Alerts, models.Alert{City: el, Date: time.Now()})
			}
		}
	})
}
