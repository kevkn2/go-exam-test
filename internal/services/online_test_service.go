package services

import (
	"encoding/json"
	"exam-test/internal/models"
	"exam-test/internal/repositories"

	"gorm.io/datatypes"
)

type OnlineTestService interface {
	CreateTest(onlineTest models.OnlineTest) (*models.OnlineTest, error)
	CreateMCQQuestion(
		text string,
		options []string,
		answer string,
		order int,
	) (models.Question, error)
	CreateTOFQuestion(
		text string,
		answer string,
		statements []string,
		order int,
	) (models.Question, error)
	CreateEssayQuestion(
		text string,
		answer string,
		order int,
	) (models.Question, error)
}

type onlineTestService struct {
	onlineTestRepo repositories.OnlineTestRepo
	questionRepo   repositories.QuestionRepo
}

// CreateEssayQuestion implements OnlineTestService.
func (o *onlineTestService) CreateEssayQuestion(
	text string,
	answer string,
	order int,
) (models.Question, error) {
	return models.Question{
		Type:   "essay",
		Text:   text,
		Answer: answer,
		Order:  order,
		Meta:   nil,
	}, nil
}

// CreateTOFQuestion implements OnlineTestService.
func (o *onlineTestService) CreateTOFQuestion(
	text string,
	answer string,
	statements []string,
	order int,
) (models.Question, error) {
	tofMeta := models.TOFMeta{
		Statements: statements,
	}

	tofMetaJSON, err := json.Marshal(tofMeta)
	if err != nil {
		return models.Question{}, err
	}

	return models.Question{
		Type:   "tof",
		Text:   text,
		Answer: answer,
		Order:  order,
		Meta:   datatypes.JSON(tofMetaJSON),
	}, nil

}

// CreateMCQQuestion implements OnlineTestService.
func (o *onlineTestService) CreateMCQQuestion(
	text string,
	options []string,
	answer string,
	order int,
) (models.Question, error) {
	mcqMeta := models.MCQMeta{
		Options: options,
	}

	mcqMetaJSON, err := json.Marshal(mcqMeta)
	if err != nil {
		return models.Question{}, err
	}

	return models.Question{
		Type:   "mcq",
		Text:   text,
		Answer: answer,
		Meta:   datatypes.JSON(mcqMetaJSON),
		Order:  order,
	}, nil
}

// CreateTest implements OnlineTestService.
func (o *onlineTestService) CreateTest(
	onlineTest models.OnlineTest,
) (*models.OnlineTest, error) {
	return o.onlineTestRepo.CreateTest(onlineTest)
}

func NewOnlineTestService(
	onlineTestRepo repositories.OnlineTestRepo,
	questionRepo repositories.QuestionRepo,
) OnlineTestService {
	return &onlineTestService{
		onlineTestRepo: onlineTestRepo,
		questionRepo:   questionRepo,
	}
}
