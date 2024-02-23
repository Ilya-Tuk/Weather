package memory

import (
	"encoding/json"
	"os"

	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/Ilya-Tuk/Weather/internal/models"
)

type Repository []models.FullUser

var userDoesntExists error
var cfg = config.Read()

func (rep *Repository) AddUser(user models.User) bool {
	*rep = append(*rep, models.FullUser{User: user, Favourites: []string{}, Alerts: []models.Alert{}})
	return true
}

func (rep *Repository) UserExist(name string) bool {
	for _, el := range *rep {
		if el.User.Name == name {
			return true
		}
	}
	return false
}

func (rep *Repository) FindUser(name string) (models.FullUser, bool) {
	for _, el := range *rep {
		if el.User.Name == name {
			return el, true
		}
	}
	return models.FullUser{User: models.User{Name: "", Password: ""}, Alerts: []models.Alert{}, Favourites: []string{}}, false
}

func (rep *Repository) GetUsersFavourite(name string) ([]string, error) {
	for _, el := range *rep {
		if el.User.Name == name {
			return el.Favourites, nil
		}
	}
	return []string{}, userDoesntExists
}

func (rep *Repository) SetUsersFavourite(name string, favs []string) error {
	for i := range *rep {
		if (*rep)[i].User.Name == name {
			(*rep)[i].Favourites = favs
			return nil
		}
	}
	return userDoesntExists
}

func (rep *Repository) WalkByUsers(operation func(*models.FullUser)) {
	for i := range *rep {
		operation(&(*rep)[i])
	}
}

func (rep *Repository) Init() {
	if cfg.Db.MemoryFileName == "" {
		return
	}
	base, _ := os.Open(cfg.Db.MemoryFileName)
	defer base.Close()
	var buffer []byte

	base.Read(buffer)

	json.Unmarshal(buffer, rep)
}

func (rep *Repository) Close() {
	base, _ := os.Open(cfg.Db.MemoryFileName)
	defer base.Close()
	buffer, _ := json.Marshal(*rep)
	_, _ = base.Write(buffer)
}

func New() *Repository {
	new := Repository{}
	new.Init()
	return &new
}
