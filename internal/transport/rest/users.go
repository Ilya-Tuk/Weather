package rest

import (
	"encoding/json"
	"net/http"

	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/Ilya-Tuk/Weather/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var client = resty.New()

func (s *Rest) createUser(ctx *gin.Context) {
	var name struct{ name string }
	err := ctx.BindJSON(&name)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ok := s.service.CreateNewUser(name.name)

	if !ok {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func (s *Rest) userExists(ctx *gin.Context) {
	ok := s.service.UserExists(ctx.Param("name"))

	ctx.JSON(http.StatusOK, struct {
		Status bool
	}{
		Status: ok,
	})
}

func (s *Rest) usersFavourites(ctx *gin.Context) {
	favs, ok := s.service.GetUsersFavourites(ctx.Param("name"))

	ctx.JSON(http.StatusOK, struct {
		Favourites []models.Note
		Status     bool
	}{
		Favourites: favs,
		Status:     ok,
	})
}

func (s *Rest) addUsersFavourites(ctx *gin.Context) {
	var note models.Note
	err := ctx.BindJSON(&note)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ok := s.service.AddUsersFavourite(ctx.Param("name"), note)

	ctx.JSON(http.StatusOK, struct {
		Status bool
	}{
		Status: ok,
	})
}

func (s *Rest) deleteUsersFavourite(ctx *gin.Context) {
	var city struct {
		city string
	}
	err := ctx.BindJSON(&city)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ok := s.service.DeleteUsersFavourite(ctx.Param("name"), city.city)

	ctx.JSON(http.StatusOK, struct {
		Status bool
	}{
		Status: ok,
	})
}

func (s *Rest) getWeather(ctx *gin.Context) {
	city := ctx.Query("city")

	qC := map[string]string{"key": config.API_KEY, "city name": city}
	qF := map[string]string{"key": config.API_KEY, "city name": city, "days": "7"}

	respCurr := make(map[string]string)
	respFore := make(map[string]string)
	userCurr := make(map[string]string)
	userFore := []map[string]string

	respCurrent, err := client.R().
		SetQueryParams(qC).
		Get("http://api.weatherapi.com/v1/current.json")

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	if respCurrent.IsError() {
		ctx.AbortWithError(http.StatusBadRequest, nil)
	}

	respForecast, err := client.R().
		SetQueryParams(qF).
		Get("http://api.weatherapi.com/v1/forecast.json")

	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
	}
	if respForecast.IsError() {
		ctx.AbortWithError(http.StatusBadRequest, nil)
	}

	json.Unmarshal(respCurrent.Body(),&respCurr)
	json.Unmarshal(respForecast.Body(),&respFore)

	userCurr["temperature celsius"] = respCurr["temp_c"]
	userCurr["feelslike temperature celsius"] = respCurr["feelslike_c"]
	userCurr["wind speed"] = respCurr["wind_kph"]
	userCurr["wind direction"] = respCurr["wind_dir"]

}
