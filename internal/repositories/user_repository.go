package repositories

import (
	"exam-test/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetUser(email string) (*models.User, error)
	GetAdmin(id uint) (*models.User, error)
	GetStudent(id uint) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

// CreateUser implements UserRepository.
func (u *userRepository) CreateUser(user *models.User) (*models.User, error) {
	result := u.db.Create(&user)
	return user, result.Error
}

// GetAdmin implements UserRepository.
func (u *userRepository) GetAdmin(id uint) (*models.User, error) {
	var admin models.User
	err := u.db.Where("authority = ?", "admin").Find(&admin, id).Error

	return &admin, err
}

// GetStudent implements UserRepository.
func (u *userRepository) GetStudent(id uint) (*models.User, error) {
	var student models.User
	err := u.db.Where("authority = ?", "student").Find(&student, id).Error

	return &student, err
}

// GetUser implements UserRepository.
func (u *userRepository) GetUser(email string) (*models.User, error) {
	var user models.User
	err := u.db.Where("email = ?", email).Find(&user).Error
	return &user, err
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}
