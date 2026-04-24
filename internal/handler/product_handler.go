package handler

import (
	"errors"
	"log"
	"net/http"

	"product-api-go/internal/domain"
	"product-api-go/internal/usecase"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HTTPError struct {
	Message string `json:"message" example:"The requested resource was not found"`
}

// DeleteSuccessResponse represents the successful response for a delete operation.
type DeleteSuccessResponse struct {
	Message string `json:"message" example:"deleted"`
}

type ProductHandler struct {
	usecase *usecase.ProductUsecase
}

func NewProductHandler(u *usecase.ProductUsecase) *ProductHandler {
	return &ProductHandler{usecase: u}
}

// GetProducts godoc
// @Summary Get all products
// @Description Get a list of all products
// @Tags products
// @Produce json
// @Success 200 {array} domain.Product
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.usecase.GetProducts()
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get a single product by ID
// @Description Get a single product by its unique ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} domain.Product
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /products/{id} [get]
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	p, err := h.usecase.GetProduct(id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Add a new product to the database
// @Tags products
// @Accept json
// @Produce json
// @Param product body domain.Product true "Product object"
// @Success 201 {object} domain.Product
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /products [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var p domain.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{Message: "invalid request body"})
		return
	}
	if err := h.usecase.CreateProduct(&p); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, p)
}

// UpdateProduct godoc
// @Summary Update an existing product
// @Description Update details of an existing product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body domain.Product true "Product object to update"
// @Success 200 {object} domain.Product
// @Failure 400 {object} HTTPError
// @Failure 404 {object} HTTPError
// @Router /products/{id} [put]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var p domain.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, HTTPError{Message: "invalid request body"})
		return
	}
	p.ID = id
	if err := h.usecase.UpdateProduct(&p); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, p)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product by its ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} DeleteSuccessResponse
// @Failure 404 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /products/{id} [delete]
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if err := h.usecase.DeleteProduct(id); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, DeleteSuccessResponse{Message: "deleted"})
}

// handleError centralizes error handling for all handlers in the package.
func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, HTTPError{Message: "product not found"})
	case errors.As(err, &usecase.ValidationError{}):
		c.JSON(http.StatusBadRequest, HTTPError{Message: err.Error()})
	case errors.Is(err, usecase.ErrInsufficientStock):
		c.JSON(http.StatusConflict, HTTPError{Message: err.Error()}) // 409 Conflict is a good status for this
	case errors.Is(err, usecase.ErrProductNotFoundForPurchase):
		c.JSON(http.StatusNotFound, HTTPError{Message: err.Error()})
	default:
		// Log the unexpected error for debugging.
		log.Printf("An unexpected error occurred: %v", err)
		c.JSON(http.StatusInternalServerError, HTTPError{Message: "an internal server error occurred"})
	}
}
