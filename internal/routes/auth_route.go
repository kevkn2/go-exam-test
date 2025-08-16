package routes

import (
	"exam-test/internal/handlers"

	"github.com/gin-gonic/gin"
)

type AuthRoute interface {
	Routes(route *gin.Engine)
}

type authRoute struct {
	authHandler handlers.AuthHandler
}

// authRoutes implements AuthRoute.
func (a *authRoute) Routes(route *gin.Engine) {
	auth := route.Group("/api/v1/auth")

	auth.POST("/login", a.authHandler.AuthenticateUser)
	auth.POST("/registerAdmin", a.authHandler.RegisterAdmin)
	auth.POST("/registerStudent", a.authHandler.RegisterStudent)
	auth.GET("/validateAdmin", a.authHandler.ValidAdmin)
}

func NewAuthRoute(authHandler handlers.AuthHandler) AuthRoute {
	return &authRoute{
		authHandler: authHandler,
	}
}
