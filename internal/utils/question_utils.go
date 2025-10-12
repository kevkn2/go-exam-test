package utils

import (
	"encoding/json"
	"exam-test/internal/models"
	"exam-test/internal/schemas"

	"gorm.io/datatypes"
)

type questionUtils struct{}

func CreateMCQQuestion(
	text string,
	options []string,
	answer string,
	order int,
	onlineTestID uint,
) (*models.Question, error) {
	mcqMeta := models.MCQMeta{
		Options: options,
	}

	mcqMetaJSON, err := json.Marshal(mcqMeta)
	if err != nil {
		return nil, err
	}

	return &models.Question{
		Type:         "mcq",
		Text:         text,
		Answer:       answer,
		Meta:         datatypes.JSON(mcqMetaJSON),
		Order:        order,
		OnlineTestID: onlineTestID,
	}, nil
}

func CreateTOFQuestion(
	text string,
	answer string,
	statements []string,
	order int,
	onlineTestID uint,
) (*models.Question, error) {
	tofMeta := models.TOFMeta{
		Statements: statements,
	}

	tofMetaJSON, err := json.Marshal(tofMeta)
	if err != nil {
		return nil, err
	}

	return &models.Question{
		Type:   "tof",
		Text:   text,
		Answer: answer,
		Order:  order,
		Meta:   datatypes.JSON(tofMetaJSON),
	}, nil
}

// CreateQuestionModel implements QuestionUtils.
func (q *questionUtils) CreateQuestionModel(
	question *schemas.QuestionSchema,
	onlineTestID uint,
) (*models.Question, error) {
	var qModel *models.Question
	var err error

	switch question.Type {
	case "mcq":
		qModel, err = CreateMCQQuestion(
			question.Text,
			question.Options,
			question.Answer,
			question.Order,
			onlineTestID,
		)
	case "tof":
		qModel, err = CreateTOFQuestion(
			question.Text,
			question.Answer,
			question.Statements,
			question.Order,
			onlineTestID,
		)
	}
	if err != nil {
		return nil, err
	}

	return qModel, nil
}

func CreateMCQQuestionSchema(
	text string,
	answer string,
	metadataOptions datatypes.JSON,
) (*schemas.QuestionSchema, error) {
	var metadata models.MCQMeta
	err := json.Unmarshal([]byte(metadataOptions), &metadata)
	if err != nil {
		return nil, err
	}

	return &schemas.QuestionSchema{
		Type:       "mcq",
		Text:       text,
		Answer:     answer,
		Options:    metadata.Options,
		Statements: []string{},
	}, nil
}

func CreateTOFQuestionSchema(
	text string,
	answer string,
	metadataStatements datatypes.JSON,
) (*schemas.QuestionSchema, error) {
	var metadata models.TOFMeta
	err := json.Unmarshal([]byte(metadataStatements), &metadata)
	if err != nil {
		return nil, err
	}
	return &schemas.QuestionSchema{
		Type:       "tof",
		Text:       text,
		Answer:     answer,
		Options:    []string{},
		Statements: metadata.Statements,
	}, nil
}

// GenerateQuestionsSchema implements QuestionUtils.
func (q *questionUtils) GenerateQuestionSchema(
	question *models.Question,
) (*schemas.QuestionSchema, error) {
	var questionSchema *schemas.QuestionSchema
	var err error

	switch question.Type {
	case "mcq":
		questionSchema, err = CreateMCQQuestionSchema(
			question.Text,
			question.Answer,
			question.Meta,
		)

	case "tof":
		questionSchema, err = CreateTOFQuestionSchema(
			question.Text,
			question.Answer,
			question.Meta,
		)
	}
	if err != nil {
		return nil, err
	}

	return questionSchema, nil
}

type QuestionUtils interface {
	CreateQuestionModel(
		question *schemas.QuestionSchema,
		onlineTestID uint,
	) (*models.Question, error)
	GenerateQuestionSchema(
		question *models.Question,
	) (*schemas.QuestionSchema, error)
}

func NewQuestionUtils() QuestionUtils {
	return &questionUtils{}
}
