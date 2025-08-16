package repositories

import (
	"exam-test/internal/models"

	"gorm.io/gorm"
)

type StudentRepository interface {
	WithTx(tx *gorm.DB) StudentRepository
	CreateStudent(student models.Student) (*models.Student, error)
	GetStudent(id uint) (*models.Student, error)
}

type studentRepository struct {
	db *gorm.DB
}

// WithTx implements StudentRepository.
func (s *studentRepository) WithTx(tx *gorm.DB) StudentRepository {
	return &studentRepository{db: tx}
}

// CreateStudent implements StudentRepository.
func (s *studentRepository) CreateStudent(student models.Student) (*models.Student, error) {
	result := s.db.Create(&student)
	return &student, result.Error
}

// GetStudent implements StudentRepository.
func (s *studentRepository) GetStudent(id uint) (*models.Student, error) {
	var student models.Student
	err := s.db.Where("user_id", id).Find(&student).Error
	return &student, err
}

func NewStudentRepository(db *gorm.DB) StudentRepository {
	return &studentRepository{db: db}
}
