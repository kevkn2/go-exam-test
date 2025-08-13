package services

import (
	"exam-test/internal/models"
	"exam-test/internal/repositories"
	"exam-test/internal/schemas"
)

type UserService interface {
	GetUser(email string) (*models.User, error)
	CreateAdmin(registerSchema schemas.RegisterRequestSchema) (*models.User, error)
	AuthorizeAdmin(id uint) models.User
	AuthorizeStudent(id uint) models.User
}

type userService struct {
	userRepo repositories.UserRepository
}

// CreateAdmin implements UserService.
func (u *userService) CreateAdmin(reg schemas.RegisterRequestSchema) (*models.User, error) {
	user := models.User{
		Email:     reg.Email,
		Password:  reg.Password,
		Authority: "admin",
	}

	return u.userRepo.CreateUser(&user)
}

// AuthenticateUser implements UserService.
func (u *userService) GetUser(email string) (*models.User, error) {
	return u.userRepo.GetUser(email)
}

// AuthorizeAdmin implements UserService.
func (u *userService) AuthorizeAdmin(id uint) models.User {
	panic("unimplemented")
}

// AuthorizeStudent implements UserService.
func (u *userService) AuthorizeStudent(id uint) models.User {
	panic("unimplemented")
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}
