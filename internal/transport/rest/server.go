package rest

import (
	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/Ilya-Tuk/Weather/internal/services"
	"github.com/gin-gonic/gin"
)

type Rest struct {
	service services.Service
}

func NewServer(service services.Service) *gin.Engine {
	if config.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	rest := Rest{service}

	r.POST("/users", rest.createUser)
	r.GET("/users/:name/exists", rest.userExists)
	r.GET("/users/:name/favourites", rest.usersFavourites)
	r.POST("/users/:name/favourites", rest.addUsersFavourites)
	r.GET("/weather/current", rest.getWeather)
	r.DELETE("/users/:name/favourites", rest.deleteUsersFavourite)

	return r
}
