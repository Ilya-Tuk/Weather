package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ilya-Tuk/Weather/internal/config"
	"github.com/Ilya-Tuk/Weather/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

var client = resty.New()

func (s *Rest) createUser(ctx *gin.Context) {
	var name map[string]string
	err := ctx.BindJSON(&name)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, ok := name["name"]

	if !ok {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ok = s.service.CreateNewUser(name["name"])

	ctx.JSON(http.StatusOK, struct {
		Status bool
	}{
		Status: ok,
	})
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
	var note map[string]string
	err := ctx.BindJSON(&note)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	_, ok := note["City"]
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	_, ok = note["Note"]
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ok = s.service.AddUsersFavourite(ctx.Param("name"), models.Note{City: note["City"], Note: note["Note"]})

	ctx.JSON(http.StatusOK, struct {
		Status bool
	}{
		Status: ok,
	})
}

func (s *Rest) deleteUsersFavourite(ctx *gin.Context) {
	var note map[string]string
	err := ctx.BindJSON(&note)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	_, ok := note["City"]
	if !ok {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ok = s.service.DeleteUsersFavourite(ctx.Param("name"), note["City"])

	ctx.JSON(http.StatusOK, struct {
		Status bool
	}{
		Status: ok,
	})
}

func (s *Rest) getWeather(ctx *gin.Context) {
	city := ctx.Query("city")

	qC := map[string]string{"key": config.API_KEY, "q": city}
	qF := map[string]string{"key": config.API_KEY, "q": city, "days": "7"}

	respCurr := make(map[string]interface{})
	respFore := make(map[string]interface{})
	userCurr := make(map[string]string)
	userFore := []map[string]string{}

	respCurrent, _ := client.R().
		SetQueryParams(qC).
		Get("http://api.weatherapi.com/v1/current.json")
	if respCurrent.IsError() {
		ctx.AbortWithStatus(400)
		return
	}

	respForecast, _ := client.R().
		SetQueryParams(qF).
		Get("http://api.weatherapi.com/v1/forecast.json")
	if respForecast.IsError() {
		ctx.AbortWithStatus(400)
		return
	}

	json.Unmarshal(respCurrent.Body(), &respCurr)
	json.Unmarshal(respForecast.Body(), &respFore)

	decodeCurr := []string{"temp_c", "feelslike_c", "wind_kph", "wind_dir", "pressure_mb", "precip_mm", "pressure_mb"}

	for _, el := range decodeCurr {
		var val = respCurr["current"].(map[string]interface{})[el]
		switch val.(type) {
		case int:
			userCurr[el] = fmt.Sprint(val.(int))
		case float64:
			userCurr[el] = fmt.Sprint(val.(float64))
		case string:
			userCurr[el] = val.(string)
		}
	}

	for _, el := range respFore["forecast"].(map[string]interface{})["forecastday"].([]interface{}) {
		dayweath := el.(map[string]interface{})["day"].(map[string]interface{})
		tempmap := make(map[string]string)

		tempmap["avgtemp_c"] = fmt.Sprint(dayweath["avgtemp_c"].(float64))
		tempmap["totalprecip_mm"] = fmt.Sprint(dayweath["totalprecip_mm"].(float64))
		tempmap["maxwind_kph"] = fmt.Sprint(dayweath["maxwind_kph"].(float64))

		userFore = append(userFore, tempmap)
	}

	ctx.JSON(http.StatusAccepted, struct {
		UserCurr map[string]string
		UserFore []map[string]string
	}{
		UserCurr: userCurr,
		UserFore: userFore,
	})
}
