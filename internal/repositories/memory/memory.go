package memory

import (
	"encoding/json"
	"os"

	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/Ilya-Tuk/Weather/internal/models"
)

type Repository []models.User

func (rep *Repository) AddUser(name string) bool {
	*rep = append(*rep, models.User{Name: name, Favourites: []models.Note{}})
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

func (rep *Repository) GetUsersFavourite(name string) ([]models.Note, bool) {
	for _, el := range *rep {
		if el.Name == name {
			return el.Favourites, true
		}
	}
	return []models.Note{}, false
}

func (rep *Repository) SetUsersFavourite(name string, favs []models.Note) bool {
	for i := range *rep {
		if (*rep)[i].Name == name {
			(*rep)[i].Favourites = favs
			return true
		}
	}
	return false
}

func (rep *Repository) Init() {
	base, _ := os.Open(config.MemoryFileName)
	defer base.Close()
	var buffer []byte

	base.Read(buffer)

	json.Unmarshal(buffer, rep)
}

func (rep *Repository) Close() {
	base, _ := os.Open(config.MemoryFileName)
	defer base.Close()
	buffer, _ := json.Marshal(*rep)
	_, _ = base.Write(buffer)
}
