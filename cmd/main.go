package main

import (
	"exam-test/config"
	"exam-test/internal/handlers"
	"exam-test/internal/repositories"
	"exam-test/internal/routes"
	"exam-test/internal/services"
	"exam-test/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	env := config.NewEnvConfig()
	databaseConfig := config.NewDatabaseConfig(env)
	db := databaseConfig.Connect()

	jwtUtils := utils.NewJWTUtils(env)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	authHandler := handlers.NewAuthHandler(userService, jwtUtils)
	authRoutes := routes.NewAuthRoute(authHandler)

	api := r.Group("/api/v1")

	api.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(
			http.StatusOK,
			gin.H{"message": "Server running on PORT 8080"},
		)
	})

	authRoutes.Routes(r)
	r.Run(":8080")
}
