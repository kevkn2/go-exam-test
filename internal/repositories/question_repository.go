package repositories

import (
	"exam-test/internal/models"

	"gorm.io/gorm"
)

type QuestionRepo interface {
	CreateQuestion(question models.Question) (*models.Question, error)
	GetQuestion(id uint) (*models.Question, error)
	GetQuestionsByTestID(onlineTestID uint) ([]*models.Question, error)
	UpdateQuestion(
		onlineTestID uint,
		order int,
		question *models.Question,
	) error
}

type questionRepo struct {
	db *gorm.DB
}

// UpdateQuestion implements QuestionRepo.
func (q *questionRepo) UpdateQuestion(
	onlineTestID uint,
	order int,
	question *models.Question,
) error {
	if err := q.db.Where(
		`online_test_id = ? AND id = ?`,
		onlineTestID,
		order,
	).Delete(question).Error; err != nil {
		return err
	}

	if err := q.db.Create(question).Error; err != nil {
		return err
	}

	return nil
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
func (q *questionRepo) GetQuestionsByTestID(onlineTestID uint) ([]*models.Question, error) {
	var questions []*models.Question

	err := q.db.Where("online_test_id = ?", onlineTestID).Find(&questions).Error
	if err != nil {
		return nil, err
	}

	return questions, nil
}

func NewQuestionRepo(db *gorm.DB) QuestionRepo {
	return &questionRepo{
		db: db,
	}
}
