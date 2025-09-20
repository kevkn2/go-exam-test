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
}

type onlineTestHandler struct {
	onlineTestService services.OnlineTestService
	jwtUtils          utils.JWTUtils
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

	var req schemas.OnlineTestRequestSchema

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
