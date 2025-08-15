package services

import (
	"errors"
	"exam-test/internal/models"
	"exam-test/internal/repositories"
	"exam-test/internal/schemas"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type AuthService interface {
	GetUser(email string) (*models.User, error)
	CreateAdmin(registerSchema schemas.RegisterRequestSchema) (*models.User, error)
	CreateStudent(registerSchema schemas.RegisterStudentRequestSchema) (struct {
		*models.User
		*models.Student
	}, error)
	AuthorizeAdmin(id uint) (*models.User, error)
	AuthorizeStudent(id uint) (*models.User, error)
}

type authService struct {
	userRepo    repositories.UserRepository
	studentRepo repositories.StudentRepository
}

// CreateStudent implements UserService.
func (u *authService) CreateStudent(reg schemas.RegisterStudentRequestSchema) (struct {
	*models.User
	*models.Student
}, error) {
	user := models.User{
		Email:     reg.Email,
		Password:  reg.Password,
		Authority: "student",
	}

	resultUser, err := u.userRepo.CreateUser(&user)

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return struct {
				*models.User
				*models.Student
			}{}, fmt.Errorf("student already exists")
		}
	}

	student := models.Student{
		Name:   reg.Name,
		School: reg.School,
		UserID: resultUser.ID,
	}

	resultStudent, err := u.studentRepo.CreateStudent(student)
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return struct {
				*models.User
				*models.Student
			}{}, fmt.Errorf("student already exists")
		}
	}

	return struct {
		*models.User
		*models.Student
	}{resultUser, resultStudent}, nil
}

// CreateAdmin implements UserService.
func (u *authService) CreateAdmin(reg schemas.RegisterRequestSchema) (*models.User, error) {
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
func (u *authService) GetUser(email string) (*models.User, error) {
	return u.userRepo.GetUser(email)

}

// AuthorizeAdmin implements UserService.
func (u *authService) AuthorizeAdmin(id uint) (*models.User, error) {
	return u.userRepo.GetAdmin(id)
}

// AuthorizeStudent implements UserService.
func (u *authService) AuthorizeStudent(id uint) (*models.User, error) {
	return u.userRepo.GetStudent(id)
}

func NewAuthService(
	userRepo repositories.UserRepository,
	studentRepo repositories.StudentRepository,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		studentRepo: studentRepo,
	}
}
