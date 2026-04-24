package infrastructure

import (
	_ "product-api-go/docs"
	"product-api-go/internal/handler"
	"product-api-go/internal/repository"
	"product-api-go/internal/usecase"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Database connection
	db := NewPostgresDB()

	// Dependency Injection
	repo := repository.NewProductRepo(db)
	productUsecase := usecase.NewProductUsecase(repo)
	productHandler := handler.NewProductHandler(productUsecase)

	// New Sales components
	saleRepo := repository.NewSaleRepo(db)
	saleUsecase := usecase.NewSaleUsecase(repo, saleRepo, db)
	saleHandler := handler.NewSaleHandler(saleUsecase)

	// New Purchases components
	purchaseRepo := repository.NewPurchaseRepo(db)
	purchaseUsecase := usecase.NewPurchaseUsecase(repo, purchaseRepo, db)
	purchaseHandler := handler.NewPurchaseHandler(purchaseUsecase)

	// Routes
	api := r.Group("/api/products")
	{
		api.GET("", productHandler.GetProducts)
		api.GET("/:id", productHandler.GetProduct)
		api.POST("", productHandler.CreateProduct)
		api.PUT("/:id", productHandler.UpdateProduct)
		api.DELETE("/:id", productHandler.DeleteProduct)
	}

	// Sales Routes
	r.POST("/api/sales", saleHandler.CreateSale)

	// Purchases Routes
	r.POST("/api/purchases", purchaseHandler.CreatePurchase)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
