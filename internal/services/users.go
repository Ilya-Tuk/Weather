package services

import (
	"strconv"

	"github.com/Ilya-Tuk/Weather/internal/models"
)

type UsersRepository interface {
	AddUser(int64) bool
	FindUser(int64) bool
	Close()
	SetUsersFavourite(token int64, favs []models.Note) bool
	GetUsersFavourite(token int64) ([]models.Note, bool)
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
	intname, err := strconv.Atoi(name)
	if err != nil {
		return false
	}
	token := int64(intname)
	if serv.repo.FindUser(token) {
		return false
	}

	return serv.repo.AddUser(token)
}

func (serv *Service) UserExists(name string) bool {
	intname, err := strconv.Atoi(name)
	if err != nil {
		return false
	}
	token := int64(intname)

	return serv.repo.FindUser(token)
}

func (serv *Service) GetUsersFavourites(name string) ([]models.Note, bool) {
	intname, err := strconv.Atoi(name)
	if err != nil {
		return []models.Note{}, false
	}
	token := int64(intname)
	return serv.repo.GetUsersFavourite(token)
}

func (serv *Service) AddUsersFavourite(name string, fav models.Note) bool {
	intname, err := strconv.Atoi(name)
	if err != nil {
		return false
	}

	token := int64(intname)

	if !serv.repo.FindUser(token){
		return false
	}

	favs,_:=  serv.repo.GetUsersFavourite(token)
	favs = append(favs, fav)

	return serv.repo.SetUsersFavourite(token, favs)
}

func (serv *Service) DeleteUsersFavourite(name string, city string) bool {
	intname, err := strconv.Atoi(name)
	if err != nil {
		return false
	}

	token := int64(intname)

	if !serv.repo.FindUser(token){
		return false
	}

	favs,_:=  serv.repo.GetUsersFavourite(token)

	aimFav := -1

	for i,el := range favs{
		if el.City == city{
			aimFav = i
			break
		}
	}

	if aimFav == -1{
		return false
	}

	favs = append(favs[:aimFav], favs[aimFav+1:]...)

	return serv.repo.SetUsersFavourite(token, favs)
}

func (serv *Service) Close() {
	serv.repo.Close()
}

func (serv *Service) Init() {
	serv.repo.Init()
}
