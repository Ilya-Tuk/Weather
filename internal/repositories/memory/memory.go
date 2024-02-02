package memory

import (
	"encoding/json"
	"os"

	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/Ilya-Tuk/Weather/internal/models"
)

type Repository []models.User

func (rep *Repository) AddUser(token int64) bool {
	*rep = append(*rep, models.User{Token: token, Favourites: []models.Note{}})
	return true
}

func (rep *Repository) FindUser(token int64) bool {
	for _, el := range *rep {
		if el.Token == token {
			return true
		}
	}
	return false
}

func (rep *Repository) GetUsersFavourite(token int64) ([]models.Note, bool) {
	for _, el := range *rep {
		if el.Token == token {
			return el.Favourites, true
		}
	}
	return []models.Note{}, false
}

func (rep *Repository) SetUsersFavourite(token int64, favs []models.Note) bool {
	for _, el := range *rep {
		if el.Token == token {
			el.Favourites = favs
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
