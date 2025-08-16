package services

import (
	"errors"
	"exam-test/internal/models"
	"exam-test/internal/repositories"
	"exam-test/internal/schemas"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type AuthService interface {
	GetUser(email string) (*models.User, error)
	CreateAdmin(registerSchema schemas.RegisterRequestSchema) (*models.User, error)
	CreateStudent(registerSchema schemas.RegisterStudentRequestSchema) (struct {
		*models.User
		*models.Student
	}, error)
	AuthorizeAdmin(id uint) (*models.User, error)
	AuthorizeStudent(id uint) (struct {
		*models.User
		*models.Student
	}, error)
}

type authService struct {
	userRepo    repositories.UserRepository
	studentRepo repositories.StudentRepository
	db          *gorm.DB
}

// CreateStudent implements UserService.
func (u *authService) CreateStudent(reg schemas.RegisterStudentRequestSchema) (struct {
	*models.User
	*models.Student
}, error) {
	var result struct {
		*models.User
		*models.Student
	}

	err := u.db.Transaction(func(tx *gorm.DB) error {
		user := models.User{
			Email:     reg.Email,
			Password:  reg.Password,
			Authority: "student",
		}

		resultUser, err := u.userRepo.WithTx(tx).CreateUser(&user)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return fmt.Errorf("student already exists")
			}
			return err
		}

		student := models.Student{
			Name:   reg.Name,
			School: reg.School,
			UserID: resultUser.ID,
		}

		resultStudent, err := u.studentRepo.WithTx(tx).CreateStudent(student)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == "23505" {
				return fmt.Errorf("student already exists")
			}
			return err
		}

		// store results for returning after transaction
		result.User = resultUser
		result.Student = resultStudent
		return nil
	})

	if err != nil {
		return struct {
			*models.User
			*models.Student
		}{}, err
	}

	return result, nil
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
func (u *authService) AuthorizeStudent(id uint) (struct {
	*models.User
	*models.Student
}, error) {
	user, err := u.userRepo.GetStudent(id)
	if err != nil {
		return struct {
			*models.User
			*models.Student
		}{}, err
	}

	student, err := u.studentRepo.GetStudent(id)
	if err != nil {
		return struct {
			*models.User
			*models.Student
		}{}, err
	}

	return struct {
		*models.User
		*models.Student
	}{
		User:    user,
		Student: student,
	}, nil
}

func NewAuthService(
	userRepo repositories.UserRepository,
	studentRepo repositories.StudentRepository,
	db *gorm.DB,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		studentRepo: studentRepo,
		db:          db,
	}
}
