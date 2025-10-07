package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rezbow/ecommerce/internal/app/models"
	"github.com/rezbow/ecommerce/internal/app/services"
)

type ProductHandler struct {
	productSvc services.IProductSvc
}

func NewProductHandler(productSvc services.IProductSvc) *ProductHandler {
	return &ProductHandler{
		productSvc: productSvc,
	}
}

func (handler *ProductHandler) GetProduct(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "product not found",
		})
		return
	}

	product, err := handler.productSvc.GetProduct(id)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "product not found",
			})
			return

		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ProductToProductResponse(*product))

}

func (handler *ProductHandler) ListProducts(ctx *gin.Context) {
	pagination := ExtractPagination(ctx)
	products, err := handler.productSvc.ListProducts(&pagination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	response := gin.H{
		"data": models.ProductsToProductsResponse(products),
		"metadata": gin.H{
			"page":  pagination.Page,
			"limit": pagination.Limit,
		},
	}
	ctx.JSON(http.StatusOK, response)
}

func (handler *ProductHandler) CreateProduct(ctx *gin.Context) {
	var productCreate models.ProductCreate
	if err := ctx.ShouldBindJSON(&productCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if valid, errs := productCreate.Validate(); !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errs})
		return
	}

	product, err := handler.productSvc.CreateProduct(&productCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, models.ProductToProductResponse(*product))
}

func (handler *ProductHandler) UpdateProduct(ctx *gin.Context) {
	productId, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var productUpdateRequest models.ProductUpdateRequest
	if err := ctx.ShouldBindJSON(&productUpdateRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if valid, errs := productUpdateRequest.Validate(); !valid {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errs})
		return
	}

	product, err := handler.productSvc.UpdateProduct(productId, &productUpdateRequest)
	if err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, models.ProductToProductResponse(*product))
}
