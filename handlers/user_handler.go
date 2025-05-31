package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/theweird-kid/property-list/models"
	"github.com/theweird-kid/property-list/services/user_service"
)

func (api *API) GetUsers(ctx *gin.Context) {
	users, err := user_service.GetUsers(ctx, api.MongoDB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (api *API) RegisterUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := user_service.RegisterUser(ctx, user, api.MongoDB)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, "Registered")
}

func (api *API) LoginUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	if user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, "Required email and password!")
		return
	}

	token, err := user_service.LoginUser(ctx, &user, api.MongoDB)
	if err != nil && err.Error() == "invalid" {
		ctx.JSON(http.StatusNotFound, "user with credentials doesn't exist")
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Internal problem")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user.Email,
	})
}
