package routes

import (
	"exam-test/internal/handlers"

	"github.com/gin-gonic/gin"
)

type OnlineTestRoute interface {
	Routes(route *gin.Engine)
}

type onlineTestRoute struct {
	onlineTestHandler handlers.OnlineTestHandler
}

// Routes implements OnlineTestRoute.
func (o *onlineTestRoute) Routes(route *gin.Engine) {
	onlineTest := route.Group("/api/v1/onlineTest")

	onlineTest.POST("/createTest", o.onlineTestHandler.CreateTest)
	onlineTest.GET("/getTest/:testID", o.onlineTestHandler.GetTest)
	onlineTest.GET("/getAllTests", o.onlineTestHandler.GetAllTest)
	onlineTest.PATCH("/updateTestData/:testID", o.onlineTestHandler.UpdateTestData)
}

func NewOnlineTestRoute(onlineTestHandler handlers.OnlineTestHandler) OnlineTestRoute {
	return &onlineTestRoute{
		onlineTestHandler: onlineTestHandler,
	}
}
