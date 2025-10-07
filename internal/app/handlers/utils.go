// internal/app/handlers/pagination_util.go
package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rezbow/ecommerce/internal/app/models"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
	MaxLimit     = 100 // Maximum number of items allowed per page
)

func ExtractPagination(c *gin.Context) models.Pagination {

	// Default values
	page := DefaultPage
	limit := DefaultLimit

	// 1. Parse 'page' parameter
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// 2. Parse 'limit' parameter
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// 3. Enforce MaxLimit
	if limit > MaxLimit {
		limit = MaxLimit
	}

	// 4. Calculate Offset
	// Offset = (Page - 1) * Limit
	offset := (page - 1) * limit
	if offset < 0 {
		offset = 0 // Should only happen if DefaultPage was changed to 0 or less
	}

	return models.Pagination{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	}
}
