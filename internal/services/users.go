package services

import (
	"fmt"

	"github.com/Ilya-Tuk/Weather/internal/models"
)

type UsersRepository interface {
	AddUser(string) bool
	FindUser(string) bool
	Close()
	SetUsersFavourite(string, []models.Note) bool
	GetUsersFavourite(string) ([]models.Note, bool)
	Init()
}

type Service struct {
	repo UsersRepository
}

func New(repo UsersRepository) Service {
	return Service{
		repo: repo,
	}
}

func (serv *Service) CreateNewUser(name string) bool {
	if serv.repo.FindUser(name) {
		fmt.Println("found same user")
		return false
	}

	return serv.repo.AddUser(name)
}

func (serv *Service) UserExists(name string) bool {
	return serv.repo.FindUser(name)
}

func (serv *Service) GetUsersFavourites(name string) ([]models.Note, bool) {
	return serv.repo.GetUsersFavourite(name)
}

func (serv *Service) AddUsersFavourite(name string, fav models.Note) bool {
	if !serv.repo.FindUser(name) {
		return false
	}

	favs, _ := serv.repo.GetUsersFavourite(name)
	favs = append(favs, fav)

	return serv.repo.SetUsersFavourite(name, favs)
}

func (serv *Service) DeleteUsersFavourite(name string, city string) bool {
	if !serv.repo.FindUser(name) {
		return false
	}

	favs, _ := serv.repo.GetUsersFavourite(name)

	aimFav := -1

	for i, el := range favs {
		if el.City == city {
			aimFav = i
			break
		}
	}

	if aimFav == -1 {
		return false
	}

	favs = append(favs[:aimFav], favs[aimFav+1:]...)

	return serv.repo.SetUsersFavourite(name, favs)
}

func (serv *Service) Close() {
	serv.repo.Close()
}

func (serv *Service) Init() {
	serv.repo.Init()
}
