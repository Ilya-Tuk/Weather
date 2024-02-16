package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Ilya-Tuk/Weather/internal/models"
	"github.com/Ilya-Tuk/Weather/internal/outerApis"
	"github.com/gin-gonic/gin"
)

func (s *Rest) createUser(ctx *gin.Context) {
	var user models.User
	err := ctx.BindJSON(&user)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ok := s.service.CreateNewUser(ctx.GetHeader("Name"), ctx.GetHeader("Password"))

	ctx.JSON(http.StatusOK, models.StandartRequest{Status: ok})
}

func (s *Rest) userExists(ctx *gin.Context) {
	ok := s.service.UserExists(ctx.Param("name"))

	ctx.JSON(http.StatusOK, models.StandartRequest{Status: ok})
}

func (s *Rest) usersFavourites(ctx *gin.Context) {
	favs, ok := s.service.GetUsersFavourites(ctx.Param("name"))

	ctx.JSON(http.StatusOK, struct {
		Favs   []string
		Status bool
	}{
		Favs:   favs,
		Status: ok,
	})
}

func (s *Rest) addUsersFavourites(ctx *gin.Context) {
	var note string
	err := ctx.BindJSON(&note)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ok := s.service.AddUsersFavourite(ctx.Param("name"), note)

	ctx.JSON(http.StatusOK, models.StandartRequest{Status: ok})
}

func (s *Rest) deleteUsersFavourite(ctx *gin.Context) {
	var note string
	err := ctx.BindJSON(&note)

	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ok := s.service.DeleteUsersFavourite(ctx.Param("name"), note)

	ctx.JSON(http.StatusOK, models.StandartRequest{Status: ok})
}

func (s *Rest) getWeather(ctx *gin.Context) {
	city := ctx.Query("city")

	reqCurr, reqFore := outerApis.GetWeather(ctx, city)

	if reqCurr.IsError() {
		ctx.AbortWithStatus(400)
		return
	}

	if reqFore.IsError() {
		ctx.AbortWithStatus(400)
		return
	}

	respCurr := make(map[string]interface{})
	respFore := make(map[string]interface{})
	userCurr := make(map[string]string)
	userFore := []map[string]string{}

	json.Unmarshal(reqCurr.Body(), &respCurr)
	json.Unmarshal(reqFore.Body(), &respFore)

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
