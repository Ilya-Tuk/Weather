package memory

import (
	"encoding/json"
	"os"

	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/Ilya-Tuk/Weather/internal/models"
)

type Repository []models.User

var userDoesntExists error
var cfg = config.Read()

func (rep *Repository) AddUser(user models.User) bool {
	*rep = append(*rep, user)
	return true
}

func (rep *Repository) FindUser(name string) bool {
	for _, el := range *rep {
		if el.Name == name {
			return true
		}
	}
	return false
}

func (rep *Repository) GetUsersFavourite(name string) ([]string, error) {
	for _, el := range *rep {
		if el.Name == name {
			return el.Favourites, nil
		}
	}
	return []string{}, userDoesntExists
}

func (rep *Repository) SetUsersFavourite(name string, favs []string) error {
	for i := range *rep {
		if (*rep)[i].Name == name {
			(*rep)[i].Favourites = favs
			return nil
		}
	}
	return userDoesntExists
}

func (rep *Repository) Init() {
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
