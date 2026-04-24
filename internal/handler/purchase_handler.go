package handler

import (
	"net/http"
	"product-api-go/internal/domain"
	"product-api-go/internal/usecase"

	"github.com/gin-gonic/gin"
)

type PurchaseHandler struct {
	usecase *usecase.PurchaseUsecase
}

func NewPurchaseHandler(u *usecase.PurchaseUsecase) *PurchaseHandler {
	return &PurchaseHandler{usecase: u}
}

// CreatePurchase godoc
// @Summary Process a new product purchase
// @Description Add stock and record a purchase
// @Tags purchases
// @Accept json
// @Produce json
// @Param purchase body domain.Purchase true "Purchase object"
// @Success 201 {object} domain.Purchase
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /purchases [post]
func (h *PurchaseHandler) CreatePurchase(c *gin.Context) {
	var purchase domain.Purchase
	if err := c.ShouldBindJSON(&purchase); err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{Message: "invalid request body"})
		return
	}
	if err := h.usecase.ProcessPurchase(&purchase); err != nil {
		handleError(c, err) // Re-use the existing error handler
		return
	}
	c.JSON(http.StatusCreated, purchase)
}
