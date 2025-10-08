package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rezbow/ecommerce/internal/app/handlers"
	"github.com/rezbow/ecommerce/internal/app/services"
	"github.com/rezbow/ecommerce/internal/platform/cache"
	"github.com/rezbow/ecommerce/internal/platform/config"
	"github.com/rezbow/ecommerce/internal/platform/database"
	"github.com/rezbow/ecommerce/internal/platform/middlewares"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	db, err := database.ConnectDB(cfg)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	redis, err := cache.ConnectRedis()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// repo
	userRepo := database.NewUserRepo(db)
	productRepo := database.NewProductRepo(db)
	cartRepo := database.NewCartRepoRedis(redis)

	// services
	userSvc := services.NewUserService(userRepo, cfg.JWTSecret)
	productSvc := services.NewProductService(productRepo)
	cartSvc := services.NewCartService(cartRepo, productRepo)

	// handler
	userHandler := handlers.NewUserHandler(userSvc)
	productHandler := handlers.NewProductHandler(productSvc)
	cartHandler := handlers.NewCartHandler(cartSvc)

	// middlewares
	authMiddleware := middlewares.AuthMiddleware(cfg)

	router := gin.Default()

	router.GET("/products/:id", productHandler.GetProduct)
	router.GET("/products", productHandler.ListProducts)

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	protected := router.Group("/")
	protected.Use(authMiddleware)
	{
		protected.POST("/profile", userHandler.Profile)
		// endpoints for cart operations
		protected.GET("/cart", cartHandler.GetCart)                // getting user's cart information
		protected.POST("/cart", cartHandler.AddToCart)             // adding an item to cart
		protected.PUT("/cart/:id", cartHandler.UpdateItemQuantity) // update quantity of an item in cart
		protected.DELETE("/cart/:id", cartHandler.DeleteItem)      // delete a specific item with id in the cart
		protected.DELETE("/cart", cartHandler.ClearCart)           // delete the entire cart
		// endpoint for user's order details
		protected.GET("/orders")
		protected.GET("/orders/:id")
		// endpoint for checkingout the cart
		protected.POST("/checkout")
	}

	admin := router.Group("/admin")
	admin.Use(authMiddleware, middlewares.AdminMiddleware())
	{
		admin.POST("/products", productHandler.CreateProduct)
		admin.PUT("/products/:id", productHandler.UpdateProduct)
	}

	router.Run(":8080")
}
