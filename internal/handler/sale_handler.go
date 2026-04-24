package handler

import (
	"net/http"
	"product-api-go/internal/domain"
	"product-api-go/internal/usecase"

	"github.com/gin-gonic/gin"
)

type SaleHandler struct {
	usecase *usecase.SaleUsecase
}

func NewSaleHandler(u *usecase.SaleUsecase) *SaleHandler {
	return &SaleHandler{usecase: u}
}

// CreateSale godoc
// @Summary Process a new product sale
// @Description Deduct stock and record a sale
// @Tags sales
// @Accept json
// @Produce json
// @Param sale body domain.Sale true "Sale object"
// @Success 201 {object} domain.Sale
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /sales [post]
func (h *SaleHandler) CreateSale(c *gin.Context) {
	var sale domain.Sale
	if err := c.ShouldBindJSON(&sale); err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{Message: "invalid request body"})
		return
	}
	if err := h.usecase.ProcessSale(&sale); err != nil {
		handleError(c, err) // Re-use the existing error handler
		return
	}
	c.JSON(http.StatusCreated, sale)
}
