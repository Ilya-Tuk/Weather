package rest

import (
	"io"
	"net/http"

	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/Ilya-Tuk/Weather/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Rest struct {
	lg      *zap.SugaredLogger
	service services.Service
}

func NewServer(cfg config.Config, service services.Service, lg *zap.SugaredLogger) *http.Server {
	if cfg.ServCfg.DebugMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	gin.DefaultWriter = io.Discard
	r := gin.Default()

	rest := Rest{
		service: service,
		lg:      lg,
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}

	r.Use(func(ctx *gin.Context) {
		lg.Info("http request", ctx.Request.URL.Path)
	})

	r.Use(cors.New(config))

	r.POST("/users", rest.createUser)
	r.GET("/users/:name/exists", rest.userExists)
	r.GET("/users/:name/favourites", rest.usersFavourites)
	r.POST("/users/:name/afavourites", rest.addUsersFavourites)
	r.GET("/weather/current", rest.getWeather)
	r.DELETE("/users/:name/dfavourites", rest.deleteUsersFavourite)

	return &http.Server{
		Addr:    cfg.ServCfg.ServerHost,
		Handler: r,
	}
}
