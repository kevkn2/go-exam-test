package repositories

import (
	"exam-test/internal/models"

	"gorm.io/gorm"
)

type OnlineTestRepo interface {
	CreateTest(onlineTest models.OnlineTest) (*models.OnlineTest, error)
	GetTest(testID string) (*models.OnlineTest, error)
	GetTests() ([]models.OnlineTest, error)
}

type onlineTestRepo struct {
	db *gorm.DB
}

// GetTests implements OnlineTestRepo.
func (o *onlineTestRepo) GetTests() ([]models.OnlineTest, error) {
	panic("unimplemented")
}

// CreateTest implements OnlineTestRepo.
func (o *onlineTestRepo) CreateTest(onlineTest models.OnlineTest) (*models.OnlineTest, error) {
	err := o.db.Create(&onlineTest).Error

	if err != nil {
		return nil, err
	}

	return &onlineTest, nil
}

// GetTest implements OnlineTestRepo.
func (o *onlineTestRepo) GetTest(testID string) (*models.OnlineTest, error) {
	panic("unimplemented")
}

func NewOnlineTestRepo(db *gorm.DB) OnlineTestRepo {
	return &onlineTestRepo{db: db}
}
