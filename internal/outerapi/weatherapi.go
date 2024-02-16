package outerapi

import (
	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/gin-gonic/gin"

	"github.com/go-resty/resty/v2"
)

var client = resty.New()

func GetWeather(ctx gin.Context, city string) (resty.Request, resty.Request, err) {

	respCurrent, _ := client.R().
		SetQueryParams(map[string]string{"key": config.API_KEY, "q": city}).
		Get("http://api.weatherapi.com/v1/current.json")

	respForecast, _ := client.R().
		SetQueryParams(map[string]string{"key": config.API_KEY, "q": city, "days": "3"}).
		Get("http://api.weatherapi.com/v1/forecast.json")

	return respCurrent, respForecast
}
