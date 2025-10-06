package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rezbow/ecommerce/internal/app/handlers"
	"github.com/rezbow/ecommerce/internal/app/services"
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

	// repo
	userRepo := database.NewUserRepo(db)

	// services
	userSvc := services.NewUserService(userRepo, cfg.JWTSecret)

	// handler
	userHandler := handlers.NewUserHandler(userSvc)

	router := gin.Default()

	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	protected := router.Group("/", middlewares.AuthMiddleware(cfg))
	{
		protected.POST("/profile", userHandler.Profile)
	}
	router.Run(":8080")
}
