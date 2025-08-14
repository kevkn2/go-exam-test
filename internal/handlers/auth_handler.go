package handlers

import (
	"exam-test/internal/schemas"
	"exam-test/internal/services"
	"exam-test/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler interface {
	AuthenticateUser(ctx *gin.Context)
	RegisterAdmin(ctx *gin.Context)
	ValidAdmin(ctx *gin.Context)
}

type authHandler struct {
	userService services.UserService
	jwtUtils    utils.JWTUtils
}

// ValidAdmin implements AuthHandler.
func (a *authHandler) ValidAdmin(ctx *gin.Context) {
	userId, err := a.jwtUtils.TokenValid(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	user, err := a.userService.AuthorizeAdmin(userId)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "user not found"},
		)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// RegisterAdmin implements AuthHandler.
func (a *authHandler) RegisterAdmin(ctx *gin.Context) {
	var req schemas.RegisterRequestSchema

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Failed to hash password"},
		)
		return
	}

	user, err := a.userService.CreateAdmin(
		schemas.RegisterRequestSchema{
			Email:    req.Email,
			Password: hashedPassword,
		},
	)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		gin.H{
			"id":        user.ID,
			"email":     user.Email,
			"authority": user.Authority,
		},
	)
}

// AuthenticateUser implements AuthHandler.
func (a *authHandler) AuthenticateUser(ctx *gin.Context) {
	var req schemas.LoginRequestSchema

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.userService.GetUser(req.Email)
	if err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Invalid Credentials"},
		)
		return
	}

	var cp = utils.ComparePasswords{
		HashedPassword: user.Password,
		Password:       req.Password,
	}
	// check password
	if utils.ComparePasswordFail(cp) {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": "Invalid Credentials"},
		)
		return
	}

	token, err := a.jwtUtils.GenerateToken(user.ID)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "Invalid Credentials"},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"token": token,
			"type":  "Bearer",
		},
	)
}

func NewAuthHandler(
	userService services.UserService,
	jwtUtils utils.JWTUtils,
) AuthHandler {
	return &authHandler{
		userService: userService,
		jwtUtils:    jwtUtils,
	}
}
