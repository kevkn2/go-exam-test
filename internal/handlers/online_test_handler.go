package handlers

import (
	"exam-test/internal/models"
	"exam-test/internal/schemas"
	"exam-test/internal/services"
	"exam-test/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OnlineTestHandler interface {
	CreateTest(ctx *gin.Context)
	GetTest(ctx *gin.Context)
	GetAllTest(ctx *gin.Context)
	UpdateTestData(ctx *gin.Context)
	UpdateQuestions(ctx *gin.Context)
}

type onlineTestHandler struct {
	onlineTestService services.OnlineTestService
	jwtUtils          utils.JWTUtils
}

// UpdateQuestions implements OnlineTestHandler.
func (o *onlineTestHandler) UpdateQuestions(ctx *gin.Context) {
	_, err := o.jwtUtils.TokenValid(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": err.Error()},
		)
		return
	}

	testID := ctx.Param("testID")
	if testID == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "testID parameter is required"},
		)
		return
	}

	var questions schemas.QuestionsSchema
	if err := ctx.ShouldBindJSON(&questions); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	onlineTestUpdated, err := o.onlineTestService.UpdateQuestions(testID, questions)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Failed to update questions: %v", err)},
		)
		return
	}

	ctx.JSON(http.StatusOK, onlineTestUpdated)
}

// UpdateTestData implements OnlineTestHandler.
func (o *onlineTestHandler) UpdateTestData(ctx *gin.Context) {
	_, err := o.jwtUtils.TokenValid(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": err.Error()},
		)
		return
	}

	testID := ctx.Param("testID")
	if testID == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "testID parameter is required"},
		)
		return
	}

	var onlineTestData schemas.UpdateOnlineTestSchema
	if err := ctx.ShouldBindJSON(&onlineTestData); err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": err.Error()},
		)
		return
	}

	onlineTest, err := o.onlineTestService.UpdateTestData(testID, onlineTestData)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Failed to update test data: %v", err)},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		schemas.OnlineTestSummarySchema{
			ID:       onlineTest.ID,
			Title:    onlineTest.Title,
			TestID:   onlineTest.TestID,
			Duration: onlineTest.Duration,
		},
	)
}

// GetAllTest implements OnlineTestHandler.
func (o *onlineTestHandler) GetAllTest(ctx *gin.Context) {
	_, err := o.jwtUtils.TokenValid(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": err.Error()},
		)
		return
	}

	onlineTests, err := o.onlineTestService.GetAllTest()
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Failed to get tests: %v", err)},
		)
		return
	}

	var onlineTestSummaries schemas.AllOnlineTestsSchema
	for _, test := range onlineTests {
		summary := schemas.OnlineTestSummarySchema{
			ID:       test.ID,
			Title:    test.Title,
			TestID:   test.TestID,
			Duration: test.Duration,
		}

		onlineTestSummaries.Tests = append(onlineTestSummaries.Tests, summary)
	}

	ctx.JSON(http.StatusOK, onlineTestSummaries)
}

// GetTest implements OnlineTestHandler.
func (o *onlineTestHandler) GetTest(ctx *gin.Context) {
	_, err := o.jwtUtils.TokenValid(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": err.Error()},
		)
		return
	}

	testID := ctx.Param("testID")
	if testID == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "testID parameter is required"},
		)
		return
	}

	onlineTest, err := o.onlineTestService.GetTest(testID)
	if err != nil {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": fmt.Sprintf("Failed to get test: %v", err)},
		)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": onlineTest})
}

// CreateTest implements OnlineTestHandler.
func (o *onlineTestHandler) CreateTest(ctx *gin.Context) {
	_, err := o.jwtUtils.TokenValid(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusUnauthorized,
			gin.H{"error": err.Error()},
		)
		return
	}

	var req schemas.OnlineTestSchema

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var questions []models.Question

	for questionIndex, question := range req.Questions {
		var questionModel models.Question
		switch question.Type {
		case "mcq":
			questionModel, err = o.onlineTestService.CreateMCQQuestion(
				question.Text,
				question.Options,
				question.Answer,
				questionIndex+1,
			)
			if err != nil {
				ctx.JSON(
					http.StatusInternalServerError,
					gin.H{"error": fmt.Sprintf("Failed to create question %d: %v", questionIndex+1, err)},
				)
				return
			}
		case "tof":
			questionModel, err = o.onlineTestService.CreateTOFQuestion(
				question.Text,
				question.Answer,
				question.Statements,
				questionIndex+1,
			)
			if err != nil {
				ctx.JSON(
					http.StatusInternalServerError,
					gin.H{"error": fmt.Sprintf("Failed to create question %d: %v", questionIndex+1, err)},
				)
				return
			}
		default:
			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": fmt.Sprintf("Invalid question type for question %d", questionIndex+1)},
			)
			return
		}
		questions = append(questions, questionModel)
	}

	onlineTest, err := o.onlineTestService.CreateTest(
		models.OnlineTest{
			Title:     req.Title,
			TestID:    req.TestID,
			Duration:  req.Duration,
			Questions: questions,
		},
	)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"error": fmt.Sprintf("Failed to create test: %v", err)},
		)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": onlineTest})
}

func NewOnlineTestHandler(
	onlineTestService services.OnlineTestService,
	jwtUtils utils.JWTUtils,
) OnlineTestHandler {
	return &onlineTestHandler{
		onlineTestService: onlineTestService,
		jwtUtils:          jwtUtils,
	}
}
