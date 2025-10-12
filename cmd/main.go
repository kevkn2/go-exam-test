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
	if env.MODE == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	databaseConfig := config.NewDatabaseConfig(env)
	db := databaseConfig.Connect()

	jwtUtils := utils.NewJWTUtils(env)
	questionUtils := utils.NewQuestionUtils()

	userRepository := repositories.NewUserRepository(db)
	studentRepository := repositories.NewStudentRepository(db)
	userService := services.NewAuthService(userRepository, studentRepository, db)
	authHandler := handlers.NewAuthHandler(userService, jwtUtils)
	authRoutes := routes.NewAuthRoute(authHandler)

	onlineTestRepo := repositories.NewOnlineTestRepo(db)
	questionRepo := repositories.NewQuestionRepo(db)
	onlineTestService := services.NewOnlineTestService(onlineTestRepo, questionRepo, questionUtils)
	onlineTestHandler := handlers.NewOnlineTestHandler(onlineTestService, jwtUtils)
	onlineTestRoute := routes.NewOnlineTestRoute(onlineTestHandler)

	api := r.Group("/api/v1")

	api.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(
			http.StatusOK,
			gin.H{"message": "Server running on PORT 8080"},
		)
	})

	authRoutes.Routes(r)
	onlineTestRoute.Routes(r)
	r.Run(":8080")
}
