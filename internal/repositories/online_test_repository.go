package repositories

import (
	"exam-test/internal/models"
	"exam-test/internal/schemas"

	"gorm.io/gorm"
)

type OnlineTestRepo interface {
	CreateTest(onlineTest models.OnlineTest) (*models.OnlineTest, error)
	GetTest(testID string) (*models.OnlineTest, error)
	GetTests() ([]models.OnlineTest, error)
	UpdateTestData(
		testID string,
		onlineTestData schemas.UpdateOnlineTestSchema,
	) (*models.OnlineTest, error)
}

type onlineTestRepo struct {
	db *gorm.DB
}

// UpdateTestData implements OnlineTestRepo.
func (o *onlineTestRepo) UpdateTestData(
	testID string,
	onlineTestData schemas.UpdateOnlineTestSchema,
) (*models.OnlineTest, error) {
	var onlineTest models.OnlineTest
	err := o.db.First(&onlineTest, "test_id = ?", testID).Error
	if err != nil {
		return nil, err
	}

	onlineTest.Title = *onlineTestData.Title
	onlineTest.Duration = *onlineTestData.Duration

	err = o.db.Save(&onlineTest).Error
	if err != nil {
		return nil, err
	}

	return &onlineTest, nil
}

// GetTests implements OnlineTestRepo.
func (o *onlineTestRepo) GetTests() ([]models.OnlineTest, error) {
	var onlineTests []models.OnlineTest
	err := o.db.Find(&onlineTests).Error
	if err != nil {
		return nil, err
	}
	return onlineTests, nil
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
	var onlineTest models.OnlineTest
	err := o.db.First(&onlineTest, "test_id = ?", testID).Error
	if err != nil {
		return nil, err
	}
	return &onlineTest, nil
}

func NewOnlineTestRepo(db *gorm.DB) OnlineTestRepo {
	return &onlineTestRepo{db: db}
}
