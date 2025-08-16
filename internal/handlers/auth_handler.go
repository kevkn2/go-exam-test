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
	RegisterStudent(ctx *gin.Context)
	ValidAdmin(ctx *gin.Context)
	ValidStudent(ctx *gin.Context)
}

type authHandler struct {
	authService services.AuthService
	jwtUtils    utils.JWTUtils
}

// RegisterStudent implements AuthHandler.
func (a *authHandler) RegisterStudent(ctx *gin.Context) {
	var req schemas.RegisterStudentRequestSchema

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

	user, err := a.authService.CreateStudent(
		schemas.RegisterStudentRequestSchema{
			Email:    req.Email,
			Password: hashedPassword,
			Name:     req.Name,
			School:   req.School,
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
			"userId":    user.User.ID,
			"email":     user.User.Email,
			"authority": user.User.Authority,
			"name":      user.Student.Name,
			"school":    user.Student.School,
		},
	)
}

// ValidStudent implements AuthHandler.
func (a *authHandler) ValidStudent(ctx *gin.Context) {
	userId, err := a.jwtUtils.TokenValid(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	user, err := a.authService.AuthorizeStudent(userId)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "user not found"},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{
			"userId":    user.User.ID,
			"email":     user.User.Email,
			"authority": user.User.Authority,
			"name":      user.Student.Name,
			"school":    user.Student.School,
		},
	)
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

	user, err := a.authService.AuthorizeAdmin(userId)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "user not found"},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		schemas.UserInfoSchema{
			ID:        user.ID,
			Email:     user.Email,
			Authority: user.Authority,
		},
	)
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

	user, err := a.authService.CreateAdmin(
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
		schemas.UserInfoSchema{
			ID:        user.ID,
			Email:     user.Email,
			Authority: user.Authority,
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

	user, err := a.authService.GetUser(req.Email)
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
		schemas.AuthResponseSchema{
			Token: token,
			Type:  "Bearer",
		},
	)
}

func NewAuthHandler(
	authService services.AuthService,
	jwtUtils utils.JWTUtils,
) AuthHandler {
	return &authHandler{
		authService: authService,
		jwtUtils:    jwtUtils,
	}
}
