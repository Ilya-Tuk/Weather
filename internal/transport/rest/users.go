package rest

import (
	"net/http"

	"github.com/Ilya-Tuk/Weather/internal/models"
	"github.com/gin-gonic/gin"
)

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
	favs,ok := s.service.GetUsersFavourites(ctx.Param("name"))

	ctx.JSON(http.StatusOK, struct {
		Favourites []models.Note
		Status bool
	}{
		Favourites: favs,
		Status: ok,
	})
}

func (s *Rest)addUsersFavourites(ctx *gin.Context) {
	var note models.Note
	err := ctx.BindJSON(&note)
	
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}

	ok := s.service.AddUsersFavourite(ctx.Param("name"),note)

	ctx.JSON(http.StatusOK, struct {
		Status bool
	}{
		Status: ok,
	})
}

