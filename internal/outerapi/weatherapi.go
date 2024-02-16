package outerapi

import (
	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/gin-gonic/gin"

	"github.com/go-resty/resty/v2"
)

var client = resty.New()
var cfg = config.Read()

func GetWeather(ctx *gin.Context, city string) (resty.Response, resty.Response) {

	respCurrent, _ := client.R().
		SetQueryParams(map[string]string{"key": cfg.ServCfg.API_KEY, "q": city}).
		Get(cfg.ServCfg.WeatherApiUrl + "current.json")

	respForecast, _ := client.R().
		SetQueryParams(map[string]string{"key": cfg.ServCfg.API_KEY, "q": city, "days": "3"}).
		Get(cfg.ServCfg.WeatherApiUrl + "forecast.json")

	return *respCurrent, *respForecast
}
