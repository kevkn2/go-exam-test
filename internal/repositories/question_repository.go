package repositories

import (
	"exam-test/internal/models"

	"gorm.io/gorm"
)

type QuestionRepo interface {
	CreateQuestion(question models.Question) (*models.Question, error)
	GetQuestion(id uint) (*models.Question, error)
	GetQuestionsByTestID(onlineTestID uint) ([]models.Question, error)
}

type questionRepo struct {
	db *gorm.DB
}

// CreateQuestion implements QuestionRepo.
func (q *questionRepo) CreateQuestion(question models.Question) (*models.Question, error) {
	err := q.db.Create(&question).Error
	if err != nil {
		return nil, err
	}

	return &question, nil
}

// GetQuestion implements QuestionRepo.
func (q *questionRepo) GetQuestion(id uint) (*models.Question, error) {
	panic("unimplemented")
}

// GetQuestionsByTestID implements QuestionRepo.
func (q *questionRepo) GetQuestionsByTestID(onlineTestID uint) ([]models.Question, error) {
	panic("unimplemented")
}

func NewQuestionRepo(db *gorm.DB) QuestionRepo {
	return &questionRepo{
		db: db,
	}
}
