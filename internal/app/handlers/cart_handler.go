package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezbow/ecommerce/internal/app/models"
	"github.com/rezbow/ecommerce/internal/app/services"
)

type CartHandler struct {
	cartSvc services.ICartService
}

func NewCartHandler(cartSvc services.ICartService) *CartHandler {
	return &CartHandler{
		cartSvc: cartSvc,
	}
}

func (handler *CartHandler) GetCart(ctx *gin.Context) {
	value, _ := ctx.Get("userId")
	userId, ok := value.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthroized user",
		})
		return
	}

	cart, err := handler.cartSvc.GetUserCart(userId)
	if err != nil {
		if errors.Is(err, services.ErrCartNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "no cart found for the user"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	ctx.JSON(http.StatusOK, cart)
}

func (handler *CartHandler) AddToCart(ctx *gin.Context) {
	value, _ := ctx.Get("userId")
	userId, ok := value.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthroized user",
		})
		return
	}

	var itemCart models.ItemCartRequest
	if err := ctx.ShouldBindJSON(&itemCart); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if valid, errs := itemCart.Validate(); !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": errs,
		})
		return
	}

	if err := handler.cartSvc.AddToUserCart(userId, &itemCart); err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return

		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (handler *CartHandler) DeleteItem(ctx *gin.Context) {
	value, _ := ctx.Get("userId")
	userId, ok := value.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthroized user",
		})
		return
	}
	productId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := handler.cartSvc.RemoveItemFromCart(userId, productId); err != nil {
		if errors.Is(err, services.ErrItemNotFound) {
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

	ctx.Status(http.StatusNoContent)
}

func (handler *CartHandler) ClearCart(ctx *gin.Context) {
	value, _ := ctx.Get("userId")
	userId, ok := value.(uuid.UUID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthroized user",
		})
		return
	}
	if err := handler.cartSvc.ClearCart(userId); err != nil {
		if errors.Is(err, services.ErrCartNotFound) {
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
	ctx.Status(http.StatusNoContent)
}
