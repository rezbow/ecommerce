package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezbow/ecommerce/internal/app/models"
	"github.com/rezbow/ecommerce/internal/app/services"
)

type UserHandler struct {
	userSvc models.UserSvc
}

func NewUserHandler(userSvc models.UserSvc) *UserHandler {
	return &UserHandler{
		userSvc: userSvc,
	}
}

func (handler *UserHandler) Register(ctx *gin.Context) {
	var req models.RegisterUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := handler.userSvc.RegisterUser(&req)
	if err != nil {
		if errors.Is(err, services.ErrDuplicateEmail) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	response := gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}

	ctx.JSON(http.StatusCreated, response)
}

func (handler *UserHandler) Login(ctx *gin.Context) {
	var loginData models.Login
	if err := ctx.ShouldBindJSON(&loginData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := handler.userSvc.Authenticate(&loginData)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}
	response := gin.H{
		"token": token,
	}
	ctx.JSON(http.StatusOK, response)
}

func (handler *UserHandler) Profile(ctx *gin.Context) {
	userIdStr, _ := ctx.Get("userId")
	userId, _ := userIdStr.(uuid.UUID)

	user, err := handler.userSvc.GetUser(userId)
	if err != nil {
		if errors.Is(err, services.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	response := gin.H{
		"id":         user.ID,
		"email":      user.Email,
		"is_admin":   user.IsAdmin,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
	ctx.JSON(http.StatusOK, response)
}
