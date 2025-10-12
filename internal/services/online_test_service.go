package services

import (
	"exam-test/internal/models"
	"exam-test/internal/repositories"
	"exam-test/internal/schemas"
	"exam-test/internal/utils"
	"strconv"
)

type OnlineTestService interface {
	CreateTest(
		title string,
		testID string,
		duration int,
		questions []*schemas.QuestionSchema,
	) (*models.OnlineTest, error)
	GetTest(testID string) (*schemas.OnlineTestSchema, error)
	GetAllTest() ([]*models.OnlineTest, error)
	UpdateTestData(
		testID string,
		onlineTestData *schemas.UpdateOnlineTestSchema,
	) (*models.OnlineTest, error)
	UpdateQuestions(
		testID string,
		questions *schemas.QuestionsSchema,
	) (*schemas.OnlineTestSchema, error)
}

type onlineTestService struct {
	onlineTestRepo repositories.OnlineTestRepo
	questionRepo   repositories.QuestionRepo
	questionUtils  utils.QuestionUtils
}

// UpdateQuestions implements OnlineTestService.
func (o *onlineTestService) UpdateQuestions(
	testID string,
	questions *schemas.QuestionsSchema,
) (*schemas.OnlineTestSchema, error) {
	onlineTest, err := o.onlineTestRepo.GetTest(testID)
	if err != nil {
		return nil, err
	}

	for _, question := range questions.Questions {
		if question.Order < 1 {
			continue
		}

		qModel, err := o.questionUtils.CreateQuestionModel(
			question,
			onlineTest.ID,
		)
		if err != nil {
			return nil, err
		}

		if err := o.questionRepo.UpdateQuestion(
			onlineTest.ID,
			question.Order,
			qModel,
		); err != nil {
			return nil, err
		}
	}

	questionModels, err := o.questionRepo.GetQuestionsByTestID(onlineTest.ID)
	if err != nil {
		return nil, err
	}

	var updatedQuestions []*schemas.QuestionSchema

	for _, questionModel := range questionModels {
		questionSchema, err := o.questionUtils.GenerateQuestionSchema(questionModel)
		if err != nil {
			return nil, err
		}

		updatedQuestions = append(updatedQuestions, questionSchema)
	}

	onlineTestFinal := schemas.OnlineTestSchema{
		Title:     onlineTest.Title,
		TestID:    onlineTest.TestID,
		Duration:  onlineTest.Duration,
		Questions: updatedQuestions,
	}

	return &onlineTestFinal, nil
}

// UpdateTestData implements OnlineTestService.
func (o *onlineTestService) UpdateTestData(
	testID string,
	onlineTestData *schemas.UpdateOnlineTestSchema,
) (*models.OnlineTest, error) {
	onlineTest, err := o.onlineTestRepo.UpdateTestData(testID, onlineTestData)
	if err != nil {
		return nil, err
	}

	return onlineTest, nil
}

// GetAllTest implements OnlineTestService.
func (o *onlineTestService) GetAllTest() ([]*models.OnlineTest, error) {
	onlineTests, err := o.onlineTestRepo.GetTests()
	if err != nil {
		return nil, err
	}

	return onlineTests, nil
}

// GetTest implements OnlineTestService.
func (o *onlineTestService) GetTest(testID string) (*schemas.OnlineTestSchema, error) {
	onlineTest, err := o.onlineTestRepo.GetTest(testID)
	if err != nil {
		return nil, err
	}

	questionModels, err := o.questionRepo.GetQuestionsByTestID(onlineTest.ID)
	if err != nil {
		return nil, err
	}

	var questions []*schemas.QuestionSchema
	for _, questionModel := range questionModels {
		question, err := o.questionUtils.GenerateQuestionSchema(questionModel)
		if err != nil {
			return nil, err
		}

		questions = append(questions, question)
	}

	var onlineTestWithQuestions schemas.OnlineTestSchema
	onlineTestWithQuestions.Title = onlineTest.Title
	onlineTestWithQuestions.TestID = onlineTest.TestID
	onlineTestWithQuestions.Duration = onlineTest.Duration
	onlineTestWithQuestions.Questions = questions

	return &onlineTestWithQuestions, nil
}

// CreateTest implements OnlineTestService.
func (o *onlineTestService) CreateTest(
	title string,
	testID string,
	duration int,
	questions []*schemas.QuestionSchema,
) (*models.OnlineTest, error) {
	var questionModels []*models.Question

	for _, question := range questions {
		testIDUint, err := strconv.ParseUint(testID, 10, 64)
		if err != nil {
			return nil, err
		}

		questionModel, err := o.questionUtils.CreateQuestionModel(
			question,
			uint(testIDUint),
		)
		if err != nil {
			return nil, err
		}

		questionModels = append(questionModels, questionModel)
	}

	onlineTest := &models.OnlineTest{
		Title:     title,
		TestID:    testID,
		Duration:  duration,
		Questions: questionModels,
	}

	return o.onlineTestRepo.CreateTest(onlineTest)
}

func NewOnlineTestService(
	onlineTestRepo repositories.OnlineTestRepo,
	questionRepo repositories.QuestionRepo,
	questionUtils utils.QuestionUtils,
) OnlineTestService {
	return &onlineTestService{
		onlineTestRepo: onlineTestRepo,
		questionRepo:   questionRepo,
		questionUtils:  questionUtils,
	}
}
