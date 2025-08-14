package services

import (
	"errors"
	"exam-test/internal/models"
	"exam-test/internal/repositories"
	"exam-test/internal/schemas"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserService interface {
	GetUser(email string) (*models.User, error)
	CreateAdmin(registerSchema schemas.RegisterRequestSchema) (*models.User, error)
	AuthorizeAdmin(id uint) (*models.User, error)
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

	resultUser, err := u.userRepo.CreateUser(&user)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return nil, fmt.Errorf("admin already exists")
		}
	}

	return resultUser, nil
}

// AuthenticateUser implements UserService.
func (u *userService) GetUser(email string) (*models.User, error) {
	return u.userRepo.GetUser(email)

}

// AuthorizeAdmin implements UserService.
func (u *userService) AuthorizeAdmin(id uint) (*models.User, error) {
	return u.userRepo.GetAdmin(id)
}

// AuthorizeStudent implements UserService.
func (u *userService) AuthorizeStudent(id uint) models.User {
	panic("unimplemented")
}

func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}
